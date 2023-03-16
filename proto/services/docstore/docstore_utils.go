package docstore

import (
	context "context"
	"errors"
	"fmt"

	database "github.com/galexrt/arpanet/proto/resources/common/database"
	"github.com/galexrt/arpanet/proto/resources/documents"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) getDocumentsQuery(where jet.BoolExpression, onlyColumns jet.ProjectionList, additionalColumns jet.ProjectionList, userId int32, job string, jobGrade int32) jet.SelectStatement {
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
						dUserAccess.Access.NOT_EQ(jet.Int32(int32(documents.DOC_ACCESS_BLOCKED))),
					),
					jet.AND(
						dUserAccess.Access.IS_NULL(),
						dJobAccess.Access.IS_NOT_NULL(),
						dJobAccess.Access.NOT_EQ(jet.Int32(int32(documents.DOC_ACCESS_BLOCKED))),
					),
				),
			),
		),
	}

	if where != nil {
		wheres = append(wheres, where)
	}

	u := u.AS("creator")
	var q jet.SelectStatement
	if onlyColumns != nil {
		q = docs.SELECT(
			onlyColumns,
		)
	} else {
		if additionalColumns == nil {
			q = docs.SELECT(
				docs.AllColumns,
				dCategory.ID,
				dCategory.Name,
				u.ID,
				u.Identifier,
				u.Job,
				u.JobGrade,
				u.Firstname,
				u.Lastname,
			)
		} else {
			additionalColumns = append(jet.ProjectionList{
				dCategory.Name,
				dCategory.ID,
				u.ID,
				u.Identifier,
				u.Job,
				u.JobGrade,
				u.Firstname,
				u.Lastname,
			}, additionalColumns)
			q = docs.SELECT(
				docs.AllColumns,
				additionalColumns...,
			)
		}
	}

	return q.
		FROM(
			docs.LEFT_JOIN(dUserAccess,
				dUserAccess.DocumentID.EQ(docs.ID).
					AND(dUserAccess.UserID.EQ(jet.Int32(userId)))).
				LEFT_JOIN(dJobAccess,
					dJobAccess.DocumentID.EQ(docs.ID).
						AND(dJobAccess.Job.EQ(jet.String(job))).
						AND(dJobAccess.MinimumGrade.LT_EQ(jet.Int32(jobGrade))),
				).
				LEFT_JOIN(u,
					docs.CreatorID.EQ(u.ID),
				).
				LEFT_JOIN(dCategory,
					docs.CategoryID.EQ(dCategory.ID),
				),
		).WHERE(
		jet.AND(
			wheres...,
		),
	).
		ORDER_BY(docs.CreatedAt.DESC()).
		LIMIT(database.PaginationLimit)
}

func (s *Server) checkIfUserHasAccessToDoc(ctx context.Context, documentID uint64, userId int32, job string, jobGrade int32, access documents.DOC_ACCESS) (bool, error) {
	out, err := s.checkIfUserHasAccessToDocIDs(ctx, userId, job, jobGrade, access, documentID)
	return len(out) > 0, err
}

func (s *Server) checkIfUserHasAccessToDocs(ctx context.Context, userId int32, job string, jobGrade int32, access documents.DOC_ACCESS, documentIDs ...uint64) (bool, error) {
	out, err := s.checkIfUserHasAccessToDocIDs(ctx, userId, job, jobGrade, access, documentIDs...)
	return len(out) == len(documentIDs), err
}

func (s *Server) checkIfUserHasAccessToDocIDs(ctx context.Context, userId int32, job string, jobGrade int32, access documents.DOC_ACCESS, documentIDs ...uint64) ([]uint64, error) {
	ids := make([]jet.Expression, len(documentIDs))
	for i := 0; i < len(documentIDs); i++ {
		ids[i] = jet.Uint64(documentIDs[i])
	}

	stmt := docs.SELECT(
		docs.ID,
	).
		FROM(
			docs.LEFT_JOIN(dUserAccess,
				dUserAccess.DocumentID.EQ(docs.ID).
					AND(dUserAccess.UserID.EQ(jet.Int32(userId)))).
				LEFT_JOIN(dJobAccess,
					dJobAccess.DocumentID.EQ(docs.ID).
						AND(dJobAccess.Job.EQ(jet.String(job))).
						AND(dJobAccess.MinimumGrade.LT_EQ(jet.Int32(jobGrade))),
				),
		).WHERE(
		jet.AND(
			docs.ID.IN(ids...),
			docs.DeletedAt.IS_NULL(),
			jet.OR(
				docs.CreatorID.EQ(jet.Int32(userId)),
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
		),
	)

	fmt.Println(stmt.DebugSql())

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
