package userinfo

import "go.uber.org/fx"

var Module = fx.Module("userinfo",
	fx.Provide(
		NewUIRetriever,
	),
)
