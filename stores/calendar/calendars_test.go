package calendarstore

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	calendarresource "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCountCalendarsReturnsCount(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(testParams(db))
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_calendar AS calendar`) + `(?s).*` + regexp.QuoteMeta(`calendar.deleted_at IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(9)))

	total, err := store.CountCalendars(t.Context(), ListQuery{})
	require.NoError(t, err)
	assert.Equal(t, int64(9), total)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCountCalendarsStmtUsesAliasedCalendarColumnForBirthdayAccess(t *testing.T) {
	t.Parallel()

	store := New(testParams(new(sql.DB))).(*Store)
	stmt := store.countCalendarsStmt(
		ListQuery{
			UserInfo: &userinfo.UserInfo{
				UserId:    7,
				Job:       "police",
				Superuser: true,
			},
		},
	)

	sql, _ := stmt.Sql()
	assert.Contains(t, sql, "calendar.id")
	assert.NotContains(t, sql, "fivenet_calendar.id")
}

func TestListCalendarsStmtOrdersByCalendarIds(t *testing.T) {
	t.Parallel()

	store := New(testParams(new(sql.DB))).(*Store)
	stmt := store.listCalendarsStmt(
		ListQuery{
			UserInfo:    &userinfo.UserInfo{UserId: 7, Superuser: true},
			CalendarIDs: []int64{4, 9},
		},
		7,
		20,
		40,
	)

	sql, args := stmt.Sql()
	assert.Contains(t, sql, "calendar.id IN (?, ?) DESC")
	assert.Contains(t, sql, "LIMIT ?")
	assert.Contains(t, sql, "OFFSET ?")
	assert.NotEmpty(t, args)
}

func TestListCalendarsStmtFiltersOnlyPublicWhenRequested(t *testing.T) {
	t.Parallel()

	store := New(testParams(new(sql.DB))).(*Store)
	stmt := store.listCalendarsStmt(
		ListQuery{UserInfo: &userinfo.UserInfo{UserId: 7, Job: "police"}, OnlyPublic: true},
		7,
		20,
		40,
	)

	sql, _ := stmt.Sql()
	assert.Contains(t, sql, "calendar.public IS TRUE")
}

func TestGetCalendarStmtIncludesCreatorJoins(t *testing.T) {
	t.Parallel()

	store := New(testParams(new(sql.DB))).(*Store)
	stmt := store.getCalendarStmt(&userinfo.UserInfo{UserId: 7}, mysql.Bool(true))

	sql, _ := stmt.Sql()
	assert.Contains(t, sql, "LEFT JOIN fivenet_calendar_subs AS calendar_sub ON")
	assert.Contains(t, sql, "LIMIT ?")
}

func TestCreateCalendarIncludesIconColumn(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(testParams(db)).(*Store)
	icon := "CalendarStarIcon"
	kind := calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_UNSPECIFIED

	mock.ExpectQuery(`(?s)SELECT fivenet_calendar.id AS "id" FROM fivenet_calendar .*`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectExec(`(?s)INSERT INTO fivenet_calendar .*system_kind.*icon.*`).
		WillReturnResult(sqlmock.NewResult(42, 1))

	cal := &calendarresource.Calendar{
		Name:       "Test calendar",
		Color:      "blue",
		Icon:       &icon,
		CreatorJob: "police",
	}
	cal.SetSystemKind(kind)

	id, err := store.CreateCalendar(
		t.Context(),
		db,
		cal,
		&userinfo.UserInfo{UserId: 9, Job: "police"},
	)
	require.NoError(t, err)
	assert.Equal(t, int64(42), id)
	require.True(t, cal.HasSystemKind())
	assert.Equal(t, kind, cal.GetSystemKind())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateCalendarIncludesIconColumn(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(testParams(db)).(*Store)
	icon := "CalendarStarIcon"

	mock.ExpectExec(`(?s)UPDATE fivenet_calendar .*icon .*`).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = store.UpdateCalendar(
		t.Context(),
		db,
		&calendarresource.Calendar{
			Id:         77,
			Name:       "Updated calendar",
			Color:      "green",
			Icon:       &icon,
			CreatorJob: "police",
		},
	)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
