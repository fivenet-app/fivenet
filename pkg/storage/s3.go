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
)

func init() {
	storageFactories["s3"] = NewS3
}

type S3 struct {
	IStorage

	s3         *minio.Client
	bucketName string
	prefix     string
}

func NewS3(cfg *config.Config) (IStorage, error) {
	// Initialize minio client object.
	mc, err := minio.New(cfg.Storage.S3.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Storage.S3.AccessKeyID, cfg.Storage.S3.SecretAccessKey, ""),
		Secure: cfg.Storage.S3.UseSSL,
		Region: cfg.Storage.S3.Region,
	})
	if err != nil {
		return nil, err
	}

	s := &S3{
		s3:         mc,
		bucketName: cfg.Storage.S3.BucketName,
	}

	return s, nil
}

func (s *S3) WithPrefix(prefix string) (IStorage, error) {
	return &S3{
		s3:         s.s3,
		bucketName: s.bucketName,
		prefix:     prefix,
	}, nil
}

func (s *S3) Get(ctx context.Context, filePath string) (IObject, IObjectInfo, error) {
	filePath = path.Join(s.prefix, filePath)
	object, err := s.s3.GetObject(ctx, s.bucketName, filePath, minio.GetObjectOptions{})
	if err != nil {
		return nil, nil, err
	}

	// Retrieve object info
	info, err := object.Stat()
	if err != nil {
		return nil, nil, err
	}

	return object, &ObjectInfo{
		contentType: info.ContentType,
		size:        info.Size,
		expiration:  info.Expiration,
	}, nil
}

func (s *S3) Stat(ctx context.Context, filePath string) (IObjectInfo, error) {
	filePath = path.Join(s.prefix, filePath)

	info, err := s.s3.StatObject(ctx, s.bucketName, filePath, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return &ObjectInfo{
		contentType: info.ContentType,
		size:        info.Size,
		expiration:  info.Expiration,
	}, nil
}

// Put the file path must end with a file extension (e.g., `jpg`, `png`)
func (s *S3) Put(ctx context.Context, filePath string, reader io.Reader, size int64, contentType string) (string, error) {
	return s.PutWithTTL(ctx, filePath, reader, size, contentType, time.Time{})
}

func (s *S3) PutWithTTL(ctx context.Context, filePath string, reader io.Reader, size int64, contentType string, ttl time.Time) (string, error) {
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

func (s *S3) Delete(ctx context.Context, filePath string) error {
	filePath = path.Join(s.prefix, filePath)

	if err := s.s3.RemoveObject(ctx, s.bucketName, filePath, minio.RemoveObjectOptions{}); err != nil {
		return err
	}

	return nil
}
