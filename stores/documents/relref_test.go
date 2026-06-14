package documentsstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	documentsreferences "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/references"
	documentsrelations "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/relations"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreDocumentReferences(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_documents_references AS document_reference`)+`(?s).*`+regexp.QuoteMeta(`document_reference.id = ?`)+`(?s).*`+regexp.QuoteMeta(`LIMIT ?`)).
		WithArgs(int64(7), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"document_reference.id", "document_reference.source_document_id", "document_reference.target_document_id", "document_reference.reference", "document_reference.creator_id"}).AddRow(int64(7), int64(42), int64(99), 1, int32(3)))

	ref, err := store.GetDocumentReference(t.Context(), 7)
	require.NoError(t, err)
	require.NotNil(t, ref)
	assert.Equal(t, int64(7), ref.GetId())
	assert.Equal(t, int64(42), ref.GetSourceDocumentId())
	assert.Equal(t, int64(99), ref.GetTargetDocumentId())
	assert.Equal(t, int32(3), ref.GetCreatorId())

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_documents_references AS document_reference`)+`(?s).*`+regexp.QuoteMeta(`document_reference.deleted_at IS NULL`)+`(?s).*`+regexp.QuoteMeta(`LIMIT ?`)).
		WithArgs(int64(42), int64(42), int64(25)).
		WillReturnRows(sqlmock.NewRows([]string{}))

	refs, err := store.ListDocumentReferences(t.Context(), 42)
	require.NoError(t, err)
	assert.Empty(t, refs)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_documents_references`)).
		WithArgs(int64(42), int32(1), int64(99), int32(3)).
		WillReturnResult(sqlmock.NewResult(7, 1))
	lastID, err := store.CreateDocumentReference(
		t.Context(),
		db,
		&documentsreferences.DocumentReference{
			SourceDocumentId: 42,
			TargetDocumentId: 99,
			Reference:        documentsreferences.DocReference_DOC_REFERENCE_LINKED,
			CreatorId:        func() *int32 { v := int32(3); return &v }(),
		},
	)
	require.NoError(t, err)
	assert.Equal(t, int64(7), lastID)

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE fivenet_documents_references`)).
		WithArgs(int64(7), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, store.DeleteDocumentReference(t.Context(), db, 7))

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreDocumentRelations(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_documents_relations AS document_relation`)+`(?s).*`+regexp.QuoteMeta(`document_relation.id = ?`)+`(?s).*`+regexp.QuoteMeta(`LIMIT ?`)).
		WithArgs(int64(11), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"document_relation.id", "document_relation.document_id", "document_relation.source_user_id", "document_relation.relation", "document_relation.target_user_id"}).AddRow(int64(11), int64(42), int32(3), int32(1), int32(8)))

	rel, err := store.GetDocumentRelation(t.Context(), 11)
	require.NoError(t, err)
	require.NotNil(t, rel)
	assert.Equal(t, int64(11), rel.GetId())
	assert.Equal(t, int64(42), rel.GetDocumentId())
	assert.Equal(t, int32(3), rel.GetSourceUserId())
	assert.Equal(t, int32(8), rel.GetTargetUserId())

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_documents_relations AS document_relation`)+`(?s).*`+regexp.QuoteMeta(`document_relation.deleted_at IS NULL`)+`(?s).*`+regexp.QuoteMeta(`LIMIT ?`)).
		WithArgs(int64(42), int64(25)).
		WillReturnRows(sqlmock.NewRows([]string{}))

	rels, err := store.ListDocumentRelations(t.Context(), 42)
	require.NoError(t, err)
	assert.Empty(t, rels)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_documents_relations`)).
		WithArgs(int64(42), int32(3), int32(1), int32(8)).
		WillReturnResult(sqlmock.NewResult(12, 1))
	lastID, created, err := store.CreateDocumentRelation(
		t.Context(),
		db,
		&documentsrelations.DocumentRelation{
			DocumentId:   42,
			SourceUserId: 3,
			Relation:     documentsrelations.DocRelation_DOC_RELATION_MENTIONED,
			TargetUserId: 8,
		},
	)
	require.NoError(t, err)
	assert.True(t, created)
	assert.Equal(t, int64(12), lastID)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_documents_relations`)).
		WithArgs(int64(42), int32(3), int32(1), int32(8)).
		WillReturnError(&mysql.MySQLError{Number: 1062, Message: "duplicate"})
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_documents_relations`)+`(?s).*`+regexp.QuoteMeta(`document_id = ?`)+`(?s).*`+regexp.QuoteMeta(`relation = ?`)+`(?s).*`+regexp.QuoteMeta(`target_user_id = ?`)+`(?s).*`+regexp.QuoteMeta(`LIMIT ?`)).
		WithArgs(int64(42), int32(1), int32(8), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(44)))
	lastID, created, err = store.CreateDocumentRelation(
		t.Context(),
		db,
		&documentsrelations.DocumentRelation{
			DocumentId:   42,
			SourceUserId: 3,
			Relation:     documentsrelations.DocRelation_DOC_RELATION_MENTIONED,
			TargetUserId: 8,
		},
	)
	require.NoError(t, err)
	assert.False(t, created)
	assert.Equal(t, int64(44), lastID)

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE fivenet_documents_relations`)).
		WithArgs(int64(11), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, store.DeleteDocumentRelation(t.Context(), db, 11))

	require.NoError(t, mock.ExpectationsWereMet())
}
