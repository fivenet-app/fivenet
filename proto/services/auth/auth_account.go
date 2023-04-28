package auth

import (
	"context"
	"errors"

	"github.com/galexrt/fivenet/pkg/auth"
	accounts "github.com/galexrt/fivenet/proto/resources/accounts"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	oAuth2Accounts = table.FivenetOauth2Accounts

	GenericAccountErr = status.Error(codes.Internal, "Failed to get/update your account, please try again.")
)

func (s *Server) GetAccountInfo(ctx context.Context, req *GetAccountInfoRequest) (*GetAccountInfoResponse, error) {
	claims, err := s.tm.ParseWithClaims(auth.MustGetTokenFromGRPCContext(ctx))
	if err != nil {
		return nil, GenericAccountErr
	}

	// Load account
	acc, err := s.getAccountFromDB(ctx, account.ID.EQ(jet.Uint64(claims.AccountID)))
	if err != nil {
		return nil, GenericAccountErr
	}
	if acc.ID == 0 {
		return nil, GenericAccountErr
	}

	stmt := oAuth2Accounts.
		SELECT(
			oAuth2Accounts.AllColumns,
		).
		FROM(
			oAuth2Accounts,
		).
		WHERE(
			oAuth2Accounts.AccountID.EQ(jet.Uint64(acc.ID)),
		).
		LIMIT(3)

	oauth2Conns := []*accounts.OAuth2Account{}
	if err := stmt.QueryContext(ctx, s.db, oauth2Conns); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, GenericAccountErr
		}
	}

	return &GetAccountInfoResponse{
		Account:           accounts.ConvertFromAcc(acc),
		Oauth2Connections: oauth2Conns,
	}, nil
}

func (s *Server) OAuth2Disconnect(ctx context.Context, req *OAuth2DisconnectRequest) (*OAuth2DisconnectResponse, error) {
	claims, err := s.tm.ParseWithClaims(auth.MustGetTokenFromGRPCContext(ctx))
	if err != nil {
		return nil, GenericAccountErr
	}

	// TODO validate provider name in some way..

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
