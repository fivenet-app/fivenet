package documents

import (
	context "context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2025/services/documents/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

var tUserActivity = table.FivenetUserActivity

func (s *Server) ListUserDocuments(
	ctx context.Context,
	req *pbdocuments.ListUserDocumentsRequest,
) (*pbdocuments.ListUserDocumentsResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.user_id", req.GetUserId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	userCondition := jet.Bool(true)
	if !userInfo.GetSuperuser() {
		userCondition = jet.OR(
			tDAccess.UserID.EQ(jet.Int32(userInfo.GetUserId())),
			jet.AND(
				tDAccess.Job.EQ(jet.String(userInfo.GetJob())),
				tDAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.GetJobGrade())),
			),
		)
	}

	condition := jet.AND(
		jet.OR(
			tDocRel.SourceUserID.EQ(jet.Int32(req.GetUserId())),
			tDocRel.TargetUserID.EQ(jet.Int32(req.GetUserId())),
		),
		tDocRel.DeletedAt.IS_NULL(),
		tDocument.DeletedAt.IS_NULL(),
		jet.OR(
			tDocument.Public.IS_TRUE(),
			jet.AND(
				tDocument.CreatorID.EQ(jet.Int32(userInfo.GetUserId())),
				tDocument.CreatorJob.EQ(jet.String(userInfo.GetJob())),
			),
			userCondition,
		),
	)

	if req.Closed != nil {
		condition = condition.AND(tDocument.Closed.EQ(
			jet.Bool(req.GetClosed()),
		))
	}
	if len(req.GetRelations()) > 0 {
		types := []jet.Expression{}
		for _, t := range req.GetRelations() {
			types = append(types, jet.Int32(int32(*t.Enum())))
		}

		condition = condition.AND(tDocRel.Relation.IN(types...))
	} else {
		return &pbdocuments.ListUserDocumentsResponse{
			Pagination: &database.PaginationResponse{},
			Relations:  []*documents.DocumentRelation{},
		}, nil
	}

	countStmt := tDocRel.
		SELECT(
			jet.COUNT(jet.DISTINCT(tDocRel.DocumentID)).AS("data_count.total"),
		).
		FROM(
			tDocRel.
				INNER_JOIN(tDocument,
					tDocument.ID.EQ(tDocRel.DocumentID),
				).
				LEFT_JOIN(tDAccess,
					tDAccess.TargetID.EQ(tDocRel.DocumentID).
						AND(tDAccess.Access.GT_EQ(jet.Int32(int32(documents.AccessLevel_ACCESS_LEVEL_VIEW)))),
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count.Total, 20)
	resp := &pbdocuments.ListUserDocumentsResponse{
		Pagination: pag,
		Relations:  []*documents.DocumentRelation{},
	}
	if count.Total <= 0 {
		return resp, nil
	}

	idStmt := tDocRel.
		SELECT(
			tDocRel.ID,
		).
		FROM(
			tDocRel.
				INNER_JOIN(tDocument,
					tDocument.ID.EQ(tDocRel.DocumentID),
				).
				LEFT_JOIN(tDAccess,
					tDAccess.TargetID.EQ(tDocRel.DocumentID).
						AND(tDAccess.Access.GT_EQ(jet.Int32(int32(documents.AccessLevel_ACCESS_LEVEL_VIEW)))),
				),
		).
		WHERE(condition).
		OFFSET(req.GetPagination().GetOffset()).
		ORDER_BY(
			tDocRel.CreatedAt.DESC(),
		).
		GROUP_BY(tDocRel.ID).
		LIMIT(limit)

	var dbRelIds []int64
	if err := idStmt.QueryContext(ctx, s.db, &dbRelIds); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	if len(dbRelIds) == 0 {
		return resp, nil
	}

	// Convert proto sort to db sorting
	orderBys := []jet.OrderByClause{}
	if req.GetSort() != nil {
		var column jet.Column
		switch req.GetSort().GetColumn() {
		case "createdAt":
			fallthrough
		default:
			column = tDocument.CreatedAt
		}

		if req.GetSort().GetDirection() == database.AscSortDirection {
			orderBys = append(orderBys,
				column.ASC(),
			)
		} else {
			orderBys = append(orderBys,
				column.DESC(),
			)
		}
	} else {
		orderBys = append(orderBys, tDocument.CreatedAt.DESC())
	}

	rIds := make([]jet.Expression, len(dbRelIds))
	for i := range dbRelIds {
		rIds[i] = jet.Int64(dbRelIds[i])
	}

	tCreator := tables.User().AS("creator")
	tASource := tCreator.AS("source_user")

	stmt := tDocRel.
		SELECT(
			tDocRel.ID,
			tDocRel.CreatedAt,
			tDocRel.DeletedAt,
			tDocRel.DocumentID,
			tDocRel.SourceUserID,
			tDocument.ID,
			tDocument.CreatedAt,
			tDocument.UpdatedAt,
			tDocument.CategoryID,
			tDocument.CreatorID,
			tDocument.State,
			tDocument.Closed,
			tDocument.Draft,
			tDocument.Public,
			tDocument.Title,
			tDCategory.ID,
			tDCategory.Name,
			tDCategory.Description,
			tDCategory.Color,
			tDCategory.Icon,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tDocRel.SourceUserID,
			tASource.ID,
			tASource.Job,
			tASource.JobGrade,
			tASource.Firstname,
			tASource.Lastname,
			tASource.Dateofbirth,
			tDocRel.Relation,
			tDocRel.TargetUserID,
		).
		FROM(
			tDocRel.
				LEFT_JOIN(tDocument,
					tDocRel.DocumentID.EQ(tDocument.ID),
				).
				LEFT_JOIN(tDCategory,
					tDocument.CategoryID.EQ(tDCategory.ID).
						AND(tDCategory.DeletedAt.IS_NULL()),
				).
				LEFT_JOIN(tCreator,
					tDocument.CreatorID.EQ(tCreator.ID),
				).
				LEFT_JOIN(tASource,
					tASource.ID.EQ(tDocRel.SourceUserID),
				),
		).
		WHERE(jet.AND(
			tDocument.DeletedAt.IS_NULL(),
			tDocRel.ID.IN(rIds...),
		)).
		ORDER_BY(orderBys...).
		GROUP_BY(tDocument.ID).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Relations); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetRelations() {
		if resp.GetRelations()[i].GetSourceUser() != nil {
			jobInfoFn(resp.GetRelations()[i].GetSourceUser())
		}
		if resp.GetRelations()[i].GetDocument() != nil &&
			resp.GetRelations()[i].GetDocument().GetCreator() != nil {
			jobInfoFn(resp.GetRelations()[i].GetDocument().GetCreator())
		}
	}

	resp.GetPagination().Update(len(resp.GetRelations()))

	return resp, nil
}

func (s *Server) addUserActivity(
	ctx context.Context,
	tx qrm.DB,
	userId int32,
	targetUserId int32,
	aType users.UserActivityType,
	reason string,
	data *users.UserActivityData,
) error {
	reasonField := jet.NULL
	if reason != "" {
		reasonField = jet.String(reason)
	}

	stmt := tUserActivity.
		INSERT(
			tUserActivity.SourceUserID,
			tUserActivity.TargetUserID,
			tUserActivity.Type,
			tUserActivity.Reason,
			tUserActivity.Data,
		).
		VALUES(
			userId,
			targetUserId,
			int32(aType),
			reasonField,
			data,
		)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}
