package auth

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
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
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	errorsauth "github.com/fivenet-app/fivenet/v2026/services/auth/errors"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
)

func (s *Server) Login(
	ctx context.Context,
	req *pbauth.LoginRequest,
) (*pbauth.LoginResponse, error) {
	req.Username = normalizeUsername(req.GetUsername())

	logging.InjectFields(ctx, logging.Fields{"fivenet.auth.username", req.GetUsername()})

	account, err := s.store.GetLoginAccountByUsername(ctx, req.GetUsername())
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

	acc, err := s.store.GetPasswordResetAccountByRegToken(ctx, req.GetRegToken())
	if err != nil {
		s.logger.Error(
			"failed to get account from database by registration token",
			zap.Error(err),
			zap.String("reg_token", req.GetRegToken()),
		)
		return nil, errswrap.NewError(
			err,
			errorsauth.ErrAccountCreateFailed(map[string]any{"code": "404"}),
		)
	}

	if acc.Username != nil || acc.Password != nil {
		return nil, errorsauth.ErrAccountExistsFailed
	}

	req.Username = normalizeUsername(req.GetUsername())

	hashedPassword, err := hashPassword(req.GetPassword())
	if err != nil {
		return nil, errswrap.NewError(
			err,
			errorsauth.ErrAccountCreateFailed(map[string]any{"code": "500"}),
		)
	}

	if err := s.store.ActivateAccount(
		ctx,
		acc.ID,
		req.GetRegToken(),
		req.GetUsername(),
		hashedPassword,
	); err != nil {
		if dbutils.IsDuplicateError(err) {
			return nil, errswrap.NewError(err, errorsauth.ErrAccountDuplicate)
		}

		s.logger.Error(
			"failed to update account in database during account creation",
			zap.Error(err),
		)
		return nil, errswrap.NewError(
			err,
			errorsauth.ErrAccountCreateFailed(map[string]any{"code": "409"}),
		)
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
		return nil, errswrap.NewError(
			err,
			errorsauth.ErrChangePassword(map[string]any{"code": "401"}),
		)
	}

	acc, err := s.store.GetAccountByIDAndUsername(ctx, claims.AccID, claims.Username, true)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsgrpcauth.ErrNoUserInfo)
		}
		return nil, errswrap.NewError(
			err,
			errorsauth.ErrChangePassword(map[string]any{"code": "500"}),
		)
	}

	// Account has no password set
	if acc.Password == nil {
		return nil, errswrap.NewError(err, errorsauth.ErrNoAccount)
	}

	// Password check logic
	if err := checkPassword(*acc.Password, req.GetCurrentPassword()); err != nil {
		return nil, errswrap.NewError(
			err,
			errorsauth.ErrChangePassword(map[string]any{"code": "401"}),
		)
	}

	hashedPassword, err := hashPassword(req.GetNewPassword())
	if err != nil {
		return nil, errswrap.NewError(
			err,
			errorsauth.ErrChangePassword(map[string]any{"code": "500"}),
		)
	}

	pass := hashedPassword
	acc.Password = &pass

	if err := s.store.UpdatePassword(ctx, acc.ID, pass); err != nil {
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

	// Make sure current and new username aren't the same
	currentUsername := req.GetCurrentUsername()
	newUsername := normalizeUsername(req.GetNewUsername())
	if strings.EqualFold(currentUsername, newUsername) {
		return nil, errorsauth.ErrBadUsername
	}

	// Retrieve account from db using account ID and username
	acc, err := s.store.GetAccountByIDAndUsername(ctx, claims.AccID, claims.Username, true)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsgrpcauth.ErrNoUserInfo)
		}
		return nil, errswrap.NewError(err, errorsauth.ErrChangeUsername)
	}

	// No username nor password set on account, fail
	if acc.Username == nil || acc.Password == nil {
		return nil, errorsauth.ErrNoAccount
	}

	// Make sure current username matches the sent current username
	if !strings.EqualFold(*acc.Username, currentUsername) {
		return nil, errorsauth.ErrBadUsername
	}

	// New username is same as current username.. just return here.
	resp := &pbauth.ChangeUsernameResponse{}
	if *acc.Username == newUsername {
		return nil, errorsauth.ErrBadUsername
	}

	// If there is an account with the new username, fail
	newAcc, err := s.store.GetAccountByUsername(ctx, newUsername, false)
	if err != nil && !errors.Is(err, qrm.ErrNoRows) {
		// Other database error
		return nil, errswrap.NewError(err, errorsauth.ErrBadUsername)
	}
	// An account with the requested username was found, fail
	if newAcc != nil {
		return nil, errorsauth.ErrBadUsername
	}

	acc.Username = &newUsername

	if err := s.store.UpdateUsername(ctx, acc.ID, newUsername); err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrUpdateAccount)
	}

	// Destroy session
	if err := s.destroyCookies(ctx); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) ForgotPassword(
	ctx context.Context,
	req *pbauth.ForgotPasswordRequest,
) (*pbauth.ForgotPasswordResponse, error) {
	acc, err := s.store.GetAccountByRegToken(ctx, req.GetRegToken(), true)
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

	if err := s.store.ForgotPassword(ctx, acc.ID, pass); err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrForgotPassword)
	}

	// Destroy session
	if err := s.destroyCookies(ctx); err != nil {
		return nil, err
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
	acc, err := s.store.GetAccountByIDAndUsername(ctx, accClaims.AccID, accClaims.Username, false)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsgrpcauth.ErrNoUserInfo)
		}
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}
	if acc.ID <= 0 {
		return nil, errorsauth.ErrGenericLogin
	}

	resp := &pbauth.GetCharactersResponse{}
	resp.Chars, err = s.store.ListCharacters(ctx, accClaims.AccID, acc.License)
	if err != nil {
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
	acc, err := s.store.GetAccountByIDAndUsername(
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
			return p.GetNamespace() == string(permsauth.Namespace) &&
				p.GetService() == string(permsauth.AuthServicePerm) &&
				p.GetName() == string(permsauth.AuthServiceChooseCharacterPerm)
		})) {
		return nil, errorsauth.ErrUnableToChooseChar
	}

	grpc_audit.SetUser(ctx, char.UserId, char.Job)

	// Ensure can be superuser is set on account claims
	accClaims := auth.MapAccountToClaims(account)
	if canBeSuperuser && currentAccClaims.CanBeSuperuser {
		accClaims.CanBeSuperuser = currentAccClaims.CanBeSuperuser
	}

	// Ensure superuser is set on user claims
	userClaims := auth.MapUserToClaims(account.Id, char)
	if canBeSuperuser && currentUserClaims != nil && currentUserClaims.Superuser != nil &&
		*currentUserClaims.Superuser {
		userClaims.Superuser = currentUserClaims.Superuser
	}
	// Based on impersonation, (re-)set the job and grade of the user
	if currentUserClaims != nil && currentUserClaims.OriginalJob != nil {
		if userClaims.Superuser != nil && *userClaims.Superuser {
			// Set "original" job of user
			if currentUserClaims.Job != nil {
				userClaims.Job = currentUserClaims.Job
				char.Job = *currentUserClaims.Job
			} else {
				userClaims.Job = nil
			}
			if currentUserClaims.JobGrade != nil {
				userClaims.JobGrade = currentUserClaims.JobGrade
				char.JobGrade = *currentUserClaims.JobGrade
			} else {
				userClaims.JobGrade = nil
			}

			userClaims.OriginalJob = currentUserClaims.OriginalJob

			s.enricher.EnrichJobInfo(char)
		}
	} else {
		userClaims.OriginalJob = nil
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
		userPs = append(userPs, perms.PermCanBeSuperuser)

		if isSuperuserActive {
			userPs = append(userPs, perms.PermSuperuser)
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
	acc, err := s.store.GetAccountByIDAndUsername(ctx, accClaims.AccID, accClaims.Username, false)
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

	userClaims := auth.MapUserToClaims(acc.ID, char)
	// Carry over superuser status
	if accClaims.CanBeSuperuser {
		superuser := userInfo.GetSuperuser()
		userClaims.Superuser = &superuser
	}

	if imp {
		// Set original job of user
		userClaims.OriginalJob = &authclaims.UserJobInfo{
			Job:      char.GetJob(),
			JobGrade: char.GetJobGrade(),
		}

		// Set user job to requested impersonation job + grade
		char.Job = job.GetName()
		char.JobGrade = grade.GetGrade()
	} else {
		userClaims.OriginalJob = nil
	}

	canBeSuperuser := account.Groups.ContainsAnyGroup(s.superuserGroups) ||
		slices.Contains(s.superuserUsers, accClaims.Subject)

	ps, attrs, err := s.listUserPerms(ctx, char, canBeSuperuser, userInfo.Superuser)
	if err != nil {
		return nil, err
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

	origJob := char.GetJob()
	origGrade := char.GetJobGrade()

	var jobProps *jobsprops.JobProps
	var ps []*permissionspermissions.Permission
	var attrs []*permissionsattributes.RoleAttribute

	// Reset job when switching off superuser mode
	if req.GetSuperuser() {
		// Only set job if requested
		job, jobGrade, jProps, err := s.getJobWithProps(ctx, req.GetJob())
		if err != nil {
			return nil, errswrap.NewError(
				fmt.Errorf("failed to get job props for %q job. %w", req.GetJob(), err),
				errorsauth.ErrGenericLogin,
			)
		}
		jobProps = jProps

		userInfo.Superuser = true

		char.Job = job.GetName()
		char.JobGrade = jobGrade
		s.enricher.EnrichJobInfo(char)

		ps = []*permissionspermissions.Permission{
			perms.PermCanBeSuperuser,
			perms.PermSuperuser,
		}
	} else {
		ps, attrs, err = s.listUserPerms(ctx, char, true, false)
		if err != nil {
			return nil, fmt.Errorf("failed to get user perms. %w", err)
		}

		userInfo.Superuser = false
	}

	// Load account data for token creation
	account, err := s.store.GetAccountByIDAndUsername(
		ctx,
		accClaims.AccID,
		accClaims.Username,
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

	if req.GetSuperuser() {
		// Set original job of user
		userClaims.OriginalJob = &authclaims.UserJobInfo{
			Job:      origJob,
			JobGrade: origGrade,
		}
	} else {
		userClaims.OriginalJob = nil
	}

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
