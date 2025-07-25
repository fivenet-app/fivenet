package settings

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	pbsettings "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorssettings "github.com/fivenet-app/fivenet/v2025/services/settings/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

const AuditLogPageSize = 30

var tAuditLog = table.FivenetAuditLog.AS("audit_entry")

func (s *Server) ViewAuditLog(ctx context.Context, req *pbsettings.ViewAuditLogRequest) (*pbsettings.ViewAuditLogResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	defer s.aud.Log(&audit.AuditEntry{
		Service: pbsettings.SettingsService_ServiceDesc.ServiceName,
		Method:  "ViewAuditLog",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_VIEWED,
	}, req)

	condition := jet.Bool(true)
	if !userInfo.Superuser {
		condition = jet.AND(
			tAuditLog.UserJob.EQ(jet.String(userInfo.Job)).
				OR(tAuditLog.TargetUserJob.EQ(jet.String(userInfo.Job))),
		)
	}

	if len(req.UserIds) > 0 {
		ids := make([]jet.Expression, len(req.UserIds))
		for i := range req.UserIds {
			ids[i] = jet.Int32(req.UserIds[i])
		}
		condition = condition.AND(tAuditLog.UserID.IN(ids...))
	}
	if req.From != nil {
		condition = condition.AND(tAuditLog.CreatedAt.GT_EQ(
			jet.TimestampT(req.From.AsTime()),
		))
	}
	if req.To != nil {
		condition = condition.AND(tAuditLog.CreatedAt.LT_EQ(
			jet.TimestampT(req.To.AsTime()),
		))
	}
	if len(req.Services) > 0 {
		svcs := make([]jet.Expression, len(req.Services))
		for i := range req.Services {
			svcs[i] = jet.String(req.Services[i])
		}
		condition = condition.AND(tAuditLog.Service.IN(svcs...))
	}
	if len(req.Methods) > 0 {
		methods := make([]jet.Expression, len(req.Methods))
		for i := range req.Methods {
			methods[i] = jet.String(req.Methods[i])
		}
		condition = condition.AND(tAuditLog.Method.IN(methods...))
	}
	if len(req.States) > 0 {
		states := make([]jet.Expression, len(req.States))
		for i := range req.States {
			states[i] = jet.Int32(int32(req.States[i]))
		}
		condition = condition.AND(tAuditLog.State.IN(states...))
	}
	if req.Search != nil && *req.Search != "" {
		condition = condition.AND(jet.BoolExp(
			jet.Raw("MATCH(`data`) AGAINST ($search IN BOOLEAN MODE)",
				jet.RawArgs{"$search": *req.Search}),
		))
	}

	countStmt := tAuditLog.
		SELECT(
			jet.COUNT(tAuditLog.ID).AS("data_count.total"),
		).
		FROM(tAuditLog).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.Total, AuditLogPageSize)
	resp := &pbsettings.ViewAuditLogResponse{
		Pagination: pag,
	}
	if count.Total <= 0 {
		return resp, nil
	}

	// Convert proto sort to db sorting
	orderBys := []jet.OrderByClause{}
	if req.Sort != nil {
		var column jet.Column
		switch req.Sort.Column {
		case "service":
			column = tAuditLog.Service
		case "state":
			column = tAuditLog.State
		case "createdAt":
			fallthrough
		default:
			column = tAuditLog.CreatedAt
		}

		if req.Sort.Direction == database.AscSortDirection {
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
		OFFSET(req.Pagination.Offset).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Logs); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	resp.Pagination.Update(len(resp.Logs))

	return resp, nil
}
