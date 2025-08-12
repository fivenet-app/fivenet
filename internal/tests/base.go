package tests

import (
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)

	// Root folder of the repository to be able to easily load,
	// e.g., database migration files.
	Root = filepath.Join(filepath.Dir(b), "../..")

	TestDataSQLPath = filepath.Join(Root, "internal/tests/testdata/sql")
)
