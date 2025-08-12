package sanitizer

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errSanitizeTest = errors.New("sanitize test error")

type callCount interface {
	GetCount() int
}

type reqSanitizeBase struct {
	called int
}

func (r *reqSanitizeBase) GetCount() int {
	return r.called
}

type reqNoSanitize struct {
	reqSanitizeBase
}

type reqSanitize struct {
	reqSanitizeBase

	err error
}

func (r *reqSanitize) Sanitize() error {
	r.called++
	return r.err
}

var tests = []struct {
	msg           string
	req           callCount
	err           error
	sanitizeCount int
}{
	{
		msg:           "Sanitize without error",
		req:           &reqSanitize{},
		err:           nil,
		sanitizeCount: 3,
	},
	{
		msg:           "No sanitization should be triggered",
		req:           &reqNoSanitize{},
		err:           nil,
		sanitizeCount: 0,
	},
	{
		msg: "Sanitize with error",
		req: &reqSanitize{
			err: errSanitizeTest,
		},
		err:           status.Error(codes.InvalidArgument, errSanitizeTest.Error()),
		sanitizeCount: 3,
	},
}

func TestUnaryServerInterceptor(t *testing.T) {
	ctx := t.Context()
	interceptor := UnaryServerInterceptor()

	for _, test := range tests {
		for range 3 {
			_, err := interceptor(ctx, test.req, nil, unaryHandler)
			if test.err == nil {
				require.NoError(t, err, test.msg)
			} else {
				assert.Equal(t, test.err.Error(), err.Error(), test.msg)
			}
		}
		assert.Equal(t, test.req.GetCount(), test.sanitizeCount)
	}
}

func unaryHandler(ctx context.Context, req any) (any, error) {
	return nil, nil
}
