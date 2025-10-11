package database

// HasPaginationRequest that carry pagination information.
type HasPaginationRequest interface {
	GetPagination() *PaginationRequest
}

// HasPaginationResponse that carry pagination and can report item count.
type HasPaginationResponse interface {
	GetPagination() *PaginationResponse
	ItemsLen() int
}

func (p *PaginationRequest) GetResponse(totalCount int64) (*PaginationResponse, int64) {
	return p.GetResponseWithPageSize(totalCount, DefaultMaxPageSize)
}

func (p *PaginationRequest) GetResponseWithPageSize(
	totalCount int64,
	maxPageSize int64,
) (*PaginationResponse, int64) {
	if p == nil {
		return &PaginationResponse{
			TotalCount: totalCount,
			Offset:     0,
			End:        0,
			PageSize:   maxPageSize,
		}, maxPageSize
	}

	pageSize := p.normalizedPageSize(maxPageSize)
	offset := ensureOffsetInRange(p.GetOffset(), pageSize, totalCount)

	// If you really want to persist normalized values back:
	p.PageSize = &pageSize
	p.Offset = offset

	p.Offset = ensureOffsetInRange(p.GetOffset(), p.GetPageSize(), totalCount)

	return &PaginationResponse{
		TotalCount: totalCount,
		Offset:     p.GetOffset(),
		End:        0,
		PageSize:   p.GetPageSize(),
	}, pageSize
}

func (p *PaginationRequest) normalizedPageSize(maxSize int64) int64 {
	ps := p.GetPageSize()
	if ps <= 0 || ps > maxSize {
		ps = maxSize
	}
	return ps
}

func ensureOffsetInRange(offset int64, pageSize int64, totalCount int64) int64 {
	if totalCount != NoTotalCount && totalCount >= 0 {
		// If offset is beyond or exactly at totalCount, snap to last page start.
		if offset >= totalCount {
			offset = totalCount - pageSize
		}
	}

	// Make sure offset is at least 0
	if offset < 0 {
		offset = 0
	}

	return offset
}

func (p *PaginationResponse) Update(length int) {
	p.Offset = ensureOffsetInRange(p.GetOffset(), p.GetPageSize(), p.GetTotalCount())

	end := p.GetOffset() + int64(length)
	if p.GetTotalCount() != NoTotalCount && p.GetTotalCount() >= 0 && end > p.GetTotalCount() {
		end = p.GetTotalCount()
	}
	p.End = end
}

func (p *PaginationResponse) UpdateWithTotalCount(totalCount int64, length int) {
	p.TotalCount = totalCount

	p.Update(length)
}
