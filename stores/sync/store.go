package syncstore

import (
	"context"
	"database/sql"

	syncdata "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/data"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/dispatches"
	citizensstore "github.com/fivenet-app/fivenet/v2026/stores/citizens"
	jobsstore "github.com/fivenet-app/fivenet/v2026/stores/jobs"
	"go.uber.org/zap"
)

type IStore interface {
	AddActivity(ctx context.Context, req *pbsync.AddActivityRequest) (*pbsync.AddActivityResponse, error)
	AddUserUpdate(ctx context.Context, req *pbsync.AddUserUpdateRequest) (*pbsync.AddActivityResponse, error)
	AddUserActivity(ctx context.Context, req *pbsync.AddUserActivityRequest) (*pbsync.AddActivityResponse, error)
	AddUserProps(ctx context.Context, req *pbsync.AddUserPropsRequest) (*pbsync.AddActivityResponse, error)
	AddColleagueActivity(ctx context.Context, req *pbsync.AddColleagueActivityRequest) (*pbsync.AddActivityResponse, error)
	AddColleagueProps(ctx context.Context, req *pbsync.AddColleaguePropsRequest) (*pbsync.AddActivityResponse, error)
	AddJobTimeclock(ctx context.Context, req *pbsync.AddJobTimeclockRequest) (*pbsync.AddActivityResponse, error)
	AddDispatch(ctx context.Context, req *pbsync.AddDispatchRequest) (*pbsync.AddActivityResponse, error)

	SendUsers(ctx context.Context, data []*syncdata.DataUser) (int64, error)
	DeleteUsers(ctx context.Context, userIDs []int32) (*pbsync.DeleteDataResponse, error)
	SendVehicles(ctx context.Context, req *pbsync.SendVehiclesRequest) (*pbsync.SendDataResponse, error)
	DeleteVehicles(ctx context.Context, plates []string) (*pbsync.DeleteDataResponse, error)
	SendJobs(ctx context.Context, req *pbsync.SendJobsRequest) (*pbsync.SendDataResponse, error)
	SendLicenses(ctx context.Context, req *pbsync.SendLicensesRequest) (*pbsync.SendDataResponse, error)
	SendAccounts(ctx context.Context, req *pbsync.SendAccountsRequest) (*pbsync.SendDataResponse, error)
	SendUserLocations(ctx context.Context, req *pbsync.SendUserLocationsRequest) (*pbsync.SendDataResponse, error)
	SetLastCharID(ctx context.Context, req *pbsync.SetLastCharIDRequest) (*pbsync.SendDataResponse, error)
	SendData(ctx context.Context, req *pbsync.SendDataRequest) (*pbsync.SendDataResponse, error)
	DeleteData(ctx context.Context, req *pbsync.DeleteDataRequest) (*pbsync.DeleteDataResponse, error)

	RegisterAccount(ctx context.Context, req *pbsync.RegisterAccountRequest) (*pbsync.RegisterAccountResponse, error)
	TransferAccount(ctx context.Context, req *pbsync.TransferAccountRequest) (*pbsync.TransferAccountResponse, error)
	AddAccountUpdate(ctx context.Context, req *pbsync.AddAccountUpdateRequest) (*pbsync.AddActivityResponse, error)
	AddUserOAuth2Conn(ctx context.Context, req *pbsync.AddUserOAuth2ConnRequest) (*pbsync.AddActivityResponse, error)

	CountJobs(ctx context.Context) (int64, error)
	CountAccounts(ctx context.Context) (int64, error)
	CountUsers(ctx context.Context) (int64, error)
	CountVehicles(ctx context.Context) (int64, error)
	CountLicenses(ctx context.Context) (int64, error)
}

type Store struct {
	db *sql.DB

	logger *zap.Logger
	cfg    *config.Config

	dispatches    *dispatches.DispatchDB
	citizensStore citizensstore.IStore
	jobsStore     jobsstore.IStore
}

func New(
	db *sql.DB,
	logger *zap.Logger,
	cfg *config.Config,
	dispatches *dispatches.DispatchDB,
	citizensStore citizensstore.IStore,
	jobsStore jobsstore.IStore,
) IStore {
	return &Store{
		db:            db,
		logger:        logger,
		cfg:           cfg,
		dispatches:    dispatches,
		citizensStore: citizensStore,
		jobsStore:     jobsStore,
	}
}
