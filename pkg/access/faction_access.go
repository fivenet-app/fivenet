package access

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
)

func CheckIfHasAccess(levels *permissions.StringList, userInfo *userinfo.UserInfo, creatorJob string, creator *users.UserShort) bool {
	if userInfo.Superuser {
		return true
	}

	// If the document creator job is not equal to the creator's current job, normal access checks need to be applied
	// and not the rank attributes checks
	if creatorJob != userInfo.Job {
		return true
	}

	// If the creator is nil, treat it like a normal doc access check
	if creator == nil {
		return true
	}

	// If no levels set, assume "Own" as a safe default
	if levels.Len() == 0 {
		return creator.UserId == userInfo.UserId
	}

	if levels.Contains("Any") {
		return true
	}
	if levels.Contains("Lower_Rank") {
		if creator.JobGrade < userInfo.JobGrade {
			return true
		}
	}
	if levels.Contains("Same_Rank") {
		if creator.JobGrade <= userInfo.JobGrade {
			return true
		}
	}
	if levels.Contains("Own") {
		if creator.UserId == userInfo.UserId {
			return true
		}
	}

	return false
}
