package completor

import (
	context "context"
	"errors"
	"slices"

	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	pbcompletor "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/completor"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorscompletor "github.com/fivenet-app/fivenet/v2026/services/completor/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) CompleteCitizens(
	ctx context.Context,
	req *pbcompletor.CompleteCitizensRequest,
) (*pbcompletor.CompleteCitizensResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tUsers := table.FivenetUser.AS("user_short")

	orderBys := []mysql.OrderByClause{}
	condition := s.customDB.Conditions.User.GetFilter(tUsers.Alias())

	currentJob := req.CurrentJob != nil && req.GetCurrentJob()

	if currentJob {
		condition = condition.AND(
			tUsers.Job.EQ(mysql.String(userInfo.GetJob())),
		)
	}

	search := dbutils.PrepareForLikeSearch(req.GetSearch())
	if search != "" {
		condition = condition.AND(
			mysql.CONCAT(tUsers.Firstname, mysql.String(" "), tUsers.Lastname).
				LIKE(mysql.String(search)),
		)
	}

	if len(req.GetUserIds()) > 0 {
		userIds := []mysql.Expression{}
		for _, v := range req.GetUserIds() {
			userIds = append(userIds, mysql.Int32(v))
		}

		if req.GetUserIdsOnly() {
			condition = condition.OR(tUsers.ID.IN(userIds...))
		}

		// Make sure to sort by the user IDs if provided
		orderBys = append(orderBys, tUsers.ID.IN(userIds...).DESC())
	}

	orderBys = append(orderBys, tUsers.Lastname.ASC())

	columns := mysql.ProjectionList{
		tUsers.ID,
		tUsers.Firstname,
		tUsers.Lastname,
		tUsers.Dateofbirth,
	}
	if currentJob {
		columns = append(columns, tUsers.Job, tUsers.JobGrade)
	}

	stmt := tUsers.
		SELECT(
			columns[0],
			columns[1:]...,
		).
		OPTIMIZER_HINTS(mysql.OptimizerHint("idx_users_firstname_lastname_fulltext")).
		FROM(tUsers).
		WHERE(condition).
		ORDER_BY(orderBys...).
		LIMIT(15)

	var dest []*usershort.UserShort
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscompletor.ErrFailedSearch)
		}
	}

	if req.OnDuty != nil && req.GetOnDuty() {
		dest = slices.DeleteFunc(dest, func(us *usershort.UserShort) bool {
			return !s.tracker.IsUserOnDuty(us.GetUserId())
		})
	}

	if currentJob {
		jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
		for i := range dest {
			jobInfoFn(dest[i])
		}
	}

	return &pbcompletor.CompleteCitizensResponse{
		Users: dest,
	}, nil
}
