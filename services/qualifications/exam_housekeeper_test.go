package qualifications

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	resqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
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
) (bool, error) {
	s.events = append(s.events, "expire")
	return true, nil
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
		Snapshot: &qualificationsexam.ExamSnapshot{
			Exam: &qualificationsexam.ExamQuestions{},
		},
	})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, []string{
		"get-qualification",
		"expire",
		"get-responses",
		"grade:REQUEST_STATUS_EXAM_GRADING",
	}, store.events)
}
