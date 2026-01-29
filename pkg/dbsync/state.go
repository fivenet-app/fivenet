package dbsync

import (
	"errors"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type DBSyncState struct {
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

func NewDBSyncState(logger *zap.Logger, filepath string) *DBSyncState {
	d := &DBSyncState{
		mu: sync.RWMutex{},

		logger: logger.Named("dbsync.state"),

		filepath: filepath,
	}

	d.Jobs = &TableSyncState{
		dss: d,
	}
	d.JobGrades = &TableSyncState{
		dss: d,
	}
	d.Licenses = &TableSyncState{
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

func (s *DBSyncState) Load() error {
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

	return nil
}

func (s *DBSyncState) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.save()
}

func (s *DBSyncState) save() error {
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
	dss *DBSyncState

	LastCheck *time.Time `yaml:"lastCheck"`
	Offset    int64      `yaml:"offset"`
	LastID    *string    `yaml:"lastId"`
	SyncedUp  bool       `yaml:"syncedUp"`
}

func (s *TableSyncState) Set(offset int64, lastId *string) {
	s.dss.mu.Lock()
	defer s.dss.mu.Unlock()

	now := time.Now()
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
