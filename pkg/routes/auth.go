package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/galexrt/arpanet/model"
	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/session"
	"github.com/galexrt/arpanet/query"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
}

func (r *Auth) Register(e *gin.Engine) {
	g := e.Group("/auth")
	{
		g.GET("", r.IndexGET)
		g.POST("", r.IndexGET)
		g.POST("/login", r.LoginPOST)
		g.POST("/logout", r.LogoutPOST)
		g.POST("/token", r.TokenPOST)
	}
}

// Base auth page to show authentication status
func (r *Auth) IndexGET(c *gin.Context) {
	s := sessions.DefaultMany(c, session.UserSession)

	userLoggedIn := s.Get(session.UserIDKey)
	if userLoggedIn == nil || userLoggedIn.(string) == "" {
		c.JSON(http.StatusForbidden, "You are not authenticated!")
		return
	}

	c.JSON(http.StatusOK, "You are already authenticated!")
}

func (r *Auth) createTokenForAccount(account *model.Account, charIndex int) (string, error) {
	return session.Tokens.NewWithClaims(&session.UserInfoClaims{
		AccountID: account.ID,
		CharIndex: charIndex,
		RegisteredClaims: jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "arpanet",
			Subject:   account.License,
			ID:        strconv.FormatUint(uint64(account.ID), 10),
			Audience:  []string{"arpanet"},
		},
	})
}

func (r *Auth) getAccountFromDB(userID string) (*model.Account, error) {
	accounts := query.Account
	account, err := accounts.Where(accounts.Enabled.Is(true), accounts.Username.Eq(userID)).Limit(1).First()
	if err != nil {
		return nil, err

	}

	return account, nil
}

type AuthLoginPOSTForm struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

// User login
func (r *Auth) LoginPOST(c *gin.Context) {
	var form AuthLoginPOSTForm
	if err := c.ShouldBind(&form); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	account, err := r.getAccountFromDB(form.Username)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	s := sessions.DefaultMany(c, session.UserSession)
	s.Set(session.UserIDKey, account.Username)
	if err := s.Save(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	token, err := r.createTokenForAccount(account, 1)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &TokenResponse{
		Token: token,
	})
}

// User logout
func (r *Auth) LogoutPOST(c *gin.Context) {
	// TODO handle JWT token "invalidation" in the future
	s := sessions.DefaultMany(c, session.UserSession)
	s.Clear()
	if err := s.Save(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, "Your are logged out!")
}

type TokenPOSTRequest struct {
	Token     string `json:"token"`
	CharIndex int    `json:"char_index"`
}

// Basically a way to "switch" characters by updating the claim
func (r *Auth) TokenPOST(c *gin.Context) {
	var form TokenPOSTRequest
	if err := c.ShouldBind(&form); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	claims, err := session.Tokens.ParseWithClaims(form.Token)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Find user info for the new/old char index in the claim
	users := query.User
	if _, err := users.Where(users.Identifier.Like(auth.BuildIdentifierFromLicense(form.CharIndex, claims.Subject))).Limit(1).First(); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Update claims to the new char index
	claims.CharIndex = form.CharIndex

	token, err := session.Tokens.NewWithClaims(claims)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, &TokenResponse{
		Token: token,
	})
}
