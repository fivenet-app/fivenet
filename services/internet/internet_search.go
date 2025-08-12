package internet

import (
	"context"
	"errors"

	internet "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/internet"
	pbinternet "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/internet"
	errorsinternet "github.com/fivenet-app/fivenet/v2025/services/internet/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) Search(
	ctx context.Context,
	req *pbinternet.SearchRequest,
) (*pbinternet.SearchResponse, error) {
	resp := &pbinternet.SearchResponse{
		Results: []*internet.SearchResult{},
	}

	if req.GetSearch() == "" {
		return resp, nil
	}

	tPage := tPage.AS("search_result")

	condition := jet.OR(
		// Search title and description fields
		jet.BoolExp(
			jet.Raw("MATCH(`title`) AGAINST ($search IN BOOLEAN MODE)",
				jet.RawArgs{"$search": req.GetSearch()}),
		),
		jet.BoolExp(
			jet.Raw("MATCH(`description`) AGAINST ($search IN BOOLEAN MODE)",
				jet.RawArgs{"$search": req.GetSearch()}),
		),
	)

	if req.DomainId != nil && req.GetDomainId() > 0 {
		condition = condition.AND(
			tPage.DomainID.EQ(jet.Uint64(req.GetDomainId())),
		)
	}

	stmt := tPage.
		SELECT(
			tPage.ID,
			tPage.Title,
			tPage.Description,
			tPage.DomainID,
			tPage.Path,
			tDomains.ID,
			tDomains.Name,
			tTLDs.ID,
			tTLDs.Name,
		).
		FROM(
			tPage.
				INNER_JOIN(tDomains,
					tDomains.ID.EQ(tPage.DomainID),
				).
				INNER_JOIN(tTLDs,
					tTLDs.ID.EQ(tDomains.TldID),
				),
		).
		WHERE(condition).
		LIMIT(15)

	if err := stmt.QueryContext(ctx, s.db, &resp.Results); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errorsinternet.ErrFailedSearch
		}
	}

	return resp, nil
}
