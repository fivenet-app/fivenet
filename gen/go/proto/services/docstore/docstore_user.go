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
	userAct = table.FivenetUserActivity
)

func (s *Server) ListUserDocuments(ctx context.Context, req *ListUserDocumentsRequest) (*ListUserDocumentsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	resp := &ListUserDocumentsResponse{}
	// An user can never see their own activity on their own "profile"
	if userInfo.UserId == req.UserId {
		return resp, nil
	}

	condition := jet.AND(
		docRel.DeletedAt.IS_NULL(),
		jet.OR(
			docRel.SourceUserID.EQ(jet.Int32(req.UserId)),
			docRel.TargetUserID.EQ(jet.Int32(req.UserId)),
		),
	)

	var docIds []uint64
	idStmt := docRel.
		SELECT(
			docRel.DocumentID,
		).
		FROM(
			docRel,
		).
		WHERE(jet.AND(
			condition,
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

	dCreator := user.AS("creator")
	uSource := user.AS("source_user")
	uTarget := user.AS("target_user")
	stmt := docRel.
		SELECT(
			docRel.ID,
			docRel.CreatedAt,
			docRel.DeletedAt,
			docRel.DocumentID,
			docRel.SourceUserID,
			docs.ID,
			docs.CreatedAt,
			docs.UpdatedAt,
			docs.CategoryID,
			docs.CreatorID,
			docs.State,
			docs.Closed,
			docs.Title,
			dCategory.ID,
			dCategory.Name,
			dCategory.Description,
			dCreator.ID,
			dCreator.Identifier,
			dCreator.Job,
			dCreator.JobGrade,
			dCreator.Firstname,
			dCreator.Lastname,
			docRel.SourceUserID,
			uSource.ID,
			uSource.Identifier,
			uSource.Job,
			uSource.JobGrade,
			uSource.Firstname,
			uSource.Lastname,
			docRel.Relation,
			docRel.TargetUserID,
			uTarget.ID,
			uTarget.Identifier,
			uTarget.Job,
			uTarget.JobGrade,
			uTarget.Firstname,
			uTarget.Lastname,
		).
		FROM(
			docRel.
				LEFT_JOIN(docs,
					docRel.DocumentID.EQ(docs.ID),
				).
				LEFT_JOIN(dCategory,
					docs.CategoryID.EQ(dCategory.ID),
				).
				LEFT_JOIN(dCreator,
					docs.CreatorID.EQ(dCreator.ID),
				).
				LEFT_JOIN(uSource,
					uSource.ID.EQ(docRel.SourceUserID),
				).
				LEFT_JOIN(uTarget,
					uTarget.ID.EQ(docRel.TargetUserID),
				),
		).
		WHERE(
			jet.AND(
				docRel.DocumentID.IN(dIds...),
				condition,
				docs.DeletedAt.IS_NULL(),
			),
		).
		ORDER_BY(
			docRel.CreatedAt.DESC(),
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

	stmt := userAct.
		INSERT(
			userAct.SourceUserID,
			userAct.TargetUserID,
			userAct.Type,
			userAct.Key,
			userAct.OldValue,
			userAct.NewValue,
			userAct.Reason,
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
