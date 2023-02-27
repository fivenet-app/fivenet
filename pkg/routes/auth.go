package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/galexrt/arpanet/model"
	"github.com/galexrt/arpanet/pkg/config"
	"github.com/galexrt/arpanet/pkg/session"
	"github.com/galexrt/arpanet/query"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	jwtSigningKey []byte
}

func NewAuth() *Auth {
	return &Auth{
		jwtSigningKey: []byte(config.C.JWT.Secret),
	}
}

func (r *Auth) Register(e *gin.Engine) {
	g := e.Group("/auth")
	{
		g.GET("", r.IndexGET)
		g.POST("", r.IndexGET)
		g.POST("/login", r.LoginPOST)
		g.POST("/logout", r.LogoutPOST)
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

func (r *Auth) createTokenForAccount(account *model.Account) (string, error) {
	claims := &session.UserInfoClaims{
		License: account.License,
		RegisteredClaims: jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "arpanet",
			Subject:   "somebody",
			ID:        strconv.FormatUint(uint64(account.ID), 10),
			Audience:  []string{"somebody_else"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(r.jwtSigningKey)
	if err != nil {
		return "", err
	}
	return ss, nil
}

func (r *Auth) getAccountFromDB(userID string) (*model.Account, error) {
	accounts := query.Accounts
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

type AuthLoginPOSTResponse struct {
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

	token, err := r.createTokenForAccount(account)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &AuthLoginPOSTResponse{
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
