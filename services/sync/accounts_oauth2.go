package sync

import (
	"context"
	"time"

	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
)

func (s *Server) AddUserOAuth2Conn(
	ctx context.Context,
	req *pbsync.AddUserOAuth2ConnRequest,
) (*pbsync.AddActivityResponse, error) {
	s.lastSyncedActivity.Store(time.Now().Unix())
	return s.store.AddUserOAuth2Conn(ctx, req)
}
