package vehiclesstore

import (
	"context"
	"errors"

	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	_ "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/vehicles"
	propertyregistry "github.com/fivenet-app/fivenet/v2026/pkg/properties"
)

// IsGenericPropertyKey reports whether a key belongs to the generic property
// store. System-origin definitions remain owned by their dedicated behavior.
func IsGenericPropertyKey(key string) bool {
	return propertyregistry.IsCustomPropertyKey(
		properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
		key,
	)
}

func (s *Store) IsSystemPropertyEnabled(ctx context.Context, key string) (bool, error) {
	enabled, err := s.propertyCache.SystemPropertyEnabled(
		ctx, s.db, properties.PropertyScope_PROPERTY_SCOPE_VEHICLE, key,
	)
	if errors.Is(err, propertyregistry.ErrSystemPropertyNotFound) {
		return false, ErrInvalidProperty
	}
	return enabled, err
}
