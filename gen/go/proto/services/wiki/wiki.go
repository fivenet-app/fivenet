package wiki

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"strings"

	database "github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/wiki"
	errorswiki "github.com/fivenet-app/fivenet/gen/go/proto/services/wiki/errors"
	permswiki "github.com/fivenet-app/fivenet/gen/go/proto/services/wiki/perms"
	"github.com/fivenet-app/fivenet/pkg/access"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/html/htmldiffer"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/server/audit"
	"github.com/fivenet-app/fivenet/pkg/utils"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/gosimple/slug"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	tPage        = table.FivenetWikiPages.AS("page")
	tPageShort   = table.FivenetWikiPages.AS("pageshort")
	tPJobAccess  = table.FivenetWikiPageJobAccess.AS("job_access")
	tPUserAccess = table.FivenetWikiPageUserAccess.AS("user_access")
	tUsers       = table.Users
	tCreator     = tUsers.AS("creator")
	tJobProps    = table.FivenetJobProps
)

type Server struct {
	WikiServiceServer

	logger   *zap.Logger
	db       *sql.DB
	aud      audit.IAuditer
	perms    perms.Permissions
	enricher *mstlystcdata.UserAwareEnricher
	htmlDiff *htmldiffer.Differ

	access *access.Grouped[wiki.PageJobAccess, *wiki.PageJobAccess, wiki.PageUserAccess, *wiki.PageUserAccess, wiki.AccessLevel]
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger     *zap.Logger
	DB         *sql.DB
	Audit      audit.IAuditer
	Perms      perms.Permissions
	Enricher   *mstlystcdata.UserAwareEnricher
	HTMLDiffer *htmldiffer.Differ
}

func NewServer(p Params) *Server {
	s := &Server{
		logger:   p.Logger.Named("wiki"),
		db:       p.DB,
		aud:      p.Audit,
		perms:    p.Perms,
		enricher: p.Enricher,
		htmlDiff: p.HTMLDiffer,

		access: access.NewGrouped(
			p.DB,
			table.FivenetWikiPages,
			&access.TargetTableColumns{
				ID:         table.FivenetWikiPages.ID,
				DeletedAt:  table.FivenetWikiPages.DeletedAt,
				CreatorID:  table.FivenetWikiPages.CreatorID,
				CreatorJob: table.FivenetWikiPages.Job,
			},
			access.NewJobs[wiki.PageJobAccess, *wiki.PageJobAccess, wiki.AccessLevel](
				table.FivenetWikiPageJobAccess,
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:        table.FivenetWikiPageJobAccess.ID,
						CreatedAt: table.FivenetWikiPageJobAccess.CreatedAt,
						TargetID:  table.FivenetWikiPageJobAccess.PageID,
						Access:    table.FivenetWikiPageJobAccess.Access,
					},
					Job:          table.FivenetWikiPageJobAccess.Job,
					MinimumGrade: table.FivenetWikiPageJobAccess.MinimumGrade,
				},
				table.FivenetWikiPageJobAccess.AS("page_job_access"),
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:        table.FivenetWikiPageJobAccess.AS("page_job_access").ID,
						CreatedAt: table.FivenetWikiPageJobAccess.AS("page_job_access").CreatedAt,
						TargetID:  table.FivenetWikiPageJobAccess.AS("page_job_access").PageID,
						Access:    table.FivenetWikiPageJobAccess.AS("page_job_access").Access,
					},
					Job:          table.FivenetWikiPageJobAccess.AS("page_job_access").Job,
					MinimumGrade: table.FivenetWikiPageJobAccess.AS("page_job_access").MinimumGrade,
				},
			),
			access.NewUsers[wiki.PageUserAccess, *wiki.PageUserAccess, wiki.AccessLevel](
				table.FivenetWikiPageUserAccess,
				&access.UserAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:        table.FivenetWikiPageUserAccess.ID,
						CreatedAt: table.FivenetWikiPageUserAccess.CreatedAt,
						TargetID:  table.FivenetWikiPageUserAccess.PageID,
						Access:    table.FivenetWikiPageUserAccess.Access,
					},
					UserId: table.FivenetWikiPageUserAccess.UserID,
				},
				table.FivenetWikiPageUserAccess.AS("page_user_access"),
				&access.UserAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:        table.FivenetWikiPageUserAccess.AS("page_user_access").ID,
						CreatedAt: table.FivenetWikiPageUserAccess.AS("page_user_access").CreatedAt,
						TargetID:  table.FivenetWikiPageUserAccess.AS("page_user_access").PageID,
						Access:    table.FivenetWikiPageUserAccess.AS("page_user_access").Access,
					},
					UserId: table.FivenetWikiPageUserAccess.AS("page_user_access").UserID,
				},
			),
		),
	}

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterWikiServiceServer(srv, s)
}

func (s *Server) ListPages(ctx context.Context, req *ListPagesRequest) (*ListPagesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: WikiService_ServiceDesc.ServiceName,
		Method:  "ListPages",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	groupBys := []jet.GroupByClause{tPageShort.ID}
	if req.RootOnly != nil && *req.RootOnly {
		groupBys = []jet.GroupByClause{tPageShort.Job}
	}

	condition := jet.Bool(true)
	if req.Search != nil && *req.Search != "" {
		*req.Search = strings.TrimRight(*req.Search, "*") + "*"

		condition = jet.OR(
			jet.BoolExp(
				jet.Raw("MATCH(`title`) AGAINST ($search IN BOOLEAN MODE)",
					jet.RawArgs{"$search": *req.Search}),
			),
			jet.BoolExp(
				jet.Raw("MATCH(`content`) AGAINST ($search IN BOOLEAN MODE)",
					jet.RawArgs{"$search": *req.Search}),
			),
		)
	}

	if !userInfo.SuperUser {
		condition = condition.AND(jet.AND(
			tPageShort.DeletedAt.IS_NULL(),
			jet.OR(
				tPageShort.Public.IS_TRUE(),
				tPageShort.CreatorID.EQ(jet.Int32(userInfo.UserId)),
				jet.OR(
					jet.AND(
						tPUserAccess.Access.IS_NOT_NULL(),
						tPUserAccess.Access.GT(jet.Int32(int32(wiki.AccessLevel_ACCESS_LEVEL_BLOCKED))),
					),
					jet.AND(
						tPUserAccess.Access.IS_NULL(),
						tPJobAccess.Access.IS_NOT_NULL(),
						tPJobAccess.Access.GT(jet.Int32(int32(wiki.AccessLevel_ACCESS_LEVEL_BLOCKED))),
					),
				),
			),
		))
	}

	if req.Job != nil {
		if *req.Job == "" {
			*req.Job = userInfo.Job
		}
		condition = condition.AND(tPageShort.Job.EQ(jet.String(*req.Job)))
	}

	countStmt := tPageShort.
		SELECT(
			jet.COUNT(jet.DISTINCT(tPageShort.ID)).AS("datacount.totalcount"),
		).
		FROM(
			tPageShort.
				LEFT_JOIN(tPUserAccess,
					tPUserAccess.PageID.EQ(tPageShort.ID).
						AND(tPUserAccess.UserID.EQ(jet.Int32(userInfo.UserId))),
				).
				LEFT_JOIN(tPJobAccess,
					tPJobAccess.PageID.EQ(tPageShort.ID).
						AND(tPJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tPJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, 40)
	resp := &ListPagesResponse{
		Pagination: pag,
		Pages:      []*wiki.PageShort{},
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	columns := []jet.Projection{
		tPageShort.Job,
		tPageShort.ParentID,
		tPageShort.Slug,
		tPageShort.Title,
		tPageShort.Description,
	}
	if req.RootOnly != nil && *req.RootOnly {
		columns = append(columns,
			tJobProps.LogoURL.AS("page_root_info.logo"),
		)
	}

	stmt := tPageShort.
		SELECT(
			tPageShort.ID,
			columns...,
		).
		FROM(
			tPageShort.
				LEFT_JOIN(tPUserAccess,
					tPUserAccess.PageID.EQ(tPageShort.ID).
						AND(tPUserAccess.UserID.EQ(jet.Int32(userInfo.UserId))),
				).
				LEFT_JOIN(tPJobAccess,
					tPJobAccess.PageID.EQ(tPageShort.ID).
						AND(tPJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tPJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				).
				LEFT_JOIN(tJobProps,
					tJobProps.Job.EQ(tPageShort.Job),
				),
		).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		ORDER_BY(tPageShort.ParentID.ASC().NULLS_FIRST(), tPageShort.Title.ASC()).
		GROUP_BY(groupBys...).
		LIMIT(limit)

	pages := []*wiki.PageShort{}
	if err := stmt.QueryContext(ctx, s.db, &pages); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errorswiki.ErrFailedQuery
		}
	}

	for i := range pages {
		s.enricher.EnrichJobName(pages[i])
	}

	if req.Search == nil && (req.RootOnly == nil || !*req.RootOnly) {
		for _, page := range mapPagesToNavItems(pages) {
			resp.Pages = append(resp.Pages, page)
		}
	} else {
		resp.Pages = pages
	}

	return resp, nil
}

func (s *Server) GetPage(ctx context.Context, req *GetPageRequest) (*GetPageResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: WikiService_ServiceDesc.ServiceName,
		Method:  "GetPage",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	condition := jet.AND(
		tPage.ID.EQ(jet.Uint64(req.Id)),
		tPage.DeletedAt.IS_NULL(),
		jet.OR(
			tPage.Public.IS_TRUE(),
			tPage.CreatorID.EQ(jet.Int32(userInfo.UserId)),
			jet.AND(
				tPJobAccess.Access.IS_NOT_NULL(),
				tPJobAccess.Access.GT(jet.Int32(int32(wiki.AccessLevel_ACCESS_LEVEL_BLOCKED))),
			),
		),
	)

	stmt := tPage.
		SELECT(
			tPage.ID,
			tPage.Job,
			tPage.ParentID,
			tPage.CreatedAt.AS("page_meta.created_at"),
			tPage.UpdatedAt.AS("page_meta.updated_at"),
			tPage.DeletedAt.AS("page_meta.deleted_at"),
			tPage.Slug.AS("page_meta.slug"),
			tPage.Title.AS("page_meta.title"),
			tPage.Description.AS("page_meta.description"),
			tPage.CreatorID.AS("page_meta.creator_id"),
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tPage.ContentType.AS("page_meta.content_Type"),
			tPage.Content,
			tPage.Data,
		).
		FROM(
			tPage.
				LEFT_JOIN(tPJobAccess,
					tPJobAccess.PageID.EQ(tPage.ID).
						AND(tPJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tPJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				).
				LEFT_JOIN(tCreator,
					tPage.CreatorID.EQ(tCreator.ID),
				),
		).
		WHERE(condition).
		LIMIT(1)

	resp := &GetPageResponse{
		Page: &wiki.Page{},
	}

	if err := stmt.QueryContext(ctx, s.db, resp.Page); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
	}

	if resp.Page.Id <= 0 {
		resp.Page = nil
	} else {
		s.enricher.EnrichJobName(resp.Page)

		access, err := s.getPageAccess(ctx, userInfo, req.Id)
		if err != nil {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
		resp.Page.Access = access

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_VIEWED)
	}

	return resp, nil
}

func (s *Server) getPageAccess(ctx context.Context, userInfo *userinfo.UserInfo, pageId uint64) (*wiki.PageAccess, error) {
	jobsAccess, err := s.access.Jobs.List(ctx, s.db, pageId)
	if err != nil {
		return nil, errorswiki.ErrFailedQuery
	}

	usersAccess, err := s.access.Users.List(ctx, s.db, pageId)
	if err != nil {
		return nil, errorswiki.ErrFailedQuery
	}

	for i := 0; i < len(jobsAccess); i++ {
		s.enricher.EnrichJobInfo(jobsAccess[i])
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := 0; i < len(usersAccess); i++ {
		if usersAccess[i].User != nil {
			jobInfoFn(usersAccess[i].User)
		}
	}

	return &wiki.PageAccess{
		Jobs:  jobsAccess,
		Users: usersAccess,
	}, nil
}

func (s *Server) getPage(ctx context.Context, pageId uint64, withContent bool, withAccess bool, userInfo *userinfo.UserInfo) (*wiki.Page, error) {
	columns := []jet.Projection{
		tPage.ID,
		tPage.Job,
		tPage.ParentID,
		tPage.CreatedAt.AS("page_meta.created_at"),
		tPage.UpdatedAt.AS("page_meta.updated_at"),
		tPage.DeletedAt.AS("page_meta.deleted_at"),
		tPage.Slug.AS("page_meta.slug"),
		tPage.Title.AS("page_meta.title"),
		tPage.Description.AS("page_meta.description"),
		tPage.CreatorID.AS("page_meta.creator_id"),
		tCreator.ID,
		tCreator.Job,
		tCreator.JobGrade,
		tCreator.Firstname,
		tCreator.Lastname,
		tCreator.Dateofbirth,
		tPage.ContentType.AS("page_meta.content_Type"),
		tPage.Toc.AS("page_meta.toc"),
		tPage.Public.AS("page_meta.public"),
	}
	if withContent {
		columns = append(columns,
			tPage.Content,
			tPage.Data,
		)
	}

	stmt := tPage.
		SELECT(
			tPage.ID,

			columns...,
		).
		FROM(
			tPage.
				LEFT_JOIN(tCreator,
					tPage.CreatorID.EQ(tCreator.ID),
				),
		).
		WHERE(jet.AND(
			tPage.ID.EQ(jet.Uint64(pageId)),
		)).
		LIMIT(1)

	var dest wiki.Page
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.Id == 0 {
		return nil, nil
	}

	s.enricher.EnrichJobName(&dest)

	if withAccess {
		access, err := s.getPageAccess(ctx, userInfo, pageId)
		if err != nil {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
		dest.Access = access
	}

	return &dest, nil
}

func (s *Server) CreatePage(ctx context.Context, req *CreatePageRequest) (*CreatePageResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: WikiService_ServiceDesc.ServiceName,
		Method:  "CreatePage",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	if req.Page.ParentId == nil || *req.Page.ParentId <= 0 {
		countStmt := tPage.
			SELECT(
				jet.COUNT(tPage.ID).AS("datacount.totalcount"),
			).
			FROM(tPage).
			WHERE(jet.AND(
				tPage.Job.EQ(jet.String(userInfo.Job)),
				tPage.DeletedAt.IS_NULL(),
			))

		var count database.DataCount
		if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
			}
		}

		if count.TotalCount > 0 {
			return nil, errorswiki.ErrPageDenied
		}
	} else {
		p, err := s.getPage(ctx, *req.Page.ParentId, false, false, nil)
		if err != nil {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}

		if p.Job != userInfo.Job {
			return nil, errorswiki.ErrPageDenied
		}

		parentCheck, err := s.access.CanUserAccessTarget(ctx, *req.Page.ParentId, userInfo, wiki.AccessLevel_ACCESS_LEVEL_VIEW)
		if err != nil {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
		if !parentCheck {
			return nil, errorswiki.ErrPageDenied
		}
	}

	// Field Permission Check
	fieldsAttr, err := s.perms.Attr(userInfo, permswiki.WikiServicePerm, permswiki.WikiServiceCreatePagePerm, permswiki.WikiServiceCreatePageFieldsPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}
	if !slices.Contains(fields, "Public") {
		req.Page.Meta.Public = false
	}

	tPage := table.FivenetWikiPages
	stmt := tPage.
		INSERT(
			tPage.Job,
			tPage.ParentID,
			tPage.ContentType,
			tPage.Toc,
			tPage.Public,
			tPage.Slug,
			tPage.Title,
			tPage.Description,
			tPage.Content,
			tPage.Data,
			tPage.CreatorID,
		).
		VALUES(
			userInfo.Job,
			req.Page.ParentId,
			req.Page.Meta.ContentType,
			req.Page.Meta.Toc,
			req.Page.Meta.Public,
			slug.Make(utils.StringFirstN(req.Page.Meta.Title, 100)),
			req.Page.Meta.Title,
			req.Page.Meta.Description,
			req.Page.Content,
			nil,
			userInfo.UserId,
		)

	result, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	req.Page.Id = uint64(lastId)

	if _, err := s.addPageActivity(ctx, tx, &wiki.PageActivity{
		PageId:       uint64(lastId),
		ActivityType: wiki.PageActivityType_PAGE_ACTIVITY_TYPE_CREATED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.Job,
	}); err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	if err := s.handlePageAccessChange(ctx, tx, req.Page.Id, userInfo, req.Page.Access, false); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	page, err := s.getPage(ctx, req.Page.Id, true, true, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	return &CreatePageResponse{
		Page: page,
	}, nil
}

func (s *Server) UpdatePage(ctx context.Context, req *UpdatePageRequest) (*UpdatePageResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: WikiService_ServiceDesc.ServiceName,
		Method:  "UpdatePage",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	if req.Page.ParentId == nil || *req.Page.ParentId <= 0 {
		stmt := tPage.
			SELECT(
				tPage.ID.AS("id"),
			).
			FROM(tPage).
			WHERE(jet.AND(
				tPage.Job.EQ(jet.String(userInfo.Job)),
				tPage.DeletedAt.IS_NULL(),
			))

		var ids struct {
			ID uint64 `alias:"id"`
		}
		if err := stmt.QueryContext(ctx, s.db, &ids); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
			}
		}

		if ids.ID != req.Page.Id {
			return nil, errorswiki.ErrPageDenied
		}
	} else {
		p, err := s.getPage(ctx, *req.Page.ParentId, false, false, nil)
		if err != nil {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}

		if p.Job != userInfo.Job {
			return nil, errorswiki.ErrPageDenied
		}

		parentCheck, err := s.access.CanUserAccessTarget(ctx, *req.Page.ParentId, userInfo, wiki.AccessLevel_ACCESS_LEVEL_VIEW)
		if err != nil {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
		// Reset parent id to current one
		if !parentCheck {
			*req.Page.ParentId = *p.ParentId
		}
	}

	tPage := table.FivenetWikiPages

	check, err := s.access.CanUserAccessTarget(ctx, req.Page.Id, userInfo, wiki.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	if !check {
		return nil, errorswiki.ErrPageDenied
	}

	page, err := s.getPage(ctx, req.Page.Id, true, true, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	stmt := tPage.
		UPDATE(
			tPage.ParentID,
			tPage.ContentType,
			tPage.Toc,
			tPage.Public,
			tPage.Slug,
			tPage.Title,
			tPage.Description,
			tPage.Content,
			tPage.Data,
		).
		SET(
			req.Page.ParentId,
			req.Page.Meta.ContentType,
			req.Page.Meta.Toc,
			req.Page.Meta.Public,
			slug.Make(utils.StringFirstN(req.Page.Meta.Title, 100)),
			req.Page.Meta.Title,
			req.Page.Meta.Description,
			req.Page.Content,
			nil,
		).
		WHERE(jet.AND(
			tPage.ID.EQ(jet.Uint64(req.Page.Id)),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	diff, err := s.generatePageDiff(page, &wiki.Page{
		Meta: &wiki.PageMeta{
			Title:       req.Page.Meta.Title,
			Description: req.Page.Meta.Description,
		},
		Content: req.Page.Content,
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	if _, err := s.addPageActivity(ctx, tx, &wiki.PageActivity{
		PageId:       req.Page.Id,
		ActivityType: wiki.PageActivityType_PAGE_ACTIVITY_TYPE_UPDATED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.Job,
		Data: &wiki.PageActivityData{
			Data: &wiki.PageActivityData_Updated{
				Updated: diff,
			},
		},
	}); err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	if err := s.handlePageAccessChange(ctx, tx, req.Page.Id, userInfo, req.Page.Access, true); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	page, err = s.getPage(ctx, req.Page.Id, true, true, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	return &UpdatePageResponse{
		Page: page,
	}, nil
}

func (s *Server) handlePageAccessChange(ctx context.Context, tx qrm.DB, pageId uint64, userInfo *userinfo.UserInfo, access *wiki.PageAccess, addActivity bool) error {
	changes, err := s.access.HandleAccessChanges(ctx, tx, pageId, access.Jobs, access.Users)
	if err != nil {
		if dbutils.IsDuplicateError(err) {
			return errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
		return errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	if addActivity && !changes.IsEmpty() {
		if _, err := s.addPageActivity(ctx, tx, &wiki.PageActivity{
			PageId:       pageId,
			ActivityType: wiki.PageActivityType_PAGE_ACTIVITY_TYPE_ACCESS_UPDATED,
			CreatorId:    &userInfo.UserId,
			CreatorJob:   userInfo.Job,
			Data: &wiki.PageActivityData{
				Data: &wiki.PageActivityData_AccessUpdated{
					AccessUpdated: &wiki.PageAccessUpdated{
						Jobs: &wiki.PageAccessJobsDiff{
							ToCreate: changes.Jobs.ToCreate,
							ToUpdate: changes.Jobs.ToUpdate,
							ToDelete: changes.Jobs.ToDelete,
						},
						Users: &wiki.PageAccessUsersDiff{
							ToCreate: changes.Users.ToCreate,
							ToUpdate: changes.Users.ToUpdate,
							ToDelete: changes.Users.ToDelete,
						},
					},
				},
			},
		}); err != nil {
			return errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
	}

	return nil
}

func (s *Server) DeletePage(ctx context.Context, req *DeletePageRequest) (*DeletePageResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: WikiService_ServiceDesc.ServiceName,
		Method:  "DeletePage",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.Id, userInfo, wiki.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	if !check {
		return nil, errorswiki.ErrPageDenied
	}

	stmt := tPage.
		UPDATE(
			tPage.DeletedAt,
		).
		SET(
			tPage.DeletedAt.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(jet.AND(
			tPage.ID.EQ(jet.Uint64(req.Id)),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeletePageResponse{}, nil
}
