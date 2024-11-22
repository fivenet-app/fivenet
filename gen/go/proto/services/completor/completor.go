package completor

import (
	context "context"
	"database/sql"
	"errors"
	"slices"
	"strings"

	users "github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	errorscompletor "github.com/fivenet-app/fivenet/gen/go/proto/services/completor/errors"
	permscompletor "github.com/fivenet-app/fivenet/gen/go/proto/services/completor/perms"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/tracker"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
	grpc "google.golang.org/grpc"
)

var (
	tUsers                = table.Users.AS("usershort")
	tDCategory            = table.FivenetDocumentsCategories.AS("category")
	tJobCitizenAttributes = table.FivenetJobCitizenAttributes.AS("citizen_attribute")
)

type Server struct {
	CompletorServiceServer

	db       *sql.DB
	p        perms.Permissions
	data     *mstlystcdata.Cache
	tracker  tracker.ITracker
	enricher *mstlystcdata.UserAwareEnricher

	customDB config.CustomDB
}

type Params struct {
	fx.In

	DB       *sql.DB
	Perms    perms.Permissions
	Data     *mstlystcdata.Cache
	Tracker  tracker.ITracker
	Enricher *mstlystcdata.UserAwareEnricher
	Config   *config.Config
}

func NewServer(p Params) *Server {
	s := &Server{
		db:       p.DB,
		p:        p.Perms,
		data:     p.Data,
		tracker:  p.Tracker,
		enricher: p.Enricher,

		customDB: p.Config.Database.Custom,
	}

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterCompletorServiceServer(srv, s)
}

func (s *Server) CompleteCitizens(ctx context.Context, req *CompleteCitizensRequest) (*CompleteCitizensRespoonse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := s.customDB.Conditions.User.GetFilter(tUsers.Alias())

	currentJob := false
	if req.CurrentJob != nil && *req.CurrentJob {
		currentJob = true
	}

	if req.UserId == nil {
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
		condition = condition.AND(tUsers.ID.EQ(jet.Int32(*req.UserId)))
	}

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
		ORDER_BY(
			tUsers.Lastname.ASC(),
		).
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
		for i := 0; i < len(dest); i++ {
			jobInfoFn(dest[i])
		}
	}

	return &CompleteCitizensRespoonse{
		Users: dest,
	}, nil
}

func (s *Server) CompleteJobs(ctx context.Context, req *CompleteJobsRequest) (*CompleteJobsResponse, error) {
	var search string
	if req.Search != nil {
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

	resp := &CompleteJobsResponse{}
	var err error
	resp.Jobs, err = s.data.GetSearcher().SearchJobs(ctx, search, exactMatch)
	if err != nil {
		return nil, errswrap.NewError(err, errorscompletor.ErrFailedSearch)
	}

	return resp, nil
}

func (s *Server) CompleteDocumentCategories(ctx context.Context, req *CompleteDocumentCategoriesRequest) (*CompleteDocumentCategoriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	jobsAttr, err := s.p.Attr(userInfo, permscompletor.CompletorServicePerm, permscompletor.CompletorServiceCompleteDocumentCategoriesPerm, permscompletor.CompletorServiceCompleteDocumentCategoriesJobsPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorscompletor.ErrFailedSearch)
	}
	var jobs perms.StringList
	if jobsAttr != nil {
		jobs = jobsAttr.([]string)
	}

	if len(jobs) == 0 {
		jobs = append(jobs, userInfo.Job)
	}

	req.Search = strings.TrimSpace(req.Search)
	req.Search = strings.ReplaceAll(req.Search, "%", "")
	req.Search = strings.ReplaceAll(req.Search, " ", "%")

	jobsExp := make([]jet.Expression, len(jobs))
	for i := 0; i < len(jobs); i++ {
		jobsExp[i] = jet.String(jobs[i])
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
			tDCategory.Name.ASC(),
		).
		LIMIT(15)

	resp := &CompleteDocumentCategoriesResponse{}
	if err := stmt.QueryContext(ctx, s.db, &resp.Categories); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscompletor.ErrFailedSearch)
		}
	}

	return resp, nil
}

func (s *Server) ListLawBooks(ctx context.Context, req *ListLawBooksRequest) (*ListLawBooksResponse, error) {
	return &ListLawBooksResponse{
		Books: s.data.GetLawBooks(),
	}, nil
}

func (s *Server) CompleteCitizenAttributes(ctx context.Context, req *CompleteCitizenAttributesRequest) (*CompleteCitizenAttributesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	jobsAttr, err := s.p.Attr(userInfo, permscompletor.CompletorServicePerm, permscompletor.CompletorServiceCompleteCitizenAttributesPerm, permscompletor.CompletorServiceCompleteCitizenAttributesJobsPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorscompletor.ErrFailedSearch)
	}
	var jobs perms.StringList
	if jobsAttr != nil {
		jobs = jobsAttr.([]string)
	}

	if len(jobs) == 0 {
		jobs = append(jobs, userInfo.Job)
	}

	req.Search = strings.TrimSpace(req.Search)
	req.Search = strings.ReplaceAll(req.Search, "%", "")
	req.Search = strings.ReplaceAll(req.Search, " ", "%")

	jobsExp := make([]jet.Expression, len(jobs))
	for i := 0; i < len(jobs); i++ {
		jobsExp[i] = jet.String(jobs[i])
	}

	condition := tJobCitizenAttributes.Job.IN(jobsExp...)

	if req.Search != "" {
		req.Search = "%" + req.Search + "%"
		condition = condition.AND(tJobCitizenAttributes.Name.LIKE(jet.String(req.Search)))
	}

	stmt := tJobCitizenAttributes.
		SELECT(
			tJobCitizenAttributes.ID,
			tJobCitizenAttributes.Name,
			tJobCitizenAttributes.Color,
		).
		FROM(tJobCitizenAttributes).
		WHERE(condition).
		ORDER_BY(
			tJobCitizenAttributes.Name.ASC(),
		).
		LIMIT(10)

	resp := &CompleteCitizenAttributesResponse{
		Attributes: []*users.CitizenAttribute{},
	}
	if err := stmt.QueryContext(ctx, s.db, &resp.Attributes); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscompletor.ErrFailedSearch)
		}
	}

	return resp, nil
}
