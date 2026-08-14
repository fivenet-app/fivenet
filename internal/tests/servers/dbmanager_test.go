package servers

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDBServerSharedSeedCloneIsolation(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	left := NewDBServer(ctx, t, true)
	right := NewDBServer(ctx, t, true)

	leftDB, err := left.DB()
	require.NoError(t, err)

	rightDB, err := right.DB()
	require.NoError(t, err)

	var leftName string
	require.NoError(t, leftDB.QueryRowContext(ctx, "SELECT DATABASE()").Scan(&leftName))
	require.NotEmpty(t, leftName)

	var rightName string
	require.NoError(t, rightDB.QueryRowContext(ctx, "SELECT DATABASE()").Scan(&rightName))
	require.NotEmpty(t, rightName)
	require.NotEqual(t, leftName, rightName)

	var leftUsers int
	require.NoError(
		t,
		leftDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM fivenet_user").Scan(&leftUsers),
	)
	require.Equal(t, 5, leftUsers)

	var rightUsers int
	require.NoError(
		t,
		rightDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM fivenet_user").Scan(&rightUsers),
	)
	require.Equal(t, 5, rightUsers)

	_, err = leftDB.ExecContext(ctx, "CREATE TABLE dbmanager_clone_marker (id INT PRIMARY KEY)")
	require.NoError(t, err)

	_, err = leftDB.ExecContext(ctx, "INSERT INTO dbmanager_clone_marker (id) VALUES (1)")
	require.NoError(t, err)

	var markerCount int
	require.NoError(
		t,
		leftDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM dbmanager_clone_marker").
			Scan(&markerCount),
	)
	require.Equal(t, 1, markerCount)

	require.Error(
		t,
		rightDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM dbmanager_clone_marker").
			Scan(&markerCount),
	)

	_, err = leftDB.ExecContext(
		ctx,
		`INSERT INTO fivenet_documents
			(category_id, title, summary, content_type, content, creator_job, state, public)
			VALUES (999999, 'fk check', 'fk check', 0, 'fk check', 'ambulance', 'Open', 0)`,
	)
	require.Error(t, err)

	left.Stop()

	var stillRightUsers int
	require.NoError(
		t,
		rightDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM fivenet_user").Scan(&stillRightUsers),
	)
	require.Equal(t, 5, stillRightUsers)

	right.Stop()

	fresh := NewDBServer(ctx, t, true)
	defer fresh.Stop()

	freshDB, err := fresh.DB()
	require.NoError(t, err)

	var freshUsers int
	require.NoError(
		t,
		freshDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM fivenet_user").Scan(&freshUsers),
	)
	require.Equal(t, 5, freshUsers)
}

func TestDBServerGeneratedColumnClone(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	srv := NewDBServer(ctx, t, true)
	defer srv.Stop()

	db, err := srv.DB()
	require.NoError(t, err)

	var extra string
	require.NoError(
		t,
		db.QueryRowContext(
			ctx,
			`SELECT EXTRA
			 FROM INFORMATION_SCHEMA.COLUMNS
			 WHERE TABLE_SCHEMA = DATABASE()
			   AND TABLE_NAME = 'fivenet_documents_categories'
			   AND COLUMN_NAME = 'sort_key'`,
		).Scan(&extra),
	)
	require.Contains(t, extra, "GENERATED")

	var name string
	var sortKey string
	require.NoError(
		t,
		db.QueryRowContext(
			ctx,
			`SELECT name, sort_key
			 FROM fivenet_documents_categories
			 WHERE id = 1`,
		).Scan(&name, &sortKey),
	)
	require.NotEmpty(t, name)
	require.Equal(t, name, sortKey)
}

func TestDBServerConcurrentSetup(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	probe := NewDBServer(ctx, t, true)
	probe.Stop()

	type result struct {
		name  string
		users int
		err   error
	}

	var (
		mu      sync.Mutex
		results []result
	)

	t.Cleanup(func() {
		require.Len(t, results, 2)

		names := map[string]int{}
		for _, res := range results {
			require.NoError(t, res.err)
			require.NotEmpty(t, res.name)
			require.Equal(t, 5, res.users)
			names[res.name] = res.users
		}

		require.Len(t, names, 2)
	})

	run := func(name string) {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			srv := NewDBServer(ctx, t, true)
			defer srv.Stop()

			db, err := srv.DB()
			if err != nil {
				mu.Lock()
				results = append(results, result{err: err})
				mu.Unlock()
				return
			}

			var dbName string
			if err := db.QueryRowContext(ctx, "SELECT DATABASE()").Scan(&dbName); err != nil {
				mu.Lock()
				results = append(results, result{err: err})
				mu.Unlock()
				return
			}

			var users int
			if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM fivenet_user").
				Scan(&users); err != nil {
				mu.Lock()
				results = append(results, result{err: err})
				mu.Unlock()
				return
			}

			mu.Lock()
			results = append(results, result{name: dbName, users: users})
			mu.Unlock()
		})
	}

	run("one")
	run("two")
}
