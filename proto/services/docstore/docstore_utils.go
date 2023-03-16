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
			ad.Public.IS_TRUE(),
			ad.CreatorID.EQ(jet.Int32(userID)),
		),
		jet.OR(
			jet.AND(
				adua.Access.IS_NOT_NULL(),
				adua.Access.NOT_EQ(jet.Int32(int32(documents.DOC_ACCESS_BLOCKED))),
			),
			jet.AND(
				adua.Access.IS_NULL(),
				adja.Access.IS_NOT_NULL(),
				adja.Access.NOT_EQ(jet.Int32(int32(documents.DOC_ACCESS_BLOCKED))),
			),
		),
	)}

	if where != nil {
		wheres = append(wheres, where)
	}

	u := u.AS("creator")
	var q jet.SelectStatement
	if onlyColumns != nil {
		q = ad.SELECT(
			onlyColumns,
		)
	} else {
		if additionalColumns == nil {
			q = ad.SELECT(
				ad.AllColumns,
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
			q = ad.SELECT(
				ad.AllColumns,
				additionalColumns...,
			)
		}
	}

	return q.
		FROM(
			ad.LEFT_JOIN(adua,
				adua.DocumentID.EQ(ad.ID).
					AND(adua.UserID.EQ(jet.Int32(userID)))).
				LEFT_JOIN(adja,
					adja.DocumentID.EQ(ad.ID).
						AND(adja.Job.EQ(jet.String(job))).
						AND(adja.MinimumGrade.LT_EQ(jet.Int32(jobGrade))),
				).
				LEFT_JOIN(u,
					ad.CreatorID.EQ(u.ID),
				).
				LEFT_JOIN(dCategory,
					ad.CategoryID.EQ(dCategory.ID),
				),
		).WHERE(
		jet.AND(
			wheres...,
		),
	).
		ORDER_BY(ad.CreatedAt.DESC()).
		LIMIT(database.PaginationLimit)
}

func (s *Server) checkIfUserHasAccessToDoc(ctx context.Context, userID int32, job string, jobGrade int32, access documents.DOC_ACCESS) (bool, error) {
	checkStmt := ad.SELECT(
		ad.ID,
	).
		FROM(
			ad.LEFT_JOIN(adua,
				adua.DocumentID.EQ(ad.ID).
					AND(adua.UserID.EQ(jet.Int32(userID)))).
				LEFT_JOIN(adja,
					adja.DocumentID.EQ(ad.ID).
						AND(adja.Job.EQ(jet.String(job))).
						AND(adja.MinimumGrade.LT_EQ(jet.Int32(jobGrade))),
				),
		).WHERE(
		jet.OR(
			ad.CreatorID.EQ(jet.Int32(userID)),
			jet.AND(
				adua.Access.IS_NOT_NULL(),
				adua.Access.GT_EQ(jet.Int32(int32(access))),
			),
			jet.AND(
				adua.Access.IS_NULL(),
				adja.Access.IS_NOT_NULL(),
				adja.Access.GT_EQ(jet.Int32(int32(access))),
			),
		),
	).
		LIMIT(1)

	var dest struct {
		ID uint64 `alias:"document.id"`
	}
	if err := checkStmt.QueryContext(ctx, s.db, &dest); err != nil {
		return false, err
	}

	return dest.ID > 0, nil
}
