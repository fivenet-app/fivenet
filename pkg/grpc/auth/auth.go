package auth

import (
	"context"

	permissionspermissions "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/permissions"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	errorsgrpcauth "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth/errors"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/pkg/userinfo"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/fx"
)

const (
	AuthAccIDCtxTag        = "auth.accid"
	AuthActiveCharIDCtxTag = "auth.usrid"
	AuthSubCtxTag          = "auth.sub"

	AuthUserJobCtxTag      = "auth.usrjob"
	AuthUserJobGradeCtxTag = "auth.usrjobg"

	AuthUserImpersonateCtxTag         = "auth.imp"
	AuthUserImpersonateJobCtxTag      = "auth.impjob"
	AuthUserImpersonateJobGradeCtxTag = "auth.impjobg"
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
	PermCanBeSuperuser = &permissionspermissions.Permission{
		Category:  PermSuperuserCategory,
		Name:      PermCanBeSuperuserName,
		GuardName: PermCanBeSuperuserGuardName,
	}

	PermSuperuser = &permissionspermissions.Permission{
		Category:  PermSuperuserCategory,
		Name:      PermSuperuserName,
		GuardName: PermSuperuserGuardName,
	}
)

var AuthModule = fx.Module("grpc.auth",
	fx.Provide(
		NewGRPCAuth,
	),
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
	accToken, userToken, err := GetTokensFromGRPCContext(ctx)
	if err != nil {
		return nil, err
	}
	if accToken == "" || userToken == "" {
		return nil, errorsgrpcauth.ErrNoToken
	}

	// Parse tokens only returns the token info when the token is still valid
	accClaims, err := g.tm.ParseAccToken(accToken)
	if err != nil {
		return nil, errorsgrpcauth.ErrInvalidToken
	}

	userClaims, err := g.tm.ParseUserToken(userToken)
	if err != nil {
		return nil, errorsgrpcauth.ErrInvalidToken
	}

	userInfo, err := g.ui.GetUserInfoFromClaims(ctx, userClaims, accClaims)
	if err != nil {
		// Inject logging fields for better debugging
		if accClaims != nil {
			logging.InjectFields(ctx, logging.Fields{
				AuthSubCtxTag, accClaims.Subject,
				AuthAccIDCtxTag, accClaims.AccID,
				AuthActiveCharIDCtxTag, userClaims.UserID,
			})
		}

		return nil, errswrap.NewError(err, errorsgrpcauth.ErrNoUserInfo)
	}

	fields := logging.Fields{
		AuthSubCtxTag, accClaims.Subject,
		AuthAccIDCtxTag, accClaims.AccID,
		AuthActiveCharIDCtxTag, userClaims.UserID,
		AuthUserJobCtxTag, userInfo.GetJob(),
		AuthUserJobGradeCtxTag, userInfo.GetJobGrade(),
	}
	if userClaims.Impersonate != nil {
		fields = append(fields,
			AuthUserImpersonateCtxTag, true,
			AuthUserImpersonateJobCtxTag, userClaims.Impersonate.Job,
			AuthUserImpersonateJobGradeCtxTag, userClaims.Impersonate.JobGrade,
		)
	}

	newCtx := logging.InjectFields(ctx, fields)

	if userClaims.Impersonate != nil {
		userInfo.Job = userClaims.Impersonate.Job
		userInfo.JobGrade = userClaims.Impersonate.JobGrade
	}

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
	t, err := GetAccTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, err
	}
	if t == "" {
		return nil, errorsgrpcauth.ErrNoToken
	}

	// Parse acc token only returns the token info when the token is still valid
	accClaims, err := g.tm.ParseAccToken(t)
	if err != nil {
		return nil, errorsgrpcauth.ErrInvalidToken
	}

	ctx = logging.InjectFields(ctx, logging.Fields{
		AuthSubCtxTag, accClaims.Subject,
		AuthAccIDCtxTag, accClaims.AccID,
		AuthActiveCharIDCtxTag, 0,
	})

	return ctx, nil
}
