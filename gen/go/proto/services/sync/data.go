package sync

import "context"

func (s *Server) SyncData(ctx context.Context, req *SyncDataRequest) (*SyncDataResponse, error) {
	resp := &SyncDataResponse{
		AffectedRows: 0,
	}

	switch d := req.Data.(type) {
	case *SyncDataRequest_Users:
		_ = d
		// TODO

	case *SyncDataRequest_Jobs:
		// TODO

	case *SyncDataRequest_Licenses:
		// TODO

	case *SyncDataRequest_UserLicenses:
		// TODO

	case *SyncDataRequest_Vehicles:
		// TODO
	}

	return resp, nil
}
