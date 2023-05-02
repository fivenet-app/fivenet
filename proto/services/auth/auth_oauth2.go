package auth

import (
	"context"

	"github.com/galexrt/fivenet/pkg/auth"
	jet "github.com/go-jet/jet/v2/mysql"
)

func (s *Server) OAuth2Disconnect(ctx context.Context, req *OAuth2DisconnectRequest) (*OAuth2DisconnectResponse, error) {
	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, auth.InvalidTokenErr
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, GenericAccountErr
	}

	stmt := oAuth2Accounts.
		DELETE().
		WHERE(jet.AND(
			oAuth2Accounts.AccountID.EQ(jet.Uint64(claims.AccountID)),
			oAuth2Accounts.Provider.EQ(jet.String(req.Provider)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, GenericAccountErr
	}

	return &OAuth2DisconnectResponse{
		Success: true,
	}, nil
}
