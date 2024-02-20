package storage

import (
	"context"
	"io"
	"mime"
	"path"
	"time"

	"github.com/galexrt/fivenet/pkg/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/fx"
)

var Module = fx.Module("storage",
	fx.Provide(New),
)

type IStorage interface {
	Get(ctx context.Context, filePath string) (Object, error)
	Put(ctx context.Context, filePath string, reader io.Reader, size int64, contentType string) (string, error)
	Delete(ctx context.Context, filePath string) error
}

type Storage struct {
	s3         *minio.Client
	bucketName string
	prefix     string
}

func New(cfg *config.Config) (*Storage, error) {
	if !cfg.Storage.Enabled {
		return nil, nil
	}

	// Initialize minio client object.
	mc, err := minio.New(cfg.Storage.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Storage.AccessKeyID, cfg.Storage.SecretAccessKey, ""),
		Secure: cfg.Storage.UseSSL,
		Region: cfg.Storage.Region,
	})
	if err != nil {
		return nil, err
	}

	st := &Storage{
		s3:         mc,
		bucketName: cfg.Storage.BucketName,
	}

	return st, nil
}

func (s *Storage) WithPrefix(prefix string) *Storage {
	return &Storage{
		s3:         s.s3,
		bucketName: s.bucketName,
		prefix:     prefix,
	}
}

func (s *Storage) Get(ctx context.Context, filePath string) (Object, error) {
	filePath = path.Join(s.prefix, filePath)
	object, err := s.s3.GetObject(ctx, s.bucketName, filePath, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return object, nil
}

// Put the file path must end with a file extension (e.g., `jpg`, `png`)
func (s *Storage) Put(ctx context.Context, filePath string, reader io.Reader, size int64, contentType string) (string, error) {
	return s.PutWithTTL(ctx, filePath, reader, size, contentType, time.Time{})
}

func (s *Storage) PutWithTTL(ctx context.Context, filePath string, reader io.Reader, size int64, contentType string, ttl time.Time) (string, error) {
	filePath = path.Join(s.prefix, filePath)
	uploadInfo, err := s.s3.PutObject(ctx, s.bucketName, filePath, reader, size, minio.PutObjectOptions{
		ContentType: mime.TypeByExtension(filePath),
		Expires:     ttl,
	})
	if err != nil {
		return "", err
	}

	return uploadInfo.Key, nil
}

func (s *Storage) Delete(ctx context.Context, filePath string) error {
	filePath = path.Join(s.prefix, filePath)
	return s.s3.RemoveObject(ctx, s.bucketName, filePath, minio.RemoveObjectOptions{})
}
