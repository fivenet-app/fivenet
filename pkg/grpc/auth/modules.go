package auth

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"go.uber.org/fx"
)

var TokenMgrModule = fx.Module("tokenMgr",
	fx.Provide(
		NewTokenMgrFromConfig,
	),
)

func NewTokenMgrFromConfig(cfg *config.Config) *TokenMgr {
	return NewTokenMgr(cfg.JWT.Secret)
}

var AuthModule = fx.Module("grpc_auth",
	fx.Provide(
		NewGRPCAuth,
	),
)

var PermsModule = fx.Module("grpc_perms",
	fx.Provide(
		NewGRPCPerms,
	),
)
