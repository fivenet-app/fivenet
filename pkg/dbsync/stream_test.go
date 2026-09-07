package dbsync

import (
	"context"
	"testing"
	"time"

	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type recordingUsersSyncer struct {
	userIDs     []int32
	identifiers []string
}

func (s *recordingUsersSyncer) Sync(context.Context) (int64, int64, string, *time.Time, error) {
	return 0, 0, "", nil, nil
}

func (s *recordingUsersSyncer) SyncUser(_ context.Context, userID int32) error {
	s.userIDs = append(s.userIDs, userID)
	return nil
}

func (s *recordingUsersSyncer) SyncUserByIdentifier(_ context.Context, identifier string) error {
	s.identifiers = append(s.identifiers, identifier)
	return nil
}

func TestProcessUserSyncRequestProcessesIDsAndIdentifiers(t *testing.T) {
	t.Parallel()

	users := &recordingUsersSyncer{}
	syncer := &Sync{
		logger: zap.NewNop(),
		users:  users,
	}

	syncer.processUserSyncRequest(t.Context(), &pbsync.UserSyncRequest{
		UserIds:     []int32{42, 43},
		Identifiers: []string{"license:one", "license:two"},
	})

	require.Equal(t, []int32{42, 43}, users.userIDs)
	require.Equal(t, []string{"license:one", "license:two"}, users.identifiers)
}
