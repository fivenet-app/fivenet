package citizens

import (
	context "context"
	"database/sql"
	"math"
	"slices"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	users "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	pbcitizens "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/citizens"
	permscitizens "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/access"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/filestore"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/audit"
	"github.com/fivenet-app/fivenet/v2025/pkg/storage"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/fx"
	grpc "google.golang.org/grpc"
)

type Server struct {
	pbcitizens.CitizensServiceServer

	db       *sql.DB
	ps       perms.Permissions
	enricher *mstlystcdata.UserAwareEnricher
	aud      audit.IAuditer
	st       storage.IStorage
	appCfg   appconfig.IConfig
	cfg      *config.Config
	customDB config.CustomDB
	notifi   notifi.INotifi

	avatarHandler  *filestore.Handler[int32]
	mugshotHandler *filestore.Handler[int32]
}

type Params struct {
	fx.In

	DB        *sql.DB
	P         perms.Permissions
	Enricher  *mstlystcdata.UserAwareEnricher
	Aud       audit.IAuditer
	Config    *config.Config
	Storage   storage.IStorage
	AppConfig appconfig.IConfig
	Notifi    notifi.INotifi
}

func NewServer(p Params) *Server {
	tUserProps := table.FivenetUserProps

	avatarHandler := filestore.NewHandler(p.Storage, p.DB, tUserProps, tUserProps.UserID, tUserProps.AvatarFileID, 3<<20, func(parentId int32) jet.BoolExpression {
		return tUserProps.UserID.EQ(jet.Int32(parentId))
	}, filestore.UpdateJoinRow, true)
	mugshotHandler := filestore.NewHandler(p.Storage, p.DB, tUserProps, tUserProps.UserID, tUserProps.MugshotFileID, 3<<20, func(parentId int32) jet.BoolExpression {
		return tUserProps.UserID.EQ(jet.Int32(parentId))
	}, filestore.UpdateJoinRow, true)

	s := &Server{
		db:       p.DB,
		ps:       p.P,
		enricher: p.Enricher,
		aud:      p.Aud,
		st:       p.Storage,
		appCfg:   p.AppConfig,
		cfg:      p.Config,
		customDB: p.Config.Database.Custom,
		notifi:   p.Notifi,

		avatarHandler:  avatarHandler,
		mugshotHandler: mugshotHandler,
	}

	access.RegisterAccess("citizen", &access.GroupedAccessAdapter{
		CanUserAccessTargetFn: func(ctx context.Context, targetId uint64, userInfo *userinfo.UserInfo, access int32) (bool, error) {
			if !s.ps.Can(userInfo, permscitizens.CitizensServicePerm, permscitizens.CitizensServiceGetUserPerm) {
				return false, nil
			}

			if targetId > uint64(math.MaxInt32) {
				return false, nil // targetId is too large to fit in int32
			}
			userId := int32(targetId)
			tUser := tables.User().AS("user")

			// Retrieve user job from database
			stmt := tUser.
				SELECT(
					tUser.ID,
					tUser.Job,
				).
				FROM(tUser).
				WHERE(tUser.ID.EQ(jet.Int32(userId))).
				LIMIT(1)

			user := &users.User{}
			if err := stmt.QueryContext(ctx, s.db, user); err != nil {
				return false, err
			}

			if slices.Contains(s.appCfg.Get().JobInfo.PublicJobs, user.Job) ||
				slices.Contains(s.appCfg.Get().JobInfo.HiddenJobs, user.Job) {
				// Make sure user has permission to see that grade
				check, err := s.checkIfUserCanAccess(userInfo, user.Job, user.JobGrade)
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
}

func (s *Server) GetPermsRemap() map[string]string {
	return pbcitizens.PermsRemap
}
