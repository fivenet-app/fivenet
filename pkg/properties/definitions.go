package properties

import (
	"errors"
	"regexp"
	"slices"
	"sort"
	"strings"

	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"google.golang.org/protobuf/proto"
)

var (
	ErrInvalidPropertyKey        = errors.New("invalid property key")
	ErrInvalidPropertyDefinition = errors.New("invalid property definition")
)

var propertyKeyPattern = regexp.MustCompile(`^[a-z][a-z0-9_.-]*$`)

func ScopeName(scope propertypb.PropertyScope) (string, bool) {
	switch scope {
	case propertypb.PropertyScope_PROPERTY_SCOPE_USER:
		return "user", true
	case propertypb.PropertyScope_PROPERTY_SCOPE_VEHICLE:
		return "vehicle", true
	default:
		return "", false
	}
}

func ScopePrefix(scope propertypb.PropertyScope) (string, bool) {
	name, ok := ScopeName(scope)
	if !ok {
		return "", false
	}
	return name + ".", true
}

func ValidateKey(key string) error {
	if len(key) == 0 || len(key) > 64 || !propertyKeyPattern.MatchString(key) {
		return ErrInvalidPropertyKey
	}
	return nil
}

func ValidateScopedKey(scope propertypb.PropertyScope, key string) error {
	prefix, ok := ScopePrefix(scope)
	if !ok || !strings.HasPrefix(key, prefix) || len(key) == len(prefix) {
		return ErrInvalidPropertyKey
	}
	return ValidateKey(key)
}

func IsSystemProperty(scope propertypb.PropertyScope, key string) bool {
	definition, ok := SystemPropertyDefinition(scope, key)
	return ok &&
		definition.GetSource() == propertypb.PropertyDefinitionSource_PROPERTY_DEFINITION_SOURCE_SYSTEM
}

func IsCustomPropertyKey(scope propertypb.PropertyScope, key string) bool {
	return ValidateScopedKey(scope, key) == nil && !IsSystemProperty(scope, key)
}

// IsSearchablePropertyKey reports whether a key is structurally valid and may
// be used by the generic property-filter API. Custom keys are accepted here;
// their persisted definition and searchable flag are validated when the
// filter is compiled.
func IsSearchablePropertyKey(scope propertypb.PropertyScope, key string) bool {
	if ValidateScopedKey(scope, key) != nil {
		return false
	}
	if definition, ok := SystemPropertyDefinition(scope, key); ok {
		return definition.GetSearchable()
	}
	return true
}

func SystemPropertyDefinitionsForScope(
	scope propertypb.PropertyScope,
) []*propertypb.PropertyDefinition {
	definitions := SystemPropertyDefinitions()
	filtered := make([]*propertypb.PropertyDefinition, 0, len(definitions))
	for _, definition := range definitions {
		if definition.GetScope() == scope {
			filtered = append(filtered, definition)
		}
	}
	return filtered
}

// NormalizeDefinition returns a normalized copy and never mutates the
// caller-owned protobuf.
func NormalizeDefinition(definition *propertypb.PropertyDefinition) *propertypb.PropertyDefinition {
	if definition == nil {
		return nil
	}
	normalized := proto.Clone(definition).(*propertypb.PropertyDefinition)
	if normalized.GetValueFormat() == propertypb.PropertyValueFormat_PROPERTY_VALUE_FORMAT_UNSPECIFIED {
		normalized.ValueFormat = propertypb.PropertyValueFormat_PROPERTY_VALUE_FORMAT_PLAIN
	}
	if normalized.GetBehavior() == propertypb.PropertyBehavior_PROPERTY_BEHAVIOR_UNSPECIFIED {
		normalized.Behavior = propertypb.PropertyBehavior_PROPERTY_BEHAVIOR_GENERIC
	}
	if normalized.GetRenderer() == propertypb.PropertyRenderer_PROPERTY_RENDERER_UNSPECIFIED {
		normalized.Renderer = propertypb.PropertyRenderer_PROPERTY_RENDERER_DEFAULT
	}
	if normalized.GetDisplayMode() == propertypb.PropertyDisplayMode_PROPERTY_DISPLAY_MODE_UNSPECIFIED {
		normalized.DisplayMode = propertypb.PropertyDisplayMode_PROPERTY_DISPLAY_MODE_DEFAULT
	}
	return normalized
}

func IsValueFormatCompatible(
	valueType propertypb.PropertyValueType,
	valueFormat propertypb.PropertyValueFormat,
) bool {
	switch valueFormat {
	case propertypb.PropertyValueFormat_PROPERTY_VALUE_FORMAT_PLAIN:
		return true
	case propertypb.PropertyValueFormat_PROPERTY_VALUE_FORMAT_MONEY:
		return valueType == propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_INTEGER
	case propertypb.PropertyValueFormat_PROPERTY_VALUE_FORMAT_PERCENTAGE,
		propertypb.PropertyValueFormat_PROPERTY_VALUE_FORMAT_DISTANCE:
		return valueType == propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_INTEGER ||
			valueType == propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_NUMBER
	case propertypb.PropertyValueFormat_PROPERTY_VALUE_FORMAT_DATETIME:
		return valueType == propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_TIMESTAMP
	case propertypb.PropertyValueFormat_PROPERTY_VALUE_FORMAT_DATE:
		return valueType == propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_DATE
	case propertypb.PropertyValueFormat_PROPERTY_VALUE_FORMAT_DURATION:
		return valueType == propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_INTEGER
	default:
		return false
	}
}

func ValidateCustomDefinition(definition *propertypb.PropertyDefinition) error {
	_, err := NormalizeAndValidateCustomDefinition(definition)
	return err
}

// NormalizeAndValidateCustomDefinition returns a normalized copy suitable for
// persistence, keeping validation free of caller-visible side effects.
func NormalizeAndValidateCustomDefinition(
	definition *propertypb.PropertyDefinition,
) (*propertypb.PropertyDefinition, error) {
	normalized := NormalizeDefinition(definition)
	if normalized == nil ||
		normalized.GetSource() != propertypb.PropertyDefinitionSource_PROPERTY_DEFINITION_SOURCE_CUSTOM ||
		normalized.GetValueType() == propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_UNSPECIFIED ||
		normalized.GetBehavior() != propertypb.PropertyBehavior_PROPERTY_BEHAVIOR_GENERIC ||
		!IsValueFormatCompatible(normalized.GetValueType(), normalized.GetValueFormat()) ||
		!IsRendererCompatible(normalized.GetValueType(), normalized.GetRenderer()) ||
		!IsRendererConfigCompatible(normalized) ||
		strings.TrimSpace(normalized.GetLabel()) == "" ||
		!IsCustomPropertyKey(normalized.GetScope(), normalized.GetKey()) {
		return nil, ErrInvalidPropertyDefinition
	}
	return normalized, nil
}

// IsRendererCompatible reports whether a renderer can consume the definition's
// scalar value type. Generic renderers are intended for custom definitions;
// special renderers remain reserved for system-managed properties.
func IsRendererCompatible(
	valueType propertypb.PropertyValueType,
	renderer propertypb.PropertyRenderer,
) bool {
	switch renderer {
	case propertypb.PropertyRenderer_PROPERTY_RENDERER_DEFAULT:
		return true
	case propertypb.PropertyRenderer_PROPERTY_RENDERER_STRING:
		return valueType == propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_STRING
	case propertypb.PropertyRenderer_PROPERTY_RENDERER_INTEGER:
		return valueType == propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_INTEGER
	case propertypb.PropertyRenderer_PROPERTY_RENDERER_NUMBER:
		return valueType == propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_NUMBER
	case propertypb.PropertyRenderer_PROPERTY_RENDERER_BOOLEAN:
		return valueType == propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_BOOL
	case propertypb.PropertyRenderer_PROPERTY_RENDERER_DATE:
		return valueType == propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_DATE
	case propertypb.PropertyRenderer_PROPERTY_RENDERER_DATETIME:
		return valueType == propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_TIMESTAMP
	default:
		return false
	}
}

// IsRendererConfigCompatible validates the oneof selected by a renderer and
// keeps renderer-specific options from silently being ignored.
func IsRendererConfigCompatible(definition *propertypb.PropertyDefinition) bool {
	if definition == nil {
		return false
	}
	config := definition.GetRendererConfig()
	if config == nil {
		return true
	}
	switch definition.GetRenderer() {
	case propertypb.PropertyRenderer_PROPERTY_RENDERER_STRING:
		return config.GetString() != nil
	case propertypb.PropertyRenderer_PROPERTY_RENDERER_INTEGER:
		return config.GetInteger() != nil && config.GetInteger().GetDecimalPlaces() <= 9
	case propertypb.PropertyRenderer_PROPERTY_RENDERER_NUMBER:
		return config.GetNumber() != nil && config.GetNumber().GetDecimalPlaces() <= 9
	default:
		return false
	}
}

func SortDefinitions(
	definitions []*propertypb.PropertyDefinition,
) []*propertypb.PropertyDefinition {
	result := append([]*propertypb.PropertyDefinition(nil), definitions...)
	sort.SliceStable(result, func(i, j int) bool {
		left, right := result[i], result[j]
		leftOrder, rightOrder := left.GetGroupSortOrder(), right.GetGroupSortOrder()
		if leftOrder != rightOrder {
			return leftOrder < rightOrder
		}
		if left.GetGroupKey() != right.GetGroupKey() {
			return left.GetGroupKey() < right.GetGroupKey()
		}
		if left.GetSortOrder() != right.GetSortOrder() {
			return left.GetSortOrder() < right.GetSortOrder()
		}
		if left.GetSource() != right.GetSource() {
			return left.GetSource() < right.GetSource()
		}
		return left.GetKey() < right.GetKey()
	})
	return result
}

// SortProperties orders resolved property values by their definition order.
// Values without a definition are kept after defined values and sorted by key.
func SortProperties(values []*propertypb.Property) []*propertypb.Property {
	result := append([]*propertypb.Property(nil), values...)
	sort.SliceStable(result, func(i, j int) bool {
		left, right := result[i], result[j]
		leftGroupOrder, rightGroupOrder := int32(1<<30), int32(1<<30)
		leftOrder, rightOrder := int32(1<<30), int32(1<<30)
		if left.GetDefinition() != nil {
			leftGroupOrder = left.GetDefinition().GetGroupSortOrder()
			leftOrder = left.GetDefinition().GetSortOrder()
		}
		if right.GetDefinition() != nil {
			rightGroupOrder = right.GetDefinition().GetGroupSortOrder()
			rightOrder = right.GetDefinition().GetSortOrder()
		}
		if leftGroupOrder != rightGroupOrder {
			return leftGroupOrder < rightGroupOrder
		}
		if leftOrder != rightOrder {
			return leftOrder < rightOrder
		}
		return left.GetKey() < right.GetKey()
	})
	return result
}

// ResolveProperties attaches the effective definitions and returns the values
// in their configured display order.
func ResolveProperties(
	values []*propertypb.Property,
	definitions map[string]*propertypb.PropertyDefinition,
) []*propertypb.Property {
	result := make([]*propertypb.Property, 0, len(values)+len(definitions))
	resolved := make(map[string]struct{}, len(values))
	for _, value := range values {
		if value == nil {
			continue
		}
		definition, ok := definitions[value.GetKey()]
		if !ok || definition == nil || !definition.GetEnabled() {
			continue
		}
		value.Definition = definition
		result = append(result, value)
		resolved[value.GetKey()] = struct{}{}
	}
	// Include empty custom definitions so callers can create their first value
	// without separately loading the definition catalog. System properties
	// retain their dedicated resource-specific editing flows.
	for key, definition := range definitions {
		if definition == nil || !definition.GetEnabled() ||
			definition.GetSource() != propertypb.PropertyDefinitionSource_PROPERTY_DEFINITION_SOURCE_CUSTOM {
			continue
		}
		if _, ok := resolved[key]; ok {
			continue
		}
		result = append(result, &propertypb.Property{Key: key, Definition: definition})
	}
	return SortProperties(result)
}

// ResolvePropertiesResource builds the wire representation used by services
// after values and definitions have been filtered for the caller.
func ResolvePropertiesResource(
	values []*propertypb.Property,
	definitions map[string]*propertypb.PropertyDefinition,
) *propertypb.Properties {
	return &propertypb.Properties{
		Properties: ResolveProperties(values, definitions),
		Groups:     ResolvePropertyGroups(definitions),
	}
}

func ResolvePropertiesResourceWithGroups(
	values []*propertypb.Property,
	definitions map[string]*propertypb.PropertyDefinition,
	groups []*propertypb.PropertyGroup,
) *propertypb.Properties {
	return &propertypb.Properties{
		Properties: ResolveProperties(values, definitions),
		Groups:     groups,
	}
}

// ResolvePropertyGroups returns the effective groups represented by the
// supplied definitions. Definitions are already resolved for the caller's
// scope and permissions, including instance-specific group ordering.
func ResolvePropertyGroups(
	definitions map[string]*propertypb.PropertyDefinition,
) []*propertypb.PropertyGroup {
	groups := make(map[string]*propertypb.PropertyGroup)
	for _, definition := range definitions {
		if definition == nil || !definition.GetEnabled() || definition.GetGroupKey() == "" {
			continue
		}
		key := definition.GetGroupKey()
		if existing, ok := groups[key]; ok {
			if definition.GetGroupSortOrder() < existing.GetSortOrder() {
				existing.SortOrder = definition.GetGroupSortOrder()
			}
			continue
		}
		group := &propertypb.PropertyGroup{
			Scope:     definition.GetScope(),
			Key:       key,
			SortOrder: definition.GetGroupSortOrder(),
		}
		if systemGroup, ok := SystemPropertyGroup(definition.GetScope(), key); ok {
			group.DisplayKey = systemGroup.GetDisplayKey()
		}
		groups[key] = group
	}

	result := make([]*propertypb.PropertyGroup, 0, len(groups))
	for _, group := range groups {
		result = append(result, group)
	}
	slices.SortFunc(result, func(a, b *propertypb.PropertyGroup) int {
		if a.GetSortOrder() != b.GetSortOrder() {
			return int(a.GetSortOrder() - b.GetSortOrder())
		}
		return strings.Compare(a.GetKey(), b.GetKey())
	})
	return result
}

func MergeDefinitions(
	system []*propertypb.PropertyDefinition,
	custom []*propertypb.PropertyDefinition,
	enabledOverrides map[string]bool,
) []*propertypb.PropertyDefinition {
	definitions := make([]*propertypb.PropertyDefinition, 0, len(system)+len(custom))
	for _, definition := range append(append([]*propertypb.PropertyDefinition{}, system...), custom...) {
		if definition == nil {
			continue
		}
		copyDefinition := proto.Clone(definition).(*propertypb.PropertyDefinition)
		if copyDefinition.GetSource() == propertypb.PropertyDefinitionSource_PROPERTY_DEFINITION_SOURCE_SYSTEM {
			if enabled, ok := enabledOverrides[copyDefinition.GetKey()]; ok {
				copyDefinition.Enabled = enabled
			}
		}
		if group, ok := SystemPropertyGroup(
			copyDefinition.GetScope(),
			copyDefinition.GetGroupKey(),
		); ok {
			copyDefinition.GroupSortOrder = group.GetSortOrder()
		} else {
			copyDefinition.GroupSortOrder = 1 << 30
		}
		definitions = append(definitions, copyDefinition)
	}
	return SortDefinitions(definitions)
}

// ApplySortOrderOverrides applies instance-specific ordering to system
// definitions before the normal group/property sort is performed.
func ApplySortOrderOverrides(
	definitions []*propertypb.PropertyDefinition,
	overrides map[string]int32,
) []*propertypb.PropertyDefinition {
	result := make([]*propertypb.PropertyDefinition, 0, len(definitions))
	for _, definition := range definitions {
		copyDefinition := proto.Clone(definition).(*propertypb.PropertyDefinition)
		if order, ok := overrides[copyDefinition.GetKey()]; ok {
			copyDefinition.SortOrder = order
		}
		result = append(result, copyDefinition)
	}
	return SortDefinitions(result)
}
