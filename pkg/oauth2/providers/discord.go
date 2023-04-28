package providers

import (
	"context"
	"fmt"
	"strconv"
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

func (p *Discord) GetUserInfo(code string) (*UserInfo, error) {
	token, err := p.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	// Use the access token, here we use it to get the logged in user's info.
	res, err := p.oauthConfig.Client(context.Background(), token).Get("https://discord.com/api/users/@me")
	if err != nil || res.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get user info: %+q", err)
	}
	defer res.Body.Close()

	var dest discordUser
	if err := json.NewDecoder(res.Body).Decode(&dest); err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(dest.ID)
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		ID:       int64(id),
		Username: dest.Username + "#" + dest.Discriminator,
		Avatar:   fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/user_avatar.png", dest.Avatar),
	}, nil
}
