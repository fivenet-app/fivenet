package providers

import (
	"fmt"

	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/server/oauth2/types"
)

type ProviderFactory func(cfg *config.OAuth2Provider) types.IProvider

var providerFactories = map[string]ProviderFactory{}

func GetProvider(cfg *config.OAuth2Provider) (types.IProvider, error) {
	fn, ok := providerFactories[string(cfg.Type)]
	if !ok {
		return nil, fmt.Errorf("invalid oauth2 provider %q type given", cfg.Type)
	}

	return fn(cfg), nil
}

func RegisterProvider(name string, factory ProviderFactory) {
	providerFactories[name] = factory
}
