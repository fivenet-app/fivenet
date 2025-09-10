package citizens

import (
	context "context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	users "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	pbcitizens "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/citizens"
	permscitizens "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorscitizens "github.com/fivenet-app/fivenet/v2025/services/citizens/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) ListUserActivity(
	ctx context.Context,
	req *pbcitizens.ListUserActivityRequest,
) (*pbcitizens.ListUserActivityResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.citizens.user_id", req.GetUserId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	resp := &pbcitizens.ListUserActivityResponse{
		Activity: []*users.UserActivity{},
	}

	// User can't see their own activities, unless they have "Own" perm attribute, or are a superuser
	fields, err := s.ps.AttrStringList(
		userInfo,
		permscitizens.CitizensServicePerm,
		permscitizens.CitizensServiceListUserActivityPerm,
		permscitizens.CitizensServiceListUserActivityFieldsPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if userInfo.GetUserId() == req.GetUserId() {
		// If isn't superuser or doesn't have 'Own' activity feed access
		if !userInfo.GetSuperuser() && !fields.Contains("Own") {
			return resp, nil
		}
	}

	tUserActivity := table.FivenetUserActivity.AS("user_activity")

	condition := tUserActivity.TargetUserID.EQ(mysql.Int32(req.GetUserId()))

	if len(req.GetTypes()) > 0 {
		types := []mysql.Expression{}
		for _, t := range req.GetTypes() {
			types = append(types, mysql.Int32(int32(*t.Enum())))
		}

		condition = condition.AND(tUserActivity.Type.IN(types...))
	}

	// Get total count of values
	countStmt := tUserActivity.
		SELECT(
			mysql.COUNT(tUserActivity.ID).AS("data_count.total"),
		).
		FROM(tUserActivity).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count.Total, 20)
	resp.Pagination = pag
	if count.Total <= 0 {
		return resp, nil
	}

	tUTarget := tables.User().AS("target_user")
	tUSource := tUTarget.AS("source_user")

	// Convert proto sort to db sorting
	orderBys := []mysql.OrderByClause{}
	if req.GetSort() != nil {
		var column mysql.Column
		switch req.GetSort().GetColumn() {
		case "createdAt":
			fallthrough
		default:
			column = tUserActivity.CreatedAt
		}

		if req.GetSort().GetDirection() == database.AscSortDirection {
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
		OFFSET(req.GetPagination().GetOffset()).
		ORDER_BY(orderBys...).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Activity); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetActivity() {
		if resp.GetActivity()[i].GetSourceUser() != nil {
			jobInfoFn(resp.GetActivity()[i].GetSourceUser())
		}
		if resp.GetActivity()[i].GetTargetUser() != nil {
			jobInfoFn(resp.GetActivity()[i].GetTargetUser())
		}
	}

	resp.GetPagination().Update(len(resp.GetActivity()))

	return resp, nil
}
