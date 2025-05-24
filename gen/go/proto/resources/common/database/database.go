package database

const (
	DefaultMaxPageSize int64 = 20
	NoTotalCount       int64 = -1
)

const (
	AscSortDirection  = "asc"
	DescSortDirection = "desc"
)

type DataCount struct {
	Total int64 `alias:"total"`
}

func (p *PaginationRequest) GetResponse(totalCount int64) (*PaginationResponse, int64) {
	return p.GetResponseWithPageSize(totalCount, DefaultMaxPageSize)
}

func (p *PaginationRequest) GetResponseWithPageSize(totalCount int64, pageSize int64) (*PaginationResponse, int64) {
	if p.PageSize != nil {
		if *p.PageSize <= 0 {
			p.PageSize = &pageSize
		} else if *p.PageSize > pageSize {
			p.PageSize = &pageSize
		}
	} else {
		p.PageSize = &pageSize
	}

	p.Offset = ensureOffsetInRage(p.Offset, *p.PageSize, totalCount)

	return &PaginationResponse{
		TotalCount: totalCount,
		Offset:     p.Offset,
		End:        0,
		PageSize:   *p.PageSize,
	}, *p.PageSize
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
	p.Offset = ensureOffsetInRage(p.Offset, p.PageSize, p.TotalCount)

	p.End = p.Offset + int64(length)
}

func (p *PaginationResponse) UpdateWithTotalCount(totalCount int64, length int) {
	p.TotalCount = totalCount

	p.Update(length)
}
