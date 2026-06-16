package settingsstore

import (
	"context"
	"database/sql"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	audit "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	jobsprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/props"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/laws"
	resourcesettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/settings"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/settings"
	"github.com/go-jet/jet/v2/mysql"
)

type ViewAuditLogOptions struct {
	Pagination *database.PaginationRequest
	Sort       *database.Sort
	UserIDs    []int32
	From       *timestamp.Timestamp
	To         *timestamp.Timestamp
	Services   []string
	Methods    []string
	Actions    []audit.EventAction
	Results    []audit.EventResult
	Search     string
}

type ListAccountsOptions struct {
	Pagination   *database.PaginationRequest
	Sort         *database.Sort
	License      string
	OnlyDisabled bool
	Username     string
	ExternalID   string
	Group        string
}

type IStore interface {
	ViewAuditLog(
		ctx context.Context,
		opts ViewAuditLogOptions,
	) (*pbsettings.ViewAuditLogResponse, error)
	ListAccounts(
		ctx context.Context,
		opts ListAccountsOptions,
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
	db             *sql.DB
	accountSorter  *database.SorterBuilder
	auditLogSorter *database.SorterBuilder
}

func New(db *sql.DB) IStore {
	return &Store{
		db: db,
		accountSorter: database.New(
			database.SpecMap{
				"license":  database.Column{Col: tAccounts.License},
				"username": database.Column{Col: tAccounts.Username},
				"id":       database.Column{Col: tAccounts.ID},
			},
			[]mysql.OrderByClause{tAccounts.CreatedAt.DESC()},
			nil,
			"id",
			3,
		),
		auditLogSorter: database.New(
			database.SpecMap{
				"service":   database.Column{Col: tAuditLog.Service},
				"action":    database.Column{Col: tAuditLog.Action},
				"createdAt": database.Column{Col: tAuditLog.CreatedAt},
			},
			[]mysql.OrderByClause{tAuditLog.CreatedAt.DESC()},
			nil,
			"createdAt",
			3,
		),
	}
}
