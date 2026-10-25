package vehicles

import (
	"context"
	"database/sql"

	pbvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/vehicles"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/filestore"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/storage"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	citizenshydrator "github.com/fivenet-app/fivenet/v2026/stores/citizens/hydrator"
	vehiclesstore "github.com/fivenet-app/fivenet/v2026/stores/vehicles"
	"github.com/go-jet/jet/v2/mysql"
	"go.uber.org/fx"
	grpc "google.golang.org/grpc"
)

func lockVehicleFileLimitParent(ctx context.Context, tx *sql.Tx, plate string) error {
	tVehicle := table.FivenetVehicle
	var vehicle struct {
		Plate string `alias:"plate"`
	}
	return tVehicle.
		SELECT(tVehicle.Plate).
		FROM(tVehicle).
		WHERE(tVehicle.Plate.EQ(mysql.String(plate))).
		FOR(mysql.UPDATE()).
		QueryContext(ctx, tx, &vehicle)
}

func init() {
	filestore.RegisterFileReference(table.FivenetVehicleFiles, table.FivenetVehicleFiles.FileID)
}

type Server struct {
	pbvehicles.VehiclesServiceServer

	perms        perms.Permissions
	enricher     mstlystcdata.IUserAwareEnricher
	hydrator     citizenshydrator.IHydrator
	store        vehiclesstore.IStore
	appCfg       appconfig.IConfig
	db           *sql.DB
	filesHandler *filestore.Handler[string]
}

type Params struct {
	fx.In

	DB      *sql.DB
	Storage storage.IStorage

	Perms     perms.Permissions
	Enricher  mstlystcdata.IUserAwareEnricher
	Hydrator  citizenshydrator.IHydrator
	Store     vehiclesstore.IStore
	AppConfig appconfig.IConfig
}

func NewServer(p Params) *Server {
	vehicleFiles := table.FivenetVehicleFiles
	filesHandler := filestore.NewHandler(
		p.Storage,
		p.DB,
		vehicleFiles,
		vehicleFiles.Plate,
		vehicleFiles.FileID,
		5<<20,
		vehiclesstore.MaxVehicleFiles,
		func(plate string) mysql.BoolExpression { return vehicleFiles.Plate.EQ(mysql.String(plate)) },
		vehicleFileInserter,
		false,
	).WithUploadFilter(filestore.NewImageUploadFilter()).
		WithFileLimitLock(lockVehicleFileLimitParent)
	return &Server{
		db:           p.DB,
		perms:        p.Perms,
		enricher:     p.Enricher,
		hydrator:     p.Hydrator,
		store:        p.Store,
		appCfg:       p.AppConfig,
		filesHandler: filesHandler,
	}
}

func vehicleFileInserter(
	ctx context.Context,
	tx *sql.Tx,
	_ mysql.Table,
	_ mysql.Column,
	_ mysql.ColumnInteger,
	plate string,
	_ mysql.BoolExpression,
	fileID int64,
) error {
	files := table.FivenetVehicleFiles
	_, err := files.INSERT(files.Plate, files.FileID, files.Kind).
		VALUES(plate, fileID, "attachment").
		ExecContext(ctx, tx)
	return err
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbvehicles.RegisterVehiclesServiceServer(srv, s)
}
