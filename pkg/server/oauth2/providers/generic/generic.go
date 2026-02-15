package generic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/server/oauth2/providers"
	"github.com/fivenet-app/fivenet/v2026/pkg/server/oauth2/types"
)

func init() {
	providers.RegisterProvider("generic", New)
}

type Generic struct {
	types.BaseProvider
}

func New(cfg *config.OAuth2Provider) types.IProvider {
	bp := types.NewBaseProvider(cfg)
	return &Generic{
		BaseProvider: bp,
	}
}

func (p *Generic) GetUserInfo(ctx context.Context, code string) (*types.UserInfo, error) {
	oauthCfg := p.GetOAuthConfig()
	token, err := oauthCfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed. %w", err)
	}

	//nolint:noctx // The context is already passed in to the oauth2 http client during creation.
	res, err := oauthCfg.Client(ctx, token).Get(p.UserInfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info. %+w", err)
	}
	defer res.Body.Close()

	user, err := p.decodeUserInfo(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode user info. %w", err)
	}

	return user, nil
}

func (p *Generic) decodeUserInfo(data io.Reader) (*types.UserInfo, error) {
	var dest map[string]any
	if err := json.NewDecoder(data).Decode(&dest); err != nil {
		return nil, err
	}

	mapping := p.GetMapping()
	// External ID
	sub, ok := dest[mapping.ID]
	if !ok {
		return nil, errors.New("failed to get id from user info")
	}
	subId, ok := sub.(float64)
	if !ok {
		return nil, errors.New("failed to convert id to float64")
	}
	if subId <= 0 {
		return nil, errors.New("invalid external user id given")
	}

	// Username
	usernameRaw, ok := dest[mapping.Username]
	if !ok {
		return nil, errors.New("failed to get username from user info")
	}
	if usernameRaw == nil {
		return nil, errors.New("no username found in user info")
	}
	username, ok := usernameRaw.(string)
	if !ok {
		return nil, errors.New("failed to get username from user info")
	}

	// Profile Picture
	profilePictureRaw, ok := dest[mapping.Avatar]
	if !ok {
		return nil, errors.New("failed to get profile_picture from user info")
	}
	if profilePictureRaw == nil {
		profilePictureRaw = p.DefaultAvatar
	}
	profilePicture, ok := profilePictureRaw.(string)
	if !ok {
		return nil, errors.New("failed to get profile_picture from user info")
	}

	return &types.UserInfo{
		ID:       strconv.FormatInt(int64(subId), 10),
		Username: username,
		Avatar:   profilePicture,
	}, nil
}
