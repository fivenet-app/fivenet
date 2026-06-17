package wikistore

import (
	"context"
	"errors"
	"strings"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/content"
	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	reswiki "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki"
	wikiaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki/access"
	wikiactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki/activity"
	objectaccess "github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorswiki "github.com/fivenet-app/fivenet/v2026/services/wiki/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/gosimple/slug"
)

const defaultWikiUpperLimit = 250

var wikiPageSubjectAccessOptions = objectaccess.SubjectAccessOptions{
	BlockedAccess: int32(wikiaccess.AccessLevel_ACCESS_LEVEL_BLOCKED),
	DeniedAccessLevels: []int32{
		int32(wikiaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		int32(wikiaccess.AccessLevel_ACCESS_LEVEL_ACCESS),
		int32(wikiaccess.AccessLevel_ACCESS_LEVEL_EDIT),
	},
}

type ListPagesQuery struct {
	Search string
	Job    string

	RootOnly   bool
	Superuser  bool
	UserInfo   *userinfo.UserInfo
	Pagination *database.PaginationRequest
}

type ListPagesResult struct {
	Pagination *database.PaginationResponse
	Pages      []*reswiki.PageShort
}

func (s *Store) ListPages(ctx context.Context, q ListPagesQuery) (*ListPagesResult, error) {
	tPageShort := table.FivenetWikiPages.AS("page_short")
	tPage := table.FivenetWikiPages
	tJobProps := table.FivenetJobProps
	tFiles := table.FivenetFiles.AS("logo")

	columns := mysql.ProjectionList{
		tPageShort.Job,
		tPageShort.ParentID,
		tPageShort.Slug,
		tPageShort.Title,
		tPageShort.Description,
		tPageShort.Draft,
		tPageShort.Public,
		tPageShort.Startpage,
	}

	if q.RootOnly {
		columns = append(columns,
			tJobProps.LogoFileID.AS("page_root_info.logo_file_id"),
			tFiles.ID,
			tFiles.FilePath,
		)
	}

	searchConditionBase := wikiPageSearchCondition(tPage, q.Search)
	searchConditionShort := wikiPageSearchCondition(tPageShort, q.Search)
	jobCondition := mysql.Bool(true)
	if q.Job != "" {
		jobCondition = jobCondition.AND(tPageShort.Job.EQ(mysql.String(q.Job)))
	}
	outerCondition := mysql.AND(searchConditionShort, jobCondition)

	if q.RootOnly {
		return s.listRootPages(
			ctx,
			q,
			tPageShort,
			tJobProps,
			tFiles,
			columns,
			outerCondition,
			searchConditionBase,
		)
	}

	if q.Superuser {
		return s.listPagesSuperuser(
			ctx,
			q,
			tPageShort,
			tJobProps,
			tFiles,
			columns,
			outerCondition,
		)
	}
	return s.listVisiblePages(
		ctx,
		q,
		tPageShort,
		tJobProps,
		tFiles,
		columns,
		outerCondition,
		searchConditionBase,
	)
}

func wikiPageSearchCondition(
	page *table.FivenetWikiPagesTable,
	search string,
) mysql.BoolExpression {
	condition := mysql.Bool(true)
	if search = strings.TrimRight(search, "*"); search != "" {
		search += "*"
		condition = condition.AND(mysql.OR(
			dbutils.MATCH(page.Title, mysql.String(search)),
			dbutils.MATCH(page.Content, mysql.String(search)),
		))
	}
	return condition
}

func (s *Store) listVisiblePages(
	ctx context.Context,
	q ListPagesQuery,
	tPageShort *table.FivenetWikiPagesTable,
	tJobProps *table.FivenetJobPropsTable,
	tFiles *table.FivenetFilesTable,
	columns mysql.ProjectionList,
	outerCondition mysql.BoolExpression,
	searchConditionBase mysql.BoolExpression,
) (*ListPagesResult, error) {
	visibleIDs := s.access.VisibleIDsByConditionQuery(
		q.UserInfo,
		int32(wikiaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		false,
		searchConditionBase,
	)
	ctes := visibleIDs.CTEs
	visiblePageID := mysql.IntegerColumn("id").From(visibleIDs.Table)

	var countStmt mysql.Statement = tPageShort.
		SELECT(mysql.COUNT(visiblePageID).AS("data_count.total")).
		FROM(visibleIDs.Table)

	if len(ctes) > 0 {
		countStmt = mysql.WITH(ctes...)(countStmt)
	}

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	var paginationReq *database.PaginationRequest
	if q.Pagination != nil {
		paginationReq = q.Pagination
	}
	pagination, limit := paginationReq.GetResponseWithPageSize(count.Total, defaultWikiUpperLimit)
	result := &ListPagesResult{Pagination: pagination, Pages: []*reswiki.PageShort{}}
	if count.Total <= 0 {
		return result, nil
	}

	var stmt mysql.Statement = tPageShort.
		SELECT(tPageShort.ID, columns...).
		FROM(
			visibleIDs.Table.
				INNER_JOIN(tPageShort,
					tPageShort.ID.EQ(visiblePageID),
				).
				LEFT_JOIN(tJobProps,
					tJobProps.Job.EQ(tPageShort.Job),
				).
				LEFT_JOIN(tFiles,
					tFiles.ID.EQ(tJobProps.LogoFileID),
				),
		).
		WHERE(outerCondition).
		ORDER_BY(
			tPageShort.Job.ASC(),
			tPageShort.Startpage.DESC(),
			tPageShort.ParentID.ASC().NULLS_FIRST(),
			tPageShort.SortRank.ASC(),
			tPageShort.Draft.ASC(),
			tPageShort.ID.ASC(),
		).
		OFFSET(pagination.GetOffset()).
		LIMIT(limit)

	if len(ctes) > 0 {
		stmt = mysql.WITH(ctes...)(stmt)
	}

	if err := stmt.QueryContext(ctx, s.db, &result.Pages); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return result, nil
}

func (s *Store) listRootPages(
	ctx context.Context,
	q ListPagesQuery,
	tPageShort *table.FivenetWikiPagesTable,
	tJobProps *table.FivenetJobPropsTable,
	tFiles *table.FivenetFilesTable,
	columns mysql.ProjectionList,
	outerCondition mysql.BoolExpression,
	searchConditionBase mysql.BoolExpression,
) (*ListPagesResult, error) {
	subPage := table.FivenetWikiPages.AS("sub_page")
	var rankedPages mysql.SelectTable
	var ctes []mysql.CommonTableExpression

	if q.Superuser {
		rankedPages = subPage.
			SELECT(
				subPage.ID,
				subPage.Job,
				mysql.ROW_NUMBER().OVER(
					mysql.PARTITION_BY(subPage.Job).
						ORDER_BY(
							subPage.Startpage.DESC(),
							subPage.ParentID.ASC().NULLS_FIRST(),
							subPage.SortRank.ASC(),
							subPage.Draft.ASC(),
							subPage.ID.ASC(),
						),
				).AS("rn"),
			).
			FROM(subPage).
			AsTable("ranked_pages")
	} else {
		visibleIDs := s.access.VisibleIDsByConditionQuery(
			q.UserInfo,
			int32(wikiaccess.AccessLevel_ACCESS_LEVEL_VIEW),
			false,
			searchConditionBase,
		)
		ctes = visibleIDs.CTEs
		visiblePageID := mysql.IntegerColumn("id").From(visibleIDs.Table)
		rankedPages = subPage.
			SELECT(
				subPage.ID,
				subPage.Job,
				mysql.ROW_NUMBER().OVER(
					mysql.PARTITION_BY(subPage.Job).
						ORDER_BY(
							subPage.Startpage.DESC(),
							subPage.ParentID.ASC().NULLS_FIRST(),
							subPage.SortRank.ASC(),
							subPage.Draft.ASC(),
							subPage.ID.ASC(),
						),
				).AS("rn"),
			).
			FROM(
				visibleIDs.Table.
					INNER_JOIN(subPage,
						subPage.ID.EQ(visiblePageID),
					),
			).
			AsTable("ranked_pages")
	}

	rootIDs := rankedPages.
		SELECT(mysql.IntegerColumn("sub_page.id")).
		FROM(rankedPages).
		WHERE(mysql.RawBool("ranked_pages.rn = 1"))

	var countStmt mysql.Statement = tPageShort.
		SELECT(mysql.COUNT(tPageShort.ID).AS("data_count.total")).
		FROM(
			tPageShort.
				LEFT_JOIN(tJobProps,
					tJobProps.Job.EQ(tPageShort.Job),
				).
				LEFT_JOIN(tFiles,
					tFiles.ID.EQ(tJobProps.LogoFileID),
				),
		).
		WHERE(mysql.AND(outerCondition, tPageShort.ID.IN(rootIDs)))

	if len(ctes) > 0 {
		countStmt = mysql.WITH(ctes...)(countStmt)
	}

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	var paginationReq *database.PaginationRequest
	if q.Pagination != nil {
		paginationReq = q.Pagination
	}
	pagination, limit := paginationReq.GetResponseWithPageSize(count.Total, defaultWikiUpperLimit)
	result := &ListPagesResult{Pagination: pagination, Pages: []*reswiki.PageShort{}}
	if count.Total <= 0 {
		return result, nil
	}

	var stmt mysql.Statement = tPageShort.
		SELECT(tPageShort.ID, columns...).
		FROM(
			tPageShort.
				LEFT_JOIN(tJobProps,
					tJobProps.Job.EQ(tPageShort.Job),
				).
				LEFT_JOIN(tFiles,
					tFiles.ID.EQ(tJobProps.LogoFileID),
				),
		).
		WHERE(mysql.AND(outerCondition, tPageShort.ID.IN(rootIDs))).
		ORDER_BY(
			tPageShort.Job.ASC(),
			tPageShort.Startpage.DESC(),
			tPageShort.ParentID.ASC().NULLS_FIRST(),
			tPageShort.SortRank.ASC(),
			tPageShort.Draft.ASC(),
			tPageShort.ID.ASC(),
		).
		OFFSET(pagination.GetOffset()).
		LIMIT(limit)

	if len(ctes) > 0 {
		stmt = mysql.WITH(ctes...)(stmt)
	}

	if err := stmt.QueryContext(ctx, s.db, &result.Pages); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return result, nil
}

func (s *Store) listPagesSuperuser(
	ctx context.Context,
	q ListPagesQuery,
	tPageShort *table.FivenetWikiPagesTable,
	tJobProps *table.FivenetJobPropsTable,
	tFiles *table.FivenetFilesTable,
	columns mysql.ProjectionList,
	outerCondition mysql.BoolExpression,
) (*ListPagesResult, error) {
	columns = append(columns, tPageShort.DeletedAt)

	var countStmt mysql.Statement = tPageShort.
		SELECT(mysql.COUNT(tPageShort.ID).AS("data_count.total")).
		FROM(
			tPageShort.
				LEFT_JOIN(tJobProps,
					tJobProps.Job.EQ(tPageShort.Job),
				).
				LEFT_JOIN(tFiles,
					tFiles.ID.EQ(tJobProps.LogoFileID),
				),
		).
		WHERE(outerCondition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	var paginationReq *database.PaginationRequest
	if q.Pagination != nil {
		paginationReq = q.Pagination
	}
	pagination, limit := paginationReq.GetResponseWithPageSize(count.Total, defaultWikiUpperLimit)
	result := &ListPagesResult{Pagination: pagination, Pages: []*reswiki.PageShort{}}
	if count.Total <= 0 {
		return result, nil
	}

	var stmt mysql.Statement = tPageShort.
		SELECT(tPageShort.ID, columns...).
		FROM(
			tPageShort.
				LEFT_JOIN(tJobProps,
					tJobProps.Job.EQ(tPageShort.Job),
				).
				LEFT_JOIN(tFiles,
					tFiles.ID.EQ(tJobProps.LogoFileID),
				),
		).
		WHERE(outerCondition).
		ORDER_BY(
			tPageShort.Job.ASC(),
			tPageShort.Startpage.DESC(),
			tPageShort.ParentID.ASC().NULLS_FIRST(),
			tPageShort.SortRank.ASC(),
			tPageShort.Draft.ASC(),
			tPageShort.ID.ASC(),
		).
		OFFSET(pagination.GetOffset()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &result.Pages); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return result, nil
}

func (s *Store) GetPage(
	ctx context.Context,
	pageID int64,
	withContent bool,
) (*reswiki.Page, error) {
	tPage := table.FivenetWikiPages.AS("page")
	tCreator := table.FivenetUser.AS("creator")

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
		tPage.Startpage.AS("page_meta.startpage"),
	}
	if withContent {
		columns = append(columns,
			tPage.Content.AS("page.content"),
			tPage.Data,
		)
	}

	stmt := tPage.
		SELECT(tPage.ID, columns...).
		FROM(
			tPage.
				LEFT_JOIN(tCreator,
					tPage.CreatorID.EQ(tCreator.ID),
				),
		).
		WHERE(mysql.AND(
			tPage.ID.EQ(mysql.Int64(pageID)),
		)).
		LIMIT(1)

	dest := &reswiki.Page{}
	if err := stmt.QueryContext(ctx, s.db, dest); err != nil {
		return nil, err
	}

	return dest, nil
}

func (s *Store) CreatePage(
	ctx context.Context,
	tx qrm.DB,
	userInfo *userinfo.UserInfo,
	parentID *int64,
	contentType content.ContentType,
	pageAccess *wikiaccess.PageAccess,
) (int64, *wikiaccess.PageAccess, error) {
	normalizedAccess, err := objectaccess.NormalizeAccess(pageAccess, nil, nil, 15)
	if err != nil {
		return 0, nil, err
	}

	var parentPageID *int64
	if parentID != nil && *parentID > 0 {
		parentPageID = parentID
	}

	sortRank, err := s.NextPageGroupRank(ctx, tx, userInfo.GetJob(), parentPageID, false, 0)
	if err != nil {
		return 0, nil, err
	}

	tPage := table.FivenetWikiPages
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
			parentID,
			sortRank,
			int32(contentType),
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
		return 0, nil, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, nil, err
	}

	if _, err := s.AddPageActivity(ctx, tx, &wikiactivity.PageActivity{
		PageId:       lastID,
		ActivityType: wikiactivity.PageActivityType_PAGE_ACTIVITY_TYPE_CREATED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.GetJob(),
	}); err != nil {
		return lastID, nil, err
	}

	if err := s.handlePageAccessChange(
		ctx,
		tx,
		lastID,
		userInfo,
		pageAccess,
		false,
	); err != nil {
		return lastID, nil, err
	}

	if err := s.access.RefreshTargetVisibility(ctx, tx, lastID); err != nil {
		return lastID, nil, err
	}

	return lastID, normalizedAccess, nil
}

func (s *Store) UpdatePage(
	ctx context.Context,
	tx qrm.DB,
	userInfo *userinfo.UserInfo,
	page *reswiki.Page,
	sortRank string,
) (*wikiaccess.PageAccess, error) {
	pageAccess := page.GetAccess()
	if pageAccess == nil || pageAccess.IsEmpty() {
		pageAccess = &wikiaccess.PageAccess{
			Jobs: []*wikiaccess.PageJobAccess{
				{
					Job:          userInfo.GetJob(),
					MinimumGrade: userInfo.GetJobGrade(),
					Access:       int32(wikiaccess.AccessLevel_ACCESS_LEVEL_EDIT),
				},
			},
		}
	}

	normalizedAccess, err := objectaccess.NormalizeAccess(pageAccess, nil, nil, 15)
	if err != nil {
		return nil, err
	}

	tPage := table.FivenetWikiPages
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
			page.ParentId,
			sortRank,
			int32(page.GetContent().GetContentType()),
			page.GetMeta().GetToc(),
			page.GetMeta().GetDraft(),
			page.GetMeta().GetPublic(),
			page.GetMeta().GetStartpage(),
			slug.Make(utils.StringFirstN(page.GetMeta().GetTitle(), 100)),
			page.GetMeta().GetTitle(),
			page.GetMeta().GetDescription(),
			page.GetContent(),
			nil,
		).
		WHERE(mysql.AND(
			tPage.ID.EQ(mysql.Int64(page.GetId())),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, err
	}

	if err := s.handlePageAccessChange(
		ctx,
		tx,
		page.GetId(),
		userInfo,
		pageAccess,
		true,
	); err != nil {
		return nil, err
	}

	if err := s.access.RefreshTargetVisibility(ctx, tx, page.GetId()); err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	return normalizedAccess, nil
}

func (s *Store) handlePageAccessChange(
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
		if _, err := s.AddPageActivity(ctx, tx, &wikiactivity.PageActivity{
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

func (s *Store) DeletePage(
	ctx context.Context,
	tx qrm.DB,
	pageId int64,
	deletedAtTime *timestamp.Timestamp,
	parentId int64,
) error {
	// Check if page has any un-deleted child pages
	tPage := table.FivenetWikiPages

	condition := tPage.ID.EQ(mysql.Int64(pageId))
	// Restore the page's parent page if any
	if deletedAtTime == nil && parentId > 0 {
		condition = condition.OR(
			tPage.ID.EQ(mysql.Int64(parentId)),
		)
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
		return err
	}

	return nil
}
