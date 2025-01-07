package dbsync

import (
	"context"
)

type UserLicensesSync struct {
	*syncer

	state *TableSyncState
}

func NewUserLicensesSync(s *syncer, state *TableSyncState) (ISyncer, error) {
	return &UserLicensesSync{
		syncer: s,
		state:  state,
	}, nil
}

func (s *UserLicensesSync) Sync(ctx context.Context) error {
	if !s.cfg.Tables.UserLicenses.Enabled {
		return nil
	}

	// TODO

	return nil
}
