package perms

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("perms",
	fx.Provide(
		New,
	),
	fx.Decorate(wrapLogger),
)

func wrapLogger(log *zap.Logger) *zap.Logger {
	return log.Named("perms")
}

var TestModule = fx.Module("perms",
	fx.Provide(func(p Params) (Permissions, error) {
		ps, err := New(p)
		if err != nil {
			return nil, err
		}

		// Enable dev mode
		ps.(*Perms).devMode = true

		return ps, nil
	}),
	fx.Decorate(wrapLogger),
)
