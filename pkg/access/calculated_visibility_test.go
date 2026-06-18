package access

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollapseCalculatedVisibilityACLRows(t *testing.T) {
	t.Parallel()

	rows := []calculatedVisibilityACLRow{
		{SubjectID: 10, Access: 2, Effect: true},
		{SubjectID: 10, Access: 4, Effect: true},
		{SubjectID: 10, Access: 3, Effect: false},
		{SubjectID: 10, Access: 5, Effect: false},
	}

	out := collapseCalculatedVisibilityACLRows(rows)

	require.Len(t, out, 2)
	assert.Equal(t, int64(10), out[0].SubjectID)
	assert.True(t, out[0].Effect)
	assert.Equal(t, int32(4), out[0].Access)
	assert.Equal(t, int64(10), out[1].SubjectID)
	assert.False(t, out[1].Effect)
	assert.Equal(t, int32(5), out[1].Access)
}
