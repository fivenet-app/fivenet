package types

import (
	"context"
	"time"

	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"golang.org/x/oauth2"
)

type IProvider interface {
	GetName() string

	GetRedirect(state string) (string, error)
	GetUserInfo(ctx context.Context, code string) (*UserInfo, error)
}

type UserInfo struct {
	ID       string
	Username string
	Avatar   string

	RefreshToken *string
	AccessToken  *string
	Scope        *string
	TokenType    *string
	ExpiresIn    *int64
	ObtainedAt   *time.Time
}

type BaseProvider struct {
	name string

	oauthConfig *oauth2.Config
	mapping     *config.OAuth2Mapping

	DefaultAvatar string
	UserInfoURL   string
}

func NewBaseProvider(cfg *config.OAuth2Provider) BaseProvider {
	return BaseProvider{
		name:          cfg.Name,
		oauthConfig:   cfg.GetOAuth2Config(),
		mapping:       cfg.Mapping,
		DefaultAvatar: cfg.DefaultAvatar,
		UserInfoURL:   cfg.Endpoints.UserInfoURL,
	}
}

func (b *BaseProvider) GetName() string {
	return b.name
}

func (b *BaseProvider) SetName(name string) {
	b.name = name
}

func (b *BaseProvider) GetOAuthConfig() *oauth2.Config {
	return b.oauthConfig
}

func (b *BaseProvider) SetOAuthConfig(cfg *oauth2.Config) {
	b.oauthConfig = cfg
}

func (b *BaseProvider) GetMapping() *config.OAuth2Mapping {
	return b.mapping
}

func (b *BaseProvider) SetMapping(mapping *config.OAuth2Mapping) {
	b.mapping = mapping
}

func (b *BaseProvider) GetRedirect(state string) (string, error) {
	return b.GetOAuthConfig().AuthCodeURL(state), nil
}
