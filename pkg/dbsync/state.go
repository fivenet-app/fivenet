package dbsync

import (
	"errors"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type DBSyncState struct {
	mu sync.Mutex

	filepath string `yaml:"-"`

	States map[string]*TableSyncState `yaml:"states"`
}

type TableSyncState struct {
	// Used to track if the last id/strign needs to be reset
	IDField string `yaml:"idField"`

	Offset uint64  `yaml:"offset"`
	LastID *string `yaml:"lastId"`
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

	if err := os.WriteFile(s.filepath, out, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (s *DBSyncState) UpdateState(key string, state *TableSyncState) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.States[key] = state
}
