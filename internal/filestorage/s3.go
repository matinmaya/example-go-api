package filestorage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Storage struct {
	client *s3.Client
	bucket string
	region string
}

func NewS3Storage(cfg *config.LoadOptions, bucket, region string) (*S3Storage, error) {
	// NOTE: cfg param is unused here in signature for flexibility; we'll use default config loader.
	return nil, errors.New("not implemented: use NewS3StorageFromEnv")
}

func NewS3StorageFromEnv(ctx context.Context) (*S3Storage, error) {
	bucket := os.Getenv("S3_BUCKET")
	region := os.Getenv("AWS_REGION")
	if bucket == "" || region == "" {
		return nil, fmt.Errorf("S3_BUCKET and AWS_REGION must be set for s3 storage")
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(cfg)
	return &S3Storage{client: client, bucket: bucket, region: region}, nil
}

func (s *S3Storage) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	uploader := manager.NewUploader(s.client)
	input := &s3.PutObjectInput{
		Bucket:      &s.bucket,
		Key:         &key,
		Body:        io.NopCloser(bytes.NewReader(data)),
		ContentType: &contentType,
		ACL:         types.ObjectCannedACLPublicRead,
	}
	if _, err := uploader.Upload(ctx, input); err != nil {
		return "", err
	}
	// Build public URL (this assumes public bucket/object or appropriate presigned policy)
	u := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, url.PathEscape(key))
	return u, nil
}

func (s *S3Storage) Read(ctx context.Context, key string) ([]byte, string, error) {
	get, err := s.client.GetObject(ctx, &s3.GetObjectInput{Bucket: &s.bucket, Key: &key})
	if err != nil {
		return nil, "", err
	}
	defer get.Body.Close()
	b, err := io.ReadAll(get.Body)
	if err != nil {
		return nil, "", err
	}
	ct := ""
	if get.ContentType != nil {
		ct = *get.ContentType
	}
	return b, ct, nil
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: &s.bucket, Key: &key})
	return err
}

func (s *S3Storage) URL(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("empty key")
	}
	u := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, url.PathEscape(key))
	return u, nil
}
