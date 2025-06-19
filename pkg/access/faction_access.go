package access

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
)

// CheckIfHasOwnJobAccess determines if a user has access to a resource based on permission levels, user info, and creator details.
//
// If the document creator job is not equal to the creator's current job, normal access checks need to be applied
// and not the rank attributes checks. If the creator is nil, treat it like a normal doc access check.
// If no levels set, assume "Own" as a safe default.
//
// Returns true if the user has access, otherwise false.
func CheckIfHasOwnJobAccess(
	levels *permissions.StringList, // List of access levels (e.g., Any, Lower_Rank, Same_Rank, Own)
	userInfo *userinfo.UserInfo, // Information about the user requesting access
	creatorJob string, // Job of the document creator
	creator *users.UserShort, // Short info about the creator (may be nil)
) bool {
	// Superusers always have access
	if userInfo.Superuser {
		return true
	}

	// If the document creator job is not equal to the creator's current job, normal access checks need to be applied
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

	// Grant access if any level is "Any"
	if levels.Contains("Any") {
		return true
	}
	// Grant access if user has a higher rank than the creator
	if levels.Contains("Lower_Rank") {
		if creator.JobGrade < userInfo.JobGrade {
			return true
		}
	}
	// Grant access if user has the same or higher rank than the creator
	if levels.Contains("Same_Rank") {
		if creator.JobGrade <= userInfo.JobGrade {
			return true
		}
	}
	// Grant access if user is the creator
	if levels.Contains("Own") {
		if creator.UserId == userInfo.UserId {
			return true
		}
	}

	return false
}
