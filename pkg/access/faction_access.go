package access

import (
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
)

// CheckIfHasOwnJobAccess determines if a user has access to a resource based on permission levels, user info, and creator details.
//
// If the document creator job is not equal to the creator's current job, normal access checks need to be applied
// and not the rank attributes checks. If the creator is nil, treat it like a normal doc access check.
// If no levels set, assume "Own" as a safe default.
//
// Returns true if the user has access, otherwise false.
func CheckIfHasOwnJobAccess(
	levels *permissionsattributes.StringList, // List of access levels (e.g., Any, Lower_Rank, Same_Rank, Own)
	userInfo *pbuserinfo.UserInfo, // Information about the user requesting access
	creatorJob string, // Job of the document creator
	creator *usershort.UserShort, // Short info about the creator (may be nil)
) bool {
	// Superusers always have access
	if userInfo.GetSuperuser() {
		return true
	}

	// If the document creator job is not equal to the creator's current job, normal access checks need to be applied
	if creatorJob != userInfo.GetJob() {
		return true
	}

	// If the creator is nil, treat it like a normal doc access check
	if creator == nil {
		return true
	}

	// If no levels set, assume "Own" as a safe default
	if levels.Len() == 0 {
		return creator.GetUserId() == userInfo.GetUserId()
	}

	// Grant access if any level is "Any"
	if levels.Contains("Any") {
		return true
	}
	// Grant access if user has a higher rank than the creator
	if levels.Contains("Lower_Rank") {
		if creator.GetJobGrade() < userInfo.GetJobGrade() {
			return true
		}
	}
	// Grant access if user has the same or higher rank than the creator
	if levels.Contains("Same_Rank") {
		if creator.GetJobGrade() <= userInfo.GetJobGrade() {
			return true
		}
	}
	// Grant access if user is the creator
	if levels.Contains("Own") {
		if creator.GetUserId() == userInfo.GetUserId() {
			return true
		}
	}

	return false
}
