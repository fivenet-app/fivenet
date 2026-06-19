package calendarstore

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
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

	store := New(db)
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_calendar AS calendar`) + `(?s).*` + regexp.QuoteMeta(`LEFT JOIN fivenet_user AS creator ON`) + `(?s).*` + regexp.QuoteMeta(`calendar.deleted_at IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(9)))

	total, err := store.CountCalendars(t.Context(), ListQuery{})
	require.NoError(t, err)
	assert.Equal(t, int64(9), total)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCountCalendarsStmtUsesAliasedCalendarColumnForBirthdayAccess(t *testing.T) {
	t.Parallel()

	store := New(new(sql.DB)).(*Store)
	stmt := store.countCalendarsStmt(
		ListQuery{
			UserInfo: &userinfo.UserInfo{
				UserId:    7,
				Job:       "police",
				Superuser: true,
			},
		},
		table.FivenetUser.AS("creator"),
	)

	sql, _ := stmt.Sql()
	assert.Contains(t, sql, "calendar.id")
	assert.NotContains(t, sql, "fivenet_calendar.id")
}

func TestListCalendarsStmtOrdersByCalendarIds(t *testing.T) {
	t.Parallel()

	store := New(new(sql.DB)).(*Store)
	stmt := store.listCalendarsStmt(
		ListQuery{
			UserInfo:    &userinfo.UserInfo{UserId: 7, Superuser: true},
			CalendarIDs: []int64{4, 9},
		},
		7,
		table.FivenetUser.AS("creator"),
		table.FivenetFiles.AS("profile_picture"),
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

	store := New(new(sql.DB)).(*Store)
	stmt := store.listCalendarsStmt(
		ListQuery{UserInfo: &userinfo.UserInfo{UserId: 7, Job: "police"}, OnlyPublic: true},
		7,
		table.FivenetUser.AS("creator"),
		table.FivenetFiles.AS("profile_picture"),
		20,
		40,
	)

	sql, _ := stmt.Sql()
	assert.Contains(t, sql, "calendar.public IS TRUE")
}

func TestGetCalendarStmtIncludesCreatorJoins(t *testing.T) {
	t.Parallel()

	store := New(new(sql.DB)).(*Store)
	stmt := store.getCalendarStmt(&userinfo.UserInfo{UserId: 7}, mysql.Bool(true))

	sql, _ := stmt.Sql()
	assert.Contains(t, sql, "LEFT JOIN fivenet_user AS creator ON")
	assert.Contains(t, sql, "LEFT JOIN fivenet_calendar_subs AS calendar_sub ON")
	assert.Contains(t, sql, "LIMIT ?")
}
