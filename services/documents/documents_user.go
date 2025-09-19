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
	"github.com/go-jet/jet/v2/mysql"
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

	var userCondition mysql.BoolExpression
	if !userInfo.GetSuperuser() {
		userCondition = mysql.EXISTS(
			mysql.
				SELECT(mysql.Int(1)).
				FROM(tDAccess).
				WHERE(
					mysql.AND(
						tDAccess.TargetID.EQ(tDocument.ID),
						mysql.OR(
							// Direct user access
							tDAccess.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
							// or job + grade access
							mysql.AND(
								tDAccess.Job.EQ(mysql.String(userInfo.GetJob())),
								tDAccess.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),
							),
						),
						tDAccess.Access.GT_EQ(
							mysql.Int32(int32(documents.AccessLevel_ACCESS_LEVEL_VIEW)),
						),
					),
				),
		)
	} else {
		userCondition = mysql.Bool(true)
	}

	condition := mysql.AND(
		mysql.OR(
			tDocRel.SourceUserID.EQ(mysql.Int32(req.GetUserId())),
			tDocRel.TargetUserID.EQ(mysql.Int32(req.GetUserId())),
		),
		tDocRel.DeletedAt.IS_NULL(),
		tDocument.DeletedAt.IS_NULL(),
		mysql.OR(
			tDocument.Public.IS_TRUE(),
			mysql.AND(
				tDocument.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
				tDocument.CreatorJob.EQ(mysql.String(userInfo.GetJob())),
			),
			userCondition,
		),
	)

	if req.Closed != nil {
		condition = condition.AND(tDocument.Closed.EQ(
			mysql.Bool(req.GetClosed()),
		))
	}
	if len(req.GetRelations()) > 0 {
		types := []mysql.Expression{}
		for _, t := range req.GetRelations() {
			types = append(types, mysql.Int32(int32(*t.Enum())))
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
			mysql.COUNT(mysql.DISTINCT(tDocRel.DocumentID)).AS("data_count.total"),
		).
		FROM(
			tDocRel.
				INNER_JOIN(tDocument,
					tDocument.ID.EQ(tDocRel.DocumentID),
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

	// Convert proto sort to db sorting
	orderBys := []mysql.OrderByClause{}
	if req.GetSort() != nil {
		var column mysql.Column
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

	tCreator := tables.User().AS("creator")
	tASource := tCreator.AS("source_user")

	docRel := tDocRel.
		SELECT(
			tDocRel.ID.AS("id"),
			tDocRel.DocumentID.AS("document_id"),
		).
		FROM(
			tDocRel.
				INNER_JOIN(tDocument,
					tDocument.ID.EQ(tDocRel.DocumentID),
				),
		).
		WHERE(condition).
		OFFSET(req.GetPagination().GetOffset()).
		ORDER_BY(tDocRel.CreatedAt.DESC()).
		LIMIT(limit).AsTable("doc_rel")

	stmt := docRel.
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
			docRel.
				INNER_JOIN(tDocRel,
					tDocRel.ID.EQ(mysql.RawInt("doc_rel.id")),
				).
				INNER_JOIN(tDocument,
					tDocument.ID.EQ(mysql.RawInt("doc_rel.document_id")),
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
		WHERE(mysql.AND(
			tDocument.DeletedAt.IS_NULL(),
		)).
		ORDER_BY(orderBys...).
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
	reasonField := mysql.NULL
	if reason != "" {
		reasonField = mysql.String(reason)
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
