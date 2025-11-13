package dbsync

import (
	"testing"
	"time"
)

func TestBuildQueryFromColumns(t *testing.T) {
	tests := []struct {
		tableName     string
		columns       map[string]string
		conditions    []string
		offset        int64
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
			offset:        10,
			limit:         50,
			expectedQuery: "SELECT `user_email` AS `email`, `user_id` AS `id`, `user_name` AS `username`\nFROM `users`\nWHERE `updated_at` >= '2023-01-01 00:00:00'\nLIMIT 50 OFFSET 10",
		},
		{
			tableName: "products",
			columns: map[string]string{
				"id":    "product_id",
				"name":  "product_name",
				"price": "product_price",
			},
			conditions:    []string{"`price` > 100"},
			offset:        0,
			limit:         20,
			expectedQuery: "SELECT `product_id` AS `id`, `product_name` AS `name`, `product_price` AS `price`\nFROM `products`\nWHERE `price` > 100\nLIMIT 20 OFFSET 0",
		},
		{
			tableName: "user_licenses",
			columns: map[string]string{
				"license.type": "type",
				"license.name": "name_but_different",
			},
			conditions:    []string{},
			offset:        10,
			limit:         25,
			expectedQuery: "SELECT `name_but_different` AS `license.name`, `type` AS `license.type`\nFROM `user_licenses`\nLIMIT 25 OFFSET 10",
		},
	}

	for _, test := range tests {
		query := buildQueryFromColumns(
			test.tableName,
			test.columns,
			test.conditions,
			test.offset,
			test.limit,
		)

		if query != test.expectedQuery {
			t.Errorf(
				"For table %s, expected query:\n%s\nGot:\n%s",
				test.tableName,
				test.expectedQuery,
				query,
			)
		}
	}
}

func TestPrepareStringQuery(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		table         DBSyncTable
		state         *TableSyncState
		offset        int64
		limit         int64
		expectedQuery string
	}{
		{
			name:          "No state provided",
			query:         "SELECT * FROM `users` $whereCondition LIMIT $limit OFFSET $offset",
			table:         DBSyncTable{},
			state:         nil,
			offset:        10,
			limit:         50,
			expectedQuery: "SELECT * FROM `users`  LIMIT 50 OFFSET 10",
		},
		{
			name:  "State with zero LastCheck",
			query: "SELECT * FROM `users` $whereCondition LIMIT $limit OFFSET $offset",
			table: DBSyncTable{
				UpdatedTimeColumn: ptr("updated_at"),
			},
			state: &TableSyncState{
				LastCheck: nil,
			},
			offset:        0,
			limit:         20,
			expectedQuery: "SELECT * FROM `users`  LIMIT 20 OFFSET 0",
		},
		{
			name:  "State with valid LastCheck",
			query: "SELECT * FROM `users` $whereCondition LIMIT $limit OFFSET $offset",
			table: DBSyncTable{
				UpdatedTimeColumn: ptr("updated_at"),
			},
			state: &TableSyncState{
				LastCheck: parseTime("2023-01-01 00:00:00"),
			},
			offset:        5,
			limit:         15,
			expectedQuery: "SELECT * FROM `users` WHERE `updated_at` >= '2023-01-01 00:00:00'\n LIMIT 15 OFFSET 5",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := prepareStringQuery(
				test.query,
				test.table,
				test.state,
				test.offset,
				test.limit,
			)
			if result != test.expectedQuery {
				t.Errorf("Expected query:\n%s\nGot:\n%s", test.expectedQuery, result)
			}
		})
	}
}

func ptr(s string) *string {
	return &s
}

func parseTime(value string) *time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", value)
	return &t
}
