package proto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CompareGRPCError(t *testing.T, expected error, err error) {
	expectedStatus := status.FromContextError(expected)
	errStatus := status.FromContextError(err)
	assert.Equal(t, expectedStatus.Code(), errStatus.Code())
	assert.Equal(t, expectedStatus.Message(), errStatus.Message())
}

func CompareGRPCStatusCode(t *testing.T, code codes.Code, err error) {
	errStatus := status.FromContextError(err)
	assert.Equal(t, code, errStatus.Code())
}
