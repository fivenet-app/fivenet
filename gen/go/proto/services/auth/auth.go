package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/accounts"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
	users "github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	errorsauth "github.com/fivenet-app/fivenet/gen/go/proto/services/auth/errors"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/server/audit"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	grpc "google.golang.org/grpc"
)

var (
	tAccounts  = table.FivenetAccounts
	tUsers     = table.Users.AS("user")
	tUserProps = table.FivenetUserProps.AS("user_props")
	tJobs      = table.Jobs
	tJobGrades = table.JobGrades
	tJobProps  = table.FivenetJobProps.AS("jobprops")
)

type Server struct {
	AuthServiceServer

	logger   *zap.Logger
	db       *sql.DB
	auth     *auth.GRPCAuth
	tm       *auth.TokenMgr
	ps       perms.Permissions
	enricher *mstlystcdata.Enricher
	aud      audit.IAuditer
	ui       userinfo.UserInfoRetriever
	appCfg   appconfig.IConfig
	js       *events.JSWrapper

	domain          string
	oauth2Providers []*config.OAuth2Provider
	customDB        config.CustomDB
	superuserGroups []string
	superuserUsers  []string
}

type Params struct {
	fx.In

	Logger    *zap.Logger
	DB        *sql.DB
	Auth      *auth.GRPCAuth
	TM        *auth.TokenMgr
	Perms     perms.Permissions
	Enricher  *mstlystcdata.Enricher
	Aud       audit.IAuditer
	UI        userinfo.UserInfoRetriever
	Config    *config.Config
	AppConfig appconfig.IConfig
	JS        *events.JSWrapper
}

func NewServer(p Params) *Server {
	return &Server{
		logger:   p.Logger.Named("grpc.auth"),
		db:       p.DB,
		auth:     p.Auth,
		tm:       p.TM,
		ps:       p.Perms,
		enricher: p.Enricher,
		aud:      p.Aud,
		ui:       p.UI,
		appCfg:   p.AppConfig,
		js:       p.JS,

		domain:          p.Config.HTTP.Sessions.Domain,
		oauth2Providers: p.Config.OAuth2.Providers,
		customDB:        p.Config.Database.Custom,
		superuserGroups: p.Config.Auth.SuperuserGroups,
		superuserUsers:  p.Config.Auth.SuperuserUsers,
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterAuthServiceServer(srv, s)
}

// AuthFuncOverride is called instead of exampleAuthFunc
func (s *Server) AuthFuncOverride(ctx context.Context, fullMethod string) (context.Context, error) {
	// Skip authentication for the anon accessible endpoints
	if fullMethod == "/services.auth.AuthService/CreateAccount" ||
		fullMethod == "/services.auth.AuthService/Login" ||
		fullMethod == "/services.auth.AuthService/ForgotPassword" {
		return ctx, nil
	}

	if fullMethod == "/services.auth.AuthService/Logout" {
		c, _ := s.auth.GRPCAuthFunc(ctx, fullMethod)
		if c != nil {
			return c, nil
		}
		return ctx, nil
	}

	if fullMethod == "/services.auth.AuthService/SetSuperUserMode" {
		return s.auth.GRPCAuthFunc(ctx, fullMethod)
	}

	return s.auth.GRPCAuthFuncWithoutUserInfo(ctx, fullMethod)
}

func (s *Server) PermissionUnaryFuncOverride(ctx context.Context, info *grpc.UnaryServerInfo) (context.Context, error) {
	// Skip permission check for the auth services
	return ctx, nil
}

func (s *Server) createTokenFromAccountAndChar(account *model.FivenetAccounts, activeChar *users.User) (string, *auth.CitizenInfoClaims, error) {
	claims := auth.BuildTokenClaimsFromAccount(account, activeChar)
	token, err := s.tm.NewWithClaims(claims)
	return token, claims, err
}

func (s *Server) getAccountFromDB(ctx context.Context, condition jet.BoolExpression) (*model.FivenetAccounts, error) {
	stmt := tAccounts.
		SELECT(
			tAccounts.ID,
			tAccounts.CreatedAt,
			tAccounts.UpdatedAt,
			tAccounts.Enabled,
			tAccounts.Username,
			tAccounts.Password,
			tAccounts.License,
			tAccounts.RegToken,
			tAccounts.LastChar,
			tAccounts.OverrideJob,
			tAccounts.OverrideJobGrade,
			tAccounts.Superuser,
			tAccounts.LastChar,
		).
		FROM(tAccounts).
		WHERE(
			tAccounts.Enabled.IS_TRUE().
				AND(condition),
		).
		LIMIT(1)

	var acc model.FivenetAccounts
	if err := stmt.QueryContext(ctx, s.db, &acc); err != nil {
		return nil, err
	}

	return &acc, nil
}

func (s *Server) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	req.Username = strings.TrimSpace(req.Username)

	account, err := s.getAccountFromDB(ctx, jet.AND(
		tAccounts.Username.EQ(jet.String(req.Username)),
		tAccounts.RegToken.IS_NULL(),
		tAccounts.Password.IS_NOT_NULL(),
	))
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrInvalidLogin)
	}

	// No password set
	if account.Password == nil {
		return nil, errorsauth.ErrNoAccount
	}

	// Password check logic
	if err := bcrypt.CompareHashAndPassword([]byte(*account.Password), []byte(req.Password)); err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrInvalidLogin)
	}

	token, claims, err := s.createTokenFromAccountAndChar(account, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}

	var chooseCharResp *ChooseCharacterResponse
	if s.appCfg.Get().Auth.LastCharLock && account.LastChar != nil {
		ctx = auth.SetTokenInGRPCContext(ctx, token)
		chooseCharResp, err = s.ChooseCharacter(ctx, &ChooseCharacterRequest{
			CharId: *account.LastChar,
		})
		if err != nil {
			chooseCharResp = nil
		}
	}

	s.setTokenCookie(ctx, token)

	return &LoginResponse{
		Expires:   timestamp.New(claims.ExpiresAt.Time),
		AccountId: account.ID,
		Char:      chooseCharResp,
	}, nil
}

func (s *Server) Logout(ctx context.Context, req *LogoutRequest) (*LogoutResponse, error) {
	s.destroyTokenCookie(ctx)

	return &LogoutResponse{
		Success: true,
	}, nil
}

func (s *Server) CreateAccount(ctx context.Context, req *CreateAccountRequest) (*CreateAccountResponse, error) {
	if !s.appCfg.Get().Auth.SignupEnabled {
		return nil, errorsauth.ErrSignupDisabled
	}

	acc, err := s.getAccountFromDB(ctx, tAccounts.RegToken.EQ(jet.String(req.RegToken)))
	if err != nil {
		s.logger.Error("failed to get account from database by registration token", zap.Error(err), zap.String("reg_token", req.RegToken))
		return nil, errswrap.NewError(err, errorsauth.ErrAccountCreateFailed)
	}

	if acc.Username != nil || acc.Password != nil {
		return nil, errorsauth.ErrAccountExistsFailed
	}

	req.Username = strings.TrimSpace(req.Username)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrAccountCreateFailed)
	}

	stmt := tAccounts.
		UPDATE(
			tAccounts.Username,
			tAccounts.Password,
			tAccounts.RegToken,
		).
		SET(
			tAccounts.Username.SET(jet.String(req.Username)),
			tAccounts.Password.SET(jet.String(string(hashedPassword))),
			tAccounts.RegToken.SET(jet.StringExp(jet.NULL)),
		).
		WHERE(
			jet.AND(
				tAccounts.ID.EQ(jet.Uint64(acc.ID)),
				tAccounts.RegToken.EQ(jet.String(req.RegToken)),
			),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		if dbutils.IsDuplicateError(err) {
			return nil, errswrap.NewError(err, errorsauth.ErrAccountDuplicate)
		}

		s.logger.Error("failed to update account in database during account creation", zap.Error(err))
		return nil, errswrap.NewError(err, errorsauth.ErrAccountCreateFailed)
	}

	return &CreateAccountResponse{
		AccountId: acc.ID,
	}, nil
}

func (s *Server) ChangePassword(ctx context.Context, req *ChangePasswordRequest) (*ChangePasswordResponse, error) {
	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, auth.ErrInvalidToken)
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrChangePassword)
	}

	acc, err := s.getAccountFromDB(ctx, tAccounts.ID.EQ(jet.Uint64(claims.AccID)).
		AND(tAccounts.Username.EQ(jet.String(claims.Username))))
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrChangePassword)
	}

	// No password set
	if acc.Password == nil {
		return nil, errswrap.NewError(err, errorsauth.ErrNoAccount)
	}

	// Password check logic
	if err := bcrypt.CompareHashAndPassword([]byte(*acc.Password), []byte(req.Current)); err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrChangePassword)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.New), 14)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrAccountCreateFailed)
	}

	var char *users.User
	if claims.CharID > 0 {
		char, _, _, err = s.getCharacter(ctx, claims.CharID)
		if err != nil {
			return nil, errswrap.NewError(err, errorsauth.ErrChangePassword)
		}
	}

	pass := string(hashedPassword)
	acc.Password = &pass

	stmt := tAccounts.
		UPDATE(
			tAccounts.Password,
		).
		SET(
			tAccounts.Password.SET(jet.String(pass)),
		).
		WHERE(
			tAccounts.ID.EQ(jet.Uint64(acc.ID)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrUpdateAccount)
	}

	newToken, newClaims, err := s.createTokenFromAccountAndChar(acc, char)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrChangePassword)
	}

	s.setTokenCookie(ctx, newToken)

	return &ChangePasswordResponse{
		Expires: timestamp.New(newClaims.ExpiresAt.Time),
	}, nil
}

func (s *Server) ChangeUsername(ctx context.Context, req *ChangeUsernameRequest) (*ChangeUsernameResponse, error) {
	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, auth.ErrInvalidToken)
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrChangeUsername)
	}

	acc, err := s.getAccountFromDB(ctx, tAccounts.ID.EQ(jet.Uint64(claims.AccID)).
		AND(tAccounts.Username.EQ(jet.String(claims.Username))))
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrChangeUsername)
	}

	// No username nor password set on account, fail
	if acc.Username == nil || acc.Password == nil {
		return nil, errswrap.NewError(err, errorsauth.ErrNoAccount)
	}

	// Make sure current username matches the sent current username
	if !strings.EqualFold(*acc.Username, req.Current) {
		return nil, errorsauth.ErrBadUsername
	}

	req.New = strings.TrimSpace(req.New)
	username := req.New

	// New username is same as current username.. just return here.
	resp := &ChangeUsernameResponse{}
	if *acc.Username == username {
		return nil, errorsauth.ErrBadUsername
	}

	// If there is an account with the new username, fail
	newAcc, err := s.getAccountFromDB(ctx, tAccounts.Username.EQ(jet.String(username)))
	if err != nil && !errors.Is(err, qrm.ErrNoRows) {
		// Other database error
		return nil, errswrap.NewError(err, errorsauth.ErrBadUsername)
	}
	// An account with the requested username was found, fail
	if newAcc != nil {
		return nil, errorsauth.ErrBadUsername
	}

	acc.Username = &username

	stmt := tAccounts.
		UPDATE(
			tAccounts.Username,
		).
		SET(
			tAccounts.Username.SET(jet.String(username)),
		).
		WHERE(
			tAccounts.ID.EQ(jet.Uint64(acc.ID)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrUpdateAccount)
	}

	return resp, nil
}

func (s *Server) ForgotPassword(ctx context.Context, req *ForgotPasswordRequest) (*ForgotPasswordResponse, error) {
	acc, err := s.getAccountFromDB(ctx, jet.AND(
		tAccounts.RegToken.EQ(jet.String(req.RegToken)),
		tAccounts.Username.IS_NOT_NULL(),
		tAccounts.Password.IS_NULL(),
	))
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrForgotPassword)
	}

	// We expect the account to not have a password for a "forgot password" request via token
	if acc == nil || acc.Password != nil {
		return nil, errorsauth.ErrNoAccount
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.New), 14)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrForgotPassword)
	}

	pass := string(hashedPassword)
	acc.Password = &pass

	stmt := tAccounts.
		UPDATE(
			tAccounts.Password,
			tAccounts.RegToken,
		).
		SET(
			tAccounts.Password.SET(jet.String(pass)),
			tAccounts.RegToken.SET(jet.StringExp(jet.NULL)),
		).
		WHERE(
			tAccounts.ID.EQ(jet.Uint64(acc.ID)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrForgotPassword)
	}

	return &ForgotPasswordResponse{}, nil
}

func (s *Server) GetCharacters(ctx context.Context, req *GetCharactersRequest) (*GetCharactersResponse, error) {
	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, auth.ErrInvalidToken)
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}

	// Load account to make sure it (still) exists
	acc, err := s.getAccountFromDB(ctx, tAccounts.ID.EQ(jet.Uint64(claims.AccID)).
		AND(tAccounts.Username.EQ(jet.String(claims.Username))))
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}
	if acc.ID == 0 {
		return nil, errorsauth.ErrGenericLogin
	}

	// Load chars from database
	stmt := tUsers.
		SELECT(
			tUsers.ID,
			dbutils.Columns{
				tUsers.Identifier,
				tUsers.Job,
				tJobs.Label.AS("user.job_label"),
				tUsers.JobGrade,
				tJobGrades.Label.AS("user.job_grade_label"),
				tUsers.Firstname,
				tUsers.Lastname,
				tUsers.Dateofbirth,
				tUsers.Sex,
				tUsers.Height,
				tUsers.PhoneNumber,
				tUserProps.Avatar.AS("user.avatar"),
				s.customDB.Columns.User.GetVisum(tUsers.Alias()),
				s.customDB.Columns.User.GetPlaytime(tUsers.Alias()),
				tUsers.Group.AS("character.group"),
			}.Get()...,
		).
		FROM(tUsers.
			LEFT_JOIN(tJobs,
				tJobs.Name.EQ(tUsers.Job),
			).
			LEFT_JOIN(tJobGrades,
				jet.AND(
					tJobGrades.Grade.EQ(tUsers.JobGrade),
					tJobGrades.JobName.EQ(tUsers.Job),
				),
			).
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tUsers.ID),
			),
		).
		WHERE(
			tUsers.Identifier.LIKE(jet.String(buildCharSearchIdentifier(claims.Subject))),
		).
		ORDER_BY(tUsers.ID).
		LIMIT(10)

	resp := &GetCharactersResponse{}
	if err := stmt.QueryContext(ctx, s.db, &resp.Chars); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsauth.ErrNoCharFound)
		}

		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}

	// If last char lock is enabled ensure to mark the one char as available only
	if s.appCfg.Get().Auth.LastCharLock && acc.LastChar != nil {
		for i := 0; i < len(resp.Chars); i++ {
			if resp.Chars[i].Char.UserId == *acc.LastChar ||
				slices.Contains(s.superuserGroups, resp.Chars[i].Group) || slices.Contains(s.superuserUsers, claims.Subject) {
				resp.Chars[i].Available = true
				continue
			}
		}

		// Sort chars for convience of the user
		slices.SortFunc(resp.Chars, func(a, b *accounts.Character) int {
			switch {
			case !a.Available && b.Available:
				return +1
			case a.Available && !b.Available:
				return -1
			default:
				return 0
			}
		})
	} else {
		for i := 0; i < len(resp.Chars); i++ {
			resp.Chars[i].Available = true
		}
	}

	return resp, nil
}

func buildCharSearchIdentifier(license string) string {
	return "%" + license
}

func (s *Server) getCharacter(ctx context.Context, charId int32) (*users.User, *users.JobProps, string, error) {
	stmt := tUsers.
		SELECT(
			tUsers.ID,
			tUsers.Identifier,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Dateofbirth,
			tUserProps.Avatar.AS("user.avatar"),
			tUsers.Group.AS("group"),
			tJobs.Label.AS("user.job_label"),
			tJobGrades.Label.AS("user.job_grade_label"),
			tJobProps.Theme,
			tJobProps.RadioFrequency,
			tJobProps.QuickButtons,
			tJobProps.LogoURL,
		).
		FROM(
			tUsers.
				LEFT_JOIN(tJobs,
					tJobs.Name.EQ(tUsers.Job),
				).
				LEFT_JOIN(tJobGrades,
					jet.AND(
						tJobGrades.Grade.EQ(tUsers.JobGrade),
						tJobGrades.JobName.EQ(tUsers.Job),
					),
				).
				LEFT_JOIN(tJobProps,
					tJobProps.Job.EQ(tJobs.Name),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUsers.ID),
				),
		).
		WHERE(
			tUsers.ID.EQ(jet.Int32(charId)),
		).
		LIMIT(1)

	var dest struct {
		users.User
		Group    string
		JobProps *users.JobProps
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, nil, "", errorsauth.ErrNoCharFound
		}
		return nil, nil, "", err
	}

	if dest.JobProps != nil {
		s.enricher.EnrichJobName(dest.JobProps)
	}

	return &dest.User, dest.JobProps, dest.Group, nil
}

func (s *Server) ChooseCharacter(ctx context.Context, req *ChooseCharacterRequest) (*ChooseCharacterResponse, error) {
	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, auth.ErrInvalidToken)
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, errswrap.NewError(err, auth.ErrInvalidToken)
	}

	// Load account data for token creation
	account, err := s.getAccountFromDB(ctx, tAccounts.ID.EQ(jet.Uint64(claims.AccID)).
		AND(tAccounts.Username.EQ(jet.String(claims.Username))))
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrNoCharFound)
	}

	char, jProps, userGroup, err := s.getCharacter(ctx, req.CharId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrNoCharFound)
	}

	// Make sure the user isn't sending us a different char ID than their own
	if !strings.HasSuffix(*char.Identifier, claims.Subject) {
		s.logger.Error("user sent bad char!", zap.String("expected", *char.Identifier), zap.String("current", claims.Subject))
		return nil, errorsauth.ErrUnableToChooseChar
	}

	isSuperUser := slices.Contains(s.superuserGroups, userGroup) || slices.Contains(s.superuserUsers, claims.Subject)

	// If char lock is active, make sure that the user is choosing the correct char
	if !isSuperUser &&
		s.appCfg.Get().Auth.LastCharLock &&
		account.LastChar != nil &&
		*account.LastChar != req.CharId {
		return nil, errorsauth.ErrCharLock
	}

	// Reset override jobs when char is not a superuser but has an override set..
	if !isSuperUser &&
		((account.Superuser != nil && *account.Superuser) ||
			account.OverrideJob != nil) {
		account.OverrideJob = nil
		account.OverrideJobGrade = nil

		if err := s.ui.SetUserInfo(ctx, claims.AccID, false, account.OverrideJob, account.OverrideJobGrade); err != nil {
			return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
		}

		not := false
		account.Superuser = &not
	} else if isSuperUser &&
		(account.Superuser != nil && *account.Superuser) && account.OverrideJob != nil && account.OverrideJobGrade != nil {
		char.Job = *account.OverrideJob
		char.JobGrade = *account.OverrideJobGrade

		s.enricher.EnrichJobInfo(char)

		_, _, jProps, err = s.getJobWithProps(ctx, char.Job)
		if err != nil {
			return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
		}
	}

	newToken, newClaims, err := s.createTokenFromAccountAndChar(account, char)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}

	// Load permissions of user
	userPs, err := s.ps.GetPermissionsOfUser(&userinfo.UserInfo{
		UserId:   char.UserId,
		Job:      char.Job,
		JobGrade: char.JobGrade,
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrUnableToChooseChar)
	}

	ps := userPs.GuardNames()
	if isSuperUser {
		ps = append(ps, auth.PermCanBeSuperKey)

		if account.Superuser != nil && *account.Superuser {
			ps = append(ps, auth.PermSuperUserKey)
		}
	}

	attrs, err := s.ps.FlattenRoleAttributes(char.Job, char.JobGrade)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}
	ps = append(ps, attrs...)

	if len(ps) == 0 || (!isSuperUser && !slices.Contains(ps, "authservice-choosecharacter")) {
		return nil, errorsauth.ErrUnableToChooseChar
	}

	s.aud.Log(&model.FivenetAuditLog{
		Service: AuthService_ServiceDesc.ServiceName,
		Method:  "ChooseCharacter",
		UserID:  char.UserId,
		UserJob: char.Job,
		State:   int16(rector.EventType_EVENT_TYPE_VIEWED),
	}, char.UserShort())

	s.setTokenCookie(ctx, newToken)

	return &ChooseCharacterResponse{
		Expires:     timestamp.New(newClaims.ExpiresAt.Time),
		Permissions: ps,
		JobProps:    jProps,
		Char:        char,
		Username:    *account.Username,
	}, nil
}

func (s *Server) SetSuperUserMode(ctx context.Context, req *SetSuperUserModeRequest) (*SetSuperUserModeResponse, error) {
	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, auth.ErrInvalidToken
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, errswrap.NewError(err, auth.ErrInvalidToken)
	}

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if !userInfo.CanBeSuper {
		return nil, errswrap.NewError(err, errorsauth.ErrNoCharFound)
	}

	char, _, _, err := s.getCharacter(ctx, claims.CharID)
	if err != nil {
		return nil, errswrap.NewError(fmt.Errorf("failed to get char %d. %w", claims.CharID, err), errorsauth.ErrNoCharFound)
	}

	var jobProps *users.JobProps

	// Reset override job when switching off superuser mode
	if !req.Superuser {
		userInfo.Job = char.Job
		userInfo.JobGrade = char.JobGrade

		userInfo.OverrideJob = nil
		userInfo.OverrideJobGrade = nil

		// Send original char job props to user
		_, _, jProps, err := s.getJobWithProps(ctx, char.Job)
		if err != nil {
			return nil, errswrap.NewError(fmt.Errorf("failed to get job props from %s job. %w", char.Job, err), errorsauth.ErrGenericLogin)
		}
		jobProps = jProps
	} else if req.Job != nil {
		// Only set job if requested
		job, jobGrade, jProps, err := s.getJobWithProps(ctx, *req.Job)
		if err != nil {
			return nil, errswrap.NewError(fmt.Errorf("failed to get job props from %s job. %w", *req.Job, err), errorsauth.ErrGenericLogin)
		}
		jobProps = jProps

		userInfo.Job = job.Name
		userInfo.JobGrade = jobGrade
		userInfo.OverrideJob = &job.Name
		userInfo.OverrideJobGrade = &jobGrade

		char.Job = job.Name
		char.JobGrade = jobGrade
		s.enricher.EnrichJobInfo(char)
	}

	if err := s.ui.SetUserInfo(ctx, claims.AccID, req.Superuser, userInfo.OverrideJob, userInfo.OverrideJobGrade); err != nil {
		return nil, errswrap.NewError(fmt.Errorf("failed to set user info. %w", err), errorsauth.ErrGenericLogin)
	}

	userInfo.SuperUser = req.Superuser

	// Load account data for token creation
	account, err := s.getAccountFromDB(ctx, tAccounts.Username.EQ(jet.String(claims.Username)))
	if err != nil {
		return nil, errswrap.NewError(fmt.Errorf("failed to get account from db. %w", err), errorsauth.ErrGenericLogin)
	}

	newToken, newClaims, err := s.createTokenFromAccountAndChar(account, char)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}

	s.setTokenCookie(ctx, newToken)

	return &SetSuperUserModeResponse{
		Expires:  timestamp.New(newClaims.ExpiresAt.Time),
		JobProps: jobProps,
		Char:     char,
	}, nil
}

func (s *Server) getJobWithProps(ctx context.Context, jobName string) (*users.Job, int32, *users.JobProps, error) {
	js := tJobs.AS("job")
	stmt := js.
		SELECT(
			js.Name,
			js.Label,
			tJobGrades.Grade.AS("job_grade"),
			tJobProps.Job,
			tJobProps.UpdatedAt,
			tJobProps.Theme,
			tJobProps.LivemapMarkerColor,
			tJobProps.RadioFrequency,
			tJobProps.QuickButtons,
			tJobProps.LogoURL,
		).
		FROM(
			js.
				INNER_JOIN(tJobGrades,
					tJobGrades.JobName.EQ(js.Name),
				).
				LEFT_JOIN(tJobProps,
					tJobProps.Job.EQ(js.Name)),
		).
		WHERE(
			js.Name.EQ(jet.String(jobName)),
		).
		ORDER_BY(tJobGrades.Grade.DESC()).
		LIMIT(1)

	var dest struct {
		Job      *users.Job
		JobGrade int32 `alias:"job_grade"`
		JobProps *users.JobProps
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, 0, nil, err
	}

	if dest.JobProps != nil {
		s.enricher.EnrichJobName(dest.JobProps)
	}

	return dest.Job, dest.JobGrade, dest.JobProps, nil
}
