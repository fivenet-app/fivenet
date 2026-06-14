package sync

import (
	"context"
	"time"

	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
)

func (s *Server) SendJobs(
	ctx context.Context,
	req *pbsync.SendJobsRequest,
) (*pbsync.SendDataResponse, error) {
	s.lastSyncedData.Store(time.Now().Unix())
	return s.store.SendJobs(ctx, req)
}
