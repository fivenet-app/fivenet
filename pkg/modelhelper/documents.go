package modelhelper

type DocumentType string

const (
	HTMLDocumentType  DocumentType = "html"
	PlainDocumentType DocumentType = "plain"
)

const AnyJobGradeHasAccess = 0

type AccessRole string

const (
	BlockedAccessRole = "blocked"
	ViewAccessRole    = "view"
	EditAccessRole    = "edit"
	LeaderAccessRole  = "leader"
	AdminAccessRole   = "admin"
)
