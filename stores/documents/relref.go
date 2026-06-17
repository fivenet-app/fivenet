package documentsstore

import (
	"context"
	"errors"

	documentsreferences "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/references"
	documentsrelations "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/relations"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) GetDocumentReference(
	ctx context.Context,
	id int64,
	includeDeleted bool,
) (*documentsreferences.DocumentReference, error) {
	tRef := table.FivenetDocumentsReferences.AS("document_reference")
	stmt := tRef.
		SELECT(
			tRef.ID,
			tRef.CreatedAt,
			tRef.SourceDocumentID,
			tRef.Reference,
			tRef.TargetDocumentID,
			tRef.CreatorID,
		).
		FROM(tRef).
		WHERE(mysql.AND(
			tRef.ID.EQ(mysql.Int64(id)),
			mysql.OR(
				mysql.Bool(includeDeleted),
				tRef.DeletedAt.IS_NULL(),
			),
		)).
		LIMIT(1)

	var dest documentsreferences.DocumentReference
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return &dest, nil
}

func (s *Store) ListDocumentReferences(
	ctx context.Context,
	documentID int64,
	includeDeleted bool,
) ([]*documentsreferences.DocumentReference, error) {
	tRef := table.FivenetDocumentsReferences.AS("document_reference")
	tSourceDoc := table.FivenetDocuments.AS("source_document")
	tTargetDoc := table.FivenetDocuments.AS("target_document")
	tCreator := table.FivenetUser.AS("creator")
	tRefCreator := tCreator.AS("ref_creator")

	stmt := tRef.
		SELECT(
			tRef.ID,
			tRef.CreatedAt,
			tRef.SourceDocumentID,
			tRef.Reference,
			tRef.TargetDocumentID,
			tRef.CreatorID,
			tSourceDoc.ID,
			tSourceDoc.CreatedAt,
			tSourceDoc.UpdatedAt,
			tSourceDoc.CategoryID,
			tSourceDoc.Title,
			tSourceDoc.CreatorID,
			tSourceDoc.State,
			tSourceDoc.Closed,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tTargetDoc.ID,
			tTargetDoc.CreatedAt,
			tTargetDoc.UpdatedAt,
			tTargetDoc.CategoryID,
			tTargetDoc.Title,
			tTargetDoc.CreatorID,
			tTargetDoc.State,
			tTargetDoc.Closed,
			tRefCreator.ID,
			tRefCreator.Job,
			tRefCreator.JobGrade,
			tRefCreator.Firstname,
			tRefCreator.Lastname,
			tRefCreator.Dateofbirth,
		).
		FROM(
			tRef.
				LEFT_JOIN(tSourceDoc,
					tRef.SourceDocumentID.EQ(tSourceDoc.ID),
				).
				LEFT_JOIN(tTargetDoc,
					tRef.TargetDocumentID.EQ(tTargetDoc.ID),
				).
				LEFT_JOIN(tCreator,
					tSourceDoc.CreatorID.EQ(tCreator.ID),
				).
				LEFT_JOIN(tRefCreator,
					tRef.CreatorID.EQ(tRefCreator.ID),
				),
		).
		WHERE(mysql.AND(
			mysql.OR(
				tRef.SourceDocumentID.EQ(mysql.Int64(documentID)),
				tRef.TargetDocumentID.EQ(mysql.Int64(documentID)),
			),
			mysql.OR(
				mysql.Bool(includeDeleted),
				tRef.DeletedAt.IS_NULL(),
			),
		)).
		ORDER_BY(tRef.CreatedAt.DESC()).
		LIMIT(25)

	var dest []*documentsreferences.DocumentReference
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}

func (s *Store) CreateDocumentReference(
	ctx context.Context,
	tx qrm.DB,
	ref *documentsreferences.DocumentReference,
) (int64, error) {
	tRef := table.FivenetDocumentsReferences
	stmt := tRef.
		INSERT(
			tRef.SourceDocumentID,
			tRef.Reference,
			tRef.TargetDocumentID,
			tRef.CreatorID,
		).
		VALUES(
			ref.GetSourceDocumentId(),
			ref.GetReference(),
			ref.GetTargetDocumentId(),
			ref.GetCreatorId(),
		)

	result, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (s *Store) DeleteDocumentReference(ctx context.Context, tx qrm.DB, id int64) error {
	tRef := table.FivenetDocumentsReferences
	stmt := tRef.
		UPDATE(tRef.DeletedAt).
		SET(tRef.DeletedAt.SET(mysql.CURRENT_TIMESTAMP())).
		WHERE(tRef.ID.EQ(mysql.Int64(id))).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) GetDocumentRelation(
	ctx context.Context,
	id int64,
	includeDeleted bool,
) (*documentsrelations.DocumentRelation, error) {
	tRel := table.FivenetDocumentsRelations.AS("document_relation")
	stmt := tRel.
		SELECT(
			tRel.ID,
			tRel.CreatedAt,
			tRel.DocumentID,
			tRel.SourceUserID,
			tRel.Relation,
			tRel.TargetUserID,
		).
		FROM(tRel).
		WHERE(mysql.AND(
			tRel.ID.EQ(mysql.Int64(id)),
			mysql.OR(
				mysql.Bool(includeDeleted),
				tRel.DeletedAt.IS_NULL(),
			),
		)).
		LIMIT(1)

	var dest documentsrelations.DocumentRelation
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return &dest, nil
}

func (s *Store) ListDocumentRelations(
	ctx context.Context,
	documentID int64,
	includeDeleted bool,
) ([]*documentsrelations.DocumentRelation, error) {
	tRel := table.FivenetDocumentsRelations.AS("document_relation")
	tDocument := table.FivenetDocuments.AS("document")
	tCategory := table.FivenetDocumentsCategories.AS("category")
	tSourceUser := table.FivenetUser.AS("source_user")
	tTargetUser := tSourceUser.AS("target_user")

	stmt := tRel.
		SELECT(
			tRel.ID,
			tRel.CreatedAt,
			tRel.DocumentID,
			tRel.SourceUserID,
			tRel.Relation,
			tRel.TargetUserID,
			tDocument.ID,
			tDocument.CreatedAt,
			tDocument.UpdatedAt,
			tDocument.CategoryID,
			tDocument.Title,
			tDocument.CreatorID,
			tDocument.State,
			tDocument.Closed,
			tDocument.Draft,
			tDocument.Public,
			tCategory.ID,
			tCategory.Name,
			tCategory.Description,
			tCategory.Color,
			tCategory.Icon,
			tSourceUser.ID,
			tSourceUser.Job,
			tSourceUser.JobGrade,
			tSourceUser.Firstname,
			tSourceUser.Lastname,
			tSourceUser.Dateofbirth,
			tTargetUser.ID,
			tTargetUser.Job,
			tTargetUser.JobGrade,
			tTargetUser.Firstname,
			tTargetUser.Lastname,
			tTargetUser.Dateofbirth,
		).
		FROM(
			tRel.
				LEFT_JOIN(tDocument,
					tDocument.ID.EQ(tRel.DocumentID),
				).
				LEFT_JOIN(tCategory,
					mysql.AND(
						tDocument.CategoryID.EQ(tCategory.ID),
						mysql.OR(
							mysql.Bool(includeDeleted),
							tCategory.DeletedAt.IS_NULL(),
						),
					),
				).
				LEFT_JOIN(tSourceUser,
					tSourceUser.ID.EQ(tRel.SourceUserID),
				).
				LEFT_JOIN(tTargetUser,
					tTargetUser.ID.EQ(tRel.TargetUserID),
				),
		).
		WHERE(mysql.AND(
			tRel.DocumentID.EQ(mysql.Int64(documentID)),
			mysql.OR(
				mysql.Bool(includeDeleted),
				tRel.DeletedAt.IS_NULL(),
			),
		)).
		ORDER_BY(tRel.CreatedAt.DESC()).
		LIMIT(25)

	var dest []*documentsrelations.DocumentRelation
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}

func (s *Store) CreateDocumentRelation(
	ctx context.Context,
	tx qrm.DB,
	rel *documentsrelations.DocumentRelation,
) (int64, bool, error) {
	tRel := table.FivenetDocumentsRelations
	stmt := tRel.
		INSERT(
			tRel.DocumentID,
			tRel.SourceUserID,
			tRel.Relation,
			tRel.TargetUserID,
		).
		VALUES(
			rel.GetDocumentId(),
			rel.GetSourceUserId(),
			rel.GetRelation(),
			rel.GetTargetUserId(),
		)

	result, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		if !dbutils.IsDuplicateError(err) {
			return 0, false, err
		}

		stmt := tRel.
			SELECT(tRel.ID.AS("id")).
			FROM(tRel).
			WHERE(mysql.AND(
				tRel.DocumentID.EQ(mysql.Int64(rel.GetDocumentId())),
				tRel.Relation.EQ(mysql.Int32(int32(rel.GetRelation()))),
				tRel.TargetUserID.EQ(mysql.Int32(rel.GetTargetUserId())),
			)).
			LIMIT(1)

		var dest struct{ ID int64 }
		if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
			return 0, false, err
		}

		return dest.ID, false, nil
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, false, err
	}

	return lastID, true, nil
}

func (s *Store) DeleteDocumentRelation(ctx context.Context, tx qrm.DB, id int64) error {
	tRel := table.FivenetDocumentsRelations
	stmt := tRel.
		UPDATE(tRel.DeletedAt).
		SET(tRel.DeletedAt.SET(mysql.CURRENT_TIMESTAMP())).
		WHERE(tRel.ID.EQ(mysql.Int64(id))).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}
