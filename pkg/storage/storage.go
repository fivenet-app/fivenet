package storage

import (
	"context"
	"io"
	"mime"

	"github.com/galexrt/fivenet/pkg/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/fx"
)

var Module = fx.Module("storage",
	fx.Provide(New),
)

type Storage struct {
	s3         *minio.Client
	bucketName string
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

func (s *Storage) Get(ctx context.Context, filePath string) (Object, error) {
	object, err := s.s3.GetObject(ctx, s.bucketName, filePath, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return object, nil
}

// Put the file path must end with a file extension (e.g., `jpg`, `png`)
func (s *Storage) Put(ctx context.Context, filePath string, reader io.Reader, size int64, contentType string) (string, error) {
	uploadInfo, err := s.s3.PutObject(ctx, s.bucketName, filePath, reader, size, minio.PutObjectOptions{
		ContentType: mime.TypeByExtension(filePath),
	})
	if err != nil {
		return "", err
	}

	return uploadInfo.Key, nil
}

func (s *Storage) Delete(ctx context.Context, filePath string) error {
	return s.s3.RemoveObject(ctx, s.bucketName, filePath, minio.RemoveObjectOptions{})
}
