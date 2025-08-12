package providers

import (
	"context"
	"errors"
	"fmt"
	"strconv"
)

type Generic struct {
	BaseProvider
}

func (p *Generic) GetUserInfo(ctx context.Context, code string) (*UserInfo, error) {
	token, err := p.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	res, err := p.oauthConfig.Client(ctx, token).Get(p.UserInfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %+w", err)
	}
	defer res.Body.Close()
	var dest map[string]any
	if err := json.NewDecoder(res.Body).Decode(&dest); err != nil {
		return nil, err
	}

	mapping := p.mapping
	sub, ok := dest[mapping.ID]
	if !ok {
		return nil, errors.New("failed to get id from user info")
	}
	subId := sub.(float64)
	if subId <= 0 {
		return nil, errors.New("invalid external user id given")
	}

	usernameRaw, ok := dest[mapping.Username]
	if !ok {
		return nil, errors.New("failed to get username from user info")
	}
	if usernameRaw == nil {
		return nil, errors.New("no username found in user info")
	}

	avatarRaw, ok := dest[mapping.Avatar]
	if !ok {
		return nil, errors.New("failed to get avatar from user info")
	}

	username := usernameRaw.(string)
	if avatarRaw == nil {
		avatarRaw = p.DefaultAvatar
	}
	avatar := avatarRaw.(string)

	user := &UserInfo{
		ID:       strconv.FormatInt(int64(subId), 10),
		Username: username,
		Avatar:   avatar,
	}

	return user, nil
}
