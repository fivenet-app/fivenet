package calendar

import (
	"testing"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalendarEntriesQueryOmitLimitWhenNil(t *testing.T) {
	t.Parallel()

	stmt := calendarEntriesQuery(
		&userinfo.UserInfo{UserId: 1},
		mysql.Bool(true),
		mysql.Bool(true),
		nil,
	)

	sql, _ := stmt.Sql()
	assert.NotContains(t, sql, "LIMIT", "expected no limit in query, got %s", sql)
}

func TestCalendarEntriesQueryUsesExplicitLimit(t *testing.T) {
	t.Parallel()

	stmt := calendarEntriesQuery(
		&userinfo.UserInfo{UserId: 1},
		mysql.Bool(true),
		mysql.Bool(true),
		new(maxCalendarEntriesLimit),
	)

	sql, args := stmt.Sql()
	require.Contains(t, sql, "LIMIT ?", "expected explicit limit placeholder in query, got %s", sql)
	require.NotEmpty(t, args, "expected limit arguments")
	assert.Equal(
		t,
		maxCalendarEntriesLimit,
		args[len(args)-1],
		"expected limit argument %d, got %#v",
		maxCalendarEntriesLimit,
		args,
	)
}
