package storage

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/utils"
	"github.com/h2non/filetype"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/fx"
)

func init() {
	storageFactories[config.StorageTypeS3] = NewS3
}

type S3 struct {
	IStorage

	s3           *minio.Client
	bucketName   string
	prefix       string
	usePresigned bool
}

func NewS3(p Params) (IStorage, error) {
	transport := otelhttp.NewTransport(
		http.DefaultClient.Transport,
		otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
			return otelhttptrace.NewClientTrace(ctx, otelhttptrace.WithTracerProvider(p.TP))
		}),
	)

	// Initialize minio client object.
	mc, err := minio.New(p.Cfg.Storage.S3.Endpoint, &minio.Options{
		Creds:      credentials.NewStaticV4(p.Cfg.Storage.S3.AccessKeyID, p.Cfg.Storage.S3.SecretAccessKey, ""),
		Secure:     p.Cfg.Storage.S3.UseSSL,
		Region:     p.Cfg.Storage.S3.Region,
		MaxRetries: p.Cfg.Storage.S3.Retries,
		Transport:  transport,
	})
	if err != nil {
		return nil, err
	}

	s := &S3{
		s3:           mc,
		bucketName:   p.Cfg.Storage.S3.BucketName,
		prefix:       p.Cfg.Storage.S3.Prefix,
		usePresigned: p.Cfg.Storage.S3.UsePreSigned,
	}

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		exists, err := s.s3.BucketExists(ctx, s.bucketName)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("storage: s3 bucket '%s' doesn't exist/can't access", s.bucketName)
		}

		return nil
	}))

	return s, nil
}

func (s *S3) WithPrefix(prefix string) (IStorage, error) {
	return &S3{
		s3:         s.s3,
		bucketName: s.bucketName,
		prefix:     path.Join(s.prefix, prefix),
	}, nil
}

func (s *S3) Get(ctx context.Context, filePathIn string) (IObject, IObjectInfo, error) {
	filePath, ok := utils.CleanFilePath(filePathIn)
	if !ok {
		return nil, nil, ErrInvalidPath
	}
	filePath = path.Join(s.prefix, filePath)

	object, err := s.s3.GetObject(ctx, s.bucketName, filePath, minio.GetObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return nil, nil, ErrNotFound
		}
		return nil, nil, err
	}

	// Retrieve object info
	info, err := object.Stat()
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return nil, nil, ErrNotFound
		}

		return nil, nil, err
	}

	return object, &ObjectInfo{
		name:         info.Key,
		extension:    strings.TrimPrefix(filepath.Ext(info.Key), "."),
		contentType:  info.ContentType,
		size:         info.Size,
		lastModified: info.LastModified,
		expiration:   info.Expiration,
	}, nil
}

func (s *S3) GetURL(ctx context.Context, filePath string, expires time.Duration, reqParams url.Values) (*string, error) {
	if !s.usePresigned {
		return nil, nil
	}

	filePath, ok := utils.CleanFilePath(filePath)
	if !ok {
		return nil, ErrInvalidPath
	}
	filePath = path.Join(s.prefix, filePath)

	u, err := s.s3.PresignedGetObject(ctx, s.bucketName, filePath, expires, reqParams)
	if err != nil {
		return nil, err
	}

	url := u.String()
	return &url, nil
}

func (s *S3) Stat(ctx context.Context, filePathIn string) (IObjectInfo, error) {
	filePath, ok := utils.CleanFilePath(filePathIn)
	if !ok {
		return nil, ErrInvalidPath
	}
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
func (s *S3) Put(ctx context.Context, filePathIn string, reader io.Reader, size int64, contentType string) (string, error) {
	return s.PutWithTTL(ctx, filePathIn, reader, size, contentType, time.Time{})
}

func (s *S3) PutWithTTL(ctx context.Context, filePathIn string, reader io.Reader, size int64, contentType string, ttl time.Time) (string, error) {
	filePath, ok := utils.CleanFilePath(filePathIn)
	if !ok {
		return "", ErrInvalidPath
	}
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

func (s *S3) Delete(ctx context.Context, filePathIn string) error {
	filePath, ok := utils.CleanFilePath(filePathIn)
	if !ok {
		return ErrInvalidPath
	}
	filePath = path.Join(s.prefix, filePath)

	if err := s.s3.RemoveObject(ctx, s.bucketName, filePath, minio.RemoveObjectOptions{}); err != nil {
		return err
	}

	return nil
}

func (s *S3) List(ctx context.Context, filePathIn string, offset int, pageSize int) ([]*FileInfo, error) {
	filePath, ok := utils.CleanFilePath(filePathIn)
	if !ok {
		return nil, ErrInvalidPath
	}
	filePath = path.Join(s.prefix, filePath)
	if filePath == "." {
		filePath = ""
	}

	i := 0
	files := []*FileInfo{}

	opts := minio.ListObjectsOptions{
		Recursive: true,
		Prefix:    filePath,
	}
	for object := range s.s3.ListObjects(ctx, s.bucketName, opts) {
		if object.Err != nil {
			return nil, object.Err
		}

		if i < offset {
			i++
			continue
		}

		// Verify if we have listed page size count of objects.
		if i == offset+pageSize {
			break
		}

		contentType := filetype.GetType(strings.TrimPrefix(filepath.Ext(object.Key), "."))

		files = append(files, &FileInfo{
			Name:         object.Key,
			LastModified: object.LastModified,
			Size:         object.Size,
			ContentType:  contentType.MIME.Value,
		})

		i++
	}

	return files, nil
}
