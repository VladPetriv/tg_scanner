package cache

import "context"

type Store interface {
	Get(ctx context.Context, data interface{}) (string, error)
	Set(ctx context.Context, data interface{}, value bool) error
}
