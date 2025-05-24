package centrum

import (
	"context"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	centrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	pbcentrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	errorscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/errors"
	"go.uber.org/zap"
)

func (s *Server) TakeControl(ctx context.Context, req *pbcentrum.TakeControlRequest) (*pbcentrum.TakeControlResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "TakeControl",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	if err := s.state.DisponentSignOn(ctx, userInfo.Job, userInfo.UserId, req.Signon); err != nil {
		return nil, err
	}

	if req.Signon {
		auditEntry.State = audit.EventType_EVENT_TYPE_CREATED
	} else {
		auditEntry.State = audit.EventType_EVENT_TYPE_DELETED
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbcentrum.TakeControlResponse{}, nil
}

func (s *Server) sendLatestState(srv pbcentrum.CentrumService_StreamServer, job string, userId int32) error {
	ctx := srv.Context()

	settings := s.state.GetSettings(ctx, job)
	disponents, _ := s.state.GetDisponents(ctx, job)
	ownUnitId, _ := s.state.GetUserUnitID(ctx, userId)
	var pOwnUnitId *uint64
	if ownUnitId > 0 {
		pOwnUnitId = &ownUnitId
	}

	units, _ := s.state.ListUnits(ctx, job)
	dispatches := s.state.FilterDispatches(ctx, job, nil, []centrum.StatusDispatch{
		centrum.StatusDispatch_STATUS_DISPATCH_ARCHIVED,
		centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED,
		centrum.StatusDispatch_STATUS_DISPATCH_COMPLETED,
	})

	// Send initial state to client
	if err := srv.Send(&pbcentrum.StreamResponse{
		Change: &pbcentrum.StreamResponse_LatestState{
			LatestState: &pbcentrum.LatestState{
				ServerTime: timestamp.Now(),
				Settings:   settings,
				Disponents: disponents,
				OwnUnitId:  pOwnUnitId,
				Units:      units,
				Dispatches: dispatches,
			},
		},
	}); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stream(req *pbcentrum.StreamRequest, srv pbcentrum.CentrumService_StreamServer) error {
	userInfo := *auth.MustGetUserInfoFromContext(srv.Context())

	for {
		if err := s.sendLatestState(srv, userInfo.Job, userInfo.UserId); err != nil {
			return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		if err := s.stream(srv, userInfo.Job, userInfo.UserId); err != nil {
			return err
		}

		select {
		case <-srv.Context().Done():
			return nil

		case <-time.After(50 * time.Millisecond):
		}
	}
}

func (s *Server) stream(srv pbcentrum.CentrumService_StreamServer, job string, userId int32) error {
	s.logger.Debug("getting centrum job broker", zap.String("job", job), zap.Int32("user_id", userId))
	broker, ok := s.brokers.GetJobBroker(job)
	if !ok {
		return errorscentrum.ErrDisabled
	}

	stream := broker.Subscribe()
	defer broker.Unsubscribe(stream)
	s.logger.Debug("starting broker watch", zap.String("job", job), zap.Int32("user_id", userId))

	// Watch for events from message queue
	for {
		select {
		case <-srv.Context().Done():
			return nil

		case msg, more := <-stream:
			if !more {
				return errorscentrum.ErrDisabled
			}

			resp := &pbcentrum.StreamResponse{
				Change: msg.Change,
			}
			if err := srv.Send(resp); err != nil {
				return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}
		}
	}
}
