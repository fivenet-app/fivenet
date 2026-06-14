package completor

import (
	"context"
	"errors"

	documentscategory "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/category"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type DocumentCategoriesQuery struct {
	Search      string
	CategoryIDs []int64
	Jobs        []string
	CurrentJob  string
}

func (s *Store) CompleteDocumentCategories(
	ctx context.Context,
	q DocumentCategoriesQuery,
) ([]*documentscategory.Category, error) {
	tDCategory := table.FivenetDocumentsCategories.AS("category")

	jobsExp := make([]mysql.Expression, len(q.Jobs))
	for i := range q.Jobs {
		jobsExp[i] = mysql.String(q.Jobs[i])
	}

	orderBys := []mysql.OrderByClause{
		tDCategory.Job.EQ(mysql.String(q.CurrentJob)).DESC(),
	}
	condition := tDCategory.Job.IN(jobsExp...)

	if search := dbutils.PrepareForLikeSearch(q.Search); search != "" {
		condition = condition.AND(tDCategory.Name.LIKE(mysql.String(search)))
	}

	if len(q.CategoryIDs) > 0 {
		categoryIDs := []mysql.Expression{}
		for _, v := range q.CategoryIDs {
			categoryIDs = append(categoryIDs, mysql.Int64(v))
		}

		// Make sure to sort by the category IDs if provided
		orderBys = append([]mysql.OrderByClause{
			tDCategory.ID.IN(categoryIDs...).DESC(),
		}, orderBys...)
	}

	orderBys = append(orderBys, tDCategory.SortKey.ASC())

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
		ORDER_BY(orderBys...).
		LIMIT(15)

	dest := []*documentscategory.Category{}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}
