package api

type ClientConfig struct {
	Version string      `json:"version"`
	Login   LoginConfig `json:"login"`
	Discord *Discord    `json:"discord"`
	Links   Links       `json:"links"`
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

type Links struct {
	Imprint       *string `json:"imprint"`
	PrivacyPolicy *string `json:"privacyPolicy"`
}

type Version struct {
	Version string `json:"version"`
}
