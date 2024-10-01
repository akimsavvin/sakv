package listener

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"sakv/internal/database/config"
	"sakv/pkg/memory"
	"sakv/pkg/semaphore"
	"sakv/pkg/sl"
	"time"
)

type QueryHandler interface {
	HandleQuery(ctx context.Context, query string) string
}

type Listener struct {
	log *slog.Logger
	qh  QueryHandler

	addr        string
	idleTimeout time.Duration
	maxMsgSize  uint

	sm *semaphore.Semaphore
}

var (
	ErrIdleTimeout = errors.New("idle timeout")
)

func New(log *slog.Logger, cfg config.Network, qh QueryHandler) (*Listener, error) {
	idleTimeout, err := time.ParseDuration(cfg.IdleTimeout)
	if err != nil {
		return nil, err
	}

	maxMsgSize, err := memory.ParseSize(cfg.MaxMsgSize)
	if err != nil {
		return nil, err
	}

	return &Listener{
		log:         log.With(sl.Comp("listener.Listener")),
		qh:          qh,
		addr:        cfg.Addr,
		idleTimeout: idleTimeout,
		maxMsgSize:  maxMsgSize,
		sm:          semaphore.New(cfg.MaxConns),
	}, nil
}

func (l *Listener) StartListening(ctx context.Context) error {
	l.log.DebugContext(ctx, "starting listening", slog.Any("addr", l.addr))

	netl, err := net.Listen("tcp", l.addr)
	if err != nil {
		return err
	}

	go func() error {
		<-ctx.Done()
		return netl.Close()
	}()

	for {
		conn, err := netl.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				l.log.InfoContext(ctx, "listening stopped")
				return ctx.Err()
			}

			l.log.ErrorContext(ctx, "listening failed, an error occurred while accepting connection", sl.Err(err))
			continue
		}
		l.handleConn(ctx, conn)
	}
}

func (l *Listener) handleConn(ctx context.Context, conn net.Conn) {
	log := l.log.With(slog.String("conn_addr", conn.RemoteAddr().String()))
	log.InfoContext(ctx, "handling connection")

	l.sm.Acquire()
	go func() {
		defer l.sm.Release()
		l.listenConn(ctx, conn)
	}()

	log.InfoContext(ctx, "connection accepted")
}

func (l *Listener) listenConn(ctx context.Context, conn net.Conn) error {
	log := l.log.With(slog.String("conn_addr", conn.RemoteAddr().String()))

	log.DebugContext(ctx, "starting connection listening")

	defer func() {
		log.DebugContext(ctx, "closing connection")
		if err := conn.Close(); err != nil {
			log.WarnContext(ctx, "connection closed with error", sl.Err(err))
		} else {
			log.InfoContext(ctx, "connection closed")
		}
	}()

	t := time.NewTicker(l.idleTimeout)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			log.InfoContext(ctx, "context done")
			return ctx.Err()
		case <-t.C:
			log.InfoContext(ctx, "idle timeout")
			return ErrIdleTimeout
		default:
			t.Reset(l.idleTimeout)

			buf := make([]byte, l.maxMsgSize)
			n, err := conn.Read(buf)
			if err != nil {
				log.ErrorContext(ctx, "an error occurred while reading connection query", sl.Err(err))
				return err
			}
			buf = buf[:n]
			query := string(buf)

			log.DebugContext(ctx, "received query", slog.String("query", query))

			resp := l.qh.HandleQuery(ctx, query)
			_, err = conn.Write([]byte(resp))
			if err != nil {
				log.ErrorContext(ctx, "an error occurred while writing connection response", sl.Err(err))
				return err
			}

			log.InfoContext(ctx, "response sent", slog.String("response", resp))
		}
	}
}
