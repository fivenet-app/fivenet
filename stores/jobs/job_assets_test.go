package jobsstore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreListJobAssetsScansSelectedColumns(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)
	createdAt := time.Date(2026, time.July, 31, 12, 0, 0, 0, time.UTC)

	mock.ExpectQuery(`(?s)FROM fivenet_job_assets .*INNER JOIN fivenet_files .*`).
		WithArgs("police", int64(10), int64(5)).
		WillReturnRows(sqlmock.NewRows([]string{
			"job_asset.id",
			"job_asset.created_by_user_id",
			"job_asset.display_name",
			"job_asset.file_path",
			"job_asset.byte_size",
			"job_asset.content_type",
			"job_asset.created_at",
		}).AddRow(
			int64(42),
			int32(7),
			"briefing.jpg",
			"jobassets/police/123-briefing.jpg",
			int64(2048),
			"image/jpeg",
			createdAt,
		))

	assets, err := store.ListJobAssets(t.Context(), store.db, "police", 5, 10)
	require.NoError(t, err)
	require.Len(t, assets, 1)

	asset := assets[0]
	assert.Equal(t, int64(42), asset.GetId())
	assert.Equal(t, int32(7), asset.GetCreatedByUserId())
	assert.Equal(t, "briefing.jpg", asset.GetDisplayName())
	assert.Equal(t, "jobassets/police/123-briefing.jpg", asset.GetFile().GetFilePath())
	assert.Equal(t, int64(2048), asset.GetFile().GetByteSize())
	assert.Equal(t, "image/jpeg", asset.GetFile().GetContentType())
	assert.Equal(t, createdAt, asset.GetCreatedAt().AsTime())
	assert.Equal(t, createdAt, asset.GetFile().GetCreatedAt().AsTime())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreGetJobAssetScansSelectedColumns(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)
	createdAt := time.Date(2026, time.July, 31, 12, 0, 0, 0, time.UTC)

	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_job_assets INNER JOIN fivenet_files")).
		WithArgs("police", int64(42), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"job_asset.id",
			"job_asset.created_by_user_id",
			"job_asset.display_name",
			"job_asset.file_path",
			"job_asset.byte_size",
			"job_asset.content_type",
			"job_asset.created_at",
		}).AddRow(
			int64(42),
			int32(7),
			"briefing.jpg",
			"jobassets/police/123-briefing.jpg",
			int64(2048),
			"image/jpeg",
			createdAt,
		))

	asset, err := store.GetJobAsset(t.Context(), store.db, "police", 42)
	require.NoError(t, err)
	require.NotNil(t, asset)
	assert.Equal(t, int64(42), asset.GetId())
	assert.Equal(t, int32(7), asset.GetCreatedByUserId())
	assert.Equal(t, "briefing.jpg", asset.GetDisplayName())
	require.NoError(t, mock.ExpectationsWereMet())
}
