package auth

import (
	"context"

	pbauth "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tOAuth2Accs = table.FivenetOauth2Accounts

func (s *Server) DeleteOAuth2Connection(ctx context.Context, req *pbauth.DeleteOAuth2ConnectionRequest) (*pbauth.DeleteOAuth2ConnectionResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.String("fivenet.auth.oauth2_provider", req.Provider))

	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, auth.ErrInvalidToken
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, errswrap.NewError(err, ErrGenericAccount)
	}

	stmt := tOAuth2Accs.
		DELETE().
		WHERE(jet.AND(
			tOAuth2Accs.AccountID.EQ(jet.Uint64(claims.AccID)),
			tOAuth2Accs.Provider.EQ(jet.String(req.Provider)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, ErrGenericAccount)
	}

	return &pbauth.DeleteOAuth2ConnectionResponse{
		Success: true,
	}, nil
}
