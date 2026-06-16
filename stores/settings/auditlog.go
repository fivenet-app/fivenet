package settingsstore

import (
	"context"
	"errors"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	pbsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

const AuditLogPageSize = 30

var tAuditLog = table.FivenetAuditLog.AS("audit_entry")

func (s *Store) ViewAuditLog(
	ctx context.Context,
	opts ViewAuditLogOptions,
) (*pbsettings.ViewAuditLogResponse, error) {
	condition := mysql.Bool(true)
	if len(opts.UserIDs) > 0 {
		ids := make([]mysql.Expression, len(opts.UserIDs))
		for i := range opts.UserIDs {
			ids[i] = mysql.Int32(opts.UserIDs[i])
		}
		condition = condition.AND(tAuditLog.UserID.IN(ids...))
	}
	if opts.From != nil {
		condition = condition.AND(
			tAuditLog.CreatedAt.GT_EQ(mysql.DateTimeT(opts.From.AsTime())),
		)
	}
	if opts.To != nil {
		condition = condition.AND(tAuditLog.CreatedAt.LT_EQ(mysql.DateTimeT(opts.To.AsTime())))
	}
	if len(opts.Services) > 0 {
		svcs := make([]mysql.Expression, len(opts.Services))
		for i, svc := range opts.Services {
			svcs[i] = mysql.String(svc)
		}
		condition = condition.AND(tAuditLog.Service.IN(svcs...))
	}
	if len(opts.Methods) > 0 {
		methods := make([]mysql.Expression, len(opts.Methods))
		for i := range opts.Methods {
			methods[i] = mysql.String(opts.Methods[i])
		}
		condition = condition.AND(tAuditLog.Method.IN(methods...))
	}
	if len(opts.Actions) > 0 {
		actions := make([]mysql.Expression, len(opts.Actions))
		for i := range opts.Actions {
			actions[i] = mysql.Int32(int32(opts.Actions[i]))
		}
		condition = condition.AND(tAuditLog.Action.IN(actions...))
	}
	if len(opts.Results) > 0 {
		results := make([]mysql.Expression, len(opts.Results))
		for i := range opts.Results {
			results[i] = mysql.Int32(int32(opts.Results[i]))
		}
		condition = condition.AND(tAuditLog.Result.IN(results...))
	}
	if opts.Search != "" {
		condition = condition.AND(dbutils.MATCH(tAuditLog.Data, mysql.String(opts.Search)))
	}

	countStmt := tAuditLog.
		SELECT(mysql.COUNT(tAuditLog.ID).AS("data_count.total")).
		FROM(tAuditLog).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	pag, limit := opts.Pagination.GetResponseWithPageSize(count.Total, AuditLogPageSize)
	resp := &pbsettings.ViewAuditLogResponse{Pagination: pag}
	if count.Total <= 0 {
		return resp, nil
	}

	orderBys := s.auditLogSorter.Build(opts.Sort)

	tUser := table.FivenetUser.AS("user_short")
	stmt := tAuditLog.
		SELECT(
			tAuditLog.ID,
			tAuditLog.CreatedAt,
			tAuditLog.UserID,
			tAuditLog.UserJob,
			tAuditLog.TargetUserID,
			tAuditLog.Service,
			tAuditLog.Method,
			tAuditLog.Action,
			tAuditLog.Result,
			tAuditLog.Meta,
			tAuditLog.Data,
			tUser.ID,
			tUser.Identifier,
			tUser.Job,
			tUser.JobGrade,
			tUser.Firstname,
			tUser.Lastname,
			tUser.Dateofbirth,
		).
		FROM(
			tAuditLog.
				LEFT_JOIN(tUser,
					tUser.ID.EQ(tAuditLog.UserID),
				),
		).
		WHERE(condition).
		ORDER_BY(orderBys...).
		OFFSET(opts.Pagination.GetOffset()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Logs); err != nil {
		return nil, err
	}

	return resp, nil
}
