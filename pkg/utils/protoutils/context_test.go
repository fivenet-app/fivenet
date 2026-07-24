package protoutils

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsContextCanceled(t *testing.T) {
	t.Parallel()

	assert.False(t, IsContextCanceled(nil))
	assert.True(t, IsContextCanceled(context.Canceled))
	assert.True(t, IsContextCanceled(fmt.Errorf("wrapped: %w", context.Canceled)))
	assert.False(t, IsContextCanceled(context.DeadlineExceeded))
}
