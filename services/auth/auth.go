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
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	users "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	pbauth "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsauth "github.com/fivenet-app/fivenet/v2025/services/auth/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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
	req.Username = strings.TrimSpace(req.Username)

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
	if err := bcrypt.CompareHashAndPassword([]byte(*account.Password), []byte(req.Password)); err != nil {
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
	acc, err := s.getAccountFromDB(ctx, tAccounts.ID.EQ(jet.Uint64(claims.AccID)).
		AND(tAccounts.Username.EQ(jet.String(claims.Username))))
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}
	if acc.ID == 0 {
		return nil, errorsauth.ErrGenericLogin
	}

	// Load chars from database
	tUsers := tables.User().AS("user")
	tJobs := tables.Jobs()
	tJobsGrades := tables.JobsGrades()

	stmt := tUsers.
		SELECT(
			tUsers.ID,
			dbutils.Columns{
				tUsers.Identifier,
				tUsers.Job,
				tJobs.Label.AS("user.job_label"),
				tUsers.JobGrade,
				tJobsGrades.Label.AS("user.job_grade_label"),
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
			LEFT_JOIN(tJobsGrades,
				jet.AND(
					tJobsGrades.Grade.EQ(tUsers.JobGrade),
					tJobsGrades.JobName.EQ(tUsers.Job),
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
	tJobs := tables.Jobs()
	tJobsGrades := tables.JobsGrades()

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
			tJobsGrades.Label.AS("user.job_grade_label"),
			tJobProps.DeletedAt,
			tJobProps.LivemapMarkerColor,
			tJobProps.QuickButtons,
			tJobProps.RadioFrequency,
			tJobProps.LogoURL,
		).
		FROM(
			tUsers.
				LEFT_JOIN(tJobs,
					tJobs.Name.EQ(tUsers.Job),
				).
				LEFT_JOIN(tJobsGrades,
					jet.AND(
						tJobsGrades.Grade.EQ(tUsers.JobGrade),
						tJobsGrades.JobName.EQ(tUsers.Job),
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

	return &dest.User, dest.JobProps, dest.Group, nil
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

	isSuperuser := slices.Contains(s.superuserGroups, userGroup) || slices.Contains(s.superuserUsers, claims.Subject)

	// If char lock is active, make sure that the user is choosing the active char
	if !isSuperuser &&
		s.appCfg.Get().Auth.LastCharLock &&
		account.LastChar != nil &&
		*account.LastChar != req.CharId {
		return nil, errorsauth.ErrCharLock
	}

	// Reset override jobs when char is not a superuser but has an override set..
	if !isSuperuser &&
		((account.Superuser != nil && *account.Superuser) ||
			account.OverrideJob != nil) {
		account.OverrideJob = nil
		account.OverrideJobGrade = nil

		if err := s.ui.SetUserInfo(ctx, claims.AccID, false, account.OverrideJob, account.OverrideJobGrade); err != nil {
			return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
		}

		not := false
		account.Superuser = &not
	} else if isSuperuser &&
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

	ps, err := s.listUserPerms(account, char, isSuperuser)
	if err != nil {
		return nil, err
	}

	if len(ps) == 0 || (!isSuperuser && !slices.Contains(ps, "authservice-choosecharacter")) {
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
		Permissions: ps,
		JobProps:    jProps,
		Char:        char,
		Username:    *account.Username,
	}, nil
}

func (s *Server) listUserPerms(account *model.FivenetAccounts, char *users.User, isSuperuser bool) ([]string, error) {
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
	if isSuperuser {
		ps = append(ps, auth.PermCanBeSuperKey)

		if account.Superuser != nil && *account.Superuser {
			ps = append(ps, auth.PermSuperuserKey)
		}
	}

	attrs, err := s.ps.FlattenRoleAttributes(char.Job, char.JobGrade)
	if err != nil {
		return nil, errswrap.NewError(err, errorsauth.ErrGenericLogin)
	}
	ps = append(ps, attrs...)

	return ps, nil
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

	if !userInfo.CanBeSuper {
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
	var ps []string

	// Reset override job when switching off superuser mode
	if !req.Superuser {
		userInfo.Job = char.Job
		userInfo.JobGrade = char.JobGrade

		userInfo.OverrideJob = nil
		userInfo.OverrideJobGrade = nil

		// Send original char job props to user
		_, _, jProps, err := s.getJobWithProps(ctx, char.Job)
		if err != nil {
			return nil, errswrap.NewError(fmt.Errorf("failed to get job props for '%s' job. %w", char.Job, err), errorsauth.ErrGenericLogin)
		}
		jobProps = jProps

		not := false
		ps, err = s.listUserPerms(&model.FivenetAccounts{
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

		ps = []string{auth.PermCanBeSuperKey, auth.PermSuperuserKey}
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
	}, nil
}

func (s *Server) getJobWithProps(ctx context.Context, jobName string) (*jobs.Job, int32, *jobs.JobProps, error) {
	tJobs := tables.Jobs().AS("job")
	tJobsGrades := tables.JobsGrades()

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
			tJobProps.LogoURL,
		).
		FROM(
			tJobs.
				INNER_JOIN(tJobsGrades,
					tJobsGrades.JobName.EQ(tJobs.Name),
				).
				LEFT_JOIN(tJobProps,
					tJobProps.Job.EQ(tJobs.Name)),
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
