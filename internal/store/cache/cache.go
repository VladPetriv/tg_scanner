package cache

import "context"

type Store interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value bool) error
}
