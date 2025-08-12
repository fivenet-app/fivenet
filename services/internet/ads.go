package internet

import (
	"context"

	internet "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/internet"
	pbinternet "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/internet"
)

func (s *Server) GetAds(
	ctx context.Context,
	req *pbinternet.GetAdsRequest,
) (*pbinternet.GetAdsResponse, error) {
	// TODO

	resp := &pbinternet.GetAdsResponse{
		Ads: []*internet.Ad{},
	}
	return resp, nil
}
