package modelhelper

const AnyJobGradeHasAccess = 0

type AccessRole string

const (
	BlockedAccessRole = "blocked"
	ViewAccessRole    = "view"
	EditAccessRole    = "edit"
	LeaderAccessRole  = "leader"
	AdminAccessRole   = "admin"
)
