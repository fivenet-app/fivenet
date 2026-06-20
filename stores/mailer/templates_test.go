package mailerstore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	mailertemplates "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/templates"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreListTemplates(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)
	now := time.Unix(0, 0).UTC()

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_mailer_templates AS template`) +
		`(?s).*` + regexp.QuoteMeta(`template.email_id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)
	mock.ExpectQuery(expectedQuery).
		WithArgs(int64(7), int64(25)).
		WillReturnRows(sqlmock.NewRows([]string{
			"template.id",
			"template.created_at",
			"template.updated_at",
			"template.deleted_at",
			"template.email_id",
			"template.title",
			"template.content",
		}).AddRow(int64(11), now, now, nil, int64(7), "Welcome", nil))

	templates, err := store.ListTemplates(t.Context(), db, 7, 25)
	require.NoError(t, err)
	require.Len(t, templates, 1)
	assert.Equal(t, int64(11), templates[0].GetId())
	assert.Equal(t, "Welcome", templates[0].GetTitle())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreGetTemplateGlobal(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)
	now := time.Unix(0, 0).UTC()

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_mailer_templates AS template`) +
		`(?s).*` + regexp.QuoteMeta(`template.id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`template.email_id IS NULL`) +
		`(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)
	mock.ExpectQuery(expectedQuery).
		WithArgs(int64(22), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"template.id",
			"template.created_at",
			"template.updated_at",
			"template.deleted_at",
			"template.email_id",
			"template.title",
			"template.content",
			"template.creator_job",
			"template.creator_id",
		}).AddRow(int64(22), now, now, nil, nil, "Global", nil, nil, nil))

	template, err := store.GetTemplate(t.Context(), db, 22, nil)
	require.NoError(t, err)
	require.NotNil(t, template)
	assert.Equal(t, int64(22), template.GetId())
	assert.Equal(t, "Global", template.GetTitle())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreCountTemplatesByCreatorJob(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_mailer_templates AS template`) +
		`(?s).*` + regexp.QuoteMeta(`template.creator_job = ?`)
	mock.ExpectQuery(expectedQuery).
		WithArgs("police").
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(4)))

	total, err := store.CountTemplatesByCreatorJob(t.Context(), db, "police")
	require.NoError(t, err)
	assert.Equal(t, int64(4), total)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreCreateTemplate(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)
	template := &mailertemplates.Template{
		EmailId:    7,
		Title:      "Welcome",
		Content:    nil,
		CreatorJob: func() *string { v := "police"; return &v }(),
	}

	expectedQuery := regexp.QuoteMeta(`INSERT INTO fivenet_mailer_templates`) +
		`(?s).*` + regexp.QuoteMeta(`VALUES (?, ?, ?, ?, ?)`)
	mock.ExpectExec(expectedQuery).
		WithArgs(int64(7), "Welcome", nil, "police", int32(3)).
		WillReturnResult(sqlmock.NewResult(11, 1))

	lastID, err := store.CreateTemplate(t.Context(), db, template, 3)
	require.NoError(t, err)
	assert.Equal(t, int64(11), lastID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreUpdateTemplate(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)
	template := &mailertemplates.Template{Id: 11, Title: "Updated", Content: nil}

	expectedQuery := regexp.QuoteMeta(
		`UPDATE fivenet_mailer_templates SET title = ?, content = ? WHERE fivenet_mailer_templates.id = ? LIMIT ?;`,
	)
	mock.ExpectExec(expectedQuery).
		WithArgs("Updated", nil, int64(11), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.UpdateTemplate(t.Context(), db, template))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreDeleteTemplate(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	expectedQuery := regexp.QuoteMeta(
		`UPDATE fivenet_mailer_templates SET deleted_at = CURRENT_TIMESTAMP WHERE fivenet_mailer_templates.id = ? LIMIT ?;`,
	)
	mock.ExpectExec(expectedQuery).
		WithArgs(int64(11), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.DeleteTemplate(t.Context(), db, 11))
	require.NoError(t, mock.ExpectationsWereMet())
}
