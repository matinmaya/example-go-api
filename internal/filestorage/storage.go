package filestorage

import "context"

// Storage defines storage backend operations.
type Storage interface {
	// Upload bytes: returns publicly accessible path or key.
	Upload(ctx context.Context, key string, data []byte, contentType string) (string, error)
	Read(ctx context.Context, key string) ([]byte, string, error) // returns data, contentType
	Delete(ctx context.Context, key string) error
	URL(ctx context.Context, key string) (string, error) // optional public URL
}
