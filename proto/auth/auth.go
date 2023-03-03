package auth

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/galexrt/arpanet/model"
	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/helpers"
	"github.com/galexrt/arpanet/pkg/session"
	"github.com/galexrt/arpanet/proto/common"
	"github.com/galexrt/arpanet/query"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	AccountServiceServer
}

func NewServer() *Server {
	return &Server{}
}

// AuthFuncOverride is called instead of exampleAuthFunc
func (s *Server) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	// Skip authentication for the login endpoint
	if fullMethodName == "/gen.auth.AccountService/Login" {
		return ctx, nil
	}

	return auth.GRPCAuthFunc(ctx)
}

func (s *Server) createTokenFromAccountAndChar(account *model.Account, activeChar *common.Character) (string, error) {
	claims := &session.CitizenInfoClaims{
		AccountID: account.ID,
		Username:  account.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "arpanet",
			Subject:   account.License,
			ID:        strconv.FormatUint(uint64(account.ID), 10),
			Audience:  []string{"arpanet"},
		},
	}

	if activeChar != nil {
		claims.ActiveChar = activeChar.Identifier
		claims.ActiveCharID = uint(activeChar.Id)
	} else {
		claims.ActiveChar = "N/A"
		claims.ActiveCharID = 0
	}

	return session.Tokens.NewWithClaims(claims)
}

func (s *Server) getAccountFromDB(userID string) (*model.Account, error) {
	accounts := query.Account
	account, err := accounts.Where(accounts.Enabled.Is(true), accounts.Username.Eq(userID)).Limit(1).First()
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *Server) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	resp := &LoginResponse{}
	account, err := s.getAccountFromDB(req.Username)
	if err != nil {
		return nil, err
	}

	if !account.CheckPassword(req.Password) {
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
	resp := &GetCharactersResponse{}

	claims, err := session.Tokens.ParseWithClaims(auth.MustGetTokenFromGRPCContext(ctx))
	if err != nil {
		return resp, nil
	}

	// Load chars and add them to the response
	licenseSearch := helpers.BuildCharSearchIdentifier(claims.Subject)

	u := query.User
	users, err := u.Preload(u.UserLicenses.RelationField).
		Where(u.Identifier.Like(licenseSearch)).
		Limit(5).
		Find()
	if err != nil {
		return nil, nil
	}

	resp.Chars = helpers.ConvertModelUserListToCommonCharacterList(users)

	return resp, nil
}

func (s *Server) ChooseCharacter(ctx context.Context, req *ChooseCharacterRequest) (*ChooseCharacterResponse, error) {
	resp := &ChooseCharacterResponse{}

	claims, err := session.Tokens.ParseWithClaims(auth.MustGetTokenFromGRPCContext(ctx))
	if err != nil {
		return resp, nil
	}

	// Make sure the user isn't sending us a different char identifier than his own
	if !strings.Contains(req.Identifier, claims.Subject) {
		return nil, status.Error(codes.OutOfRange, "That's not your character!")
	}

	account, err := s.getAccountFromDB(claims.Username)
	if err != nil {
		return nil, err
	}

	u := query.User
	char, err := u.Where(u.Identifier.Eq(req.Identifier)).First()
	if err != nil {
		return nil, err
	}

	token, err := s.createTokenFromAccountAndChar(account, &common.Character{
		Id:         uint64(char.ID),
		Identifier: req.Identifier,
	})
	if err != nil {
		return nil, err
	}
	resp.Token = token

	// Load permissions of user
	perms, err := query.Perms.GetAllPermissionsOfUser(uint(char.ID))
	if err != nil {
		return nil, err
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
