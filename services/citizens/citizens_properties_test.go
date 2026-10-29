package citizens

import (
	"testing"

	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	userinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	users "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	usersprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/props"
	pbcitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens"
	errorscitizens "github.com/fivenet-app/fivenet/v2026/services/citizens/errors"
	"github.com/stretchr/testify/require"
)

func TestDerivedCitizenProperties(t *testing.T) {
	t.Parallel()

	sex := "f"
	height := float32(170.5)
	phoneNumber := "555-0100"
	visum := int32(2)
	user := &users.User{
		Dateofbirth: "2000-01-02",
		Sex:         &sex,
		Height:      &height,
		PhoneNumber: &phoneNumber,
		Visum:       &visum,
	}
	user.Props = &usersprops.UserProps{}
	user.Props.SetEmail("citizen@example.invalid")

	properties := derivedCitizenProperties(user, &permissionsattributes.PropertyList{Keys: []string{
		"user.dateofbirth",
		"user.sex",
		"user.height",
		"user.phone_number",
		"user.email",
		"user.visum",
	}})

	require.Len(t, properties, 6)
	require.Equal(t, "2000-01-02", properties[0].GetValue().GetStringValue())
	require.Equal(t, "f", properties[1].GetValue().GetStringValue())
	require.InEpsilon(t, 170.5, properties[2].GetValue().GetNumberValue(), 0.0001)
	require.Equal(t, "555-0100", properties[3].GetValue().GetStringValue())
	require.Equal(t, "citizen@example.invalid", properties[4].GetValue().GetStringValue())
	require.Equal(t, int64(2), properties[5].GetValue().GetIntegerValue())
}

func TestCitizenFilterAccessRequiresEnabledDefinitionAndGrant(t *testing.T) {
	t.Parallel()

	userInfo := &userinfo.UserInfo{}
	allowed := &permissionsattributes.PropertyList{Keys: []string{"user.height", "user.risk_level"}}
	definitions := map[string]*properties.PropertyDefinition{
		"user.height":     {Enabled: false},
		"user.risk_level": {Enabled: false},
	}

	req := &pbcitizens.ListCitizensRequest{}
	req.SetMinHeight(160)
	err := requireCitizenLegacyFilterKeys(userInfo, allowed, definitions, req)
	require.ErrorIs(t, err, errorscitizens.ErrPropsPropertiesDenied)

	value := &properties.PropertyValue{}
	value.SetStringValue("high-risk")
	err = requireCitizenPropertyKeys(userInfo, allowed, definitions, []*properties.PropertyFilter{{
		Key:       "user.risk_level",
		Operation: properties.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_EQUALS,
		Value:     value,
	}})
	require.ErrorIs(t, err, errorscitizens.ErrPropsPropertiesDenied)

	definitions["user.height"].Enabled = true
	definitions["user.risk_level"].Enabled = true
	definitions["user.risk_level"].Searchable = true
	require.NoError(t, requireCitizenLegacyFilterKeys(userInfo, allowed, definitions, req))
	require.NoError(
		t,
		requireCitizenPropertyKeys(userInfo, allowed, definitions, []*properties.PropertyFilter{{
			Key:       "user.risk_level",
			Operation: properties.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_EQUALS,
			Value:     value,
		}}),
	)
	require.ErrorIs(
		t,
		requireCitizenPropertyKeys(
			userInfo,
			&permissionsattributes.PropertyList{},
			definitions,
			[]*properties.PropertyFilter{{
				Key:       "user.risk_level",
				Operation: properties.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_EQUALS,
				Value:     value,
			}},
		),
		errorscitizens.ErrPropsPropertiesDenied,
	)
}

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
