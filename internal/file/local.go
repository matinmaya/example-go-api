package file

import (
	"context"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	base string
}

func NewLocalStorage(baseDir string) (*LocalStorage, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}
	return &LocalStorage{base: baseDir}, nil
}

func (l *LocalStorage) fullPath(key string) string {
	// key might start with /, clean it
	k := filepath.Clean(key)
	return filepath.Join(l.base, k)
}

func (l *LocalStorage) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	p := l.fullPath(key)
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return "", err
	}
	if err := os.WriteFile(p, data, 0644); err != nil {
		return "", err
	}
	return filepath.ToSlash(key), nil
}

func (l *LocalStorage) Read(ctx context.Context, key string) ([]byte, string, error) {
	p := l.fullPath(key)
	b, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, "", fs.ErrNotExist
		}
		return nil, "", err
	}
	ct := http.DetectContentType(b)
	return b, ct, nil
}

func (l *LocalStorage) Delete(ctx context.Context, key string) error {
	p := l.fullPath(key)
	if err := os.Remove(p); err != nil {
		return err
	}
	return nil
}

func (l *LocalStorage) URL(ctx context.Context, key string) (string, error) {
	return filepath.ToSlash(key), nil
}
