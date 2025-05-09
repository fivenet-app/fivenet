package citizenstore

import (
	context "context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	users "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	pbcitizenstore "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/citizenstore"
	permscitizenstore "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/citizenstore/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorscitizenstore "github.com/fivenet-app/fivenet/v2025/services/citizenstore/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *Server) ListUserActivity(ctx context.Context, req *pbcitizenstore.ListUserActivityRequest) (*pbcitizenstore.ListUserActivityResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.citizenstore.user_id", int64(req.UserId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	resp := &pbcitizenstore.ListUserActivityResponse{
		Activity: []*users.UserActivity{},
	}

	// User can't see their own activities, unless they have "Own" perm attribute, or are a superuser
	fields, err := s.ps.AttrStringList(userInfo, permscitizenstore.CitizenStoreServicePerm, permscitizenstore.CitizenStoreServiceListUserActivityPerm, permscitizenstore.CitizenStoreServiceListUserActivityFieldsPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizenstore.ErrFailedQuery)
	}

	if userInfo.UserId == req.UserId {
		// If isn't superuser or doesn't have 'Own' activity feed access
		if !userInfo.SuperUser && !fields.Contains("Own") {
			return resp, nil
		}
	}

	tUserActivity := table.FivenetUserActivity

	condition := tUserActivity.TargetUserID.EQ(jet.Int32(req.UserId))

	if len(req.Types) > 0 {
		types := []jet.Expression{}
		for _, t := range req.Types {
			types = append(types, jet.Int16(int16(*t.Enum())))
		}

		condition = condition.AND(tUserActivity.Type.IN(types...))
	}

	// Get total count of values
	countStmt := tUserActivity.
		SELECT(
			jet.COUNT(tUserActivity.ID).AS("datacount.totalcount"),
		).
		FROM(tUserActivity).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscitizenstore.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, 16)
	resp.Pagination = pag
	if count.TotalCount <= 0 {
		return resp, nil
	}

	tUTarget := tables.Users().AS("target_user")
	tUSource := tUTarget.AS("source_user")

	// Convert proto sort to db sorting
	orderBys := []jet.OrderByClause{}
	if req.Sort != nil {
		var column jet.Column
		switch req.Sort.Column {
		case "createdAt":
			fallthrough
		default:
			column = tUserActivity.CreatedAt
		}

		if req.Sort.Direction == database.AscSortDirection {
			orderBys = append(orderBys, column.ASC())
		} else {
			orderBys = append(orderBys, column.DESC())
		}
	} else {
		orderBys = append(orderBys, tUserActivity.CreatedAt.DESC())
	}

	stmt := tUserActivity.
		SELECT(
			tUserActivity.ID,
			tUserActivity.CreatedAt,
			tUserActivity.SourceUserID,
			tUserActivity.TargetUserID,
			tUserActivity.Type,
			tUserActivity.Reason,
			tUserActivity.Data,
			tUTarget.ID,
			tUTarget.Job,
			tUTarget.JobGrade,
			tUTarget.Firstname,
			tUTarget.Lastname,
			tUSource.ID,
			tUSource.Job,
			tUSource.JobGrade,
			tUSource.Firstname,
			tUSource.Lastname,
		).
		FROM(
			tUserActivity.
				INNER_JOIN(tUTarget,
					tUTarget.ID.EQ(tUserActivity.TargetUserID),
				).
				LEFT_JOIN(tUSource,
					tUSource.ID.EQ(tUserActivity.SourceUserID),
				),
		).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		ORDER_BY(orderBys...).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Activity); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscitizenstore.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.Activity {
		if resp.Activity[i].SourceUser != nil {
			jobInfoFn(resp.Activity[i].SourceUser)
		}
		if resp.Activity[i].TargetUser != nil {
			jobInfoFn(resp.Activity[i].TargetUser)
		}
	}

	resp.Pagination.Update(len(resp.Activity))

	return resp, nil
}
