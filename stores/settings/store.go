package settingsstore

import (
	"context"
	"database/sql"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	jobsprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/props"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/laws"
	resourcesettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/settings"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/settings"
)

type IStore interface {
	ViewAuditLog(
		ctx context.Context,
		req *pbsettings.ViewAuditLogRequest,
	) (*pbsettings.ViewAuditLogResponse, error)
	ListAccounts(
		ctx context.Context,
		req *pbsettings.ListAccountsRequest,
	) (*pbsettings.ListAccountsResponse, error)
	UpdateAccount(
		ctx context.Context,
		req *pbsettings.UpdateAccountRequest,
	) (*pbsettings.UpdateAccountResponse, error)
	DisconnectSocialLogin(ctx context.Context, accountID int64, provider string) error
	GetAccountByID(ctx context.Context, accountID int64) (*accounts.Account, error)
	DeleteAccount(
		ctx context.Context,
		accountID int64,
		deletedAtTime *timestamp.Timestamp,
	) (*pbsettings.DeleteAccountResponse, error)
	UpdateAppConfig(ctx context.Context, cfg *resourcesettings.AppConfig) error
	ListLawBooks(ctx context.Context, superuser bool) (*pbsettings.ListLawBooksResponse, error)
	CreateOrUpdateLawBook(
		ctx context.Context,
		req *pbsettings.CreateOrUpdateLawBookRequest,
		superuser bool,
	) (*laws.LawBook, error)
	DeleteLawBook(ctx context.Context, lawbookID int64, deletedAtTime *timestamp.Timestamp) error
	ReorderLawBooks(ctx context.Context, req *pbsettings.ReorderLawBooksRequest) error
	GetLawBook(ctx context.Context, lawbookID int64) (*laws.LawBook, error)
	CreateOrUpdateLaw(
		ctx context.Context,
		req *pbsettings.CreateOrUpdateLawRequest,
		superuser bool,
	) (*laws.Law, []int64, error)
	DeleteLaw(ctx context.Context, lawID int64, deletedAtTime *timestamp.Timestamp) error
	ReorderLaws(ctx context.Context, req *pbsettings.ReorderLawsRequest) error
	GetLaw(ctx context.Context, lawId int64) (*laws.Law, error)
	SetJobProps(ctx context.Context, props *jobsprops.JobProps) error
	DeleteJobProps(ctx context.Context, job string) error
}

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) IStore {
	return &Store{db: db}
}
