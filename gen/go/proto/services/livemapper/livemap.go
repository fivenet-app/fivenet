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
	"github.com/galexrt/fivenet/pkg/grpc/errswrap"
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
	grpc "google.golang.org/grpc"
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
	tracker  tracker.ITracker
	auditer  audit.IAuditer

	markersCache *xsync.MapOf[string, []*livemap.MarkerMarker]

	broker *utils.Broker[*livemap.UsersUpdateEvent]

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
	Tracker  tracker.ITracker
	Audit    audit.IAuditer
}

func NewServer(p Params) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	broker := utils.NewBroker[*livemap.UsersUpdateEvent](ctx)
	s := &Server{
		ctx:    ctx,
		logger: p.Logger,

		tracer:   p.TP.Tracer("livemapper-cache"),
		db:       p.DB,
		ps:       p.Perms,
		enricher: p.Enricher,
		tracker:  p.Tracker,
		auditer:  p.Audit,

		markersCache: xsync.NewMapOf[string, []*livemap.MarkerMarker](),

		broker: broker,

		refreshTime: p.Config.Game.Livemap.RefreshTime,
		trackedJobs: p.Config.Game.Livemap.Jobs,
	}

	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		go broker.Start()
		go s.start()

		go func() {
			subCh := s.tracker.Subscribe()
			defer s.tracker.Unsubscribe(subCh)

			for {
				select {
				case <-s.ctx.Done():
					return

				case msg := <-subCh:
					if s.broker.SubCount() <= 0 {
						continue
					}

					s.logger.Debug("received msg from user tracker")
					s.broker.Publish(msg)
				}
			}
		}()

		return nil
	}))
	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		return nil
	}))

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterLivemapperServiceServer(srv, s)
}

func (s *Server) start() {
	for {
		s.refreshCache()

		select {
		case <-s.ctx.Done():
			return
		case <-time.After(4 * time.Second):
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

	s.logger.Debug("starting livemap stream")
	markerJobsAttr, err := s.ps.Attr(userInfo, permslivemapper.LivemapperServicePerm, permslivemapper.LivemapperServiceStreamPerm, permslivemapper.LivemapperServiceStreamMarkersPermField)
	if err != nil {
		return errswrap.NewError(ErrStreamFailed, err)
	}
	userJobsAttr, err := s.ps.Attr(userInfo, permslivemapper.LivemapperServicePerm, permslivemapper.LivemapperServiceStreamPerm, permslivemapper.LivemapperServiceStreamPlayersPermField)
	if err != nil {
		return errswrap.NewError(ErrStreamFailed, err)
	}

	var markersJobs []string
	if markerJobsAttr != nil {
		markersJobs = markerJobsAttr.([]string)
	}
	if userInfo.SuperUser {
		s.markersCache.Range(func(job string, _ []*livemap.MarkerMarker) bool {
			markersJobs = append(markersJobs, job)
			return true
		})
		markersJobs = utils.RemoveDuplicates(markersJobs)
	}

	var usersJobs map[string]int32
	if userJobsAttr != nil {
		usersJobs, _ = userJobsAttr.(map[string]int32)
	}

	if userInfo.SuperUser {
		usersJobs = map[string]int32{}
		for _, j := range s.trackedJobs {
			usersJobs[j] = -1
		}
	}

	resp := &StreamResponse{
		Data: &StreamResponse_Jobs{},
	}
	if len(markersJobs) == 0 && len(usersJobs) == 0 {
		if err := srv.Send(resp); err != nil {
			return err
		}

		return nil
	}

	// Prepare jobs for client response
	jobs := &StreamResponse_Jobs{
		Jobs: &JobsList{
			Markers: make([]*users.Job, len(markersJobs)),
			Users:   []*users.Job{},
		},
	}
	for i := 0; i < len(markersJobs); i++ {
		jobs.Jobs.Markers[i] = &users.Job{
			Name: markersJobs[i],
		}
		s.enricher.EnrichJobName(jobs.Jobs.Markers[i])
	}
	for job := range usersJobs {
		j := &users.Job{
			Name: job,
		}
		s.enricher.EnrichJobName(j)
		jobs.Jobs.Users = append(jobs.Jobs.Users, j)
	}
	resp.Data = jobs

	if err := srv.Send(resp); err != nil {
		return err
	}

	// TODO send out marker markers changes only
	if end, err := s.sendMarkerMarkers(srv, markersJobs, userInfo); end || err != nil {
		return err
	}

	if end, err := s.sendChunkedUserMarkers(srv, usersJobs, userInfo); end || err != nil {
		return err
	}

	ticker := time.NewTicker(4 * s.refreshTime)
	defer ticker.Stop()

	signalCh := s.broker.Subscribe()
	defer s.broker.Unsubscribe(signalCh)

	for {
		select {
		case <-srv.Context().Done():
			return nil

		case <-ticker.C:
			if end, err := s.sendMarkerMarkers(srv, markersJobs, userInfo); end || err != nil {
				return err
			}

		case <-time.After(s.refreshTime):
			if end, err := s.sendChunkedUserMarkers(srv, usersJobs, userInfo); end || err != nil {
				return err
			}
		}
	}
}

func (s *Server) sendChunkedUserMarkers(srv LivemapperService_StreamServer, usersJobs map[string]int32, userInfo *userinfo.UserInfo) (bool, error) {
	// Send out chunked current users
	chunkSize := 25
	userMarkers, _, err := s.getUserLocations(usersJobs, userInfo)
	if err != nil {
		return true, errswrap.NewError(ErrStreamFailed, err)
	}

	parts := int32(len(userMarkers) / chunkSize)
	for chunkSize < len(userMarkers) {
		resp := &StreamResponse{
			Data: &StreamResponse_Users{
				Users: &UserMarkersUpdates{
					Users: userMarkers[0:chunkSize:chunkSize],
					Part:  parts,
				},
			},
		}
		parts--

		if err := srv.Send(resp); err != nil {
			return true, err
		}

		userMarkers = userMarkers[chunkSize:]
		select {
		case <-srv.Context().Done():
			return true, nil
		case <-time.After(175 * time.Millisecond):
		}
	}

	if len(userMarkers) > 0 {
		resp := &StreamResponse{
			Data: &StreamResponse_Users{
				Users: &UserMarkersUpdates{
					Users: userMarkers,
					Part:  0,
				},
			},
		}

		if err := srv.Send(resp); err != nil {
			return true, err
		}
	}

	return false, nil
}

func (s *Server) sendMarkerMarkers(srv LivemapperService_StreamServer, jobs []string, userInfo *userinfo.UserInfo) (bool, error) {
	markers, err := s.getMarkerMarkers(jobs)
	if err != nil {
		return true, errswrap.NewError(ErrStreamFailed, err)
	}

	// Send current markers
	resp := &StreamResponse{
		Data: &StreamResponse_Markers{
			Markers: &MarkerMarkersUpdates{
				Markers: markers,
			},
		},
	}
	if err := srv.Send(resp); err != nil {
		return true, err
	}

	return false, nil
}

func (s *Server) getUserLocations(jobs map[string]int32, userInfo *userinfo.UserInfo) ([]*livemap.UserMarker, bool, error) {
	ds := []*livemap.UserMarker{}

	found := false
	if userInfo.SuperUser {
		found = true
	}

	for job, grade := range jobs {
		markers, ok := s.tracker.GetUsersByJob(job)
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
