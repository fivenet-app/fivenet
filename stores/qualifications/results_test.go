package qualificationsstore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	resqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreGetQualificationResult(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	expectedQuery := regexp.QuoteMeta(
		`FROM fivenet_qualifications_results AS qualification_result`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`LEFT JOIN fivenet_user AS user ON`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`WHERE`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`GROUP BY qualification_result.id`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`LIMIT ?;`,
	)

	now := time.Unix(0, 0).UTC()
	mock.ExpectQuery(expectedQuery).
		WithArgs(true, int64(9), int64(42), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"qualification_result.id",
			"qualification_result.created_at",
			"qualification_result.deleted_at",
			"qualification_result.qualification_id",
			"qualification_result.user_id",
			"user.id",
			"user.job",
			"user.job_grade",
			"user.firstname",
			"user.lastname",
			"user.dateofbirth",
			"qualification_result.status",
			"qualification_result.score",
			"qualification_result.summary",
			"qualification_result.creator_id",
			"qualification_result.creator_job",
			"creator.id",
			"creator.job",
			"creator.job_grade",
			"creator.firstname",
			"creator.lastname",
			"creator.dateofbirth",
		}).AddRow(
			int64(9), now, nil, int64(42), int32(7),
			int32(7), "police", int32(5), "Jane", "Doe", now,
			int32(
				resqualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL,
			), nil, "ok", int32(1), "police",
			int32(1), "police", int32(5), "Jane", "Doe", now,
		),
		)

	result, err := store.GetQualificationResult(
		t.Context(),
		42,
		9,
		nil,
		&userinfo.UserInfo{UserId: 7, Job: "police"},
		7,
		false,
	)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, int64(9), result.GetId())
	assert.Equal(t, int64(42), result.GetQualificationId())
	require.NoError(t, mock.ExpectationsWereMet())
}
