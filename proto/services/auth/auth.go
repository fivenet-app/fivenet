package auth

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/perms"
	"github.com/galexrt/arpanet/pkg/session"
	users "github.com/galexrt/arpanet/proto/resources/users"
	"github.com/galexrt/arpanet/query"
	"github.com/galexrt/arpanet/query/arpanet/model"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	a   = table.ArpanetAccounts
	u   = table.Users.AS("user")
	aup = table.ArpanetUserProps
)

type Server struct {
	AuthServiceServer
}

func NewServer() *Server {
	return &Server{}
}

// AuthFuncOverride is called instead of exampleAuthFunc
func (s *Server) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	// Skip authentication for the login endpoint
	if fullMethodName == "/services.auth.AuthService/Login" {
		return ctx, nil
	}

	return auth.GRPCAuthFunc(ctx)
}

func (s *Server) createTokenFromAccountAndChar(account *model.ArpanetAccounts, activeChar *users.User) (string, error) {
	claims := &session.CitizenInfoClaims{
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
		claims.ActiveCharID = activeChar.UserID
		claims.ActiveCharJob = activeChar.Job
		claims.ActiveCharJobGrade = activeChar.JobGrade
	} else {
		claims.ActiveCharID = 0
		claims.ActiveCharJob = ""
		claims.ActiveCharJobGrade = 0
	}

	return session.Tokens.NewWithClaims(claims)
}

func (s *Server) getAccountFromDB(ctx context.Context, username string) (*model.ArpanetAccounts, error) {
	var account model.ArpanetAccounts
	stmt := a.SELECT(
		a.AllColumns,
	).
		FROM(a).
		WHERE(
			a.Enabled.IS_TRUE().
				AND(a.Username.EQ(jet.String(username))),
		).LIMIT(1)
	if err := stmt.QueryContext(ctx, query.DB, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *Server) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	resp := &LoginResponse{}
	account, err := s.getAccountFromDB(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	// Password check logic
	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(req.Password)); err != nil {
		return &LoginResponse{}, errors.New("wrong username or password")
	}

	token, err := s.createTokenFromAccountAndChar(account, nil)
	if err != nil {
		return nil, err
	}
	resp.Token = token

	return resp, nil
}

func (s *Server) GetCharacters(ctx context.Context, req *GetCharactersRequest) (*GetCharactersResponse, error) {
	claims, err := session.Tokens.ParseWithClaims(auth.MustGetTokenFromGRPCContext(ctx))
	if err != nil {
		return nil, err
	}

	resp := &GetCharactersResponse{}
	// Load chars from database
	stmt := u.SELECT(
		u.AllColumns,
	).
		FROM(u.LEFT_JOIN(aup, aup.UserID.EQ(u.ID))).
		WHERE(u.Identifier.LIKE(jet.String(buildCharSearchIdentifier(claims.Subject)))).
		ORDER_BY(u.ID).
		LIMIT(10)

	if err := stmt.QueryContext(ctx, query.DB, &resp.Chars); err != nil {
		return nil, err
	}

	return resp, nil
}

func buildCharSearchIdentifier(license string) string {
	return "char%:" + license
}

func (s *Server) ChooseCharacter(ctx context.Context, req *ChooseCharacterRequest) (*ChooseCharacterResponse, error) {
	resp := &ChooseCharacterResponse{}

	claims, err := session.Tokens.ParseWithClaims(auth.MustGetTokenFromGRPCContext(ctx))
	if err != nil {
		return nil, err
	}

	var char users.User
	stmt := u.SELECT(
		u.ID,
		u.Identifier,
		u.Job,
		u.JobGrade,
	).
		FROM(u).
		WHERE(u.ID.EQ(jet.Int32(req.UserID))).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, query.DB, &char); err != nil {
		return nil, err
	}

	// Make sure the user isn't sending us a different char ID than their own
	if !strings.Contains(char.Identifier, claims.Subject) {
		return nil, status.Error(codes.OutOfRange, "That's not your character!")
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
	resp.Token = token

	// Load permissions of user
	perms, err := perms.P.GetAllPermissionsOfUser(char.UserID)
	if err != nil {
		return nil, err
	}

	if len(perms) == 0 {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to select this character!")
	}

	resp.Permissions = perms.GuardNames()

	return resp, nil
}

func (s *Server) Logout(ctx context.Context, req *LogoutRequest) (*LogoutResponse, error) {
	// TODO till we have a JWT token manager "blocking" users when they logout, nothing todo here
	resp := &LogoutResponse{
		Success: true,
	}

	return resp, nil
}
