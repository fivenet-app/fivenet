package database

const (
	DefaultPageLimit int64 = 20
)

type DataCount struct {
	TotalCount int64 // alias:"total_count"
}

func EmptyPaginationResponse(offset int64) *PaginationResponse {
	return &PaginationResponse{
		TotalCount: 0,
		Offset:     offset,
		End:        0,
		PageSize:   DefaultPageLimit,
	}
}

func PaginationHelper(pag *PaginationResponse, totalCount int64, offset int64, length int) {
	PaginationHelperWithPageSize(pag, totalCount, offset, length, DefaultPageLimit)
}

func PaginationHelperWithPageSize(pag *PaginationResponse, totalCount int64, offset int64, length int, pageSize int64) {
	pag.TotalCount = totalCount
	pag.PageSize = pageSize

	if offset >= pag.TotalCount {
		pag.Offset = 0
	} else {
		pag.Offset = offset
	}

	pag.End = pag.Offset + int64(length)
}
