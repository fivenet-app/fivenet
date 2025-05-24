package auth

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	AuthAccIDCtxTag         = "auth.accid"
	AuthActiveCharIDCtxTag  = "auth.chrid"
	AuthSubCtxTag           = "auth.sub"
	AuthActiveCharJobCtxTag = "auth.chrjob"
)

const (
	PermCanBeSuper    = "CanBeSuper"
	PermCanBeSuperKey = "canbesuper"

	PermSuperuser    = "Superuser"
	PermSuperuserKey = "superuser"

	PermAny = "Any"
)

type userInfoCtxMarker struct{}

var userInfoCtxMarkerKey = &userInfoCtxMarker{}

type GRPCAuth struct {
	ui     userinfo.UserInfoRetriever
	tm     *TokenMgr
	appCfg appconfig.IConfig
}

func NewGRPCAuth(ui userinfo.UserInfoRetriever, tm *TokenMgr, appConfig appconfig.IConfig) *GRPCAuth {
	return &GRPCAuth{
		ui:     ui,
		tm:     tm,
		appCfg: appConfig,
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
		AuthActiveCharJobCtxTag, userInfo.Job,
	})

	trace.SpanFromContext(ctx).SetAttributes(
		attribute.Int64("fivenet.auth.acc_id", int64(tInfo.AccID)),
		attribute.Int("fivenet.auth.char_id", int(tInfo.CharID)),
		attribute.String("fivenet.job", userInfo.Job),
	)

	if userInfo.LastChar != nil && *userInfo.LastChar != userInfo.UserId && g.appCfg.Get().Auth.LastCharLock {
		if !userInfo.CanBeSuper && !userInfo.Superuser {
			return nil, ErrCharLock
		}
	}

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
