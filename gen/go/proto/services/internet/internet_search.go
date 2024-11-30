package internet

import (
	"context"

	internet "github.com/fivenet-app/fivenet/gen/go/proto/resources/internet"
)

func (s *Server) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	// TODO

	resp := &SearchResponse{
		Results: []*internet.SearchResult{},
	}
	return resp, nil
}
