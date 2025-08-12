package userinfo

import (
	"context"
	"errors"

	pbuserinfo "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
)

// MockUserInfoRetriever is a mock implementation of a user info retriever for testing purposes.
type MockUserInfoRetriever struct {
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
	ctx context.Context,
	userId int32,
	accountId uint64,
) (*pbuserinfo.UserInfo, error) {
	if userInfo, ok := ui.UserInfo[userId]; ok {
		return userInfo, nil
	}

	return nil, errors.New("no user info found")
}

// GetUserInfoWithoutAccountId retrieves the UserInfo for a given userId without requiring an accountId.
func (ui *MockUserInfoRetriever) GetUserInfoWithoutAccountId(
	ctx context.Context,
	userId int32,
) (*pbuserinfo.UserInfo, error) {
	if userInfo, ok := ui.UserInfo[userId]; ok {
		return userInfo, nil
	}

	return nil, errors.New("no user info found")
}

// SetUserInfo is a mock method that does nothing and always returns nil.
func (ui *MockUserInfoRetriever) SetUserInfo(
	ctx context.Context,
	accountId uint64,
	superuser bool,
	job *string,
	jobGrade *int32,
) error {
	return nil
}

// RefreshUserInfo is a mock method that does nothing and always returns nil.
func (ui *MockUserInfoRetriever) RefreshUserInfo(
	ctx context.Context,
	userId int32,
	accountId uint64,
) error {
	return nil
}
