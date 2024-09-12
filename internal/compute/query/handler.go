package query

import (
	"context"
	"sakv/internal/storage/engine"
)

type Handler struct {
	e engine.Engine
}

func NewHandler(e engine.Engine) *Handler {
	return &Handler{
		e: e,
	}
}

func (h *Handler) HandleQuery(ctx context.Context, queryStr string) string {
	q, err := ParseQueryStr(queryStr)
	if err != nil {
		return "[ERROR]: " + err.Error()
	}

	switch q.Command() {
	case CommandGET:
		res, err := h.e.GET(ctx, q.Arg(0))
		if err != nil {
			return "[ERROR]: " + err.Error()
		}
		return "[OK]: " + res
	case CommandSET:
		err = h.e.SET(ctx, q.Arg(0), q.Arg(1))
		if err != nil {
			return "[ERROR]: " + err.Error()
		}
		return "[OK]"
	case CommandDEL:
		err = h.e.DEL(ctx, q.Arg(0))
		if err != nil {
			return "[ERROR]: " + err.Error()
		}
		return "[OK]"
	default:
		panic("something went wrong")
	}
}
