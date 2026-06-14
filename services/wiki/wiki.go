package wiki

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/content"
	notificationsclientview "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications/clientview"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki"
	wikiaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki/access"
	wikiactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki/activity"
	pbwiki "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/wiki"
	permswiki "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/wiki/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorswiki "github.com/fivenet-app/fivenet/v2026/services/wiki/errors"
	wikistore "github.com/fivenet-app/fivenet/v2026/stores/wiki"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/gosimple/slug"
	logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

var wikiPageSubjectAccessOptions = access.SubjectAccessOptions{
	BlockedAccess: int32(wikiaccess.AccessLevel_ACCESS_LEVEL_BLOCKED),
	DeniedAccessLevels: []int32{
		int32(wikiaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		int32(wikiaccess.AccessLevel_ACCESS_LEVEL_ACCESS),
		int32(wikiaccess.AccessLevel_ACCESS_LEVEL_EDIT),
	},
}

func (s *Server) ListPages(
	ctx context.Context,
	req *pbwiki.ListPagesRequest,
) (*pbwiki.ListPagesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	if req.Job != nil && req.GetJob() == "" {
		*req.Job = userInfo.GetJob()
	}

	tPageShort := table.FivenetWikiPages.AS("page_short")
	subPage := table.FivenetWikiPages.AS("sub_page")
	viewCondition := s.access.ACLAccessExistsCondition(
		tPageShort.ID,
		userInfo,
		int32(wikiaccess.AccessLevel_ACCESS_LEVEL_VIEW),
	)
	rootViewCondition := s.access.ACLAccessExistsCondition(
		subPage.ID,
		userInfo,
		int32(wikiaccess.AccessLevel_ACCESS_LEVEL_VIEW),
	)

	result, err := s.store.ListPages(ctx, wikistore.ListPagesQuery{
		Search:            req.GetSearch(),
		Job:               req.GetJob(),
		RootOnly:          req.GetRootOnly(),
		Superuser:         userInfo.GetSuperuser(),
		UserJob:           userInfo.GetJob(),
		UserID:            userInfo.GetUserId(),
		ViewCondition:     viewCondition,
		RootViewCondition: rootViewCondition,
		Pagination:        req.GetPagination(),
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	for i := range result.Pages {
		s.enricher.EnrichJobName(result.Pages[i])
	}

	resp := &pbwiki.ListPagesResponse{
		Pagination: result.Pagination,
		Pages:      []*wiki.PageShort{},
	}
	if req.Search == nil && (req.RootOnly == nil || !req.GetRootOnly()) {
		resp.Pages = append(resp.Pages, mapPagesToNavItems(result.Pages)...)
	} else {
		resp.Pages = result.Pages
	}

	return resp, nil
}

func (s *Server) GetPage(
	ctx context.Context,
	req *pbwiki.GetPageRequest,
) (*pbwiki.GetPageResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.wiki.page_id", req.GetId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetId(),
		userInfo,
		int32(wikiaccess.AccessLevel_ACCESS_LEVEL_VIEW),
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

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)

	return resp, nil
}

func (s *Server) getPageAccess(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	pageId int64,
) (*wikiaccess.PageAccess, error) {
	access, err := s.access.ListTargetAccess(ctx, s.db, pageId, wikiPageSubjectAccessOptions)
	if err != nil {
		return nil, errorswiki.ErrFailedQuery
	}
	for i := range access.GetJobs() {
		s.enricher.EnrichJobInfo(access.GetJobs()[i])
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range access.GetUsers() {
		if access.GetUsers()[i].GetUser() != nil {
			jobInfoFn(access.GetUsers()[i].GetUser())
		}
	}

	return access, nil
}

func (s *Server) getPage(
	ctx context.Context,
	pageId int64,
	withContent bool,
	withAccess bool,
	userInfo *userinfo.UserInfo,
) (*wiki.Page, error) {
	dest, err := s.store.GetPage(ctx, pageId, withContent)
	if err != nil {
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

	// Check that the user has access to the parent page.
	if req.GetParentId() > 0 {
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

		parentCheck, err := s.access.CanUserAccessTarget(
			ctx,
			req.GetParentId(),
			userInfo,
			int32(wikiaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
		if !parentCheck {
			return nil, errorswiki.ErrPageDenied
		}
	}

	job := s.enricher.GetJobByName(userInfo.GetJob())

	pageAccess := &wikiaccess.PageAccess{
		Jobs: []*wikiaccess.PageJobAccess{
			{
				Job:          userInfo.GetJob(),
				MinimumGrade: userInfo.GetJobGrade(),
				Access:       int32(wikiaccess.AccessLevel_ACCESS_LEVEL_EDIT),
			},
		},
	}
	if job != nil && len(job.GetGrades()) > 0 {
		highestGrade := job.GetGrades()[len(job.GetGrades())-1]

		if highestGrade.GetGrade() > userInfo.GetJobGrade() {
			// If the user's job grade is lower than the highest grade, add an access entry for the highest grade
			pageAccess.Jobs = append(pageAccess.Jobs, &wikiaccess.PageJobAccess{
				Job:          job.GetName(),
				MinimumGrade: highestGrade.GetGrade(),
				Access:       int32(wikiaccess.AccessLevel_ACCESS_LEVEL_EDIT),
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
	var parentID *int64
	if req.GetParentId() > 0 {
		parentID = req.ParentId
	}

	sortRank, err := s.store.NextPageGroupRank(ctx, tx, userInfo.GetJob(), parentID, false, 0)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	stmt := tPage.
		INSERT(
			tPage.Job,
			tPage.ParentID,
			tPage.SortRank,
			tPage.ContentType,
			tPage.Toc,
			tPage.Draft,
			tPage.Public,
			tPage.Startpage,
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
			sortRank,
			int32(content.ContentType_CONTENT_TYPE_TIPTAP_JSON),
			true,
			true,
			false,
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

	if _, err := s.store.AddPageActivity(ctx, tx, &wikiactivity.PageActivity{
		PageId:       resp.GetId(),
		ActivityType: wikiactivity.PageActivityType_PAGE_ACTIVITY_TYPE_CREATED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.GetJob(),
	}); err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	if err := s.handlePageAccessChange(
		ctx,
		tx,
		resp.GetId(),
		userInfo,
		pageAccess,
		false,
	); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

	return resp, nil
}

func (s *Server) UpdatePage(
	ctx context.Context,
	req *pbwiki.UpdatePageRequest,
) (*pbwiki.UpdatePageResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.wiki.page_id", req.GetPage().GetId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetPage().GetId(),
		userInfo,
		int32(wikiaccess.AccessLevel_ACCESS_LEVEL_EDIT),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	if !check {
		return nil, errorswiki.ErrPageDenied
	}

	if req.GetPage().GetParentId() > 0 {
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

		parentCheck, err := s.access.CanUserAccessTarget(
			ctx,
			req.GetPage().GetParentId(),
			userInfo,
			int32(wikiaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
		// Reset parent id to current one
		if !parentCheck {
			*req.Page.ParentId = p.GetParentId()
		}
	}

	// If the parent ID is set to itself or unset, set it to nil
	if req.GetPage().GetParentId() == req.GetPage().GetId() ||
		req.GetPage().GetParentId() <= 0 {
		req.Page.ParentId = nil
	}
	if req.GetPage().GetParentId() > 0 {
		req.Page.Meta.Startpage = false
	}

	oldPage, err := s.getPage(ctx, req.GetPage().GetId(), true, true, userInfo)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, errorswiki.ErrPageNotFound
		}
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	oldOrder, err := s.store.GetPageOrderInfo(ctx, s.db, req.GetPage().GetId())
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
	fields, err := permswiki.WikiService.UpdatePage.FieldsTyped.Get(s.perms, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	if !fields.Contains(permswiki.WikiServiceUpdatePageFieldsPermValuePublic) {
		req.Page.Meta.Public = oldPage.GetMeta().GetPublic()
	}

	if req.GetPage().GetAccess().IsEmpty() {
		// Ensure at least one access entry allowing the user's rank and higher to "edit" the page
		req.Page.Access.Jobs = append(req.Page.Access.Jobs, &wikiaccess.PageJobAccess{
			Job:          userInfo.GetJob(),
			MinimumGrade: userInfo.GetJobGrade(),
			Access:       int32(wikiaccess.AccessLevel_ACCESS_LEVEL_EDIT),
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
	sortRank := oldOrder.SortRank
	newParentID := req.GetPage().ParentId
	groupChanged := func() bool {
		switch {
		case oldOrder.ParentID == nil && newParentID == nil:
			return oldOrder.Startpage != req.GetPage().GetMeta().GetStartpage()
		case oldOrder.ParentID != nil && newParentID != nil:
			return *oldOrder.ParentID != *newParentID
		default:
			return true
		}
	}()
	if groupChanged {
		sortRank, err = s.store.NextPageGroupRank(
			ctx,
			tx,
			userInfo.GetJob(),
			newParentID,
			req.GetPage().GetMeta().GetStartpage(),
			req.GetPage().GetId(),
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
	}

	stmt := tPage.
		UPDATE(
			tPage.ParentID,
			tPage.SortRank,
			tPage.ContentType,
			tPage.Toc,
			tPage.Draft,
			tPage.Public,
			tPage.Startpage,
			tPage.Slug,
			tPage.Title,
			tPage.Description,
			tPage.Content,
			tPage.Data,
		).
		SET(
			req.GetPage().ParentId,
			sortRank,
			int32(req.GetPage().GetContent().GetContentType()),
			req.GetPage().GetMeta().Toc,
			req.GetPage().GetMeta().GetDraft(),
			req.GetPage().GetMeta().GetPublic(),
			req.GetPage().GetMeta().GetStartpage(),
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

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
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
		diff.FilesChange = &wikiactivity.PageFilesChange{
			Added:   added,
			Deleted: deleted,
		}
	}

	// Only store activity if there are actual changes
	if diff.HasChanges() {
		if _, err := s.store.AddPageActivity(ctx, tx, &wikiactivity.PageActivity{
			PageId:       req.GetPage().GetId(),
			ActivityType: wikiactivity.PageActivityType_PAGE_ACTIVITY_TYPE_UPDATED,
			CreatorId:    &userInfo.UserId,
			CreatorJob:   userInfo.GetJob(),
			Data: &wikiactivity.PageActivityData{
				Data: &wikiactivity.PageActivityData_Updated{
					Updated: diff,
				},
			},
		}); err != nil {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
	}

	if err := s.handlePageAccessChange(
		ctx,
		tx,
		req.GetPage().GetId(),
		userInfo,
		req.GetPage().GetAccess(),
		true,
	); err != nil {
		return nil, err
	}

	if oldPage.GetMeta().GetDraft() != req.GetPage().GetMeta().GetDraft() {
		if _, err := s.store.AddPageActivity(ctx, tx, &wikiactivity.PageActivity{
			PageId:       req.GetPage().GetId(),
			ActivityType: wikiactivity.PageActivityType_PAGE_ACTIVITY_TYPE_DRAFT_TOGGLED,
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

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	page, err := s.getPage(ctx, req.GetPage().GetId(), true, true, userInfo)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, errorswiki.ErrPageNotFound
		}
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	s.collabServer.SendTargetSaved(ctx, page.GetId())

	s.notifi.SendObjectEvent(ctx, &notificationsclientview.ObjectEvent{
		Type:      notificationsclientview.ObjectType_OBJECT_TYPE_WIKI_PAGE,
		Id:        &page.Id,
		EventType: notificationsclientview.ObjectEventType_OBJECT_EVENT_TYPE_UPDATED,

		UserId: &userInfo.UserId,
		Job:    &userInfo.Job,
	})

	return &pbwiki.UpdatePageResponse{
		Page: page,
	}, nil
}

func (s *Server) MovePage(
	ctx context.Context,
	req *pbwiki.MovePageRequest,
) (*pbwiki.MovePageResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.wiki.page_id", req.GetPageId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetPageId(),
		userInfo,
		int32(wikiaccess.AccessLevel_ACCESS_LEVEL_EDIT),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	if !check {
		return nil, errorswiki.ErrPageDenied
	}

	if req.GetBeforeId() > 0 && req.GetAfterId() > 0 {
		return nil, errorswiki.ErrFailedQuery
	}

	pageOrder, err := s.store.GetPageOrderInfo(ctx, s.db, req.GetPageId())
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, errorswiki.ErrPageNotFound
		}
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	defer tx.Rollback()

	sortRank, err := s.store.InsertPageGroupRank(
		ctx,
		tx,
		pageOrder.Job,
		pageOrder.ParentID,
		pageOrder.Startpage,
		req.GetPageId(),
		req.BeforeId,
		req.AfterId,
	)
	if err != nil {
		return nil, err
	}

	tPage := table.FivenetWikiPages
	if _, err := tPage.
		UPDATE(
			tPage.SortRank,
		).
		SET(
			sortRank,
		).
		WHERE(mysql.AND(
			tPage.ID.EQ(mysql.Int64(req.GetPageId())),
			tPage.Job.EQ(mysql.String(pageOrder.Job)),
			tPage.DeletedAt.IS_NULL(),
		)).
		LIMIT(1).
		ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	s.collabServer.SendTargetSaved(ctx, req.GetPageId())
	s.notifi.SendObjectEvent(ctx, &notificationsclientview.ObjectEvent{
		Type:      notificationsclientview.ObjectType_OBJECT_TYPE_WIKI_PAGE,
		Id:        &req.PageId,
		EventType: notificationsclientview.ObjectEventType_OBJECT_EVENT_TYPE_UPDATED,

		UserId: &userInfo.UserId,
		Job:    &userInfo.Job,
	})

	return &pbwiki.MovePageResponse{}, nil
}

func (s *Server) handlePageAccessChange(
	ctx context.Context,
	tx qrm.DB,
	pageId int64,
	userInfo *userinfo.UserInfo,
	pageAccess *wikiaccess.PageAccess,
	addActivity bool,
) error {
	changes, err := s.access.ReplaceTargetAccess(
		ctx,
		tx,
		s.accessResolver,
		pageId,
		pageAccess,
		wikiPageSubjectAccessOptions,
	)
	if err != nil {
		if dbutils.IsDuplicateError(err) {
			return errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
		return errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	if addActivity && !changes.IsEmpty() {
		if _, err := s.store.AddPageActivity(ctx, tx, &wikiactivity.PageActivity{
			PageId:       pageId,
			ActivityType: wikiactivity.PageActivityType_PAGE_ACTIVITY_TYPE_ACCESS_UPDATED,
			CreatorId:    &userInfo.UserId,
			CreatorJob:   userInfo.GetJob(),
			Data: &wikiactivity.PageActivityData{
				Data: &wikiactivity.PageActivityData_AccessUpdated{
					AccessUpdated: &wikiactivity.PageAccessUpdated{
						Jobs: &wikiactivity.PageAccessJobsDiff{
							ToCreate: changes.Jobs.ToCreate,
							ToUpdate: changes.Jobs.ToUpdate,
							ToDelete: changes.Jobs.ToDelete,
						},
						Users: &wikiactivity.PageAccessUsersDiff{
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

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetId(),
		userInfo,
		int32(wikiaccess.AccessLevel_ACCESS_LEVEL_EDIT),
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

	// Check if page has any un-deleted child pages
	tPage := table.FivenetWikiPages
	count, err := s.store.CountPageChildren(ctx, page.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	if count > 0 {
		return nil, errorswiki.ErrPageHasChildren
	}

	condition := tPage.ID.EQ(mysql.Int64(req.GetId()))

	var deletedAtTime *timestamp.Timestamp
	if page.GetMeta() == nil || page.GetMeta().GetDeletedAt() == nil || !userInfo.GetSuperuser() {
		deletedAtTime = timestamp.Now()
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	} else {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_RESTORED)

		// Restore the page's parent page if any
		if page.GetParentId() > 0 {
			condition = condition.OR(
				tPage.ID.EQ(mysql.Int64(page.GetParentId())),
			)
		}
	}

	stmt := tPage.
		UPDATE(
			tPage.DeletedAt,
		).
		SET(
			tPage.DeletedAt.SET(dbutils.TimestampToMySQL(deletedAtTime)),
		).
		WHERE(condition).
		LIMIT(2)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)

	return &pbwiki.DeletePageResponse{}, nil
}
