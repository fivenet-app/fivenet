package filestore

import (
	"testing"

	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUploadFilterCloneIsIndependent(t *testing.T) {
	t.Parallel()

	original := NewUploadFilter([]string{"image/png"}, []string{"png"})
	clone := original.clone()
	require.NotSame(t, original, clone)

	clone.allowedContentTypes["image/jpeg"] = struct{}{}
	clone.allowedExtensions["jpg"] = struct{}{}
	original.allowedContentTypes["image/gif"] = struct{}{}
	original.allowedExtensions["gif"] = struct{}{}

	assert.Equal(t, []string{"image/gif", "image/png"}, original.AllowedContentTypes())
	assert.Equal(t, []string{"gif", "png"}, original.AllowedExtensions())
	assert.Equal(t, []string{"image/jpeg", "image/png"}, clone.AllowedContentTypes())
	assert.Equal(t, []string{"jpg", "png"}, clone.AllowedExtensions())
}

func TestUploadFilterCloneNil(t *testing.T) {
	t.Parallel()

	assert.Nil(t, (*UploadFilter)(nil).clone())
}

func TestBuildKeyUsesNamespaceDateUUIDAndSanitizedName(t *testing.T) {
	t.Parallel()

	key := buildKey("documents", "report_final.pdf")
	assert.Regexp(t, `^documents/[0-9]{8}/[0-9a-f-]{36}-report_final\.pdf$`, key)
}

//nolint:paralleltest // This test modifies a package-global variable, so it shouldn't be run in parallel with other tests.
func TestAddTableIgnoresNilAndRegistersValidTables(t *testing.T) {
	// tablesList is package-global, so isolate this test from other registrations.
	tableListsMu.Lock()
	previous := append([]joinInfo(nil), tablesList...)
	tablesList = nil
	tableListsMu.Unlock()
	t.Cleanup(func() {
		tableListsMu.Lock()
		tablesList = previous
		tableListsMu.Unlock()
	})

	AddTable(joinInfo{})
	tableListsMu.Lock()
	assert.Empty(t, tablesList)
	tableListsMu.Unlock()

	AddTable(
		joinInfo{Table: table.FivenetDocumentsFiles, FileCol: table.FivenetDocumentsFiles.FileID},
	)
	tableListsMu.Lock()
	defer tableListsMu.Unlock()
	require.Len(t, tablesList, 1)
	assert.Equal(t, table.FivenetDocumentsFiles, tablesList[0].Table)
}
