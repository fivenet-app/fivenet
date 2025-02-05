package internet

import (
	"context"
	"errors"
	"path"
	"strings"

	internet "github.com/fivenet-app/fivenet/gen/go/proto/resources/internet"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var tPage = table.FivenetInternetPages.AS("page")

func cleanupPath(in string) string {
	in = strings.ToLower(
		strings.TrimSpace(path.Clean(in)),
	)
	if in == "" || in == "." {
		return "/"
	}
	return in
}

func (s *Server) getPageByDomainAndPath(ctx context.Context, domainId uint64, path string) (*internet.Page, error) {
	path = cleanupPath(path)

	return s.getPageByCondition(ctx,
		tPage.DomainID.EQ(jet.Uint64(domainId)).
			AND(tPage.Path.EQ(jet.String(path))),
	)
}

func (s *Server) getPageByCondition(ctx context.Context, condition jet.BoolExpression) (*internet.Page, error) {
	stmt := tPage.
		SELECT(
			tPage.ID,
			tPage.CreatedAt,
			tPage.UpdatedAt,
			tPage.DeletedAt,
			tPage.DomainID,
			tPage.Path,
			tPage.Title,
			tPage.Description,
			tPage.Data,
			tPage.CreatorJob,
			tPage.CreatorID,
		).
		FROM(tPage).
		WHERE(condition).
		LIMIT(1)

	dest := &internet.Page{}
	if err := stmt.QueryContext(ctx, s.db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.Id == 0 || dest.DomainId == 0 {
		return nil, nil
	}

	return dest, nil
}
