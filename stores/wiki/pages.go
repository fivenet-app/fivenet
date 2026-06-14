package wikistore

import (
	"context"
	"errors"
	"strings"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	reswiki "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

const defaultWikiUpperLimit = 250

type ListPagesQuery struct {
	Search string
	Job    string

	RootOnly  bool
	Superuser bool
	UserJob   string
	UserID    int32

	ViewCondition     mysql.BoolExpression
	RootViewCondition mysql.BoolExpression
	Pagination        *database.PaginationRequest
}

type ListPagesResult struct {
	Pagination *database.PaginationResponse
	Pages      []*reswiki.PageShort
}

func (s *Store) ListPages(ctx context.Context, q ListPagesQuery) (*ListPagesResult, error) {
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

	if q.RootOnly {
		columns = append(columns,
			tJobProps.LogoFileID.AS("page_root_info.logo_file_id"),
			tFiles.ID,
			tFiles.FilePath,
		)

		subPage := table.FivenetWikiPages.AS("sub_page")
		rootCondition := mysql.AND(subPage.DeletedAt.IS_NULL())
		if q.RootViewCondition != nil {
			rootCondition = rootCondition.AND(mysql.OR(
				subPage.Public.IS_TRUE(),
				mysql.AND(
					subPage.Job.EQ(mysql.String(q.UserJob)),
					subPage.CreatorID.EQ(mysql.Int32(q.UserID)),
				),
				q.RootViewCondition,
			))
		} else {
			rootCondition = rootCondition.AND(mysql.OR(
				subPage.Public.IS_TRUE(),
				mysql.AND(
					subPage.Job.EQ(mysql.String(q.UserJob)),
					subPage.CreatorID.EQ(mysql.Int32(q.UserID)),
				),
			))
		}

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
		pageCondition := mysql.AND(tPageShort.DeletedAt.IS_NULL())
		if q.ViewCondition != nil {
			pageCondition = pageCondition.AND(mysql.OR(
				tPageShort.Public.IS_TRUE(),
				mysql.AND(
					tPageShort.Job.EQ(mysql.String(q.UserJob)),
					tPageShort.CreatorID.EQ(mysql.Int32(q.UserID)),
				),
				q.ViewCondition,
			))
		} else {
			pageCondition = pageCondition.AND(mysql.OR(
				tPageShort.Public.IS_TRUE(),
				mysql.AND(
					tPageShort.Job.EQ(mysql.String(q.UserJob)),
					tPageShort.CreatorID.EQ(mysql.Int32(q.UserID)),
				),
			))
		}
		condition = condition.AND(pageCondition)
	}

	if q.Job != "" {
		condition = condition.AND(tPageShort.Job.EQ(mysql.String(q.Job)))
	}

	if q.Superuser {
		columns = append(columns, tPageShort.DeletedAt)
	}

	countStmt := tPageShort.
		SELECT(mysql.COUNT(tPageShort.ID).AS("data_count.total")).
		FROM(tPageShort).
		WHERE(condition)

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

	stmt := tPageShort.
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
