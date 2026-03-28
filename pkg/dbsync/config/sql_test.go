package dbsyncconfig

import (
	"testing"
	"time"

	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestBuildQueryFromColumns(t *testing.T) {
	t.Parallel()
	tests := []struct {
		tableName     string
		columns       map[string]string
		conditions    []string
		orderBy       []string
		limit         int64
		expectedQuery string
	}{
		{
			tableName: "users",
			columns: map[string]string{
				"id":       "user_id",
				"username": "user_name",
				"email":    "user_email",
			},
			conditions:    []string{"`updated_at` >= '2023-01-01 00:00:00'"},
			limit:         50,
			expectedQuery: "SELECT `user_email` AS `email`, `user_id` AS `id`, `user_name` AS `username`\nFROM `users`\nWHERE `updated_at` >= '2023-01-01 00:00:00'\nLIMIT 50;",
		},
		{
			tableName: "products",
			columns: map[string]string{
				"id":    "product_id",
				"name":  "product_name",
				"price": "product_price",
			},
			conditions:    []string{"`price` > 100"},
			limit:         20,
			expectedQuery: "SELECT `product_id` AS `id`, `product_name` AS `name`, `product_price` AS `price`\nFROM `products`\nWHERE `price` > 100\nLIMIT 20;",
		},
		{
			tableName: "user_licenses",
			columns: map[string]string{
				"license.type": "type",
				"license.name": "name_but_different",
			},
			conditions:    []string{},
			orderBy:       []string{"license.type", "license.name"},
			limit:         25,
			expectedQuery: "SELECT `name_but_different` AS `license.name`, `type` AS `license.type`\nFROM `user_licenses`\nORDER BY license.type, license.name\nLIMIT 25;",
		},
		{
			tableName: "vehicles",
			columns: map[string]string{
				"plate": "plate",
				"model": "-",
			},
			conditions:    []string{"`updated_at` >= '2023-01-01 00:00:00'"},
			orderBy:       []string{"plate"},
			limit:         50,
			expectedQuery: "SELECT `plate` AS `plate`\nFROM `vehicles`\nWHERE `updated_at` >= '2023-01-01 00:00:00'\nORDER BY plate\nLIMIT 50;",
		},
	}

	for _, test := range tests {
		query := buildQueryFromColumns(
			test.tableName,
			test.columns,
			test.conditions,
			test.limit,
			test.orderBy,
		)

		assert.Equal(t, test.expectedQuery, query, "Query did not match expected output")
	}
}

func TestPrepareStringQuery(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		query         string
		table         DBSyncTable
		state         *TableSyncState
		limit         int64
		expectedQuery string
	}{
		{
			name:          "No state provided",
			query:         "SELECT * FROM `users` $whereCondition LIMIT $limit",
			table:         DBSyncTable{},
			state:         nil,
			limit:         50,
			expectedQuery: "SELECT * FROM `users`  LIMIT 50",
		},
		{
			name:  "State with zero LastCheck",
			query: "SELECT * FROM `users` $whereCondition LIMIT $limit",
			table: DBSyncTable{
				UpdatedTimeColumn: utils.StrPtr("updated_at"),
			},
			state: &TableSyncState{
				LastCheck: nil,
			},
			limit:         20,
			expectedQuery: "SELECT * FROM `users`  LIMIT 20",
		},
		{
			name:  "State with valid LastCheck",
			query: "SELECT * FROM `users` $whereCondition LIMIT $limit",
			table: DBSyncTable{
				UpdatedTimeColumn: utils.StrPtr("updated_at"),
			},
			state: &TableSyncState{
				LastCheck: parseTime("2023-01-01 00:00:00"),
			},
			limit:         15,
			expectedQuery: "SELECT * FROM `users` WHERE `updated_at` >= '2023-01-01 00:00:00.000'\n LIMIT 15",
		},
		{
			name:  "State with cursor tuple",
			query: "SELECT * FROM `users` $whereCondition LIMIT $limit",
			table: DBSyncTable{
				UpdatedTimeColumn: utils.StrPtr("updated_at"),
			},
			state: &TableSyncState{
				LastCheck: parseTime("2023-01-01 00:00:00"),
				LastID:    utils.StrPtr("42"),
			},
			limit:         10,
			expectedQuery: "SELECT * FROM `users` WHERE (`updated_at` > '2023-01-01 00:00:00.000' OR (`updated_at` = '2023-01-01 00:00:00.000' AND `id` > 42))\n LIMIT 10",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := prepareStringQuery(
				test.query,
				test.table,
				test.state,
				test.limit,
				"id",
			)
			if result != test.expectedQuery {
				t.Errorf("Expected query:\n%s\nGot:\n%s", test.expectedQuery, result)
			}
		})
	}
}

func TestGetWhereConditionIDOnly(t *testing.T) {
	t.Parallel()
	table := DBSyncTable{}
	state := &TableSyncState{
		LastID: utils.StrPtr("XYZ-100"),
	}

	where := getWhereCondition(table, state, "plate")
	assert.Equal(t, "`plate` > 'XYZ-100'\n", where)
}

func TestGetWhereConditionBacktickedColumns(t *testing.T) {
	t.Parallel()
	t.Run("cursor column already backticked", func(t *testing.T) {
		t.Parallel()
		table := DBSyncTable{}
		state := &TableSyncState{
			LastID: utils.StrPtr("XYZ-100"),
		}

		where := getWhereCondition(table, state, "`plate`")
		assert.Equal(t, "`plate` > 'XYZ-100'\n", where)
	})

	t.Run("updated time column already backticked", func(t *testing.T) {
		t.Parallel()
		table := DBSyncTable{
			UpdatedTimeColumn: utils.StrPtr("`updated_at`"),
		}
		state := &TableSyncState{
			LastCheck: parseTime("2023-01-01 00:00:00"),
		}

		where := getWhereCondition(table, state, "id")
		assert.Equal(t, "`updated_at` >= '2023-01-01 00:00:00.000'\n", where)
	})

	t.Run("cursor and updated time columns already backticked", func(t *testing.T) {
		t.Parallel()
		table := DBSyncTable{
			UpdatedTimeColumn: utils.StrPtr("`updated_at`"),
		}
		state := &TableSyncState{
			LastCheck: parseTime("2023-01-01 00:00:00"),
			LastID:    utils.StrPtr("42"),
		}

		where := getWhereCondition(table, state, "`id`")
		assert.Equal(
			t,
			"(`updated_at` > '2023-01-01 00:00:00.000' OR (`updated_at` = '2023-01-01 00:00:00.000' AND `id` > 42))\n",
			where,
		)
	})
}

func parseTime(value string) *time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", value)
	return &t
}
