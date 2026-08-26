package statusmapper

import (
	"time"

	dbsyncstate "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/dbsync"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/settings"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	dbsyncconfig "github.com/fivenet-app/fivenet/v2026/pkg/dbsync/config"
)

type Snapshot struct {
	Tables []*dbsyncstate.DBSyncTableSyncState
}

func FromRuntimeState(state *dbsyncconfig.State) *Snapshot {
	if state == nil {
		return nil
	}

	return &Snapshot{
		Tables: collectTables(
			tableSpec{"jobs", state.Jobs},
			tableSpec{"licenses", state.Licenses},
			tableSpec{"accounts", state.Accounts},
			tableSpec{"users", state.Users},
			tableSpec{"users_resync", state.UsersResync},
			tableSpec{"vehicles", state.Vehicles},
			tableSpec{"vehicles_resync", state.VehiclesResync},
		),
	}
}

func FromClientSyncState(state *pbsync.ClientSyncState) *Snapshot {
	if state == nil {
		return nil
	}

	tables := state.GetTables()
	out := make([]*dbsyncstate.DBSyncTableSyncState, 0, len(tables))
	for _, table := range tables {
		if table == nil {
			continue
		}

		out = append(out, table)
	}

	return &Snapshot{Tables: out}
}

func (s *Snapshot) ToClientSyncState() *pbsync.ClientSyncState {
	if s == nil || len(s.Tables) == 0 {
		return nil
	}

	return &pbsync.ClientSyncState{Tables: s.Tables}
}

func (s *Snapshot) ToSettingsSyncState() *settings.DBSyncSyncState {
	if s == nil || len(s.Tables) == 0 {
		return nil
	}

	return &settings.DBSyncSyncState{Tables: s.Tables}
}

type tableSpec struct {
	name  string
	state *dbsyncconfig.TableSyncState
}

func collectTables(specs ...tableSpec) []*dbsyncstate.DBSyncTableSyncState {
	tables := make([]*dbsyncstate.DBSyncTableSyncState, 0, len(specs))
	for _, spec := range specs {
		if table := fromConfigTable(spec.name, spec.state); table != nil {
			tables = append(tables, table)
		}
	}

	return tables
}

func fromConfigTable(
	name string,
	state *dbsyncconfig.TableSyncState,
) *dbsyncstate.DBSyncTableSyncState {
	if state == nil {
		return nil
	}

	out := &dbsyncstate.DBSyncTableSyncState{
		Table: name,
	}

	lastCheck := state.GetLastCheck()
	lastID := state.GetLastID()
	if lastCheck != nil || lastID != nil {
		out.Checkpoint = &dbsyncstate.DBSyncCheckpoint{
			LastId: lastID,
		}
		if lastCheck != nil {
			out.Checkpoint.LastCheck = toTimestamp(lastCheck)
		}
	}

	if lastSyncedAt := state.GetLastSyncedAt(); lastSyncedAt != nil {
		out.LastSyncedAt = toTimestamp(lastSyncedAt)
	}

	if lastAttemptAt := state.GetLastAttemptAt(); lastAttemptAt != nil {
		out.LastAttemptAt = toTimestamp(lastAttemptAt)
	}

	if lastError := state.GetLastError(); lastError != nil && *lastError != "" {
		out.LastError = lastError
	}

	return out
}

func toTimestamp(t *time.Time) *timestamp.Timestamp {
	if t == nil {
		return nil
	}

	return timestamp.New(*t)
}
