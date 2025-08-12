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
	AscSortDirection = "asc"
	// DescSortDirection descending sort direction.
	DescSortDirection = "desc"
)

// DataCount is a struct that holds the total count of data, pages, etc.
type DataCount struct {
	Total int64 `alias:"total"`
}

func (p *PaginationRequest) GetResponse(totalCount int64) (*PaginationResponse, int64) {
	return p.GetResponseWithPageSize(totalCount, DefaultMaxPageSize)
}

func (p *PaginationRequest) GetResponseWithPageSize(
	totalCount int64,
	pageSize int64,
) (*PaginationResponse, int64) {
	if p.PageSize != nil {
		if p.GetPageSize() <= 0 {
			p.PageSize = &pageSize
		} else if p.GetPageSize() > pageSize {
			p.PageSize = &pageSize
		}
	} else {
		p.PageSize = &pageSize
	}

	p.Offset = ensureOffsetInRage(p.GetOffset(), p.GetPageSize(), totalCount)

	return &PaginationResponse{
		TotalCount: totalCount,
		Offset:     p.GetOffset(),
		End:        0,
		PageSize:   p.GetPageSize(),
	}, p.GetPageSize()
}

func ensureOffsetInRage(offset int64, pageSize int64, totalCount int64) int64 {
	if totalCount != 0 && totalCount != NoTotalCount && offset > totalCount {
		// Set offset to "last" page
		offset = totalCount - pageSize
	}
	// Make sure offset is at least 0
	if offset < 0 {
		offset = 0
	}

	return offset
}

func (p *PaginationResponse) Update(length int) {
	p.Offset = ensureOffsetInRage(p.GetOffset(), p.GetPageSize(), p.GetTotalCount())

	p.End = p.GetOffset() + int64(length)
}

func (p *PaginationResponse) UpdateWithTotalCount(totalCount int64, length int) {
	p.TotalCount = totalCount

	p.Update(length)
}
