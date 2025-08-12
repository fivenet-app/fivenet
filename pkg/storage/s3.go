package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
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

// S3 implements IStorage for S3-compatible object storage backends (e.g., MinIO, AWS S3).
type S3 struct {
	IStorage

	// s3 is the MinIO client instance for S3 operations.
	s3 *minio.Client
	// bucketName is the S3 bucket used for storage.
	bucketName string
	// prefix is an optional prefix for namespacing objects in the bucket.
	prefix string
}

// NewS3 creates a new S3 storage backend using the provided parameters and MinIO client.
// It optionally checks bucket existence on startup and registers a start hook for lifecycle management.
func NewS3(p Params) (IStorage, error) {
	transport := otelhttp.NewTransport(
		http.DefaultClient.Transport,
		otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
			return otelhttptrace.NewClientTrace(ctx, otelhttptrace.WithTracerProvider(p.TP))
		}),
	)

	// Initialize minio client object.
	mc, err := minio.New(p.Cfg.Storage.S3.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(
			p.Cfg.Storage.S3.AccessKeyID,
			p.Cfg.Storage.S3.SecretAccessKey,
			"",
		),
		Secure:     p.Cfg.Storage.S3.UseSSL,
		Region:     p.Cfg.Storage.S3.Region,
		MaxRetries: p.Cfg.Storage.S3.Retries,
		Transport:  transport,
	})
	if err != nil {
		return nil, err
	}

	s := &S3{
		s3:         mc,
		bucketName: p.Cfg.Storage.S3.BucketName,
		prefix:     p.Cfg.Storage.S3.Prefix,
	}

	if p.Cfg.Storage.S3.CheckOnStartup {
		p.LC.Append(fx.StartHook(func(ctx context.Context) error {
			exists, err := s.s3.BucketExists(ctx, s.bucketName)
			if err != nil {
				return err
			}
			if !exists {
				return fmt.Errorf(
					"storage: s3 bucket '%s' doesn't exist/can't access",
					s.bucketName,
				)
			}

			return nil
		}))
	}

	return s, nil
}

// WithPrefix returns a new S3 instance with the given prefix, sharing the same client and bucket.
func (s *S3) WithPrefix(prefix string) (IStorage, error) {
	return &S3{
		s3:         s.s3,
		bucketName: s.bucketName,
		prefix:     path.Join(s.prefix, prefix),
	}, nil
}

// Get retrieves an object and its metadata from S3 storage.
// Returns an open object and ObjectInfo, or an error if not found or invalid.
func (s *S3) Get(ctx context.Context, keyIn string) (IObject, IObjectInfo, error) {
	key, ok := utils.CleanFilePath(keyIn)
	if !ok {
		return nil, nil, ErrInvalidPath
	}
	key = path.Join(s.prefix, key)

	object, err := s.s3.GetObject(ctx, s.bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == minio.NoSuchKey {
			return nil, nil, ErrNotFound
		}
		return nil, nil, err
	}

	// Retrieve object info
	info, err := object.Stat()
	if err != nil {
		if minio.ToErrorResponse(err).Code == minio.NoSuchKey {
			return nil, nil, ErrNotFound
		}

		return nil, nil, err
	}

	return object, &ObjectInfo{
		name:         strings.TrimPrefix(info.Key, s.prefix),
		extension:    strings.TrimPrefix(filepath.Ext(info.Key), "."),
		contentType:  info.ContentType,
		size:         info.Size,
		lastModified: info.LastModified,
		expiration:   info.Expiration,
	}, nil
}

// Stat returns metadata for an object in S3 storage, or an error if not found or invalid.
func (s *S3) Stat(ctx context.Context, keyIn string) (IObjectInfo, error) {
	key, ok := utils.CleanFilePath(keyIn)
	if !ok {
		return nil, ErrInvalidPath
	}
	key = path.Join(s.prefix, key)

	info, err := s.s3.StatObject(ctx, s.bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == minio.NoSuchKey {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &ObjectInfo{
		name:         strings.TrimPrefix(info.Key, s.prefix),
		lastModified: info.LastModified,
		contentType:  info.ContentType,
		size:         info.Size,
		expiration:   info.Expiration,
	}, nil
}

// Put uploads an object to S3 storage. The file path must end with a file extension (e.g., `jpg`, `png`).
func (s *S3) Put(
	ctx context.Context,
	keyIn string,
	reader io.Reader,
	size int64,
	contentType string,
) (string, error) {
	return s.PutWithTTL(ctx, keyIn, reader, size, contentType, time.Time{})
}

// PutWithTTL uploads an object to S3 storage with an optional expiration time (TTL).
func (s *S3) PutWithTTL(
	ctx context.Context,
	keyIn string,
	reader io.Reader,
	size int64,
	contentType string,
	ttl time.Time,
) (string, error) {
	key, ok := utils.CleanFilePath(keyIn)
	if !ok {
		return "", ErrInvalidPath
	}
	key = path.Join(s.prefix, key)

	putOpts := minio.PutObjectOptions{
		ContentType: contentType,
		Expires:     ttl,
	}
	if size < 0 || size > 5<<20 { // 5 MiB
		putOpts.PartSize = 5 << 20 // 5 MiB
	}

	info, err := s.s3.PutObject(ctx, s.bucketName, key, reader, size, putOpts)
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(info.Key, s.prefix), nil
}

// Delete removes an object from S3 storage. Returns nil if the object does not exist.
func (s *S3) Delete(ctx context.Context, keyIn string) error {
	key, ok := utils.CleanFilePath(keyIn)
	if !ok {
		return ErrInvalidPath
	}
	key = path.Join(s.prefix, key)

	if err := s.s3.RemoveObject(ctx, s.bucketName, key, minio.RemoveObjectOptions{}); err != nil {
		return err
	}

	return nil
}

// List returns a list of objects and their metadata from S3 storage, supporting offset and page size.
// Returns an error if the prefix is invalid or listing fails.
func (s *S3) List(
	ctx context.Context,
	keyIn string,
	offset int,
	pageSize int,
) ([]*FileInfo, error) {
	key, ok := utils.CleanFilePath(keyIn)
	if !ok {
		return nil, ErrInvalidPath
	}
	key = path.Join(s.prefix, key)
	if key == "." {
		key = ""
	}

	i := 0
	files := []*FileInfo{}

	opts := minio.ListObjectsOptions{
		Recursive: true,
		Prefix:    key,
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
			Name:         strings.TrimPrefix(object.Key, s.prefix),
			LastModified: object.LastModified,
			Size:         object.Size,
			ContentType:  contentType.MIME.Value,
		})

		i++
	}

	return files, nil
}

// GetSpaceUsage calculates the total space used by objects in the S3 bucket.
// It iterates through all objects, summing their sizes.
// Returns the total size in bytes or an error if listing fails.
func (s *S3) GetSpaceUsage(ctx context.Context) (int64, error) {
	var totalSize int64

	opts := minio.ListObjectsOptions{
		Recursive: true,
		Prefix:    s.prefix,
	}
	for object := range s.s3.ListObjects(ctx, s.bucketName, opts) {
		if object.Err != nil {
			return 0, object.Err
		}
		totalSize += object.Size
	}

	return totalSize, nil
}
