package properties

import (
	"slices"
	"strings"
	"sync"

	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"google.golang.org/protobuf/proto"
)

var (
	muSystemProperties   sync.RWMutex
	systemProperties     []*propertypb.PropertyDefinition
	systemPropertyGroups []*propertypb.PropertyGroup
)

// AddSystemPropertyGroups is initialization-only. Generated package init
// functions register their static metadata before the application starts; do
// not call this from request handling or any runtime mutation path.
func AddSystemPropertyGroups(groups []*propertypb.PropertyGroup) {
	muSystemProperties.Lock()
	defer muSystemProperties.Unlock()
	for _, group := range groups {
		if group == nil {
			continue
		}
		pgCopy := proto.Clone(group).(*propertypb.PropertyGroup)
		index := slices.IndexFunc(
			systemPropertyGroups,
			func(existing *propertypb.PropertyGroup) bool {
				return existing.GetScope() == group.GetScope() &&
					existing.GetKey() == group.GetKey()
			},
		)
		if index >= 0 {
			systemPropertyGroups[index] = pgCopy
		} else {
			systemPropertyGroups = append(systemPropertyGroups, pgCopy)
		}
	}
}

func SystemPropertyGroups() []*propertypb.PropertyGroup {
	muSystemProperties.RLock()
	defer muSystemProperties.RUnlock()
	groups := make([]*propertypb.PropertyGroup, 0, len(systemPropertyGroups))
	for _, group := range systemPropertyGroups {
		groups = append(groups, proto.Clone(group).(*propertypb.PropertyGroup))
	}
	slices.SortFunc(groups, func(a, b *propertypb.PropertyGroup) int {
		if a.GetScope() != b.GetScope() {
			return int(a.GetScope() - b.GetScope())
		}
		if a.GetSortOrder() != b.GetSortOrder() {
			return int(a.GetSortOrder() - b.GetSortOrder())
		}
		return strings.Compare(a.GetKey(), b.GetKey())
	})
	return groups
}

func SystemPropertyGroup(
	scope propertypb.PropertyScope,
	key string,
) (*propertypb.PropertyGroup, bool) {
	muSystemProperties.RLock()
	defer muSystemProperties.RUnlock()
	for _, group := range systemPropertyGroups {
		if group.GetScope() == scope && group.GetKey() == key {
			return proto.Clone(group).(*propertypb.PropertyGroup), true
		}
	}
	return nil, false
}

// AddSystemPropertyDefinitions is initialization-only. Generated package init
// functions register static metadata before the application starts; do not
// call this from request handling or any runtime mutation path. Registration
// is idempotent by scope and key so generated packages can contribute safely.
func AddSystemPropertyDefinitions(definitions []*propertypb.PropertyDefinition) {
	muSystemProperties.Lock()
	defer muSystemProperties.Unlock()

	for _, definition := range definitions {
		if definition == nil {
			continue
		}
		index := slices.IndexFunc(
			systemProperties,
			func(existing *propertypb.PropertyDefinition) bool {
				return existing.GetScope() == definition.GetScope() &&
					existing.GetKey() == definition.GetKey()
			},
		)
		spCopy := proto.Clone(definition).(*propertypb.PropertyDefinition)
		if index >= 0 {
			systemProperties[index] = spCopy
			continue
		}
		systemProperties = append(systemProperties, spCopy)
	}
}

// SystemPropertyDefinitions returns a defensive copy of all generated system
// property definitions.
func SystemPropertyDefinitions() []*propertypb.PropertyDefinition {
	muSystemProperties.RLock()
	defer muSystemProperties.RUnlock()

	definitions := make([]*propertypb.PropertyDefinition, 0, len(systemProperties))
	for _, definition := range systemProperties {
		definitions = append(definitions, proto.Clone(definition).(*propertypb.PropertyDefinition))
	}
	return definitions
}

// SystemPropertyDefinition returns a defensive copy of one generated system
// definition. It is used by query compilers that must distinguish system
// definitions from database-backed custom definitions.
func SystemPropertyDefinition(
	scope propertypb.PropertyScope,
	key string,
) (*propertypb.PropertyDefinition, bool) {
	muSystemProperties.RLock()
	defer muSystemProperties.RUnlock()
	for _, definition := range systemProperties {
		if definition.GetScope() == scope && definition.GetKey() == key {
			return proto.Clone(definition).(*propertypb.PropertyDefinition), true
		}
	}
	return nil, false
}
