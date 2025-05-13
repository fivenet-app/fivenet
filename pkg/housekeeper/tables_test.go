package housekeeper

import (
	"strings"
	"testing"

	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/stretchr/testify/assert"
)

func TestGetQuery_NoSource_HardDelete(t *testing.T) {
	targetTable := jet.NewTable("", "target_table", "")
	targetJobColumn := jet.StringColumn("job_column")

	table := &JobTable{
		TargetTable:     targetTable,
		TargetJobColumn: targetJobColumn,
	}

	query := table.GetQuery("test_job")
	querySQL, args := query.Sql()
	querySQL = strings.ReplaceAll(strings.TrimSpace(querySQL), "\n", " ")
	expectedSQL := "DELETE FROM target_table WHERE (job_column = ?) LIMIT ?;"

	assert.Equal(t, expectedSQL, querySQL)
	assert.Equal(t, []interface{}{"test_job", int64(1000)}, args)
}

func TestGetQuery_NoSource_SoftDelete(t *testing.T) {
	targetTable := jet.NewTable("", "target_table", "")
	targetJobColumn := jet.StringColumn("job_column")
	targetDeletedAtColumn := jet.TimestampColumn("deleted_at")

	table := &JobTable{
		TargetTable:           targetTable,
		TargetJobColumn:       targetJobColumn,
		TargetDeletedAtColumn: targetDeletedAtColumn,
	}

	query := table.GetQuery("test_job")
	querySQL, args := query.Sql()
	querySQL = strings.ReplaceAll(strings.TrimSpace(querySQL), "\n", " ")
	expectedSQL := "UPDATE target_table SET deleted_at = CURRENT_TIMESTAMP WHERE (           (job_column = ?)               AND deleted_at IS NULL       ) LIMIT ?;"

	assert.Equal(t, expectedSQL, querySQL)
	assert.Equal(t, []interface{}{"test_job", int64(1000)}, args)
}

func TestGetQuery_WithSource_HardDelete(t *testing.T) {
	sourceTable := jet.NewTable("", "source_table", "")
	sourceJobColumn := jet.StringColumn("source_job_column")
	sourceIDColumn := jet.IntegerColumn("source_id_column")

	targetTable := jet.NewTable("", "target_table", "")
	targetSourceIDColumn := jet.IntegerColumn("target_source_id_column")

	table := &JobTable{
		Source: &JobTableSource{
			SourceTable:     sourceTable,
			SourceJobColumn: sourceJobColumn,
			SourceIDColumn:  sourceIDColumn,
		},
		TargetTable:          targetTable,
		TargetSourceIDColumn: targetSourceIDColumn,
	}

	query := table.GetQuery("test_job")
	querySQL, args := query.Sql()
	querySQL = strings.ReplaceAll(strings.TrimSpace(querySQL), "\n", " ")
	expectedSQL := "DELETE FROM target_table USING source_table WHERE (           (source_job_column = ?)               AND (source_id_column = target_source_id_column)       ) LIMIT ?;"

	assert.Equal(t, expectedSQL, querySQL)
	assert.Equal(t, []interface{}{"test_job", int64(1000)}, args)
}

func TestGetQuery_WithSource_SoftDelete(t *testing.T) {
	sourceTable := jet.NewTable("", "source_table", "")
	sourceJobColumn := jet.StringColumn("source_job_column")
	sourceIDColumn := jet.IntegerColumn("source_id_column")
	sourceDeletedAtColumn := jet.TimestampColumn("source_deleted_at")

	targetTable := jet.NewTable("", "target_table", "")
	targetSourceIDColumn := jet.IntegerColumn("target_source_id_column")
	targetDeletedAtColumn := jet.TimestampColumn("target_deleted_at")

	table := &JobTable{
		Source: &JobTableSource{
			SourceTable:           sourceTable,
			SourceJobColumn:       sourceJobColumn,
			SourceIDColumn:        sourceIDColumn,
			SourceDeletedAtColumn: sourceDeletedAtColumn,
		},
		TargetTable:           targetTable,
		TargetSourceIDColumn:  targetSourceIDColumn,
		TargetDeletedAtColumn: targetDeletedAtColumn,
	}

	query := table.GetQuery("test_job")
	querySQL, args := query.Sql()
	querySQL = strings.ReplaceAll(strings.TrimSpace(querySQL), "\n", " ")
	expectedSQL := "UPDATE target_table INNER JOIN source_table ON (source_id_column = target_source_id_column) SET target_deleted_at = CURRENT_TIMESTAMP WHERE (           ((source_job_column = ?) AND source_deleted_at IS NULL)               AND (source_id_column = target_source_id_column)       );"

	assert.Equal(t, expectedSQL, querySQL)
	assert.Equal(t, []interface{}{"test_job", int64(1000)}, args)
}
