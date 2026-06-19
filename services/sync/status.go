package sync

import (
	"context"
	"time"

	syncdata "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/data"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
)

func (s *Server) GetStatus(
	ctx context.Context,
	req *pbsync.GetStatusRequest,
) (*pbsync.GetStatusResponse, error) {
	resp := &pbsync.GetStatusResponse{}

	if v := s.lastSyncedData.Load(); v > 0 {
		resp.LastSyncedData = timestamp.New(time.Unix(v, 0))
	}

	if v := s.lastSyncedActivity.Load(); v > 0 {
		resp.LastSyncedActivity = timestamp.New(time.Unix(v, 0))
	}

	jobsCount, err := s.store.CountJobs(ctx)
	if err != nil {
		return nil, err
	}
	resp.Jobs = &syncdata.DataStatus{
		Count: jobsCount,
	}

	accountsCount, err := s.store.CountAccounts(ctx)
	if err != nil {
		return nil, err
	}
	resp.Accounts = &syncdata.DataStatus{
		Count: accountsCount,
	}

	usersCount, err := s.store.CountUsers(ctx)
	if err != nil {
		return nil, err
	}
	resp.Users = &syncdata.DataStatus{
		Count: usersCount,
	}

	vehiclesCount, err := s.store.CountVehicles(ctx)
	if err != nil {
		return nil, err
	}
	resp.Vehicles = &syncdata.DataStatus{
		Count: vehiclesCount,
	}

	licensesCount, err := s.store.CountLicenses(ctx)
	if err != nil {
		return nil, err
	}
	resp.Licenses = &syncdata.DataStatus{
		Count: licensesCount,
	}

	return resp, nil
}
