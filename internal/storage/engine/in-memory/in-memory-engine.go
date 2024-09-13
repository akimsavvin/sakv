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

func (e *Engine) GET(_ context.Context, key string) (string, error) {
	val, ok := e.data[key]
	if !ok {
		return "", ErrNotFound
	}

	return val, nil
}

func (e *Engine) SET(_ context.Context, key, value string) error {
	e.data[key] = value
	return nil
}

func (e *Engine) DEL(_ context.Context, key string) error {
	delete(e.data, key)
	return nil
}
