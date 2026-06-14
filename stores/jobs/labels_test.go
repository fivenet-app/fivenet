package jobsstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	jobslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/labels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreValidateLabels(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_labels AS label`)).
		WithArgs("police", int64(7), int64(2), int64(10)).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(2)))

	valid, err := store.ValidateLabels(
		t.Context(),
		store.db,
		"police",
		[]*jobslabels.Label{{Id: 7}, {Id: 2}},
	)
	require.NoError(t, err)
	assert.True(t, valid)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreGetUsersLabels(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_colleague_labels`)).
		WithArgs(int32(1), "police").
		WillReturnRows(sqlmock.NewRows([]string{"label.id", "label.job", "label.name", "label.color", "label.icon", "label.sort_order"}).
			AddRow(int64(11), "police", "alpha", "#111111", nil, int32(0)))
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_colleague_labels`)).
		WithArgs(int32(2), "police").
		WillReturnRows(sqlmock.NewRows([]string{
			"label.id",
			"label.job",
			"label.name",
			"label.color",
			"label.icon",
			"label.sort_order",
		}).
			AddRow(int64(12), "police", "beta", "#222222", nil, int32(0)))

	labels, err := store.GetUsersLabels(t.Context(), store.db, "police", []int32{1, 2})
	require.NoError(t, err)
	require.Len(t, labels, 2)
	require.NotNil(t, labels[0].Labels)
	require.NotNil(t, labels[1].Labels)
	assert.Equal(t, int32(1), labels[0].UserId)
	assert.Equal(t, "alpha", labels[0].Labels.GetList()[0].GetName())
	assert.Equal(t, int32(2), labels[1].UserId)
	assert.Equal(t, "beta", labels[1].Labels.GetList()[0].GetName())
	require.NoError(t, mock.ExpectationsWereMet())
}
