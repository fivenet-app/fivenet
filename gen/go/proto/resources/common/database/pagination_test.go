package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaginationRequestGetResponseNil(t *testing.T) {
	t.Parallel()

	resp, limit := (*PaginationRequest)(nil).GetResponse(42)

	require.NotNil(t, resp)
	assert.Equal(t, int64(42), resp.GetTotalCount())
	assert.Equal(t, int64(0), resp.GetOffset())
	assert.Equal(t, int64(0), resp.GetEnd())
	assert.Equal(t, DefaultMaxPageSize, resp.GetPageSize())
	assert.Equal(t, DefaultMaxPageSize, limit)
}

func TestPaginationRequestGetResponseWithPageSizeNormalizesPageSize(t *testing.T) {
	t.Parallel()

	req := &PaginationRequest{Offset: 7}
	resp, limit := req.GetResponseWithPageSize(100, 20)

	require.NotNil(t, resp)
	assert.Equal(t, int64(100), resp.GetTotalCount())
	assert.Equal(t, int64(7), resp.GetOffset())
	assert.Equal(t, int64(0), resp.GetEnd())
	assert.Equal(t, int64(20), resp.GetPageSize())
	assert.Equal(t, int64(20), limit)
}

func TestPaginationRequestGetResponseWithPageSizeClampsOffset(t *testing.T) {
	t.Parallel()

	pageSize := int64(20)
	req := &PaginationRequest{Offset: 150, PageSize: &pageSize}
	resp, limit := req.GetResponseWithPageSize(95, 50)

	require.NotNil(t, resp)
	assert.Equal(t, int64(95), resp.GetTotalCount())
	assert.Equal(t, int64(75), resp.GetOffset())
	assert.Equal(t, int64(0), resp.GetEnd())
	assert.Equal(t, int64(20), resp.GetPageSize())
	assert.Equal(t, int64(20), limit)
}

func TestPaginationRequestGetResponseWithPageSizeKeepsNoTotalCountOffset(t *testing.T) {
	t.Parallel()

	pageSize := int64(15)
	req := &PaginationRequest{Offset: 33, PageSize: &pageSize}
	resp, limit := req.GetResponseWithPageSize(NoTotalCount, 50)

	require.NotNil(t, resp)
	assert.Equal(t, NoTotalCount, resp.GetTotalCount())
	assert.Equal(t, int64(33), resp.GetOffset())
	assert.Equal(t, int64(0), resp.GetEnd())
	assert.Equal(t, int64(15), resp.GetPageSize())
	assert.Equal(t, int64(15), limit)
}

func TestPaginationResponseUpdateClampsEnd(t *testing.T) {
	t.Parallel()

	resp := &PaginationResponse{
		TotalCount: 95,
		Offset:     100,
		PageSize:   20,
	}

	resp.Update(20)

	assert.Equal(t, int64(75), resp.GetOffset())
	assert.Equal(t, int64(95), resp.GetEnd())
}

func TestPaginationResponseUpdateWithNoTotalCountKeepsEnd(t *testing.T) {
	t.Parallel()

	resp := &PaginationResponse{
		TotalCount: NoTotalCount,
		Offset:     12,
		PageSize:   20,
	}

	resp.Update(7)

	assert.Equal(t, int64(12), resp.GetOffset())
	assert.Equal(t, int64(19), resp.GetEnd())
}
