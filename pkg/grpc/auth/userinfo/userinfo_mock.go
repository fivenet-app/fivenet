package userinfo

import (
	"context"
	"fmt"
)

type MockUserInfoRetriever struct {
	userInfo map[int32]*UserInfo
}

func NewMockUserInfoRetriever(userInfo map[int32]*UserInfo) *MockUserInfoRetriever {
	return &MockUserInfoRetriever{
		userInfo: userInfo,
	}
}

func (ui *MockUserInfoRetriever) GetUserInfo(ctx context.Context, userId int32, accountId uint64) (*UserInfo, error) {
	if userInfo, ok := ui.userInfo[userId]; ok {
		return userInfo, nil
	}

	return nil, fmt.Errorf("no user info found")
}

func (ui *MockUserInfoRetriever) SetUserInfo(ctx context.Context, accountId uint64, job string, jobGrade int32) error {
	return nil
}
