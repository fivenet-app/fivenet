package qualifications

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/file"
	resqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	qualificationsactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/activity"
	qualificationsexam "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/exam"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	qualificationsstore "github.com/fivenet-app/fivenet/v2026/stores/qualifications"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type examHousekeeperTestStore struct {
	qualificationsstore.IStore

	events    []string
	responses *qualificationsexam.ExamResponses
	attemptID string
	questions *qualificationsexam.ExamQuestions
	imageIDs  []int64
	cutoff    time.Time
}

func (s *examHousekeeperTestStore) GetQualification(
	_ context.Context,
	_ int64,
	_ *userinfo.UserInfo,
	_ bool,
) (*resqualifications.Qualification, error) {
	s.events = append(s.events, "get-qualification")
	return &resqualifications.Qualification{
		ExamMode: qualificationsexam.QualificationExamMode_QUALIFICATION_EXAM_MODE_ENABLED,
		ExamSettings: &qualificationsexam.QualificationExamSettings{
			AutoGrade: false,
		},
	}, nil
}

func (s *examHousekeeperTestStore) ExpireExamUser(
	_ context.Context,
	_ qrm.DB,
	_ int64,
	_ int32,
	attemptID string,
) (bool, error) {
	s.events = append(s.events, "expire")
	s.attemptID = attemptID
	return true, nil
}

func (s *examHousekeeperTestStore) CreateQualificationActivity(
	_ context.Context,
	_ qrm.DB,
	activity *qualificationsactivity.QualificationActivity,
) error {
	s.events = append(s.events, "activity:"+activity.GetType().String())
	return nil
}

func (s *examHousekeeperTestStore) GetExamResponses(
	_ context.Context,
	_ qrm.DB,
	_ string,
) (*qualificationsexam.ExamResponses, *qualificationsexam.ExamGrading, error) {
	s.events = append(s.events, "get-responses")
	return s.responses, nil, nil
}

func (s *examHousekeeperTestStore) UpdateRequestStatus(
	_ context.Context,
	_ qrm.DB,
	_ int64,
	_ int32,
	status resqualifications.RequestStatus,
) error {
	s.events = append(s.events, "grade:"+status.String())
	return nil
}

func (s *examHousekeeperTestStore) GetExamQuestions(
	_ context.Context,
	_ qrm.DB,
	_ int64,
	_ bool,
) (*qualificationsexam.ExamQuestions, error) {
	return s.questions, nil
}

func (s *examHousekeeperTestStore) UnlinkStaleExamQuestionFiles(
	_ context.Context,
	_ *sql.Tx,
	_ int64,
	referencedFileIDs []int64,
	olderThan time.Time,
) (int64, error) {
	s.imageIDs = referencedFileIDs
	s.cutoff = olderThan
	return int64(len(referencedFileIDs)), nil
}

func TestExamHousekeeperGradesResponsesAfterExpiryClaim(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := &examHousekeeperTestStore{
		responses: &qualificationsexam.ExamResponses{Responses: []*qualificationsexam.ExamResponse{{
			QuestionId: 1,
		}}},
	}
	housekeeper := &ExamHousekeeper{
		store:  store,
		server: &Server{db: db, store: store, logger: zap.NewNop()},
	}

	mock.ExpectBegin()
	mock.ExpectCommit()
	err = housekeeper.completeExpiredExam(t.Context(), &qualificationsexam.ExamUser{
		QualificationId: 42,
		UserId:          7,
		AttemptId:       "expired-attempt",
		Snapshot: &qualificationsexam.ExamSnapshot{
			Exam: &qualificationsexam.ExamQuestions{},
		},
	})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, []string{
		"get-qualification",
		"expire",
		"activity:QUALIFICATION_ACTIVITY_TYPE_EXAM_EXPIRED",
		"get-responses",
		"grade:REQUEST_STATUS_EXAM_GRADING",
	}, store.events)
	assert.Equal(t, "expired-attempt", store.attemptID)
}

func TestExamHousekeeperCleanupExamQuestionFilesUsesPersistedImageReferences(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := &examHousekeeperTestStore{
		questions: &qualificationsexam.ExamQuestions{
			Questions: []*qualificationsexam.ExamQuestion{
				{
					Data: &qualificationsexam.ExamQuestionData{
						Data: &qualificationsexam.ExamQuestionData_Image{
							Image: &qualificationsexam.ExamQuestionImage{
								Image: &file.File{Id: 101},
							},
						},
					},
				},
				{Data: &qualificationsexam.ExamQuestionData{}},
			},
		},
	}
	housekeeper := &ExamHousekeeper{
		store:  store,
		server: &Server{db: db, store: store, logger: zap.NewNop()},
	}

	mock.ExpectBegin()
	mock.ExpectCommit()
	cutoff := time.Now().Add(-24 * time.Hour)
	require.NoError(t, housekeeper.cleanupExamQuestionFiles(t.Context(), 42, cutoff))
	require.NoError(t, mock.ExpectationsWereMet())

	assert.Equal(t, []int64{101}, store.imageIDs)
	assert.Equal(t, cutoff, store.cutoff)
}
