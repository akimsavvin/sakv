package enginemock

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type Engine struct {
	mock.Mock
}

func NewEngine() *Engine {
	return new(Engine)
}

func (e *Engine) GET(ctx context.Context, key string) (string, error) {
	args := e.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (e *Engine) SET(ctx context.Context, key, value string) error {
	args := e.Called(ctx, key, value)
	return args.Error(0)
}

func (e *Engine) DEL(ctx context.Context, key string) error {
	args := e.Called(ctx, key)
	return args.Error(0)
}
