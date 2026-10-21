package settings

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	pbsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorssettings "github.com/fivenet-app/fivenet/v2026/services/settings/errors"
	settingsstore "github.com/fivenet-app/fivenet/v2026/stores/settings"
)

// ListPropertyScopes is intentionally small today. Adding a new supported
// entity kind later is a server change, not an API or UI redesign.
func (s *Server) ListPropertyScopes(
	ctx context.Context,
	_ *pbsettings.ListPropertyScopesRequest,
) (*pbsettings.ListPropertyScopesResponse, error) {
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)
	return &pbsettings.ListPropertyScopesResponse{
		Scopes: []properties.PropertyScope{
			properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
			properties.PropertyScope_PROPERTY_SCOPE_USER,
		},
	}, nil
}

func (s *Server) ListPropertyDefinitions(
	ctx context.Context,
	req *pbsettings.ListPropertyDefinitionsRequest,
) (*pbsettings.ListPropertyDefinitionsResponse, error) {
	definitions, err := s.store.ListPropertyDefinitions(
		ctx,
		req.GetScope(),
		req.GetIncludeDeleted(),
	)
	if err != nil {
		return nil, propertyDefinitionError(err)
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)
	groups, err := s.store.ListPropertyGroups(ctx, req.GetScope(), req.GetIncludeDeleted())
	if err != nil {
		return nil, propertyDefinitionError(err)
	}
	return &pbsettings.ListPropertyDefinitionsResponse{
		Definitions: definitions,
		Groups:      groups,
	}, nil
}

func (s *Server) CreatePropertyGroup(
	ctx context.Context,
	req *pbsettings.CreatePropertyGroupRequest,
) (*pbsettings.CreatePropertyGroupResponse, error) {
	group, err := s.store.CreatePropertyGroup(ctx, req.GetGroup())
	if err != nil {
		return nil, propertyDefinitionError(err)
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)
	return &pbsettings.CreatePropertyGroupResponse{Group: group}, nil
}

func (s *Server) UpdatePropertyGroup(
	ctx context.Context,
	req *pbsettings.UpdatePropertyGroupRequest,
) (*pbsettings.UpdatePropertyGroupResponse, error) {
	group, err := s.store.UpdatePropertyGroup(ctx, req.GetGroup())
	if err != nil {
		return nil, propertyDefinitionError(err)
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	return &pbsettings.UpdatePropertyGroupResponse{Group: group}, nil
}

func (s *Server) DeletePropertyGroup(
	ctx context.Context,
	req *pbsettings.DeletePropertyGroupRequest,
) (*pbsettings.DeletePropertyGroupResponse, error) {
	if err := s.store.DeletePropertyGroup(ctx, req.GetScope(), req.GetKey()); err != nil {
		return nil, propertyDefinitionError(err)
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	return &pbsettings.DeletePropertyGroupResponse{}, nil
}

func (s *Server) RestorePropertyGroup(
	ctx context.Context,
	req *pbsettings.RestorePropertyGroupRequest,
) (*pbsettings.RestorePropertyGroupResponse, error) {
	group, err := s.store.RestorePropertyGroup(ctx, req.GetScope(), req.GetKey())
	if err != nil {
		return nil, propertyDefinitionError(err)
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	return &pbsettings.RestorePropertyGroupResponse{Group: group}, nil
}

func (s *Server) ReorderPropertyGroups(
	ctx context.Context,
	req *pbsettings.ReorderPropertyGroupsRequest,
) (*pbsettings.ReorderPropertyGroupsResponse, error) {
	if err := s.store.ReorderPropertyGroups(ctx, req.GetScope(), req.GetGroupKeys()); err != nil {
		return nil, propertyDefinitionError(err)
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	return &pbsettings.ReorderPropertyGroupsResponse{}, nil
}

func (s *Server) ReorderPropertyDefinitions(
	ctx context.Context,
	req *pbsettings.ReorderPropertyDefinitionsRequest,
) (*pbsettings.ReorderPropertyDefinitionsResponse, error) {
	if err := s.store.ReorderPropertyDefinitions(
		ctx,
		req.GetScope(),
		req.GetGroupKey(),
		req.GetPropertyKeys(),
	); err != nil {
		return nil, propertyDefinitionError(err)
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	return &pbsettings.ReorderPropertyDefinitionsResponse{}, nil
}

func (s *Server) CreatePropertyDefinition(
	ctx context.Context,
	req *pbsettings.CreatePropertyDefinitionRequest,
) (*pbsettings.CreatePropertyDefinitionResponse, error) {
	definition, err := s.store.CreatePropertyDefinition(ctx, req.GetDefinition())
	if err != nil {
		return nil, propertyDefinitionError(err)
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)
	return &pbsettings.CreatePropertyDefinitionResponse{Definition: definition}, nil
}

func (s *Server) UpdatePropertyDefinition(
	ctx context.Context,
	req *pbsettings.UpdatePropertyDefinitionRequest,
) (*pbsettings.UpdatePropertyDefinitionResponse, error) {
	definition, err := s.store.UpdatePropertyDefinition(ctx, req.GetDefinition())
	if err != nil {
		return nil, propertyDefinitionError(err)
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	return &pbsettings.UpdatePropertyDefinitionResponse{Definition: definition}, nil
}

func (s *Server) DeletePropertyDefinition(
	ctx context.Context,
	req *pbsettings.DeletePropertyDefinitionRequest,
) (*pbsettings.DeletePropertyDefinitionResponse, error) {
	if err := s.store.DeletePropertyDefinition(ctx, req.GetScope(), req.GetKey()); err != nil {
		return nil, propertyDefinitionError(err)
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	return &pbsettings.DeletePropertyDefinitionResponse{}, nil
}

func (s *Server) RestorePropertyDefinition(
	ctx context.Context,
	req *pbsettings.RestorePropertyDefinitionRequest,
) (*pbsettings.RestorePropertyDefinitionResponse, error) {
	definition, err := s.store.RestorePropertyDefinition(ctx, req.GetScope(), req.GetKey())
	if err != nil {
		return nil, propertyDefinitionError(err)
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	return &pbsettings.RestorePropertyDefinitionResponse{Definition: definition}, nil
}

func (s *Server) SetSystemPropertyEnabled(
	ctx context.Context,
	req *pbsettings.SetSystemPropertyEnabledRequest,
) (*pbsettings.SetSystemPropertyEnabledResponse, error) {
	definition, err := s.store.SetSystemPropertyEnabled(
		ctx,
		req.GetScope(),
		req.GetKey(),
		req.GetEnabled(),
	)
	if err != nil {
		return nil, propertyDefinitionError(err)
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	return &pbsettings.SetSystemPropertyEnabledResponse{Definition: definition}, nil
}

func propertyDefinitionError(err error) error {
	if errors.Is(err, settingsstore.ErrInvalidPropertyDefinition) ||
		errors.Is(err, settingsstore.ErrPropertyDefinitionNotFound) ||
		errors.Is(err, settingsstore.ErrPropertyGroupNotFound) {
		return errorssettings.ErrInvalidRequest
	}
	return errswrap.NewError(err, errorssettings.ErrFailedQuery)
}
