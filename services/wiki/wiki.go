package wiki

import (
	"context"
	"errors"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/wiki"
	pbwiki "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/wiki"
	permswiki "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/wiki/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorswiki "github.com/fivenet-app/fivenet/v2025/services/wiki/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/gosimple/slug"
	logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) ListPages(
	ctx context.Context,
	req *pbwiki.ListPagesRequest,
) (*pbwiki.ListPagesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbwiki.WikiService_ServiceDesc.ServiceName,
		Method:  "ListPages",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	tPageShort := table.FivenetWikiPages.AS("page_short")
	tPAccess := table.FivenetWikiPagesAccess.AS("access")
	tJobProps := table.FivenetJobProps

	condition := mysql.Bool(true)

	if req.GetRootOnly() {
		condition = condition.AND(tPageShort.ParentID.IS_NULL())
	}
	if req.Search != nil && req.GetSearch() != "" {
		*req.Search = strings.TrimRight(req.GetSearch(), "*") + "*"

		condition = mysql.OR(
			dbutils.MATCH(tPageShort.Title, mysql.String(req.GetSearch())),
			dbutils.MATCH(tPageShort.Content, mysql.String(req.GetSearch())),
		)
	}

	if !userInfo.GetSuperuser() {
		accessExists := mysql.EXISTS(
			mysql.
				SELECT(mysql.Int(1)).
				FROM(tPAccess).
				WHERE(mysql.AND(
					tPAccess.TargetID.EQ(tPageShort.ID),
					tPAccess.Access.IS_NOT_NULL(),
					tPAccess.Access.GT_EQ(
						mysql.Int32(int32(wiki.AccessLevel_ACCESS_LEVEL_VIEW)),
					),

					mysql.OR(
						tPAccess.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
						mysql.AND(
							tPAccess.Job.EQ(mysql.String(userInfo.GetJob())),
							tPAccess.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),
						),
					),
				),
				))

		condition = condition.AND(mysql.AND(
			tPageShort.DeletedAt.IS_NULL(),
			mysql.OR(
				tPageShort.Public.IS_TRUE(),
				mysql.AND(
					tPageShort.Job.EQ(mysql.String(userInfo.GetJob())),
					tPageShort.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
				),
				accessExists,
			),
		))
	}

	if req.Job != nil {
		if req.GetJob() == "" {
			*req.Job = userInfo.GetJob()
		}
		condition = condition.AND(tPageShort.Job.EQ(mysql.String(req.GetJob())))
	}

	countStmt := tPageShort.
		SELECT(
			mysql.COUNT(mysql.DISTINCT(tPageShort.ID)).AS("data_count.total"),
		).
		FROM(tPageShort).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
	}

	pag, _ := req.GetPagination().GetResponseWithPageSize(count.Total, defaultWikiUpperLimit)
	resp := &pbwiki.ListPagesResponse{
		Pagination: pag,
		Pages:      []*wiki.PageShort{},
	}
	if count.Total <= 0 {
		return resp, nil
	}

	tFiles := table.FivenetFiles.AS("logo")

	columns := mysql.ProjectionList{
		tPageShort.Job,
		tPageShort.ParentID,
		tPageShort.Slug,
		tPageShort.Title,
		tPageShort.Description,
		tPageShort.Draft,
		tPageShort.Public,
	}
	if req.GetRootOnly() {
		columns = append(columns,
			tJobProps.LogoFileID.AS("page_root_info.logo_file_id"),
			tFiles.ID,
			tFiles.FilePath,
		)

		subPage := table.FivenetWikiPages.AS("sub_page")

		// Use a subquery to get the first page per job (by min ID)
		condition = condition.AND(tPageShort.ID.IN(
			subPage.
				SELECT(
					mysql.MIN(subPage.ID).AS("min_id"),
				).
				WHERE(subPage.ParentID.IS_NULL()).
				GROUP_BY(subPage.Job)),
		)
	}
	if userInfo.GetSuperuser() {
		columns = append(columns, tPageShort.DeletedAt)
	}

	stmt := tPageShort.
		SELECT(
			tPageShort.ID,
			columns...,
		).
		FROM(
			tPageShort.
				LEFT_JOIN(tJobProps,
					tJobProps.Job.EQ(tPageShort.Job),
				).
				LEFT_JOIN(tFiles,
					tFiles.ID.EQ(tJobProps.LogoFileID),
				),
		).
		WHERE(condition).
		ORDER_BY(tPageShort.ParentID.ASC(), tPageShort.Draft.ASC(), tPageShort.SortKey.ASC()).
		OFFSET(req.GetPagination().GetOffset()).
		LIMIT(defaultWikiUpperLimit)

	pages := []*wiki.PageShort{}
	if err := stmt.QueryContext(ctx, s.db, &pages); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
	}

	for i := range pages {
		s.enricher.EnrichJobName(pages[i])
	}

	if req.Search == nil && (req.RootOnly == nil || !req.GetRootOnly()) {
		for _, page := range mapPagesToNavItems(pages) {
			resp.Pages = append(resp.Pages, page)
		}
	} else {
		resp.Pages = pages
	}

	return resp, nil
}

func (s *Server) GetPage(
	ctx context.Context,
	req *pbwiki.GetPageRequest,
) (*pbwiki.GetPageResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.wiki.page_id", req.GetId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbwiki.WikiService_ServiceDesc.ServiceName,
		Method:  "GetPage",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetId(),
		userInfo,
		wiki.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	page, err := s.getPage(ctx, req.GetId(), true, true, userInfo)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, errorswiki.ErrPageNotFound
		}
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	if page.GetId() <= 0 || page.GetMeta() == nil {
		return nil, errorswiki.ErrPageNotFound
	}

	if !check && !page.GetMeta().GetPublic() {
		return nil, errorswiki.ErrPageDenied
	}

	resp := &pbwiki.GetPageResponse{
		Page: page,
	}

	s.enricher.EnrichJobName(resp.GetPage())

	access, err := s.getPageAccess(ctx, userInfo, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	resp.Page.Access = access

	files, err := s.fHandler.ListFilesForParentID(ctx, resp.GetPage().GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	resp.Page.Files = files

	auditEntry.State = audit.EventType_EVENT_TYPE_VIEWED

	return resp, nil
}

func (s *Server) getPageAccess(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	pageId int64,
) (*wiki.PageAccess, error) {
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
		if usersAccess[i].GetUser() != nil {
			jobInfoFn(usersAccess[i].GetUser())
		}
	}

	return &wiki.PageAccess{
		Jobs:  jobsAccess,
		Users: usersAccess,
	}, nil
}

func (s *Server) getPage(
	ctx context.Context,
	pageId int64,
	withContent bool,
	withAccess bool,
	userInfo *userinfo.UserInfo,
) (*wiki.Page, error) {
	tPage := table.FivenetWikiPages.AS("page")
	tCreator := tables.User().AS("creator")

	columns := mysql.ProjectionList{
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
		WHERE(mysql.AND(
			tPage.ID.EQ(mysql.Int64(pageId)),
		)).
		LIMIT(1)

	dest := &wiki.Page{}
	if err := stmt.QueryContext(ctx, s.db, dest); err != nil {
		return nil, err
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

func (s *Server) CreatePage(
	ctx context.Context,
	req *pbwiki.CreatePageRequest,
) (*pbwiki.CreatePageResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbwiki.WikiService_ServiceDesc.ServiceName,
		Method:  "CreatePage",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	// No parent ID?
	// If so, check if there are any existing pages for the user's job and use one as the parent.
	if req.GetParentId() <= 0 {
		tPageShort := table.FivenetWikiPages.AS("page_short")

		parentStmt := tPageShort.
			SELECT(
				tPageShort.ID.AS("id"),
			).
			FROM(tPageShort).
			WHERE(mysql.AND(
				tPageShort.Job.EQ(mysql.String(userInfo.GetJob())),
				tPageShort.DeletedAt.IS_NULL(),
			)).
			ORDER_BY(tPageShort.ParentID.ASC(), tPageShort.Draft.ASC(), tPageShort.SortKey.ASC()).
			LIMIT(1)

		ids := struct{ ID int64 }{}
		if err := parentStmt.QueryContext(ctx, s.db, &ids); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
			}
		}

		// Found a potential parent page
		if ids.ID > 0 {
			req.ParentId = &ids.ID
			logging.InjectFields(ctx, logging.Fields{"fivenet.wiki.parent_id", req.GetParentId()})
		} else {
			req.ParentId = nil
		}
	} else {
		logging.InjectFields(ctx, logging.Fields{"fivenet.wiki.parent_id", req.GetParentId()})

		p, err := s.getPage(ctx, req.GetParentId(), false, false, nil)
		if err != nil {
			if errors.Is(err, qrm.ErrNoRows) {
				return nil, errorswiki.ErrPageNotFound
			}
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}

		if p.GetJob() != userInfo.GetJob() {
			return nil, errorswiki.ErrPageDenied
		}

		parentCheck, err := s.access.CanUserAccessTarget(ctx, req.GetParentId(), userInfo, wiki.AccessLevel_ACCESS_LEVEL_VIEW)
		if err != nil {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
		if !parentCheck {
			return nil, errorswiki.ErrPageDenied
		}
	}

	job := s.enricher.GetJobByName(userInfo.GetJob())

	pageAccess := &wiki.PageAccess{
		Jobs: []*wiki.PageJobAccess{
			{
				Job:          userInfo.GetJob(),
				MinimumGrade: userInfo.GetJobGrade(),
				Access:       wiki.AccessLevel_ACCESS_LEVEL_EDIT,
			},
		},
	}
	if job != nil && len(job.GetGrades()) > 0 {
		highestGrade := job.GetGrades()[len(job.GetGrades())-1]

		if highestGrade.GetGrade() > userInfo.GetJobGrade() {
			// If the user's job grade is lower than the highest grade, add an access entry for the highest grade
			pageAccess.Jobs = append(pageAccess.Jobs, &wiki.PageJobAccess{
				Job:          job.GetName(),
				MinimumGrade: highestGrade.GetGrade(),
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
			userInfo.GetJob(),
			req.ParentId,
			req.GetContentType(),
			true,
			true,
			false,
			"",
			"",
			"",
			"",
			nil,
			userInfo.GetUserId(),
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
		Job: userInfo.GetJob(),
		Id:  lastId,
	}

	if _, err := s.addPageActivity(ctx, tx, &wiki.PageActivity{
		PageId:       resp.GetId(),
		ActivityType: wiki.PageActivityType_PAGE_ACTIVITY_TYPE_CREATED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.GetJob(),
	}); err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	if err := s.handlePageAccessChange(ctx, tx, resp.GetId(), userInfo, pageAccess, false); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_CREATED

	return resp, nil
}

func (s *Server) UpdatePage(
	ctx context.Context,
	req *pbwiki.UpdatePageRequest,
) (*pbwiki.UpdatePageResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.wiki.page_id", req.GetPage().GetId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbwiki.WikiService_ServiceDesc.ServiceName,
		Method:  "UpdatePage",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetPage().GetId(),
		userInfo,
		wiki.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	if !check {
		return nil, errorswiki.ErrPageDenied
	}

	if req.Page.ParentId == nil || req.GetPage().GetParentId() <= 0 {
		tPage := table.FivenetWikiPages.AS("page")

		stmt := tPage.
			SELECT(
				tPage.ID.AS("id"),
			).
			FROM(tPage).
			WHERE(mysql.AND(
				tPage.Job.EQ(mysql.String(userInfo.GetJob())),
				tPage.DeletedAt.IS_NULL(),
				tPage.ParentID.IS_NULL(),
			))

		var ids struct {
			ID int64 `alias:"id"`
		}
		if err := stmt.QueryContext(ctx, s.db, &ids); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
			}
		}

		if ids.ID != req.GetPage().GetId() && !userInfo.GetSuperuser() {
			return nil, errorswiki.ErrPageDenied
		}
	} else {
		p, err := s.getPage(ctx, req.GetPage().GetParentId(), false, false, nil)
		if err != nil {
			if errors.Is(err, qrm.ErrNoRows) {
				return nil, errorswiki.ErrPageNotFound
			}
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}

		if p.GetJob() != userInfo.GetJob() && !userInfo.GetSuperuser() {
			return nil, errorswiki.ErrPageDenied
		}

		parentCheck, err := s.access.CanUserAccessTarget(ctx, req.GetPage().GetParentId(), userInfo, wiki.AccessLevel_ACCESS_LEVEL_VIEW)
		if err != nil {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
		// Reset parent id to current one
		if !parentCheck {
			*req.Page.ParentId = p.GetParentId()
		}
	}

	if req.Page.ParentId != nil {
		if req.GetPage().GetParentId() == req.GetPage().GetId() {
			req.Page.ParentId = nil
		} else if req.GetPage().GetParentId() <= 0 {
			// If the parent ID is not set, we can set it to nil
			req.Page.ParentId = nil
		}
	}

	oldPage, err := s.getPage(ctx, req.GetPage().GetId(), true, true, userInfo)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, errorswiki.ErrPageNotFound
		}
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	// A page can only be switched to published once
	if !oldPage.GetMeta().GetDraft() &&
		oldPage.GetMeta().GetDraft() != req.GetPage().GetMeta().GetDraft() {
		// Allow a super user to change the draft state
		if !userInfo.GetSuperuser() {
			req.Page.Meta.Draft = oldPage.GetMeta().GetDraft()
		}
	}

	// Field Permission Check
	fields, err := s.perms.AttrStringList(
		userInfo,
		permswiki.WikiServicePerm,
		permswiki.WikiServiceUpdatePagePerm,
		permswiki.WikiServiceUpdatePageFieldsPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	if !fields.Contains("Public") {
		req.Page.Meta.Public = oldPage.GetMeta().GetPublic()
	}

	if req.GetPage().GetAccess().IsEmpty() {
		// Ensure at least one access entry allowing the user's rank and higher to "edit" the page
		req.Page.Access.Jobs = append(req.Page.Access.Jobs, &wiki.PageJobAccess{
			Job:          userInfo.GetJob(),
			MinimumGrade: userInfo.GetJobGrade(),
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
			req.GetPage().ParentId,
			req.GetPage().GetMeta().GetContentType(),
			req.GetPage().GetMeta().Toc,
			req.GetPage().GetMeta().GetDraft(),
			req.GetPage().GetMeta().GetPublic(),
			slug.Make(utils.StringFirstN(req.GetPage().GetMeta().GetTitle(), 100)),
			req.GetPage().GetMeta().GetTitle(),
			req.GetPage().GetMeta().GetDescription(),
			req.GetPage().GetContent(),
			nil,
		).
		WHERE(mysql.AND(
			tPage.ID.EQ(mysql.Int64(req.GetPage().GetId())),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	diff, err := s.generatePageDiff(oldPage, &wiki.Page{
		Meta: &wiki.PageMeta{
			Title:       req.GetPage().GetMeta().GetTitle(),
			Description: req.GetPage().GetMeta().GetDescription(),
		},
		Content: req.GetPage().GetContent(),
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	added, deleted, err := s.fHandler.HandleFileChangesForParent(
		ctx,
		tx,
		req.GetPage().GetId(),
		req.GetPage().GetFiles(),
	)
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
		PageId:       req.GetPage().GetId(),
		ActivityType: wiki.PageActivityType_PAGE_ACTIVITY_TYPE_UPDATED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.GetJob(),
		Data: &wiki.PageActivityData{
			Data: &wiki.PageActivityData_Updated{
				Updated: diff,
			},
		},
	}); err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	if err := s.handlePageAccessChange(ctx, tx, req.GetPage().GetId(), userInfo, req.GetPage().GetAccess(), true); err != nil {
		return nil, err
	}

	if oldPage.GetMeta().GetDraft() != req.GetPage().GetMeta().GetDraft() {
		if _, err := s.addPageActivity(ctx, tx, &wiki.PageActivity{
			PageId:       req.GetPage().GetId(),
			ActivityType: wiki.PageActivityType_PAGE_ACTIVITY_TYPE_DRAFT_TOGGLED,
			CreatorId:    &userInfo.UserId,
			CreatorJob:   userInfo.GetJob(),
		}); err != nil {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	page, err := s.getPage(ctx, req.GetPage().GetId(), true, true, userInfo)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, errorswiki.ErrPageNotFound
		}
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	s.collabServer.SendTargetSaved(ctx, page.GetId())

	s.notifi.SendObjectEvent(ctx, &notifications.ObjectEvent{
		Type:      notifications.ObjectType_OBJECT_TYPE_WIKI_PAGE,
		Id:        &page.Id,
		EventType: notifications.ObjectEventType_OBJECT_EVENT_TYPE_UPDATED,

		UserId: &userInfo.UserId,
		Job:    &userInfo.Job,
	})

	return &pbwiki.UpdatePageResponse{
		Page: page,
	}, nil
}

func (s *Server) handlePageAccessChange(
	ctx context.Context,
	tx qrm.DB,
	pageId int64,
	userInfo *userinfo.UserInfo,
	access *wiki.PageAccess,
	addActivity bool,
) error {
	changes, err := s.access.HandleAccessChanges(
		ctx,
		tx,
		pageId,
		access.GetJobs(),
		access.GetUsers(),
		nil,
	)
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
			CreatorJob:   userInfo.GetJob(),
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

func (s *Server) DeletePage(
	ctx context.Context,
	req *pbwiki.DeletePageRequest,
) (*pbwiki.DeletePageResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.wiki.page_id", req.GetId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbwiki.WikiService_ServiceDesc.ServiceName,
		Method:  "DeletePage",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetId(),
		userInfo,
		wiki.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	if !check {
		return nil, errorswiki.ErrPageDenied
	}

	page, err := s.getPage(ctx, req.GetId(), false, false, userInfo)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, errorswiki.ErrPageNotFound
		}
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	tPage := table.FivenetWikiPages

	// Ensure page has no children
	countStmt := tPage.
		SELECT(
			mysql.COUNT(tPage.ID).AS("data_count.total"),
		).
		FROM(tPage).
		WHERE(tPage.ParentID.EQ(mysql.Int64(page.GetId())))

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
	}

	if count.Total > 0 {
		return nil, errorswiki.ErrPageHasChildren
	}

	deletedAtTime := mysql.CURRENT_TIMESTAMP()
	if page.GetMeta() != nil && page.GetMeta().GetDeletedAt() != nil && userInfo.GetSuperuser() {
		deletedAtTime = mysql.TimestampExp(mysql.NULL)
	}

	stmt := tPage.
		UPDATE(
			tPage.DeletedAt,
		).
		SET(
			tPage.DeletedAt.SET(deletedAtTime),
		).
		WHERE(mysql.AND(
			tPage.ID.EQ(mysql.Int64(req.GetId())),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return &pbwiki.DeletePageResponse{}, nil
}
