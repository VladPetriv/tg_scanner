package store

import "context"

type (
	CacheStore interface {
		Get(ctx context.Context, key string) (string, error)
		Set(ctx context.Context, key string, value bool) error
		GenerateKey(value interface{}) string
	}

	ImageStore interface {
		Send(ctx context.Context, path string, objectName string) (string, error)
	}
)
