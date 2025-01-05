package sync

import "context"

func (s *Server) SyncData(ctx context.Context, req *SyncDataRequest) (*SyncDataResponse, error) {
	// TODO handle sync data request

	return &SyncDataResponse{
		AffectedRows: 0,
	}, nil
}
