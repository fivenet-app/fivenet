package qualificationsstore

import (
	"database/sql/driver"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	qualificationsactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/activity"
	"github.com/stretchr/testify/require"
)

func TestStoreCreateQualificationActivityIsIdempotentForExamAttempt(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(testParams(db))
	activity := &qualificationsactivity.QualificationActivity{
		QualificationId: 42,
		Type:            qualificationsactivity.QualificationActivityType_QUALIFICATION_ACTIVITY_TYPE_EXAM_STARTED,
		ActorUserId:     int32Ptr(7),
		TargetUserId:    int32Ptr(7),
		AttemptId:       "attempt-1",
	}
	expectInsert := func(result driver.Result) {
		mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO fivenet_qualifications_activity")+
				`(?s).*`+regexp.QuoteMeta("ON DUPLICATE KEY UPDATE"),
		).
			WithArgs(
				int64(42),
				int32(
					qualificationsactivity.QualificationActivityType_QUALIFICATION_ACTIVITY_TYPE_EXAM_STARTED,
				),
				int32(7),
				int32(7),
				nil,
				"attempt-1",
			).
			WillReturnResult(result)
	}
	expectInsert(sqlmock.NewResult(1, 1))
	expectInsert(sqlmock.NewResult(1, 0))

	for range 2 {
		err = store.CreateQualificationActivity(t.Context(), db, activity)
		require.NoError(t, err)
	}
	require.NoError(t, mock.ExpectationsWereMet())
}

func int32Ptr(value int32) *int32 {
	return &value
}
