package userinfo

import (
	"context"
	"errors"

	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	authclaims "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth/claims"
)

// MockUserInfoRetriever is a mock implementation of a user info retriever for testing purposes.
type MockUserInfoRetriever struct {
	UserInfoRetriever

	// UserInfo maps user IDs to their corresponding UserInfo.
	UserInfo map[int32]*pbuserinfo.UserInfo
}

// NewMockUserInfoRetriever creates a new MockUserInfoRetriever with the provided user info map.
func NewMockUserInfoRetriever(userInfo map[int32]*pbuserinfo.UserInfo) *MockUserInfoRetriever {
	return &MockUserInfoRetriever{
		UserInfo: userInfo,
	}
}

// GetUserInfo retrieves the UserInfo for a given userId and accountId.
func (ui *MockUserInfoRetriever) GetUserInfo(
	_ context.Context,
	userId int32,
) (*pbuserinfo.UserInfo, error) {
	if userInfo, ok := ui.UserInfo[userId]; ok {
		return userInfo, nil
	}

	return nil, errors.New("no user info found")
}

// GetUserInfoFromClaims retrieves the UserInfo based on user and account claims.
func (ui *MockUserInfoRetriever) GetUserInfoFromClaims(
	ctx context.Context,
	userClaims *authclaims.UserInfoClaims,
	accClaims *authclaims.AccountInfoClaims,
) (*pbuserinfo.UserInfo, error) {
	return ui.GetUserInfo(ctx, userClaims.UserID)
}

// RefreshUserInfo is a mock method that does nothing and always returns nil.
func (ui *MockUserInfoRetriever) RefreshUserInfo(
	ctx context.Context,
	userId int32,
) error {
	return nil
}
