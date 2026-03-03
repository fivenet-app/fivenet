package dbsyncconfig

import (
	"context"
	"errors"
	"fmt"
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

	LC fx.Lifecycle

	Logger *zap.Logger
	Config *Config
}

type State struct {
	mu sync.RWMutex

	logger *zap.Logger

	filepath string `yaml:"-"`

	Jobs     *TableSyncState `yaml:"jobs"`
	Licenses *TableSyncState `yaml:"licenses"`

	Accounts *TableSyncState `yaml:"accounts"`

	Users       *TableSyncState `yaml:"users"`
	UsersResync *TableSyncState `yaml:"usersResync"`

	Vehicles       *TableSyncState `yaml:"vehicles"`
	VehiclesResync *TableSyncState `yaml:"vehiclesResync"`
}

func NewState(p StateParams) (*State, error) {
	s := &State{
		mu: sync.RWMutex{},

		logger: p.Logger.Named("dbsync.state"),

		filepath: p.Config.Load().StateFile,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		// Load dbsync state from file if exists
		if err := s.load(); err != nil {
			return fmt.Errorf("failed to load dbsync state. %w", err)
		}

		return nil
	}))
	p.LC.Append(fx.StopHook(func(ctxStartup context.Context) error {
		if err := s.Save(); err != nil {
			return fmt.Errorf("failed to save state. %w", err)
		}

		return nil
	}))

	return s, nil
}

func (s *State) load() error {
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
	} else {
		s.Jobs.dss = s
	}
	if s.Licenses == nil {
		s.Licenses = &TableSyncState{dss: s}
	} else {
		s.Licenses.dss = s
	}

	if s.Accounts == nil {
		s.Accounts = &TableSyncState{dss: s}
	} else {
		s.Accounts.dss = s
	}

	if s.Users == nil {
		s.Users = &TableSyncState{dss: s}
	} else {
		s.Users.dss = s
	}
	if s.UsersResync == nil {
		s.UsersResync = &TableSyncState{dss: s}
	} else {
		s.UsersResync.dss = s
	}

	if s.Vehicles == nil {
		s.Vehicles = &TableSyncState{dss: s}
	} else {
		s.Vehicles.dss = s
	}
	if s.VehiclesResync == nil {
		s.VehiclesResync = &TableSyncState{dss: s}
	} else {
		s.VehiclesResync.dss = s
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
	LastID    *string    `yaml:"lastId"`
}

func (s *TableSyncState) SetCursor(lastCheck *time.Time, lastId *string) {
	if s == nil {
		return
	}
	if s.dss == nil {
		s.LastCheck = lastCheck
		s.LastID = lastId
		return
	}

	s.dss.mu.Lock()
	defer s.dss.mu.Unlock()

	s.LastCheck = lastCheck
	s.LastID = lastId

	if err := s.dss.save(); err != nil {
		s.dss.logger.Error("failed to save state", zap.Error(err))
	}
}

func (s *TableSyncState) ResetCursor() {
	s.SetCursor(nil, nil)
}

func (s *TableSyncState) SetLastCheck(t *time.Time) {
	if s == nil {
		return
	}
	if s.dss == nil {
		s.LastCheck = t
		return
	}

	s.dss.mu.Lock()
	defer s.dss.mu.Unlock()
	s.LastCheck = t

	if err := s.dss.save(); err != nil {
		s.dss.logger.Error("failed to save state", zap.Error(err))
	}
}

func (s *TableSyncState) GetLastCheck() *time.Time {
	if s == nil {
		return nil
	}
	if s.dss == nil {
		if s.LastCheck == nil {
			return nil
		}
		t := *s.LastCheck
		return &t
	}

	s.dss.mu.RLock()
	defer s.dss.mu.RUnlock()

	if s.LastCheck == nil {
		return nil
	}

	t := *s.LastCheck
	return &t
}

func (s *TableSyncState) GetLastID() *string {
	if s == nil {
		return nil
	}
	if s.dss == nil {
		if s.LastID == nil {
			return nil
		}
		v := *s.LastID
		return &v
	}

	s.dss.mu.RLock()
	defer s.dss.mu.RUnlock()

	if s.LastID == nil {
		return nil
	}

	v := *s.LastID
	return &v
}
