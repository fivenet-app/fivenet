package auth

import (
	"context"

	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	AuthAccIDCtxTag        = "auth.accid"
	AuthActiveCharIDCtxTag = "auth.chrid"
	AuthSubCtxTag          = "auth.sub"
)

const (
	PermSuperUser = "SuperUser"
	PermAny       = "Any"
)

type userInfoCtxMarker struct{}

var userInfoCtxMarkerKey = &userInfoCtxMarker{}

var (
	ErrNoToken          = status.Errorf(codes.Unauthenticated, "errors.pkg-auth.ErrNoToken")
	ErrInvalidToken     = status.Error(codes.Unauthenticated, "errors.pkg-auth.ErrInvalidToken")
	ErrCheckToken       = status.Error(codes.Unauthenticated, "errors.pkg-auth.ErrCheckToken")
	ErrUserNoPerms      = status.Error(codes.PermissionDenied, "errors.pkg-auth.ErrUserNoPerms")
	ErrNoUserInfo       = status.Error(codes.Unauthenticated, "errors.pkg-auth.ErrNoUserInfo")
	ErrPermissionDenied = status.Errorf(codes.PermissionDenied, "errors.pkg-auth.ErrPermissionDenied")
)

type GRPCAuth struct {
	ui userinfo.UserInfoRetriever
	tm *TokenMgr
}

func NewGRPCAuth(ui userinfo.UserInfoRetriever, tm *TokenMgr) *GRPCAuth {
	return &GRPCAuth{
		ui: ui,
		tm: tm,
	}
}

func (g *GRPCAuth) GRPCAuthFunc(ctx context.Context, fullMethod string) (context.Context, error) {
	t, err := GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, err
	}

	if t == "" {
		return nil, ErrNoToken
	}

	// Parse token only returns the token info when the token is still valid
	tInfo, err := g.tm.ParseWithClaims(t)
	if err != nil {
		return nil, ErrInvalidToken
	}

	userInfo, err := g.ui.GetUserInfo(ctx, tInfo.CharID, tInfo.AccID)
	if err != nil {
		return nil, err
	}

	ctx = logging.InjectFields(ctx, logging.Fields{
		AuthSubCtxTag, tInfo.Subject,
		AuthAccIDCtxTag, tInfo.CharID,
		AuthActiveCharIDCtxTag, tInfo.CharID,
	})

	trace.SpanFromContext(ctx).SetAttributes(
		attribute.Int64("fivenet.auth.acc_id", int64(tInfo.AccID)),
		attribute.Int("fivenet.auth.char_id", int(tInfo.CharID)),
	)

	return context.WithValue(ctx, userInfoCtxMarkerKey, userInfo), nil
}

func (g *GRPCAuth) GRPCAuthFuncWithoutUserInfo(ctx context.Context, fullMethod string) (context.Context, error) {
	t, err := GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, err
	}

	if t == "" {
		return nil, ErrNoToken
	}

	// Parse token only returns the token info when the token is still valid
	tInfo, err := g.tm.ParseWithClaims(t)
	if err != nil {
		return nil, ErrInvalidToken
	}

	ctx = logging.InjectFields(ctx, logging.Fields{
		AuthSubCtxTag, tInfo.Subject,
		AuthAccIDCtxTag, tInfo.CharID,
		AuthActiveCharIDCtxTag, tInfo.CharID,
	})

	return ctx, nil
}
