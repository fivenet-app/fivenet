package sync

import (
	"context"
	"time"

	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
)

func (s *Server) AddActivity(
	ctx context.Context,
	req *pbsync.AddActivityRequest,
) (*pbsync.AddActivityResponse, error) {
	s.lastSyncedActivity.Store(time.Now().Unix())
	return s.store.AddActivity(ctx, req)
}

func (s *Server) AddUserActivity(
	ctx context.Context,
	req *pbsync.AddUserActivityRequest,
) (*pbsync.AddActivityResponse, error) {
	s.lastSyncedActivity.Store(time.Now().Unix())
	return s.store.AddUserActivity(ctx, req)
}

func (s *Server) AddUserProps(
	ctx context.Context,
	req *pbsync.AddUserPropsRequest,
) (*pbsync.AddActivityResponse, error) {
	s.lastSyncedActivity.Store(time.Now().Unix())
	return s.store.AddUserProps(ctx, req)
}

func (s *Server) AddColleagueActivity(
	ctx context.Context,
	req *pbsync.AddColleagueActivityRequest,
) (*pbsync.AddActivityResponse, error) {
	s.lastSyncedActivity.Store(time.Now().Unix())
	return s.store.AddColleagueActivity(ctx, req)
}

func (s *Server) AddColleagueProps(
	ctx context.Context,
	req *pbsync.AddColleaguePropsRequest,
) (*pbsync.AddActivityResponse, error) {
	s.lastSyncedActivity.Store(time.Now().Unix())
	return s.store.AddColleagueProps(ctx, req)
}

func (s *Server) AddJobTimeclock(
	ctx context.Context,
	req *pbsync.AddJobTimeclockRequest,
) (*pbsync.AddActivityResponse, error) {
	s.lastSyncedActivity.Store(time.Now().Unix())
	return s.store.AddJobTimeclock(ctx, req)
}

func (s *Server) AddDispatch(
	ctx context.Context,
	req *pbsync.AddDispatchRequest,
) (*pbsync.AddActivityResponse, error) {
	s.lastSyncedActivity.Store(time.Now().Unix())
	return s.store.AddDispatch(ctx, req)
}
