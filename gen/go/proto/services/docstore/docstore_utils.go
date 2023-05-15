package docstore

import (
	context "context"
	"errors"

	"github.com/galexrt/fivenet/gen/go/proto/resources/common"
	"github.com/galexrt/fivenet/gen/go/proto/resources/documents"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) getDocumentsQuery(where jet.BoolExpression, onlyColumns jet.ProjectionList, contentLength int, userId int32, job string, jobGrade int32) jet.SelectStatement {
	wheres := []jet.BoolExpression{
		jet.AND(
			docs.DeletedAt.IS_NULL(),
			jet.OR(
				jet.OR(
					docs.Public.IS_TRUE(),
					docs.CreatorID.EQ(jet.Int32(userId)),
				),
				jet.OR(
					jet.AND(
						dUserAccess.Access.IS_NOT_NULL(),
						dUserAccess.Access.NOT_EQ(jet.Int32(int32(documents.ACCESS_LEVEL_BLOCKED))),
					),
					jet.AND(
						dUserAccess.Access.IS_NULL(),
						dJobAccess.Access.IS_NOT_NULL(),
						dJobAccess.Access.NOT_EQ(jet.Int32(int32(documents.ACCESS_LEVEL_BLOCKED))),
					),
				),
			),
		),
	}

	if where != nil {
		wheres = append(wheres, where)
	}

	var q jet.SelectStatement
	if onlyColumns != nil {
		q = docs.
			SELECT(
				onlyColumns,
			)
	} else {
		columns := jet.ProjectionList{
			docs.ID,
			docs.CreatedAt,
			docs.UpdatedAt,
			docs.DeletedAt,
			docs.CategoryID,
			dCategory.ID,
			dCategory.Name,
			dCategory.Description,
			dCategory.Job,
			docs.Title,
			docs.ContentType,
			docs.Data,
			docs.CreatorID,
			uCreator.ID,
			uCreator.Identifier,
			uCreator.Job,
			uCreator.JobGrade,
			uCreator.Firstname,
			uCreator.Lastname,
			docs.CreatorJob,
			docs.State,
			docs.Closed,
			docs.Public,
		}
		if contentLength > 0 {
			columns = append(columns, jet.LEFT(docs.Content, jet.Int(int64(contentLength))).AS("document.content"))
		} else {
			columns = append(columns, docs.Content)
		}
		q = docs.SELECT(columns[0], columns[1:])
	}

	return q.
		FROM(
			docs.
				LEFT_JOIN(dUserAccess,
					dUserAccess.DocumentID.EQ(docs.ID).
						AND(dUserAccess.UserID.EQ(jet.Int32(userId)))).
				LEFT_JOIN(dJobAccess,
					dJobAccess.DocumentID.EQ(docs.ID).
						AND(dJobAccess.Job.EQ(jet.String(job))).
						AND(dJobAccess.MinimumGrade.LT_EQ(jet.Int32(jobGrade))),
				).
				LEFT_JOIN(dCategory,
					docs.CategoryID.EQ(dCategory.ID),
				).
				LEFT_JOIN(uCreator,
					docs.CreatorID.EQ(uCreator.ID),
				),
		).
		WHERE(
			jet.AND(
				wheres...,
			),
		).
		ORDER_BY(
			docs.CreatedAt.DESC(),
			docs.UpdatedAt.DESC(),
		)
}

func (s *Server) checkIfUserHasAccessToDoc(ctx context.Context, documentId uint64, userInfo *userinfo.UserInfo, publicOk bool, access documents.ACCESS_LEVEL) (bool, error) {
	out, err := s.checkIfUserHasAccessToDocIDs(ctx, userInfo, publicOk, access, documentId)
	return len(out) > 0, err
}

func (s *Server) checkIfUserHasAccessToDocs(ctx context.Context, userInfo *userinfo.UserInfo, publicOk bool, access documents.ACCESS_LEVEL, documentIds ...uint64) (bool, error) {
	out, err := s.checkIfUserHasAccessToDocIDs(ctx, userInfo, publicOk, access, documentIds...)
	return len(out) == len(documentIds), err
}

func (s *Server) checkIfUserHasAccessToDocIDs(ctx context.Context, userInfo *userinfo.UserInfo, publicOk bool, access documents.ACCESS_LEVEL, documentIds ...uint64) ([]uint64, error) {
	if len(documentIds) == 0 {
		return documentIds, nil
	}

	// Allow superusers access to any docs
	if s.p.Can(userInfo, common.SuperuserCategoryPerm, common.SuperuserAnyAccessName) {
		return documentIds, nil
	}

	ids := make([]jet.Expression, len(documentIds))
	for i := 0; i < len(documentIds); i++ {
		ids[i] = jet.Uint64(documentIds[i])
	}

	condition := jet.AND(
		docs.ID.IN(ids...),
		docs.DeletedAt.IS_NULL(),
		jet.OR(
			docs.Public.IS_TRUE(),
			docs.CreatorID.EQ(jet.Int32(userInfo.UserId)),
			jet.AND(
				dUserAccess.Access.IS_NOT_NULL(),
				dUserAccess.Access.GT_EQ(jet.Int32(int32(access))),
			),
			jet.AND(
				dUserAccess.Access.IS_NULL(),
				dJobAccess.Access.IS_NOT_NULL(),
				dJobAccess.Access.GT_EQ(jet.Int32(int32(access))),
			),
		),
	)

	stmt := docs.
		SELECT(
			docs.ID,
		).
		FROM(
			docs.
				LEFT_JOIN(dUserAccess,
					dUserAccess.DocumentID.EQ(docs.ID).
						AND(dUserAccess.UserID.EQ(jet.Int32(userInfo.UserId)))).
				LEFT_JOIN(dJobAccess,
					dJobAccess.DocumentID.EQ(docs.ID).
						AND(dJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(dJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				),
		).
		WHERE(condition).
		GROUP_BY(docs.ID).
		ORDER_BY(docs.ID.DESC(), dJobAccess.MinimumGrade)

	var dest struct {
		IDs []uint64 `alias:"document.id"`
	}
	if err := stmt.QueryContext(ctx, s.db, &dest.IDs); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	return dest.IDs, nil
}
