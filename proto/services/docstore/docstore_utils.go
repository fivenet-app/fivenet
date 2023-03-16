package docstore

import (
	context "context"

	database "github.com/galexrt/arpanet/proto/resources/common/database"
	"github.com/galexrt/arpanet/proto/resources/documents"
	jet "github.com/go-jet/jet/v2/mysql"
)

func (s *Server) getDocumentsQuery(where jet.BoolExpression, onlyColumns jet.ProjectionList, additionalColumns jet.ProjectionList, userID int32, job string, jobGrade int32) jet.SelectStatement {
	wheres := []jet.BoolExpression{jet.OR(
		jet.OR(
			docs.Public.IS_TRUE(),
			docs.CreatorID.EQ(jet.Int32(userID)),
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
	)}

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
					AND(dUserAccess.UserID.EQ(jet.Int32(userID)))).
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

func (s *Server) checkIfUserHasAccessToDoc(ctx context.Context, documentID uint64, userID int32, job string, jobGrade int32, access documents.DOC_ACCESS) (bool, error) {
	stmt := docs.SELECT(
		docs.ID,
	).
		FROM(
			docs.LEFT_JOIN(dUserAccess,
				dUserAccess.DocumentID.EQ(docs.ID).
					AND(dUserAccess.UserID.EQ(jet.Int32(userID)))).
				LEFT_JOIN(dJobAccess,
					dJobAccess.DocumentID.EQ(docs.ID).
						AND(dJobAccess.Job.EQ(jet.String(job))).
						AND(dJobAccess.MinimumGrade.LT_EQ(jet.Int32(jobGrade))),
				),
		).WHERE(
		jet.AND(
			docs.ID.EQ(jet.Uint64(documentID)),
			jet.OR(
				docs.CreatorID.EQ(jet.Int32(userID)),
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
	).
		LIMIT(1)

	var dest struct {
		ID uint64 `alias:"document.id"`
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return false, err
	}

	return dest.ID > 0, nil
}
