package settingsstore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/stretchr/testify/require"
)

func TestStoreViewAuditLogAppliesFiltersAndSort(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	search := "needle"
	pageSize := int64(20)
	req := ViewAuditLogOptions{
		Pagination: &database.PaginationRequest{PageSize: &pageSize},
		UserIDs:    []int32{3},
		From:       timestamp.New(time.Unix(100, 0)),
		To:         timestamp.New(time.Unix(200, 0)),
		Services:   []string{"settings"},
		Methods:    []string{"UpdateAppConfig"},
		Search:     search,
		Actions:    []audit.EventAction{audit.EventAction_EVENT_ACTION_UPDATED},
		Results:    []audit.EventResult{audit.EventResult_EVENT_RESULT_SUCCEEDED},
		Sort:       &database.Sort{Columns: []*database.SortByColumn{{Id: "service", Desc: true}}},
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(audit_entry.id) AS "data_count.total" FROM fivenet_audit_log AS audit_entry WHERE`) + `(?s).*` + regexp.QuoteMeta(`audit_entry.user_id IN (?)`) + `(?s).*` + regexp.QuoteMeta(`audit_entry.created_at >= CAST(? AS DATETIME)`) + `(?s).*` + regexp.QuoteMeta(`audit_entry.created_at <= CAST(? AS DATETIME)`) + `(?s).*` + regexp.QuoteMeta(`audit_entry.service IN (?)`) + `(?s).*` + regexp.QuoteMeta(`audit_entry.method IN (?)`) + `(?s).*` + regexp.QuoteMeta(`audit_entry.action IN (?)`) + `(?s).*` + regexp.QuoteMeta(`audit_entry.result IN (?)`) + `(?s).*` + regexp.QuoteMeta(`MATCH(audit_entry.data) AGAINST (? IN BOOLEAN MODE)`)).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(2)))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT audit_entry.id AS "audit_entry.id", audit_entry.created_at AS "audit_entry.created_at", audit_entry.user_id AS "audit_entry.user_id", audit_entry.user_job AS "audit_entry.user_job", audit_entry.target_user_id AS "audit_entry.target_user_id", audit_entry.service AS "audit_entry.service", audit_entry.method AS "audit_entry.method", audit_entry.action AS "audit_entry.action", audit_entry.result AS "audit_entry.result", audit_entry.meta AS "audit_entry.meta", audit_entry.data AS "audit_entry.data", user_short.id AS "user_short.id", user_short.identifier AS "user_short.identifier", user_short.job AS "user_short.job", user_short.job_grade AS "user_short.job_grade", user_short.firstname AS "user_short.firstname", user_short.lastname AS "user_short.lastname", user_short.dateofbirth AS "user_short.dateofbirth" FROM fivenet_audit_log AS audit_entry LEFT JOIN fivenet_user AS user_short ON`) + `(?s).*` + regexp.QuoteMeta(`ORDER BY audit_entry.service DESC`)).
		WillReturnRows(sqlmock.NewRows([]string{}))

	resp, err := store.ViewAuditLog(t.Context(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.GetPagination())
	require.Empty(t, resp.GetLogs())
	require.NoError(t, mock.ExpectationsWereMet())
}
