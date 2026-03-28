package grpcws

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTruncateBytesToMaxLen_NoTruncation(t *testing.T) {
	t.Parallel()

	input := []byte("grpc-web")
	output, truncated := truncateBytesToMaxLen(input, uint64(len(input)))

	require.False(t, truncated)
	require.Equal(t, input, output)
}

func TestTruncateBytesToMaxLen_Truncation(t *testing.T) {
	t.Parallel()

	input := []byte("grpc-web")
	output, truncated := truncateBytesToMaxLen(input, 4)

	require.True(t, truncated)
	require.Equal(t, []byte("grpc"), output)
}
