package auth

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/galexrt/arpanet/model"
	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/session"
	"github.com/galexrt/arpanet/query"
	"github.com/golang-jwt/jwt/v5"
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

func (s *Server) createTokenForAccount(account *model.Account, charIndex int) (string, error) {
	return session.Tokens.NewWithClaims(&session.UserInfoClaims{
		AccountID: account.ID,
		Username:  account.Username,
		CharIndex: charIndex,
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
	})
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

	token, err := s.createTokenForAccount(account, 0)
	if err != nil {
		return nil, err
	}
	resp.Token = token

	// Load chars and add them to the response
	chars, err := auth.GetCharsByLicense(auth.BuildCharSearchIdentifier(account.License))
	if err != nil {
		return resp, err
	}

	for _, char := range chars {
		resp.Chars = append(resp.Chars, &Character{
			Identifier:  char.Identifier,
			Job:         *char.Job,
			JobGrade:    int32(char.JobGrade),
			Firstname:   *char.Firstname,
			Lastname:    *char.Lastname,
			Dateofbirth: *char.Dateofbirth,
			Sex:         string(char.Sex),
			Height:      *char.Height,
			Visum:       int64(*char.Visum),
			Playtime:    int64(*char.Playtime),
		})
	}

	return resp, nil
}

func (s *Server) ChooseCharacter(ctx context.Context, req *ChooseCharacterRequest) (*ChooseCharacterResponse, error) {
	resp := &ChooseCharacterResponse{}

	claims, err := session.Tokens.ParseWithClaims(req.Token)
	if err != nil {
		return resp, nil
	}

	account, err := s.getAccountFromDB(claims.Username)
	if err != nil {
		return nil, err
	}

	token, err := s.createTokenForAccount(account, int(req.Index))
	if err != nil {
		return nil, err
	}
	resp.Token = token
	return resp, nil
}

func (s *Server) Logout(ctx context.Context, req *LogoutRequest) (*LogoutResponse, error) {
	resp := &LogoutResponse{}
	// TODO till we have a JWT token manager "blocking" users when they logout, nothing todo here
	return resp, nil
}
