package qualificationsstore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreGetQualificationRequirements(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	expectedQuery := regexp.QuoteMeta(
		`FROM fivenet_qualifications_requirements AS qualification_requirement`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`INNER JOIN fivenet_qualifications AS target_qualification ON`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`qualification_requirement.qualification_id = ?`,
	)

	mock.ExpectQuery(expectedQuery).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{
			"qualification_requirement.id",
			"qualification_requirement.created_at",
			"qualification_requirement.target_qualification_id",
			"target_qualification.id",
			"target_qualification.abbreviation",
			"target_qualification.title",
		}).AddRow(int64(7), time.Unix(0, 0).UTC(), int64(99), int64(99), "ABC", "Target"))

	reqs, err := store.GetQualificationRequirements(t.Context(), 42)
	require.NoError(t, err)
	require.Len(t, reqs, 1)
	assert.Equal(t, int64(7), reqs[0].GetId())
	assert.Equal(t, int64(99), reqs[0].GetTargetQualificationId())
	require.NoError(t, mock.ExpectationsWereMet())
}
