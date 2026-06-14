package livemapstore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	livemapmarkers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/livemap/markers"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreCreateMarker(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db, &config.CustomDB{})
	now := time.Unix(100, 0).UTC()
	expiresAt := timestamp.New(now.Add(24 * time.Hour))
	marker := &livemapmarkers.MarkerMarker{
		ExpiresAt: expiresAt,
		Name:      "Marker",
		Description: func() *string {
			s := "desc"
			return &s
		}(),
		X: 1.25,
		Y: 2.5,
		Postal: func() *string {
			p := "12345"
			return &p
		}(),
		Color: func() *string {
			c := "#fff"
			return &c
		}(),
		Type: livemapmarkers.MarkerType_MARKER_TYPE_DOT,
	}

	expectedQuery := regexp.QuoteMeta(`INSERT INTO fivenet_centrum_markers`) + `(?s).*` + regexp.QuoteMeta(`expires_at`) + `(?s).*` + regexp.QuoteMeta(`creator_id`)
	mock.ExpectExec(expectedQuery).
		WithArgs(marker.GetExpiresAt(), "police", marker.GetName(), marker.Description, marker.GetX(), marker.GetY(), marker.Postal, marker.Color, marker.GetType(), marker.GetData(), int32(3)).
		WillReturnResult(sqlmock.NewResult(55, 1))

	id, err := store.CreateMarker(t.Context(), marker, 3, "police")
	require.NoError(t, err)
	assert.Equal(t, int64(55), id)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreUpdateMarker(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db, &config.CustomDB{})
	marker := &livemapmarkers.MarkerMarker{
		Id:   42,
		Name: "Updated",
		X:    8.1,
		Y:    9.2,
	}

	expectedQuery := regexp.QuoteMeta(`UPDATE fivenet_centrum_markers SET`) + `(?s).*` + regexp.QuoteMeta(`fivenet_centrum_markers.job = ?`) + `(?s).*` + regexp.QuoteMeta(`fivenet_centrum_markers.id = ?`) + `(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)
	mock.ExpectExec(expectedQuery).
		WithArgs(nil, marker.GetName(), marker.Description, marker.GetX(), marker.GetY(), marker.Postal, marker.Color, marker.GetType(), marker.GetData(), "police", int64(42), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.UpdateMarker(t.Context(), marker, "police"))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreSoftDeleteMarker(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db, &config.CustomDB{})

	expectedQuery := regexp.QuoteMeta(`UPDATE fivenet_centrum_markers SET`) + `(?s).*` + regexp.QuoteMeta(`deleted_at = CURRENT_TIMESTAMP`) + `(?s).*` + regexp.QuoteMeta(`fivenet_centrum_markers.id = ?`) + `(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)
	mock.ExpectExec(expectedQuery).
		WithArgs(int64(99), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.SoftDeleteMarker(t.Context(), 99))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreGetMarker(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db, &config.CustomDB{})
	now := time.Unix(100, 0).UTC()
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_centrum_markers AS marker_marker`)+`(?s).*`+regexp.QuoteMeta(`LEFT JOIN fivenet_user AS user_short ON`)+`(?s).*`+regexp.QuoteMeta(`marker_marker.id = ?`)+`(?s).*`+regexp.QuoteMeta(`LIMIT ?;`)).
		WithArgs(int64(42), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"marker_marker.id",
			"marker_marker.created_at",
			"marker_marker.updated_at",
			"marker_marker.deleted_at",
			"marker_marker.expires_at",
			"marker_marker.job",
			"marker_marker.name",
			"marker_marker.description",
			"marker_marker.x",
			"marker_marker.y",
			"marker_marker.postal",
			"marker_marker.color",
			"marker_marker.marker_type",
			"marker_marker.marker_data",
			"marker_marker.creator_id",
			"user_short.id",
			"user_short.job",
			"user_short.job_grade",
			"user_short.firstname",
			"user_short.lastname",
			"user_short.sex",
			"user_short.phone_number",
		}).AddRow(
			int64(42),
			now,
			now,
			nil,
			now.Add(24*time.Hour),
			"police",
			"Marker",
			"desc",
			1.25,
			2.5,
			"12345",
			"#fff",
			int32(1),
			nil,
			int32(3),
			int32(3),
			"police",
			int32(16),
			"Jane",
			"Doe",
			"f",
			"555-1234",
		))

	marker, err := store.GetMarker(t.Context(), 42)
	require.NoError(t, err)
	require.NotNil(t, marker)
	assert.Equal(t, int64(42), marker.GetId())
	assert.Equal(t, "police", marker.GetJob())
	require.NotNil(t, marker.GetCreator())
	assert.Equal(t, "Jane", marker.GetCreator().GetFirstname())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListActiveMarkers(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db, &config.CustomDB{})
	now := time.Unix(100, 0).UTC()
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_centrum_markers AS marker_marker`) + `(?s).*` + regexp.QuoteMeta(`marker_marker.deleted_at IS NULL`) + `(?s).*` + regexp.QuoteMeta(`ORDER BY marker_marker.job ASC, marker_marker.id ASC`)).
		WillReturnRows(sqlmock.NewRows([]string{
			"marker_marker.id",
			"marker_marker.created_at",
			"marker_marker.updated_at",
			"marker_marker.deleted_at",
			"marker_marker.expires_at",
			"marker_marker.job",
			"marker_marker.name",
			"marker_marker.description",
			"marker_marker.x",
			"marker_marker.y",
			"marker_marker.postal",
			"marker_marker.color",
			"marker_marker.marker_type",
			"marker_marker.marker_data",
			"marker_marker.creator_id",
			"user_short.id",
			"user_short.job",
			"user_short.job_grade",
			"user_short.firstname",
			"user_short.lastname",
			"user_short.sex",
			"user_short.phone_number",
		}).AddRow(
			int64(42),
			now,
			now,
			nil,
			now.Add(24*time.Hour),
			"police",
			"Marker",
			"desc",
			1.25,
			2.5,
			"12345",
			"#fff",
			int32(1),
			nil,
			int32(3),
			int32(3),
			"police",
			int32(16),
			"Jane",
			"Doe",
			"f",
			"555-1234",
		))

	markers, err := store.ListActiveMarkers(t.Context())
	require.NoError(t, err)
	require.Len(t, markers, 1)
	assert.Equal(t, "police", markers[0].GetJob())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListDeletedMarkers(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db, &config.CustomDB{})
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_centrum_markers AS marker_marker`) + `(?s).*` + regexp.QuoteMeta(`marker_marker.deleted_at IS NOT NULL`) + `(?s).*` + regexp.QuoteMeta(`ORDER BY marker_marker.id ASC`)).
		WillReturnRows(sqlmock.NewRows([]string{
			"marker_marker.id",
			"marker_marker.job",
		}).AddRow(int64(99), "police").AddRow(int64(100), "ems"))

	markers, err := store.ListDeletedMarkers(t.Context())
	require.NoError(t, err)
	require.Len(t, markers, 2)
	assert.Equal(t, "police", markers[0].GetJob())
	assert.Equal(t, int64(100), markers[1].GetId())
	require.NoError(t, mock.ExpectationsWereMet())
}
