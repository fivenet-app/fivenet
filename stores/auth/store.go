package authstore

import (
	"context"
	"database/sql"

	accounts "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	accountsoauth2 "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts/oauth2"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	jobsprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/props"
	users "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/model"
)

type IStore interface {
	GetAccountByID(
		ctx context.Context,
		accountID int64,
		withPassword bool,
	) (*model.FivenetAccounts, error)
	GetAccountByUsername(
		ctx context.Context,
		username string,
		withPassword bool,
	) (*model.FivenetAccounts, error)
	GetLoginAccountByUsername(ctx context.Context, username string) (*model.FivenetAccounts, error)
	GetAccountByIDAndUsername(
		ctx context.Context,
		accountID int64,
		username string,
		withPassword bool,
	) (*model.FivenetAccounts, error)
	GetAccountByRegToken(
		ctx context.Context,
		regToken string,
		withPassword bool,
	) (*model.FivenetAccounts, error)
	GetPasswordResetAccountByRegToken(
		ctx context.Context,
		regToken string,
	) (*model.FivenetAccounts, error)
	ActivateAccount(
		ctx context.Context,
		accountID int64,
		regToken, username, hashedPassword string,
	) error
	UpdatePassword(ctx context.Context, accountID int64, hashedPassword string) error
	UpdateUsername(ctx context.Context, accountID int64, username string) error
	ForgotPassword(ctx context.Context, accountID int64, hashedPassword string) error
	ListCharacters(
		ctx context.Context,
		accountID int64,
		license string,
	) ([]*accounts.Character, error)
	GetCharacter(ctx context.Context, charID int32) (*users.User, *jobsprops.JobProps, error)
	GetJobWithProps(
		ctx context.Context,
		jobName string,
	) (*jobs.Job, int32, *jobsprops.JobProps, error)
	ListOAuth2Connections(
		ctx context.Context,
		accountID int64,
	) ([]*accountsoauth2.OAuth2Account, error)
	DeleteSocialLogin(ctx context.Context, accountID int64, provider string) error
}

type Store struct {
	db       *sql.DB
	customDB *config.CustomDB
}

func New(db *sql.DB, customDB *config.CustomDB) IStore {
	return &Store{
		db:       db,
		customDB: customDB,
	}
}
