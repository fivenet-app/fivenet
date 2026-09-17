package qualifications

import (
	"context"
	"testing"

	qualificationspb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	qualificationsstore "github.com/fivenet-app/fivenet/v2026/stores/qualifications"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/stretchr/testify/require"
)

type resultAttemptCleanupTestStore struct {
	qualificationsstore.IStore

	deletedRequestAttemptID  string
	deletedResponseAttemptID string
}

func (s *resultAttemptCleanupTestStore) DeleteQualificationRequestByAttemptID(
	_ context.Context,
	_ qrm.DB,
	attemptID string,
) error {
	s.deletedRequestAttemptID = attemptID
	return nil
}

func (s *resultAttemptCleanupTestStore) DeleteExamResponses(
	_ context.Context,
	_ qrm.DB,
	attemptID string,
) error {
	s.deletedResponseAttemptID = attemptID
	return nil
}

func TestDeleteQualificationResultAttemptCleansOnlyMatchingRequestAndResponses(t *testing.T) {
	t.Parallel()

	store := &resultAttemptCleanupTestStore{}
	server := &Server{store: store}
	attemptID := "attempt-old"
	result := &qualificationspb.QualificationResult{ExamAttemptId: &attemptID}

	require.NoError(t, server.deleteQualificationResultAttempt(t.Context(), nil, result))
	require.Equal(t, "attempt-old", store.deletedRequestAttemptID)
	require.Equal(t, "attempt-old", store.deletedResponseAttemptID)
}

func TestDeleteQualificationResultAttemptIgnoresLegacyUnlinkedResult(t *testing.T) {
	t.Parallel()

	store := &resultAttemptCleanupTestStore{}
	server := &Server{store: store}

	require.NoError(t, server.deleteQualificationResultAttempt(
		t.Context(), nil, &qualificationspb.QualificationResult{},
	))
	require.Empty(t, store.deletedRequestAttemptID)
	require.Empty(t, store.deletedResponseAttemptID)
}
