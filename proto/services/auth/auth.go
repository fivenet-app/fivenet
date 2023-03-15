package auth

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/perms"
	users "github.com/galexrt/arpanet/proto/resources/users"
	"github.com/galexrt/arpanet/query/arpanet/model"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	a   = table.ArpanetAccounts
	u   = table.Users.AS("user")
	aup = table.ArpanetUserProps
)

var (
	InvalidLoginErr       = status.Error(codes.NotFound, "Wrong username or password")
	NoCharacterFoundErr   = status.Error(codes.NotFound, "No character found for your account")
	GenericLoginErr       = status.Error(codes.Internal, "Failed to login you in")
	UnableToChooseCharErr = status.Error(codes.PermissionDenied, "You don't have permission to select this character!")
)

type Server struct {
	AuthServiceServer

	db   *sql.DB
	auth *auth.GRPCAuth
	tm   *auth.TokenManager
	p    perms.Permissions
}

func NewServer(db *sql.DB, auth *auth.GRPCAuth, tm *auth.TokenManager, p perms.Permissions) *Server {
	return &Server{
		db:   db,
		auth: auth,
		tm:   tm,
		p:    p,
	}
}

// AuthFuncOverride is called instead of exampleAuthFunc
func (s *Server) AuthFuncOverride(ctx context.Context, fullMethod string) (context.Context, error) {
	// Skip authentication for the login endpoint
	if fullMethod == "/services.auth.AuthService/Login" {
		return ctx, nil
	}

	return s.auth.GRPCAuthFunc(ctx, fullMethod)
}

func (s *Server) PermissionUnaryFuncOverride(ctx context.Context, info *grpc.UnaryServerInfo) (context.Context, error) {
	// Skip permission check for the auth services
	return ctx, nil
}

func (s *Server) createTokenFromAccountAndChar(account *model.ArpanetAccounts, activeChar *users.User) (string, error) {
	claims := &auth.CitizenInfoClaims{
		AccountID: account.ID,
		Username:  account.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "arpanet",
			Subject:   account.License,
			ID:        strconv.FormatUint(uint64(account.ID), 10),
			Audience:  []string{"arpanet"},
		},
	}

	if activeChar != nil {
		claims.ActiveCharID = activeChar.UserId
		claims.ActiveCharJob = activeChar.Job
		claims.ActiveCharJobGrade = activeChar.JobGrade
	} else {
		claims.ActiveCharID = 0
		claims.ActiveCharJob = ""
		claims.ActiveCharJobGrade = 0
	}

	return s.tm.NewWithClaims(claims)
}

func (s *Server) getAccountFromDB(ctx context.Context, username string) (*model.ArpanetAccounts, error) {
	stmt := a.SELECT(
		a.AllColumns,
	).
		FROM(a).
		WHERE(
			a.Enabled.IS_TRUE().
				AND(a.Username.EQ(jet.String(username))),
		).LIMIT(1)

	var account model.ArpanetAccounts
	if err := stmt.QueryContext(ctx, s.db, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *Server) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	account, err := s.getAccountFromDB(ctx, req.Username)
	if err != nil {
		if errors.Is(qrm.ErrNoRows, err) {
			return nil, InvalidLoginErr
		}

		return nil, err
	}

	// Password check logic
	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(req.Password)); err != nil {
		return nil, InvalidLoginErr
	}

	token, err := s.createTokenFromAccountAndChar(account, nil)
	if err != nil {
		return nil, InvalidLoginErr
	}

	return &LoginResponse{
		Token: token,
	}, nil
}

func (s *Server) GetCharacters(ctx context.Context, req *GetCharactersRequest) (*GetCharactersResponse, error) {
	claims, err := s.tm.ParseWithClaims(auth.MustGetTokenFromGRPCContext(ctx))
	if err != nil {
		return nil, err
	}

	// Load chars from database
	stmt := u.SELECT(
		u.AllColumns,
	).
		FROM(u.LEFT_JOIN(aup, aup.UserID.EQ(u.ID))).
		WHERE(u.Identifier.LIKE(jet.String(buildCharSearchIdentifier(claims.Subject)))).
		ORDER_BY(u.ID).
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

func (s *Server) ChooseCharacter(ctx context.Context, req *ChooseCharacterRequest) (*ChooseCharacterResponse, error) {
	claims, err := s.tm.ParseWithClaims(auth.MustGetTokenFromGRPCContext(ctx))
	if err != nil {
		return nil, err
	}

	stmt := u.SELECT(
		u.ID,
		u.Identifier,
		u.Job,
		u.JobGrade,
	).
		FROM(u).
		WHERE(u.ID.EQ(jet.Int32(req.CharId))).
		LIMIT(1)

	var char users.User
	if err := stmt.QueryContext(ctx, s.db, &char); err != nil {
		if errors.Is(qrm.ErrNoRows, err) {
			return nil, NoCharacterFoundErr
		}
		return nil, err
	}

	// Make sure the user isn't sending us a different char ID than their own
	if !strings.Contains(char.Identifier, claims.Subject) {
		return nil, UnableToChooseCharErr
	}

	// Load account data for token creation
	account, err := s.getAccountFromDB(ctx, claims.Username)
	if err != nil {
		return nil, err
	}

	token, err := s.createTokenFromAccountAndChar(account, &char)
	if err != nil {
		return nil, err
	}

	// Load permissions of user
	perms, err := s.p.GetAllPermissionsOfUser(char.UserId)
	if err != nil {
		return nil, err
	}

	if len(perms) == 0 {
		return nil, UnableToChooseCharErr
	}

	return &ChooseCharacterResponse{
		Token:       token,
		Permissions: perms.GuardNames(),
	}, nil
}

func (s *Server) Logout(ctx context.Context, req *LogoutRequest) (*LogoutResponse, error) {
	// TODO till we have a JWT token manager "blocking" users when they logout, nothing todo here
	return &LogoutResponse{
		Success: true,
	}, nil
}
