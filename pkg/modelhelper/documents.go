package modelhelper

const AnyJobGradeHasAccess = 0

type AccessRole string

const (
	BlockedAccessRole AccessRole = "blocked"
	ViewAccessRole               = "view"
	EditAccessRole               = "edit"
	LeaderAccessRole             = "leader"
	AdminAccessRole              = "admin"
)
