package sync

import (
	"context"
	"time"

	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
)

func (s *Server) SendVehicles(
	ctx context.Context,
	req *pbsync.SendVehiclesRequest,
) (*pbsync.SendDataResponse, error) {
	s.lastSyncedData.Store(time.Now().Unix())
	return s.store.SendVehicles(ctx, req)
}

func (s *Server) DeleteVehicles(
	ctx context.Context,
	req *pbsync.DeleteVehiclesRequest,
) (*pbsync.DeleteDataResponse, error) {
	return s.store.DeleteVehicles(ctx, req.GetPlates())
}
