package file

import (
	"bytes"
	"context"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/chai2010/webp"
	"github.com/google/uuid"
)

// Service provides higher level file operations for handlers.
type Service struct {
	store Storage
}

func NewService(s Storage) *Service { return &Service{store: s} }

// SaveUpload reads a multipart file and stores it. Returns key or URL.
func (s *Service) SaveUpload(ctx context.Context, prefixPath string, fh *multipart.FileHeader) (string, error) {
	f, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	uuidFileName := uuid.New().String() + filepath.Ext(fh.Filename)
	key := prefixPath + filepath.Join(time.Now().Format("2006/01/02"), uuidFileName)

	contentType := fh.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(buf)
	}

	return s.store.Upload(ctx, key, buf, contentType)
}

func (s *Service) SaveImageUpload(ctx context.Context, prefixPath string, fh *multipart.FileHeader, toWebp bool) (string, error) {
	f, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	contentType := fh.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(buf)
	}

	if contentType == "" || contentType[:6] != "image/" {
		return "", http.ErrNotSupported
	}

	var data []byte
	var ext string
	if toWebp && contentType != "image/gif" {
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return "", err
		}
		var b bytes.Buffer
		// err = webp.Encode(&b, img, &webp.Options{Lossless: true}) // big size good qty
		err = webp.Encode(&b, img, &webp.Options{Lossless: false, Quality: 100})
		if err != nil {
			return "", err
		}
		data = b.Bytes()
		ext = ".webp"
		contentType = "image/webp"
	} else {
		data = buf
		ext = filepath.Ext(fh.Filename)
	}

	uuidFileName := uuid.New().String() + ext
	key := prefixPath + filepath.Join(time.Now().Format("2006/01/02"), uuidFileName)

	return s.store.Upload(ctx, key, data, contentType)
}

func (s *Service) Read(ctx context.Context, key string) ([]byte, string, error) {
	return s.store.Read(ctx, key)
}

func (s *Service) Delete(ctx context.Context, key string) error {
	return s.store.Delete(ctx, key)
}

func (s *Service) URL(ctx context.Context, key string) (string, error) {
	return s.store.URL(ctx, key)
}
