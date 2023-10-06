package livemapper

import (
	"context"
	"database/sql"
	"time"

	"github.com/galexrt/fivenet/gen/go/proto/resources/livemap"
	users "github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/server/audit"
	"github.com/galexrt/fivenet/pkg/tracker"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/puzpuzpuz/xsync/v2"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DispatchMarkerLimit = 60
)

var (
	ErrStreamFailed = status.Error(codes.Internal, "errors.LivemapperService.ErrStreamFailed")
)

var (
	tUsers = table.Users.AS("user")
)

type Server struct {
	LivemapperServiceServer

	ctx      context.Context
	logger   *zap.Logger
	tracer   trace.Tracer
	db       *sql.DB
	ps       perms.Permissions
	enricher *mstlystcdata.Enricher
	tracker  *tracker.Tracker
	auditer  audit.IAuditer

	markersCache *xsync.MapOf[string, []*livemap.Marker]

	broker *utils.Broker[interface{}]

	refreshTime time.Duration
	visibleJobs []string
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger   *zap.Logger
	TP       *tracesdk.TracerProvider
	DB       *sql.DB
	Perms    perms.Permissions
	Enricher *mstlystcdata.Enricher
	Config   *config.Config
	Tracker  *tracker.Tracker
	Audit    audit.IAuditer
}

func NewServer(p Params) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	broker := utils.NewBroker[interface{}](ctx)
	s := &Server{
		ctx:    ctx,
		logger: p.Logger,

		tracer:   p.TP.Tracer("livemapper-cache"),
		db:       p.DB,
		ps:       p.Perms,
		enricher: p.Enricher,
		tracker:  p.Tracker,
		auditer:  p.Audit,

		markersCache: xsync.NewMapOf[[]*livemap.Marker](),

		broker: broker,

		refreshTime: p.Config.Game.Livemap.RefreshTime,
		visibleJobs: p.Config.Game.Livemap.Jobs,
	}

	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		go broker.Start()
		go s.start()
		return nil
	}))
	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()
		return nil
	}))

	return s
}

func (s *Server) start() {
	for {
		s.refreshCache()

		select {
		case <-s.ctx.Done():
			s.broker.Stop()
			return
		case <-time.After(s.refreshTime):
		}
	}
}

func (s *Server) refreshCache() {
	ctx, span := s.tracer.Start(s.ctx, "livemap-refresh-cache")
	defer span.End()

	if err := s.refreshMarkers(ctx); err != nil {
		s.logger.Error("failed to refresh livemap markers cache", zap.Error(err))
	}

	s.broker.Publish(nil)
}

func (s *Server) Stream(req *StreamRequest, srv LivemapperService_StreamServer) error {
	userInfo := auth.MustGetUserInfoFromContext(srv.Context())

	dispatchesAttr, err := s.ps.Attr(userInfo, LivemapperServicePerm, LivemapperServiceStreamPerm, LivemapperServiceStreamMarkersPermField)
	if err != nil {
		return ErrStreamFailed
	}
	playersAttr, err := s.ps.Attr(userInfo, LivemapperServicePerm, LivemapperServiceStreamPerm, LivemapperServiceStreamPlayersPermField)
	if err != nil {
		return ErrStreamFailed
	}

	var dispatchesJobs []string
	if dispatchesAttr != nil {
		dispatchesJobs = dispatchesAttr.([]string)
	}
	if userInfo.SuperUser {
		dispatchesJobs = s.visibleJobs
	}

	var playersJobs map[string]int32
	if playersAttr != nil {
		playersJobs, _ = playersAttr.(map[string]int32)
	}

	if userInfo.SuperUser {
		playersJobs = map[string]int32{}
		for _, j := range s.visibleJobs {
			playersJobs[j] = -1
		}
	}

	resp := &StreamResponse{}

	if len(dispatchesJobs) == 0 && len(playersJobs) == 0 {
		if err := srv.Send(resp); err != nil {
			return err
		}

		return nil
	}

	// Add jobs to list visible jobs list
	resp.JobsMarkers = make([]*users.Job, len(dispatchesJobs))
	for i := 0; i < len(dispatchesJobs); i++ {
		resp.JobsMarkers[i] = &users.Job{
			Name: dispatchesJobs[i],
		}
		s.enricher.EnrichJobName(resp.JobsMarkers[i])
	}
	resp.JobsUsers = []*users.Job{}
	for job := range playersJobs {
		j := &users.Job{
			Name: job,
		}
		s.enricher.EnrichJobName(j)
		resp.JobsUsers = append(resp.JobsUsers, j)
	}

	signalCh := s.broker.Subscribe()
	defer s.broker.Unsubscribe(signalCh)

	for {
		userMarkers, _, err := s.getUserLocations(playersJobs, userInfo)
		if err != nil {
			return ErrStreamFailed
		}
		resp.Users = userMarkers

		markers, err := s.getMarkers(dispatchesJobs)
		if err != nil {
			return ErrStreamFailed
		}
		resp.Markers = markers

		if err := srv.Send(resp); err != nil {
			return err
		}

		select {
		case <-srv.Context().Done():
			return nil
		case <-signalCh:
		}
	}
}

func (s *Server) getUserLocations(jobs map[string]int32, userInfo *userinfo.UserInfo) ([]*livemap.UserMarker, bool, error) {
	ds := []*livemap.UserMarker{}

	found := false
	if userInfo.SuperUser {
		found = true
	}

	for job, grade := range jobs {
		markers, ok := s.tracker.GetUsers(job)
		if !ok {
			continue
		}

		markers.Range(func(key int32, marker *livemap.UserMarker) bool {
			// SuperUser returns grade as `-1`, job has access to that grade or it is the user itself
			if grade == -1 || (marker.User.JobGrade <= grade || key == userInfo.UserId) {
				ds = append(ds, marker)
			}

			// If the user is found in the list of user markers, set found state
			if !found && (userInfo.Job == job && key == userInfo.UserId) {
				found = true
			}

			return true
		})
	}

	if found {
		return ds, true, nil
	}

	return nil, false, nil
}

func (s *Server) getMarkers(jobs []string) ([]*livemap.Marker, error) {
	ds := []*livemap.Marker{}

	for _, job := range jobs {
		markers, ok := s.markersCache.Load(job)
		if !ok {
			continue
		}

		ds = append(ds, markers...)
	}

	return ds, nil
}

func (s *Server) refreshMarkers(ctx context.Context) error {
	tUsers := tUsers.AS("creator")
	tMarkers := tMarkers.AS("marker")
	stmt := tMarkers.
		SELECT(
			tMarkers.ID.AS("markerinfo.id"),
			tMarkers.Job.AS("markerinfo.job"),
			tMarkers.Name.AS("markerinfo.name"),
			tMarkers.Description.AS("markerinfo.description"),
			tMarkers.X.AS("markerinfo.x"),
			tMarkers.Y.AS("markerinfo.y"),
			tMarkers.Postal.AS("markerinfo.postal"),
			tMarkers.Color.AS("markerinfo.color"),
			tMarkers.Icon.AS("markerinfo.icon"),
			tMarkers.MarkerType,
			tMarkers.MarkerData,
			tMarkers.CreatorID,
			tUsers.ID,
			tUsers.Identifier,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Firstname,
			tUsers.Lastname,
		).
		FROM(
			tMarkers.
				LEFT_JOIN(tUsers,
					tMarkers.CreatorID.EQ(tUsers.ID),
				),
		).
		WHERE(
			tMarkers.CreatedAt.GT_EQ(
				jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(60, jet.MINUTE)),
			),
		)

	var dest []*livemap.Marker
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return err
	}

	markers := map[string][]*livemap.Marker{}
	for _, job := range s.visibleJobs {
		markers[job] = []*livemap.Marker{}
	}
	for _, m := range dest {
		if _, ok := markers[m.Info.Job]; !ok {
			markers[m.Info.Job] = []*livemap.Marker{}
		}

		markers[m.Info.Job] = append(markers[m.Info.Job], m)
	}

	for job, ms := range markers {
		if len(ms) == 0 {
			s.markersCache.Delete(job)
		} else {
			s.markersCache.Store(job, ms)
		}
	}

	return nil
}
