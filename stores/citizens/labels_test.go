package citizensstore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	citizenslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/citizens/labels"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreListLabels(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db, &config.CustomDB{})

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_user_labels_job AS label`) +
		`(?s).*` + regexp.QuoteMeta(`label.deleted_at IS NULL`) +
		`(?s).*` + regexp.QuoteMeta(`label.job = ?`) +
		`(?s).*` + regexp.QuoteMeta(`ORDER BY label.sort_order ASC, label.sort_key ASC LIMIT ?;`)
	mock.ExpectQuery(expectedQuery).
		WithArgs("police", int64(20)).
		WillReturnRows(sqlmock.NewRows([]string{
			"label.id",
			"label.created_at",
			"label.sort_order",
			"label.name",
			"label.color",
			"label.icon",
			"label.settings",
		}).AddRow(
			int64(7),
			time.Unix(0, 0).UTC(),
			int32(3),
			"Patrol",
			"#00ff00",
			nil,
			nil,
		))

	labels, err := store.ListLabels(
		t.Context(),
		db,
		&userinfo.UserInfo{Superuser: true, Job: "police"},
		"",
		false,
		false,
		0,
		false,
	)
	require.NoError(t, err)
	require.Len(t, labels.GetList(), 1)
	assert.Equal(t, int64(7), labels.GetList()[0].GetId())
	assert.Equal(t, "Patrol", labels.GetList()[0].GetName())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListLabelsUsesVisibilityForNonSuperuser(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db, &config.CustomDB{})

	expectedQuery := `(?s).*WITH user_subjects AS.*visible_sources AS.*winning_visibility AS.*` +
		regexp.QuoteMeta(`SELECT label.id AS "label.id"`) +
		`.*` + regexp.QuoteMeta(`FROM ( SELECT DISTINCT fivenet_user_labels_job.id AS "id"`) +
		`.*` + regexp.QuoteMeta(`INNER JOIN winning_visibility ON (winning_visibility.target_id = fivenet_user_labels_job.id)`) +
		`.*` + regexp.QuoteMeta(`fivenet_user_labels_job.deleted_at IS NULL`) +
		`.*` + regexp.QuoteMeta(`fivenet_user_labels_job.name LIKE ?`) +
		`.*` + regexp.QuoteMeta(`ORDER BY label.sort_order ASC, label.sort_key ASC LIMIT ?;`)
	mock.ExpectQuery(expectedQuery).
		WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"label.id",
			"label.created_at",
			"label.sort_order",
			"label.name",
			"label.color",
			"label.icon",
			"label.settings",
		}))

	labels, err := store.ListLabels(
		t.Context(),
		db,
		&userinfo.UserInfo{UserId: 3, Job: "police", JobGrade: 16},
		"patrol",
		false,
		false,
		1,
		false,
	)
	require.NoError(t, err)
	require.Empty(t, labels.GetList())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreInsertLabel(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db, &config.CustomDB{})
	job := "police"
	icon := "shield"
	label := &citizenslabels.Label{
		Job:       &job,
		SortOrder: 4,
		Name:      "Active",
		Color:     "#ffffff",
		Icon:      &icon,
	}

	expectedQuery := regexp.QuoteMeta(`INSERT INTO fivenet_user_labels_job`) + `(?s).*`
	mock.ExpectExec(expectedQuery).
		WithArgs("police", int32(4), "Active", "#ffffff", "shield", nil).
		WillReturnResult(sqlmock.NewResult(42, 1))

	lastInsertID, err := store.InsertLabel(t.Context(), db, label)
	require.NoError(t, err)
	assert.Equal(t, int64(42), lastInsertID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreUpdateLabel(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db, &config.CustomDB{})
	job := "police"
	icon := "shield"
	label := &citizenslabels.Label{
		Id:    7,
		Job:   &job,
		Name:  "Updated",
		Color: "#123456",
		Icon:  &icon,
	}

	expectedQuery := regexp.QuoteMeta(
		`UPDATE fivenet_user_labels_job`,
	) + `(?s).*` + regexp.QuoteMeta(
		`WHERE`,
	) + `(?s).*`
	mock.ExpectExec(expectedQuery).
		WithArgs("Updated", "#123456", "shield", nil, int64(7), "police", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.UpdateLabel(t.Context(), db, label, job))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreDeleteLabel(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db, &config.CustomDB{})
	deletedAt := timestamp.Now()

	expectedQuery := regexp.QuoteMeta(
		`UPDATE fivenet_user_labels_job`,
	) + `(?s).*` + regexp.QuoteMeta(
		`deleted_at`,
	) + `(?s).*`
	mock.ExpectExec(expectedQuery).
		WithArgs(sqlmock.AnyArg(), int64(7), "police", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.DeleteLabel(t.Context(), db, "police", 7, deletedAt))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreReorderLabels(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db, &config.CustomDB{})
	labelIDs := []int64{7, 4, 9}

	lookupQuery := regexp.QuoteMeta(`FROM fivenet_user_labels_job`) +
		`(?s).*` + regexp.QuoteMeta(`fivenet_user_labels_job.job = ?`) +
		`(?s).*` + regexp.QuoteMeta(`fivenet_user_labels_job.deleted_at IS NULL`) +
		`(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)
	mock.ExpectQuery(lookupQuery).
		WithArgs("police", int64(3)).
		WillReturnRows(sqlmock.NewRows([]string{"citizen_label.id"}).AddRow(int64(7)).AddRow(int64(4)).AddRow(int64(9)))

	mock.ExpectBegin()
	for idx, labelID := range labelIDs {
		execQuery := regexp.QuoteMeta(
			`UPDATE fivenet_user_labels_job SET sort_order = ? WHERE ( (fivenet_user_labels_job.id = ?) AND (fivenet_user_labels_job.job = ?) AND (fivenet_user_labels_job.deleted_at IS NULL) ) LIMIT ?;`,
		)
		mock.ExpectExec(execQuery).
			WithArgs(int32(idx), labelID, "police", int64(1)).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}
	mock.ExpectCommit()

	require.NoError(t, store.ReorderLabels(t.Context(), "police", labelIDs))
	require.NoError(t, mock.ExpectationsWereMet())
}
