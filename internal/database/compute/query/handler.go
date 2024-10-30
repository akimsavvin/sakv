package query

import (
	"context"
	"github.com/akimsavvin/sakv/internal/database/storage/engine"
	"github.com/akimsavvin/sakv/internal/database/storage/wal"
	"github.com/akimsavvin/sakv/pkg/sl"
	"log/slog"
)

type Handler struct {
	log *slog.Logger
	wal *wal.WAL
	e   engine.Engine
}

func NewHandler(log *slog.Logger, wal *wal.WAL, e engine.Engine) *Handler {
	return &Handler{
		log: log.With(sl.Comp("query.Handler")),
		wal: wal,
		e:   e,
	}
}

func (h *Handler) HandleQuery(ctx context.Context, queryStr string) string {
	log := h.log.With(slog.String("string_query", queryStr))
	log.DebugContext(ctx, "handling query")

	query, err := ParseQueryStr(queryStr)
	if err != nil {
		log.WarnContext(ctx, "failed to parse query", sl.Err(err))
		return h.err(err)
	}

	log = h.log.With("query", slog.Group("query",
		slog.String("command", string(query.Command())),
		slog.Any("args", query.args),
	))

	if h.wal != nil && query.Command() != CommandGET {
		if err := h.wal.Write(ctx, queryStr); err != nil {
			log.ErrorContext(ctx, "failed to write query to WAL", sl.Err(err))
			return h.err(err)
		}
	}

	switch query.Command() {
	case CommandGET:
		log.InfoContext(ctx, "handling GET command")

		res, err := h.e.GET(ctx, query.Arg(0))
		if err != nil {
			return h.err(err)
		}

		return h.ok(res)
	case CommandSET:
		log.InfoContext(ctx, "handling SET command")

		err := h.e.SET(ctx, query.Arg(0), query.Arg(1))
		if err != nil {
			return h.err(err)
		}

		return h.ok()
	case CommandDEL:
		log.InfoContext(ctx, "handling DEL command")

		err := h.e.DEL(ctx, query.Arg(0))
		if err != nil {
			return h.err(err)
		}

		return h.ok()
	default:
		log.ErrorContext(ctx, "received query with unknown command")
		panic("received query with unknown command")
	}
}

func (h *Handler) err(err error) string {
	return "[ERROR]: " + err.Error()
}

func (h *Handler) ok(msg ...string) string {
	if len(msg) > 0 {
		return "[OK]: " + msg[0]
	}

	return "[OK]"
}
