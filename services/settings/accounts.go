package settings

import (
	"context"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorssettings "github.com/fivenet-app/fivenet/v2026/services/settings/errors"
)

func (s *Server) ListAccounts(
	ctx context.Context,
	req *pbsettings.ListAccountsRequest,
) (*pbsettings.ListAccountsResponse, error) {
	resp, err := s.store.ListAccounts(ctx, req)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	return resp, nil
}

func (s *Server) UpdateAccount(
	ctx context.Context,
	req *pbsettings.UpdateAccountRequest,
) (*pbsettings.UpdateAccountResponse, error) {
	resp, err := s.store.UpdateAccount(ctx, req)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return resp, nil
}

func (s *Server) DisconnectSocialLogin(
	ctx context.Context,
	req *pbsettings.DisconnectSocialLoginRequest,
) (*pbsettings.DisconnectSocialLoginResponse, error) {
	if err := s.store.DisconnectSocialLogin(ctx, req.GetId(), req.GetProviderName()); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)

	return nil, nil
}

func (s *Server) DeleteAccount(
	ctx context.Context,
	req *pbsettings.DeleteAccountRequest,
) (*pbsettings.DeleteAccountResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if userInfo.GetAccountId() == req.GetId() {
		return nil, errorssettings.ErrCannotDeleteOwnAccount
	}

	account, err := s.store.GetAccountByID(ctx, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}
	if account == nil {
		return &pbsettings.DeleteAccountResponse{}, nil
	}

	var deletedAtTime *timestamp.Timestamp
	if account.GetDeletedAt() == nil || !userInfo.GetSuperuser() {
		deletedAtTime = timestamp.Now()
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	} else {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_RESTORED)
	}

	resp, err := s.store.DeleteAccount(ctx, req.GetId(), deletedAtTime)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	return resp, nil
}
