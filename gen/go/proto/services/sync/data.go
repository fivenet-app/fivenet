package sync

import "context"

func (s *Server) SendData(ctx context.Context, req *SendDataRequest) (*SendDataResponse, error) {
	resp := &SendDataResponse{
		AffectedRows: 0,
	}

	switch d := req.Data.(type) {
	case *SendDataRequest_Users:
		_ = d
		// TODO

	case *SendDataRequest_Jobs:
		// TODO

	case *SendDataRequest_Licenses:
		// TODO

	case *SendDataRequest_Vehicles:
		// TODO
	}

	return resp, nil
}
