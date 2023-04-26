package auth

import (
	"context"

	"github.com/galexrt/fivenet/pkg/auth"
	accounts "github.com/galexrt/fivenet/proto/resources/accounts"
	jet "github.com/go-jet/jet/v2/mysql"
)

func (s *Server) GetAccountInfo(ctx context.Context, req *GetAccountInfoRequest) (*GetAccountInfoResponse, error) {
	claims, err := s.tm.ParseWithClaims(auth.MustGetTokenFromGRPCContext(ctx))
	if err != nil {
		return nil, GenericLoginErr
	}

	// Load account
	acc, err := s.getAccountFromDB(ctx, account.ID.EQ(jet.Uint64(claims.AccountID)))
	if err != nil {
		return nil, GenericLoginErr
	}
	if acc.ID == 0 {
		return nil, GenericLoginErr
	}

	return &GetAccountInfoResponse{
		Account: accounts.ConvertFromAcc(acc),
	}, nil
}
