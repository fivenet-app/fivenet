package auth

import (
	"context"

	pbauth "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	errorsgrpcauth "github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/errors"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) DeleteOAuth2Connection(
	ctx context.Context,
	req *pbauth.DeleteOAuth2ConnectionRequest,
) (*pbauth.DeleteOAuth2ConnectionResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.auth.oauth2_provider", req.GetProvider()})

	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, errorsgrpcauth.ErrInvalidToken
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, errswrap.NewError(err, ErrGenericAccount)
	}

	tOAuth2Accs := table.FivenetAccountsOauth2

	stmt := tOAuth2Accs.
		DELETE().
		WHERE(jet.AND(
			tOAuth2Accs.AccountID.EQ(jet.Int64(claims.AccID)),
			tOAuth2Accs.Provider.EQ(jet.String(req.GetProvider())),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, ErrGenericAccount)
	}

	return &pbauth.DeleteOAuth2ConnectionResponse{
		Success: true,
	}, nil
}
