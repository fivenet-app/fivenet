package citizens

import (
	context "context"
	"database/sql"
	"math"
	"slices"

	citizenslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/citizens/labels"
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
	"github.com/go-jet/jet/v2/mysql"
	"go.uber.org/fx"
	grpc "google.golang.org/grpc"
)

type Server struct {
	pbcitizens.CitizensServiceServer
	pbcitizens.LabelsServiceServer

	db       *sql.DB
	ps       perms.Permissions
	enricher *mstlystcdata.UserAwareEnricher
	st       storage.IStorage
	appCfg   appconfig.IConfig
	cfg      *config.Config
	customDB config.CustomDB
	notifi   notifi.INotifi

	profilePictureHandler *filestore.Handler[int32]
	mugshotHandler        *filestore.Handler[int32]

	labelsAccess *access.Grouped[citizenslabels.JobAccess, *citizenslabels.JobAccess, access.DummyUserAccess[citizenslabels.AccessLevel], *access.DummyUserAccess[citizenslabels.AccessLevel], access.DummyQualificationAccess[citizenslabels.AccessLevel], *access.DummyQualificationAccess[citizenslabels.AccessLevel], citizenslabels.AccessLevel]
}

type Params struct {
	fx.In

	DB        *sql.DB
	P         perms.Permissions
	Enricher  *mstlystcdata.UserAwareEnricher
	Config    *config.Config
	Storage   storage.IStorage
	AppConfig appconfig.IConfig
	Notifi    notifi.INotifi
}

func NewServer(p Params) *Server {
	tUserProps := table.FivenetUserProps

	profilePictureHandler := filestore.NewHandler(
		p.Storage,
		p.DB,
		tUserProps,
		tUserProps.UserID,
		tUserProps.AvatarFileID,
		3<<20,
		1,
		func(parentId int32) mysql.BoolExpression {
			return tUserProps.UserID.EQ(mysql.Int32(parentId))
		},
		filestore.UpdateJoinRow,
		true,
	).WithUploadFilter(filestore.NewImageUploadFilter())
	mugshotHandler := filestore.NewHandler(
		p.Storage,
		p.DB,
		tUserProps,
		tUserProps.UserID,
		tUserProps.MugshotFileID,
		3<<20,
		1,
		func(parentId int32) mysql.BoolExpression {
			return tUserProps.UserID.EQ(mysql.Int32(parentId))
		},
		filestore.UpdateJoinRow,
		true,
	).WithUploadFilter(filestore.NewImageUploadFilter())

	s := &Server{
		db:       p.DB,
		ps:       p.P,
		enricher: p.Enricher,
		st:       p.Storage,
		appCfg:   p.AppConfig,
		cfg:      p.Config,
		customDB: p.Config.Database.Custom,
		notifi:   p.Notifi,

		profilePictureHandler: profilePictureHandler,
		mugshotHandler:        mugshotHandler,

		labelsAccess: access.NewGrouped[citizenslabels.JobAccess, *citizenslabels.JobAccess, access.DummyUserAccess[citizenslabels.AccessLevel], *access.DummyUserAccess[citizenslabels.AccessLevel], access.DummyQualificationAccess[citizenslabels.AccessLevel], *access.DummyQualificationAccess[citizenslabels.AccessLevel], citizenslabels.AccessLevel](
			p.DB,
			table.FivenetUserLabelsJob,
			&access.TargetTableColumns{
				ID:         table.FivenetUserLabelsJob.ID,
				DeletedAt:  table.FivenetUserLabelsJob.DeletedAt,
				CreatorID:  nil,
				CreatorJob: table.FivenetUserLabelsJob.Job,
			},
			access.NewJobs[citizenslabels.JobAccess, *citizenslabels.JobAccess, citizenslabels.AccessLevel](
				table.FivenetUserLabelsJobJobAccess,
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetUserLabelsJobJobAccess.ID,
						TargetID: table.FivenetUserLabelsJobJobAccess.TargetID,
						Access:   table.FivenetUserLabelsJobJobAccess.Access,
					},
					Job:          table.FivenetUserLabelsJobJobAccess.Job,
					MinimumGrade: table.FivenetUserLabelsJobJobAccess.MinimumGrade,
				},
				table.FivenetUserLabelsJobJobAccess.AS("job_access"),
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID: table.FivenetUserLabelsJobJobAccess.AS("job_access").ID,
						TargetID: table.FivenetUserLabelsJobJobAccess.AS(
							"job_access",
						).TargetID,
						Access: table.FivenetUserLabelsJobJobAccess.AS("job_access").Access,
					},
					Job: table.FivenetUserLabelsJobJobAccess.AS("job_access").Job,
					MinimumGrade: table.FivenetUserLabelsJobJobAccess.AS(
						"job_access",
					).MinimumGrade,
				},
			),
			nil,
			nil,
		),
	}

	access.RegisterAccess("citizen", &access.GroupedAccessAdapter{
		CanUserAccessTargetFn: func(ctx context.Context, targetId int64, userInfo *userinfo.UserInfo, access int32) (bool, error) {
			if !s.ps.Can(
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

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbcitizens.RegisterCitizensServiceServer(srv, s)
	pbcitizens.RegisterLabelsServiceServer(srv, s)
}
