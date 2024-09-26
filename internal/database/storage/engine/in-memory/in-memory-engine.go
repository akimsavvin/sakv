package inmemory

import (
	"context"
	"errors"
	"log/slog"
	"sakv/pkg/sl"
	"sync"
)

type Engine struct {
	log *slog.Logger

	mx   sync.Mutex
	data map[string]string
}

var (
	ErrNotFound = errors.New("no value set for key")
)

func NewEngine(log *slog.Logger) *Engine {
	return &Engine{
		log:  log.With(sl.Comp("inmemory.Engine")),
		data: make(map[string]string),
	}
}

func (e *Engine) GET(ctx context.Context, key string) (string, error) {
	log := e.log.With(slog.String("key", key))
	log.DebugContext(ctx, "getting a value")

	e.mx.Lock()
	defer e.mx.Unlock()

	val, ok := e.data[key]
	if !ok {
		log.WarnContext(ctx, "no value found for key")
		return "", ErrNotFound
	}

	log.InfoContext(ctx, "value retrieved")
	return val, nil
}

func (e *Engine) SET(ctx context.Context, key, value string) error {
	log := e.log.With(slog.String("key", key), slog.String("value", value))
	log.DebugContext(ctx, "setting a value")

	e.mx.Lock()
	defer e.mx.Unlock()

	e.data[key] = value
	log.InfoContext(ctx, "value set")

	return nil
}

func (e *Engine) DEL(ctx context.Context, key string) error {
	log := e.log.With(slog.String("key", key))
	log.DebugContext(ctx, "deleting a value")

	e.mx.Lock()
	defer e.mx.Unlock()

	delete(e.data, key)
	log.InfoContext(ctx, "value deleted")

	return nil
}
