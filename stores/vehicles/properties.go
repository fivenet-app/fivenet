package vehiclesstore

import (
	"context"
	"errors"

	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	vehiclesactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/activity"
	propertyregistry "github.com/fivenet-app/fivenet/v2026/pkg/properties"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/protobuf/proto"
)

const vehiclePropertyValuesTable = "fivenet_vehicle_property_values"

var ErrInvalidProperty = errors.New("invalid vehicle property")

type propertyValueType = properties.PropertyValueType

const (
	propertyValueTypeString    = properties.PropertyValueType_PROPERTY_VALUE_TYPE_STRING
	propertyValueTypeInteger   = properties.PropertyValueType_PROPERTY_VALUE_TYPE_INTEGER
	propertyValueTypeNumber    = properties.PropertyValueType_PROPERTY_VALUE_TYPE_NUMBER
	propertyValueTypeBool      = properties.PropertyValueType_PROPERTY_VALUE_TYPE_BOOL
	propertyValueTypeTimestamp = properties.PropertyValueType_PROPERTY_VALUE_TYPE_TIMESTAMP
	propertyValueTypeDate      = properties.PropertyValueType_PROPERTY_VALUE_TYPE_DATE
)

func (s *Store) GetProperties(
	ctx context.Context,
	plate string,
) ([]*properties.Property, error) {
	byPlate, err := s.ListProperties(ctx, []string{plate})
	if err != nil {
		return nil, err
	}
	return byPlate[plate], nil
}

func (s *Store) applyPropertiesWithActivity(
	ctx context.Context,
	tx qrm.DB,
	plate string,
	values []*properties.PropertyChange,
	creatorID *int32,
	creatorJob, reason string,
) error {
	if len(values) == 0 {
		return nil
	}
	previousByPlate, err := loadVehicleProperties(ctx, tx, []string{plate})
	if err != nil {
		return err
	}
	previous := previousByPlate[plate]
	old := make(map[string]*properties.Property, len(previous))
	for _, value := range previous {
		old[value.GetKey()] = value
	}
	tValues := table.FivenetVehiclePropertyValues
	for _, value := range values {
		if value.GetValue() == nil {
			if before := old[value.GetKey()]; before != nil {
				tValues := table.FivenetVehiclePropertyValues
				if _, err := tValues.DELETE().
					WHERE(tValues.Plate.EQ(mysql.String(plate)).AND(tValues.PropertyKey.EQ(mysql.String(value.GetKey())))).
					ExecContext(ctx, tx); err != nil {
					return err
				}
				key := value.GetKey()
				activity := &vehiclesactivity.VehicleActivity{
					Plate:        plate,
					ActivityType: vehiclesactivity.VehicleActivityType_VEHICLE_ACTIVITY_TYPE_PROPERTY,
					CreatorId:    creatorID,
					CreatorJob:   creatorJob,
					Reason:       &reason,
					Key:          &key,
					Data: &vehiclesactivity.VehicleActivityData{
						Data: &vehiclesactivity.VehicleActivityData_PropertyChange{
							PropertyChange: &vehiclesactivity.PropertyChange{
								Key:      key,
								OldValue: proto.Clone(before.GetValue()).(*properties.PropertyValue),
							},
						},
					},
				}
				if _, err := s.addVehicleActivity(ctx, tx, activity); err != nil {
					return err
				}
			}
			continue
		}
		property := &properties.Property{Key: value.GetKey(), Value: value.GetValue()}
		valueType, stringVal, integerVal, numberVal, boolVal, timeVal, dateVal, err := propertyValueColumns(
			property,
		)
		if err != nil {
			return err
		}
		if _, err := tValues.INSERT(tValues.Plate, tValues.PropertyKey, tValues.ValueType, tValues.ValueString, tValues.ValueInteger, tValues.ValueNumber, tValues.ValueBool, tValues.ValueTimestamp, tValues.ValueDate).
			VALUES(plate, property.GetKey(), valueType, stringVal, integerVal, numberVal, boolVal, timeVal, dateVal).
			ON_DUPLICATE_KEY_UPDATE(tValues.ValueType.SET(mysql.IntExp(mysql.RawString("VALUES(`value_type`)"))), tValues.ValueString.SET(mysql.StringExp(mysql.RawString("VALUES(`value_string`)"))), tValues.ValueInteger.SET(mysql.IntExp(mysql.RawString("VALUES(`value_integer`)"))), tValues.ValueNumber.SET(mysql.FloatExp(mysql.RawString("VALUES(`value_number`)"))), tValues.ValueBool.SET(mysql.BoolExp(mysql.RawString("VALUES(`value_bool`)"))), tValues.ValueTimestamp.SET(mysql.TimestampExp(mysql.RawString("VALUES(`value_timestamp`)"))), tValues.ValueDate.SET(mysql.DateExp(mysql.RawString("VALUES(`value_date`)")))).
			ExecContext(ctx, tx); err != nil {
			return err
		}
		if before := old[property.GetKey()]; before != nil &&
			proto.Equal(before.GetValue(), property.GetValue()) {
			continue
		}
		key := property.GetKey()
		activity := &vehiclesactivity.VehicleActivity{
			Plate:        plate,
			ActivityType: vehiclesactivity.VehicleActivityType_VEHICLE_ACTIVITY_TYPE_PROPERTY,
			CreatorId:    creatorID,
			CreatorJob:   creatorJob,
			Reason:       &reason,
			Key:          &key,
			Data: &vehiclesactivity.VehicleActivityData{
				Data: &vehiclesactivity.VehicleActivityData_PropertyChange{
					PropertyChange: &vehiclesactivity.PropertyChange{
						Key:      key,
						NewValue: proto.Clone(property.GetValue()).(*properties.PropertyValue),
					},
				},
			},
		}
		if before := old[key]; before != nil {
			activity.Data.GetPropertyChange().OldValue = proto.Clone(before.GetValue()).(*properties.PropertyValue)
		}
		if _, err := s.addVehicleActivity(ctx, tx, activity); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) ListProperties(
	ctx context.Context,
	plates []string,
) (map[string][]*properties.Property, error) {
	return loadVehicleProperties(ctx, s.db, plates)
}

func (s *Store) ListCustomPropertyKeys(ctx context.Context) ([]string, error) {
	tDefinitions := table.FivenetPropertyDefinitions
	stmt := tDefinitions.SELECT(tDefinitions.PropertyKey).
		FROM(tDefinitions).
		WHERE(mysql.AND(
			tDefinitions.Scope.EQ(mysql.String("vehicle")),
			tDefinitions.Enabled.IS_TRUE(),
			tDefinitions.DeletedAt.IS_NULL(),
		)).
		ORDER_BY(tDefinitions.PropertyKey.ASC())

	var keys []string
	if err := stmt.QueryContext(ctx, s.db, &keys); err != nil {
		return nil, err
	}
	return keys, nil
}

func (s *Store) ResolvePropertyDefinitions(
	ctx context.Context,
	keys []string,
) (map[string]*properties.PropertyDefinition, error) {
	return s.propertyCache.Resolve(
		ctx,
		s.db,
		properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
		keys,
	)
}

func (s *Store) ResolvePropertyGroups(
	ctx context.Context,
	definitions map[string]*properties.PropertyDefinition,
) ([]*properties.PropertyGroup, error) {
	scope := properties.PropertyScope_PROPERTY_SCOPE_VEHICLE
	return s.propertyCache.Groups(ctx, s.db, scope, definitions)
}

func (s *Store) ResolvePropertyCatalog(
	ctx context.Context,
	keys []string,
) (*propertyregistry.ResolvedCatalog, error) {
	return s.propertyCache.ResolveCatalog(
		ctx, s.db, properties.PropertyScope_PROPERTY_SCOPE_VEHICLE, keys,
	)
}

func loadVehicleProperties(
	ctx context.Context,
	db qrm.DB,
	plates []string,
) (map[string][]*properties.Property, error) {
	result := make(map[string][]*properties.Property, len(plates))
	if len(plates) == 0 {
		return result, nil
	}

	plateExpressions := make([]mysql.Expression, 0, len(plates))
	for _, plate := range plates {
		plateExpressions = append(plateExpressions, mysql.String(plate))
	}

	tValues := table.FivenetVehiclePropertyValues
	tDefinitions := table.FivenetPropertyDefinitions.AS("property_definition")
	deletedDefinitionExists := mysql.EXISTS(
		mysql.SELECT(mysql.Int(1)).FROM(tDefinitions).WHERE(mysql.AND(
			tDefinitions.Scope.EQ(mysql.String("vehicle")),
			tDefinitions.PropertyKey.EQ(tValues.PropertyKey),
			tDefinitions.DeletedAt.IS_NOT_NULL(),
		)),
	)
	stmt := tValues.
		SELECT(tValues.AllColumns).
		FROM(tValues).
		WHERE(tValues.Plate.IN(plateExpressions...).AND(mysql.NOT(deletedDefinitionExists))).
		ORDER_BY(tValues.Plate.ASC(), tValues.PropertyKey.ASC())
	rows := []model.FivenetVehiclePropertyValues{}
	if err := stmt.QueryContext(ctx, db, &rows); err != nil {
		return nil, err
	}

	for _, row := range rows {
		property, err := vehiclePropertyFromModel(row)
		if err != nil {
			return nil, err
		}
		result[row.Plate] = append(result[row.Plate], property)
	}

	return result, nil
}

func vehiclePropertyFromModel(
	stored model.FivenetVehiclePropertyValues,
) (*properties.Property, error) {
	value, err := propertyregistry.DecodeStoredValue(propertyregistry.StoredValue{
		Type: properties.PropertyValueType(
			stored.ValueType,
		),
		String:    stored.ValueString,
		Integer:   stored.ValueInteger,
		Number:    stored.ValueNumber,
		Bool:      stored.ValueBool,
		Timestamp: stored.ValueTimestamp,
		Date:      stored.ValueDate,
	})
	if err != nil {
		return nil, ErrInvalidProperty
	}

	return &properties.Property{
		Key:       stored.PropertyKey,
		Value:     value,
		UpdatedAt: timestamp.New(stored.UpdatedAt),
	}, nil
}

func propertyValueColumns(
	property *properties.Property,
) (propertyValueType, any, any, any, any, any, any, error) {
	if property == nil || !isVehiclePropertyKey(property.GetKey()) || property.GetValue() == nil {
		return 0, nil, nil, nil, nil, nil, nil, ErrInvalidProperty
	}

	typeValue, stringValue, integerValue, numberValue, boolValue, timestampValue, dateValue, err := propertyregistry.ValueColumns(
		property.GetValue(),
	)
	if err != nil {
		return 0, nil, nil, nil, nil, nil, nil, ErrInvalidProperty
	}
	return typeValue, stringValue, integerValue, numberValue, boolValue, timestampValue, dateValue, nil
}

func isVehiclePropertyKey(key string) bool {
	return IsGenericPropertyKey(key)
}

func propertyFilterCondition(
	filter *properties.PropertyFilter,
	vehicleAlias string,
	index int,
) (mysql.BoolExpression, error) {
	return propertyregistry.FilterCondition(
		filter,
		propertyregistry.FilterConfig{
			Scope:             properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
			ScopeName:         "vehicle",
			ValuesTable:       vehiclePropertyValuesTable,
			ParentAlias:       vehicleAlias,
			ValueParentColumn: "plate", ParentColumn: "plate",
		},
		index,
	)
}
