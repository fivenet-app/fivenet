package api

type ClientConfig struct {
	Version   string      `json:"version"`
	SentryDSN string      `json:"sentryDSN"`
	Login     LoginConfig `json:"login"`
	Discord   *Discord    `json:"discord"`
}

type LoginConfig struct {
	SignupEnabled bool              `json:"signupEnabled"`
	Providers     []*ProviderConfig `json:"providers"`
}

type ProviderConfig struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

type Discord struct {
	BotInviteURL string `json:"botInviteURL"`
}

type Version struct {
	Version string `json:"version"`
}
