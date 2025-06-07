package wiki

import (
	"context"
	"errors"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/wiki"
	pbwiki "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/wiki"
	permswiki "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/wiki/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorswiki "github.com/fivenet-app/fivenet/v2025/services/wiki/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/gosimple/slug"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *Server) ListPages(ctx context.Context, req *pbwiki.ListPagesRequest) (*pbwiki.ListPagesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbwiki.WikiService_ServiceDesc.ServiceName,
		Method:  "ListPages",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

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

	groupBys := []jet.GroupByClause{tPageShort.ID}
	if req.RootOnly != nil && *req.RootOnly {
		groupBys = []jet.GroupByClause{tPageShort.Job}
	}

	if !userInfo.Superuser {
		condition = condition.AND(jet.AND(
			tPageShort.DeletedAt.IS_NULL(),
			jet.OR(
				tPageShort.Public.IS_TRUE(),
				tPageShort.CreatorID.EQ(jet.Int32(userInfo.UserId)),
				jet.OR(
					tPAccess.UserID.EQ(jet.Int32(userInfo.UserId)),
					jet.AND(
						tPAccess.Job.EQ(jet.String(userInfo.Job)),
						tPAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade)),
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
			jet.COUNT(jet.DISTINCT(tPageShort.ID)).AS("data_count.total"),
		).
		FROM(
			tPageShort.
				LEFT_JOIN(tPAccess,
					tPAccess.TargetID.EQ(tPageShort.ID).
						AND(tPAccess.Access.GT_EQ(jet.Int32(int32(wiki.AccessLevel_ACCESS_LEVEL_VIEW)))),
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.Total, defaultWikiUpperLimit)
	resp := &pbwiki.ListPagesResponse{
		Pagination: pag,
		Pages:      []*wiki.PageShort{},
	}
	if count.Total <= 0 {
		return resp, nil
	}

	tFiles := table.FivenetFiles.AS("logo")

	columns := []jet.Projection{
		tPageShort.Job,
		tPageShort.ParentID,
		tPageShort.Slug,
		tPageShort.Title,
		tPageShort.Description,
		tPageShort.Draft,
	}
	if req.RootOnly != nil && *req.RootOnly {
		columns = append(columns,
			tJobProps.LogoFileID.AS("page_root_info.logo"),
			tFiles.ID,
			tFiles.FilePath,
		)
	}
	if userInfo.Superuser {
		columns = append(columns, tPageShort.DeletedAt)
	}

	stmt := tPageShort.
		SELECT(
			tPageShort.ID,
			columns...,
		).
		FROM(
			tPageShort.
				LEFT_JOIN(tPAccess,
					tPAccess.TargetID.EQ(tPageShort.ID).
						AND(tPAccess.Access.GT_EQ(jet.Int32(int32(wiki.AccessLevel_ACCESS_LEVEL_VIEW)))),
				).
				LEFT_JOIN(tJobProps,
					tJobProps.Job.EQ(tPageShort.Job),
				).
				LEFT_JOIN(tFiles,
					tFiles.ID.EQ(tJobProps.LogoFileID),
				),
		).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		// .NULLS_FIRST()
		ORDER_BY(tPageShort.ParentID.ASC(), tPageShort.Draft.ASC(), tPageShort.SortKey.ASC()).
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

func (s *Server) GetPage(ctx context.Context, req *pbwiki.GetPageRequest) (*pbwiki.GetPageResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.wiki.page_id", int64(req.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbwiki.WikiService_ServiceDesc.ServiceName,
		Method:  "GetPage",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.Id, userInfo, wiki.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	page, err := s.getPage(ctx, req.Id, true, true, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	if page == nil || page.Id <= 0 {
		return nil, errorswiki.ErrPageNotFound
	}

	if !check && !page.Meta.Public {
		return nil, errorswiki.ErrPageDenied
	}

	resp := &pbwiki.GetPageResponse{
		Page: page,
	}

	s.enricher.EnrichJobName(resp.Page)

	access, err := s.getPageAccess(ctx, userInfo, req.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	resp.Page.Access = access

	files, err := s.fHandler.ListFilesForParentID(ctx, resp.Page.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	resp.Page.Files = files

	auditEntry.State = audit.EventType_EVENT_TYPE_VIEWED

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
	for i := range jobsAccess {
		s.enricher.EnrichJobInfo(jobsAccess[i])
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range usersAccess {
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
	tCreator := tables.User().AS("creator")

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
		tPage.Draft.AS("page_meta.draft"),
	}
	if withContent {
		columns = append(columns,
			tPage.Content.AS("page.content"),
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

	dest := &wiki.Page{}
	if err := stmt.QueryContext(ctx, s.db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.Id == 0 {
		return nil, nil
	}

	s.enricher.EnrichJobName(dest)

	if withAccess {
		access, err := s.getPageAccess(ctx, userInfo, pageId)
		if err != nil {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
		dest.Access = access
	}

	if withContent {
		files, err := s.fHandler.ListFilesForParentID(ctx, pageId)
		if err != nil {
			return nil, err
		}
		dest.Files = files
	}

	return dest, nil
}

func (s *Server) CreatePage(ctx context.Context, req *pbwiki.CreatePageRequest) (*pbwiki.CreatePageResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbwiki.WikiService_ServiceDesc.ServiceName,
		Method:  "CreatePage",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	// No parent ID?
	// If so, check if there are any existing pages for the user's job and use one as the parent.
	if req.ParentId == nil || *req.ParentId <= 0 {
		parentStmt := tPageShort.
			SELECT(
				tPageShort.ID.AS("id"),
			).
			FROM(tPageShort).
			WHERE(jet.AND(
				tPageShort.Job.EQ(jet.String(userInfo.Job)),
				tPageShort.DeletedAt.IS_NULL(),
			)).
			ORDER_BY(tPageShort.ParentID.ASC(), tPageShort.Draft.ASC(), tPageShort.SortKey.ASC()).
			LIMIT(1)

		ids := struct{ ID uint64 }{}
		if err := parentStmt.QueryContext(ctx, s.db, &ids); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
			}
		}

		// Found a potential parent page
		if ids.ID > 0 {
			req.ParentId = &ids.ID
			trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.wiki.parent_id", int64(*req.ParentId)))
		} else {
			req.ParentId = nil
		}
	} else {
		trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.wiki.parent_id", int64(*req.ParentId)))

		p, err := s.getPage(ctx, *req.ParentId, false, false, nil)
		if err != nil {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}

		if p.Job != userInfo.Job {
			return nil, errorswiki.ErrPageDenied
		}

		parentCheck, err := s.access.CanUserAccessTarget(ctx, *req.ParentId, userInfo, wiki.AccessLevel_ACCESS_LEVEL_VIEW)
		if err != nil {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
		if !parentCheck {
			return nil, errorswiki.ErrPageDenied
		}
	}

	job := s.enricher.GetJobByName(userInfo.Job)

	pageAccess := &wiki.PageAccess{
		Jobs: []*wiki.PageJobAccess{
			{
				Job:          userInfo.Job,
				MinimumGrade: userInfo.JobGrade,
				Access:       wiki.AccessLevel_ACCESS_LEVEL_EDIT,
			},
		},
	}
	if job != nil && len(job.Grades) > 0 {
		highestGrade := job.Grades[len(job.Grades)-1]

		if highestGrade.Grade > userInfo.JobGrade {
			// If the user's job grade is lower than the highest grade, add an access entry for the highest grade
			pageAccess.Jobs = append(pageAccess.Jobs, &wiki.PageJobAccess{
				Job:          job.Name,
				MinimumGrade: highestGrade.Grade,
				Access:       wiki.AccessLevel_ACCESS_LEVEL_EDIT,
			})
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	tPage := table.FivenetWikiPages
	stmt := tPage.
		INSERT(
			tPage.Job,
			tPage.ParentID,
			tPage.ContentType,
			tPage.Toc,
			tPage.Draft,
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
			req.ParentId,
			req.ContentType,
			true,
			true,
			false,
			"",
			"",
			"",
			"",
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
	resp := &pbwiki.CreatePageResponse{
		Job: userInfo.Job,
		Id:  uint64(lastId),
	}

	if _, err := s.addPageActivity(ctx, tx, &wiki.PageActivity{
		PageId:       resp.Id,
		ActivityType: wiki.PageActivityType_PAGE_ACTIVITY_TYPE_CREATED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.Job,
	}); err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	if err := s.handlePageAccessChange(ctx, tx, resp.Id, userInfo, pageAccess, false); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_CREATED

	return resp, nil
}

func (s *Server) UpdatePage(ctx context.Context, req *pbwiki.UpdatePageRequest) (*pbwiki.UpdatePageResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.wiki.page_id", int64(req.Page.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbwiki.WikiService_ServiceDesc.ServiceName,
		Method:  "UpdatePage",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	if req.Page.ParentId == nil || *req.Page.ParentId <= 0 {
		stmt := tPage.
			SELECT(
				tPage.ID.AS("id"),
			).
			FROM(tPage).
			WHERE(jet.AND(
				tPage.Job.EQ(jet.String(userInfo.Job)),
				tPage.DeletedAt.IS_NULL(),
				tPage.ParentID.IS_NULL(),
			))

		var ids struct {
			ID uint64 `alias:"id"`
		}
		if err := stmt.QueryContext(ctx, s.db, &ids); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
			}
		}

		if ids.ID != req.Page.Id && !userInfo.Superuser {
			return nil, errorswiki.ErrPageDenied
		}
	} else {
		p, err := s.getPage(ctx, *req.Page.ParentId, false, false, nil)
		if err != nil {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}

		if p.Job != userInfo.Job && !userInfo.Superuser {
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

	if req.Page.ParentId != nil {
		if *req.Page.ParentId == req.Page.Id {
			req.Page.ParentId = nil
		} else if *req.Page.ParentId <= 0 {
			// If the parent ID is not set, we can set it to nil
			req.Page.ParentId = nil
		}
	}

	check, err := s.access.CanUserAccessTarget(ctx, req.Page.Id, userInfo, wiki.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	if !check {
		return nil, errorswiki.ErrPageDenied
	}

	oldPage, err := s.getPage(ctx, req.Page.Id, true, true, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	// A page can only be switched to published once
	if !oldPage.Meta.Draft && oldPage.Meta.Draft != req.Page.Meta.Draft {
		// Allow a super user to change the draft state
		if !userInfo.Superuser {
			req.Page.Meta.Draft = oldPage.Meta.Draft
		}
	}

	// Field Permission Check
	fields, err := s.perms.AttrStringList(userInfo, permswiki.WikiServicePerm, permswiki.WikiServiceUpdatePagePerm, permswiki.WikiServiceUpdatePageFieldsPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	if !fields.Contains("Public") {
		req.Page.Meta.Public = oldPage.Meta.Public
	}

	if req.Page.Access.IsEmpty() {
		// Ensure at least one access entry allowing the user's rank and higher to "edit" the page
		req.Page.Access.Jobs = append(req.Page.Access.Jobs, &wiki.PageJobAccess{
			Job:          userInfo.Job,
			MinimumGrade: userInfo.JobGrade,
			Access:       wiki.AccessLevel_ACCESS_LEVEL_EDIT,
		})
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	tPage := table.FivenetWikiPages
	stmt := tPage.
		UPDATE(
			tPage.ParentID,
			tPage.ContentType,
			tPage.Toc,
			tPage.Draft,
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
			req.Page.Meta.Draft,
			req.Page.Meta.Public,
			slug.Make(utils.StringFirstN(req.Page.Meta.Title, 100)),
			req.Page.Meta.Title,
			req.Page.Meta.Description,
			req.Page.Content,
			nil,
		).
		WHERE(jet.AND(
			tPage.ID.EQ(jet.Uint64(req.Page.Id)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	diff, err := s.generatePageDiff(oldPage, &wiki.Page{
		Meta: &wiki.PageMeta{
			Title:       req.Page.Meta.Title,
			Description: req.Page.Meta.Description,
		},
		Content: req.Page.Content,
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	added, deleted, err := s.fHandler.HandleFileChangesForParent(ctx, tx, req.Page.Id, req.Page.Files)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	if added > 0 || deleted > 0 {
		diff.FilesChange = &wiki.PageFilesChange{
			Added:   added,
			Deleted: deleted,
		}
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

	if oldPage.Meta.Draft != req.Page.Meta.Draft {
		if _, err := s.addPageActivity(ctx, tx, &wiki.PageActivity{
			PageId:       req.Page.Id,
			ActivityType: wiki.PageActivityType_PAGE_ACTIVITY_TYPE_DRAFT_TOGGLED,
			CreatorId:    &userInfo.UserId,
			CreatorJob:   userInfo.Job,
		}); err != nil {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	page, err := s.getPage(ctx, req.Page.Id, true, true, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	return &pbwiki.UpdatePageResponse{
		Page: page,
	}, nil
}

func (s *Server) handlePageAccessChange(ctx context.Context, tx qrm.DB, pageId uint64, userInfo *userinfo.UserInfo, access *wiki.PageAccess, addActivity bool) error {
	changes, err := s.access.HandleAccessChanges(ctx, tx, pageId, access.Jobs, access.Users, nil)
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

func (s *Server) DeletePage(ctx context.Context, req *pbwiki.DeletePageRequest) (*pbwiki.DeletePageResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.wiki.page_id", int64(req.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbwiki.WikiService_ServiceDesc.ServiceName,
		Method:  "DeletePage",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.Id, userInfo, wiki.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	if !check {
		return nil, errorswiki.ErrPageDenied
	}

	page, err := s.getPage(ctx, req.Id, false, false, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	tPage := table.FivenetWikiPages

	// Ensure page has no children
	countStmt := tPage.
		SELECT(
			jet.COUNT(tPage.ID).AS("data_count.total"),
		).
		FROM(tPage).
		WHERE(tPage.ParentID.EQ(jet.Uint64(page.Id)))

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
	}

	if count.Total > 0 {
		return nil, errorswiki.ErrPageHasChildren
	}

	deletedAtTime := jet.CURRENT_TIMESTAMP()
	if page.Meta != nil && page.Meta.DeletedAt != nil && userInfo.Superuser {
		deletedAtTime = jet.TimestampExp(jet.NULL)
	}

	stmt := tPage.
		UPDATE(
			tPage.DeletedAt,
		).
		SET(
			tPage.DeletedAt.SET(deletedAtTime),
		).
		WHERE(jet.AND(
			tPage.ID.EQ(jet.Uint64(req.Id)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return &pbwiki.DeletePageResponse{}, nil
}
