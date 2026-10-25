package citizens

import (
	context "context"
	"database/sql"
	"errors"
	"math"
	"slices"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	users "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	pbcitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens"
	permscitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/filestore"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/storage"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	citizensstore "github.com/fivenet-app/fivenet/v2026/stores/citizens"
	citizenshydrator "github.com/fivenet-app/fivenet/v2026/stores/citizens/hydrator"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
	grpc "google.golang.org/grpc"
)

const citizenIDLogFieldKey = "fivenet.citizens.user_id"

func lockCitizenFileLimitParent(ctx context.Context, tx *sql.Tx, userID int32) error {
	tUser := table.FivenetUser
	var user struct {
		ID int32 `alias:"id"`
	}
	return tUser.
		SELECT(tUser.ID).
		FROM(tUser).
		WHERE(tUser.ID.EQ(mysql.Int32(userID))).
		FOR(mysql.UPDATE()).
		QueryContext(ctx, tx, &user)
}

func replaceUserMugshot(
	ctx context.Context,
	tx *sql.Tx,
	join mysql.Table,
	parentCol mysql.Column,
	fileCol mysql.ColumnInteger,
	userID int32,
	parentIDBoolExp mysql.BoolExpression,
	fileID int64,
) (int64, error) {
	tUserFiles := table.FivenetUserFiles
	var current struct {
		FileID int64 `alias:"file_id"`
	}
	err := tUserFiles.
		SELECT(tUserFiles.FileID.AS("file_id")).
		WHERE(mysql.AND(
			tUserFiles.UserID.EQ(mysql.Int32(userID)),
			tUserFiles.Kind.EQ(mysql.String("mugshot")),
		)).
		FOR(mysql.UPDATE()).
		QueryContext(ctx, tx, &current)
	if err != nil && !errors.Is(err, qrm.ErrNoRows) {
		return 0, err
	}

	if errors.Is(err, qrm.ErrNoRows) {
		_, err = tUserFiles.
			INSERT(tUserFiles.UserID, tUserFiles.FileID, tUserFiles.Kind).
			VALUES(userID, fileID, "mugshot").
			ExecContext(ctx, tx)
		return 0, err
	}

	_, err = tUserFiles.
		UPDATE(tUserFiles.FileID, tUserFiles.CreatedAt).
		SET(fileID, mysql.CURRENT_TIMESTAMP()).
		WHERE(mysql.AND(parentIDBoolExp, tUserFiles.Kind.EQ(mysql.String("mugshot")))).
		LIMIT(1).
		ExecContext(ctx, tx)
	return current.FileID, err
}

func findUserMugshot(
	ctx context.Context,
	tx *sql.Tx,
	join mysql.Table,
	parentCol mysql.Column,
	fileCol mysql.ColumnInteger,
	userID int32,
) (int64, error) {
	tUserFiles := table.FivenetUserFiles
	var current struct {
		FileID int64 `alias:"file_id"`
	}
	err := tUserFiles.
		SELECT(tUserFiles.FileID.AS("file_id")).
		WHERE(mysql.AND(
			tUserFiles.UserID.EQ(mysql.Int32(userID)),
			tUserFiles.Kind.EQ(mysql.String("mugshot")),
		)).
		FOR(mysql.UPDATE()).
		QueryContext(ctx, tx, &current)
	if errors.Is(err, qrm.ErrNoRows) {
		return 0, nil
	}
	return current.FileID, err
}

func init() {
	filestore.RegisterFileReference(table.FivenetUser, table.FivenetUser.AvatarFileID)
}

type Server struct {
	pbcitizens.CitizensServiceServer
	pbcitizens.LabelsServiceServer

	db       *sql.DB
	perms    perms.Permissions
	enricher mstlystcdata.IUserAwareEnricher
	st       storage.IStorage
	appCfg   appconfig.IConfig
	cfg      *config.Config
	customDB *config.CustomDB
	notifi   notifi.INotifi
	store    citizensstore.IStore
	hydrator citizenshydrator.IHydrator

	profilePictureHandler *filestore.Handler[int32]
	mugshotHandler        *filestore.Handler[int32]
	citizenFilesHandler   *filestore.Handler[int32]

	labelsAccess         *access.CitizenLabelsObjectAccess
	labelsAccessResolver *access.SubjectResolver
}

type Params struct {
	fx.In

	DB           *sql.DB
	P            perms.Permissions
	Enricher     mstlystcdata.IUserAwareEnricher
	Config       *config.Config
	Storage      storage.IStorage
	AppConfig    appconfig.IConfig
	Notifi       notifi.INotifi
	Store        citizensstore.IStore
	Hydrator     citizenshydrator.IHydrator
	LabelsAccess *access.CitizenLabelsObjectAccess
}

func NewServer(p Params) *Server {
	tUserFiles := table.FivenetUserFiles

	profilePictureHandler := filestore.NewHandler(
		p.Storage,
		p.DB,
		table.FivenetUser,
		table.FivenetUser.ID,
		table.FivenetUser.AvatarFileID,
		3<<20,
		0,
		func(parentId int32) mysql.BoolExpression {
			return table.FivenetUser.ID.EQ(mysql.Int32(parentId))
		},
		userAvatarInserter,
		true,
	).WithUploadFilter(filestore.NewImageUploadFilter())
	mugshotHandler := filestore.NewHandler(
		p.Storage,
		p.DB,
		tUserFiles,
		tUserFiles.UserID,
		tUserFiles.FileID,
		3<<20,
		citizensstore.MaxUserFiles,
		func(parentId int32) mysql.BoolExpression {
			return tUserFiles.UserID.EQ(mysql.Int32(parentId))
		},
		userFileJoinInserter("mugshot"),
		false,
	).WithUploadFilter(filestore.NewImageUploadFilter()).
		WithFileLimitLock(lockCitizenFileLimitParent)
	citizenFilesHandler := filestore.NewHandler(
		p.Storage,
		p.DB,
		tUserFiles,
		tUserFiles.UserID,
		tUserFiles.FileID,
		5<<20,
		citizensstore.MaxUserFiles,
		func(parentID int32) mysql.BoolExpression { return tUserFiles.UserID.EQ(mysql.Int32(parentID)) },
		userFileJoinInserter("attachment"),
		false,
	).WithUploadFilter(filestore.NewImageUploadFilter()).
		WithFileLimitLock(lockCitizenFileLimitParent)

	s := &Server{
		db:       p.DB,
		perms:    p.P,
		enricher: p.Enricher,
		st:       p.Storage,
		appCfg:   p.AppConfig,
		cfg:      p.Config,
		customDB: &p.Config.Database.Custom,
		notifi:   p.Notifi,
		store:    p.Store,
		hydrator: p.Hydrator,

		profilePictureHandler: profilePictureHandler,
		mugshotHandler:        mugshotHandler,
		citizenFilesHandler:   citizenFilesHandler,

		labelsAccess:         p.LabelsAccess,
		labelsAccessResolver: access.NewSubjectResolver(p.DB),
	}

	access.RegisterAccess("citizen", &access.GroupedAccessAdapter{
		CanUserAccessTargetFn: func(ctx context.Context, targetId int64, userInfo *userinfo.UserInfo, access int32) (bool, error) {
			if !s.perms.Can(
				userInfo,
				permscitizens.CitizensService.GetUser.Perm,
			) {
				return false, nil
			}

			if targetId < int64(math.MinInt32) || targetId > int64(math.MaxInt32) {
				return false, nil // targetId is too large to fit in int32
			}
			userId := int32(targetId)
			tUser := table.FivenetUser.AS("user")

			// Retrieve user job from database
			stmt := tUser.
				SELECT(
					tUser.ID,
					tUser.Job,
				).
				FROM(tUser).
				WHERE(tUser.ID.EQ(mysql.Int32(userId))).
				LIMIT(1)

			user := &users.User{}
			if err := stmt.QueryContext(ctx, s.db, user); err != nil {
				return false, err
			}

			if slices.Contains(s.appCfg.Get().JobInfo.GetPublicJobs(), user.GetJob()) ||
				slices.Contains(s.appCfg.Get().JobInfo.GetHiddenJobs(), user.GetJob()) {
				// Make sure user has permission to see that grade
				check, err := s.checkIfUserCanAccess(userInfo, user.GetJob(), user.GetJobGrade())
				if err != nil {
					return false, err
				}
				if !check {
					return false, nil
				}
			}

			return true, nil
		},
	})

	return s
}

func userAvatarInserter(
	ctx context.Context,
	tx *sql.Tx,
	_ mysql.Table,
	parentCol mysql.Column,
	fileCol mysql.ColumnInteger,
	userID int32,
	_ mysql.BoolExpression,
	fileID int64,
) error {
	user := table.FivenetUser
	_, err := user.UPDATE(fileCol).SET(mysql.IntExp(mysql.Int64(fileID))).WHERE(
		user.ID.EQ(mysql.Int32(userID)),
	).ExecContext(ctx, tx)
	return err
}

func userFileJoinInserter(kind string) filestore.JoinRowInserterFn[int32] {
	return func(
		ctx context.Context,
		tx *sql.Tx,
		_ mysql.Table,
		_ mysql.Column,
		_ mysql.ColumnInteger,
		userID int32,
		_ mysql.BoolExpression,
		fileID int64,
	) error {
		tUserFiles := table.FivenetUserFiles
		_, err := tUserFiles.
			INSERT(tUserFiles.UserID, tUserFiles.FileID, tUserFiles.Kind).
			VALUES(userID, fileID, kind).
			ExecContext(ctx, tx)
		return err
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbcitizens.RegisterCitizensServiceServer(srv, s)
	pbcitizens.RegisterLabelsServiceServer(srv, s)
}
