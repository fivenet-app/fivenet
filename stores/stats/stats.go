package statsstore

import (
	"context"
	"errors"
	"fmt"

	resstats "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/stats"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/multierr"
)

type Stats = map[string]*resstats.Stat

type statQuery struct {
	key  string
	stmt mysql.Statement
}

func (s *Store) LoadPublicStats(ctx context.Context) (Stats, error) {
	tAccounts := table.FivenetAccounts
	tDocuments := table.FivenetDocuments
	tDispatches := table.FivenetCentrumDispatches
	tCitizenActivity := table.FivenetUserActivity
	tJobUserTimeclock := table.FivenetJobTimeclock
	tUsers := table.FivenetUser

	queries := []statQuery{
		{
			key: "users_registered",
			stmt: tAccounts.SELECT(mysql.COUNT(tAccounts.ID).AS("value")).
				WHERE(tAccounts.DeletedAt.IS_NULL()),
		},
		{
			key: "documents_created",
			stmt: tDocuments.SELECT(mysql.COUNT(tDocuments.ID).AS("value")).
				WHERE(tDocuments.DeletedAt.IS_NULL()),
		},
		{
			key:  "dispatches_created",
			stmt: tDispatches.SELECT(mysql.MAX(tDispatches.ID).AS("value")),
		},
		{
			key:  "citizen_activity",
			stmt: tCitizenActivity.SELECT(mysql.COUNT(tCitizenActivity.ID).AS("value")),
		},
		{
			key: "timeclock_tracked",
			stmt: tJobUserTimeclock.SELECT(
				mysql.CAST(mysql.SUM(tJobUserTimeclock.SpentTime)).AS_SIGNED().AS("value"),
			),
		},
		{
			key:  "citizens_total",
			stmt: tUsers.SELECT(mysql.COUNT(tUsers.ID).AS("value")),
		},
	}

	data := make(Stats, len(queries))
	var errs error
	zero := int32(0)

	for _, query := range queries {
		dest := &struct {
			Value *int32 `alias:"value"`
		}{
			Value: &zero,
		}

		if err := query.stmt.QueryContext(ctx, s.db, dest); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				errs = multierr.Append(
					errs,
					fmt.Errorf("error during %q stats query. %w", query.key, err),
				)
				continue
			}
		}

		data[query.key] = &resstats.Stat{Value: dest.Value}
	}

	return data, errs
}
