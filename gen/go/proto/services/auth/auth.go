package auth

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/galexrt/fivenet/gen/go/proto/resources/common"
	"github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
	users "github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/audit"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"golang.org/x/crypto/bcrypt"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	tAccounts  = table.FivenetAccounts
	tUsers     = table.Users.AS("user")
	tJobs      = table.Jobs
	tJobGrades = table.JobGrades
	tJobProps  = table.FivenetJobProps.AS("jobprops")
)

var (
	AccountCreateFailedErr = status.Error(codes.InvalidArgument, "Please check the token you have typed in!")
	InvalidLoginErr        = status.Error(codes.InvalidArgument, "Wrong username or password.")
	NoAccountErr           = status.Error(codes.InvalidArgument, "You need to use your registration token to create an account first.")
	NoCharacterFoundErr    = status.Error(codes.NotFound, "No character(s) found for your account.")
	GenericLoginErr        = status.Error(codes.Internal, "Failed to login you in, please try again.")
	UnableToChooseCharErr  = status.Error(codes.PermissionDenied, "You don't have permission to select this character!")
	UpdateAccountErr       = status.Error(codes.InvalidArgument, "Failed to update your account!")
	ChangePasswordErr      = status.Error(codes.InvalidArgument, "Failed to change your password!")
	ForgotPasswordErr      = status.Error(codes.InvalidArgument, "Failed to reset password!")
)

type Server struct {
	AuthServiceServer

	db   *sql.DB
	auth *auth.GRPCAuth
	tm   *auth.TokenMgr
	p    perms.Permissions
	c    *mstlystcdata.Enricher
	a    audit.IAuditer
	ui   userinfo.UserInfoRetriever
}

func NewServer(db *sql.DB, auth *auth.GRPCAuth, tm *auth.TokenMgr, p perms.Permissions, c *mstlystcdata.Enricher, aud audit.IAuditer, ui userinfo.UserInfoRetriever) *Server {
	return &Server{
		db:   db,
		auth: auth,
		tm:   tm,
		p:    p,
		c:    c,
		a:    aud,
		ui:   ui,
	}
}

// AuthFuncOverride is called instead of exampleAuthFunc
func (s *Server) AuthFuncOverride(ctx context.Context, fullMethod string) (context.Context, error) {
	// Skip authentication for the anon accessible endpoints
	if fullMethod == "/services.auth.AuthService/CreateAccount" ||
		fullMethod == "/services.auth.AuthService/Login" ||
		fullMethod == "/services.auth.AuthService/Logout" ||
		fullMethod == "/services.auth.AuthService/ForgotPassword" {
		return ctx, nil
	}

	if fullMethod == "/services.auth.AuthService/SetJob" {
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
			tAccounts.AllColumns,
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
		if errors.Is(qrm.ErrNoRows, err) {
			return nil, InvalidLoginErr
		}

		return nil, err
	}

	// No password set
	if account.Password == nil {
		return nil, NoAccountErr
	}

	// Password check logic
	if err := bcrypt.CompareHashAndPassword([]byte(*account.Password), []byte(req.Password)); err != nil {
		return nil, InvalidLoginErr
	}

	token, claims, err := s.createTokenFromAccountAndChar(account, nil)
	if err != nil {
		return nil, InvalidLoginErr
	}

	return &LoginResponse{
		Token:   token,
		Expires: timestamp.New(claims.ExpiresAt.Time),
	}, nil
}

func (s *Server) Logout(ctx context.Context, req *LogoutRequest) (*LogoutResponse, error) {
	return &LogoutResponse{
		Success: true,
	}, nil
}

func (s *Server) CreateAccount(ctx context.Context, req *CreateAccountRequest) (*CreateAccountResponse, error) {
	acc, err := s.getAccountFromDB(ctx, jet.AND(
		tAccounts.RegToken.EQ(jet.String(req.RegToken)),
		tAccounts.Username.IS_NULL(),
		tAccounts.Password.IS_NULL(),
	))
	if err != nil {
		return nil, AccountCreateFailedErr
	}

	req.Username = strings.TrimSpace(req.Username)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		return nil, AccountCreateFailedErr
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
		return nil, AccountCreateFailedErr
	}

	return &CreateAccountResponse{}, nil
}

func (s *Server) ChangePassword(ctx context.Context, req *ChangePasswordRequest) (*ChangePasswordResponse, error) {
	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, auth.InvalidTokenErr
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, ChangePasswordErr
	}

	acc, err := s.getAccountFromDB(ctx, tAccounts.ID.EQ(jet.Uint64(claims.AccID)))
	if err != nil {
		if errors.Is(qrm.ErrNoRows, err) {
			return nil, ChangePasswordErr
		}

		return nil, err
	}

	// No password set
	if acc.Password == nil {
		return nil, NoAccountErr
	}

	// Password check logic
	if err := bcrypt.CompareHashAndPassword([]byte(*acc.Password), []byte(req.Current)); err != nil {
		return nil, ChangePasswordErr
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.New), 14)
	if err != nil {
		return nil, AccountCreateFailedErr
	}

	var char *users.User
	if claims.CharID > 0 {
		char, _, _, err = s.getCharacter(ctx, claims.CharID)
		if err != nil {
			return nil, ChangePasswordErr
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
		return nil, UpdateAccountErr
	}

	newToken, newClaims, err := s.createTokenFromAccountAndChar(acc, char)
	if err != nil {
		return nil, ChangePasswordErr
	}

	return &ChangePasswordResponse{
		Token:   newToken,
		Expires: timestamp.New(newClaims.ExpiresAt.Time),
	}, nil
}

func (s *Server) ForgotPassword(ctx context.Context, req *ForgotPasswordRequest) (*ForgotPasswordResponse, error) {
	acc, err := s.getAccountFromDB(ctx, jet.AND(
		tAccounts.RegToken.EQ(jet.String(req.RegToken)),
		tAccounts.Username.IS_NOT_NULL(),
		tAccounts.Password.IS_NULL(),
	))
	if err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, ForgotPasswordErr
		}
	}

	// We expect the account to not have a password for a "forgot password" request via token
	if acc.Password != nil {
		return nil, NoAccountErr
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.New), 14)
	if err != nil {
		return nil, ForgotPasswordErr
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
		return nil, ForgotPasswordErr
	}

	return &ForgotPasswordResponse{}, nil
}

func (s *Server) GetCharacters(ctx context.Context, req *GetCharactersRequest) (*GetCharactersResponse, error) {
	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, auth.InvalidTokenErr
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, GenericLoginErr
	}

	// Load account to make sure it (still) exists
	acc, err := s.getAccountFromDB(ctx, tAccounts.ID.EQ(jet.Uint64(claims.AccID)))
	if err != nil {
		return nil, GenericLoginErr
	}
	if acc.ID == 0 {
		return nil, GenericLoginErr
	}

	// Load chars from database
	stmt := tUsers.
		SELECT(
			tUsers.ID,
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
			tUsers.Visum,
			tUsers.Playtime,
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
			),
		).
		WHERE(
			tUsers.Identifier.LIKE(jet.String(buildCharSearchIdentifier(claims.Subject))),
		).
		ORDER_BY(tUsers.ID).
		LIMIT(10)

	resp := &GetCharactersResponse{}
	if err := stmt.QueryContext(ctx, s.db, &resp.Chars); err != nil {
		if errors.Is(qrm.ErrNoRows, err) {
			return nil, NoCharacterFoundErr
		}
		return nil, err
	}

	return resp, nil
}

func buildCharSearchIdentifier(license string) string {
	return "char%:" + license
}

func (s *Server) getCharacter(ctx context.Context, charId int32) (*users.User, *jobs.JobProps, string, error) {
	stmt := tUsers.
		SELECT(
			tUsers.ID,
			tUsers.Identifier,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Group.AS("group"),
			tJobs.Label.AS("user.job_label"),
			tJobGrades.Label.AS("user.job_grade_label"),
			tJobProps.Theme,
			tJobProps.QuickButtons,
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
				),
		).
		WHERE(
			tUsers.ID.EQ(jet.Int32(charId)),
		).
		LIMIT(1)

	var dest struct {
		users.User
		Group    string
		JobProps jobs.JobProps
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if errors.Is(qrm.ErrNoRows, err) {
			return nil, nil, "", NoCharacterFoundErr
		}
		return nil, nil, "", err
	}

	return &dest.User, &dest.JobProps, dest.Group, nil
}

func (s *Server) ChooseCharacter(ctx context.Context, req *ChooseCharacterRequest) (*ChooseCharacterResponse, error) {
	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, auth.InvalidTokenErr
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, err
	}

	char, jProps, userGroup, err := s.getCharacter(ctx, req.CharId)
	if err != nil {
		return nil, NoCharacterFoundErr
	}

	// Make sure the user isn't sending us a different char ID than their own
	if !strings.HasSuffix(char.Identifier, ":"+claims.Subject) {
		return nil, UnableToChooseCharErr
	}

	// Load account data for token creation
	account, err := s.getAccountFromDB(ctx, tAccounts.ID.EQ(jet.Uint64(claims.AccID)))
	if err != nil {
		return nil, NoCharacterFoundErr
	}

	// Reset override jobs when choosing a character
	if account.OverrideJob != nil {
		if err := s.ui.SetUserInfo(ctx, claims.AccID, "", 0); err != nil {
			return nil, GenericLoginErr
		}
	}

	newToken, newClaims, err := s.createTokenFromAccountAndChar(account, char)
	if err != nil {
		return nil, GenericLoginErr
	}

	// Load permissions of user
	userPs, err := s.p.GetPermissionsOfUser(&userinfo.UserInfo{
		UserId:   char.UserId,
		Job:      char.Job,
		JobGrade: char.JobGrade,
	})
	if err != nil {
		return nil, UnableToChooseCharErr
	}
	ps := userPs.GuardNames()

	if utils.InStringSlice(config.C.Game.SuperuserGroups, userGroup) {
		ps = append(ps, common.SuperuserPermission)
	}

	attrs, err := s.p.FlattenRoleAttributes(char.Job, char.JobGrade)
	if err != nil {
		return nil, GenericLoginErr
	}
	ps = append(ps, attrs...)

	if len(ps) == 0 {
		return nil, UnableToChooseCharErr
	} else if !utils.InStringSlice(ps, "authservice-choosecharacter") {
		return nil, UnableToChooseCharErr
	}

	s.a.AddEntryWithData(&model.FivenetAuditLog{
		Service: AuthService_ServiceDesc.ServiceName,
		Method:  "ChooseCharacter",
		UserID:  char.UserId,
		UserJob: char.Job,
		State:   int16(rector.EVENT_TYPE_VIEWED),
	}, char.UserShort())

	return &ChooseCharacterResponse{
		Token:       newToken,
		Expires:     timestamp.New(newClaims.ExpiresAt.Time),
		Permissions: ps,
		JobProps:    jProps,
		Char:        char,
	}, nil
}

func (s *Server) SetJob(ctx context.Context, req *SetJobRequest) (*SetJobResponse, error) {
	token, err := auth.GetTokenFromGRPCContext(ctx)
	if err != nil {
		return nil, auth.InvalidTokenErr
	}

	claims, err := s.tm.ParseWithClaims(token)
	if err != nil {
		return nil, err
	}

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	char, _, _, err := s.getCharacter(ctx, claims.CharID)
	if err != nil {
		return nil, err
	}

	job, jobGrade, jProps, err := s.getJobWithProps(ctx, req.Job)
	if err != nil {
		return nil, err
	}

	char.Job = job.Name
	char.JobGrade = jobGrade
	s.c.EnrichJobInfo(char)

	if err := s.ui.SetUserInfo(ctx, claims.AccID, char.Job, char.JobGrade); err != nil {
		return nil, err
	}

	userInfo.OrigJob = char.Job
	userInfo.OrigJobGrade = char.JobGrade
	userInfo.Job = job.Name
	userInfo.JobGrade = jobGrade

	// Load account data for token creation
	account, err := s.getAccountFromDB(ctx, tAccounts.Username.EQ(jet.String(claims.Username)))
	if err != nil {
		return nil, err
	}

	newToken, newClaims, err := s.createTokenFromAccountAndChar(account, char)
	if err != nil {
		return nil, err
	}

	return &SetJobResponse{
		Token:    newToken,
		Expires:  timestamp.New(newClaims.ExpiresAt.Time),
		JobProps: jProps,
		Char:     char,
	}, nil
}

func (s *Server) getJobWithProps(ctx context.Context, jobName string) (*jobs.Job, int32, *jobs.JobProps, error) {
	js := tJobs.AS("job")
	stmt := js.
		SELECT(
			js.Name,
			js.Label,
			tJobGrades.Grade.AS("job_grade"),
			tJobProps.AllColumns,
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
		Job      *jobs.Job
		JobGrade int32 `alias:"job_grade"`
		JobProps *jobs.JobProps
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, 0, nil, err
	}

	return dest.Job, dest.JobGrade, dest.JobProps, nil
}
