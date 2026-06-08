package calendar

import (
	"strings"
	"testing"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/go-jet/jet/v2/mysql"
)

func TestCalendarEntriesQueryOmitLimitWhenNil(t *testing.T) {
	stmt := calendarEntriesQuery(&userinfo.UserInfo{UserId: 1}, mysql.Bool(true), mysql.Bool(true), nil)

	sql, _ := stmt.Sql()
	if strings.Contains(sql, "LIMIT") {
		t.Fatalf("expected no limit in query, got %s", sql)
	}
}

func TestCalendarEntriesQueryUsesExplicitLimit(t *testing.T) {
	stmt := calendarEntriesQuery(&userinfo.UserInfo{UserId: 1}, mysql.Bool(true), mysql.Bool(true), int64Ptr(100))

	sql, args := stmt.Sql()
	if !strings.Contains(sql, "LIMIT ?") {
		t.Fatalf("expected explicit limit placeholder in query, got %s", sql)
	}
	if len(args) == 0 || args[len(args)-1] != int64(100) {
		t.Fatalf("expected limit argument 100, got %#v", args)
	}
}
