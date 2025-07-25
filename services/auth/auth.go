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
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsauth "github.com/fivenet-app/fivenet/v2025/services/auth/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var (
	tAccounts  = table.FivenetAccounts
	tUserProps = table.FivenetUserProps.AS("user_props")
	tJobProps  = table.FivenetJobProps.AS("job_props")
)

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

func (s *Server) Login(ctx context.Context, req *pbauth.LoginRequest) (*pbauth.LoginResponse, error) {
	req.Username = normalizeUsername(req.Username)

	trace.SpanFromContext(ctx).SetAttributes(attribute.String("fivenet.auth.username", req.Username))

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

	if err := checkPassword(*account.Password, req.Password); err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrInvalidLogin)
	}

	token, claims, err := s.createTokenFromAccountAndChar(account, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}

	var chooseCharResp *pbauth.ChooseCharacterResponse
	if s.appCfg.Get().Auth.LastCharLock && account.LastChar != nil {
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

func (s *Server) Logout(ctx context.Context, req *pbauth.LogoutRequest) (*pbauth.LogoutResponse, error) {
	s.destroyTokenCookie(ctx)

	return &pbauth.LogoutResponse{
		Success: true,
	}, nil
}

func (s *Server) CreateAccount(ctx context.Context, req *pbauth.CreateAccountRequest) (*pbauth.CreateAccountResponse, error) {
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

	req.Username = normalizeUsername(req.Username)

	hashedPassword, err := hashPassword(req.Password)
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
			tAccounts.Password.SET(jet.String(hashedPassword)),
			tAccounts.RegToken.SET(jet.StringExp(jet.NULL)),
		).
		WHERE(jet.AND(
			tAccounts.ID.EQ(jet.Uint64(acc.ID)),
			tAccounts.RegToken.EQ(jet.String(req.RegToken)),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		if dbutils.IsDuplicateError(err) {
			return nil, errswrap.NewError(err, errorsauth.ErrAccountDuplicate)
		}

		s.logger.Error("failed to update account in database during account creation", zap.Error(err))
		return nil, errswrap.NewError(err, errorsauth.ErrAccountCreateFailed)
	}

	return &pbauth.CreateAccountResponse{
		AccountId: acc.ID,
	}, nil
}

func (s *Server) ChangePassword(ctx context.Context, req *pbauth.ChangePasswordRequest) (*pbauth.ChangePasswordResponse, error) {
	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, auth.ErrInvalidToken)
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

	if err := checkPassword(*acc.Password, req.Current); err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrChangePassword)
	}

	hashedPassword, err := hashPassword(req.New)
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

	if err := s.setTokenCookie(ctx, newToken); err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}

	return &pbauth.ChangePasswordResponse{
		Expires: timestamp.New(newClaims.ExpiresAt.Time),
	}, nil
}

func (s *Server) ChangeUsername(ctx context.Context, req *pbauth.ChangeUsernameRequest) (*pbauth.ChangeUsernameResponse, error) {
	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, auth.ErrInvalidToken)
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
	if !strings.EqualFold(*acc.Username, req.Current) {
		return nil, errorsauth.ErrBadUsername
	}

	req.New = normalizeUsername(req.New)
	username := req.New

	// New username is same as current username.. just return here.
	resp := &pbauth.ChangeUsernameResponse{}
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

func (s *Server) ForgotPassword(ctx context.Context, req *pbauth.ForgotPasswordRequest) (*pbauth.ForgotPasswordResponse, error) {
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

	hashedPassword, err := hashPassword(req.New)
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
			tAccounts.Password.SET(jet.String(pass)),
			tAccounts.RegToken.SET(jet.StringExp(jet.NULL)),
		).
		WHERE(
			tAccounts.ID.EQ(jet.Uint64(acc.ID)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrForgotPassword)
	}

	return &pbauth.ForgotPasswordResponse{}, nil
}

func (s *Server) GetCharacters(ctx context.Context, req *pbauth.GetCharactersRequest) (*pbauth.GetCharactersResponse, error) {
	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, auth.ErrInvalidToken)
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
	tAvatar := table.FivenetFiles.AS("avatar")

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
				tUserProps.AvatarFileID.AS("user.avatar_file_id"),
				tAvatar.FilePath.AS("user.avatar"),
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
			tUsers.Identifier.LIKE(jet.String(buildCharSearchIdentifier(claims.Subject))),
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
	if s.appCfg.Get().Auth.LastCharLock && acc.LastChar != nil {
		for i := range resp.Chars {
			s.enricher.EnrichJobInfo(resp.Chars[i].Char)

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
		for i := range resp.Chars {
			s.enricher.EnrichJobInfo(resp.Chars[i].Char)

			resp.Chars[i].Available = true
		}
	}

	return resp, nil
}

func buildCharSearchIdentifier(license string) string {
	return "%" + license
}

func (s *Server) getCharacter(ctx context.Context, charId int32) (*users.User, *jobs.JobProps, string, error) {
	tUsers := tables.User().AS("user")
	tLogo := table.FivenetFiles.AS("logo_file")
	tAvatar := table.FivenetFiles.AS("avatar")

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
			tUserProps.AvatarFileID.AS("user.avatar_file_id"),
			tAvatar.FilePath.AS("user.avatar"),
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
			tUsers.ID.EQ(jet.Int32(charId)),
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

func (s *Server) ChooseCharacter(ctx context.Context, req *pbauth.ChooseCharacterRequest) (*pbauth.ChooseCharacterResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.auth.char_id", int64(req.CharId)))

	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, auth.ErrInvalidToken)
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, errswrap.NewError(err, auth.ErrInvalidToken)
	}

	// Load account data for token creation
	account, err := s.getAccountFromClaims(ctx, claims)
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

	isSuperuser := slices.Contains(s.superuserGroups, userGroup) || slices.Contains(s.superuserUsers, claims.Subject)

	// If char lock is active, make sure that the user is choosing the active char
	if !isSuperuser &&
		s.appCfg.Get().Auth.LastCharLock &&
		account.LastChar != nil &&
		*account.LastChar != req.CharId {
		return nil, errorsauth.ErrCharLock
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

	if len(ps) == 0 || (!isSuperuser && !slices.ContainsFunc(ps, func(p *permissions.Permission) bool {
		return p.Category == string(permsauth.AuthServicePerm) && p.Name == string(permsauth.AuthServiceChooseCharacterPerm)
	})) {
		return nil, errorsauth.ErrUnableToChooseChar
	}

	defer s.aud.Log(&audit.AuditEntry{
		Service: pbauth.AuthService_ServiceDesc.ServiceName,
		Method:  "ChooseCharacter",
		UserId:  char.UserId,
		UserJob: char.Job,
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

func (s *Server) listUserPerms(ctx context.Context, account *model.FivenetAccounts, char *users.User, isSuperuser bool) ([]*permissions.Permission, []*permissions.RoleAttribute, error) {
	// Load permissions of user
	userPs, err := s.ps.GetPermissionsOfUser(&userinfo.UserInfo{
		UserId:   char.UserId,
		Job:      char.Job,
		JobGrade: char.JobGrade,
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

	attrs, err := s.ps.GetEffectiveRoleAttributes(ctx, char.Job, char.JobGrade)
	if err != nil {
		return nil, nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}

	return userPs, attrs, nil
}

func (s *Server) SetSuperuserMode(ctx context.Context, req *pbauth.SetSuperuserModeRequest) (*pbauth.SetSuperuserModeResponse, error) {
	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, auth.ErrInvalidToken
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, errswrap.NewError(err, auth.ErrInvalidToken)
	}

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	trace.SpanFromContext(ctx).SetAttributes(attribute.Bool("fivenet.auth.superuser", req.Superuser))
	if req.Job != nil {
		trace.SpanFromContext(ctx).SetAttributes(attribute.String("fivenet.auth.superuser.job", *req.Job))
	}

	if !userInfo.CanBeSuperuser {
		return nil, errswrap.NewError(err, errorsauth.ErrNoCharFound)
	}

	// Set user's job as requested job when superuser mode is turned on
	if req.Job == nil {
		req.Job = &userInfo.Job
	}

	char, _, _, err := s.getCharacter(ctx, claims.CharID)
	if err != nil {
		return nil, errswrap.NewError(fmt.Errorf("failed to get char %d. %w", claims.CharID, err), errorsauth.ErrNoCharFound)
	}

	var jobProps *jobs.JobProps
	var ps []*permissions.Permission
	var attrs []*permissions.RoleAttribute

	// Reset override job when switching off superuser mode using centralized helper
	if !req.Superuser {
		// Fetch the account for the current user
		account, err := s.getAccountFromDB(ctx, tAccounts.Username.EQ(jet.String(claims.Username)))
		if err != nil {
			return nil, errswrap.NewError(fmt.Errorf("failed to get account from db. %w", err), errorsauth.ErrGenericLogin)
		}

		jPropsOverride, err := s.handleSuperuserOverride(ctx, account, char, claims, false)
		if err != nil {
			return nil, err
		}
		if jPropsOverride != nil {
			jobProps = jPropsOverride
		}

		userInfo.Job = char.Job
		userInfo.JobGrade = char.JobGrade
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
		job, jobGrade, jProps, err := s.getJobWithProps(ctx, *req.Job)
		if err != nil {
			return nil, errswrap.NewError(fmt.Errorf("failed to get job props for '%s' job. %w", *req.Job, err), errorsauth.ErrGenericLogin)
		}
		jobProps = jProps

		userInfo.Job = job.Name
		userInfo.JobGrade = jobGrade
		userInfo.OverrideJob = &job.Name
		userInfo.OverrideJobGrade = &jobGrade

		char.Job = job.Name
		char.JobGrade = jobGrade
		s.enricher.EnrichJobInfo(char)

		ps = []*permissions.Permission{auth.PermCanBeSuperuser, auth.PermSuperuser}
	}

	if err := s.ui.SetUserInfo(ctx, claims.AccID, req.Superuser, userInfo.OverrideJob, userInfo.OverrideJobGrade); err != nil {
		return nil, errswrap.NewError(fmt.Errorf("failed to set user info. %w", err), errorsauth.ErrGenericLogin)
	}

	userInfo.Superuser = req.Superuser

	// Load account data for token creation
	account, err := s.getAccountFromDB(ctx, tAccounts.Username.EQ(jet.String(claims.Username)))
	if err != nil {
		return nil, errswrap.NewError(fmt.Errorf("failed to get account from db. %w", err), errorsauth.ErrGenericLogin)
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

func (s *Server) getJobWithProps(ctx context.Context, jobName string) (*jobs.Job, int32, *jobs.JobProps, error) {
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
			tJobs.Name.EQ(jet.String(jobName)),
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
