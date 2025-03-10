package api

import "github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"

type Version struct {
	Version string `json:"version"`
}

type ClientConfig struct {
	Version string `json:"version"`

	DefaultLocale string `json:"defaultLocale"`

	Login        LoginConfig  `json:"login"`
	Discord      Discord      `json:"discord"`
	Website      Website      `json:"website"`
	FeatureGates FeatureGates `json:"featureGates"`
	Game         Game         `json:"game"`
	System       System       `json:"system"`
}

type LoginConfig struct {
	SignupEnabled bool              `json:"signupEnabled"`
	LastCharLock  bool              `json:"lastCharLock"`
	Providers     []*ProviderConfig `json:"providers"`
}

type ProviderConfig struct {
	Name  string  `json:"name"`
	Label string  `json:"label"`
	Icon  *string `json:"icon"`
}

type Discord struct {
	BotInviteURL *string `json:"botInviteURL"`
}

type Website struct {
	Links     Links `json:"links"`
	StatsPage bool  `json:"statsPage"`
}

type Links struct {
	Imprint       *string `json:"imprint"`
	PrivacyPolicy *string `json:"privacyPolicy"`
}

type FeatureGates struct {
	ImageProxy bool `json:"imageProxy"`
}

type Game struct {
	UnemployedJobName string `json:"unemployedJobName"`
	StartJobGrade     int32  `json:"startJobGrade"`
}

type System struct {
	BannerMessage *rector.BannerMessage `json:"bannerMessage,omitempty"`
}
