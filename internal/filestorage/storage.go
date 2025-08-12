package filestorage

import "context"

type Storage interface {
	Upload(ctx context.Context, key string, data []byte, contentType string) (string, error)
	Read(ctx context.Context, key string) ([]byte, string, error)
	Delete(ctx context.Context, key string) error
	URL(ctx context.Context, key string) (string, error)
}
