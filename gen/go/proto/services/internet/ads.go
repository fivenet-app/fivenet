package internet

import (
	"context"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/filestore"
	internet "github.com/fivenet-app/fivenet/gen/go/proto/resources/internet"
)

func (s *Server) GetAds(ctx context.Context, req *GetAdsRequest) (*GetAdsResponse, error) {
	// TODO

	cafeLogoUrl := "/api/filestore/joblogos/cafe.png"
	nagataLogoUrl := "/api/filestore/joblogos/mechanic.png"
	fivenetLogoUrl := "/api/filestore/joblogos/5net.png"
	resp := &GetAdsResponse{
		Ads: []*internet.Ad{
			{
				Id:          1,
				Title:       "Katzen Cafe",
				Description: "Cafe, Heiße Schokolade, Törtchen und mehr im Katzen Cafe. Suck on our Boba!",
				Image: &filestore.File{
					Url: &cafeLogoUrl,
				},
			},
			{
				Id:          2,
				Title:       "Nagata Performance",
				Description: "Heiße Karren, leuchtende Lacke, Tuning vom Feinsten, und mehr bei Nagata Performance. Eat our Exhaust!",
				Image: &filestore.File{
					Url: &nagataLogoUrl,
				},
			},
			{
				Id:          3,
				Title:       "FiveNet",
				Description: "Tablet? Tablet!",
				Image: &filestore.File{
					Url: &fivenetLogoUrl,
				},
			},
		},
	}
	return resp, nil
}
