package auth

import (
	"go.uber.org/fx"
)

var Module = fx.Module("proto-services-auth",
	fx.Provide(
		NewServer,
	),
)
