package image

import "context"

type Store interface {
	Send(ctx context.Context, path string, objectName string) (string, error)
}
