package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fivenet-app/fivenet/v2026/cmd/envs"
	"github.com/stretchr/testify/require"
)

func TestLoadRejectsLongOAuth2ProviderName(t *testing.T) {
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "config.yaml")
	providerName := strings.Repeat("a", OAuth2ProviderNameMaxLen+1)

	err := os.WriteFile(cfgFile, []byte(`
oauth2:
  providers:
    - name: "`+providerName+`"
      label: "Example"
      homepage: "https://example.com"
      type: "generic"
      redirectURL: "https://example.com/api/oauth2/callback/example"
      clientID: "client-id"
      clientSecret: "client-secret"
      scopes:
        - openid
      endpoints:
        authURL: "https://example.com/oauth2/authorize"
        tokenURL: "https://example.com/oauth2/token"
        userInfoURL: "https://example.com/oauth2/userinfo"
`), 0o600)
	require.NoError(t, err)

	t.Setenv(envs.ConfigFileEnvVar, cfgFile)

	_, err = Load()
	require.Error(t, err)
	require.Contains(t, err.Error(), "Config.OAuth2.Providers[0].Name")
	require.Contains(t, err.Error(), "max")
}
