package grpc_database

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"google.golang.org/grpc"
)

// PaginationInterceptor is a gRPC interceptor that normalizes pagination parameters in requests.
// It ensures that the page size does not exceed the specified maximum and that offsets are non-negative.
// If maxPageSize is less than or equal to zero, it defaults to database.DefaultMaxPage.
func PaginationInterceptor(maxPageSize int64) grpc.UnaryServerInterceptor {
	if maxPageSize <= 0 {
		maxPageSize = database.DefaultMaxPageSize
	}

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// PRE: normalize request pagination
		if r, ok := req.(database.HasPaginationRequest); ok {
			if p := r.GetPagination(); p != nil {
				ps := p.GetPageSize()
				// default/cap
				if ps <= 0 || ps > maxPageSize {
					ps = maxPageSize
					p.PageSize = &ps // mirrors your helper's behavior :contentReference[oaicite:2]{index=2}
				}
				// clamp negatives; "snap to last page" is total-aware and handled later by your helpers
				if p.GetOffset() < 0 {
					p.Offset = 0
				}
			}
		}

		// Call method handler
		resp, err := handler(ctx, req)
		if err != nil {
			return resp, err
		}

		// POST: finalize response pagination by updating it with the actual item count
		if r, ok := resp.(database.HasPaginationResponse); ok {
			if p := r.GetPagination(); p != nil {
				// Update() sets End = Offset + itemsLen and re-applies offset safety,
				// considering TotalCount when it's known (and leaves it when NoTotalCount). :contentReference[oaicite:3]{index=3}
				p.Update(r.ItemsLen())
			}
		}

		return resp, nil
	}
}
