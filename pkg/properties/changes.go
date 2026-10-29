package properties

import (
	"context"
	"errors"

	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
)

var (
	ErrInvalidPropertyChange = errors.New("invalid property change")
	ErrPropertyChangeDenied  = errors.New("property change denied")
)

type DefinitionResolver func(context.Context, []string) (map[string]*propertypb.PropertyDefinition, error)

// ValidateCustomPropertyChanges validates generic scalar property changes
// before the caller maps them to its resource-specific storage operation.
// An absent value clears a property and therefore does not require an active
// definition, allowing a stale value to be removed after its definition was
// disabled or deleted.
func ValidateCustomPropertyChanges(
	ctx context.Context,
	scope propertypb.PropertyScope,
	changes []*propertypb.PropertyChange,
	allowed func(string) bool,
	resolve DefinitionResolver,
) error {
	if len(changes) == 0 {
		return nil
	}
	if allowed == nil || resolve == nil {
		return ErrInvalidPropertyChange
	}

	keys := make([]string, 0, len(changes))
	seen := make(map[string]struct{}, len(changes))
	for _, change := range changes {
		if change == nil || !IsCustomPropertyKey(scope, change.GetKey()) {
			return ErrInvalidPropertyChange
		}
		if !allowed(change.GetKey()) {
			return ErrPropertyChangeDenied
		}
		if _, ok := seen[change.GetKey()]; ok {
			return ErrInvalidPropertyChange
		}
		seen[change.GetKey()] = struct{}{}
		keys = append(keys, change.GetKey())
	}

	definitions, err := resolve(ctx, keys)
	if err != nil {
		return err
	}
	for _, change := range changes {
		if change.GetValue() == nil {
			continue
		}
		valueType, _, _, _, _, _, _, err := ValueColumns(change.GetValue())
		definition := definitions[change.GetKey()]
		if err != nil || definition == nil || !definition.GetEnabled() ||
			definition.GetSource() != propertypb.PropertyDefinitionSource_PROPERTY_DEFINITION_SOURCE_CUSTOM ||
			definition.GetValueType() != valueType {
			return ErrInvalidPropertyChange
		}
	}

	return nil
}
