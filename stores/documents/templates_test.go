package documentsstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	documentstemplates "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/templates"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreListTemplates(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db)

	mock.ExpectQuery(`(?s).*AS doc_ids INNER JOIN fivenet_documents_templates AS template_short.*ORDER BY template_short\.weight DESC, template_short\.id ASC.*`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	templates, err := store.ListTemplates(t.Context(), &userinfo.UserInfo{Superuser: true})
	require.NoError(t, err)
	assert.Empty(t, templates)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListTemplatesUsesVisibilityTablesForNonSuperuser(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db)

	mock.ExpectQuery(`(?s).*WITH actor_subjects AS .*fivenet_documents_templates_visibility_creator.*fivenet_documents_templates_visibility_subject.*AS doc_ids INNER JOIN fivenet_documents_templates AS template_short.*ORDER BY template_short\.weight DESC, template_short\.id ASC.*`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	templates, err := store.ListTemplates(
		t.Context(),
		&userinfo.UserInfo{UserId: 3, Job: "doj", JobGrade: 16},
	)
	require.NoError(t, err)
	assert.Empty(t, templates)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreGetTemplate(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db)

	mock.ExpectQuery(`(?s).*FROM fivenet_documents_templates AS template.*LIMIT \?.*`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	tmpl, err := store.GetTemplate(t.Context(), 42)
	require.NoError(t, err)
	assert.Nil(t, tmpl)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreTemplateWrites(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db)
	tmpl := &documentstemplates.Template{Id: 7, Title: "Hello"}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_documents_templates`)).
		WillReturnResult(sqlmock.NewResult(7, 1))
	_, err = store.CreateTemplate(t.Context(), db, tmpl, "doj", nil)
	require.NoError(t, err)

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE fivenet_documents_templates SET`)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, store.UpdateTemplate(t.Context(), db, tmpl, nil))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE fivenet_documents_templates SET`)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, store.DeleteTemplate(t.Context(), db, 7, "doj"))

	require.NoError(t, mock.ExpectationsWereMet())
}
