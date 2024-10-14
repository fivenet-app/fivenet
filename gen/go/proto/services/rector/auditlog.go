package rector

import (
	"context"
	"errors"
	"strings"

	database "github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	rector "github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	errorsrector "github.com/fivenet-app/fivenet/gen/go/proto/services/rector/errors"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

const AuditLogPageSize = 30

var (
	tAuditLog = table.FivenetAuditLog.AS("auditentry")
	tUser     = table.Users.AS("usershort")
)

func (s *Server) ViewAuditLog(ctx context.Context, req *ViewAuditLogRequest) (*ViewAuditLogResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if req.Pagination.Offset <= 0 {
		defer s.aud.Log(&model.FivenetAuditLog{
			Service: RectorService_ServiceDesc.ServiceName,
			Method:  "ViewAuditLog",
			UserID:  userInfo.UserId,
			UserJob: userInfo.Job,
			State:   int16(rector.EventType_EVENT_TYPE_VIEWED),
		}, req)
	}

	condition := jet.Bool(true)
	if !userInfo.SuperUser {
		condition = jet.AND(
			tAuditLog.UserJob.EQ(jet.String(userInfo.Job)).
				OR(tAuditLog.TargetUserJob.EQ(jet.String(userInfo.Job))),
		)
	}

	if len(req.UserIds) > 0 {
		ids := make([]jet.Expression, len(req.UserIds))
		for i := 0; i < len(req.UserIds); i++ {
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
	if req.Service != nil && *req.Service != "" {
		service := strings.ReplaceAll(*req.Service, "%", "")
		condition = condition.AND(tAuditLog.Service.LIKE(jet.String(service + "%")))
	}
	if req.Method != nil && *req.Method != "" {
		method := strings.ReplaceAll(*req.Method, "%", "")
		condition = condition.AND(tAuditLog.Method.LIKE(jet.String(method + "%")))
	}
	if req.Search != nil && *req.Search != "" {
		condition = condition.AND(jet.BoolExp(
			jet.Raw("MATCH(`data`) AGAINST ($search IN BOOLEAN MODE)",
				jet.RawArgs{"$search": *req.Search}),
		))
	}

	countStmt := tAuditLog.
		SELECT(
			jet.COUNT(tAuditLog.ID).AS("datacount.totalcount"),
		).
		FROM(tAuditLog).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, AuditLogPageSize)
	resp := &ViewAuditLogResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
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
		return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
	}

	resp.Pagination.Update(len(resp.Logs))

	return resp, nil
}
