// Package events distributes property-definition cache invalidations.
package events

import (
	"context"
	"encoding/json"
	"fmt"

	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	properties "github.com/fivenet-app/fivenet/v2026/pkg/properties"
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const DefinitionInvalidateSubject = "properties.definition.invalidate"

type definitionInvalidation struct {
	Scope propertypb.PropertyScope `json:"scope"`
}

// Publisher is deliberately small so stores do not depend on NATS.
type Publisher interface {
	Invalidate(context.Context, propertypb.PropertyScope) error
}

type Service struct {
	logger *zap.Logger
	nc     *nats.Conn
	cache  *properties.DefinitionCache
	sub    *nats.Subscription
}

type Params struct {
	fx.In

	LC     fx.Lifecycle
	Logger *zap.Logger
	NC     *nats.Conn
	Cache  *properties.DefinitionCache
}

func New(p Params) (*Service, error) {
	s := &Service{logger: p.Logger.Named("properties.events"), nc: p.NC, cache: p.Cache}
	p.LC.Append(fx.Hook{
		OnStart: func(context.Context) error {
			sub, err := s.nc.Subscribe(DefinitionInvalidateSubject, s.handle)
			if err != nil {
				return fmt.Errorf("subscribe to property definition invalidations: %w", err)
			}
			s.sub = sub
			if err := s.nc.Flush(); err != nil {
				_ = sub.Unsubscribe()
				s.sub = nil
				return fmt.Errorf("flush property definition invalidation subscription: %w", err)
			}
			return nil
		},
		OnStop: func(context.Context) error {
			if s.sub == nil {
				return nil
			}
			err := s.sub.Unsubscribe()
			s.sub = nil
			return err
		},
	})
	return s, nil
}

func (s *Service) Invalidate(_ context.Context, scope propertypb.PropertyScope) error {
	b, err := json.Marshal(definitionInvalidation{Scope: scope})
	if err != nil {
		return err
	}
	if err := s.nc.Publish(DefinitionInvalidateSubject, b); err != nil {
		return fmt.Errorf("publish property definition invalidation: %w", err)
	}
	return nil
}

func (s *Service) handle(msg *nats.Msg) {
	var event definitionInvalidation
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		s.logger.Warn("ignoring malformed property definition invalidation", zap.Error(err))
		return
	}
	if _, ok := propertyScopeName(event.Scope); !ok {
		s.logger.Warn("ignoring property definition invalidation with invalid scope")
		return
	}
	s.cache.Invalidate(event.Scope)
}

func propertyScopeName(scope propertypb.PropertyScope) (string, bool) {
	switch scope {
	case propertypb.PropertyScope_PROPERTY_SCOPE_USER:
		return "user", true
	case propertypb.PropertyScope_PROPERTY_SCOPE_VEHICLE:
		return "vehicle", true
	default:
		return "", false
	}
}

var Module = fx.Module(
	"property_definition_events",
	fx.Provide(New, func(s *Service) Publisher { return s }),
)
