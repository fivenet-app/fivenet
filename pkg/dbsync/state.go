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
	mu sync.Mutex

	logger *zap.Logger

	filepath string `yaml:"-"`

	Users         *TableSyncState `yaml:"users"`
	Jobs          *TableSyncState `yaml:"jobs"`
	JobGrades     *TableSyncState `yaml:"jobGrades"`
	Licenses      *TableSyncState `yaml:"licenses"`
	UserLicenses  *TableSyncState `yaml:"userLicenses"`
	OwnedVehicles *TableSyncState `yaml:"ownedVehicles"`
}

type TableSyncState struct {
	dss *DBSyncState

	// Used to track if the last id/string needs to be reset
	IDField string `yaml:"idField"`

	LastCheck time.Time `yaml:"lastCheck"`
	Offset    uint64    `yaml:"offset"`
	LastID    *string   `yaml:"lastId"`
}

func NewDBSyncState(logger *zap.Logger, filepath string) *DBSyncState {
	d := &DBSyncState{
		mu:       sync.Mutex{},
		filepath: filepath,
	}

	d.Users = &TableSyncState{
		dss: d,
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
	d.UserLicenses = &TableSyncState{
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

	out, err := yaml.Marshal(s)
	if err != nil {
		return err
	}

	if err := os.WriteFile(s.filepath, out, 0o640); err != nil {
		return err
	}

	return nil
}

func (s *TableSyncState) Set(idField string, offset uint64, lastId *string) {
	s.LastCheck = time.Now()
	s.IDField = idField
	s.Offset = offset
	s.LastID = lastId

	if err := s.dss.Save(); err != nil {
		s.dss.logger.Error("failed to save state", zap.Error(err))
	}
}
