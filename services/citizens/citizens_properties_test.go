package citizens

import (
	"testing"

	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	users "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	"github.com/stretchr/testify/require"
)

func TestResolveCitizenPropertiesFiltersValuesAndIncludesEmptyCustomDefinitions(t *testing.T) {
	t.Parallel()

	visible := &properties.PropertyValue{}
	visible.SetStringValue("watchlist")
	hidden := &properties.PropertyValue{}
	hidden.SetStringValue("secret")

	definition := &properties.PropertyDefinition{
		Scope:     properties.PropertyScope_PROPERTY_SCOPE_USER,
		Source:    properties.PropertyDefinitionSource_PROPERTY_DEFINITION_SOURCE_CUSTOM,
		Key:       "user.risk_level",
		ValueType: properties.PropertyValueType_PROPERTY_VALUE_TYPE_STRING,
		Enabled:   true,
		Label:     "Risk level",
	}
	emptyDefinition := &properties.PropertyDefinition{
		Scope:     properties.PropertyScope_PROPERTY_SCOPE_USER,
		Source:    properties.PropertyDefinitionSource_PROPERTY_DEFINITION_SOURCE_CUSTOM,
		Key:       "user.clearance",
		ValueType: properties.PropertyValueType_PROPERTY_VALUE_TYPE_STRING,
		Enabled:   true,
		Label:     "Clearance",
	}
	user := &users.User{UserId: 42}
	resolveCitizenProperties(
		[]*users.User{user},
		map[int32][]*properties.Property{
			42: {
				{Key: "user.risk_level", Value: visible},
				{Key: "user.internal_note", Value: hidden},
			},
		},
		&permissionsattributes.PropertyList{Keys: []string{"user.risk_level", "user.clearance"}},
		map[string]*properties.PropertyDefinition{
			"user.risk_level": definition,
			"user.clearance":  emptyDefinition,
		},
	)

	require.Len(t, user.GetProperties().GetProperties(), 2)
	require.Equal(t, "user.clearance", user.GetProperties().GetProperties()[0].GetKey())
	require.Nil(t, user.GetProperties().GetProperties()[0].GetValue())
	require.Equal(t, "user.risk_level", user.GetProperties().GetProperties()[1].GetKey())
}
