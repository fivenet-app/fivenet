package internet

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	internet "github.com/fivenet-app/fivenet/gen/go/proto/resources/internet"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils/tables"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var tDomains = table.FivenetInternetDomains.AS("domain")

func cleanAndSplitDomain(in string) (string, string) {
	full := strings.ToLower(strings.Replace(
		strings.TrimSpace(in),
		"www.", "", 1),
	)

	idx := strings.LastIndex(full, ".")
	tld := full[idx+1:]
	domain := full[:idx]
	return domain, tld
}

func (s *Server) getDomainByCondition(ctx context.Context, tx *sql.DB, condition jet.BoolExpression) (*internet.Domain, error) {
	tCreator := tables.Users().AS("creator")

	stmt := tDomains.
		SELECT(
			tDomains.ID,
			tDomains.CreatedAt,
			tDomains.UpdatedAt,
			tDomains.DeletedAt,
			tDomains.ExpiresAt,
			tDomains.TldID,
			tDomains.Name,
			tDomains.Active,
			tDomains.TransferCode,
			tDomains.ApproverJob,
			tDomains.ApproverID,
			tDomains.CreatorJob,
			tDomains.CreatorID,
			tTLDs.ID,
			tTLDs.Name,
			tTLDs.Internal,
		).
		FROM(
			tDomains.
				INNER_JOIN(tTLDs,
					tTLDs.ID.EQ(tDomains.TldID),
				).
				LEFT_JOIN(tCreator,
					tDomains.CreatorID.EQ(tCreator.ID),
				),
		).
		WHERE(condition).
		LIMIT(1)

	dest := &internet.Domain{}
	if err := stmt.QueryContext(ctx, tx, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.Id == 0 {
		return nil, nil
	}

	return dest, nil
}

func (s *Server) getDomainByTLDAndName(ctx context.Context, tx *sql.DB, tldId uint64, name string) (*internet.Domain, error) {
	domain, _ := cleanAndSplitDomain(name)
	return s.getDomainByCondition(ctx, tx,
		jet.AND(
			tTLDs.ID.EQ(jet.Uint64(tldId)),
			tDomains.Name.EQ(jet.String(domain)),
			tDomains.DeletedAt.IS_NULL(),
		),
	)
}

func (s *Server) getDomainByName(ctx context.Context, tx *sql.DB, name string) (*internet.Domain, error) {
	domain, tld := cleanAndSplitDomain(name)
	return s.getDomainByCondition(ctx, tx,
		jet.AND(
			tTLDs.Name.EQ(jet.String(tld)),
			tDomains.Name.EQ(jet.String(domain)),
			tDomains.DeletedAt.IS_NULL(),
		),
	)
}

func (s *Server) getDomainById(ctx context.Context, tx *sql.DB, id uint64) (*internet.Domain, error) {
	return s.getDomainByCondition(ctx, tx,
		tDomains.ID.EQ(jet.Uint64(id)).
			AND(tDomains.DeletedAt.IS_NULL()),
	)
}
