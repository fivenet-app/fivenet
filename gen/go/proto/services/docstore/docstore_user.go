package docstore

import (
	context "context"
	"database/sql"
	"errors"

	"github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	"github.com/galexrt/fivenet/gen/go/proto/resources/documents"
	"github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tUserAct = table.FivenetUserActivity
)

func (s *Server) ListUserDocuments(ctx context.Context, req *ListUserDocumentsRequest) (*ListUserDocumentsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := jet.AND(
		tDocRel.DeletedAt.IS_NULL(),
		jet.OR(
			tDocRel.SourceUserID.EQ(jet.Int32(req.UserId)),
			tDocRel.TargetUserID.EQ(jet.Int32(req.UserId)),
		),
		tDocs.DeletedAt.IS_NULL(),
		jet.OR(
			jet.OR(
				tDocs.Public.IS_TRUE(),
				tDocs.CreatorID.EQ(jet.Int32(userInfo.UserId)),
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
				INNER_JOIN(tDocs,
					tDocs.ID.EQ(tDocRel.DocumentID),
				).
				LEFT_JOIN(tDUserAccess,
					tDUserAccess.DocumentID.EQ(tDocs.ID).
						AND(tDUserAccess.UserID.EQ(jet.Int32(userInfo.UserId)))).
				LEFT_JOIN(tDJobAccess,
					tDJobAccess.DocumentID.EQ(tDocs.ID).
						AND(tDJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tDJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, ErrFailedQuery
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(15)
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
				INNER_JOIN(tDocs,
					tDocs.ID.EQ(tDocRel.DocumentID),
				).
				LEFT_JOIN(tDUserAccess,
					tDUserAccess.DocumentID.EQ(tDocs.ID).
						AND(tDUserAccess.UserID.EQ(jet.Int32(userInfo.UserId)))).
				LEFT_JOIN(tDJobAccess,
					tDJobAccess.DocumentID.EQ(tDocs.ID).
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
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, ErrFailedQuery
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
			tDocs.ID,
			tDocs.CreatedAt,
			tDocs.UpdatedAt,
			tDocs.CategoryID,
			tDocs.CreatorID,
			tDocs.State,
			tDocs.Closed,
			tDocs.Title,
			tDCategory.ID,
			tDCategory.Name,
			tDCategory.Description,
			tDCreator.ID,
			tDCreator.Identifier,
			tDCreator.Job,
			tDCreator.JobGrade,
			tDCreator.Firstname,
			tDCreator.Lastname,
			tDocRel.SourceUserID,
			tASource.ID,
			tASource.Identifier,
			tASource.Job,
			tASource.JobGrade,
			tASource.Firstname,
			tASource.Lastname,
			tDocRel.Relation,
			tDocRel.TargetUserID,
		).
		FROM(
			tDocRel.
				LEFT_JOIN(tDocs,
					tDocRel.DocumentID.EQ(tDocs.ID),
				).
				LEFT_JOIN(tDCategory,
					tDocs.CategoryID.EQ(tDCategory.ID),
				).
				LEFT_JOIN(tDCreator,
					tDocs.CreatorID.EQ(tDCreator.ID),
				).
				LEFT_JOIN(tASource,
					tASource.ID.EQ(tDocRel.SourceUserID),
				),
		).
		WHERE(
			jet.AND(
				tDocs.DeletedAt.IS_NULL(),
				tDocRel.ID.IN(rIds...),
			),
		).
		ORDER_BY(
			tDocRel.CreatedAt.DESC(),
		).
		GROUP_BY(tDocs.ID).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Relations); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoFunc(userInfo)
	for i := 0; i < len(resp.Relations); i++ {
		if resp.Relations[i].SourceUser != nil {
			jobInfoFn(resp.Relations[i].SourceUser)
		}
		if resp.Relations[i].Document != nil && resp.Relations[i].Document.Creator != nil {
			jobInfoFn(resp.Relations[i].Document.Creator)
		}
	}

	resp.Pagination.Update(count.TotalCount, len(resp.Relations))

	return resp, nil
}

func (s *Server) addUserActivity(ctx context.Context, tx *sql.Tx, userId int32, targetUserId int32, activityType users.UserActivityType, key string, oldValue string, newValue string, reason string) error {
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
