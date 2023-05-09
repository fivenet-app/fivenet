package proto

import (
	"testing"

	"google.golang.org/grpc/status"
	"gotest.tools/assert"
)

func CompareGRPCError(t *testing.T, expected error, err error) {
	expectedStatus := status.FromContextError(expected)
	errStatus := status.FromContextError(err)
	assert.Equal(t, expectedStatus.Code(), errStatus.Code())
	assert.Equal(t, expectedStatus.Message(), errStatus.Message())
}
