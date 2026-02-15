package discord

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/server/oauth2/providers"
	"github.com/fivenet-app/fivenet/v2026/pkg/server/oauth2/types"
	"golang.org/x/oauth2"
)

func init() {
	providers.RegisterProvider("discord", New)
}

type Discord struct {
	types.BaseProvider
}

type discordUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"profile_picture"`
}

func New(cfg *config.OAuth2Provider) types.IProvider {
	if cfg.Endpoints.AuthURL == "" {
		cfg.Endpoints.AuthURL = "https://discord.com/api/oauth2/authorize"
	}
	if cfg.Endpoints.TokenURL == "" {
		cfg.Endpoints.TokenURL = "https://discord.com/api/oauth2/token"
	}
	if cfg.Endpoints.UserInfoURL == "" {
		cfg.Endpoints.UserInfoURL = "https://discord.com/api/users/@me"
	}

	bp := types.NewBaseProvider(cfg)
	return &Discord{
		BaseProvider: bp,
	}
}

func (p *Discord) GetUserInfo(ctx context.Context, code string) (*types.UserInfo, error) {
	obtainedAt := time.Now()

	oauthCfg := p.GetOAuthConfig()
	token, err := oauthCfg.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	// Use the access token, here we use it to get the logged in user's info.
	//nolint:noctx // The context is already passed in to the oauth2 http client during creation.
	res, err := oauthCfg.Client(ctx, token).Get(p.UserInfoURL)
	if err != nil || res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info. %+w", err)
	}
	defer res.Body.Close()

	user, err := p.decodeUserInfo(oauthCfg.Scopes, token, res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode user info. %+w", err)
	}
	user.ObtainedAt = &obtainedAt

	return user, nil
}

func (p *Discord) decodeUserInfo(
	scopes []string,
	token *oauth2.Token,
	data io.Reader,
) (*types.UserInfo, error) {
	var dest discordUser
	if err := json.NewDecoder(data).Decode(&dest); err != nil {
		return nil, err
	}

	username := dest.Username
	if dest.Discriminator != "0" {
		username = dest.Username + "#" + dest.Discriminator
	}

	scopesString := strings.Join(scopes, " ")

	return &types.UserInfo{
		ID:       dest.ID,
		Username: username,
		Avatar: fmt.Sprintf(
			"https://cdn.discordapp.com/profile_pictures/%s/%s.png",
			dest.ID,
			dest.Avatar,
		),
		TokenType:    &token.TokenType,
		RefreshToken: &token.RefreshToken,
		AccessToken:  &token.AccessToken,
		Scope:        &scopesString,
		ExpiresIn:    &token.ExpiresIn,
	}, nil
}
