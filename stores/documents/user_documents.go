package documentsstore

import (
	"context"
	"errors"

	resourcesdatabase "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	documentsrelations "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/relations"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) ListUserDocuments(
	ctx context.Context,
	q ListUserDocumentsQuery,
) (resourcesdatabase.DataCount, []*documentsrelations.DocumentRelation, error) {
	if q.UserInfo == nil {
		q.UserInfo = &userinfo.UserInfo{}
	}

	tDocument := table.FivenetDocuments.AS("document")
	tDocRel := table.FivenetDocumentsRelations.AS("document_relation")
	tDCategory := table.FivenetDocumentsCategories.AS("category")
	tCreator := table.FivenetUser.AS("creator")
	tASource := tCreator.AS("source_user")

	visibleQuery := s.subjectAccess.VisibleIDsByConditionQuery(
		q.UserInfo,
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		table.FivenetDocuments.DeletedAt.IS_NULL(),
	)
	visibleIDs := mysql.
		SELECT(mysql.IntegerColumn("id").From(visibleQuery.Table)).
		FROM(visibleQuery.Table)

	condition := tDocRel.TargetUserID.EQ(mysql.Int32(q.UserID))
	if q.IncludeCreated {
		condition = mysql.OR(condition, tDocRel.SourceUserID.EQ(mysql.Int32(q.UserID)))
	}

	condition = mysql.AND(
		condition,
		tDocRel.DeletedAt.IS_NULL(),
		tDocument.ID.IN(visibleIDs),
	)

	if q.Closed != nil {
		condition = condition.AND(tDocument.Closed.EQ(mysql.Bool(*q.Closed)))
	}

	if len(q.Relations) == 0 {
		return resourcesdatabase.DataCount{}, []*documentsrelations.DocumentRelation{}, nil
	}

	types := make([]mysql.Expression, 0, len(q.Relations))
	for _, rel := range q.Relations {
		types = append(types, mysql.Int32(int32(rel)))
	}
	condition = condition.AND(tDocRel.Relation.IN(types...))

	var countStmt mysql.Statement = tDocRel.
		SELECT(
			mysql.COUNT(mysql.DISTINCT(tDocRel.DocumentID)).AS("data_count.total"),
		).
		FROM(
			tDocRel.
				INNER_JOIN(tDocument,
					tDocument.ID.EQ(tDocRel.DocumentID),
				),
		).
		WHERE(mysql.AND(
			condition,
			tDocument.ID.IN(visibleIDs),
		))
	if len(visibleQuery.CTEs) > 0 {
		countStmt = mysql.WITH(visibleQuery.CTEs...)(countStmt)
	}

	var count resourcesdatabase.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return resourcesdatabase.DataCount{}, nil, err
		}
	}

	_, limit := q.Pagination.GetResponseWithPageSize(count.Total, 20)
	if count.Total <= 0 {
		return count, []*documentsrelations.DocumentRelation{}, nil
	}

	orderBys := s.userDocumentSorter.Build(q.Sort)

	docRel := tDocRel.
		SELECT(
			tDocRel.ID.AS("id"),
			tDocRel.DocumentID.AS("document_id"),
		).
		FROM(
			tDocRel.
				INNER_JOIN(tDocument,
					tDocument.ID.EQ(tDocRel.DocumentID),
				),
		).
		WHERE(mysql.AND(
			condition,
			tDocument.ID.IN(visibleIDs),
		)).
		OFFSET(q.Pagination.GetOffset()).
		ORDER_BY(orderBys...).
		LIMIT(limit).
		AsTable("doc_rel")

	var stmt mysql.Statement = docRel.
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
			tDocument.State.AS("meta.state"),
			tDocument.Closed.AS("meta.closed"),
			tDocument.Draft.AS("meta.draft"),
			tDocument.Public.AS("meta.public"),
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
			docRel.
				INNER_JOIN(tDocRel,
					tDocRel.ID.EQ(mysql.RawInt("doc_rel.id")),
				).
				INNER_JOIN(tDocument,
					tDocument.ID.EQ(mysql.RawInt("doc_rel.document_id")),
				).
				LEFT_JOIN(tDCategory,
					mysql.AND(
						tDocument.CategoryID.EQ(tDCategory.ID),
						tDCategory.DeletedAt.IS_NULL(),
					),
				).
				LEFT_JOIN(tCreator,
					tDocument.CreatorID.EQ(tCreator.ID),
				).
				LEFT_JOIN(tASource,
					tASource.ID.EQ(tDocRel.SourceUserID),
				),
		).
		WHERE(mysql.AND(
			tDocument.DeletedAt.IS_NULL(),
			tDocument.ID.IN(visibleIDs),
		)).
		ORDER_BY(orderBys...).
		LIMIT(limit)

	var relations []*documentsrelations.DocumentRelation
	if len(visibleQuery.CTEs) > 0 {
		stmt = mysql.WITH(visibleQuery.CTEs...)(stmt)
	}
	if err := stmt.QueryContext(ctx, s.db, &relations); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return resourcesdatabase.DataCount{}, nil, err
		}
	}

	return count, relations, nil
}
