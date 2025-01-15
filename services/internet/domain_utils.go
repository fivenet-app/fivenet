package internet

import (
	"context"
	"errors"
	"strings"

	internet "github.com/fivenet-app/fivenet/gen/go/proto/resources/internet"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var tDomains = table.FivenetInternetDomains.AS("domain")

func cleanupDomainName(in string) string {
	return strings.ToLower(strings.Replace(
		strings.TrimSpace(in),
		"www.", "", 1),
	)
}

func (s *Server) getDomainByCondition(ctx context.Context, condition jet.BoolExpression) (*internet.Domain, error) {
	stmt := tDomains.
		SELECT(
			tDomains.ID,
			tDomains.CreatedAt,
			tDomains.UpdatedAt,
			tDomains.DeletedAt,
			tDomains.Name,
			tDomains.CreatorJob,
			tDomains.CreatorID,
		).
		FROM(tDomains).
		WHERE(condition).
		LIMIT(1)

	dest := &internet.Domain{}
	if err := stmt.QueryContext(ctx, s.db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.Id == 0 {
		return nil, nil
	}

	return dest, nil
}

func (s *Server) getDomainByName(ctx context.Context, name string) (*internet.Domain, error) {
	return s.getDomainByCondition(ctx,
		tDomains.Name.EQ(jet.String(cleanupDomainName(name))).
			AND(tDomains.DeletedAt.IS_NULL()),
	)
}
