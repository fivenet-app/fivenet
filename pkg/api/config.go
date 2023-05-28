package api

type ClientConfig struct {
	SentryDSN string      `json:"sentryDSN"`
	Login     LoginConfig `json:"login"`
}

type LoginConfig struct {
	Providers []*ProviderConfig `json:"providers"`
}

type ProviderConfig struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}
