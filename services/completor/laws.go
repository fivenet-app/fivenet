package completor

import (
	context "context"

	pbcompletor "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/completor"
)

func (s *Server) ListLawBooks(
	ctx context.Context,
	req *pbcompletor.ListLawBooksRequest,
) (*pbcompletor.ListLawBooksResponse, error) {
	return &pbcompletor.ListLawBooksResponse{
		Books: s.laws.GetLawBooks(),
	}, nil
}
