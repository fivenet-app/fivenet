package vehiclesstore

import (
	"math"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	permsvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/vehicles/perms"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPropertyValueColumnsAndRowsRoundTrip(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		set  func(*properties.PropertyValue)
		want propertyValueType
	}{
		{
			"string",
			func(value *properties.PropertyValue) { value.SetStringValue("inspection") },
			propertyValueTypeString,
		},
		{
			"integer",
			func(value *properties.PropertyValue) { value.SetIntegerValue(42) },
			propertyValueTypeInteger,
		},
		{
			"number",
			func(value *properties.PropertyValue) { value.SetNumberValue(4.2) },
			propertyValueTypeNumber,
		},
		{
			"bool",
			func(value *properties.PropertyValue) { value.SetBoolValue(true) },
			propertyValueTypeBool,
		},
		{
			"timestamp",
			func(value *properties.PropertyValue) { value.SetTimestampValue(timestamp.New(now)) },
			propertyValueTypeTimestamp,
		},
		{
			"date",
			func(value *properties.PropertyValue) {
				value.SetDateValue(&properties.LocalDate{Year: 2026, Month: 10, Day: 21})
			},
			propertyValueTypeDate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := &properties.PropertyValue{}
			tt.set(value)
			property := &properties.Property{Key: "vehicle.inspection", Value: value}

			kind, stringVal, integerVal, numberVal, boolVal, timeVal, dateVal, err := propertyValueColumns(
				property,
			)
			require.NoError(t, err)
			assert.Equal(t, tt.want, kind)

			row := model.FivenetVehiclePropertyValues{
				PropertyKey: property.GetKey(), ValueType: uint8(kind), UpdatedAt: now,
			}
			switch kind {
			case propertyValueTypeString:
				value := stringVal.(string)
				row.ValueString = &value
			case propertyValueTypeInteger:
				value := integerVal.(int64)
				row.ValueInteger = &value
			case propertyValueTypeNumber:
				value := numberVal.(float64)
				row.ValueNumber = &value
			case propertyValueTypeBool:
				value := boolVal.(bool)
				row.ValueBool = &value
			case propertyValueTypeTimestamp:
				value := timeVal.(time.Time)
				row.ValueTimestamp = &value
			case propertyValueTypeDate:
				value := dateVal.(time.Time)
				row.ValueDate = &value
			}

			roundTripped, err := vehiclePropertyFromModel(row)
			require.NoError(t, err)
			assert.Equal(t, property.GetKey(), roundTripped.GetKey())
			assert.Equal(t, property.GetValue(), roundTripped.GetValue())
		})
	}
}

func TestPropertyValueColumnsRejectsReservedAndNonFiniteValues(t *testing.T) {
	t.Parallel()

	reserved := &properties.Property{
		Key:   string(permsvehicles.VehiclesServiceListVehiclesPropertiesPermValueVehicleWanted),
		Value: &properties.PropertyValue{},
	}
	reserved.Value.SetBoolValue(true)
	_, _, _, _, _, _, _, err := propertyValueColumns(reserved)
	require.ErrorIs(t, err, ErrInvalidProperty)

	nonFinite := &properties.Property{Key: "vehicle.score", Value: &properties.PropertyValue{}}
	nonFinite.Value.SetNumberValue(math.Inf(1))
	_, _, _, _, _, _, _, err = propertyValueColumns(nonFinite)
	require.ErrorIs(t, err, ErrInvalidProperty)
}

func TestIsGenericPropertyKeyExcludesSystemOrigin(t *testing.T) {
	t.Parallel()

	assert.False(
		t,
		IsGenericPropertyKey(
			string(permsvehicles.VehiclesServiceListVehiclesPropertiesPermValueVehicleWanted),
		),
	)
	assert.False(t, IsGenericPropertyKey("wanted"))
	assert.True(t, IsGenericPropertyKey("vehicle.inspection_score"))
}

func TestPropertyFilterConditionEscapesContainsValue(t *testing.T) {
	t.Parallel()

	value := &properties.PropertyValue{}
	value.SetStringValue("50%_off")
	condition, err := propertyFilterCondition(&properties.PropertyFilter{
		Key:       "vehicle.note",
		Operation: properties.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_CONTAINS,
		Value:     value,
	}, "vehicle", 0)
	require.NoError(t, err)

	vehicle := table.FivenetVehicle.AS("vehicle")
	stmt := vehicle.SELECT(vehicle.Plate).FROM(vehicle).WHERE(mysql.AND(condition))
	sql, args := stmt.Sql()
	assert.Contains(t, sql, "fivenet_vehicle_property_values AS property_filter_0")
	assert.Contains(t, sql, "LIKE ? ESCAPE '\\\\'")
	assert.Contains(t, args, "%50\\%\\_off%")
}

func TestStoreIsSystemPropertyEnabledQRMMapping(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	store := New(db, nil, nil).(*Store)

	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_property_system_overrides")).
		WithArgs("vehicle", string(permsvehicles.VehiclesServiceListVehiclesPropertiesPermValueVehicleWanted)).
		WillReturnRows(sqlmock.NewRows([]string{
			"fivenet_property_system_overrides.enabled",
		}).AddRow(false))

	enabled, err := store.IsSystemPropertyEnabled(
		t.Context(),
		string(permsvehicles.VehiclesServiceListVehiclesPropertiesPermValueVehicleWanted),
	)
	require.NoError(t, err)
	assert.False(t, enabled)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreIsSystemPropertyEnabledUsesDefinitionDefault(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	store := New(db, nil, nil).(*Store)

	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_property_system_overrides")).
		WithArgs("vehicle", string(permsvehicles.VehiclesServiceListVehiclesPropertiesPermValueVehicleWanted)).
		WillReturnRows(sqlmock.NewRows([]string{"fivenet_property_system_overrides.enabled"}))

	enabled, err := store.IsSystemPropertyEnabled(
		t.Context(),
		string(permsvehicles.VehiclesServiceListVehiclesPropertiesPermValueVehicleWanted),
	)
	require.NoError(t, err)
	assert.True(t, enabled)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListCustomPropertyKeysExcludesDisabledAndDeleted(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	store := New(db, nil, nil).(*Store)

	mock.ExpectQuery("SELECT .*property_key.*FROM fivenet_property_definitions.*").
		WithArgs("vehicle").
		WillReturnRows(sqlmock.NewRows([]string{"property_key"}).AddRow("vehicle.color"))

	keys, err := store.ListCustomPropertyKeys(t.Context())
	require.NoError(t, err)
	assert.Equal(t, []string{"vehicle.color"}, keys)
	require.NoError(t, mock.ExpectationsWereMet())
}
