package engine

import "context"

type Engine interface {
	GET(ctx context.Context, key string) (string, error)
	SET(ctx context.Context, key, value string) error
	DEL(ctx context.Context, key string) error
}
