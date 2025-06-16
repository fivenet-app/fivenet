// Package filestore provides file storage utilities and helpers.
package filestore

import "sync"

var (
	// tableListsMu protects access to the tablesList slice.
	tableListsMu = sync.Mutex{}

	// tablesList holds the list of joinInfo tables registered for the filestore.
	tablesList = []joinInfo{}
)

// AddTable adds a joinInfo table to the tablesList if the table is not nil.
// It acquires a lock to ensure thread-safe access to the tablesList.
func AddTable(tbl joinInfo) {
	if tbl.Table == nil {
		return
	}

	tableListsMu.Lock()
	defer tableListsMu.Unlock()

	tablesList = append(tablesList, tbl)
}
