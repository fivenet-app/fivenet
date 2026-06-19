package sync

import (
	"context"
	"time"

	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
)

func (s *Server) AddUserUpdate(
	ctx context.Context,
	req *pbsync.AddUserUpdateRequest,
) (*pbsync.AddActivityResponse, error) {
	s.lastSyncedActivity.Store(time.Now().Unix())
	return s.store.AddUserUpdate(ctx, req)
}

func (s *Server) SendUsers(
	ctx context.Context,
	req *pbsync.SendUsersRequest,
) (*pbsync.SendDataResponse, error) {
	s.lastSyncedData.Store(time.Now().Unix())
	rowsAffected, err := s.store.SendUsers(ctx, req.GetUsers())
	if err != nil {
		return nil, err
	}

	return &pbsync.SendDataResponse{RowsAffected: rowsAffected}, nil
}

func (s *Server) DeleteUsers(
	ctx context.Context,
	req *pbsync.DeleteUsersRequest,
) (*pbsync.DeleteDataResponse, error) {
	return s.store.DeleteUsers(ctx, req.GetUserIds())
}
