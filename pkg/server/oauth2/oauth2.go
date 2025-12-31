package oauth2

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/crypt"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/oauth2/providers"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

const (
	// AccountInfoRedirBase is the base path for account info redirects.
	AccountInfoRedirBase string = "/auth/account-info"
	// LoginRedirBase is the base path for login redirects.
	LoginRedirBase string = "/auth/login"

	// ReasonInternalError is used for generic internal error reasons.
	ReasonInternalError string = "internal_error"

	// Session keys used for OAuth2 state handling.
	SessionKeyOAuth2State = "fivenet_oauth2_state"
	SessionKeyConnectOnly = "connect-only"
	SessionKeyState       = "state"
	SessionKeyToken       = "token"
	SessionKeyRedirect    = "redirect"
)

var (
	tAccs   = table.FivenetAccounts
	tOauth2 = table.FivenetAccountsOauth2
)

// Params contains dependencies for constructing the OAuth2 handler.
type Params struct {
	fx.In

	// Logger is the zap logger instance for logging.
	Logger *zap.Logger
	// DB is the SQL database connection.
	DB *sql.DB
	// TM is the token manager for authentication tokens.
	TM *auth.TokenMgr
	// Config is the application configuration.
	Config *config.Config
	// Crypt is the cryptographic utility for secure operations.
	Crypt *crypt.Crypt
}

// OAuth2 handles OAuth2 authentication and user info management.
type OAuth2 struct {
	// logger is used for logging within the OAuth2 handler.
	logger *zap.Logger
	// db is the SQL database connection.
	db *sql.DB
	// tm is the token manager for authentication tokens.
	tm *auth.TokenMgr

	// domain is the cookie/session domain.
	domain string
	// oauthConfigs maps provider names to their configuration and logic.
	oauthConfigs map[string]providers.IProvider
	// userInfoStore handles user info storage and retrieval.
	userInfoStore userInfoStore
}

// New creates a new OAuth2 handler with all configured providers.
func New(p Params) *OAuth2 {
	if len(p.Config.OAuth2.Providers) == 0 {
		return nil
	}

	o := &OAuth2{
		logger:       p.Logger,
		db:           p.DB,
		tm:           p.TM,
		domain:       p.Config.HTTP.Sessions.Domain,
		oauthConfigs: make(map[string]providers.IProvider, len(p.Config.OAuth2.Providers)),
		userInfoStore: &oauth2UserInfo{
			db:    p.DB,
			crypt: p.Crypt,
		},
	}

	for _, p := range p.Config.OAuth2.Providers {
		cfg := &oauth2.Config{
			RedirectURL:  p.RedirectURL,
			ClientID:     p.ClientID,
			ClientSecret: p.ClientSecret,
			Scopes:       p.Scopes,
			Endpoint: oauth2.Endpoint{
				AuthURL:   p.Endpoints.AuthURL,
				TokenURL:  p.Endpoints.TokenURL,
				AuthStyle: oauth2.AuthStyleInParams,
			},
		}
		var provider providers.IProvider
		switch p.Type {
		case config.OAuth2ProviderDiscord:
			provider = &providers.Discord{
				BaseProvider: providers.BaseProvider{
					Name: p.Name,
				},
			}

		case config.OAuth2ProviderGeneric:
			fallthrough
		default:
			provider = &providers.Generic{
				BaseProvider: providers.BaseProvider{
					Name:          p.Name,
					UserInfoURL:   p.Endpoints.UserInfoURL,
					DefaultAvatar: p.DefaultAvatar,
				},
			}
		}

		provider.SetOauthConfig(cfg)
		provider.SetMapping(p.Mapping)

		o.oauthConfigs[p.Name] = provider
	}

	return o
}

// RegisterHTTP registers the OAuth2 login and callback endpoints on the given Gin engine.
func (o *OAuth2) RegisterHTTP(e *gin.Engine) {
	g := e.Group("/api/oauth2")
	{
		g.GET("/login/:provider", o.Login)
		g.POST("/login/:provider", o.Login)
		g.GET("/callback/:provider", o.Callback)
		g.POST("/callback/:provider", o.Callback)
	}
}

// isValidRedirectPath checks if the redirect path is allowed (basic check, can be extended).
func isValidRedirectPath(path string) bool {
	// Only allow internal redirects (no schema, no host, must start with /)
	return path != "" && path[0] == '/' && !containsDotDot(path)
}

func containsDotDot(path string) bool {
	return len(path) >= 2 && (path[:2] == ".." || path[len(path)-2:] == "..")
}

// handleRedirect redirects the user to the appropriate URL based on the outcome of the OAuth2 flow.
func (o *OAuth2) handleRedirect(c *gin.Context, connectOnly bool, success bool, reason string) {
	var redirURL string
	if !success {
		if connectOnly {
			redirURL = AccountInfoRedirBase + "?oauth2Connect=failed&tab=oauth2Connections"
		} else {
			redirURL = LoginRedirBase + "?oauth2Login=failed"
		}
		if reason != "" {
			redirURL = redirURL + "&reason=" + url.QueryEscape(reason)
		}
	} else {
		if connectOnly {
			redirURL = AccountInfoRedirBase + "?oauth2Connect=success&tab=oauth2Connections"
		} else {
			redirURL = LoginRedirBase + "?oauth2Login=success"
		}
	}
	c.Redirect(http.StatusTemporaryRedirect, redirURL+"#")
}

// Login initiates the OAuth2 login flow for the selected provider.
func (o *OAuth2) Login(c *gin.Context) {
	sess := sessions.DefaultMany(c, SessionKeyOAuth2State)
	connectOnly := false
	connectOnlyVal := c.Query(SessionKeyConnectOnly)
	if connectOnlyVal != "" {
		var err error
		connectOnly, err = strconv.ParseBool(connectOnlyVal)
		if err != nil {
			o.logger.Error("failed to parse connect only var", zap.Error(err))
			o.handleRedirect(c, false, false, "invalid_request_connect_only")
			return
		}
	}

	customRedirectVal := c.Query(SessionKeyRedirect)
	if customRedirectVal != "" {
		u, err := url.Parse(customRedirectVal)
		if err != nil || !isValidRedirectPath(u.Path) {
			o.logger.Error(
				"failed to parse or validate redirect url",
				zap.Error(err),
			)
			o.handleRedirect(c, false, false, "invalid_request_redirect")
			return
		}
		query := u.RawQuery
		up := u.Path
		if query != "" {
			up += "?" + query
		}
		if u.Fragment != "" {
			up += "#" + u.Fragment
		}
		sess.Set(SessionKeyRedirect, up)
	}

	tokenVal, err := c.Cookie("fivenet_token")
	if err != nil && connectOnly {
		o.logger.Error("failed to get token cookie for connect only request", zap.Error(err))
		o.handleRedirect(c, false, false, "invalid_request_token")
		return
	}
	if tokenVal != "" {
		sess.Set(SessionKeyToken, tokenVal)
	}

	state, err := utils.GenerateRandomString(64)
	if err != nil {
		o.logger.Error("failed to generate random string", zap.Error(err))
		o.handleRedirect(c, connectOnly, false, ReasonInternalError)
		return
	}

	sess.Set(SessionKeyConnectOnly, connectOnly)
	sess.Set(SessionKeyState, state)
	if err := sess.Save(); err != nil {
		o.logger.Error("failed to save session in Login", zap.Error(err))
		o.handleRedirect(c, connectOnly, false, ReasonInternalError)
		return
	}

	provider, err := o.GetProvider(c.Param("provider"))
	if err != nil {
		o.logger.Error(
			"failed to get provider",
			zap.Error(err),
		)
		o.handleRedirect(c, connectOnly, false, "invalid_provider")
		return
	}

	// Redirect to provider for login
	c.Redirect(http.StatusTemporaryRedirect, provider.GetRedirect(state))
}

// GetProvider retrieves the OAuth2 provider from the request context and name is case-insensitive.
func (o *OAuth2) GetProvider(providerName string) (providers.IProvider, error) {
	if providerName == "" {
		return nil, errors.New("no provider found in path")
	}

	for name, provider := range o.oauthConfigs {
		if name == providerName {
			return provider, nil
		}
		if name != "" && providerName != "" && name == providerName {
			return provider, nil
		}
	}

	return nil, fmt.Errorf("provider %q not configured", providerName)
}

// Callback handles the OAuth2 callback, processes user info, and issues tokens or connects accounts.
func (o *OAuth2) Callback(c *gin.Context) {
	sess := sessions.DefaultMany(c, SessionKeyOAuth2State)

	sessState, ok := sess.Get(SessionKeyState).(string)
	if !ok || sessState == "" {
		o.handleRedirect(c, false, false, "invalid_state")
		return
	}
	state := sessState

	sessConnectOnly, ok := sess.Get(SessionKeyConnectOnly).(bool)
	connectOnly := false
	if ok {
		connectOnly = sessConnectOnly
	}

	sessToken, ok := sess.Get(SessionKeyToken).(string)
	token := ""
	if ok {
		token = sessToken
	}

	sessRedirect, ok := sess.Get(SessionKeyRedirect).(string)
	redirect := ""
	if ok {
		redirect = sessRedirect
	}

	if c.Request.FormValue("state") != state {
		o.handleRedirect(c, connectOnly, false, "invalid_state_404")
		return
	}

	// Remove vars from session only after validation
	sess.Delete(SessionKeyConnectOnly)
	sess.Delete(SessionKeyState)
	sess.Delete(SessionKeyToken)
	sess.Delete(SessionKeyRedirect)
	if err := sess.Save(); err != nil {
		o.logger.Error("failed to save session in Callback", zap.Error(err))
		// Log error, but continue
	}

	provider, err := o.GetProvider(c.Param("provider"))
	if err != nil {
		o.logger.Error(
			"failed to get provider in callback",
			zap.String("provider", c.Param("provider")),
			zap.Error(err),
		)
		o.handleRedirect(c, connectOnly, false, "invalid_provider")
		return
	}

	ctx := c.Request.Context()

	userInfo, err := provider.GetUserInfo(ctx, c.Request.FormValue("code"))
	if err != nil {
		o.logger.Error(
			"failed to get userinfo from provider",
			zap.String("provider", provider.GetName()),
			zap.Error(err),
		)
		o.handleRedirect(c, connectOnly, false, "provider_failed")
		return
	}

	if connectOnly {
		o.handleConnectOnlyCallback(c, sess, token, provider, userInfo, redirect)
		return
	}
	o.handleLoginCallback(c, sess, provider, userInfo, connectOnly)
}

// handleConnectOnlyCallback processes the connect-only OAuth2 callback logic.
func (o *OAuth2) handleConnectOnlyCallback(
	c *gin.Context,
	_ sessions.Session,
	token string,
	provider providers.IProvider,
	userInfo *providers.UserInfo,
	redirect string,
) {
	claims, err := o.tm.ParseWithClaims(token)
	if err != nil {
		c.Redirect(
			http.StatusTemporaryRedirect,
			LoginRedirBase+"?oauth2Connect=failed&reason=token_invalid&tab=oauth2Connections#",
		)
		return
	}

	if err := o.userInfoStore.storeUserInfo(c.Request.Context(), claims.AccID, provider.GetName(), userInfo); err != nil {
		if dbutils.IsDuplicateError(err) {
			o.handleRedirect(c, true, false, "already_in_use")
		} else {
			o.logger.Error("failed to store user info", zap.Int64("acc_id", claims.AccID), zap.String("provider", provider.GetName()), zap.Error(err))
			o.handleRedirect(c, true, false, ReasonInternalError)
		}
		return
	}

	if redirect == "" {
		o.handleRedirect(c, true, true, "")
	} else {
		c.Redirect(http.StatusTemporaryRedirect, redirect)
	}
}

// handleLoginCallback processes the login OAuth2 callback logic.
func (o *OAuth2) handleLoginCallback(
	c *gin.Context,
	_ sessions.Session,
	provider providers.IProvider,
	userInfo interface{},
	connectOnly bool,
) {
	uInfo, ok := userInfo.(*providers.UserInfo)
	if !ok {
		o.logger.Error(
			"userInfo type assertion failed in handleLoginCallback",
			zap.String("provider", provider.GetName()),
		)
		o.handleRedirect(c, connectOnly, false, ReasonInternalError)
		return
	}

	account, err := o.userInfoStore.getAccountInfo(c.Request.Context(), provider.GetName(), uInfo)
	if err != nil {
		o.logger.Error(
			"failed to get/store userinfo in database",
			zap.String("provider", provider.GetName()),
			zap.Error(err),
		)
		o.handleRedirect(c, connectOnly, false, ReasonInternalError)
		return
	}

	if account == nil {
		c.Redirect(
			http.StatusTemporaryRedirect,
			LoginRedirBase+"?oauth2Login=failed&reason=unconnected",
		)
		return
	} else if account.ID == 0 {
		o.logger.Error("invalid account id from userinfo", zap.String("provider", provider.GetName()), zap.Error(err))
		o.handleRedirect(c, connectOnly, true, ReasonInternalError)
		return
	}

	if err := o.userInfoStore.updateUserInfo(c.Request.Context(), account.ID, provider.GetName(), uInfo); err != nil {
		o.logger.Error(
			"failed to update oauth2 user info for account id",
			zap.Int64("account_id", account.ID),
			zap.String("provider", provider.GetName()),
			zap.Error(err),
		)
		o.handleRedirect(c, connectOnly, true, ReasonInternalError)
		return
	}

	claims := auth.BuildTokenClaimsFromAccount(account, nil)
	newToken, err := o.tm.NewWithClaims(claims)
	if err != nil {
		o.logger.Error(
			"failed to create token from account",
			zap.String("provider", provider.GetName()),
			zap.Error(err),
		)
		o.handleRedirect(c, connectOnly, true, ReasonInternalError)
		return
	}

	c.SetCookie(auth.TokenCookieName, newToken, 6*24*60*60, "/", o.domain, true, true)

	c.Redirect(
		http.StatusTemporaryRedirect,
		fmt.Sprintf(LoginRedirBase+"?oauth2Login=success&u=%s&exp=%d",
			url.QueryEscape(*account.Username),
			claims.ExpiresAt.Time.UTC().UnixNano()/1e6,
		),
	)
}
