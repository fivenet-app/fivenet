package dbsyncconfig

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDBSyncConfigInitUsersCustomQueryRequiresWhereCondition(t *testing.T) {
	t.Parallel()

	missingWhere := "SELECT `id` AS `user.id` FROM `users` LIMIT $limit"
	cfg := &DBSyncConfig{
		Tables: DBSyncSourceTables{
			Users: UsersTable{
				DBSyncTable: DBSyncTable{
					Query: &missingWhere,
				},
			},
		},
	}

	err := cfg.Init()
	require.ErrorContains(
		t,
		err,
		"users table custom query must contain $whereCondition placeholder",
	)
}

func TestDBSyncConfigInitUsersCustomQueryWithWhereCondition(t *testing.T) {
	t.Parallel()

	withWhere := "SELECT `id` AS `user.id` FROM `users` $whereCondition LIMIT $limit"
	cfg := &DBSyncConfig{
		Tables: DBSyncSourceTables{
			Users: UsersTable{
				DBSyncTable: DBSyncTable{
					Query: &withWhere,
				},
			},
		},
	}

	require.NoError(t, cfg.Init())
}
