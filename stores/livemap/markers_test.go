package livemapstore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	livemapmarkers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/livemap/markers"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreCreateMarker(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(Params{DB: db, Enricher: mstlystcdata.NewDummyEnricher()})
	now := time.Unix(100, 0).UTC()
	expiresAt := timestamp.New(now.Add(24 * time.Hour))
	marker := &livemapmarkers.MarkerMarker{
		ExpiresAt:   expiresAt,
		Name:        "Marker",
		Description: new("desc"),
		X:           1.25,
		Y:           2.5,
		Postal:      new("12345"),
		Color:       new("#fff"),
		Type:        livemapmarkers.MarkerType_MARKER_TYPE_DOT,
	}

	expectedQuery := regexp.QuoteMeta(
		`INSERT INTO fivenet_centrum_markers`,
	) + `(?s).*` + regexp.QuoteMeta(
		`expires_at`,
	) + `(?s).*` + regexp.QuoteMeta(
		`creator_id`,
	)
	mock.ExpectExec(expectedQuery).
		WithArgs(marker.GetExpiresAt(), "police", marker.GetPublic(), marker.GetName(), marker.Description, marker.GetX(), marker.GetY(), marker.Postal, marker.Color, marker.GetType(), marker.GetData(), int32(3)).
		WillReturnResult(sqlmock.NewResult(55, 1))
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_centrum_markers AS marker_marker`)+`(?s).*`+regexp.QuoteMeta(`marker_marker.id = ?`)+`(?s).*`+regexp.QuoteMeta(`LIMIT ?;`)).
		WithArgs(int64(55), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"marker_marker.id",
			"marker_marker.created_at",
			"marker_marker.updated_at",
			"marker_marker.deleted_at",
			"marker_marker.expires_at",
			"marker_marker.job",
			"marker_marker.public",
			"marker_marker.name",
			"marker_marker.description",
			"marker_marker.x",
			"marker_marker.y",
			"marker_marker.postal",
			"marker_marker.color",
			"marker_marker.marker_type",
			"marker_marker.marker_data",
			"marker_marker.creator_id",
		}).AddRow(
			int64(55),
			now,
			now,
			nil,
			now.Add(24*time.Hour),
			"police",
			false,
			"Marker",
			"desc",
			1.25,
			2.5,
			"12345",
			"#fff",
			int32(1),
			nil,
			int32(3),
		))

	creatorID := int32(3)
	id, err := store.CreateMarker(t.Context(), marker, &creatorID, "police")
	require.NoError(t, err)
	assert.Equal(t, int64(55), id)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreCreateMarkerAllowsNilCreatorID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(Params{DB: db, Enricher: mstlystcdata.NewDummyEnricher()})
	now := time.Unix(100, 0).UTC()
	marker := &livemapmarkers.MarkerMarker{
		ExpiresAt: timestamp.New(now.Add(24 * time.Hour)),
		Name:      "Marker",
		X:         1.25,
		Y:         2.5,
		Type:      livemapmarkers.MarkerType_MARKER_TYPE_DOT,
	}

	expectedQuery := regexp.QuoteMeta(
		`INSERT INTO fivenet_centrum_markers`,
	) + `(?s).*` + regexp.QuoteMeta(
		`creator_id`,
	)
	mock.ExpectExec(expectedQuery).
		WithArgs(marker.GetExpiresAt(), "police", marker.GetPublic(), marker.GetName(), marker.Description, marker.GetX(), marker.GetY(), marker.Postal, marker.Color, marker.GetType(), marker.GetData(), nil).
		WillReturnResult(sqlmock.NewResult(55, 1))
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_centrum_markers AS marker_marker`)+`(?s).*`+regexp.QuoteMeta(`marker_marker.id = ?`)+`(?s).*`+regexp.QuoteMeta(`LIMIT ?;`)).
		WithArgs(int64(55), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"marker_marker.id",
			"marker_marker.created_at",
			"marker_marker.updated_at",
			"marker_marker.deleted_at",
			"marker_marker.expires_at",
			"marker_marker.job",
			"marker_marker.public",
			"marker_marker.name",
			"marker_marker.description",
			"marker_marker.x",
			"marker_marker.y",
			"marker_marker.postal",
			"marker_marker.color",
			"marker_marker.marker_type",
			"marker_marker.marker_data",
			"marker_marker.creator_id",
		}).AddRow(
			int64(55),
			now,
			now,
			nil,
			now.Add(24*time.Hour),
			"police",
			false,
			"Marker",
			nil,
			1.25,
			2.5,
			nil,
			nil,
			int32(1),
			nil,
			nil,
		))

	id, err := store.CreateMarker(t.Context(), marker, nil, "police")
	require.NoError(t, err)
	assert.Equal(t, int64(55), id)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreUpdateMarker(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(Params{DB: db, Enricher: mstlystcdata.NewDummyEnricher()})
	now := time.Unix(100, 0).UTC()
	marker := &livemapmarkers.MarkerMarker{
		Id:   42,
		Name: "Updated",
		X:    8.1,
		Y:    9.2,
	}

	expectedQuery := regexp.QuoteMeta(
		`UPDATE fivenet_centrum_markers SET`,
	) + `(?s).*` + regexp.QuoteMeta(
		`fivenet_centrum_markers.job = ?`,
	) + `(?s).*` + regexp.QuoteMeta(
		`fivenet_centrum_markers.id = ?`,
	) + `(?s).*` + regexp.QuoteMeta(
		`LIMIT ?;`,
	)
	mock.ExpectExec(expectedQuery).
		WithArgs(nil, marker.GetName(), marker.Description, marker.GetX(), marker.GetY(), marker.Postal, marker.Color, marker.GetPublic(), marker.GetType(), marker.GetData(), "police", int64(42), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_centrum_markers AS marker_marker`)+`(?s).*`+regexp.QuoteMeta(`marker_marker.id = ?`)+`(?s).*`+regexp.QuoteMeta(`LIMIT ?;`)).
		WithArgs(int64(42), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"marker_marker.id",
			"marker_marker.created_at",
			"marker_marker.updated_at",
			"marker_marker.deleted_at",
			"marker_marker.expires_at",
			"marker_marker.job",
			"marker_marker.public",
			"marker_marker.name",
			"marker_marker.description",
			"marker_marker.x",
			"marker_marker.y",
			"marker_marker.postal",
			"marker_marker.color",
			"marker_marker.marker_type",
			"marker_marker.marker_data",
			"marker_marker.creator_id",
		}).AddRow(
			int64(42),
			now,
			now,
			nil,
			now.Add(24*time.Hour),
			"police",
			false,
			"Updated",
			nil,
			8.1,
			9.2,
			nil,
			nil,
			int32(0),
			nil,
			int32(3),
		))

	require.NoError(t, store.UpdateMarker(t.Context(), marker, "police"))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreDeleteMarker(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(Params{DB: db, Enricher: mstlystcdata.NewDummyEnricher()})
	now := time.Unix(100, 0).UTC()
	deletedAt := time.Unix(0, 0).UTC()

	deleteQuery := regexp.QuoteMeta(
		`UPDATE fivenet_centrum_markers SET`,
	) + `(?s).*` + regexp.QuoteMeta(
		`deleted_at = CAST(? AS DATETIME)`,
	) + `(?s).*` + regexp.QuoteMeta(
		`fivenet_centrum_markers.id = ?`,
	) + `(?s).*` + regexp.QuoteMeta(
		`LIMIT ?;`,
	)
	mock.ExpectExec(deleteQuery).
		WithArgs(deletedAt, int64(99), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_centrum_markers AS marker_marker`)+`(?s).*`+regexp.QuoteMeta(`marker_marker.id = ?`)+`(?s).*`+regexp.QuoteMeta(`LIMIT ?;`)).
		WithArgs(int64(99), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"marker_marker.id",
			"marker_marker.created_at",
			"marker_marker.updated_at",
			"marker_marker.deleted_at",
			"marker_marker.expires_at",
			"marker_marker.job",
			"marker_marker.public",
			"marker_marker.name",
			"marker_marker.description",
			"marker_marker.x",
			"marker_marker.y",
			"marker_marker.postal",
			"marker_marker.color",
			"marker_marker.marker_type",
			"marker_marker.marker_data",
			"marker_marker.creator_id",
		}).AddRow(
			int64(99),
			now,
			now,
			deletedAt,
			now.Add(24*time.Hour),
			"police",
			false,
			"Marker",
			"desc",
			1.25,
			2.5,
			"12345",
			"#fff",
			int32(1),
			nil,
			int32(3),
		))

	require.NoError(t, store.DeleteMarker(t.Context(), 99, timestamp.New(deletedAt)))

	restoreQuery := regexp.QuoteMeta(
		`UPDATE fivenet_centrum_markers SET`,
	) + `(?s).*` + regexp.QuoteMeta(
		`deleted_at = NULL`,
	) + `(?s).*` + regexp.QuoteMeta(
		`fivenet_centrum_markers.id = ?`,
	) + `(?s).*` + regexp.QuoteMeta(
		`LIMIT ?;`,
	)
	mock.ExpectExec(restoreQuery).
		WithArgs(int64(99), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_centrum_markers AS marker_marker`)+`(?s).*`+regexp.QuoteMeta(`marker_marker.id = ?`)+`(?s).*`+regexp.QuoteMeta(`LIMIT ?;`)).
		WithArgs(int64(99), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"marker_marker.id",
			"marker_marker.created_at",
			"marker_marker.updated_at",
			"marker_marker.deleted_at",
			"marker_marker.expires_at",
			"marker_marker.job",
			"marker_marker.public",
			"marker_marker.name",
			"marker_marker.description",
			"marker_marker.x",
			"marker_marker.y",
			"marker_marker.postal",
			"marker_marker.color",
			"marker_marker.marker_type",
			"marker_marker.marker_data",
			"marker_marker.creator_id",
		}).AddRow(
			int64(99),
			now,
			now,
			nil,
			now.Add(24*time.Hour),
			"police",
			false,
			"Marker",
			"desc",
			1.25,
			2.5,
			"12345",
			"#fff",
			int32(1),
			nil,
			int32(3),
		))

	require.NoError(t, store.DeleteMarker(t.Context(), 99, nil))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreGetMarker(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(Params{DB: db, Enricher: mstlystcdata.NewDummyEnricher()})
	now := time.Unix(100, 0).UTC()
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_centrum_markers AS marker_marker`)+`(?s).*`+regexp.QuoteMeta(`marker_marker.id = ?`)+`(?s).*`+regexp.QuoteMeta(`LIMIT ?;`)).
		WithArgs(int64(42), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"marker_marker.id",
			"marker_marker.created_at",
			"marker_marker.updated_at",
			"marker_marker.deleted_at",
			"marker_marker.expires_at",
			"marker_marker.job",
			"marker_marker.public",
			"marker_marker.name",
			"marker_marker.description",
			"marker_marker.x",
			"marker_marker.y",
			"marker_marker.postal",
			"marker_marker.color",
			"marker_marker.marker_type",
			"marker_marker.marker_data",
			"marker_marker.creator_id",
		}).AddRow(
			int64(42),
			now,
			now,
			nil,
			now.Add(24*time.Hour),
			"police",
			false,
			"Marker",
			"desc",
			1.25,
			2.5,
			"12345",
			"#fff",
			int32(1),
			nil,
			int32(3),
		))

	marker, err := store.GetMarker(t.Context(), 42)
	require.NoError(t, err)
	require.NotNil(t, marker)
	assert.Equal(t, int64(42), marker.GetId())
	assert.Equal(t, "police", marker.GetJob())
	require.Nil(t, marker.GetCreator())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListActiveMarkers(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(Params{DB: db, Enricher: mstlystcdata.NewDummyEnricher()})
	now := time.Unix(100, 0).UTC()
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_centrum_markers AS marker_marker`) + `(?s).*` + regexp.QuoteMeta(`marker_marker.deleted_at IS NULL`) + `(?s).*` + regexp.QuoteMeta(`ORDER BY marker_marker.job ASC, marker_marker.id ASC`)).
		WillReturnRows(sqlmock.NewRows([]string{
			"marker_marker.id",
			"marker_marker.created_at",
			"marker_marker.updated_at",
			"marker_marker.deleted_at",
			"marker_marker.expires_at",
			"marker_marker.job",
			"marker_marker.public",
			"marker_marker.name",
			"marker_marker.description",
			"marker_marker.x",
			"marker_marker.y",
			"marker_marker.postal",
			"marker_marker.color",
			"marker_marker.marker_type",
			"marker_marker.marker_data",
			"marker_marker.creator_id",
		}).AddRow(
			int64(42),
			now,
			now,
			nil,
			now.Add(24*time.Hour),
			"police",
			false,
			"Marker",
			"desc",
			1.25,
			2.5,
			"12345",
			"#fff",
			int32(1),
			nil,
			int32(3),
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

	store := New(Params{DB: db, Enricher: mstlystcdata.NewDummyEnricher()})
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_centrum_markers AS marker_marker`) + `(?s).*` + regexp.QuoteMeta(`marker_marker.deleted_at IS NOT NULL`) + `(?s).*` + regexp.QuoteMeta(`ORDER BY marker_marker.id ASC`)).
		WillReturnRows(sqlmock.NewRows([]string{
			"marker_marker.id",
			"marker_marker.job",
			"marker_marker.public",
		}).AddRow(int64(99), "police", false).AddRow(int64(100), "ems", true))

	markers, err := store.ListDeletedMarkers(t.Context())
	require.NoError(t, err)
	require.Len(t, markers, 2)
	assert.Equal(t, "police", markers[0].GetJob())
	assert.Equal(t, int64(100), markers[1].GetId())
	require.NoError(t, mock.ExpectationsWereMet())
}
