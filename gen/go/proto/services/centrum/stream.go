package centrum

import (
	"context"
	"time"

	centrum "github.com/fivenet-app/fivenet/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
	errorscentrum "github.com/fivenet-app/fivenet/gen/go/proto/services/centrum/errors"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/utils/broker"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"go.uber.org/zap"
)

func (s *Server) TakeControl(ctx context.Context, req *TakeControlRequest) (*TakeControlResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "TakeControl",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	if err := s.state.DisponentSignOn(ctx, userInfo.Job, userInfo.UserId, req.Signon); err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	if req.Signon {
		auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)
	} else {
		auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &TakeControlResponse{}, nil
}

func (s *Server) sendLatestState(srv CentrumService_StreamServer, job string, userId int32) error {
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
	if err := srv.Send(&StreamResponse{
		Change: &StreamResponse_LatestState{
			LatestState: &LatestState{
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

func (s *Server) Stream(req *StreamRequest, srv CentrumService_StreamServer) error {
	userInfo := *auth.MustGetUserInfoFromContext(srv.Context())

	for {
		if err := s.sendLatestState(srv, userInfo.Job, userInfo.UserId); err != nil {
			return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		if err := s.stream(srv, userInfo.Job, userInfo.UserId); err != nil {
			return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		select {
		case <-srv.Context().Done():
			return nil

		case <-time.After(50 * time.Millisecond):
		}
	}
}

func (s *Server) getJobBroker(job string) (*broker.Broker[*StreamResponse], bool) {
	s.brokersMutex.RLock()
	defer s.brokersMutex.RUnlock()

	broker, ok := s.brokers[job]
	return broker, ok
}

func (s *Server) stream(srv CentrumService_StreamServer, job string, userId int32) error {
	s.logger.Debug("getting centrum job broker", zap.String("job", job), zap.Int32("user_id", userId))
	broker, ok := s.getJobBroker(job)
	if !ok {
		s.logger.Warn("no job broker found", zap.String("job", job), zap.Int32("user_id", userId))
		<-srv.Context().Done()
		return nil
	}

	s.logger.Debug("subscribing to centrum job broker", zap.String("job", job), zap.Int32("user_id", userId))
	stream := broker.Subscribe()
	defer broker.Unsubscribe(stream)
	s.logger.Debug("starting broker watch", zap.String("job", job), zap.Int32("user_id", userId))

	// Watch for events from message queue
	for {
		resp := &StreamResponse{}

		select {
		case <-srv.Context().Done():
			return nil

		case msg := <-stream:
			resp.Change = msg.Change
			if err := srv.Send(resp); err != nil {
				return err
			}
		}
	}
}
