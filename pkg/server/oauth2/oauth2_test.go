package oauth2

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/oauth2/providers"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

type MockProvider struct {
	mock.Mock
	providers.BaseProvider
}

func (m *MockProvider) GetRedirect(state string) string {
	args := m.Called(state)
	return args.String(0)
}

func (m *MockProvider) GetUserInfo(ctx context.Context, code string) (*providers.UserInfo, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*providers.UserInfo), args.Error(1)
}

func createMockProvider() *MockProvider {
	mockProvider := new(MockProvider)
	mockProvider.Name = "test-provider"
	return mockProvider
}

type MockUserInfoStore struct {
	mock.Mock
}

func (m *MockUserInfoStore) storeUserInfo(
	ctx context.Context,
	accountId int64,
	provider string,
	userInfo *providers.UserInfo,
) error {
	args := m.Called(ctx, accountId, provider, userInfo)
	return args.Error(0)
}

func (m *MockUserInfoStore) updateUserInfo(
	ctx context.Context,
	accountId int64,
	provider string,
	userInfo *providers.UserInfo,
) error {
	args := m.Called(ctx, accountId, provider, userInfo)
	return args.Error(0)
}

func (m *MockUserInfoStore) getAccountInfo(
	ctx context.Context,
	provider string,
	userInfo *providers.UserInfo,
) (*model.FivenetAccounts, error) {
	args := m.Called(ctx, provider, userInfo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.FivenetAccounts), args.Error(1)
}

func TestCallback_InvalidState(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	store := cookie.NewStore([]byte("secret"))
	sess := sessions.SessionsMany([]string{"fivenet_oauth2_state"}, store)
	router.Use(sess)

	oauth := &OAuth2{
		logger: zaptest.NewLogger(t),
	}
	router.GET("/callback/:provider", oauth.Callback)

	req := httptest.NewRequest(http.MethodGet, "/callback/test-provider?state=invalid", nil)
	c.Request = req

	sess(c)
	session := sessions.DefaultMany(c, "fivenet_oauth2_state")
	session.Set("state", "valid")
	session.Save()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Contains(t, w.Header().Get("Location"), "reason=invalid_state_404")
}

func TestCallback_InvalidProvider(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	store := cookie.NewStore([]byte("secret"))
	sess := sessions.SessionsMany([]string{"fivenet_oauth2_state"}, store)
	router.Use(sess)

	oauth := &OAuth2{
		logger: zaptest.NewLogger(t),
	}
	router.GET("/callback/:provider", oauth.Callback)

	req := httptest.NewRequest(http.MethodGet, "/callback/invalid-provider?state=valid", nil)
	c.Request = req

	sess(c)
	session := sessions.DefaultMany(c, "fivenet_oauth2_state")
	session.Set("state", "valid")
	session.Save()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Contains(t, w.Header().Get("Location"), "reason=invalid_provider")
}

func TestCallback_ProviderError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	store := cookie.NewStore([]byte("secret"))
	sess := sessions.SessionsMany([]string{"fivenet_oauth2_state"}, store)
	router.Use(sess)

	mockProvider := createMockProvider()
	mockProvider.On("GetUserInfo", mock.Anything, "test-code").
		Return(nil, errors.New("provider error"))

	oauth := &OAuth2{
		logger: zaptest.NewLogger(t),
		oauthConfigs: map[string]providers.IProvider{
			"test-provider": mockProvider,
		},
	}
	router.GET("/callback/:provider", oauth.Callback)

	req := httptest.NewRequest(
		http.MethodGet,
		"/callback/test-provider?state=valid&code=test-code",
		nil,
	)
	c.Request = req

	sess(c)
	session := sessions.DefaultMany(c, "fivenet_oauth2_state")
	session.Set("state", "valid")
	session.Save()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Contains(t, w.Header().Get("Location"), "reason=provider_failed")
	mockProvider.AssertExpectations(t)
}

func TestCallback_LoginProviderError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	store := cookie.NewStore([]byte("secret"))
	sess := sessions.SessionsMany([]string{"fivenet_oauth2_state"}, store)
	router.Use(sess)

	oauth := &OAuth2{
		logger: zaptest.NewLogger(t),
	}
	router.GET("/login/:provider", oauth.Callback)

	req := httptest.NewRequest(http.MethodGet, "/login/invalid-provider?state=valid", nil)
	c.Request = req

	sess(c)
	session := sessions.DefaultMany(c, "fivenet_oauth2_state")
	session.Set("state", "valid")
	session.Save()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Contains(t, w.Header().Get("Location"), "reason=invalid_provider")
}

func TestCallback_LoginSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	store := cookie.NewStore([]byte("secret"))
	sess := sessions.SessionsMany([]string{"fivenet_oauth2_state"}, store)
	router.Use(sess)

	mockProvider := createMockProvider()
	mockUserInfo := &providers.UserInfo{
		ID:       "123",
		Username: "testuser",
		Avatar:   "profile_picture.png",
	}
	mockProvider.On("GetUserInfo", mock.Anything, "test-code").Return(mockUserInfo, nil)

	mockUserInfoStore := &MockUserInfoStore{}
	mockUserInfoStore.On("getAccountInfo", mock.Anything, "test-provider", mockUserInfo).
		Return(&model.FivenetAccounts{
			ID:       123,
			Username: &mockUserInfo.Username,
			License:  "license",
		}, nil)
	mockUserInfoStore.On("updateUserInfo", mock.Anything, int64(123), "test-provider", mockUserInfo).
		Return(nil)

	oauth := &OAuth2{
		logger: zaptest.NewLogger(t),
		oauthConfigs: map[string]providers.IProvider{
			"test-provider": mockProvider,
		},
		userInfoStore: mockUserInfoStore,
		tm:            auth.NewTokenMgr("secret"),
	}
	router.GET("/callback/:provider", oauth.Callback)

	req := httptest.NewRequest(
		http.MethodGet,
		"/callback/test-provider?state=valid&code=test-code",
		nil,
	)
	c.Request = req

	sess(c)
	session := sessions.DefaultMany(c, "fivenet_oauth2_state")
	session.Set("state", "valid")
	session.Save()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Contains(
		t,
		w.Header().Get("Location"),
		"/auth/login?oauth2Login=success&u=testuser&exp=",
	)
	mockProvider.AssertExpectations(t)
}

func TestCallback_ConnectError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	store := cookie.NewStore([]byte("secret"))
	sess := sessions.SessionsMany([]string{"fivenet_oauth2_state"}, store)
	router.Use(sess)

	tm := auth.NewTokenMgr("secret")

	mockProvider := createMockProvider()
	mockUserInfo := &providers.UserInfo{
		ID:       "123",
		Username: "testuser",
		Avatar:   "profile_picture.png",
	}
	mockProvider.On("GetUserInfo", mock.Anything, "test-code").Return(mockUserInfo, nil)

	mockUserInfoStore := &MockUserInfoStore{}
	mockUserInfoStore.On("getAccountInfo", mock.Anything, "test-provider", mockUserInfo).
		Return(&model.FivenetAccounts{
			ID:       123,
			Username: &mockUserInfo.Username,
			License:  "license",
		}, nil)

	oauth := &OAuth2{
		logger: zaptest.NewLogger(t),
		oauthConfigs: map[string]providers.IProvider{
			"test-provider": mockProvider,
		},
		userInfoStore: mockUserInfoStore,
		tm:            tm,
	}
	router.GET("/callback/:provider", oauth.Callback)

	req := httptest.NewRequest(
		http.MethodGet,
		"/callback/test-provider?state=valid&code=test-code",
		nil,
	)
	c.Request = req

	sess(c)
	session := sessions.DefaultMany(c, "fivenet_oauth2_state")
	session.Set("state", "valid")
	session.Set("connect-only", true)
	session.Set("token", "invalid")
	session.Save()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Contains(t, w.Header().Get("Location"), "oauth2Connect=failed&reason=token_invalid")
}

func TestCallback_ConnectErrorAlreadyInUse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	store := cookie.NewStore([]byte("secret"))
	sess := sessions.SessionsMany([]string{"fivenet_oauth2_state"}, store)
	router.Use(sess)

	tm := auth.NewTokenMgr("secret")
	token, err := tm.NewWithClaims(&auth.CitizenInfoClaims{
		AccID:    123,
		Username: "testuser",
	})
	require.NoError(t, err)

	mockProvider := createMockProvider()
	mockUserInfo := &providers.UserInfo{
		ID:       "123",
		Username: "testuser",
		Avatar:   "profile_picture.png",
	}
	mockProvider.On("GetUserInfo", mock.Anything, "test-code").Return(mockUserInfo, nil)

	mockUserInfoStore := &MockUserInfoStore{}
	mockUserInfoStore.On("storeUserInfo", mock.Anything, int64(123), "test-provider", mockUserInfo).
		Return(&mysql.MySQLError{Number: 1062})

	oauth := &OAuth2{
		logger: zaptest.NewLogger(t),
		oauthConfigs: map[string]providers.IProvider{
			"test-provider": mockProvider,
		},
		userInfoStore: mockUserInfoStore,
		tm:            tm,
	}
	router.GET("/callback/:provider", oauth.Callback)

	req := httptest.NewRequest(
		http.MethodGet,
		"/callback/test-provider?state=valid&code=test-code",
		nil,
	)
	c.Request = req

	sess(c)
	session := sessions.DefaultMany(c, "fivenet_oauth2_state")
	session.Set("state", "valid")
	session.Set("connect-only", true)
	session.Set("token", token)
	session.Save()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Contains(t, w.Header().Get("Location"), "reason=already_in_use")
}

func TestCallback_ConnectFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	store := cookie.NewStore([]byte("secret"))
	sess := sessions.SessionsMany([]string{"fivenet_oauth2_state"}, store)
	router.Use(sess)

	tm := auth.NewTokenMgr("secret")
	token, err := tm.NewWithClaims(&auth.CitizenInfoClaims{
		AccID:    123,
		Username: "testuser",
	})
	require.NoError(t, err)

	mockProvider := createMockProvider()
	mockUserInfo := &providers.UserInfo{
		ID:       "123",
		Username: "testuser",
		Avatar:   "profile_picture.png",
	}
	redirectUrl := "https://example.com/redirect-url?state="
	mockProvider.On("GetRedirect", mock.Anything).Return(redirectUrl)
	mockProvider.On("GetUserInfo", mock.Anything, "test-code").Return(mockUserInfo, nil)

	mockUserInfoStore := &MockUserInfoStore{}
	mockUserInfoStore.On("getAccountInfo", mock.Anything, "test-provider", mockUserInfo).
		Return(&model.FivenetAccounts{
			ID:       123,
			Username: &mockUserInfo.Username,
			License:  "license",
		}, nil)
	mockUserInfoStore.On("storeUserInfo", mock.Anything, int64(123), "test-provider", mockUserInfo).
		Return(nil)
	oauth := &OAuth2{
		logger: zaptest.NewLogger(t),
		oauthConfigs: map[string]providers.IProvider{
			"test-provider": mockProvider,
		},
		userInfoStore: mockUserInfoStore,
		tm:            tm,
	}
	router.GET("/login/:provider", oauth.Login)
	router.GET("/callback/:provider", oauth.Callback)

	req := httptest.NewRequest(http.MethodGet, "/login/test-provider?connect-only=true", nil)
	req.AddCookie(&http.Cookie{
		Name:  "fivenet_token",
		Value: token,
	})
	c.Request = req
	sess(c)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Contains(t, w.Header().Get("Location"), redirectUrl)

	session := sessions.DefaultMany(c, "fivenet_oauth2_state")
	state := session.Get("state")
	assert.NotEmpty(t, state)
	require.NoError(t, session.Save())

	req.AddCookie(&http.Cookie{
		Name:  "fivenet_token",
		Value: token,
	})

	stateVal := state.(string)
	req.URL, err = url.Parse("/callback/test-provider?state=" + stateVal + "&code=test-code")
	require.NoError(t, err)
	c.Request = req
	sess(c)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Equal(
		t,
		"/auth/account-info?oauth2Connect=success&tab=oauth2Connections#",
		w.Header().Get("Location"),
	)
	mockProvider.AssertExpectations(t)
}
