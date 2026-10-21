package properties

import (
	"errors"

	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
)

var (
	ErrPropertyFilterNotSearchable = errors.New("property filter is not searchable")
	ErrPropertyAccessDenied        = errors.New("property access denied")
)

// EffectivePropertyAccess evaluates a caller's property grants against the
// effective definition catalog, including instance-level enablement overrides.
type EffectivePropertyAccess struct {
	scope       propertypb.PropertyScope
	allowed     *permissionsattributes.PropertyList
	definitions map[string]*propertypb.PropertyDefinition
	isAdmin     bool
}

func NewEffectivePropertyAccess(
	scope propertypb.PropertyScope,
	allowed *permissionsattributes.PropertyList,
	definitions map[string]*propertypb.PropertyDefinition,
	isAdmin bool,
) EffectivePropertyAccess {
	return EffectivePropertyAccess{
		scope:       scope,
		allowed:     allowed,
		definitions: definitions,
		isAdmin:     isAdmin,
	}
}

func (a EffectivePropertyAccess) IsEnabled(key string) bool {
	definition, ok := a.definitions[key]
	return ok && definition.GetEnabled()
}

// EnabledPropertyKeys returns keys that are present and enabled in the
// effective definition catalog. Callers remain responsible for limiting keys
// to those the caller is permitted to access.
func EnabledPropertyKeys(
	keys []string,
	definitions map[string]*propertypb.PropertyDefinition,
) []string {
	visibleKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		definition, ok := definitions[key]
		if ok && definition.GetEnabled() {
			visibleKeys = append(visibleKeys, key)
		}
	}
	return visibleKeys
}

// CanView requires an explicit property-list grant. Administrative access is
// intentionally not an implicit view grant for current property values.
func (a EffectivePropertyAccess) CanView(key string) bool {
	return a.IsEnabled(key) && a.allowed != nil && a.allowed.Contains(key)
}

// CanAccess permits property administrators to bypass role-list grants while
// still honoring the effective enabled state.
func (a EffectivePropertyAccess) CanAccess(key string) bool {
	return a.IsEnabled(key) && (a.isAdmin || (a.allowed != nil && a.allowed.Contains(key)))
}

// ValidateFilter verifies that a property is searchable, enabled, and
// available to the caller for filtering.
func (a EffectivePropertyAccess) ValidateFilter(key string) error {
	if !IsSearchablePropertyKey(a.scope, key) {
		return ErrPropertyFilterNotSearchable
	}
	definition, ok := a.definitions[key]
	if !ok || !definition.GetEnabled() || !a.CanAccess(key) {
		return ErrPropertyAccessDenied
	}
	if !definition.GetSearchable() {
		return ErrPropertyFilterNotSearchable
	}
	return nil
}
