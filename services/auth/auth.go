package auth

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	jobsprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/props"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	permissionspermissions "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/permissions"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	users "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	pbauth "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/auth"
	permsauth "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/auth/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	authclaims "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth/claims"
	errorsgrpcauth "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth/errors"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorsauth "github.com/fivenet-app/fivenet/v2026/services/auth/errors"
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

func (s *Server) getAccountFromDB(
	ctx context.Context,
	condition mysql.BoolExpression,
	withPass bool,
) (*model.FivenetAccounts, error) {
	columns := mysql.ProjectionList{
		tAccounts.ID,
		tAccounts.CreatedAt,
		tAccounts.UpdatedAt,
		tAccounts.Enabled,
		tAccounts.Username,
		tAccounts.License,
		tAccounts.RegToken,
		tAccounts.Groups,
		tAccounts.LastChar,
	}
	if withPass {
		columns = append(columns, tAccounts.Password)
	}

	stmt := tAccounts.
		SELECT(
			columns[0],
			columns[1:]...,
		).
		FROM(tAccounts).
		WHERE(mysql.AND(
			tAccounts.Enabled.IS_TRUE(),
			condition,
		)).
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
	), true)
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

	accClaims := auth.MapAccountToClaims(accounts.ConvertFromModelAcc(account))

	// FIXME move this to the sync API user logic
	if err := s.linkCharsToAccount(ctx, account.License, account.ID); err != nil {
		s.logger.Error(
			"failed to link chars to account",
			zap.Error(err),
		)
	}

	var chooseCharResp *pbauth.ChooseCharacterResponse
	if s.appCfg.Get().Auth.GetLastCharLock() && account.LastChar != nil {
		token, err := s.tm.FromAccClaims(accClaims)
		if err != nil {
			return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
		}

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
		if err := s.setCookies(ctx, accClaims); err != nil {
			return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
		}
	}

	return &pbauth.LoginResponse{
		Expires:   timestamp.New(accClaims.ExpiresAt.Time),
		AccountId: account.ID,
		Char:      chooseCharResp,
	}, nil
}

func (s *Server) linkCharsToAccount(
	ctx context.Context,
	identifier string,
	accId int64,
) error {
	tUsers := table.FivenetUser

	stmt := tUsers.
		UPDATE(
			tUsers.AccountID,
		).
		SET(
			accId,
		).
		WHERE(
			tUsers.Identifier.LIKE(mysql.String(fmt.Sprintf("%%%s", identifier))),
		).
		LIMIT(15)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return err
	}

	return nil
}

func (s *Server) Logout(
	ctx context.Context,
	req *pbauth.LogoutRequest,
) (*pbauth.LogoutResponse, error) {
	// No need to audit logout actions
	grpc_audit.Skip(ctx)

	err := s.destroyCookies(ctx)
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

	acc, err := s.getAccountFromDB(
		ctx,
		tAccounts.RegToken.EQ(mysql.String(req.GetRegToken())),
		true,
	)
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
	token, err := auth.GetAccTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, errorsgrpcauth.ErrInvalidToken)
	}
	claims, err := s.tm.ParseAccToken(token)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrChangePassword)
	}

	acc, err := s.getAccountFromIDAndUsername(ctx, claims.AccID, claims.Username, true)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsgrpcauth.ErrNoUserInfo)
		}
		return nil, errswrap.NewError(err, errorsauth.ErrChangePassword)
	}

	// Account has no password set
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
		).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrUpdateAccount)
	}

	// Clear session cookies after password change
	if err := s.destroyCookies(ctx); err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}

	return &pbauth.ChangePasswordResponse{}, nil
}

func (s *Server) ChangeUsername(
	ctx context.Context,
	req *pbauth.ChangeUsernameRequest,
) (*pbauth.ChangeUsernameResponse, error) {
	token, err := auth.GetAccTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, errorsgrpcauth.ErrInvalidToken)
	}

	claims, err := s.tm.ParseAccToken(token)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrChangeUsername)
	}

	acc, err := s.getAccountFromIDAndUsername(ctx, claims.AccID, claims.Username, true)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsgrpcauth.ErrNoUserInfo)
		}
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
	newAcc, err := s.getAccountFromDB(ctx, tAccounts.Username.EQ(mysql.String(username)), false)
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
		).
		LIMIT(1)

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
	), true)
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
	accToken, err := auth.GetAccTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, errorsgrpcauth.ErrInvalidToken)
	}
	accClaims, err := s.tm.ParseAccToken(accToken)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}

	// Load account to make sure it (still) exists
	acc, err := s.getAccountFromIDAndUsername(ctx, accClaims.AccID, accClaims.Username, false)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsgrpcauth.ErrNoUserInfo)
		}
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}
	if acc.ID <= 0 {
		return nil, errorsauth.ErrGenericLogin
	}

	// Load chars from database
	tUsers := table.FivenetUser.AS("user")
	tAvatar := table.FivenetFiles.AS("profile_picture")

	stmt := tUsers.
		SELECT(
			tUsers.ID,
			dbutils.Columns{
				tUsers.ID.AS("user.user_id"),
				tUsers.AccountID,
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
		WHERE(tUsers.AccountID.EQ(mysql.Int64(accClaims.AccID))).
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
				) || slices.Contains(s.superuserUsers, accClaims.Subject) {
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

func (s *Server) getCharacter(
	ctx context.Context,
	charId int32,
) (*users.User, *jobsprops.JobProps, error) {
	tUsers := table.FivenetUser.AS("user")
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
		WHERE(tUsers.ID.EQ(mysql.Int32(charId))).
		LIMIT(1)

	var dest struct {
		*users.User

		JobProps *jobsprops.JobProps
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, nil, err
	}

	if dest.JobProps != nil {
		s.enricher.EnrichJobName(dest.JobProps)
	}

	s.enricher.EnrichJobInfo(dest.User)

	return dest.User, dest.JobProps, nil
}

func (s *Server) ChooseCharacter(
	ctx context.Context,
	req *pbauth.ChooseCharacterRequest,
) (*pbauth.ChooseCharacterResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.auth.char_id", req.GetCharId()})

	accToken, err := auth.GetAccTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, errorsgrpcauth.ErrInvalidToken)
	}
	currentAccClaims, err := s.tm.ParseAccToken(accToken)
	if err != nil {
		return nil, errswrap.NewError(err, errorsgrpcauth.ErrInvalidToken)
	}

	var currentUserClaims *authclaims.UserInfoClaims
	userToken, err := auth.GetUserTokenFromGRPCContext(ctx)
	if err == nil {
		currentUserClaims, err = s.tm.ParseUserToken(userToken)
		if err != nil {
			return nil, errswrap.NewError(err, errorsgrpcauth.ErrInvalidToken)
		}
	}

	// Load account data for token creation
	acc, err := s.getAccountFromIDAndUsername(
		ctx,
		currentAccClaims.AccID,
		currentAccClaims.Username,
		false,
	)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsgrpcauth.ErrNoUserInfo)
		}
		return nil, errswrap.NewError(err, errorsauth.ErrNoCharFound)
	}
	account := accounts.ConvertFromModelAcc(acc)

	char, jProps, err := s.getCharacter(ctx, req.GetCharId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrNoCharFound)
	}

	// Make sure the user isn't sending us a different char ID than their own
	if !strings.HasSuffix(char.GetIdentifier(), currentAccClaims.Subject) {
		s.logger.Error(
			"user sent bad char!",
			zap.String("expected", char.GetIdentifier()),
			zap.String("current", currentAccClaims.Subject),
		)
		return nil, errorsauth.ErrUnableToChooseChar
	}

	canBeSuperuser := account.Groups.ContainsAnyGroup(s.superuserGroups) ||
		slices.Contains(s.superuserUsers, currentAccClaims.Subject)

	if err := s.ui.RefreshUserInfo(ctx, char.GetUserId()); err != nil {
		s.logger.Error(
			"failed to refresh user info",
			zap.Error(err),
			zap.Int32("user_id", char.GetUserId()),
		)
		return nil, errswrap.NewError(err, errorsauth.ErrUnableToChooseChar)
	}

	// If char lock is active, make sure that the user is choosing the active char
	if !canBeSuperuser &&
		s.appCfg.Get().Auth.GetLastCharLock() &&
		account.LastChar != nil &&
		*account.LastChar != req.GetCharId() {
		return nil, errorsgrpcauth.ErrCharLock
	}

	ps, attrs, err := s.listUserPerms(
		ctx,
		char,
		canBeSuperuser,
		currentUserClaims != nil && currentUserClaims.Superuser != nil &&
			*currentUserClaims.Superuser,
	)
	if err != nil {
		return nil, err
	}

	if len(ps) == 0 ||
		(!canBeSuperuser && !slices.ContainsFunc(ps, func(p *permissionspermissions.Permission) bool {
			return p.GetCategory() == string(permsauth.AuthServicePerm) && p.GetName() == string(permsauth.AuthServiceChooseCharacterPerm)
		})) {
		return nil, errorsauth.ErrUnableToChooseChar
	}

	grpc_audit.SetUser(ctx, char.UserId, char.Job)

	accClaims := auth.MapAccountToClaims(account)
	if canBeSuperuser && currentAccClaims.CanBeSuperuser {
		accClaims.CanBeSuperuser = currentAccClaims.CanBeSuperuser
	}
	userClaims := auth.MapUserToClaims(account.Id, char)
	if canBeSuperuser && currentUserClaims != nil && currentUserClaims.Superuser != nil &&
		*currentUserClaims.Superuser {
		userClaims.Superuser = currentUserClaims.Superuser
	}

	if err := s.setCookies(ctx, accClaims); err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}

	userToken, err = s.tm.FromUserClaims(userClaims)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}

	return &pbauth.ChooseCharacterResponse{
		Token:       userToken,
		Expires:     timestamp.New(currentAccClaims.ExpiresAt.Time),
		Username:    account.Username,
		JobProps:    jProps,
		Char:        char,
		Permissions: ps,
		Attributes:  attrs,
	}, nil
}

func (s *Server) listUserPerms(
	ctx context.Context,
	char *users.User,
	canBeSuperuser bool,
	isSuperuserActive bool,
) ([]*permissionspermissions.Permission, []*permissionsattributes.RoleAttribute, error) {
	// Load permissions of user
	userPs, err := s.ps.GetPermissionsOfUser(&userinfo.UserInfo{
		UserId:   char.GetUserId(),
		Job:      char.GetJob(),
		JobGrade: char.GetJobGrade(),
	})
	if err != nil {
		return nil, nil, errswrap.NewError(err, errorsauth.ErrUnableToChooseChar)
	}

	if canBeSuperuser {
		userPs = append(userPs, auth.PermCanBeSuperuser)

		if isSuperuserActive {
			userPs = append(userPs, auth.PermSuperuser)
		}
	}

	attrs, err := s.ps.GetEffectiveRoleAttributes(ctx, char.GetJob(), char.GetJobGrade())
	if err != nil {
		return nil, nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}

	return userPs, attrs, nil
}

func (s *Server) ImpersonateJob(
	ctx context.Context,
	req *pbauth.ImpersonateJobRequest,
) (*pbauth.ImpersonateJobResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	imp := true
	if req.GetJobGrade() > userInfo.GetJobGrade() {
		return nil, errorsauth.ErrImpersonateJobInvalid
	} else if req.GetJobGrade() < 0 || req.GetJobGrade() == userInfo.GetJobGrade() {
		// Disable impersonation if job grade is the same
		imp = false
	}

	job, grade := s.enricher.GetJobGrade(userInfo.GetJob(), req.GetJobGrade())
	if job == nil || grade == nil {
		return nil, errorsauth.ErrImpersonateJobInvalid
	}

	accToken, err := auth.GetAccTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, errorsgrpcauth.ErrInvalidToken)
	}
	accClaims, err := s.tm.ParseAccToken(accToken)
	if err != nil {
		return nil, errswrap.NewError(err, errorsgrpcauth.ErrInvalidToken)
	}

	// Load user's account data for token creation
	acc, err := s.getAccountFromIDAndUsername(ctx, accClaims.AccID, accClaims.Username, false)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsgrpcauth.ErrNoUserInfo)
		}
		return nil, errswrap.NewError(err, errorsauth.ErrNoCharFound)
	}
	account := accounts.ConvertFromModelAcc(acc)

	char, _, err := s.getCharacter(ctx, userInfo.GetUserId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrNoCharFound)
	}

	canBeSuperuser := account.Groups.ContainsAnyGroup(s.superuserGroups) ||
		slices.Contains(s.superuserUsers, accClaims.Subject)

	if imp {
		char.Job = job.GetName()
		char.JobGrade = grade.GetGrade()
	}

	ps, attrs, err := s.listUserPerms(ctx, char, canBeSuperuser, userInfo.Superuser)
	if err != nil {
		return nil, err
	}

	userClaims := auth.MapUserToClaims(acc.ID, char)
	if imp {
		userClaims.Impersonate = &authclaims.UserImpersonate{
			Job:      job.GetName(),
			JobGrade: grade.GetGrade(),
		}
	} else {
		userClaims.Impersonate = nil
	}

	userToken, err := s.tm.FromUserClaims(userClaims)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}

	return &pbauth.ImpersonateJobResponse{
		Token:       userToken,
		Expires:     timestamp.New(userClaims.ExpiresAt.Time),
		Char:        char,
		Permissions: ps,
		Attributes:  attrs,
		State:       imp,
	}, nil
}

func (s *Server) SetSuperuserMode(
	ctx context.Context,
	req *pbauth.SetSuperuserModeRequest,
) (*pbauth.SetSuperuserModeResponse, error) {
	accToken, err := auth.GetAccTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, errorsgrpcauth.ErrInvalidToken
	}

	accClaims, err := s.tm.ParseAccToken(accToken)
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
		job := userInfo.GetJob()
		req.Job = &job
	}

	char, _, err := s.getCharacter(ctx, userInfo.UserId)
	if err != nil {
		return nil, errswrap.NewError(
			fmt.Errorf("failed to get char by id %d. %w", userInfo.UserId, err),
			errorsauth.ErrNoCharFound,
		)
	}

	var jobProps *jobsprops.JobProps
	var ps []*permissionspermissions.Permission
	var attrs []*permissionsattributes.RoleAttribute

	// Reset job when switching off superuser mode
	if !req.GetSuperuser() {
		userInfo.Job = char.GetJob()
		userInfo.JobGrade = char.GetJobGrade()

		ps, attrs, err = s.listUserPerms(ctx, char, true, false)
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

		char.Job = job.GetName()
		char.JobGrade = jobGrade
		s.enricher.EnrichJobInfo(char)

		ps = []*permissionspermissions.Permission{
			auth.PermCanBeSuperuser,
			auth.PermSuperuser,
		}
	}

	userInfo.Superuser = req.GetSuperuser()

	// Load account data for token creation
	account, err := s.getAccountFromDB(
		ctx,
		mysql.AND(
			tAccounts.ID.EQ(mysql.Int64(accClaims.AccID)),
			tAccounts.Username.EQ(mysql.String(accClaims.Username)),
		),
		false,
	)
	if err != nil {
		return nil, errswrap.NewError(
			fmt.Errorf("failed to get account from db. %w", err),
			errorsauth.ErrGenericLogin,
		)
	}

	accClaims = auth.MapAccountToClaims(accounts.ConvertFromModelAcc(account))
	accClaims.CanBeSuperuser = userInfo.GetCanBeSuperuser()

	userClaims := auth.MapUserToClaims(account.ID, char)
	superuser := userInfo.GetSuperuser()
	userClaims.Superuser = &superuser

	userToken, err := s.tm.FromUserClaims(userClaims)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}

	if err := s.setCookies(ctx, accClaims); err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}

	return &pbauth.SetSuperuserModeResponse{
		Token:       userToken,
		Expires:     timestamp.New(accClaims.ExpiresAt.Time),
		JobProps:    jobProps,
		Char:        char,
		Permissions: ps,
		Attributes:  attrs,
	}, nil
}

func (s *Server) getJobWithProps(
	ctx context.Context,
	jobName string,
) (*jobs.Job, int32, *jobsprops.JobProps, error) {
	tJobs := table.FivenetJobs.AS("job")
	tJobsGrades := table.FivenetJobsGrades
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
		JobProps *jobsprops.JobProps
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, 0, nil, err
		}

		// Create empty dummy job if not found in DB
		dest.Job = &jobs.Job{
			Name:  jobName,
			Label: jobName,
		}
		dest.JobProps = &jobsprops.JobProps{
			Job:      jobName,
			JobLabel: &jobName,
		}
		dest.JobProps.Default(jobName)
	}

	if dest.JobProps != nil {
		s.enricher.EnrichJobName(dest.JobProps)
	}

	return dest.Job, dest.JobGrade, dest.JobProps, nil
}
