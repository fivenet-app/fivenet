package api

import (
	"net/url"
	"strconv"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/clientconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"go.uber.org/zap"
)

func (r *Routes) handleAppConfigUpdate(
	providers []*clientconfig.ProviderConfig,
	appCfg *appconfig.Cfg,
) {
	clientCfg := clientconfig.BuildClientConfig(r.cfg, providers, appCfg)
	r.clientCfg.Store(clientCfg)

	if appCfg.Discord.BotId == nil || appCfg.Discord.GetBotId() == "" {
		r.discordInviteUrl.Store(nil)
		return
	}

	clientId := appCfg.Discord.BotId
	permissions := strconv.FormatInt(appCfg.Discord.GetBotPermissions(), 10)
	redirectUri, err := url.JoinPath(r.cfg.HTTP.PublicURL, "/settings/props")
	if err != nil {
		r.logger.Warn("failed to build redirect URI for discord invite", zap.Error(err))
		return
	}
	redirectUri += "?tab=discord#"
	scopes := "bot identify"

	u, err := url.Parse("https://discord.com/oauth2/authorize")
	if err != nil {
		r.logger.Warn("failed to build discord invite URL", zap.Error(err))
		return
	}
	q := u.Query()
	q.Set("client_id", *clientId)
	q.Set("permissions", permissions)
	q.Set("scope", scopes)
	q.Set("redirect_uri", redirectUri)
	q.Set("response_type", "code")
	u.RawQuery = q.Encode()

	inviteUrl := u.String()
	r.discordInviteUrl.Store(&inviteUrl)
}
