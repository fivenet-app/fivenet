package auth

import (
	"context"

	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

var (
	tOAuth2Accs = table.FivenetOauth2Accounts
)

func (s *Server) DeleteOAuth2Connection(ctx context.Context, req *DeleteOAuth2ConnectionRequest) (*DeleteOAuth2ConnectionResponse, error) {
	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, auth.InvalidTokenErr
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, GenericAccountErr
	}

	stmt := tOAuth2Accs.
		DELETE().
		WHERE(jet.AND(
			tOAuth2Accs.AccountID.EQ(jet.Uint64(claims.AccID)),
			tOAuth2Accs.Provider.EQ(jet.String(req.Provider)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, GenericAccountErr
	}

	return &DeleteOAuth2ConnectionResponse{
		Success: true,
	}, nil
}
