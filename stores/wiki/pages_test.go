package wikistore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreListPages(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)
	pageSize := int64(5)
	query := ListPagesQuery{
		Job:        "police",
		Superuser:  true,
		UserInfo:   &userinfo.UserInfo{Superuser: true},
		Pagination: &database.PaginationRequest{PageSize: &pageSize},
	}

	countQuery := regexp.QuoteMeta(`FROM fivenet_wiki_pages AS page_short`) +
		`(?s).*` + regexp.QuoteMeta(`LEFT JOIN fivenet_job_props ON (fivenet_job_props.job = page_short.job)`) +
		`(?s).*` + regexp.QuoteMeta(`WHERE ( ? AND (? AND (page_short.job = ?)) );`)
	mock.ExpectQuery(countQuery).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(1)))

	listQuery := regexp.QuoteMeta(
		`FROM fivenet_wiki_pages AS page_short LEFT JOIN fivenet_job_props`,
	) +
		`(?s).*` + regexp.QuoteMeta(`WHERE ( ? AND (? AND (page_short.job = ?)) )`) +
		`(?s).*` + regexp.QuoteMeta(
		`ORDER BY page_short.job ASC, page_short.startpage DESC, page_short.parent_id ASC, page_short.sort_rank ASC, page_short.draft ASC, page_short.id ASC LIMIT ? OFFSET ?;`,
	)
	mock.ExpectQuery(listQuery).
		WillReturnRows(sqlmock.NewRows([]string{
			"page_short.id",
			"page_short.job",
			"page_short.parent_id",
			"page_short.slug",
			"page_short.title",
			"page_short.description",
			"page_short.draft",
			"page_short.public",
			"page_short.startpage",
		}).AddRow(int64(7), "police", nil, "page-7", "Page 7", "desc", false, true, false))

	res, err := store.ListPages(t.Context(), query)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.NotNil(t, res.Pagination)
	assert.Equal(t, int64(1), res.Pagination.GetTotalCount())
	require.Len(t, res.Pages, 1)
	assert.Equal(t, int64(7), res.Pages[0].GetId())
	assert.Equal(t, "Page 7", res.Pages[0].GetTitle())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListPagesRootOnly(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)
	pageSize := int64(5)
	query := ListPagesQuery{
		Job:        "police",
		RootOnly:   true,
		Superuser:  false,
		UserInfo:   &userinfo.UserInfo{UserId: 7, Job: "police", JobGrade: 6},
		Pagination: &database.PaginationRequest{PageSize: &pageSize},
	}

	countQuery := `(?s).*WITH user_subjects AS.*visible_sources AS.*winning_visibility AS.*page_short\.id IN.*ranked_pages\.rn = 1`
	mock.ExpectQuery(countQuery).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(1)))

	listQuery := `(?s).*WITH user_subjects AS.*visible_sources AS.*winning_visibility AS.*page_short\.id IN.*ranked_pages\.rn = 1.*` +
		regexp.QuoteMeta(
			`ORDER BY page_short.job ASC, page_short.startpage DESC, page_short.parent_id ASC, page_short.sort_rank ASC, page_short.draft ASC, page_short.id ASC LIMIT ? OFFSET ?;`,
		)
	mock.ExpectQuery(listQuery).
		WillReturnRows(sqlmock.NewRows([]string{
			"page_short.id",
			"page_short.job",
			"page_short.parent_id",
			"page_short.slug",
			"page_short.title",
			"page_short.description",
			"page_short.draft",
			"page_short.public",
			"page_short.startpage",
			"page_root_info.logo_file_id",
			"logo.id",
			"logo.file_path",
		}).AddRow(int64(7), "police", nil, "page-7", "Page 7", "desc", false, true, true, int64(11), int64(11), "/logo.png"))

	res, err := store.ListPages(t.Context(), query)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.Pages, 1)
	assert.Equal(t, int64(7), res.Pages[0].GetId())
	assert.Equal(t, "Page 7", res.Pages[0].GetTitle())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreGetPage(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)
	query := regexp.QuoteMeta(`FROM fivenet_wiki_pages AS page`) +
		`(?s).*` + regexp.QuoteMeta(`LEFT JOIN fivenet_user AS creator ON`) +
		`(?s).*` + regexp.QuoteMeta(`page.id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)

	now := time.Unix(0, 0).UTC()
	mock.ExpectQuery(query).
		WithArgs(int64(42), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"page.id",
			"page.job",
			"page.parent_id",
			"page_meta.created_at",
			"page_meta.updated_at",
			"page_meta.deleted_at",
			"page_meta.slug",
			"page_meta.title",
			"page_meta.description",
			"page_meta.creator_id",
			"creator.id",
			"creator.job",
			"creator.job_grade",
			"creator.firstname",
			"creator.lastname",
			"creator.dateofbirth",
			"page_meta.content_Type",
			"page_meta.toc",
			"page_meta.public",
			"page_meta.draft",
			"page_meta.startpage",
		}).AddRow(
			int64(42),
			"police",
			nil,
			now,
			now,
			nil,
			"page-42",
			"Page 42",
			"Description",
			int32(7),
			int32(7),
			"police",
			int32(5),
			"Jane",
			"Doe",
			now,
			int32(0),
			true,
			true,
			false,
			false,
		))

	page, err := store.GetPage(t.Context(), 42, false)
	require.NoError(t, err)
	require.NotNil(t, page)
	assert.Equal(t, int64(42), page.GetId())
	assert.Equal(t, "Page 42", page.GetMeta().GetTitle())
	assert.Nil(t, page.GetContent())
	require.NoError(t, mock.ExpectationsWereMet())
}
