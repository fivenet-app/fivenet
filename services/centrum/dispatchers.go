package centrum

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	pbcentrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2025/pkg/grpc/interceptors/audit"
	errorscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/errors"
)

func (s *Server) TakeControl(
	ctx context.Context,
	req *pbcentrum.TakeControlRequest,
) (*pbcentrum.TakeControlResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if err := s.dispatchers.SetUserState(ctx, userInfo.GetJob(), userInfo.GetUserId(), req.GetSignon()); err != nil {
		return nil, err
	}

	if req.GetSignon() {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)
	} else {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return &pbcentrum.TakeControlResponse{}, nil
}

func (s *Server) UpdateDispatchers(
	ctx context.Context,
	req *pbcentrum.UpdateDispatchersRequest,
) (*pbcentrum.UpdateDispatchersResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Sign off any requested dispatchers
	for _, userId := range req.GetToRemove() {
		if err := s.dispatchers.SetUserState(ctx, userInfo.GetJob(), userId, false); err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	// Retrieve the updated dispatchers list
	dispatchers, err := s.dispatchers.Get(ctx, userInfo.GetJob())
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	return &pbcentrum.UpdateDispatchersResponse{
		Dispatchers: &centrum.Dispatchers{
			Job:         userInfo.GetJob(),
			Dispatchers: dispatchers.GetDispatchers(),
		},
	}, nil
}
