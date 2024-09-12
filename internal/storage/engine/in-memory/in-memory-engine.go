package inmemory

import (
	"context"
	"errors"
)

type Engine struct {
	data map[string]string
}

var (
	ErrNotFound = errors.New("no value set for key")
)

func NewEngine() *Engine {
	return &Engine{
		data: make(map[string]string),
	}
}

func (e *Engine) GET(ctx context.Context, key string) (string, error) {
	_ = ctx

	val, ok := e.data[key]
	if !ok {
		return "", ErrNotFound
	}

	return val, nil
}

func (e *Engine) SET(ctx context.Context, key, value string) error {
	_ = ctx
	e.data[key] = value
	return nil
}

func (e *Engine) DEL(ctx context.Context, key string) error {
	_ = ctx
	delete(e.data, key)
	return nil
}
