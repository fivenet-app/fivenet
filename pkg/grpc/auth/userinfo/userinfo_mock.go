package userinfo

import "fmt"

type MockUserInfoRetriever struct {
	userInfo map[int32]*UserInfo
}

func NewMockUserInfoRetriever(userInfo map[int32]*UserInfo) *MockUserInfoRetriever {
	return &MockUserInfoRetriever{
		userInfo: userInfo,
	}
}

func (ui *MockUserInfoRetriever) GetUserInfo(userId int32) (*UserInfo, error) {
	if userInfo, ok := ui.userInfo[userId]; ok {
		return userInfo, nil
	}

	return nil, fmt.Errorf("no user info found")
}
