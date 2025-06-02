package filestore

import "sync"

var (
	// Mutexes to protect access to the table list maps
	tableListsMu = sync.Mutex{}

	tablesList = []joinInfo{}
)

func AddTable(tbl joinInfo) {
	tableListsMu.Lock()
	defer tableListsMu.Unlock()

	tablesList = append(tablesList, tbl)
}
