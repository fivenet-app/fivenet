package vehiclesstore

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

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

func (s *Store) applyPropertyMutationStates(
	ctx context.Context,
	tx qrm.DB,
	plate string,
	mutations []*properties.PropertyMutation,
	creatorID *int32,
	creatorJob, reason string,
	activityMetadata map[string]*vehiclesactivity.PropertyChange,
) error {
	ordered := append([]*properties.PropertyMutation(nil), mutations...)
	slices.SortFunc(ordered, func(left, right *properties.PropertyMutation) int {
		if left == nil {
			if right == nil {
				return 0
			}
			return -1
		}
		if right == nil {
			return 1
		}
		return strings.Compare(left.GetKey(), right.GetKey())
	})
	for _, mutation := range ordered {
		if mutation == nil ||
			propertyregistry.IsSystemProperty(
				properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
				mutation.GetKey(),
			) {
			continue
		}
		definitions, err := s.ResolvePropertyDefinitions(ctx, []string{mutation.GetKey()})
		if err != nil {
			return err
		}
		if err := propertyregistry.ValidateMutationState(
			definitions[mutation.GetKey()],
			mutation,
		); err != nil {
			return fmt.Errorf("invalid state for property %q: %w", mutation.GetKey(), err)
		}
		if err := propertyregistry.ValidateMutationReason(
			definitions[mutation.GetKey()],
			mutation,
		); err != nil {
			return fmt.Errorf("invalid reason for property %q: %w", mutation.GetKey(), err)
		}
		if mutation.GetState() == nil && mutation.Reason == nil && mutation.ClearReason == nil {
			continue
		}
		if mutation.GetClearValue() != nil {
			continue
		}
		// Keep the value -> state order used by expiry cleanup.
		tValues := table.FivenetVehiclePropertyValues
		value := model.FivenetVehiclePropertyValues{}
		err = tValues.SELECT(tValues.Plate).FROM(tValues).WHERE(mysql.AND(
			tValues.Plate.EQ(mysql.String(plate)),
			tValues.PropertyKey.EQ(mysql.String(mutation.GetKey())),
		)).LIMIT(1).FOR(mysql.UPDATE()).QueryContext(ctx, tx, &value)
		if errors.Is(err, qrm.ErrNoRows) && mutation.GetSetValue() == nil {
			return fmt.Errorf("cannot set state for missing property %q", mutation.GetKey())
		}
		if err != nil && !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
		var existing *properties.Property
		if mutation.GetSetValue() == nil {
			propertiesByPlate, err := loadVehicleProperties(ctx, tx, []string{plate})
			if err != nil {
				return err
			}
			for _, property := range propertiesByPlate[plate] {
				if property.GetKey() == mutation.GetKey() {
					existing = property
					break
				}
			}
		}
		tState := table.FivenetVehiclePropertyState
		current := model.FivenetVehiclePropertyState{}
		err = tState.SELECT(tState.AllColumns).
			FROM(tState).
			WHERE(mysql.AND(tState.Plate.EQ(mysql.String(plate)), tState.PropertyKey.EQ(mysql.String(mutation.GetKey())))).
			LIMIT(1).
			FOR(mysql.UPDATE()).
			QueryContext(ctx, tx, &current)
		if err != nil && !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
		stateReason, activeAt, expiresAt := current.Reason, current.ActiveAt, current.ExpiresAt
		if activeAt == nil {
			now := time.Now().UTC()
			activeAt = &now
		}
		state := mutation.GetState()
		if mutation.Reason != nil {
			stateReason = mutation.Reason
		} else if mutation.ClearReason != nil {
			stateReason = nil
		}
		if state != nil && state.GetExpiresAt() != nil {
			value := state.GetExpiresAt().AsTime().UTC()
			expiresAt = &value
		} else if state != nil && state.GetClearExpiry() != nil {
			expiresAt = nil
		}
		reasonPatched := mutation.Reason != nil || mutation.ClearReason != nil
		expiryPatched := state != nil &&
			(state.GetExpiresAt() != nil || state.GetClearExpiry() != nil)
		reasonChanged := reasonPatched &&
			!propertyregistry.SameOptionalString(current.Reason, stateReason)
		expiryChanged := expiryPatched &&
			!propertyregistry.SameOptionalTime(current.ExpiresAt, expiresAt)
		if reasonChanged || expiryChanged {
			_, err = tState.INSERT(tState.Plate, tState.PropertyKey, tState.Reason, tState.ActiveAt, tState.ExpiresAt).
				VALUES(plate, mutation.GetKey(), stateReason, activeAt, expiresAt).
				ON_DUPLICATE_KEY_UPDATE(tState.Reason.SET(mysql.StringExp(mysql.RawString("VALUES(`reason`)"))), tState.ActiveAt.SET(mysql.TimestampExp(mysql.RawString("VALUES(`active_at`)"))), tState.ExpiresAt.SET(mysql.TimestampExp(mysql.RawString("VALUES(`expires_at`)")))).
				ExecContext(ctx, tx)
			if err != nil {
				return err
			}
		}
		activityMutation := mutation
		if mutation.Reason != nil && current.Reason != nil && *mutation.Reason == *current.Reason {
			activityMutation = proto.Clone(mutation).(*properties.PropertyMutation)
			activityMutation.Reason = nil
		}
		if mutation.ClearReason != nil && current.Reason == nil {
			activityMutation = proto.Clone(mutation).(*properties.PropertyMutation)
			activityMutation.ClearReason = nil
		}
		if expiryPatched && !expiryChanged {
			activityMutation = proto.Clone(activityMutation).(*properties.PropertyMutation)
			activityMutation.State = proto.Clone(activityMutation.GetState()).(*properties.PropertyStatePatch)
			activityMutation.State.ExpiryChange = nil
		}
		lifecycle := vehicleMutationLifecycle(activityMutation)
		if scalarActivity := activityMetadata[mutation.GetKey()]; scalarActivity != nil {
			// The scalar write below emits exactly one activity with the actual
			// lifecycle diff computed under the same locks.
			scalarActivity.Lifecycle = lifecycle
			continue
		}
		if existing == nil {
			propertiesByPlate, err := loadVehicleProperties(ctx, tx, []string{plate})
			if err != nil {
				return err
			}
			for _, property := range propertiesByPlate[plate] {
				if property.GetKey() == mutation.GetKey() {
					existing = property
					break
				}
			}
		}
		if existing != nil {
			if lifecycle == nil {
				continue
			}
			key := mutation.GetKey()
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
							Key:       key,
							OldValue:  proto.Clone(existing.GetValue()).(*properties.PropertyValue),
							NewValue:  proto.Clone(existing.GetValue()).(*properties.PropertyValue),
							Lifecycle: lifecycle,
							Cause:     vehiclesactivity.PropertyChangeCause_PROPERTY_CHANGE_CAUSE_MANUAL,
						},
					},
				},
			}
			if _, err := s.addVehicleActivity(ctx, tx, activity); err != nil {
				return err
			}
		}
	}
	return nil
}

// applyPropertyMutationsWithActivity is the vehicle equivalent of the citizen
// generic mutation path: lifecycle state is resolved under the value/state
// locks, then the scalar write emits the single merged activity.
func (s *Store) applyPropertyMutationsWithActivity(
	ctx context.Context,
	tx qrm.DB,
	plate string,
	values []*propertyregistry.ScalarChange,
	mutations []*properties.PropertyMutation,
	creatorID *int32,
	creatorJob, reason string,
	activityMetadata map[string]*vehiclesactivity.PropertyChange,
) error {
	if err := s.applyPropertyMutationStates(
		ctx, tx, plate, mutations, creatorID, creatorJob, reason, activityMetadata,
	); err != nil {
		return err
	}
	return s.applyPropertiesWithActivity(
		ctx, tx, plate, values, creatorID, creatorJob, reason, activityMetadata,
	)
}

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
	values []*propertyregistry.ScalarChange,
	creatorID *int32,
	creatorJob, reason string,
	activityMetadata map[string]*vehiclesactivity.PropertyChange,
) error {
	if len(values) == 0 {
		return nil
	}
	ordered := append([]*propertyregistry.ScalarChange(nil), values...)
	slices.SortFunc(ordered, func(left, right *propertyregistry.ScalarChange) int {
		return strings.Compare(left.Key, right.Key)
	})
	// Always lock value before state, including scalar-only and clear
	// mutations. This is the same order used by expiry cleanup.
	tValues := table.FivenetVehiclePropertyValues
	tState := table.FivenetVehiclePropertyState
	for _, value := range ordered {
		if value == nil {
			return ErrInvalidProperty
		}
		lockedValue := model.FivenetVehiclePropertyValues{}
		err := tValues.SELECT(tValues.Plate).FROM(tValues).WHERE(mysql.AND(
			tValues.Plate.EQ(mysql.String(plate)), tValues.PropertyKey.EQ(mysql.String(value.Key)),
		)).LIMIT(1).FOR(mysql.UPDATE()).QueryContext(ctx, tx, &lockedValue)
		if err != nil && !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
		lockedState := model.FivenetVehiclePropertyState{}
		err = tState.SELECT(tState.AllColumns).FROM(tState).WHERE(mysql.AND(
			tState.Plate.EQ(mysql.String(plate)), tState.PropertyKey.EQ(mysql.String(value.Key)),
		)).LIMIT(1).FOR(mysql.UPDATE()).QueryContext(ctx, tx, &lockedState)
		if err != nil && !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
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
	for _, value := range ordered {
		metadata := activityMetadata[value.Key]
		var lifecycle *vehiclesactivity.PropertyLifecycleChange
		var cause vehiclesactivity.PropertyChangeCause
		if metadata != nil {
			lifecycle, cause = metadata.GetLifecycle(), metadata.GetCause()
		}
		if value.Value == nil {
			if before := old[value.Key]; before != nil {
				tValues := table.FivenetVehiclePropertyValues
				if _, err := tValues.DELETE().
					WHERE(tValues.Plate.EQ(mysql.String(plate)).AND(tValues.PropertyKey.EQ(mysql.String(value.Key)))).
					ExecContext(ctx, tx); err != nil {
					return err
				}
				if _, err := tState.DELETE().WHERE(mysql.AND(
					tState.Plate.EQ(mysql.String(plate)),
					tState.PropertyKey.EQ(mysql.String(value.Key)),
				)).ExecContext(ctx, tx); err != nil {
					return err
				}
				key := value.Key
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
								Key:       key,
								OldValue:  proto.Clone(before.GetValue()).(*properties.PropertyValue),
								Lifecycle: lifecycle,
								Cause:     cause,
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
		property := &properties.Property{Key: value.Key, Value: value.Value}
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
						Key:       key,
						NewValue:  proto.Clone(property.GetValue()).(*properties.PropertyValue),
						Lifecycle: lifecycle,
						Cause:     cause,
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

func vehicleMutationLifecycle(
	mutation *properties.PropertyMutation,
) *vehiclesactivity.PropertyLifecycleChange {
	lifecycle := propertyregistry.MutationLifecycleForMutation(mutation)
	if lifecycle == nil {
		return nil
	}
	return &vehiclesactivity.PropertyLifecycleChange{
		ExpiresAt: lifecycle.ExpiresAt, ExpiryCleared: lifecycle.ExpiryCleared,
		StateReason: lifecycle.StateReason, StateReasonCleared: lifecycle.StateReasonCleared,
	}
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

// ListExpiringCustomPropertyKeys returns custom definitions whose explicit
// lifecycle policy clears their value when its state expires.
func (s *Store) ListExpiringCustomPropertyKeys(ctx context.Context) ([]string, error) {
	tDefinitions := table.FivenetPropertyDefinitions
	stmt := tDefinitions.SELECT(tDefinitions.AllColumns).
		FROM(tDefinitions).
		WHERE(mysql.AND(
			tDefinitions.Scope.EQ(mysql.String("vehicle")),
			tDefinitions.Enabled.IS_TRUE(),
			tDefinitions.DeletedAt.IS_NULL(),
		)).
		ORDER_BY(tDefinitions.PropertyKey.ASC())

	rows := []model.FivenetPropertyDefinitions{}
	if err := stmt.QueryContext(ctx, s.db, &rows); err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(rows))
	for _, row := range rows {
		definition, err := propertyregistry.DefinitionFromModel(
			properties.PropertyScope_PROPERTY_SCOPE_VEHICLE,
			row,
		)
		if err != nil {
			return nil, err
		}
		policy := definition.GetStatePolicy()
		if policy.GetEnabled() && policy.GetSupportsExpiry() &&
			policy.GetExpiryAction() == properties.PropertyExpiryAction_PROPERTY_EXPIRY_ACTION_CLEAR_VALUE {
			keys = append(keys, definition.GetKey())
		}
	}
	return keys, nil
}

func (s *Store) ListExpiredCustomPropertyPlates(
	ctx context.Context,
	key string,
	limit int64,
) ([]string, error) {
	tValues := table.FivenetVehiclePropertyValues.AS("value")
	tState := table.FivenetVehiclePropertyState.AS("state")
	stmt := tValues.SELECT(tValues.Plate).
		FROM(tValues.INNER_JOIN(tState, mysql.AND(
			tState.Plate.EQ(tValues.Plate),
			tState.PropertyKey.EQ(tValues.PropertyKey),
		))).
		WHERE(mysql.AND(
			tValues.PropertyKey.EQ(mysql.String(key)),
			tState.ExpiresAt.IS_NOT_NULL(),
			tState.ExpiresAt.LT_EQ(mysql.CURRENT_TIMESTAMP()),
		)).
		ORDER_BY(tState.ExpiresAt.ASC()).
		LIMIT(limit)

	plates := []string{}
	if err := stmt.QueryContext(ctx, s.db, &plates); err != nil && !errors.Is(err, qrm.ErrNoRows) {
		return nil, err
	}
	return plates, nil
}

func (s *Store) ClearExpiredCustomProperty(
	ctx context.Context,
	plate string,
	key string,
) (bool, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	// Lock value then state in the same order as normal mutations. This makes
	// a renewal commit before cleanup is evaluated, or wait until cleanup has
	// completed, so an expired batch entry cannot delete a renewed property.
	tValues := table.FivenetVehiclePropertyValues
	lockedValue := model.FivenetVehiclePropertyValues{}
	err = tValues.SELECT(tValues.Plate).FROM(tValues).WHERE(mysql.AND(
		tValues.Plate.EQ(mysql.String(plate)),
		tValues.PropertyKey.EQ(mysql.String(key)),
	)).LIMIT(1).FOR(mysql.UPDATE()).QueryContext(ctx, tx, &lockedValue)
	if errors.Is(err, qrm.ErrNoRows) {
		return false, tx.Commit()
	}
	if err != nil {
		return false, err
	}
	tState := table.FivenetVehiclePropertyState
	lockedState := model.FivenetVehiclePropertyState{}
	err = tState.SELECT(tState.AllColumns).FROM(tState).WHERE(mysql.AND(
		tState.Plate.EQ(mysql.String(plate)),
		tState.PropertyKey.EQ(mysql.String(key)),
	)).LIMIT(1).FOR(mysql.UPDATE()).QueryContext(ctx, tx, &lockedState)
	if errors.Is(err, qrm.ErrNoRows) {
		return false, tx.Commit()
	}
	if err != nil {
		return false, err
	}
	if lockedState.ExpiresAt == nil || lockedState.ExpiresAt.After(time.Now()) {
		return false, tx.Commit()
	}

	valuesByPlate, err := loadVehicleProperties(ctx, tx, []string{plate})
	if err != nil {
		return false, err
	}
	values := valuesByPlate[plate]
	for _, property := range values {
		if property.GetKey() != key {
			continue
		}
		if err := s.applyPropertiesWithActivity(ctx, tx, plate,
			[]*propertyregistry.ScalarChange{{Key: key}}, nil, "", "property_expired",
			map[string]*vehiclesactivity.PropertyChange{key: {
				Lifecycle: &vehiclesactivity.PropertyLifecycleChange{
					ExpiresAt: property.GetState().GetExpiresAt(),
				},
				Cause: vehiclesactivity.PropertyChangeCause_PROPERTY_CHANGE_CAUSE_EXPIRED,
			}}); err != nil {
			return false, err
		}
		if err := tx.Commit(); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, tx.Commit()
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
	if err := attachVehiclePropertyStates(ctx, db, plates, result); err != nil {
		return nil, err
	}

	return result, nil
}

func attachVehiclePropertyStates(
	ctx context.Context,
	db qrm.DB,
	plates []string,
	values map[string][]*properties.Property,
) error {
	if len(values) == 0 {
		return nil
	}
	plateExpressions := make([]mysql.Expression, 0, len(plates))
	for _, plate := range plates {
		plateExpressions = append(plateExpressions, mysql.String(plate))
	}
	byPlateAndKey := make(map[string]*properties.Property)
	for plate, properties := range values {
		for _, property := range properties {
			byPlateAndKey[plate+"\x00"+property.GetKey()] = property
		}
	}
	tState := table.FivenetVehiclePropertyState
	states := []model.FivenetVehiclePropertyState{}
	if err := tState.SELECT(tState.AllColumns).FROM(tState).
		WHERE(tState.Plate.IN(plateExpressions...)).QueryContext(ctx, db, &states); err != nil {
		return err
	}
	for _, stored := range states {
		property := byPlateAndKey[stored.Plate+"\x00"+stored.PropertyKey]
		if property == nil {
			continue
		}
		state := &properties.PropertyState{}
		state.SetKey(stored.PropertyKey)
		if stored.Reason != nil {
			state.SetReason(*stored.Reason)
		}
		if stored.ActiveAt != nil {
			state.SetActiveAt(timestamp.New(*stored.ActiveAt))
		}
		if stored.ExpiresAt != nil {
			state.SetExpiresAt(timestamp.New(*stored.ExpiresAt))
		}
		state.SetUpdatedAt(timestamp.New(stored.UpdatedAt))
		property.SetState(state)
	}
	return nil
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
