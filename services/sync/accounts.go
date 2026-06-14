package sync

import (
	"context"
	"time"

	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
)

func (s *Server) RegisterAccount(
	ctx context.Context,
	req *pbsync.RegisterAccountRequest,
) (*pbsync.RegisterAccountResponse, error) {
	return s.store.RegisterAccount(ctx, req)
}

func (s *Server) TransferAccount(
	ctx context.Context,
	req *pbsync.TransferAccountRequest,
) (*pbsync.TransferAccountResponse, error) {
	return s.store.TransferAccount(ctx, req)
}

func (s *Server) AddAccountUpdate(
	ctx context.Context,
	req *pbsync.AddAccountUpdateRequest,
) (*pbsync.AddActivityResponse, error) {
	s.lastSyncedActivity.Store(time.Now().Unix())
	return s.store.AddAccountUpdate(ctx, req)
}
