package auth

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/accounts"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	users "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	pbauth "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/auth"
	permsauth "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/auth/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	errorsgrpcauth "github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/errors"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsauth "github.com/fivenet-app/fivenet/v2025/services/auth/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
)

var (
	tAccounts  = table.FivenetAccounts
	tUserProps = table.FivenetUserProps.AS("user_props")
	tJobProps  = table.FivenetJobProps.AS("job_props")
)

func (s *Server) createTokenFromAccountAndChar(
	account *model.FivenetAccounts,
	activeChar *users.User,
) (string, *auth.CitizenInfoClaims, error) {
	claims := auth.BuildTokenClaimsFromAccount(account, activeChar)
	token, err := s.tm.NewWithClaims(claims)
	return token, claims, err
}

func (s *Server) getAccountFromDB(
	ctx context.Context,
	condition mysql.BoolExpression,
) (*model.FivenetAccounts, error) {
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

	acc := &model.FivenetAccounts{}
	if err := stmt.QueryContext(ctx, s.db, acc); err != nil {
		return nil, err
	}

	return acc, nil
}

func (s *Server) Login(
	ctx context.Context,
	req *pbauth.LoginRequest,
) (*pbauth.LoginResponse, error) {
	req.Username = normalizeUsername(req.GetUsername())

	logging.InjectFields(ctx, logging.Fields{"fivenet.auth.username", req.GetUsername()})

	account, err := s.getAccountFromDB(ctx, mysql.AND(
		tAccounts.Username.EQ(mysql.String(req.GetUsername())),
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

	if err := checkPassword(*account.Password, req.GetPassword()); err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrInvalidLogin)
	}

	token, claims, err := s.createTokenFromAccountAndChar(account, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}

	var chooseCharResp *pbauth.ChooseCharacterResponse
	if s.appCfg.Get().Auth.GetLastCharLock() && account.LastChar != nil {
		ctx = auth.SetTokenInGRPCContext(ctx, token)
		chooseCharResp, err = s.ChooseCharacter(ctx, &pbauth.ChooseCharacterRequest{
			CharId: *account.LastChar,
		})
		if err != nil {
			chooseCharResp = nil
		}
	}

	// If choose char response is null, set the basic login token
	if chooseCharResp == nil {
		if err := s.setTokenCookie(ctx, token); err != nil {
			return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
		}
	}

	return &pbauth.LoginResponse{
		Expires:   timestamp.New(claims.ExpiresAt.Time),
		AccountId: account.ID,
		Char:      chooseCharResp,
	}, nil
}

func (s *Server) Logout(
	ctx context.Context,
	req *pbauth.LogoutRequest,
) (*pbauth.LogoutResponse, error) {
	err := s.destroyTokenCookie(ctx)
	if err != nil {
		s.logger.Error("failed to destroy token cookie", zap.Error(err))
	}

	return &pbauth.LogoutResponse{
		Success: err == nil,
	}, nil
}

func (s *Server) CreateAccount(
	ctx context.Context,
	req *pbauth.CreateAccountRequest,
) (*pbauth.CreateAccountResponse, error) {
	if !s.appCfg.Get().Auth.GetSignupEnabled() {
		return nil, errorsauth.ErrSignupDisabled
	}

	acc, err := s.getAccountFromDB(ctx, tAccounts.RegToken.EQ(mysql.String(req.GetRegToken())))
	if err != nil {
		s.logger.Error(
			"failed to get account from database by registration token",
			zap.Error(err),
			zap.String("reg_token", req.GetRegToken()),
		)
		return nil, errswrap.NewError(err, errorsauth.ErrAccountCreateFailed)
	}

	if acc.Username != nil || acc.Password != nil {
		return nil, errorsauth.ErrAccountExistsFailed
	}

	req.Username = normalizeUsername(req.GetUsername())

	hashedPassword, err := hashPassword(req.GetPassword())
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
			tAccounts.Username.SET(mysql.String(req.GetUsername())),
			tAccounts.Password.SET(mysql.String(hashedPassword)),
			tAccounts.RegToken.SET(mysql.StringExp(mysql.NULL)),
		).
		WHERE(mysql.AND(
			tAccounts.ID.EQ(mysql.Int64(acc.ID)),
			tAccounts.RegToken.EQ(mysql.String(req.GetRegToken())),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		if dbutils.IsDuplicateError(err) {
			return nil, errswrap.NewError(err, errorsauth.ErrAccountDuplicate)
		}

		s.logger.Error(
			"failed to update account in database during account creation",
			zap.Error(err),
		)
		return nil, errswrap.NewError(err, errorsauth.ErrAccountCreateFailed)
	}

	return &pbauth.CreateAccountResponse{
		AccountId: acc.ID,
	}, nil
}

func (s *Server) ChangePassword(
	ctx context.Context,
	req *pbauth.ChangePasswordRequest,
) (*pbauth.ChangePasswordResponse, error) {
	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, errorsgrpcauth.ErrInvalidToken)
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrChangePassword)
	}

	acc, err := s.getAccountFromClaims(ctx, claims)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrChangePassword)
	}

	// No password set
	if acc.Password == nil {
		return nil, errswrap.NewError(err, errorsauth.ErrNoAccount)
	}

	// Password check logic

	if err := checkPassword(*acc.Password, req.GetCurrent()); err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrChangePassword)
	}

	hashedPassword, err := hashPassword(req.GetNew())
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

	pass := hashedPassword
	acc.Password = &pass

	stmt := tAccounts.
		UPDATE(
			tAccounts.Password,
		).
		SET(
			tAccounts.Password.SET(mysql.String(pass)),
		).
		WHERE(
			tAccounts.ID.EQ(mysql.Int64(acc.ID)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrUpdateAccount)
	}

	newToken, newClaims, err := s.createTokenFromAccountAndChar(acc, char)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrChangePassword)
	}

	if err := s.setTokenCookie(ctx, newToken); err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}

	return &pbauth.ChangePasswordResponse{
		Expires: timestamp.New(newClaims.ExpiresAt.Time),
	}, nil
}

func (s *Server) ChangeUsername(
	ctx context.Context,
	req *pbauth.ChangeUsernameRequest,
) (*pbauth.ChangeUsernameResponse, error) {
	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, errorsgrpcauth.ErrInvalidToken)
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrChangeUsername)
	}

	acc, err := s.getAccountFromClaims(ctx, claims)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrChangeUsername)
	}

	// No username nor password set on account, fail
	if acc.Username == nil || acc.Password == nil {
		return nil, errswrap.NewError(err, errorsauth.ErrNoAccount)
	}

	// Make sure current username matches the sent current username
	if !strings.EqualFold(*acc.Username, req.GetCurrent()) {
		return nil, errorsauth.ErrBadUsername
	}

	req.New = normalizeUsername(req.GetNew())
	username := req.GetNew()

	// New username is same as current username.. just return here.
	resp := &pbauth.ChangeUsernameResponse{}
	if *acc.Username == username {
		return nil, errorsauth.ErrBadUsername
	}

	// If there is an account with the new username, fail
	newAcc, err := s.getAccountFromDB(ctx, tAccounts.Username.EQ(mysql.String(username)))
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
			tAccounts.Username.SET(mysql.String(username)),
		).
		WHERE(
			tAccounts.ID.EQ(mysql.Int64(acc.ID)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrUpdateAccount)
	}

	return resp, nil
}

func (s *Server) ForgotPassword(
	ctx context.Context,
	req *pbauth.ForgotPasswordRequest,
) (*pbauth.ForgotPasswordResponse, error) {
	acc, err := s.getAccountFromDB(ctx, mysql.AND(
		tAccounts.RegToken.EQ(mysql.String(req.GetRegToken())),
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

	hashedPassword, err := hashPassword(req.GetNew())
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrForgotPassword)
	}

	pass := hashedPassword
	acc.Password = &pass

	stmt := tAccounts.
		UPDATE(
			tAccounts.Password,
			tAccounts.RegToken,
		).
		SET(
			tAccounts.Password.SET(mysql.String(pass)),
			tAccounts.RegToken.SET(mysql.StringExp(mysql.NULL)),
		).
		WHERE(
			tAccounts.ID.EQ(mysql.Int64(acc.ID)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrForgotPassword)
	}

	return &pbauth.ForgotPasswordResponse{}, nil
}

func (s *Server) GetCharacters(
	ctx context.Context,
	req *pbauth.GetCharactersRequest,
) (*pbauth.GetCharactersResponse, error) {
	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, errorsgrpcauth.ErrInvalidToken)
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}

	// Load account to make sure it (still) exists
	acc, err := s.getAccountFromClaims(ctx, claims)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}
	if acc.ID == 0 {
		return nil, errorsauth.ErrGenericLogin
	}

	// Load chars from database
	tUsers := tables.User().AS("user")
	tAvatar := table.FivenetFiles.AS("profile_picture")

	stmt := tUsers.
		SELECT(
			tUsers.ID,
			dbutils.Columns{
				tUsers.ID.AS("user.user_id"),
				tUsers.Identifier,
				tUsers.Job,
				tUsers.JobGrade,
				tUsers.Firstname,
				tUsers.Lastname,
				tUsers.Dateofbirth,
				tUsers.Sex,
				tUsers.Height,
				tUsers.PhoneNumber,
				tUserProps.AvatarFileID.AS("user.profile_picture_file_id"),
				tAvatar.FilePath.AS("user.profile_picture"),
				tUsers.Group.AS("character.group"),
				s.customDB.Columns.User.GetVisum(tUsers.Alias()),
				s.customDB.Columns.User.GetPlaytime(tUsers.Alias()),
			}.Get()...,
		).
		FROM(tUsers.
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tUsers.ID),
			).
			LEFT_JOIN(tAvatar,
				tAvatar.ID.EQ(tUserProps.AvatarFileID),
			),
		).
		WHERE(
			tUsers.Identifier.LIKE(mysql.String(buildCharSearchIdentifier(claims.Subject))),
		).
		ORDER_BY(tUsers.ID).
		LIMIT(10)

	resp := &pbauth.GetCharactersResponse{}
	if err := stmt.QueryContext(ctx, s.db, &resp.Chars); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsauth.ErrNoCharFound)
		}

		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}

	// If last char lock is enabled ensure to mark the one char as available only
	if s.appCfg.Get().Auth.GetLastCharLock() && acc.LastChar != nil {
		for i := range resp.GetChars() {
			s.enricher.EnrichJobInfo(resp.GetChars()[i].GetChar())

			if resp.GetChars()[i].GetChar().GetUserId() == *acc.LastChar ||
				slices.Contains(
					s.superuserGroups,
					resp.GetChars()[i].GetGroup(),
				) || slices.Contains(s.superuserUsers, claims.Subject) {
				resp.Chars[i].Available = true
				continue
			}
		}

		// Sort chars for convience of the user
		slices.SortFunc(resp.GetChars(), func(a, b *accounts.Character) int {
			switch {
			case !a.GetAvailable() && b.GetAvailable():
				return +1
			case a.GetAvailable() && !b.GetAvailable():
				return -1
			default:
				return 0
			}
		})
	} else {
		for i := range resp.GetChars() {
			s.enricher.EnrichJobInfo(resp.GetChars()[i].GetChar())

			resp.Chars[i].Available = true
		}
	}

	return resp, nil
}

func buildCharSearchIdentifier(license string) string {
	return "%" + license
}

func (s *Server) getCharacter(
	ctx context.Context,
	charId int32,
) (*users.User, *jobs.JobProps, string, error) {
	tUsers := tables.User().AS("user")
	tLogo := table.FivenetFiles.AS("logo_file")
	tAvatar := table.FivenetFiles.AS("profile_picture")

	stmt := tUsers.
		SELECT(
			tUsers.ID,
			tUsers.ID.AS("user.user_id"),
			tUsers.Identifier,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Dateofbirth,
			tUserProps.AvatarFileID.AS("user.profile_picture_file_id"),
			tAvatar.FilePath.AS("user.profile_picture"),
			tUsers.Group.AS("group"),
			tJobProps.Job,
			tJobProps.DeletedAt,
			tJobProps.LivemapMarkerColor,
			tJobProps.QuickButtons,
			tJobProps.RadioFrequency,
			tJobProps.LogoFileID,
			tLogo.ID,
			tLogo.FilePath,
		).
		FROM(
			tUsers.
				LEFT_JOIN(tJobProps,
					tJobProps.Job.EQ(tUsers.Job),
				).
				LEFT_JOIN(tLogo,
					tLogo.ID.EQ(tJobProps.LogoFileID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUsers.ID),
				).
				LEFT_JOIN(tAvatar,
					tAvatar.ID.EQ(tUserProps.AvatarFileID),
				),
		).
		WHERE(
			tUsers.ID.EQ(mysql.Int32(charId)),
		).
		LIMIT(1)

	var dest struct {
		*users.User

		Group    string
		JobProps *jobs.JobProps
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

	s.enricher.EnrichJobInfo(dest.User)

	return dest.User, dest.JobProps, dest.Group, nil
}

func (s *Server) ChooseCharacter(
	ctx context.Context,
	req *pbauth.ChooseCharacterRequest,
) (*pbauth.ChooseCharacterResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.auth.char_id", req.GetCharId()})

	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, errorsgrpcauth.ErrInvalidToken)
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, errswrap.NewError(err, errorsgrpcauth.ErrInvalidToken)
	}

	// Load account data for token creation
	account, err := s.getAccountFromClaims(ctx, claims)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrNoCharFound)
	}

	char, jProps, userGroup, err := s.getCharacter(ctx, req.GetCharId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrNoCharFound)
	}

	// Make sure the user isn't sending us a different char ID than their own
	if !strings.HasSuffix(char.GetIdentifier(), claims.Subject) {
		s.logger.Error(
			"user sent bad char!",
			zap.String("expected", char.GetIdentifier()),
			zap.String("current", claims.Subject),
		)
		return nil, errorsauth.ErrUnableToChooseChar
	}

	isSuperuser := slices.Contains(s.superuserGroups, userGroup) ||
		slices.Contains(s.superuserUsers, claims.Subject)

	if err := s.ui.RefreshUserInfo(ctx, char.GetUserId(), claims.AccID); err != nil {
		s.logger.Error(
			"failed to refresh user info",
			zap.Error(err),
			zap.Int32("user_id", char.GetUserId()),
		)
		return nil, errswrap.NewError(err, errorsauth.ErrUnableToChooseChar)
	}

	// If char lock is active, make sure that the user is choosing the active char
	if !isSuperuser &&
		s.appCfg.Get().Auth.GetLastCharLock() &&
		account.LastChar != nil &&
		*account.LastChar != req.GetCharId() {
		return nil, errorsgrpcauth.ErrCharLock
	}

	// Centralized superuser/override logic
	if jPropsOverride, err := s.handleSuperuserOverride(ctx, account, char, claims, isSuperuser); err != nil {
		return nil, err
	} else if jPropsOverride != nil {
		jProps = jPropsOverride
	}

	newToken, newClaims, err := s.createTokenFromAccountAndChar(account, char)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}

	ps, attrs, err := s.listUserPerms(ctx, account, char, isSuperuser)
	if err != nil {
		return nil, err
	}

	if len(ps) == 0 ||
		(!isSuperuser && !slices.ContainsFunc(ps, func(p *permissions.Permission) bool {
			return p.GetCategory() == string(permsauth.AuthServicePerm) && p.GetName() == string(permsauth.AuthServiceChooseCharacterPerm)
		})) {
		return nil, errorsauth.ErrUnableToChooseChar
	}

	defer s.aud.Log(&audit.AuditEntry{
		Service: pbauth.AuthService_ServiceDesc.ServiceName,
		Method:  "ChooseCharacter",
		UserId:  char.GetUserId(),
		UserJob: char.GetJob(),
		State:   audit.EventType_EVENT_TYPE_VIEWED,
	}, char.UserShort())

	if err := s.setTokenCookie(ctx, newToken); err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}

	return &pbauth.ChooseCharacterResponse{
		Expires:     timestamp.New(newClaims.ExpiresAt.Time),
		Username:    *account.Username,
		JobProps:    jProps,
		Char:        char,
		Permissions: ps,
		Attributes:  attrs,
	}, nil
}

func (s *Server) listUserPerms(
	ctx context.Context,
	account *model.FivenetAccounts,
	char *users.User,
	isSuperuser bool,
) ([]*permissions.Permission, []*permissions.RoleAttribute, error) {
	// Load permissions of user
	userPs, err := s.ps.GetPermissionsOfUser(&userinfo.UserInfo{
		UserId:   char.GetUserId(),
		Job:      char.GetJob(),
		JobGrade: char.GetJobGrade(),
	})
	if err != nil {
		return nil, nil, errswrap.NewError(err, errorsauth.ErrUnableToChooseChar)
	}

	if isSuperuser {
		userPs = append(userPs, auth.PermCanBeSuperuser)

		if account.Superuser != nil && *account.Superuser {
			userPs = append(userPs, auth.PermSuperuser)
		}
	}

	attrs, err := s.ps.GetEffectiveRoleAttributes(ctx, char.GetJob(), char.GetJobGrade())
	if err != nil {
		return nil, nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}

	return userPs, attrs, nil
}

func (s *Server) SetSuperuserMode(
	ctx context.Context,
	req *pbauth.SetSuperuserModeRequest,
) (*pbauth.SetSuperuserModeResponse, error) {
	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, errorsgrpcauth.ErrInvalidToken
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, errswrap.NewError(err, errorsgrpcauth.ErrInvalidToken)
	}

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	logging.InjectFields(ctx, logging.Fields{"fivenet.auth.superuser", req.GetSuperuser()})
	if req.Job != nil {
		logging.InjectFields(ctx, logging.Fields{"fivenet.auth.superuser.job", req.GetJob()})
	}

	if !userInfo.GetCanBeSuperuser() {
		if !userInfo.GetSuperuser() {
			return nil, errorsauth.ErrNotSuperuser
		}

		req.Superuser = false
		req.Job = nil
	}

	// Set user's job as requested job when superuser mode is turned on
	if req.Job == nil {
		req.Job = &userInfo.Job
	}

	char, _, _, err := s.getCharacter(ctx, claims.CharID)
	if err != nil {
		return nil, errswrap.NewError(
			fmt.Errorf("failed to get char %d. %w", claims.CharID, err),
			errorsauth.ErrNoCharFound,
		)
	}

	var jobProps *jobs.JobProps
	var ps []*permissions.Permission
	var attrs []*permissions.RoleAttribute

	// Reset override job when switching off superuser mode using centralized helper
	if !req.GetSuperuser() {
		// Fetch the account for the current user
		account, err := s.getAccountFromDB(ctx, tAccounts.Username.EQ(mysql.String(claims.Username)))
		if err != nil {
			return nil, errswrap.NewError(
				fmt.Errorf("failed to get account from db. %w", err),
				errorsauth.ErrGenericLogin,
			)
		}

		jPropsOverride, err := s.handleSuperuserOverride(ctx, account, char, claims, false)
		if err != nil {
			return nil, err
		}
		if jPropsOverride != nil {
			jobProps = jPropsOverride
		}

		userInfo.Job = char.GetJob()
		userInfo.JobGrade = char.GetJobGrade()
		userInfo.OverrideJob = nil
		userInfo.OverrideJobGrade = nil

		not := false
		ps, attrs, err = s.listUserPerms(ctx, &model.FivenetAccounts{
			Superuser: &not,
		}, char, true)
		if err != nil {
			return nil, fmt.Errorf("failed to get user perms. %w", err)
		}
	} else {
		// Only set job if requested
		job, jobGrade, jProps, err := s.getJobWithProps(ctx, req.GetJob())
		if err != nil {
			return nil, errswrap.NewError(fmt.Errorf("failed to get job props for '%s' job. %w", req.GetJob(), err), errorsauth.ErrGenericLogin)
		}
		jobProps = jProps

		userInfo.Job = job.GetName()
		userInfo.JobGrade = jobGrade
		userInfo.OverrideJob = &job.Name
		userInfo.OverrideJobGrade = &jobGrade

		char.Job = job.GetName()
		char.JobGrade = jobGrade
		s.enricher.EnrichJobInfo(char)

		ps = []*permissions.Permission{auth.PermCanBeSuperuser, auth.PermSuperuser}
	}

	//nolint:protogetter // The values are needed as pointers
	if err := s.ui.SetUserInfo(ctx, claims.AccID, claims.CharID, req.GetSuperuser(), userInfo.OverrideJob, userInfo.OverrideJobGrade); err != nil {
		return nil, errswrap.NewError(
			fmt.Errorf("failed to set user info. %w", err),
			errorsauth.ErrGenericLogin,
		)
	}

	userInfo.Superuser = req.GetSuperuser()

	// Load account data for token creation
	account, err := s.getAccountFromDB(ctx, tAccounts.Username.EQ(mysql.String(claims.Username)))
	if err != nil {
		return nil, errswrap.NewError(
			fmt.Errorf("failed to get account from db. %w", err),
			errorsauth.ErrGenericLogin,
		)
	}

	newToken, newClaims, err := s.createTokenFromAccountAndChar(account, char)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}

	if err := s.setTokenCookie(ctx, newToken); err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}

	return &pbauth.SetSuperuserModeResponse{
		Expires:     timestamp.New(newClaims.ExpiresAt.Time),
		JobProps:    jobProps,
		Char:        char,
		Permissions: ps,
		Attributes:  attrs,
	}, nil
}

func (s *Server) getJobWithProps(
	ctx context.Context,
	jobName string,
) (*jobs.Job, int32, *jobs.JobProps, error) {
	tJobs := tables.Jobs().AS("job")
	tJobsGrades := tables.JobsGrades()
	tFiles := table.FivenetFiles.AS("logo_file")

	stmt := tJobs.
		SELECT(
			tJobs.Name,
			tJobs.Label,
			tJobsGrades.Grade.AS("job_grade"),
			tJobProps.Job,
			tJobProps.UpdatedAt,
			tJobProps.LivemapMarkerColor,
			tJobProps.RadioFrequency,
			tJobProps.QuickButtons,
			tJobProps.LogoFileID,
			tFiles.ID,
			tFiles.FilePath,
		).
		FROM(
			tJobs.
				INNER_JOIN(tJobsGrades,
					tJobsGrades.JobName.EQ(tJobs.Name),
				).
				LEFT_JOIN(tJobProps,
					tJobProps.Job.EQ(tJobs.Name),
				).
				LEFT_JOIN(tFiles,
					tFiles.ID.EQ(tJobProps.LogoFileID),
				),
		).
		WHERE(
			tJobs.Name.EQ(mysql.String(jobName)),
		).
		ORDER_BY(tJobsGrades.Grade.DESC()).
		LIMIT(1)

	var dest struct {
		Job      *jobs.Job
		JobGrade int32 `alias:"job_grade"`
		JobProps *jobs.JobProps
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, 0, nil, err
	}

	if dest.JobProps != nil {
		s.enricher.EnrichJobName(dest.JobProps)
	}

	return dest.Job, dest.JobGrade, dest.JobProps, nil
}
