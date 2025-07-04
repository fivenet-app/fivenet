package providers

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Discord struct {
	BaseProvider
}

type discordUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
}

func (p *Discord) GetUserInfo(ctx context.Context, code string) (*UserInfo, error) {
	obtainedAt := time.Now()
	token, err := p.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	// Use the access token, here we use it to get the logged in user's info.
	res, err := p.oauthConfig.Client(ctx, token).Get("https://discord.com/api/users/@me")
	if err != nil || res.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get user info: %+q", err)
	}
	defer res.Body.Close()

	var dest discordUser
	if err := json.NewDecoder(res.Body).Decode(&dest); err != nil {
		return nil, err
	}

	username := dest.Username
	if dest.Discriminator != "0" {
		username = dest.Username + "#" + dest.Discriminator
	}

	scopes := strings.Join(p.oauthConfig.Scopes, " ")

	return &UserInfo{
		ID:           dest.ID,
		Username:     username,
		Avatar:       fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", dest.ID, dest.Avatar),
		TokenType:    &token.TokenType,
		RefreshToken: &token.RefreshToken,
		AccessToken:  &token.AccessToken,
		Scope:        &scopes,
		ExpiresIn:    &token.ExpiresIn,
		ObtainedAt:   &obtainedAt,
	}, nil
}
