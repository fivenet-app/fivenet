package docstore

import (
	context "context"
	"database/sql"
	"errors"

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

	resp := &ListUserDocumentsResponse{}

	var docIds []uint64
	idStmt := tDocRel.
		SELECT(
			tDocRel.DocumentID,
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
		WHERE(jet.AND(
			tDocRel.DeletedAt.IS_NULL(),
			jet.OR(
				tDocRel.SourceUserID.EQ(jet.Int32(req.UserId)),
				tDocRel.TargetUserID.EQ(jet.Int32(req.UserId)),
			),
			jet.OR(
				jet.OR(
					tDocs.Public.IS_TRUE(),
					tDocs.CreatorID.EQ(jet.Int32(userInfo.UserId)),
				),
				jet.OR(
					jet.AND(
						tDUserAccess.Access.IS_NOT_NULL(),
						tDUserAccess.Access.NOT_EQ(jet.Int32(int32(documents.ACCESS_LEVEL_BLOCKED))),
					),
					jet.AND(
						tDUserAccess.Access.IS_NULL(),
						tDJobAccess.Access.IS_NOT_NULL(),
						tDJobAccess.Access.NOT_EQ(jet.Int32(int32(documents.ACCESS_LEVEL_BLOCKED))),
					),
				),
			),
		))

	if err := idStmt.QueryContext(ctx, s.db, &docIds); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	if len(docIds) == 0 {
		return resp, nil
	}

	ids, err := s.checkIfUserHasAccessToDocIDs(ctx, userInfo, true, documents.ACCESS_LEVEL_VIEW, docIds...)
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return resp, nil
	}

	dIds := make([]jet.Expression, len(ids))
	for i := 0; i < len(ids); i++ {
		dIds[i] = jet.Uint64(ids[i])
	}

	dCreator := tUsers.AS("creator")
	uSource := tUsers.AS("source_user")
	uTarget := tUsers.AS("target_user")
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
			dCreator.ID,
			dCreator.Identifier,
			dCreator.Job,
			dCreator.JobGrade,
			dCreator.Firstname,
			dCreator.Lastname,
			tDocRel.SourceUserID,
			uSource.ID,
			uSource.Identifier,
			uSource.Job,
			uSource.JobGrade,
			uSource.Firstname,
			uSource.Lastname,
			tDocRel.Relation,
			tDocRel.TargetUserID,
			uTarget.ID,
			uTarget.Identifier,
			uTarget.Job,
			uTarget.JobGrade,
			uTarget.Firstname,
			uTarget.Lastname,
		).
		FROM(
			tDocRel.
				LEFT_JOIN(tDocs,
					tDocRel.DocumentID.EQ(tDocs.ID),
				).
				LEFT_JOIN(tDCategory,
					tDocs.CategoryID.EQ(tDCategory.ID),
				).
				LEFT_JOIN(dCreator,
					tDocs.CreatorID.EQ(dCreator.ID),
				).
				LEFT_JOIN(uSource,
					uSource.ID.EQ(tDocRel.SourceUserID),
				).
				LEFT_JOIN(uTarget,
					uTarget.ID.EQ(tDocRel.TargetUserID),
				),
		).
		WHERE(
			jet.AND(
				tDocs.DeletedAt.IS_NULL(),
				tDocRel.DocumentID.IN(dIds...),
			),
		).
		ORDER_BY(
			tDocRel.CreatedAt.DESC(),
		)

	if err := stmt.QueryContext(ctx, s.db, &resp.Relations); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	for i := 0; i < len(resp.Relations); i++ {
		s.c.EnrichJobInfo(resp.Relations[i].SourceUser)
		s.c.EnrichJobInfo(resp.Relations[i].TargetUser)
	}

	return resp, nil
}

func (s *Server) addUserActivity(ctx context.Context, tx *sql.Tx, userId int32, targetUserId int32, activityType users.USER_ACTIVITY_TYPE, key string, oldValue string, newValue string, reason string) error {
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
