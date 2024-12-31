package dbsync

import (
	"context"
	"database/sql"

	"go.uber.org/fx"
)

var Module = fx.Module("dbsync",
	fx.Provide(
		New,
	),
)

type Sync struct {
	db *sql.DB
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	DB *sql.DB
}

func New(p Params) (*Sync, error) {
	s := &Sync{
		db: p.DB,
	}

	ctx, cancel := context.WithCancel(context.Background())

	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		go s.Run(ctx)

		return nil
	}))

	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
		cancel()

		return nil
	}))

	return s, nil
}

func (s *Sync) Run(ctx context.Context) {
	// TODO
}
