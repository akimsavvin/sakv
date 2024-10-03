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

	mx   sync.RWMutex
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

	e.mx.RLock()
	val, ok := e.data[key]
	e.mx.RUnlock()

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
	e.data[key] = value
	e.mx.Unlock()

	log.InfoContext(ctx, "value set")

	return nil
}

func (e *Engine) DEL(ctx context.Context, key string) error {
	log := e.log.With(slog.String("key", key))
	log.DebugContext(ctx, "deleting a value")

	e.mx.Lock()
	delete(e.data, key)
	e.mx.Unlock()

	log.InfoContext(ctx, "value deleted")

	return nil
}
