package completor

import (
	context "context"
	"errors"
	"slices"
	"strings"

	users "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	pbcompletor "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/completor"
	permscompletor "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/completor/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorscompletor "github.com/fivenet-app/fivenet/v2025/services/completor/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tDCategory         = table.FivenetDocumentsCategories.AS("category")
	tCitizensLabelsJob = table.FivenetUserLabelsJob.AS("label")
)

func (s *Server) CompleteCitizens(ctx context.Context, req *pbcompletor.CompleteCitizensRequest) (*pbcompletor.CompleteCitizensResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tUsers := tables.User().AS("user_short")

	condition := s.customDB.Conditions.User.GetFilter(tUsers.Alias())

	currentJob := false
	if req.CurrentJob != nil && *req.CurrentJob {
		currentJob = true
	}

	orderBys := []jet.OrderByClause{}
	if len(req.UserIds) == 0 {
		if currentJob {
			condition = condition.AND(
				tUsers.Job.EQ(jet.String(userInfo.Job)),
			)
		}

		req.Search = strings.TrimSpace(req.Search)
		req.Search = strings.ReplaceAll(req.Search, "%", "")
		req.Search = strings.ReplaceAll(req.Search, " ", "%")

		if req.Search != "" {
			req.Search = "%" + req.Search + "%"
			condition = jet.CONCAT(tUsers.Firstname, jet.String(" "), tUsers.Lastname).
				LIKE(jet.String(req.Search))
		}
	} else {
		userIds := []jet.Expression{}
		for _, v := range req.UserIds {
			userIds = append(userIds, jet.Int32(v))
		}

		if req.UserIdsOnly != nil && *req.UserIdsOnly {
			condition = condition.AND(tUsers.ID.IN(userIds...))
		}

		// Make sure to sort by the user IDs if provided
		orderBys = append(orderBys, tUsers.ID.IN(userIds...).DESC())
	}

	orderBys = append(orderBys, tUsers.Lastname.ASC())

	columns := jet.ProjectionList{
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
		OPTIMIZER_HINTS(jet.OptimizerHint("idx_users_firstname_lastname_fulltext")).
		FROM(tUsers).
		WHERE(condition).
		ORDER_BY(orderBys...).
		LIMIT(15)

	var dest []*users.UserShort
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscompletor.ErrFailedSearch)
		}
	}

	if req.OnDuty != nil && *req.OnDuty {
		dest = slices.DeleteFunc(dest, func(us *users.UserShort) bool {
			return !s.tracker.IsUserOnDuty(us.UserId)
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

func (s *Server) CompleteJobs(ctx context.Context, req *pbcompletor.CompleteJobsRequest) (*pbcompletor.CompleteJobsResponse, error) {
	var search string
	if req.Search != nil && *req.Search != "" {
		search = *req.Search
	}
	if req.CurrentJob != nil && *req.CurrentJob {
		userInfo := auth.MustGetUserInfoFromContext(ctx)
		search = userInfo.Job
	}
	exactMatch := false
	if req.ExactMatch != nil {
		exactMatch = *req.ExactMatch
	}

	resp := &pbcompletor.CompleteJobsResponse{}
	if search != "" {
		var err error
		resp.Jobs, err = s.jobsS.Search(ctx, search, exactMatch)
		if err != nil {
			return nil, errswrap.NewError(err, errorscompletor.ErrFailedSearch)
		}
	} else {
		resp.Jobs = s.jobsS.List()
	}

	return resp, nil
}

func (s *Server) CompleteDocumentCategories(ctx context.Context, req *pbcompletor.CompleteDocumentCategoriesRequest) (*pbcompletor.CompleteDocumentCategoriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	jobs, err := s.p.AttrJobList(userInfo, permscompletor.CompletorServicePerm, permscompletor.CompletorServiceCompleteDocumentCategoriesPerm, permscompletor.CompletorServiceCompleteDocumentCategoriesJobsPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorscompletor.ErrFailedSearch)
	}
	if jobs.Len() == 0 {
		jobs.Strings = append(jobs.Strings, userInfo.Job)
	}

	req.Search = strings.TrimSpace(req.Search)
	req.Search = strings.ReplaceAll(req.Search, "%", "")
	req.Search = strings.ReplaceAll(req.Search, " ", "%")

	jobsExp := make([]jet.Expression, jobs.Len())
	for i := range jobs.Strings {
		jobsExp[i] = jet.String(jobs.Strings[i])
	}

	condition := tDCategory.Job.IN(jobsExp...)
	if req.Search != "" {
		req.Search = "%" + req.Search + "%"
		condition = condition.AND(
			tDCategory.Name.LIKE(jet.String(req.Search)),
		)
	}

	stmt := tDCategory.
		SELECT(
			tDCategory.ID,
			tDCategory.Name,
			tDCategory.Description,
			tDCategory.Job,
			tDCategory.Color,
			tDCategory.Icon,
		).
		FROM(tDCategory).
		WHERE(condition).
		ORDER_BY(
			tDCategory.SortKey.ASC(),
		).
		LIMIT(15)

	resp := &pbcompletor.CompleteDocumentCategoriesResponse{}
	if err := stmt.QueryContext(ctx, s.db, &resp.Categories); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscompletor.ErrFailedSearch)
		}
	}

	return resp, nil
}

func (s *Server) ListLawBooks(ctx context.Context, req *pbcompletor.ListLawBooksRequest) (*pbcompletor.ListLawBooksResponse, error) {
	return &pbcompletor.ListLawBooksResponse{
		Books: s.laws.GetLawBooks(),
	}, nil
}

func (s *Server) CompleteCitizenLabels(ctx context.Context, req *pbcompletor.CompleteCitizenLabelsRequest) (*pbcompletor.CompleteCitizenLabelsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	jobs, err := s.p.AttrJobList(userInfo, permscompletor.CompletorServicePerm, permscompletor.CompletorServiceCompleteCitizenLabelsPerm, permscompletor.CompletorServiceCompleteCitizenLabelsJobsPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorscompletor.ErrFailedSearch)
	}
	if jobs.Len() == 0 {
		jobs.Strings = append(jobs.Strings, userInfo.Job)
	}

	req.Search = strings.TrimSpace(req.Search)
	req.Search = strings.ReplaceAll(req.Search, "%", "")
	req.Search = strings.ReplaceAll(req.Search, " ", "%")

	jobsExp := make([]jet.Expression, jobs.Len())
	for i := range jobs.Strings {
		jobsExp[i] = jet.String(jobs.Strings[i])
	}

	condition := tCitizensLabelsJob.Job.IN(jobsExp...)

	if req.Search != "" {
		req.Search = "%" + req.Search + "%"
		condition = condition.AND(tCitizensLabelsJob.Name.LIKE(jet.String(req.Search)))
	}

	stmt := tCitizensLabelsJob.
		SELECT(
			tCitizensLabelsJob.ID,
			tCitizensLabelsJob.Name,
			tCitizensLabelsJob.Color,
		).
		FROM(tCitizensLabelsJob).
		WHERE(condition).
		ORDER_BY(
			tCitizensLabelsJob.SortKey.ASC(),
		).
		LIMIT(10)

	resp := &pbcompletor.CompleteCitizenLabelsResponse{
		Labels: []*users.Label{},
	}
	if err := stmt.QueryContext(ctx, s.db, &resp.Labels); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscompletor.ErrFailedSearch)
		}
	}

	return resp, nil
}
