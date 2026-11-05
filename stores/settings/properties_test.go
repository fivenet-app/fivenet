package settingsstore

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	permsvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/vehicles/perms"
	propertyregistry "github.com/fivenet-app/fivenet/v2026/pkg/properties"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/protobuf/proto"
)

type jsonArgument string

func (expected jsonArgument) Match(value driver.Value) bool {
	actual, ok := value.(string)
	if !ok {
		return false
	}
	var expectedValue, actualValue any
	return json.Unmarshal([]byte(expected), &expectedValue) == nil &&
		json.Unmarshal([]byte(actual), &actualValue) == nil &&
		assert.ObjectsAreEqual(expectedValue, actualValue)
}

func TestListPropertyDefinitionsMergesSystemOverridesAndCustomDefinitions(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	store := New(db).(*Store)

	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_property_system_overrides")).
		WithArgs(vehiclePropertyScopeName).
		WillReturnRows(sqlmock.NewRows([]string{
			"fivenet_property_system_overrides.property_key",
			"fivenet_property_system_overrides.enabled",
		}).AddRow(string(permsvehicles.VehiclesServiceListVehiclesPropertiesPermValueVehicleWanted), false))
	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_property_definitions")).
		WithArgs(vehiclePropertyScopeName, false).
		WillReturnRows(sqlmock.NewRows([]string{
			"fivenet_property_definitions.property_key",
			"fivenet_property_definitions.value_type",
			"fivenet_property_definitions.label",
			"fivenet_property_definitions.description",
			"fivenet_property_definitions.enabled",
			"fivenet_property_definitions.searchable",
			"fivenet_property_definitions.value_format",
			"fivenet_property_definitions.behavior",
			"fivenet_property_definitions.unit",
			"fivenet_property_definitions.decimal_places",
			"fivenet_property_definitions.currency",
			"fivenet_property_definitions.group_key",
			"fivenet_property_definitions.sort_order",
			"fivenet_property_definitions.pinned",
			"fivenet_property_definitions.renderer",
			"fivenet_property_definitions.renderer_config",
			"fivenet_property_definitions.display_mode",
			"fivenet_property_definitions.implicit",
			"fivenet_property_definitions.state_policy",
			"fivenet_property_definitions.value_constraints",
			"fivenet_property_definitions.deleted_at",
		}).AddRow(
			"vehicle.inspection_due",
			properties.PropertyValueType_PROPERTY_VALUE_TYPE_TIMESTAMP,
			"Inspection due",
			"The next inspection date.",
			true,
			true,
			0,
			1,
			nil,
			nil,
			nil,
			"",
			0,
			false,
			properties.PropertyRenderer_PROPERTY_RENDERER_DEFAULT,
			nil,
			properties.PropertyDisplayMode_PROPERTY_DISPLAY_MODE_DEFAULT,
			false,
			nil,
			nil,
			nil,
		))
	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_property_group_overrides")).
		WithArgs(vehiclePropertyScopeName).
		WillReturnRows(sqlmock.NewRows([]string{
			"fivenet_property_group_overrides.group_key",
			"fivenet_property_group_overrides.sort_order",
		}))
	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_property_groups")).
		WithArgs(vehiclePropertyScopeName).
		WillReturnRows(sqlmock.NewRows([]string{
			"fivenet_property_groups.group_key",
			"fivenet_property_groups.sort_order",
		}))

	definitions, err := store.ListPropertyDefinitions(
		t.Context(),
		properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
		false,
	)
	require.NoError(t, err)
	require.Len(t, definitions, 3)
	assert.Equal(
		t,
		string(permsvehicles.VehiclesServiceListVehiclesPropertiesPermValueVehicleWanted),
		definitions[0].GetKey(),
	)
	assert.False(t, definitions[0].GetEnabled())
	assert.Equal(t, "vehicle.files", definitions[1].GetKey())
	assert.Equal(t, "vehicle.inspection_due", definitions[2].GetKey())
	assert.Equal(
		t,
		properties.PropertyDefinitionSource_PROPERTY_DEFINITION_SOURCE_CUSTOM,
		definitions[2].GetSource(),
	)
	assert.True(t, definitions[2].GetSearchable())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePropertyDefinitionPersistsCustomDefinition(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	store := New(db).(*Store)
	definition := &properties.PropertyDefinition{
		Scope:      properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
		Key:        "vehicle.inspection_due",
		ValueType:  properties.PropertyValueType_PROPERTY_VALUE_TYPE_TIMESTAMP,
		Source:     properties.PropertyDefinitionSource_PROPERTY_DEFINITION_SOURCE_CUSTOM,
		Enabled:    true,
		Searchable: true,
		Label:      "Inspection due",
		StatePolicy: &properties.PropertyStatePolicy{
			Enabled: true, SupportsExpiry: true,
			ExpiryAction: properties.PropertyExpiryAction_PROPERTY_EXPIRY_ACTION_CLEAR_VALUE,
		},
	}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO fivenet_property_definitions")).
		WithArgs(
			vehiclePropertyScopeName,
			definition.GetKey(),
			definition.GetValueType(),
			definition.GetLabel(),
			nil,
			true,
			true,
			properties.PropertyValueFormat_PROPERTY_VALUE_FORMAT_PLAIN,
			properties.PropertyBehavior_PROPERTY_BEHAVIOR_GENERIC,
			"",
			int32(0),
			false,
			properties.PropertyRenderer_PROPERTY_RENDERER_DEFAULT,
			nil,
			properties.PropertyDisplayMode_PROPERTY_DISPLAY_MODE_DEFAULT,
			false,
			jsonArgument(
				`{"enabled":true,"supportsExpiry":true,"expiryAction":"PROPERTY_EXPIRY_ACTION_CLEAR_VALUE"}`,
			),
			nil,
			properties.PropertyReasonMode_PROPERTY_REASON_MODE_REQUIRED,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	persisted, err := store.CreatePropertyDefinition(t.Context(), definition)
	require.NoError(t, err)
	assert.Equal(t, propertyregistry.NormalizeDefinition(definition), persisted)
	assert.Equal(
		t,
		properties.PropertyValueFormat_PROPERTY_VALUE_FORMAT_UNSPECIFIED,
		definition.GetValueFormat(),
	)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestScanPropertyDefinitionValuesUsesUserKeysetAndBoundsSamples(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	store := New(db).(*Store)
	observerCore, observed := observer.New(zap.WarnLevel)
	store.logger = zap.New(observerCore)

	definition := &properties.PropertyDefinition{
		Scope: properties.PropertyScope_PROPERTY_SCOPE_USER, Key: "user.nickname",
		ValueType:        properties.PropertyValueType_PROPERTY_VALUE_TYPE_STRING,
		ValueConstraints: &properties.PropertyValueConstraints{StringMaxLength: proto.Uint32(3)},
	}
	columns := []string{
		"fivenet_user_property_values.user_id",
		"fivenet_user_property_values.property_key",
		"fivenet_user_property_values.value_type",
		"fivenet_user_property_values.value_string",
		"fivenet_user_property_values.value_integer",
		"fivenet_user_property_values.value_number",
		"fivenet_user_property_values.value_bool",
		"fivenet_user_property_values.value_timestamp",
		"fivenet_user_property_values.value_date",
		"fivenet_user_property_values.updated_at",
	}
	first := sqlmock.NewRows(columns)
	for id := int32(1); id <= propertyConstraintScanBatchSize; id++ {
		first.AddRow(
			id,
			definition.GetKey(),
			uint8(definition.GetValueType()),
			"scanner-secret-value",
			nil,
			nil,
			nil,
			nil,
			nil,
			time.Now(),
		)
	}
	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_user_property_values")).
		WithArgs(definition.GetKey(), int32(0), propertyConstraintScanBatchSize).
		WillReturnRows(first)
	second := sqlmock.NewRows(columns).
		AddRow(int32(propertyConstraintScanBatchSize+1), definition.GetKey(), uint8(definition.GetValueType()), "scanner-secret-value", nil, nil, nil, nil, nil, time.Now())
	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_user_property_values")).
		WithArgs(definition.GetKey(), int32(propertyConstraintScanBatchSize), propertyConstraintScanBatchSize).
		WillReturnRows(second)

	require.NoError(t, store.scanPropertyDefinitionValues(t.Context(), definition))
	require.NoError(t, mock.ExpectationsWereMet())
	entries := observed.All()
	require.Len(t, entries, 1)
	fields := entries[0].ContextMap()
	assert.EqualValues(t, 501, fields["property.incompatible_value_count"])
	samples, ok := fields["property.resource_id_sample"].([]any)
	require.True(t, ok)
	assert.Len(t, samples, 10)
	assert.NotContains(t, fields, "property.value")
	assert.NotContains(t, fmt.Sprint(fields), "scanner-secret-value")
}

func TestScanPropertyDefinitionValuesCountsMalformedVehicleValues(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	store := New(db).(*Store)
	observerCore, observed := observer.New(zap.WarnLevel)
	store.logger = zap.New(observerCore)

	definition := &properties.PropertyDefinition{
		Scope: properties.PropertyScope_PROPERTY_SCOPE_VEHICLE, Key: "vehicle.note",
		ValueType:        properties.PropertyValueType_PROPERTY_VALUE_TYPE_STRING,
		ValueConstraints: &properties.PropertyValueConstraints{StringMinLength: proto.Uint32(3)},
	}
	columns := []string{
		"fivenet_vehicle_property_values.plate",
		"fivenet_vehicle_property_values.property_key",
		"fivenet_vehicle_property_values.value_type",
		"fivenet_vehicle_property_values.value_string",
		"fivenet_vehicle_property_values.value_integer",
		"fivenet_vehicle_property_values.value_number",
		"fivenet_vehicle_property_values.value_bool",
		"fivenet_vehicle_property_values.value_timestamp",
		"fivenet_vehicle_property_values.value_date",
		"fivenet_vehicle_property_values.updated_at",
	}
	rows := sqlmock.NewRows(columns).
		AddRow("ABC123", definition.GetKey(), uint8(definition.GetValueType()), "valid", nil, nil, nil, nil, nil, time.Now()).
		AddRow("BAD999", definition.GetKey(), uint8(definition.GetValueType()), nil, nil, nil, nil, nil, nil, time.Now())
	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_vehicle_property_values")).
		WithArgs(definition.GetKey(), "", propertyConstraintScanBatchSize).
		WillReturnRows(rows)

	require.NoError(t, store.scanPropertyDefinitionValues(t.Context(), definition))
	require.NoError(t, mock.ExpectationsWereMet())
	entries := observed.All()
	require.Len(t, entries, 1)
	assert.EqualValues(t, 1, entries[0].ContextMap()["property.incompatible_value_count"])
	assert.Equal(t, []any{"BAD999"}, entries[0].ContextMap()["property.resource_id_sample"])
}

func TestSetSystemPropertyEnabledRejectsCustomKeys(t *testing.T) {
	t.Parallel()

	db, _, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	store := New(db).(*Store)

	_, err = store.SetSystemPropertyEnabled(
		t.Context(),
		properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
		"vehicle.inspection_due",
		false,
	)
	require.ErrorIs(t, err, ErrInvalidPropertyDefinition)
}

func TestReconcileSystemPropertyOverridesRemovesStaleRows(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	store := New(db).(*Store)

	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_property_system_overrides")).
		WillReturnRows(sqlmock.NewRows([]string{
			"fivenet_property_system_overrides.scope",
			"fivenet_property_system_overrides.property_key",
		}).
			AddRow(vehiclePropertyScopeName, "vehicle.wanted").
			AddRow(vehiclePropertyScopeName, "vehicle.removed"))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM fivenet_property_system_overrides")).
		WithArgs(vehiclePropertyScopeName, "vehicle.removed").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_property_group_overrides")).
		WillReturnRows(sqlmock.NewRows([]string{
			"fivenet_property_group_overrides.scope",
			"fivenet_property_group_overrides.group_key",
		}))

	require.NoError(t, store.ReconcileSystemPropertyOverrides(t.Context()))
	require.NoError(t, mock.ExpectationsWereMet())
}
