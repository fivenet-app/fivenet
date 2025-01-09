package sync

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrSendDataDisabled = status.Error(codes.FailedPrecondition, "Sync API: SendData is disabled due to ESXCompat being enabled")

func (s *Server) SendData(ctx context.Context, req *SendDataRequest) (*SendDataResponse, error) {
	resp := &SendDataResponse{
		AffectedRows: 0,
	}

	if s.esxCompat {
		return nil, ErrSendDataDisabled
	}

	var err error
	switch d := req.Data.(type) {
	case *SendDataRequest_Jobs:
		if resp.AffectedRows, err = s.handleJobsData(ctx, d); err != nil {
			return nil, err
		}

	case *SendDataRequest_Licenses:
		if resp.AffectedRows, err = s.handleLicensesData(ctx, d); err != nil {
			return nil, err
		}

	case *SendDataRequest_Users:
		if resp.AffectedRows, err = s.handleUsersData(ctx, d); err != nil {
			return nil, err
		}

	case *SendDataRequest_Vehicles:
		if resp.AffectedRows, err = s.handleVehiclesData(ctx, d); err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func (s *Server) handleJobsData(ctx context.Context, data *SendDataRequest_Jobs) (int64, error) {
	// TODO

	return 0, nil
}

func (s *Server) handleLicensesData(ctx context.Context, data *SendDataRequest_Licenses) (int64, error) {
	// TODO

	return 0, nil
}

func (s *Server) handleUsersData(ctx context.Context, data *SendDataRequest_Users) (int64, error) {
	// TODO

	return 0, nil
}

func (s *Server) handleVehiclesData(ctx context.Context, data *SendDataRequest_Vehicles) (int64, error) {
	// TODO

	return 0, nil
}
