package database

const (
	DefaultPageLimit int64 = 20
)

type DataCount struct {
	TotalCount int64 // alias:"total_count"
}

func (p *PaginationRequest) GetResponse() (*PaginationResponse, int64) {
	return p.GetResponseWithPageSize(DefaultPageLimit)
}

func (p *PaginationRequest) GetResponseWithPageSize(maxPageSize int64) (*PaginationResponse, int64) {
	if p.PageSize != nil {
		if *p.PageSize <= 0 {
			p.PageSize = &maxPageSize
		} else if *p.PageSize > maxPageSize {
			p.PageSize = &maxPageSize
		}
	} else {
		p.PageSize = &maxPageSize
	}

	return &PaginationResponse{
		TotalCount: 0,
		Offset:     p.Offset,
		End:        0,
		PageSize:   *p.PageSize,
	}, *p.PageSize
}

func (p *PaginationResponse) Update(totalCount int64, length int) {
	p.TotalCount = totalCount

	if p.Offset >= p.TotalCount {
		p.Offset = 0
	}

	p.End = p.Offset + int64(length)
}
