package docstore

import (
	context "context"
	"errors"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	errorsdocstore "github.com/fivenet-app/fivenet/gen/go/proto/services/docstore/errors"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	tUserAct = table.FivenetUserActivity
)

func (s *Server) ListUserDocuments(ctx context.Context, req *ListUserDocumentsRequest) (*ListUserDocumentsResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.user_id", int64(req.UserId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := jet.AND(
		tDocRel.DeletedAt.IS_NULL(),
		jet.OR(
			tDocRel.SourceUserID.EQ(jet.Int32(req.UserId)),
			tDocRel.TargetUserID.EQ(jet.Int32(req.UserId)),
		),
		tDocument.DeletedAt.IS_NULL(),
		jet.OR(
			jet.OR(
				tDocument.Public.IS_TRUE(),
				tDocument.CreatorID.EQ(jet.Int32(userInfo.UserId)),
			),
			jet.OR(
				jet.AND(
					tDUserAccess.Access.IS_NOT_NULL(),
					tDUserAccess.Access.GT_EQ(jet.Int32(int32(documents.AccessLevel_ACCESS_LEVEL_VIEW))),
				),
				jet.AND(
					tDUserAccess.Access.IS_NULL(),
					tDJobAccess.Access.IS_NOT_NULL(),
					tDJobAccess.Access.GT_EQ(jet.Int32(int32(documents.AccessLevel_ACCESS_LEVEL_VIEW))),
				),
			),
		),
	)

	countStmt := tDocRel.
		SELECT(
			jet.COUNT(jet.DISTINCT(tDocRel.DocumentID)).AS("datacount.totalcount"),
		).
		FROM(
			tDocRel.
				INNER_JOIN(tDocument,
					tDocument.ID.EQ(tDocRel.DocumentID),
				).
				LEFT_JOIN(tDUserAccess,
					tDUserAccess.DocumentID.EQ(tDocument.ID).
						AND(tDUserAccess.UserID.EQ(jet.Int32(userInfo.UserId)))).
				LEFT_JOIN(tDJobAccess,
					tDJobAccess.DocumentID.EQ(tDocument.ID).
						AND(tDJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tDJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
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
	resp := &ListUserDocumentsResponse{
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
				LEFT_JOIN(tDUserAccess,
					tDUserAccess.DocumentID.EQ(tDocument.ID).
						AND(tDUserAccess.UserID.EQ(jet.Int32(userInfo.UserId)))).
				LEFT_JOIN(tDJobAccess,
					tDJobAccess.DocumentID.EQ(tDocument.ID).
						AND(tDJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tDJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
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
	for i := 0; i < len(dbRelIds); i++ {
		rIds[i] = jet.Uint64(dbRelIds[i])
	}

	tDCreator := tUsers.AS("creator")
	tASource := tUsers.AS("source_user")
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
			tDCreator.ID,
			tDCreator.Identifier,
			tDCreator.Job,
			tDCreator.JobGrade,
			tDCreator.Firstname,
			tDCreator.Lastname,
			tDCreator.Dateofbirth,
			tDocRel.SourceUserID,
			tASource.ID,
			tASource.Identifier,
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
				LEFT_JOIN(tDCreator,
					tDocument.CreatorID.EQ(tDCreator.ID),
				).
				LEFT_JOIN(tASource,
					tASource.ID.EQ(tDocRel.SourceUserID),
				),
		).
		WHERE(
			jet.AND(
				tDocument.DeletedAt.IS_NULL(),
				tDocRel.ID.IN(rIds...),
			),
		).
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

func (s *Server) addUserActivity(ctx context.Context, tx qrm.DB, userId int32, targetUserId int32, activityType users.UserActivityType, key string, oldValue string, newValue string, reason string) error {
	reasonField := jet.NULL
	if reason != "" {
		reasonField = jet.String(reason)
	}

	stmt := tUserAct.
		INSERT(
			tUserAct.SourceUserID,
			tUserAct.TargetUserID,
			tUserAct.Type,
			tUserAct.Key,
			tUserAct.OldValue,
			tUserAct.NewValue,
			tUserAct.Reason,
		).
		VALUES(
			userId,
			targetUserId,
			int16(activityType),
			key,
			&oldValue,
			&newValue,
			reasonField,
		)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}
