//nolint:forcetypeassert // We know the type is correct as it is validated and unmarshaled from/to the type.
package centrum

import (
	"context"
	"errors"
	"fmt"
	"time"

	centrumdispatchers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatchers"
	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	centrumsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/settings"
	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbcentrum "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	errorscentrum "github.com/fivenet-app/fivenet/v2026/services/centrum/errors"
	centrumutils "github.com/fivenet-app/fivenet/v2026/services/centrum/utils"
	"google.golang.org/protobuf/proto"
)

const feedFetch = 16

func (s *Server) sendHandshake(
	ctx context.Context,
	srv pbcentrum.CentrumService_StreamServer,
	userJob string,
	acls *centrumsettings.EffectiveAccess,
) error {
	settings, err := s.settings.Get(ctx, userJob)
	if err != nil {
		return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	if err := srv.Send(&pbcentrum.StreamResponse{
		Change: &pbcentrum.StreamResponse_Handshake{
			Handshake: &pbcentrum.StreamHandshake{
				ServerTime: timestamp.Now(),
				Settings:   settings,
				Access:     acls,
			},
		},
	}); err != nil {
		return err
	}

	return nil
}

func (s *Server) sendLatestState(
	ctx context.Context,
	srv pbcentrum.CentrumService_StreamServer,
	userInfo *userinfo.UserInfo,
	acls *centrumsettings.EffectiveAccess,
	jobList []string,
) error {
	// Dispatchers
	dispatchers := &centrumdispatchers.JobDispatchers{}
	for _, j := range acls.GetDispatches().GetJobs() {
		dispos, err := s.dispatchers.Get(ctx, j.GetJob())
		if err != nil {
			return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
		dispatchers.Dispatchers = append(dispatchers.Dispatchers, dispos)
	}

	// Own unit ID
	ownUnitMapping, ok, err := s.tracker.GetUserMapping(userInfo.GetUserId())
	if err != nil {
		return fmt.Errorf("failed to get own unit mapping: %w", err)
	}
	var pOwnUnitId *int64
	if ok && ownUnitMapping != nil && ownUnitMapping.UnitId != nil &&
		ownUnitMapping.GetUnitId() > 0 {
		pOwnUnitId = ownUnitMapping.UnitId
	}

	// Retrieve units and dispatches
	units := s.units.List(ctx, jobList)

	dispatches := s.dispatches.Filter(ctx, jobList, nil, []centrumdispatches.StatusDispatch{
		centrumdispatches.StatusDispatch_STATUS_DISPATCH_ARCHIVED,
		centrumdispatches.StatusDispatch_STATUS_DISPATCH_CANCELLED,
		centrumdispatches.StatusDispatch_STATUS_DISPATCH_COMPLETED,
		centrumdispatches.StatusDispatch_STATUS_DISPATCH_DELETED,
	})

	// Send initial state to client
	if err := srv.Send(&pbcentrum.StreamResponse{
		Change: &pbcentrum.StreamResponse_LatestState{
			LatestState: &pbcentrum.LatestState{
				Dispatchers: dispatchers,
				OwnUnitId:   pOwnUnitId,
				Units:       units,
				Dispatches:  dispatches,
			},
		},
	}); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stream(
	req *pbcentrum.StreamRequest,
	srv pbcentrum.CentrumService_StreamServer,
) error {
	userInfo := auth.MustGetUserInfoFromContext(srv.Context()).Clone()
	// Subscribe before reading current state so a change committed during the
	// database read is queued and forces a follow-up authorization snapshot.
	userInfoChanges := s.userinfoChanges.SubscribeUserInfoChanges()
	defer func() {
		s.userinfoChanges.UnsubscribeUserInfoChanges(userInfoChanges)
	}()
	if err := s.refreshStreamUserInfo(srv.Context(), userInfo); err != nil {
		return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}
	if err := s.waitForReady(srv.Context()); err != nil {
		if protoutils.IsContextCanceled(err) {
			return nil
		}
		return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	s.metrics.IncActiveStreams()
	defer s.metrics.DecActiveStreams()

	feed := s.feedBroker.Subscribe()
	defer func() {
		s.feedBroker.Unsubscribe(feed)
	}()

	for {
		// Recalculate access after a user-info or own-settings change before
		// sending another snapshot or forwarding any further feed event.
		jobList, acls, err := s.settings.GetAccessList(
			srv.Context(),
			userInfo.GetJob(),
			userInfo.GetJobGrade(),
		)
		if err != nil {
			return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		// Subscribe first, then capture the boundary. Events published while the
		// snapshot is sent stay queued and are forwarded afterwards.
		snapshotSequence := s.feedSequence.Load()
		if err := s.sendHandshake(srv.Context(), srv, userInfo.GetJob(), acls); err != nil {
			if protoutils.IsContextCanceled(err) {
				return nil
			}
			return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		if err := s.sendLatestState(srv.Context(), srv, userInfo, acls, jobList); err != nil {
			if protoutils.IsContextCanceled(err) {
				return nil
			}
			return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		if err := s.stream(
			srv.Context(),
			srv,
			userInfo,
			jobList,
			feed,
			userInfoChanges,
			snapshotSequence,
		); err != nil {
			if errors.Is(err, errFeedClosed) {
				s.metrics.IncFeedResync("slow_subscriber")
				feed = s.feedBroker.Subscribe()
				continue
			}
			if errors.Is(err, errFeedResync) {
				s.metrics.IncFeedResync("sequence_gap")
				continue
			}
			if errors.Is(err, errAccessChanged) || errors.Is(err, errUserInfoChanged) {
				if errors.Is(err, errAccessChanged) {
					s.metrics.IncFeedResync("access_change")
				} else {
					s.metrics.IncFeedResync("userinfo_change")
				}
				continue
			}
			if errors.Is(err, errUserInfoResync) {
				s.metrics.IncFeedResync("userinfo_subscriber_resync")
				replacement := s.userinfoChanges.SubscribeUserInfoChanges()
				if err := s.refreshStreamUserInfo(srv.Context(), userInfo); err != nil {
					s.userinfoChanges.UnsubscribeUserInfoChanges(replacement)
					return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
				}
				s.userinfoChanges.UnsubscribeUserInfoChanges(userInfoChanges)
				userInfoChanges = replacement
				continue
			}
			return err
		}

		select {
		case <-srv.Context().Done():
			return nil

		case <-time.After(50 * time.Millisecond):
		}
	}
}

func (s *Server) refreshStreamUserInfo(ctx context.Context, current *userinfo.UserInfo) error {
	updated, err := s.userinfo.GetUserInfo(ctx, current.GetUserId())
	if err != nil {
		return fmt.Errorf("failed to retrieve current user info. %w", err)
	}

	current.Job = updated.Job
	current.JobGrade = updated.JobGrade
	return nil
}

type feedCfg struct {
	StreamName string
	Bucket     string
	NoWildcard bool                                                                  // if true, the subjects are not a wildcard bucket (e.g., centrum_dispatchers.JOB)
	Unmarshal  func(ctx context.Context, s *Server, b []byte) (proto.Message, error) // bucket -> concrete proto
	WrapPut    func(proto.Message) *pbcentrum.StreamResponse
	WrapDelete func(key string) *pbcentrum.StreamResponse
}

var feeds = []feedCfg{
	{
		StreamName: "centrum_settings",
		Bucket:     "centrum_settings",
		NoWildcard: true,
		Unmarshal: func(ctx context.Context, _ *Server, b []byte) (proto.Message, error) {
			var u centrumsettings.Settings
			return &u, proto.Unmarshal(b, &u)
		},
		WrapPut: func(m proto.Message) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_Settings{Settings: m.(*centrumsettings.Settings)},
			}
		},
		WrapDelete: func(key string) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_SettingsDeleted{SettingsDeleted: key},
			}
		},
	},
	{
		StreamName: "centrum_dispatchers",
		Bucket:     "centrum_dispatchers",
		NoWildcard: true,
		Unmarshal: func(ctx context.Context, _ *Server, b []byte) (proto.Message, error) {
			var u centrumdispatchers.Dispatchers
			return &u, proto.Unmarshal(b, &u)
		},
		WrapPut: func(m proto.Message) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_Dispatchers{
					Dispatchers: m.(*centrumdispatchers.Dispatchers),
				},
			}
		},
		WrapDelete: func(key string) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_Dispatchers{
					Dispatchers: &centrumdispatchers.Dispatchers{Job: key},
				},
			}
		},
	},
	{
		StreamName: "centrum_units",
		Bucket:     "centrum_units.job",
		Unmarshal: func(ctx context.Context, s *Server, b []byte) (proto.Message, error) {
			var d common.IDMapping
			if err := proto.Unmarshal(b, &d); err != nil {
				return nil, fmt.Errorf("failed to unmarshal unit id mapping. %w", err)
			}
			return s.units.Get(ctx, d.GetId())
		},
		WrapPut: func(m proto.Message) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_UnitUpdated{UnitUpdated: m.(*centrumunits.Unit)},
			}
		},
		WrapDelete: func(key string) *pbcentrum.StreamResponse {
			id, err := centrumutils.ExtractID(key)
			if err != nil {
				return nil
			}

			return &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_UnitDeleted{
					UnitDeleted: id,
				},
			}
		},
	},
	{
		StreamName: "centrum_dispatches",
		Bucket:     "centrum_dispatches.job",
		Unmarshal: func(ctx context.Context, s *Server, b []byte) (proto.Message, error) {
			var d common.IDMapping
			if err := proto.Unmarshal(b, &d); err != nil {
				return nil, fmt.Errorf("failed to unmarshal dispatch id mapping. %w", err)
			}
			return s.dispatches.Get(ctx, d.GetId())
		},
		WrapPut: func(m proto.Message) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_DispatchUpdated{
					DispatchUpdated: m.(*centrumdispatches.Dispatch),
				},
			}
		},
		WrapDelete: func(key string) *pbcentrum.StreamResponse {
			id, err := centrumutils.ExtractID(key)
			if err != nil {
				return nil
			}

			return &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_DispatchDeleted{
					DispatchDeleted: id,
				},
			}
		},
	},
}
