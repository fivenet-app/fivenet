package access

import (
	"testing"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"github.com/stretchr/testify/assert"
)

func TestCheckIfHasAccess(t *testing.T) {
	tests := []struct {
		name       string
		levels     []string
		userInfo   *userinfo.UserInfo
		creatorJob string
		creator    *users.UserShort
		expected   bool
	}{
		{
			name:   "SuperUser always has access",
			levels: []string{},
			userInfo: &userinfo.UserInfo{
				UserId:    123,
				Job:       "Police",
				JobGrade:  0,
				SuperUser: true,
			},
			creatorJob: "EMS",
			expected:   true,
		},
		{
			name:   "Different creator job (creator nil)",
			levels: []string{},
			userInfo: &userinfo.UserInfo{
				UserId:   123,
				Job:      "Police",
				JobGrade: 2,
			},
			creatorJob: "EMS",
			expected:   true,
		},
		{
			name:   "Same creator job (no levels)",
			levels: []string{},
			userInfo: &userinfo.UserInfo{
				Job:      "Police",
				JobGrade: 2,
				UserId:   456,
			},
			creator: &users.UserShort{
				Job:      "Police",
				JobGrade: 2,
				UserId:   123,
			},
			creatorJob: "Police",
			expected:   false,
		},
		{
			name:   "Nil creator",
			levels: []string{},
			userInfo: &userinfo.UserInfo{
				Job:      "Police",
				JobGrade: 2,
			},
			creatorJob: "Police",
			expected:   true,
		},
		{
			name:   "No levels, default to Own",
			levels: []string{},
			userInfo: &userinfo.UserInfo{
				UserId:   123,
				Job:      "Police",
				JobGrade: 2,
			},
			creator: &users.UserShort{
				UserId:   123,
				Job:      "Police",
				JobGrade: 2,
			},
			creatorJob: "Police",
			expected:   true,
		},
		{
			name:   "Access level Any (higher rank creator)",
			levels: []string{"Any"},
			userInfo: &userinfo.UserInfo{
				UserId:   123,
				Job:      "Police",
				JobGrade: 2,
			},
			creator: &users.UserShort{
				UserId:   456,
				Job:      "Police",
				JobGrade: 10,
			},
			creatorJob: "Police",
			expected:   true,
		},
		{
			name:   "Access level Lower_Rank",
			levels: []string{"Lower_Rank"},
			userInfo: &userinfo.UserInfo{
				Job:      "Police",
				JobGrade: 2,
			},
			creator: &users.UserShort{
				Job:      "Police",
				JobGrade: 1,
			},
			creatorJob: "Police",
			expected:   true,
		},
		{
			name:   "Access level Same_Rank",
			levels: []string{"Same_Rank"},
			userInfo: &userinfo.UserInfo{
				Job:      "Police",
				JobGrade: 2,
			},
			creator: &users.UserShort{
				Job:      "Police",
				JobGrade: 2,
			},
			creatorJob: "Police",
			expected:   true,
		},
		{
			name:   "Access level Own",
			levels: []string{"Own"},
			userInfo: &userinfo.UserInfo{
				UserId:   123,
				Job:      "Police",
				JobGrade: 2,
			},
			creator: &users.UserShort{
				UserId:   123,
				Job:      "Police",
				JobGrade: 2,
			},
			creatorJob: "Police",
			expected:   true,
		},
		{
			name:   "No access (same job)",
			levels: []string{"Own"},
			userInfo: &userinfo.UserInfo{
				UserId:   123,
				Job:      "Police",
				JobGrade: 2,
			},
			creator: &users.UserShort{
				UserId:   456,
				Job:      "Police",
				JobGrade: 2,
			},
			creatorJob: "Police",
			expected:   false,
		},
		{
			name:   "No access (different jobs)",
			levels: []string{"Own"},
			userInfo: &userinfo.UserInfo{
				UserId:   123,
				Job:      "EMS",
				JobGrade: 2,
			},
			creator: &users.UserShort{
				UserId:   456,
				Job:      "Police",
				JobGrade: 2,
			},
			creatorJob: "Police",
			expected:   true,
		},
		{
			name:   "Multiple access levels",
			levels: []string{"Own", "Lower_Rank"},
			userInfo: &userinfo.UserInfo{
				UserId:   123,
				Job:      "Police",
				JobGrade: 2,
			},
			creator: &users.UserShort{
				UserId:   456,
				Job:      "Police",
				JobGrade: 2,
			},
			creatorJob: "Police",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckIfHasAccess(&permissions.StringList{
				Strings: tt.levels,
			}, tt.userInfo, tt.creatorJob, tt.creator)
			assert.Equal(t, tt.expected, result, "Test case: %s", tt.name)
		})
	}
}
