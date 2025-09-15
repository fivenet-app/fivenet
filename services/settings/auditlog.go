package settings

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	pbsettings "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorssettings "github.com/fivenet-app/fivenet/v2025/services/settings/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

const AuditLogPageSize = 30

var tAuditLog = table.FivenetAuditLog.AS("audit_entry")

func (s *Server) ViewAuditLog(
	ctx context.Context,
	req *pbsettings.ViewAuditLogRequest,
) (*pbsettings.ViewAuditLogResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	defer s.aud.Log(&audit.AuditEntry{
		Service: pbsettings.SettingsService_ServiceDesc.ServiceName,
		Method:  "ViewAuditLog",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_VIEWED,
	}, req)

	condition := mysql.Bool(true)
	if !userInfo.GetSuperuser() {
		condition = mysql.AND(
			tAuditLog.UserJob.EQ(mysql.String(userInfo.GetJob())).
				OR(tAuditLog.TargetUserJob.EQ(mysql.String(userInfo.GetJob()))),
		)
	}

	if len(req.GetUserIds()) > 0 {
		ids := make([]mysql.Expression, len(req.GetUserIds()))
		for i := range req.GetUserIds() {
			ids[i] = mysql.Int32(req.GetUserIds()[i])
		}
		condition = condition.AND(tAuditLog.UserID.IN(ids...))
	}
	if req.GetFrom() != nil {
		condition = condition.AND(tAuditLog.CreatedAt.GT_EQ(
			mysql.TimestampT(req.GetFrom().AsTime()),
		))
	}
	if req.GetTo() != nil {
		condition = condition.AND(tAuditLog.CreatedAt.LT_EQ(
			mysql.TimestampT(req.GetTo().AsTime()),
		))
	}
	if len(req.GetServices()) > 0 {
		svcs := make([]mysql.Expression, len(req.GetServices()))
		for i := range req.GetServices() {
			svcs[i] = mysql.String(req.GetServices()[i])
		}
		condition = condition.AND(tAuditLog.Service.IN(svcs...))
	}
	if len(req.GetMethods()) > 0 {
		methods := make([]mysql.Expression, len(req.GetMethods()))
		for i := range req.GetMethods() {
			methods[i] = mysql.String(req.GetMethods()[i])
		}
		condition = condition.AND(tAuditLog.Method.IN(methods...))
	}
	if len(req.GetStates()) > 0 {
		states := make([]mysql.Expression, len(req.GetStates()))
		for i := range req.GetStates() {
			states[i] = mysql.Int32(int32(req.GetStates()[i]))
		}
		condition = condition.AND(tAuditLog.State.IN(states...))
	}
	if req.Search != nil && req.GetSearch() != "" {
		condition = condition.AND(
			dbutils.MATCH(tAuditLog.Data, mysql.String(req.GetSearch())),
		)
	}

	countStmt := tAuditLog.
		SELECT(
			mysql.COUNT(tAuditLog.ID).AS("data_count.total"),
		).
		FROM(tAuditLog).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count.Total, AuditLogPageSize)
	resp := &pbsettings.ViewAuditLogResponse{
		Pagination: pag,
	}
	if count.Total <= 0 {
		return resp, nil
	}

	// Convert proto sort to db sorting
	orderBys := []mysql.OrderByClause{}
	if req.GetSort() != nil {
		var column mysql.Column
		switch req.GetSort().GetColumn() {
		case "service":
			column = tAuditLog.Service
		case "state":
			column = tAuditLog.State
		case "createdAt":
			fallthrough
		default:
			column = tAuditLog.CreatedAt
		}

		if req.GetSort().GetDirection() == database.AscSortDirection {
			orderBys = append(orderBys, column.ASC())
		} else {
			orderBys = append(orderBys, column.DESC())
		}
	} else {
		orderBys = append(orderBys, tAuditLog.CreatedAt.DESC())
	}

	tUser := tables.User().AS("user_short")

	stmt := tAuditLog.
		SELECT(
			tAuditLog.ID,
			tAuditLog.CreatedAt,
			tAuditLog.UserID,
			tAuditLog.UserJob,
			tAuditLog.TargetUserID,
			tAuditLog.Service,
			tAuditLog.Method,
			tAuditLog.State,
			tAuditLog.Data,
			tUser.ID,
			tUser.Identifier,
			tUser.Job,
			tUser.JobGrade,
			tUser.Firstname,
			tUser.Lastname,
		).
		FROM(
			tAuditLog.
				LEFT_JOIN(tUser,
					tUser.ID.EQ(tAuditLog.UserID),
				),
		).
		WHERE(condition).
		ORDER_BY(orderBys...).
		OFFSET(req.GetPagination().GetOffset()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Logs); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	resp.GetPagination().Update(len(resp.GetLogs()))

	return resp, nil
}
