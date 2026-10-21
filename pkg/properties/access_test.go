package properties

import (
	"testing"

	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEffectivePropertyAccess(t *testing.T) {
	t.Parallel()

	allowed := &permissionsattributes.PropertyList{Keys: []string{"vehicle.color"}}
	definitions := map[string]*propertypb.PropertyDefinition{
		"vehicle.color":    {Enabled: true, Searchable: true},
		"vehicle.disabled": {Enabled: false, Searchable: true},
	}
	access := NewEffectivePropertyAccess(
		propertypb.PropertyScope_PROPERTY_SCOPE_VEHICLE,
		allowed,
		definitions,
		false,
	)

	assert.True(t, access.CanView("vehicle.color"))
	assert.True(t, access.CanAccess("vehicle.color"))
	require.NoError(t, access.ValidateFilter("vehicle.color"))
	assert.False(t, access.CanView("vehicle.disabled"))
	require.ErrorIs(t, access.ValidateFilter("vehicle.disabled"), ErrPropertyAccessDenied)
	require.ErrorIs(t, access.ValidateFilter("vehicle.unknown"), ErrPropertyAccessDenied)

	adminAccess := NewEffectivePropertyAccess(
		propertypb.PropertyScope_PROPERTY_SCOPE_VEHICLE,
		&permissionsattributes.PropertyList{},
		definitions,
		true,
	)
	assert.False(t, adminAccess.CanView("vehicle.color"))
	assert.True(t, adminAccess.CanAccess("vehicle.color"))
}

func TestEnabledPropertyKeys(t *testing.T) {
	t.Parallel()

	definitions := map[string]*propertypb.PropertyDefinition{
		"vehicle.color":    {Enabled: true},
		"vehicle.disabled": {Enabled: false},
	}

	assert.Equal(t, []string{"vehicle.color"}, EnabledPropertyKeys([]string{
		"vehicle.color",
		"vehicle.disabled",
		"vehicle.unknown",
	}, definitions))
}
