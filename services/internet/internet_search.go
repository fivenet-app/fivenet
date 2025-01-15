package internet

import (
	"context"

	internet "github.com/fivenet-app/fivenet/gen/go/proto/resources/internet"
	pbinternet "github.com/fivenet-app/fivenet/gen/go/proto/services/internet"
)

func (s *Server) Search(ctx context.Context, req *pbinternet.SearchRequest) (*pbinternet.SearchResponse, error) {
	// TODO

	resp := &pbinternet.SearchResponse{
		Results: []*internet.SearchResult{},
	}
	return resp, nil
}
