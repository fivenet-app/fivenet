package docstore

import (
	context "context"
	"errors"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	pbdocstore "github.com/fivenet-app/fivenet/gen/go/proto/services/docstore"
	"github.com/fivenet-app/fivenet/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	errorsdocstore "github.com/fivenet-app/fivenet/services/docstore/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tUserActivity = table.FivenetUserActivity

func (s *Server) ListUserDocuments(ctx context.Context, req *pbdocstore.ListUserDocumentsRequest) (*pbdocstore.ListUserDocumentsResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.user_id", int64(req.UserId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	userCondition := jet.Bool(true)
	if !userInfo.SuperUser {
		userCondition = jet.OR(
			tDAccess.UserID.EQ(jet.Int32(userInfo.UserId)),
			jet.AND(
				tDAccess.Job.EQ(jet.String(userInfo.Job)),
				tDAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade)),
			),
		)
	}

	condition := jet.AND(
		jet.OR(
			tDocRel.SourceUserID.EQ(jet.Int32(req.UserId)),
			tDocRel.TargetUserID.EQ(jet.Int32(req.UserId)),
		),
		tDocRel.DeletedAt.IS_NULL(),
		tDocument.DeletedAt.IS_NULL(),
		jet.OR(
			tDocument.Public.IS_TRUE(),
			jet.AND(
				tDocument.CreatorID.EQ(jet.Int32(userInfo.UserId)),
				tDocument.CreatorJob.EQ(jet.String(userInfo.Job)),
			),
			userCondition,
		),
	)

	if req.Closed != nil {
		condition = condition.AND(tDocument.Closed.EQ(
			jet.Bool(*req.Closed),
		))
	}

	countStmt := tDocRel.
		SELECT(
			jet.COUNT(jet.DISTINCT(tDocRel.DocumentID)).AS("datacount.totalcount"),
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
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, 16)
	resp := &pbdocstore.ListUserDocumentsResponse{
		Pagination: pag,
		Relations:  []*documents.DocumentRelation{},
	}
	if count.TotalCount <= 0 {
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
		OFFSET(req.Pagination.Offset).
		ORDER_BY(
			tDocRel.CreatedAt.DESC(),
		).
		GROUP_BY(tDocRel.ID).
		LIMIT(limit)

	var dbRelIds []uint64
	if err := idStmt.QueryContext(ctx, s.db, &dbRelIds); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}
	}

	if len(dbRelIds) <= 0 {
		return resp, nil
	}

	rIds := make([]jet.Expression, len(dbRelIds))
	for i := range dbRelIds {
		rIds[i] = jet.Uint64(dbRelIds[i])
	}

	tCreator := tables.Users().AS("creator")
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
					tDocument.CategoryID.EQ(tDCategory.ID),
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
		ORDER_BY(
			tDocRel.CreatedAt.DESC(),
		).
		GROUP_BY(tDocument.ID).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Relations); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := 0; i < len(resp.Relations); i++ {
		if resp.Relations[i].SourceUser != nil {
			jobInfoFn(resp.Relations[i].SourceUser)
		}
		if resp.Relations[i].Document != nil && resp.Relations[i].Document.Creator != nil {
			jobInfoFn(resp.Relations[i].Document.Creator)
		}
	}

	resp.Pagination.Update(len(resp.Relations))

	return resp, nil
}

func (s *Server) addUserActivity(ctx context.Context, tx qrm.DB, userId int32, targetUserId int32, aType users.UserActivityType, reason string, data *users.UserActivityData) error {
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
			int16(aType),
			reasonField,
			data,
		)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}
