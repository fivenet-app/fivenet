package completor

import (
	context "context"
	"errors"
	"slices"

	citizenslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/citizens/labels"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	pbcompletor "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/completor"
	permscompletor "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/completor/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorscompletor "github.com/fivenet-app/fivenet/v2026/services/completor/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var tCitizensLabelsJob = table.FivenetUserLabelsJob.AS("label")

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

func (s *Server) CompleteCitizenLabels(
	ctx context.Context,
	req *pbcompletor.CompleteCitizenLabelsRequest,
) (*pbcompletor.CompleteCitizenLabelsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	jobs, err := s.ps.AttrJobList(
		userInfo,
		permscompletor.CompletorServicePerm,
		permscompletor.CompletorServiceCompleteCitizenLabelsPerm,
		permscompletor.CompletorServiceCompleteCitizenLabelsJobsPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscompletor.ErrFailedSearch)
	}
	if jobs.Len() == 0 {
		jobs.Strings = append(jobs.Strings, userInfo.GetJob())
	}

	jobsExp := make([]mysql.Expression, jobs.Len())
	for i := range jobs.GetStrings() {
		jobsExp[i] = mysql.String(jobs.GetStrings()[i])
	}

	condition := mysql.AND(
		tCitizensLabelsJob.Job.IN(jobsExp...),
		tCitizensLabelsJob.DeletedAt.IS_NULL(),
	)

	if search := dbutils.PrepareForLikeSearch(req.GetSearch()); search != "" {
		condition = condition.AND(tCitizensLabelsJob.Name.LIKE(mysql.String(search)))
	}

	stmt := tCitizensLabelsJob.
		SELECT(
			tCitizensLabelsJob.ID,
			tCitizensLabelsJob.CreatedAt,
			tCitizensLabelsJob.Name,
			tCitizensLabelsJob.Color,
			tCitizensLabelsJob.Icon,
		).
		FROM(tCitizensLabelsJob).
		WHERE(condition).
		ORDER_BY(
			tCitizensLabelsJob.SortKey.ASC(),
		).
		LIMIT(15)

	resp := &pbcompletor.CompleteCitizenLabelsResponse{
		Labels: []*citizenslabels.Label{},
	}
	if err := stmt.QueryContext(ctx, s.db, &resp.Labels); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscompletor.ErrFailedSearch)
		}
	}

	return resp, nil
}
