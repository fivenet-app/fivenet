package auth

import (
	"github.com/galexrt/fivenet/pkg/config"
	"go.uber.org/fx"
)

var TokenMgrModule = fx.Module("tokenMgr",
	fx.Provide(
		NewTokenMgrFromConfig,
	),
)

func NewTokenMgrFromConfig(cfg *config.BaseConfig) *TokenMgr {
	return NewTokenMgr(cfg.JWT.Secret)
}

var AuthModule = fx.Module("auth",
	fx.Provide(
		NewGRPCAuth,
	),
)
