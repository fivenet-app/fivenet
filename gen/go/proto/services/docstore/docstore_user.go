package docstore

import (
	context "context"
	"database/sql"
	"errors"

	"github.com/galexrt/fivenet/gen/go/proto/resources/documents"
	"github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/utils"
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

	idStmt := tDocRel.
		SELECT(
			tDocRel.ID.AS("id"),
			tDocRel.DocumentID.AS("documentId"),
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
			tDocRel.TargetUserID.EQ(jet.Int32(req.UserId)),
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

	var dbRelIds []struct {
		ID         uint64 `alias:"id"`
		DocumentID uint64 `alias:"documentId"`
	}
	if err := idStmt.QueryContext(ctx, s.db, &dbRelIds); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	relCount := len(dbRelIds)
	docIds := make([]uint64, relCount)
	for i := 0; i < relCount; i++ {
		docIds[i] = dbRelIds[i].DocumentID
	}
	docIds = utils.RemoveDuplicatesFromSlice(docIds)

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

	relIds := []uint64{}
	for i := 0; i < len(ids); i++ {
		for k := 0; k < len(dbRelIds); k++ {
			if dbRelIds[k].DocumentID == ids[i] {
				relIds = append(relIds, dbRelIds[k].ID)
			}
		}
	}

	relIds = utils.RemoveDuplicatesFromSlice(relIds)

	rIds := make([]jet.Expression, len(relIds))
	for i := 0; i < len(relIds); i++ {
		rIds[i] = jet.Uint64(relIds[i])
	}

	dCreator := tUsers.AS("creator")
	uSource := tUsers.AS("source_user")
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
				),
		).
		WHERE(
			jet.AND(
				tDocs.DeletedAt.IS_NULL(),
				tDocRel.ID.IN(rIds...),
				tDocRel.TargetUserID.EQ(jet.Int32(req.UserId)),
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
		if resp.Relations[i].SourceUser != nil {
			s.c.EnrichJobInfo(resp.Relations[i].SourceUser)
		}
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
