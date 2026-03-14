package stats

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/stretchr/testify/require"
)

type testExtractor struct {
	sourceKey string
	supports  bool
	metrics   []*DocumentMetric
}

func (e *testExtractor) SourceKey() string { return e.sourceKey }

func (e *testExtractor) Supports(_ *documents.Document) bool { return e.supports }

func (e *testExtractor) Extract(
	_ context.Context,
	_ *documents.Document,
) ([]*DocumentMetric, error) {
	return e.metrics, nil
}

func TestService_RebuildDocumentMetrics_ReplacesBySource(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	svc := NewService(db)
	svc.extractors = []DocumentMetricExtractor{
		&testExtractor{
			sourceKey: "penalty_calculator",
			supports:  true,
			metrics: []*DocumentMetric{
				{
					DocumentID: 42,
					Job:        "police",
					SourceKey:  "penalty_calculator",
					MetricKey:  "case_count",
					Value:      1,
					OccurredAt: time.Now().UTC(),
				},
				{
					DocumentID: 42,
					Job:        "police",
					SourceKey:  "penalty_calculator",
					MetricKey:  "law_count",
					Dimension1: ptrString("10"),
					Value:      2,
					OccurredAt: time.Now().UTC(),
				},
			},
		},
	}

	doc := &documents.Document{
		Id:         42,
		CreatorJob: "police",
		CreatedAt:  timestamp.New(time.Now().UTC()),
		Meta:       &documents.DocumentMeta{Draft: false},
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM fivenet_documents_stats_metric")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO fivenet_documents_stats_metric")).
		WillReturnResult(sqlmock.NewResult(1, 2))
	mock.ExpectCommit()

	err = svc.RebuildDocumentMetrics(t.Context(), doc)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestService_RebuildDocumentMetrics_MultiExtractorDeletesBothSources(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	svc := NewService(db)
	svc.extractors = []DocumentMetricExtractor{
		&testExtractor{
			sourceKey: "alpha",
			supports:  false,
		},
		&testExtractor{
			sourceKey: "beta",
			supports:  true,
			metrics: []*DocumentMetric{{
				DocumentID: 43,
				Job:        "police",
				SourceKey:  "beta",
				MetricKey:  "case_count",
				Value:      1,
				OccurredAt: time.Now().UTC(),
			}},
		},
	}

	doc := &documents.Document{
		Id:         43,
		CreatorJob: "police",
		CreatedAt:  timestamp.New(time.Now().UTC()),
		Meta:       &documents.DocumentMeta{Draft: false},
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM fivenet_documents_stats_metric")).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM fivenet_documents_stats_metric")).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO fivenet_documents_stats_metric")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = svc.RebuildDocumentMetrics(t.Context(), doc)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestService_RebuildDocumentMetrics_UnpublishedClearsAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	svc := NewService(db)
	doc := &documents.Document{
		Id:         44,
		CreatorJob: "police",
		Meta:       &documents.DocumentMeta{Draft: true},
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM fivenet_documents_stats_metric")).
		WillReturnResult(sqlmock.NewResult(0, 3))
	mock.ExpectCommit()

	err = svc.RebuildDocumentMetrics(t.Context(), doc)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func ptrString(v string) *string {
	return &v
}

func TestService_BuildEmployeeCountMetrics(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	svc := NewService(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM fivenet_documents_stats_daily_rollup")).
		WillReturnResult(sqlmock.NewResult(0, 3))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO fivenet_documents_stats_daily_rollup")).
		WillReturnResult(sqlmock.NewResult(3, 3))
	mock.ExpectCommit()

	err = svc.BuildEmployeeCountMetrics(t.Context())
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestService_BuildEmployeeCountMetrics_DeleteError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	svc := NewService(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM fivenet_documents_stats_daily_rollup")).
		WillReturnError(errors.New("delete failed"))
	mock.ExpectRollback()

	err = svc.BuildEmployeeCountMetrics(t.Context())
	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
