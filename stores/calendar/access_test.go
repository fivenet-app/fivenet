package calendarstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/stretchr/testify/require"
)

func TestCheckIfUserHasAccessToCalendarEntryIDsUsesEntryID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db).(*Store)
	mock.ExpectQuery(regexp.QuoteMeta(`calendar_entry.id IN (?)`)+`(?s).*`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "entry_id"}).AddRow(int64(42), int64(42)))

	ids, err := store.CheckIfUserHasAccessToCalendarEntryIDs(
		t.Context(),
		&userinfo.UserInfo{UserId: 7, Job: "police"},
		true,
		42,
	)
	require.NoError(t, err)
	require.Len(t, ids, 1)
	require.Equal(t, int64(42), ids[0])
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckIfUserHasAccessToCalendarIDsAllowsSuperuserBirthdayCalendar(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db).(*Store)
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_calendar AS calendar`)+`(?s).*`+regexp.QuoteMeta(`calendar.id IN (?)`)).
		WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"calendar.id"}).AddRow(int64(99)))

	out, err := store.CheckIfUserHasAccessToCalendarIDs(
		t.Context(),
		&userinfo.UserInfo{UserId: 1, Superuser: true},
		calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW,
		true,
		99,
	)
	require.NoError(t, err)
	require.Len(t, out, 1)
	require.Equal(t, int64(99), out[0].ID)
	require.NoError(t, mock.ExpectationsWereMet())
}
