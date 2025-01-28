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

func cleanupDomainName(in string) string {
	return strings.ToLower(strings.Replace(
		strings.TrimSpace(in),
		"www.", "", 1),
	)
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
			tDomains.ApproverJob,
			tDomains.ApproverID,
			tDomains.CreatorJob,
			tDomains.CreatorID,
		).
		FROM(
			tDomains.
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

func (s *Server) getDomainByName(ctx context.Context, tx *sql.DB, name string) (*internet.Domain, error) {
	return s.getDomainByCondition(ctx, tx,
		tDomains.Name.EQ(jet.String(cleanupDomainName(name))).
			AND(tDomains.DeletedAt.IS_NULL()),
	)
}

func (s *Server) getDomainById(ctx context.Context, tx *sql.DB, id uint64) (*internet.Domain, error) {
	return s.getDomainByCondition(ctx, tx,
		tDomains.ID.EQ(jet.Uint64(id)).
			AND(tDomains.DeletedAt.IS_NULL()),
	)
}
