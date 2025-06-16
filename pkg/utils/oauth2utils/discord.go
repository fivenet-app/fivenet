package oauth2utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// RefreshDiscordAccessToken refreshes a Discord OAuth2 access token using the provided credentials and refresh token.
// Returns the new access token, new refresh token, expiration (in seconds), or an error if the refresh fails.
func RefreshDiscordAccessToken(ctx context.Context, clientID, clientSecret, refreshToken, redirectURI string) (newAccessToken, newRefreshToken string, expiresIn int32, err error) {
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("redirect_uri", redirectURI)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://discord.com/api/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", "", 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", 0, fmt.Errorf("discord token refresh error %d: %s", resp.StatusCode, string(body))
	}

	// Inline struct for decoding Discord's token response
	var respData struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int32  `json:"expires_in"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", "", 0, err
	}

	return respData.AccessToken, respData.RefreshToken, int32(respData.ExpiresIn), nil
}
