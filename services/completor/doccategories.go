package completor

import (
	context "context"
	"errors"

	pbcompletor "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/completor"
	permsdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorscompletor "github.com/fivenet-app/fivenet/v2026/services/completor/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var tDCategory = table.FivenetDocumentsCategories.AS("category")

func (s *Server) CompleteDocumentCategories(
	ctx context.Context,
	req *pbcompletor.CompleteDocumentCategoriesRequest,
) (*pbcompletor.CompleteDocumentCategoriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	jobs, err := s.ps.AttrJobList(
		userInfo,
		permsdocuments.DocumentsServicePerm,
		permsdocuments.DocumentsServiceListCategoriesPerm,
		permsdocuments.DocumentsServiceListCategoriesJobsPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscompletor.ErrFailedSearch)
	}
	if jobs.Len() == 0 {
		jobs.Strings = append(jobs.Strings, userInfo.GetJob())
	}

	jobsExp := make([]mysql.Expression, jobs.Len())
	for i := range jobs.GetStrings() {
		jobsExp[i] = mysql.String(jobs.GetStrings()[i])
	}

	orderBys := []mysql.OrderByClause{}
	condition := tDCategory.Job.IN(jobsExp...)

	if search := dbutils.PrepareForLikeSearch(req.GetSearch()); search != "" {
		condition = condition.AND(
			tDCategory.Name.LIKE(mysql.String(search)),
		)
	}

	if len(req.GetCategoryIds()) > 0 {
		categoryIds := []mysql.Expression{}
		for _, v := range req.GetCategoryIds() {
			categoryIds = append(categoryIds, mysql.Int64(v))
		}

		// Make sure to sort by the category IDs if provided
		orderBys = append(orderBys, tDCategory.ID.IN(categoryIds...).DESC())
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

	resp := &pbcompletor.CompleteDocumentCategoriesResponse{}
	if err := stmt.QueryContext(ctx, s.db, &resp.Categories); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscompletor.ErrFailedSearch)
		}
	}

	return resp, nil
}
