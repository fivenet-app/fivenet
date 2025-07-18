package centrum

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	pbcentrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	errorscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/errors"
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

	if err := s.dispatchers.SetUserState(ctx, userInfo.Job, userInfo.UserId, req.Signon); err != nil {
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

func (s *Server) UpdateDispatchers(ctx context.Context, req *pbcentrum.UpdateDispatchersRequest) (*pbcentrum.UpdateDispatchersResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "UpdateDispatchers",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	// Sign off any requested dispatchers
	for _, userId := range req.ToRemove {
		if err := s.dispatchers.SetUserState(ctx, userInfo.Job, userId, false); err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	// Retrieve the updated dispatchers list
	dispatchers, err := s.dispatchers.Get(ctx, userInfo.Job)
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	return &pbcentrum.UpdateDispatchersResponse{
		Dispatchers: &centrum.Dispatchers{
			Job:         userInfo.Job,
			Dispatchers: dispatchers.Dispatchers,
		},
	}, nil
}
