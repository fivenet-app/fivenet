package sync

import (
	"context"
	"time"

	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
)

func (s *Server) SendUserLocations(
	ctx context.Context,
	req *pbsync.SendUserLocationsRequest,
) (*pbsync.SendDataResponse, error) {
	s.lastSyncedData.Store(time.Now().Unix())
	return s.store.SendUserLocations(ctx, req)
}
