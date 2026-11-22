package database

const (
	// DefaultMaxPageSize is the default maximum number of items per page.
	// It is used when no page size is specified in the request and no custom max page size is set.
	DefaultMaxPageSize int64 = 20

	// User-facing query page sizes. Keep per-query limits here until they become configurable.
	// Keep these in sync with the client sizes in app/utils/pagination.ts.
	FileStoreListFilesPageSize      int64 = 50
	WikiPageListPageSize            int64 = 250
	WikiPageActivityPageSize        int64 = 15
	MailerEmailPageSize             int64 = 20
	MailerThreadPageSize            int64 = 16
	MailerMessagePageSize           int64 = 20
	MailerThreadSearchPageSize      int64 = 15
	SettingsAuditLogPageSize        int64 = 30
	SettingsAccountPageSize         int64 = 30
	QualificationStorePageSize      int64 = 15
	QualificationListPageSize       int64 = 10
	QualificationActivityPageSize   int64 = 15
	DocumentListPageSize            int64 = 20
	DocumentActivityPageSize        int64 = 15
	DocumentCommentPageSize         int64 = 10
	DocumentRequestPageSize         int64 = 10
	DocumentStampPageSize           int64 = 20
	DocumentPinPageSize             int64 = 50
	DocumentUserPageSize            int64 = 20
	DocumentApprovalTasksPageSize   int64 = 20
	VehicleListPageSize             int64 = 20
	VehicleActivityPageSize         int64 = 20
	CitizenActivityPageSize         int64 = 20
	JobAssetsPageSize               int64 = 20
	JobPageSize                     int64 = 20
	JobTimeclockPageSize            int64 = 30
	JobInactiveTimeclockPageSize    int64 = 20
	CentrumUnitActivityPageSize     int64 = 15
	CentrumDispatchPageSize         int64 = 40
	CentrumDispatchActivityPageSize int64 = 15

	// NoTotalCount is a special value indicating that the total count is not available.
	NoTotalCount int64 = -1
)

const (
	// AscSortDirection ascending sort direction.
	//
	// Deprecated: the Sort message has been updated to use a boolean to indicate the sort direction.
	AscSortDirection = "asc"
	// DescSortDirection descending sort direction.
	//
	// Deprecated: the Sort message has been updated to use a boolean to indicate the sort direction.
	DescSortDirection = "desc"
)

// DataCount is a struct that holds the total count of data, pages, etc.
type DataCount struct {
	Total int64 `alias:"total"`
}

// GetColumn returns the first column in the sort.
//
// Deprecated: used to help migrating to the new sorter system.
func (x *Sort) GetColumn() string {
	if x == nil || len(x.GetColumns()) == 0 {
		return ""
	}

	return x.GetColumns()[0].GetId()
}

// GetDirection returns the first column's direction in the sort.
//
// Deprecated: used to help migrating to the new sorter system.
func (x *Sort) GetDirection() string {
	if x == nil || len(x.GetColumns()) == 0 {
		return ""
	}

	if x.GetColumns()[0].GetDesc() {
		return DescSortDirection
	}
	return AscSortDirection
}
