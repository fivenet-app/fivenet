package wikistore

import (
	"context"
	"errors"
	"strings"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	reswiki "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki"
	wikiaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

const defaultWikiUpperLimit = 250

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
	if q.UserInfo == nil {
		q.UserInfo = &userinfo.UserInfo{}
	}

	tPageShort := table.FivenetWikiPages.AS("page_short")
	tJobProps := table.FivenetJobProps
	tFiles := table.FivenetFiles.AS("logo")

	condition := mysql.Bool(true)
	if search := strings.TrimRight(q.Search, "*"); search != "" {
		search += "*"
		condition = condition.AND(mysql.OR(
			dbutils.MATCH(tPageShort.Title, mysql.String(search)),
			dbutils.MATCH(tPageShort.Content, mysql.String(search)),
		))
	}

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

	visibleQuery := s.access.VisibleIDsByConditionQuery(
		q.UserInfo,
		int32(wikiaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		table.FivenetWikiPages.DeletedAt.IS_NULL(),
	)
	visibleIDs := mysql.
		SELECT(mysql.IntegerColumn("id").From(visibleQuery.Table)).
		FROM(visibleQuery.Table)

	if q.RootOnly {
		columns = append(columns,
			tJobProps.LogoFileID.AS("page_root_info.logo_file_id"),
			tFiles.ID,
			tFiles.FilePath,
		)

		rootVisibleIDs := visibleIDs

		subPage := table.FivenetWikiPages.AS("sub_page")
		rootCondition := mysql.AND(subPage.DeletedAt.IS_NULL())
		rootCondition = rootCondition.AND(subPage.ID.IN(rootVisibleIDs))

		rankedPages := subPage.
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
			WHERE(rootCondition).
			AsTable("ranked_pages")

		condition = condition.AND(tPageShort.ID.IN(
			rankedPages.
				SELECT(mysql.IntegerColumn("sub_page.id")).
				FROM(rankedPages).
				WHERE(mysql.RawBool("ranked_pages.rn = 1")),
		))
	} else if !q.Superuser {
		condition = condition.AND(mysql.AND(
			tPageShort.DeletedAt.IS_NULL(),
			tPageShort.ID.IN(visibleIDs),
		))
	}

	if q.Job != "" {
		condition = condition.AND(tPageShort.Job.EQ(mysql.String(q.Job)))
	}

	if q.Superuser {
		columns = append(columns, tPageShort.DeletedAt)
	}

	var countStmt mysql.Statement = tPageShort.
		SELECT(mysql.COUNT(tPageShort.ID).AS("data_count.total")).
		FROM(tPageShort).
		WHERE(condition)
	if len(visibleQuery.CTEs) > 0 {
		countStmt = mysql.WITH(visibleQuery.CTEs...)(countStmt)
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
		WHERE(condition).
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
	if len(visibleQuery.CTEs) > 0 {
		stmt = mysql.WITH(visibleQuery.CTEs...)(stmt)
	}

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
