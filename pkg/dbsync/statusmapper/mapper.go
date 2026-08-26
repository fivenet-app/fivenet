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

func FromRuntimeState(state *dbsyncconfig.State, cfg *dbsyncconfig.DBSyncConfig) *Snapshot {
	if state == nil {
		return nil
	}

	var (
		usersEnabled       bool
		usersResyncIntv    *time.Duration
		vehiclesEnabled    bool
		vehiclesResyncIntv *time.Duration
	)
	if cfg != nil {
		usersEnabled = cfg.Tables.Users.Enabled
		usersResyncIntv = cfg.Tables.Users.ResyncInterval
		vehiclesEnabled = cfg.Tables.Vehicles.Enabled
		vehiclesResyncIntv = cfg.Tables.Vehicles.ResyncInterval
	}

	return &Snapshot{
		Tables: collectTables(
			tableSpec{"jobs", state.Jobs, cfg != nil && cfg.Tables.Jobs.Enabled},
			tableSpec{"licenses", state.Licenses, cfg != nil && cfg.Tables.Licenses.Enabled},
			tableSpec{"accounts", state.Accounts, cfg != nil && cfg.Tables.Accounts.Enabled},
			tableSpec{"users", state.Users, usersEnabled},
			tableSpec{"users_resync", state.UsersResync, resyncEnabled(usersEnabled, usersResyncIntv)},
			tableSpec{"vehicles", state.Vehicles, vehiclesEnabled},
			tableSpec{"vehicles_resync", state.VehiclesResync, resyncEnabled(vehiclesEnabled, vehiclesResyncIntv)},
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
	name    string
	state   *dbsyncconfig.TableSyncState
	enabled bool
}

func collectTables(specs ...tableSpec) []*dbsyncstate.DBSyncTableSyncState {
	tables := make([]*dbsyncstate.DBSyncTableSyncState, 0, len(specs))
	for _, spec := range specs {
		if table := fromConfigTable(spec); table != nil {
			tables = append(tables, table)
		}
	}

	return tables
}

func resyncEnabled(baseEnabled bool, interval *time.Duration) bool {
	return baseEnabled && interval != nil && *interval > 0
}

func fromConfigTable(spec tableSpec) *dbsyncstate.DBSyncTableSyncState {
	if spec.state == nil {
		return nil
	}

	out := &dbsyncstate.DBSyncTableSyncState{
		Table:   spec.name,
		Enabled: spec.enabled,
	}

	lastCheck := spec.state.GetLastCheck()
	lastID := spec.state.GetLastID()
	if lastCheck != nil || lastID != nil {
		out.Checkpoint = &dbsyncstate.DBSyncCheckpoint{
			LastId: lastID,
		}
		if lastCheck != nil {
			out.Checkpoint.LastCheck = toTimestamp(lastCheck)
		}
	}

	if lastSyncedAt := spec.state.GetLastSyncedAt(); lastSyncedAt != nil {
		out.LastSyncedAt = toTimestamp(lastSyncedAt)
	}

	if lastAttemptAt := spec.state.GetLastAttemptAt(); lastAttemptAt != nil {
		out.LastAttemptAt = toTimestamp(lastAttemptAt)
	}

	if lastError := spec.state.GetLastError(); lastError != nil && *lastError != "" {
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
