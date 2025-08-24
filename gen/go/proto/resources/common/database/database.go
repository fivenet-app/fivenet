package database

const (
	// DefaultMaxPageSize is the default maximum number of items per page.
	// It is used when no page size is specified in the request and no custom max page size is set.
	DefaultMaxPageSize int64 = 20
	// NoTotalCount is a special value indicating that the total count is not available.
	NoTotalCount int64 = -1
)

const (
	// AscSortDirection ascending sort direction.
	// Deprecated: the Sort message has been updated to use a boolean to indicate the sort direction.
	AscSortDirection = "asc"
	// DescSortDirection descending sort direction.
	// Deprecated: the Sort message has been updated to use a boolean to indicate the sort direction.
	DescSortDirection = "desc"
)

// DataCount is a struct that holds the total count of data, pages, etc.
type DataCount struct {
	Total int64 `alias:"total"`
}

// GetColumn returns the first column in the sort.
// Deprecated: used to help migrating to the new sorter system.
func (x *Sort) GetColumn() string {
	if x == nil || len(x.GetColumns()) == 0 {
		return ""
	}

	return x.GetColumns()[0].GetId()
}

// GetDirection returns the first column's direction in the sort.
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
