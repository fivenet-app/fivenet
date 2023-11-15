package livemapper

import (
	"context"
	"database/sql"
	"time"

	"github.com/galexrt/fivenet/gen/go/proto/resources/livemap"
	users "github.com/galexrt/fivenet/gen/go/proto/resources/users"
	permslivemapper "github.com/galexrt/fivenet/gen/go/proto/services/livemapper/perms"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/server/audit"
	"github.com/galexrt/fivenet/pkg/tracker"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/puzpuzpuz/xsync/v3"
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
	trackedJobs []string
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

		markersCache: xsync.NewMapOf[string, []*livemap.Marker](),

		broker: broker,

		refreshTime: p.Config.Game.Livemap.RefreshTime,
		trackedJobs: p.Config.Game.Livemap.Jobs,
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

	dispatchesAttr, err := s.ps.Attr(userInfo, permslivemapper.LivemapperServicePerm, permslivemapper.LivemapperServiceStreamPerm, permslivemapper.LivemapperServiceStreamMarkersPermField)
	if err != nil {
		return ErrStreamFailed
	}
	playersAttr, err := s.ps.Attr(userInfo, permslivemapper.LivemapperServicePerm, permslivemapper.LivemapperServiceStreamPerm, permslivemapper.LivemapperServiceStreamPlayersPermField)
	if err != nil {
		return ErrStreamFailed
	}

	var dispatchesJobs []string
	if dispatchesAttr != nil {
		dispatchesJobs = dispatchesAttr.([]string)
	}
	if userInfo.SuperUser {
		dispatchesJobs = s.trackedJobs
	}

	var playersJobs map[string]int32
	if playersAttr != nil {
		playersJobs, _ = playersAttr.(map[string]int32)
	}

	if userInfo.SuperUser {
		playersJobs = map[string]int32{}
		for _, j := range s.trackedJobs {
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
