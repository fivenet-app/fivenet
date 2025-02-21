package internet

import (
	"context"
	"errors"

	internet "github.com/fivenet-app/fivenet/gen/go/proto/resources/internet"
	pbinternet "github.com/fivenet-app/fivenet/gen/go/proto/services/internet"
	errorsinternet "github.com/fivenet-app/fivenet/services/internet/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) Search(ctx context.Context, req *pbinternet.SearchRequest) (*pbinternet.SearchResponse, error) {
	condition := jet.Bool(true)

	resp := &pbinternet.SearchResponse{
		Results: []*internet.SearchResult{},
	}

	if req.Search == "" {
		return resp, nil
	}

	tPage := tPage.AS("search_result")

	condition = jet.OR(
		// Search title and description fields
		jet.BoolExp(
			jet.Raw("MATCH(`title`) AGAINST ($search IN BOOLEAN MODE)",
				jet.RawArgs{"$search": req.Search}),
		),
		jet.BoolExp(
			jet.Raw("MATCH(`description`) AGAINST ($search IN BOOLEAN MODE)",
				jet.RawArgs{"$search": req.Search}),
		),
	)

	if req.DomainId != nil && *req.DomainId > 0 {
		condition = condition.AND(
			tPage.DomainID.EQ(jet.Uint64(*req.DomainId)),
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
		).
		FROM(
			tPage.
				INNER_JOIN(tDomains,
					tDomains.ID.EQ(tPage.DomainID),
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
