package dbsyncconfig

import (
	"errors"
	"os"
	"sync"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var StateModule = fx.Module("dbsync.state",
	fx.Provide(
		NewState,
	),
)

type StateParams struct {
	fx.In

	Logger *zap.Logger
	Config *Config
}

type State struct {
	mu sync.RWMutex

	logger *zap.Logger

	filepath string `yaml:"-"`

	Jobs      *TableSyncState `yaml:"jobs"`
	JobGrades *TableSyncState `yaml:"jobGrades"`
	Licenses  *TableSyncState `yaml:"licenses"`

	Accounts      *TableSyncState `yaml:"accounts"`
	Users         *TableSyncState `yaml:"users"`
	OwnedVehicles *TableSyncState `yaml:"ownedVehicles"`
}

func NewState(p StateParams) *State {
	d := &State{
		mu: sync.RWMutex{},

		logger: p.Logger.Named("dbsync.state"),

		filepath: p.Config.Load().StateFile,
	}

	d.Jobs = &TableSyncState{
		dss: d,
	}
	d.Licenses = &TableSyncState{
		dss: d,
	}

	d.Accounts = &TableSyncState{
		dss: d,
	}
	d.Users = &TableSyncState{
		dss: d,
	}
	d.OwnedVehicles = &TableSyncState{
		dss: d,
	}

	return d
}

func (s *State) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	out, err := os.ReadFile(s.filepath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	if err := yaml.Unmarshal(out, s); err != nil {
		return err
	}

	if s.Jobs == nil {
		s.Jobs = &TableSyncState{dss: s}
	}
	if s.Licenses == nil {
		s.Licenses = &TableSyncState{dss: s}
	}

	if s.Accounts == nil {
		s.Accounts = &TableSyncState{dss: s}
	}
	if s.Users == nil {
		s.Users = &TableSyncState{dss: s}
	}
	if s.OwnedVehicles == nil {
		s.OwnedVehicles = &TableSyncState{dss: s}
	}

	return nil
}

func (s *State) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.save()
}

func (s *State) save() error {
	out, err := yaml.Marshal(s)
	if err != nil {
		return err
	}

	if err := os.WriteFile(s.filepath, out, 0o640); err != nil {
		return err
	}

	return nil
}

type TableSyncState struct {
	dss *State

	LastCheck *time.Time `yaml:"lastCheck"`
	Offset    int64      `yaml:"offset"`
	LastID    *string    `yaml:"lastId"`
	SyncedUp  bool       `yaml:"syncedUp"`
}

func (s *TableSyncState) Set(offset int64, lastId *string) {
	s.dss.mu.Lock()
	defer s.dss.mu.Unlock()

	// Subtract 1 minute to account for potential clock skew and ensure no records are missed
	now := time.Now().Add(-1 * time.Minute)
	s.LastCheck = &now
	s.Offset = offset
	s.LastID = lastId

	if err := s.dss.save(); err != nil {
		s.dss.logger.Error("failed to save state", zap.Error(err))
	}
}

func (s *TableSyncState) GetSyncedUp() bool {
	s.dss.mu.RLock()
	defer s.dss.mu.RUnlock()
	return s.SyncedUp
}

func (s *TableSyncState) SetSyncedUp(syncedUp bool) {
	s.dss.mu.Lock()
	defer s.dss.mu.Unlock()
	s.SyncedUp = syncedUp
}

func (s *TableSyncState) SetLastCheck(t *time.Time) {
	s.dss.mu.Lock()
	defer s.dss.mu.Unlock()
	s.LastCheck = t
}

func (s *TableSyncState) GetOffset() int64 {
	s.dss.mu.RLock()
	defer s.dss.mu.RUnlock()
	return s.Offset
}
