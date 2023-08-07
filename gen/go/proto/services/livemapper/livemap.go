package livemapper

import (
	"context"
	"database/sql"
	"time"

	jobs "github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
	"github.com/galexrt/fivenet/gen/go/proto/resources/livemap"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/pkg/utils/syncx"
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

	logger *zap.Logger
	tracer trace.Tracer
	db     *sql.DB
	p      perms.Permissions
	c      *mstlystcdata.Enricher

	dispatchesCache syncx.Map[string, []*livemap.DispatchMarker]
	usersCache      syncx.Map[string, []*livemap.UserMarker]

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
}

func NewServer(p Params) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	broker := utils.NewBroker[interface{}](ctx)

	p.LC.Append(fx.StartHook(func(_ context.Context) {
		go broker.Start()
	}))
	p.LC.Append(fx.StopHook(func(_ context.Context) {
		cancel()
	}))

	return &Server{
		logger: p.Logger,

		tracer: p.TP.Tracer("livemapper-cache"),
		db:     p.DB,
		p:      p.Perms,
		c:      p.Enricher,

		dispatchesCache: syncx.Map[string, []*livemap.DispatchMarker]{},
		usersCache:      syncx.Map[string, []*livemap.UserMarker]{},

		broker: broker,

		refreshTime: p.Config.Cache.RefreshTime,
		visibleJobs: p.Config.Game.Livemap.Jobs,
	}
}

func (s *Server) Stream(req *StreamRequest, srv LivemapperService_StreamServer) error {
	userInfo := auth.MustGetUserInfoFromContext(srv.Context())

	dispatchesAttr, err := s.p.Attr(userInfo, LivemapperServicePerm, LivemapperServiceStreamPerm, LivemapperServiceStreamDispatchesPermField)
	if err != nil {
		return ErrStreamFailed
	}
	playersAttr, err := s.p.Attr(userInfo, LivemapperServicePerm, LivemapperServiceStreamPerm, LivemapperServiceStreamPlayersPermField)
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
	resp.JobsDispatches = make([]*jobs.Job, len(dispatchesJobs))
	for i := 0; i < len(dispatchesJobs); i++ {
		resp.JobsDispatches[i] = &jobs.Job{
			Name: dispatchesJobs[i],
		}
		s.c.EnrichJobName(resp.JobsDispatches[i])
	}
	resp.JobsUsers = []*jobs.Job{}
	for job := range playersJobs {
		j := &jobs.Job{
			Name: job,
		}
		s.c.EnrichJobName(j)
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

		dispatchMarkers, err := s.getUserDispatches(dispatchesJobs)
		if err != nil {
			return ErrStreamFailed
		}
		resp.Dispatches = dispatchMarkers

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
		markers, ok := s.usersCache.Load(job)
		if !ok {
			continue
		}

		for i := 0; i < len(markers); i++ {
			// SuperUser returns grade as `-1`, job has access to that grade or it is the user itself
			if grade == -1 || (markers[i].User.JobGrade <= grade || markers[i].User.UserId == userInfo.UserId) {
				ds = append(ds, markers[i])
			}

			// If the user is found in the list of user markers, set found state
			if !found && (userInfo.Job == job && markers[i].Marker.Id == userInfo.UserId) {
				found = true
			}
		}
	}

	if found {
		return ds, true, nil
	}

	return nil, false, nil
}

func (s *Server) getUserDispatches(jobs []string) ([]*livemap.DispatchMarker, error) {
	ds := []*livemap.DispatchMarker{}

	for _, job := range jobs {
		markers, ok := s.dispatchesCache.Load(job)
		if !ok {
			continue
		}

		ds = append(ds, markers...)
	}

	return ds, nil
}
