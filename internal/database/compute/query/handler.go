package query

import (
	"context"
	"log/slog"
	"sakv/internal/database/storage/engine"
	"sakv/pkg/sl"
)

type Handler struct {
	log *slog.Logger
	e   engine.Engine
}

func NewHandler(log *slog.Logger, e engine.Engine) *Handler {
	return &Handler{
		log: log.With(sl.Comp("query.Handler")),
		e:   e,
	}
}

func (h *Handler) HandleQuery(ctx context.Context, queryStr string) string {
	log := h.log.With(slog.String("string_query", queryStr))
	log.DebugContext(ctx, "handling query")

	q, err := ParseQueryStr(queryStr)
	if err != nil {
		log.WarnContext(ctx, "failed to parse query", sl.Err(err))
		return h.err(err)
	}

	log = log.With("query", slog.Any("query", *q))

	switch q.Command() {
	case CommandGET:
		log.InfoContext(ctx, "handling GET command")

		res, err := h.e.GET(ctx, q.Arg(0))
		if err != nil {
			return h.err(err)
		}
		return h.ok(res)
	case CommandSET:
		log.InfoContext(ctx, "handling SET command")

		err = h.e.SET(ctx, q.Arg(0), q.Arg(1))
		if err != nil {
			return h.err(err)
		}
		return h.ok()
	case CommandDEL:
		log.InfoContext(ctx, "handling DEL command")

		err = h.e.DEL(ctx, q.Arg(0))
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
