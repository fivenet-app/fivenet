package perms

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDynamicPropertyListValidValuesIncludesCustomKeys(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT fivenet_property_definitions.property_key AS "fivenet_property_definitions.property_key",
       fivenet_property_definitions.label AS "fivenet_property_definitions.label",
       fivenet_property_definitions.implicit AS "fivenet_property_definitions.implicit"
FROM fivenet_property_definitions
WHERE (fivenet_property_definitions.scope = ?) AND (fivenet_property_definitions.deleted_at IS NULL);`)).
		WithArgs("vehicle").
		WillReturnRows(sqlmock.NewRows([]string{
			"fivenet_property_definitions.property_key",
			"fivenet_property_definitions.label",
			"fivenet_property_definitions.implicit",
		}).
			AddRow("vehicle.color", "Color", false).
			AddRow("user.phone_number", "Phone number", false).
			AddRow("vehicle.color", "Color", false).
			AddRow("vehicle.secret", "Secret", true))

	base := &permissionsattributes.AttributeValues{
		ValidValues: &permissionsattributes.AttributeValues_PropertyList{
			PropertyList: &permissionsattributes.PropertyList{Keys: []string{"vehicle.wanted"}},
		},
	}

	got, err := (&Perms{db: db}).dynamicPropertyListValidValues(
		t.Context(),
		properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
		base,
	)
	require.NoError(t, err)
	assert.Equal(t, []string{"vehicle.color", "vehicle.wanted"}, got.GetPropertyList().GetKeys())
	assert.Equal(t, []*permissionsattributes.PropertyListEntry{
		{Key: "vehicle.color", Label: "Color"},
		{Key: "vehicle.wanted", Label: "vehicle.wanted"},
	}, got.GetPropertyList().GetProperties())
	assert.Equal(
		t,
		[]string{"vehicle.wanted"},
		base.GetPropertyList().GetKeys(),
		"base options should not be mutated",
	)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestIncludeExistingPropertyListGrantsPreservesDeletedGrantOnEdit(t *testing.T) {
	t.Parallel()

	validValues := &permissionsattributes.AttributeValues{
		ValidValues: &permissionsattributes.AttributeValues_PropertyList{
			PropertyList: &permissionsattributes.PropertyList{Keys: []string{"vehicle.active"}},
		},
	}
	existingValue := &permissionsattributes.AttributeValues{
		ValidValues: &permissionsattributes.AttributeValues_PropertyList{
			PropertyList: &permissionsattributes.PropertyList{Keys: []string{"vehicle.deleted"}},
		},
	}

	got := includeExistingPropertyListGrants(
		validValues,
		existingValue,
		properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
	)
	assert.Equal(t, []string{"vehicle.active", "vehicle.deleted"}, got.GetPropertyList().GetKeys())
	assert.Equal(t, []string{"vehicle.active"}, validValues.GetPropertyList().GetKeys())
}

func TestFilterPropertyListValueHidesDeletedKeys(t *testing.T) {
	t.Parallel()

	value := &permissionsattributes.AttributeValues{
		ValidValues: &permissionsattributes.AttributeValues_PropertyList{
			PropertyList: &permissionsattributes.PropertyList{
				Keys: []string{"vehicle.active", "vehicle.deleted"},
			},
		},
	}
	validValues := &permissionsattributes.AttributeValues{
		ValidValues: &permissionsattributes.AttributeValues_PropertyList{
			PropertyList: &permissionsattributes.PropertyList{Keys: []string{"vehicle.active"}},
		},
	}

	got := filterPropertyListValue(value, validValues)
	assert.Equal(t, []string{"vehicle.active"}, got.GetPropertyList().GetKeys())
	assert.Equal(
		t,
		[]string{"vehicle.active", "vehicle.deleted"},
		value.GetPropertyList().GetKeys(),
	)
}
