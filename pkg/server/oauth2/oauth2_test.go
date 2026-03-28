package oauth2

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	authclaims "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth/claims"
	"github.com/fivenet-app/fivenet/v2026/pkg/server/oauth2/types"
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
	types.BaseProvider
	mock.Mock
}

func (m *MockProvider) GetRedirect(state string) (string, error) {
	args := m.Called(state)
	return args.String(0), args.Error(1)
}

func (m *MockProvider) GetUserInfo(ctx context.Context, code string) (*types.UserInfo, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.UserInfo), args.Error(1)
}

func createMockProvider() *MockProvider {
	mockProvider := new(MockProvider)
	mockProvider.SetName("test-provider")
	return mockProvider
}

type MockUserInfoStore struct {
	mock.Mock
}

func (m *MockUserInfoStore) storeUserInfo(
	ctx context.Context,
	accountId int64,
	provider string,
	userInfo *types.UserInfo,
) error {
	args := m.Called(ctx, accountId, provider, userInfo)
	return args.Error(0)
}

func (m *MockUserInfoStore) updateUserInfo(
	ctx context.Context,
	accountId int64,
	provider string,
	userInfo *types.UserInfo,
) error {
	args := m.Called(ctx, accountId, provider, userInfo)
	return args.Error(0)
}

func (m *MockUserInfoStore) getAccountInfo(
	ctx context.Context,
	provider string,
	userInfo *types.UserInfo,
) (*accounts.Account, error) {
	args := m.Called(ctx, provider, userInfo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*accounts.Account), args.Error(1)
}

func TestCallback_InvalidState(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	store := cookie.NewStore([]byte("secret"))
	sess := sessions.SessionsMany([]string{"fivenet_oauth2_state"}, store)
	router.Use(sess)

	oauth := &OAuth2{
		logger: zaptest.NewLogger(t),
	}
	oauth.RegisterHTTP(router)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/oauth2/callback/test-provider?state=invalid",
		nil,
	)
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
	t.Parallel()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	store := cookie.NewStore([]byte("secret"))
	sess := sessions.SessionsMany([]string{"fivenet_oauth2_state"}, store)
	router.Use(sess)

	oauth := &OAuth2{
		logger: zaptest.NewLogger(t),
	}
	oauth.RegisterHTTP(router)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/oauth2/callback/invalid-provider?state=valid",
		nil,
	)
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
	t.Parallel()
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
		oauthConfigs: map[string]types.IProvider{
			"test-provider": mockProvider,
		},
	}
	oauth.RegisterHTTP(router)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/oauth2/callback/test-provider?state=valid&code=test-code",
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
	t.Parallel()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	store := cookie.NewStore([]byte("secret"))
	sess := sessions.SessionsMany([]string{"fivenet_oauth2_state"}, store)
	router.Use(sess)

	oauth := &OAuth2{
		logger: zaptest.NewLogger(t),
	}
	oauth.RegisterHTTP(router)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/oauth2/login/invalid-provider?state=valid",
		nil,
	)
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
	t.Parallel()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	store := cookie.NewStore([]byte("secret"))
	sess := sessions.SessionsMany([]string{"fivenet_oauth2_state"}, store)
	router.Use(sess)

	mockProvider := createMockProvider()
	mockUserInfo := &types.UserInfo{
		ID:       "123",
		Username: "testuser",
		Avatar:   "profile_picture.png",
	}
	mockProvider.On("GetUserInfo", mock.Anything, "test-code").Return(mockUserInfo, nil)

	mockUserInfoStore := &MockUserInfoStore{}
	mockUserInfoStore.On("getAccountInfo", mock.Anything, "test-provider", mockUserInfo).
		Return(&accounts.Account{
			Id:       123,
			Username: mockUserInfo.Username,
			License:  "license",
		}, nil)
	mockUserInfoStore.On("updateUserInfo", mock.Anything, int64(123), "test-provider", mockUserInfo).
		Return(nil)

	oauth := &OAuth2{
		logger: zaptest.NewLogger(t),
		oauthConfigs: map[string]types.IProvider{
			"test-provider": mockProvider,
		},
		userInfoStore: mockUserInfoStore,
		tm:            auth.NewTokenMgr("secret"),
	}
	oauth.RegisterHTTP(router)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/oauth2/callback/test-provider?state=valid&code=test-code",
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
		"/auth/login?oauth2Login=success&u=testuser",
	)
	mockProvider.AssertExpectations(t)
}

func TestCallback_ConnectError(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	store := cookie.NewStore([]byte("secret"))
	sess := sessions.SessionsMany([]string{"fivenet_oauth2_state"}, store)
	router.Use(sess)

	tm := auth.NewTokenMgr("secret")

	mockProvider := createMockProvider()
	mockUserInfo := &types.UserInfo{
		ID:       "123",
		Username: "testuser",
		Avatar:   "profile_picture.png",
	}
	mockProvider.On("GetUserInfo", mock.Anything, "test-code").Return(mockUserInfo, nil)

	mockUserInfoStore := &MockUserInfoStore{}
	mockUserInfoStore.On("getAccountInfo", mock.Anything, "test-provider", mockUserInfo).
		Return(&accounts.Account{
			Id:       123,
			Username: mockUserInfo.Username,
			License:  "license",
		}, nil)

	oauth := &OAuth2{
		logger: zaptest.NewLogger(t),
		oauthConfigs: map[string]types.IProvider{
			"test-provider": mockProvider,
		},
		userInfoStore: mockUserInfoStore,
		tm:            tm,
	}
	oauth.RegisterHTTP(router)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/oauth2/callback/test-provider?state=valid&code=test-code",
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
	t.Parallel()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	store := cookie.NewStore([]byte("secret"))
	sess := sessions.SessionsMany([]string{"fivenet_oauth2_state"}, store)
	router.Use(sess)

	tm := auth.NewTokenMgr("secret")
	token, err := tm.FromCombinedClaims(&authclaims.CombinedClaims{
		AccID:    123,
		Username: "testuser",
	})
	require.NoError(t, err)

	mockProvider := createMockProvider()
	mockUserInfo := &types.UserInfo{
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
		oauthConfigs: map[string]types.IProvider{
			"test-provider": mockProvider,
		},
		userInfoStore: mockUserInfoStore,
		tm:            tm,
	}
	oauth.RegisterHTTP(router)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/oauth2/callback/test-provider?state=valid&code=test-code",
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
	t.Parallel()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	store := cookie.NewStore([]byte("secret"))
	sess := sessions.SessionsMany([]string{"fivenet_oauth2_state"}, store)
	router.Use(sess)

	tm := auth.NewTokenMgr("secret")
	token, err := tm.FromCombinedClaims(&authclaims.CombinedClaims{
		AccID:    123,
		Username: "testuser",
	})
	require.NoError(t, err)

	mockProvider := createMockProvider()
	mockUserInfo := &types.UserInfo{
		ID:       "123",
		Username: "testuser",
		Avatar:   "profile_picture.png",
	}
	redirectUrl := "https://example.com/redirect-url?state="
	mockProvider.On("GetRedirect", mock.Anything).Return(redirectUrl, nil)
	mockProvider.On("GetUserInfo", mock.Anything, "test-code").Return(mockUserInfo, nil)

	mockUserInfoStore := &MockUserInfoStore{}
	mockUserInfoStore.On("getAccountInfo", mock.Anything, "test-provider", mockUserInfo).
		Return(&accounts.Account{
			Id:       123,
			Username: mockUserInfo.Username,
			License:  "license",
		}, nil)
	mockUserInfoStore.On("storeUserInfo", mock.Anything, int64(123), "test-provider", mockUserInfo).
		Return(nil)
	oauth := &OAuth2{
		logger: zaptest.NewLogger(t),
		oauthConfigs: map[string]types.IProvider{
			"test-provider": mockProvider,
		},
		userInfoStore: mockUserInfoStore,
		tm:            tm,
	}
	oauth.RegisterHTTP(router)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/oauth2/login/test-provider?connect-only=true",
		nil,
	)
	req.AddCookie(&http.Cookie{
		Name:  auth.AccCookieName,
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
	require.NotNil(t, state)
	require.NoError(t, session.Save())

	req.AddCookie(&http.Cookie{
		Name:  "fivenet_token",
		Value: token,
	})

	stateVal := state.(string)
	req.URL, err = url.Parse(
		"/api/oauth2/callback/test-provider?state=" + stateVal + "&code=test-code",
	)
	require.NoError(t, err)
	c.Request = req
	sess(c)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Equal(
		t,
		"/auth/account-info/oauth2?oauth2Connect=success&tab=oauth2Connections#",
		w.Header().Get("Location"),
	)
	mockProvider.AssertExpectations(t)
}
