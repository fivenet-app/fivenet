package sync

import (
	"context"
	"time"

	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
)

func (s *Server) SendData(
	ctx context.Context,
	req *pbsync.SendDataRequest,
) (*pbsync.SendDataResponse, error) {
	s.lastSyncedData.Store(time.Now().Unix())
	return s.store.SendData(ctx, req)
}

func (s *Server) SendLicenses(
	ctx context.Context,
	req *pbsync.SendLicensesRequest,
) (*pbsync.SendDataResponse, error) {
	s.lastSyncedData.Store(time.Now().Unix())
	return s.store.SendLicenses(ctx, req)
}

func (s *Server) SendAccounts(
	ctx context.Context,
	req *pbsync.SendAccountsRequest,
) (*pbsync.SendDataResponse, error) {
	s.lastSyncedData.Store(time.Now().Unix())
	return s.store.SendAccounts(ctx, req)
}

func (s *Server) SetLastCharID(
	ctx context.Context,
	req *pbsync.SetLastCharIDRequest,
) (*pbsync.SendDataResponse, error) {
	s.lastSyncedData.Store(time.Now().Unix())
	return s.store.SetLastCharID(ctx, req)
}

func (s *Server) DeleteData(
	ctx context.Context,
	req *pbsync.DeleteDataRequest,
) (*pbsync.DeleteDataResponse, error) {
	return s.store.DeleteData(ctx, req)
}
