package generic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
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
	subId, ok := anyToInt64(sub)
	if !ok {
		return nil, errors.New("failed to convert id to int64")
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

func anyToInt64(v any) (int64, bool) {
	switch x := v.(type) {
	case int:
		return int64(x), true
	case int8:
		return int64(x), true
	case int16:
		return int64(x), true
	case int32:
		return int64(x), true
	case int64:
		return x, true

	case uint:
		if uint64(x) > math.MaxInt64 {
			return 0, false
		}
		return int64(x), true
	case uint8:
		return int64(x), true
	case uint16:
		return int64(x), true
	case uint32:
		return int64(x), true
	case uint64:
		if x > math.MaxInt64 {
			return 0, false
		}
		return int64(x), true

	case float64:
		if x != math.Trunc(x) || x < math.MinInt64 || x > math.MaxInt64 {
			return 0, false
		}
		return int64(x), true

	case float32:
		f := float64(x)
		if f != math.Trunc(f) || f < math.MinInt64 || f > math.MaxInt64 {
			return 0, false
		}
		return int64(x), true

	case string:
		if x == "" {
			return 0, false
		}
		i, err := strconv.ParseInt(x, 10, 64)
		return i, err == nil

	case json.Number:
		i, err := x.Int64()
		return i, err == nil

	default:
		return 0, false
	}
}
