package properties

import (
	"testing"

	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/stretchr/testify/require"
)

func TestDefinitionGroupKeyIncludesResolutionInputs(t *testing.T) {
	t.Parallel()

	base := &propertypb.PropertyDefinition{
		Scope:          propertypb.PropertyScope_PROPERTY_SCOPE_USER,
		Key:            "user.example",
		Enabled:        true,
		GroupKey:       "main",
		GroupSortOrder: 10,
	}
	definitions := map[string]*propertypb.PropertyDefinition{base.GetKey(): base}

	key := definitionGroupKey(base.GetScope(), definitions)
	base.Enabled = false
	require.NotEqual(t, key, definitionGroupKey(base.GetScope(), definitions))

	base.Enabled = true
	base.GroupKey = "secondary"
	require.NotEqual(t, key, definitionGroupKey(base.GetScope(), definitions))

	base.GroupKey = "main"
	base.GroupSortOrder = 20
	require.NotEqual(t, key, definitionGroupKey(base.GetScope(), definitions))
}

func TestDefinitionCacheInvalidatesOnlyScope(t *testing.T) {
	t.Parallel()

	cache := NewDefinitionCache()
	userKey := scopeAndKeys(propertypb.PropertyScope_PROPERTY_SCOPE_USER, nil)
	vehicleKey := scopeAndKeys(propertypb.PropertyScope_PROPERTY_SCOPE_VEHICLE, nil)

	cache.definitions.Put(userKey+"user.example", map[string]*propertypb.PropertyDefinition{}, 0)
	cache.definitions.Put(
		vehicleKey+"vehicle.example",
		map[string]*propertypb.PropertyDefinition{},
		0,
	)

	cache.Invalidate(propertypb.PropertyScope_PROPERTY_SCOPE_USER)

	_, userFound := cache.definitions.Get(userKey + "user.example")
	_, vehicleFound := cache.definitions.Get(vehicleKey + "vehicle.example")
	require.False(t, userFound)
	require.True(t, vehicleFound)
}
