package auth

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	errorsgrpcauth "github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/errors"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/userinfo"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

const (
	AuthAccIDCtxTag        = "auth.accid"
	AuthActiveCharIDCtxTag = "auth.chrid"
	AuthSubCtxTag          = "auth.sub"

	AuthActiveCharJobCtxTag      = "auth.chrjob"
	AuthActiveCharJobGradeCtxTag = "auth.chrjobg"
)

const (
	PermSuperuserCategory = "Superuser"

	PermCanBeSuperuserName      = "CanBeSuperuser"
	PermCanBeSuperuserGuardName = "superuser-canbesuperuser"

	PermSuperuserName      = "Superuser"
	PermSuperuserGuardName = "superuser-superuser"

	PermAny = "Any"
)

var (
	PermCanBeSuperuser = &permissions.Permission{
		Category:  PermSuperuserCategory,
		Name:      PermCanBeSuperuserName,
		GuardName: PermCanBeSuperuserGuardName,
	}

	PermSuperuser = &permissions.Permission{
		Category:  PermSuperuserCategory,
		Name:      PermSuperuserName,
		GuardName: PermSuperuserGuardName,
	}
)

type userInfoCtxMarker struct{}

var userInfoCtxMarkerKey = &userInfoCtxMarker{}

type GRPCAuth struct {
	ui     userinfo.UserInfoRetriever
	tm     *TokenMgr
	appCfg appconfig.IConfig
}

func NewGRPCAuth(
	ui userinfo.UserInfoRetriever,
	tm *TokenMgr,
	appConfig appconfig.IConfig,
) *GRPCAuth {
	return &GRPCAuth{
		ui:     ui,
		tm:     tm,
		appCfg: appConfig,
	}
}

func (g *GRPCAuth) GRPCAuthFunc(ctx context.Context, _ string) (context.Context, error) {
	t, err := GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, err
	}

	if t == "" {
		return nil, errorsgrpcauth.ErrNoToken
	}

	// Parse token only returns the token info when the token is still valid
	tInfo, err := g.tm.ParseWithClaims(t)
	if err != nil {
		return nil, errorsgrpcauth.ErrInvalidToken
	}

	userInfo, err := g.ui.GetUserInfo(ctx, tInfo.CharID, tInfo.AccID)
	if err != nil {
		// Inject logging fields for better debugging
		if tInfo != nil {
			logging.InjectFields(ctx, logging.Fields{
				AuthSubCtxTag, tInfo.Subject,
				AuthAccIDCtxTag, tInfo.AccID,
				AuthActiveCharIDCtxTag, tInfo.CharID,
			})
		}

		return nil, errswrap.NewError(err, errorsgrpcauth.ErrNoUserInfo)
	}

	newCtx := logging.InjectFields(ctx, logging.Fields{
		AuthSubCtxTag, tInfo.Subject,
		AuthAccIDCtxTag, tInfo.AccID,
		AuthActiveCharIDCtxTag, tInfo.CharID,
		AuthActiveCharJobCtxTag, userInfo.GetJob(),
		AuthActiveCharJobGradeCtxTag, userInfo.GetJobGrade(),
	})

	if userInfo.LastChar != nil && userInfo.GetLastChar() != userInfo.GetUserId() &&
		g.appCfg.Get().Auth.GetLastCharLock() {
		if !userInfo.GetCanBeSuperuser() && !userInfo.GetSuperuser() {
			return nil, errorsgrpcauth.ErrCharLock
		}
	}

	return context.WithValue(newCtx, userInfoCtxMarkerKey, userInfo), nil
}

func (g *GRPCAuth) GRPCAuthFuncWithoutUserInfo(
	ctx context.Context,
	_ string,
) (context.Context, error) {
	t, err := GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, err
	}

	if t == "" {
		return nil, errorsgrpcauth.ErrNoToken
	}

	// Parse token only returns the token info when the token is still valid
	tInfo, err := g.tm.ParseWithClaims(t)
	if err != nil {
		return nil, errorsgrpcauth.ErrInvalidToken
	}

	ctx = logging.InjectFields(ctx, logging.Fields{
		AuthSubCtxTag, tInfo.Subject,
		AuthAccIDCtxTag, tInfo.AccID,
		AuthActiveCharIDCtxTag, tInfo.CharID,
	})

	return ctx, nil
}
