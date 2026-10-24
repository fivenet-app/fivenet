// Package filestore provides file storage utilities and helpers.
package filestore

import (
	"sync"

	"github.com/go-jet/jet/v2/mysql"
)

var (
	// tableListsMu protects access to the tablesList slice.
	tableListsMu = sync.Mutex{}

	// tablesList holds the list of joinInfo tables registered for the filestore.
	tablesList = []joinInfo{}
)

// AddTable adds a joinInfo table to the tablesList if the table is not nil.
// It acquires a lock to ensure thread-safe access to the tablesList.
func AddTable(tbl joinInfo) {
	RegisterFileReference(tbl.Table, tbl.FileCol)
}

// RegisterFileReference registers a table column that references fivenet_files.
// Registrations are deduplicated so package-level registrations and dynamically
// created file handlers can safely use the same registry.
func RegisterFileReference(tbl mysql.Table, fileCol mysql.ColumnInteger) {
	if tbl == nil || fileCol == nil {
		return
	}

	tableListsMu.Lock()
	defer tableListsMu.Unlock()

	for _, registered := range tablesList {
		if registered.Table.TableName() == tbl.TableName() &&
			registered.FileCol.Name() == fileCol.Name() {
			return
		}
	}

	tablesList = append(tablesList, joinInfo{
		Table:   tbl,
		FileCol: fileCol,
	})
}
