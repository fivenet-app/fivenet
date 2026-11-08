package citizensstore

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
	usersprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/props"
	_ "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens"
	propertyregistry "github.com/fivenet-app/fivenet/v2026/pkg/properties"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/protobuf/proto"
)

const userPropertyValuesTable = "fivenet_user_property_values"

var ErrInvalidProperty = errors.New("invalid user property")

// PropertyChange is an atomic canonical value update. File relations and
// managed properties with lifecycle state (such as user.wanted) are handled
// by their dedicated handlers rather than being represented as scalar values.
type PropertyChange struct {
	Property  *properties.Property
	Delete    bool
	Lifecycle *usersactivity.PropertyLifecycleChange
	Cause     usersactivity.PropertyChangeCause
}

type userPropertyValueType = properties.PropertyValueType

const (
	userPropertyValueTypeString    = properties.PropertyValueType_PROPERTY_VALUE_TYPE_STRING
	userPropertyValueTypeInteger   = properties.PropertyValueType_PROPERTY_VALUE_TYPE_INTEGER
	userPropertyValueTypeNumber    = properties.PropertyValueType_PROPERTY_VALUE_TYPE_NUMBER
	userPropertyValueTypeBool      = properties.PropertyValueType_PROPERTY_VALUE_TYPE_BOOL
	userPropertyValueTypeTimestamp = properties.PropertyValueType_PROPERTY_VALUE_TYPE_TIMESTAMP
	userPropertyValueTypeDate      = properties.PropertyValueType_PROPERTY_VALUE_TYPE_DATE
)

func (s *Store) GetProperties(ctx context.Context, userID int32) ([]*properties.Property, error) {
	return s.getProperties(ctx, s.db, userID)
}

func (s *Store) ListProperties(
	ctx context.Context,
	userIDs []int32,
) (map[int32][]*properties.Property, error) {
	result := make(map[int32][]*properties.Property, len(userIDs))
	if len(userIDs) == 0 {
		return result, nil
	}

	ids := make([]mysql.Expression, 0, len(userIDs))
	for _, userID := range userIDs {
		ids = append(ids, mysql.Int32(userID))
	}

	tValues := table.FivenetUserPropertyValues
	tDefinitions := table.FivenetPropertyDefinitions.AS("property_definition")
	deletedDefinitionExists := mysql.EXISTS(
		mysql.SELECT(mysql.Int(1)).FROM(tDefinitions).WHERE(mysql.AND(
			tDefinitions.Scope.EQ(mysql.String("user")),
			tDefinitions.PropertyKey.EQ(tValues.PropertyKey),
			tDefinitions.DeletedAt.IS_NOT_NULL(),
		)),
	)
	stmt := tValues.
		SELECT(tValues.AllColumns).
		FROM(tValues).
		WHERE(tValues.UserID.IN(ids...).AND(mysql.NOT(deletedDefinitionExists))).
		ORDER_BY(tValues.UserID.ASC(), tValues.PropertyKey.ASC())
	rows := []model.FivenetUserPropertyValues{}
	if err := stmt.QueryContext(ctx, s.db, &rows); err != nil {
		return nil, err
	}

	for _, stored := range rows {
		value, err := userPropertyFromModel(stored)
		if err != nil {
			return nil, err
		}
		result[stored.UserID] = append(result[stored.UserID], value)
	}
	if err := attachUserPropertyStatesForUsers(ctx, s.db, userIDs, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Store) GetPropertiesTx(
	ctx context.Context,
	tx qrm.DB,
	userID int32,
) ([]*properties.Property, error) {
	return s.getProperties(ctx, tx, userID)
}

func (s *Store) ListCustomPropertyKeys(ctx context.Context) ([]string, error) {
	tDefinitions := table.FivenetPropertyDefinitions
	stmt := tDefinitions.SELECT(tDefinitions.PropertyKey).
		FROM(tDefinitions).
		WHERE(mysql.AND(
			tDefinitions.Scope.EQ(mysql.String("user")),
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

// ListExpiringCustomPropertyKeys returns only definitions which explicitly opt
// into the generic expiry lifecycle. Managed properties keep their dedicated
// cleanup paths and are deliberately not considered here.
func (s *Store) ListExpiringCustomPropertyKeys(ctx context.Context) ([]string, error) {
	tDefinitions := table.FivenetPropertyDefinitions
	stmt := tDefinitions.SELECT(tDefinitions.AllColumns).
		FROM(tDefinitions).
		WHERE(mysql.AND(
			tDefinitions.Scope.EQ(mysql.String("user")),
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
			properties.PropertyScope_PROPERTY_SCOPE_USER,
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

func (s *Store) ListExpiredCustomPropertyUserIDs(
	ctx context.Context,
	key string,
	limit int64,
) ([]int32, error) {
	tValues := table.FivenetUserPropertyValues.AS("value")
	tState := table.FivenetUserPropertyState.AS("state")
	stmt := tValues.SELECT(tValues.UserID).
		FROM(tValues.INNER_JOIN(tState, mysql.AND(
			tState.UserID.EQ(tValues.UserID),
			tState.PropertyKey.EQ(tValues.PropertyKey),
		))).
		WHERE(mysql.AND(
			tValues.PropertyKey.EQ(mysql.String(key)),
			tState.ExpiresAt.IS_NOT_NULL(),
			tState.ExpiresAt.LT_EQ(mysql.CURRENT_TIMESTAMP()),
		)).
		ORDER_BY(tState.ExpiresAt.ASC()).
		LIMIT(limit)

	ids := []int32{}
	if err := stmt.QueryContext(ctx, s.db, &ids); err != nil && !errors.Is(err, qrm.ErrNoRows) {
		return nil, err
	}
	return ids, nil
}

// ClearExpiredCustomProperty removes an expired scalar value and its state in
// one transaction, recording the ordinary generic property-change activity.
// It rechecks the projected state so a stale batch entry is harmless.
func (s *Store) ClearExpiredCustomProperty(
	ctx context.Context,
	userID int32,
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
	tValues := table.FivenetUserPropertyValues
	lockedValue := model.FivenetUserPropertyValues{}
	err = tValues.SELECT(tValues.UserID).FROM(tValues).WHERE(mysql.AND(
		tValues.UserID.EQ(mysql.Int32(userID)),
		tValues.PropertyKey.EQ(mysql.String(key)),
	)).LIMIT(1).FOR(mysql.UPDATE()).QueryContext(ctx, tx, &lockedValue)
	if errors.Is(err, qrm.ErrNoRows) {
		return false, tx.Commit()
	}
	if err != nil {
		return false, err
	}
	tState := table.FivenetUserPropertyState
	lockedState := model.FivenetUserPropertyState{}
	err = tState.SELECT(tState.AllColumns).FROM(tState).WHERE(mysql.AND(
		tState.UserID.EQ(mysql.Int32(userID)),
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

	values, err := s.GetPropertiesTx(ctx, tx, userID)
	if err != nil {
		return false, err
	}
	var property *properties.Property
	for _, value := range values {
		if value.GetKey() == key {
			property = value
			break
		}
	}
	if property == nil {
		return false, tx.Commit()
	}

	activities, err := s.ApplyPropertyChangesWithActivity(ctx, tx, userID, []PropertyChange{{
		Property: property,
		Delete:   true,
		Lifecycle: &usersactivity.PropertyLifecycleChange{
			ExpiresAt: property.GetState().GetExpiresAt(),
		},
		Cause: usersactivity.PropertyChangeCause_PROPERTY_CHANGE_CAUSE_EXPIRED,
	}}, nil, "property_expired")
	if err != nil {
		return false, err
	}
	for _, activity := range activities {
		if err := usersactivity.CreateUserActivities(ctx, tx, activity); err != nil {
			return false, err
		}
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

func (s *Store) ResolvePropertyDefinitions(
	ctx context.Context,
	keys []string,
) (map[string]*properties.PropertyDefinition, error) {
	return s.propertyCache.Resolve(
		ctx,
		s.db,
		properties.PropertyScope_PROPERTY_SCOPE_USER,
		keys,
	)
}

func (s *Store) ResolvePropertyGroups(
	ctx context.Context,
	definitions map[string]*properties.PropertyDefinition,
) ([]*properties.PropertyGroup, error) {
	scope := properties.PropertyScope_PROPERTY_SCOPE_USER
	return s.propertyCache.Groups(ctx, s.db, scope, definitions)
}

func (s *Store) ResolvePropertyCatalog(
	ctx context.Context,
	keys []string,
) (*propertyregistry.ResolvedCatalog, error) {
	return s.propertyCache.ResolveCatalog(
		ctx, s.db, properties.PropertyScope_PROPERTY_SCOPE_USER, keys,
	)
}

func (s *Store) getProperties(
	ctx context.Context,
	db qrm.DB,
	userID int32,
) ([]*properties.Property, error) {
	tValues := table.FivenetUserPropertyValues
	tDefinitions := table.FivenetPropertyDefinitions.AS("property_definition")
	deletedDefinitionExists := mysql.EXISTS(
		mysql.SELECT(mysql.Int(1)).FROM(tDefinitions).WHERE(mysql.AND(
			tDefinitions.Scope.EQ(mysql.String("user")),
			tDefinitions.PropertyKey.EQ(tValues.PropertyKey),
			tDefinitions.DeletedAt.IS_NOT_NULL(),
		)),
	)
	stmt := tValues.
		SELECT(tValues.AllColumns).
		FROM(tValues).
		WHERE(tValues.UserID.EQ(mysql.Int32(userID)).AND(mysql.NOT(deletedDefinitionExists))).
		ORDER_BY(tValues.PropertyKey.ASC())
	rows := []model.FivenetUserPropertyValues{}
	err := stmt.QueryContext(ctx, db, &rows)
	if err != nil {
		return nil, err
	}

	values := make([]*properties.Property, 0, len(rows))
	for _, stored := range rows {
		value, err := userPropertyFromModel(stored)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	if err := attachUserPropertyStates(ctx, db, userID, values); err != nil {
		return nil, err
	}
	return values, nil
}

// attachUserPropertyStates projects lifecycle rows onto their scalar property.
// State is intentionally keyed independently in storage, but is exposed with
// the value so consumers never need a second property-state API.
func attachUserPropertyStates(
	ctx context.Context,
	db qrm.DB,
	userID int32,
	values []*properties.Property,
) error {
	return attachUserPropertyStatesForUsers(
		ctx,
		db,
		[]int32{userID},
		map[int32][]*properties.Property{userID: values},
	)
}

// attachUserPropertyStatesForUsers projects lifecycle rows onto a batch of
// scalar properties. List responses use this to avoid an N+1 state query.
func attachUserPropertyStatesForUsers(
	ctx context.Context,
	db qrm.DB,
	userIDs []int32,
	valuesByUserID map[int32][]*properties.Property,
) error {
	if len(userIDs) == 0 || len(valuesByUserID) == 0 {
		return nil
	}
	ids := make([]mysql.Expression, 0, len(userIDs))
	propertiesByUserAndKey := make(map[int32]map[string]*properties.Property, len(valuesByUserID))
	for _, userID := range userIDs {
		values := valuesByUserID[userID]
		if len(values) == 0 {
			continue
		}
		ids = append(ids, mysql.Int32(userID))
		byKey := make(map[string]*properties.Property, len(values))
		for _, value := range values {
			byKey[value.GetKey()] = value
		}
		propertiesByUserAndKey[userID] = byKey
	}
	if len(ids) == 0 {
		return nil
	}
	tState := table.FivenetUserPropertyState
	states := []model.FivenetUserPropertyState{}
	if err := tState.SELECT(tState.AllColumns).FROM(tState).
		WHERE(tState.UserID.IN(ids...)).QueryContext(ctx, db, &states); err != nil {
		return err
	}
	for _, stored := range states {
		property := propertiesByUserAndKey[stored.UserID][stored.PropertyKey]
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

// ApplyPropertyChanges applies scalar user-property changes through the
// canonical values table. The caller supplies its transaction when the update
// must be coordinated with an activity or another system mutation.
func (s *Store) ApplyPropertyChanges(
	ctx context.Context,
	db qrm.DB,
	userID int32,
	changes []PropertyChange,
) ([]*properties.Property, error) {
	if len(changes) == 0 {
		return []*properties.Property{}, nil
	}

	updated := make([]*properties.Property, 0, len(changes))
	for _, change := range changes {
		if change.Property == nil || !isUserPropertyKey(change.Property.GetKey()) {
			return nil, ErrInvalidProperty
		}

		if change.Delete {
			tValues := table.FivenetUserPropertyValues
			if _, err := tValues.DELETE().WHERE(mysql.AND(
				tValues.UserID.EQ(mysql.Int32(userID)),
				tValues.PropertyKey.EQ(mysql.String(change.Property.GetKey())),
			)).ExecContext(ctx, db); err != nil {
				return nil, err
			}
			tState := table.FivenetUserPropertyState
			if _, err := tState.DELETE().WHERE(mysql.AND(
				tState.UserID.EQ(mysql.Int32(userID)),
				tState.PropertyKey.EQ(mysql.String(change.Property.GetKey())),
			)).ExecContext(ctx, db); err != nil {
				return nil, err
			}
			continue
		}

		valueType, stringVal, integerVal, numberVal, boolVal, timeVal, dateVal, err := userPropertyValueColumns(
			change.Property,
		)
		if err != nil {
			return nil, err
		}
		if err := upsertUserPropertyValue(ctx, db, userID, change.Property.GetKey(), valueType,
			stringVal, integerVal, numberVal, boolVal, timeVal, dateVal); err != nil {
			return nil, err
		}
		updated = append(updated, change.Property)
	}

	if len(updated) == 0 {
		return updated, nil
	}
	persisted, err := s.getProperties(ctx, db, userID)
	if err != nil {
		return nil, err
	}
	byKey := make(map[string]*properties.Property, len(persisted))
	for _, property := range persisted {
		byKey[property.GetKey()] = property
	}
	for i, property := range updated {
		value, ok := byKey[property.GetKey()]
		if !ok {
			return nil, fmt.Errorf("%w: property was not persisted", ErrInvalidProperty)
		}
		updated[i] = value
	}
	return updated, nil
}

// ApplyPropertyChangesWithActivity applies scalar property changes and creates
// the corresponding user activity entries using the same transaction.
func (s *Store) ApplyPropertyChangesWithActivity(
	ctx context.Context,
	db qrm.DB,
	userID int32,
	changes []PropertyChange,
	sourceUserID *int32,
	reason string,
) ([]*usersactivity.UserActivity, error) {
	return s.applyPropertyChangesWithActivity(ctx, db, userID, changes, sourceUserID, reason, nil)
}

func (s *Store) applyPropertyChangesWithActivity(
	ctx context.Context,
	db qrm.DB,
	userID int32,
	changes []PropertyChange,
	sourceUserID *int32,
	reason string,
	lockedKeys map[string]struct{},
) ([]*usersactivity.UserActivity, error) {
	if len(changes) == 0 {
		return nil, nil
	}

	// Always take the value lock before the lifecycle lock.  Expiry cleanup
	// takes the same locks in this order, so value and state cannot be observed
	// from different versions of a mutation.
	ordered := append([]PropertyChange(nil), changes...)
	for _, change := range ordered {
		if change.Property == nil {
			return nil, ErrInvalidProperty
		}
	}
	slices.SortFunc(ordered, func(left, right PropertyChange) int {
		return strings.Compare(left.Property.GetKey(), right.Property.GetKey())
	})
	tValues := table.FivenetUserPropertyValues
	tState := table.FivenetUserPropertyState
	for _, change := range ordered {
		if _, locked := lockedKeys[change.Property.GetKey()]; locked {
			continue
		}
		locked := model.FivenetUserPropertyValues{}
		err := tValues.SELECT(tValues.UserID).FROM(tValues).WHERE(mysql.AND(
			tValues.UserID.EQ(mysql.Int32(userID)),
			tValues.PropertyKey.EQ(mysql.String(change.Property.GetKey())),
		)).LIMIT(1).FOR(mysql.UPDATE()).QueryContext(ctx, db, &locked)
		if err != nil && !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
		// Scalar-only and clear mutations also lock and read the state row. A
		// lifecycle patch later in this transaction reuses this lock.
		lockedState := model.FivenetUserPropertyState{}
		err = tState.SELECT(tState.AllColumns).FROM(tState).WHERE(mysql.AND(
			tState.UserID.EQ(mysql.Int32(userID)),
			tState.PropertyKey.EQ(mysql.String(change.Property.GetKey())),
		)).LIMIT(1).FOR(mysql.UPDATE()).QueryContext(ctx, db, &lockedState)
		if err != nil && !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}
	previous, err := s.GetPropertiesTx(ctx, db, userID)
	if err != nil {
		return nil, err
	}
	if _, err := s.ApplyPropertyChanges(ctx, db, userID, changes); err != nil {
		return nil, err
	}

	return UserPropertyActivities(userID, sourceUserID, reason, previous, changes), nil
}

// mergePropertyMutationActivities folds lifecycle changes into their scalar
// activity.  Scalar activities are created before lifecycle persistence so we
// can retain one stable audit record, while the lifecycle routine supplies the
// effective (rather than requested) lifecycle diff afterwards.
func mergePropertyMutationActivities(
	scalar, lifecycle []*usersactivity.UserActivity,
) []*usersactivity.UserActivity {
	byKey := make(map[string]*usersactivity.UserActivity, len(scalar))
	for _, activity := range scalar {
		if change := activity.GetData().GetPropertyChange(); change != nil {
			byKey[change.GetKey()] = activity
			// Requested lifecycle data is not evidence of a state change. A
			// clear is the exception: its mutation-level rationale is retained.
			if change.GetNewValue() != nil {
				change.Lifecycle = nil
			}
		}
	}
	stateOnly := make([]*usersactivity.UserActivity, 0, len(lifecycle))
	for _, activity := range lifecycle {
		change := activity.GetData().GetPropertyChange()
		if change != nil && byKey[change.GetKey()] != nil {
			byKey[change.GetKey()].GetData().GetPropertyChange().Lifecycle = change.GetLifecycle()
			continue
		}
		stateOnly = append(stateOnly, activity)
	}
	merged := make([]*usersactivity.UserActivity, 0, len(scalar)+len(stateOnly))
	for _, activity := range scalar {
		change := activity.GetData().GetPropertyChange()
		if change != nil && change.GetLifecycle() == nil &&
			proto.Equal(change.GetOldValue(), change.GetNewValue()) {
			// A requested lifecycle patch may have forced construction of this
			// scalar activity. Once the effective state diff is known to be
			// empty, suppress the resulting scalar no-op as well.
			continue
		}
		merged = append(merged, activity)
	}
	merged = append(merged, stateOnly...)
	return merged
}

// ApplyPropertyMutationsWithActivity is the single generic-property mutation
// path. It keeps value/state locking, effective lifecycle diffing, persistence
// and the corresponding one-per-property activity in one transaction API.
func (s *Store) ApplyPropertyMutationsWithActivity(
	ctx context.Context,
	tx qrm.DB,
	userID int32,
	changes []PropertyChange,
	mutations []*properties.PropertyMutation,
	sourceUserID *int32,
	reason string,
) ([]*usersactivity.UserActivity, error) {
	lockedKeys, err := s.lockPropertyMutationKeys(ctx, tx, userID, changes, mutations)
	if err != nil {
		return nil, err
	}
	mathChanges, err := s.resolveIntegerMathMutations(ctx, tx, userID, mutations)
	if err != nil {
		return nil, err
	}
	changes = append(changes, mathChanges...)
	scalarActivities, err := s.applyPropertyChangesWithActivity(
		ctx, tx, userID, changes, sourceUserID, reason, lockedKeys,
	)
	if err != nil {
		return nil, err
	}
	lifecycleActivities, err := s.applyPropertyMutationStates(
		ctx, tx, userID, mutations, sourceUserID, reason,
	)
	if err != nil {
		return nil, err
	}
	return mergePropertyMutationActivities(scalarActivities, lifecycleActivities), nil
}

func (s *Store) lockPropertyMutationKeys(
	ctx context.Context,
	tx qrm.DB,
	userID int32,
	changes []PropertyChange,
	mutations []*properties.PropertyMutation,
) (map[string]struct{}, error) {
	keys := map[string]struct{}{}
	for _, change := range changes {
		if change.Property == nil {
			return nil, ErrInvalidProperty
		}
		keys[change.Property.GetKey()] = struct{}{}
	}
	for _, mutation := range mutations {
		if mutation != nil {
			keys[mutation.GetKey()] = struct{}{}
		}
	}
	ordered := make([]string, 0, len(keys))
	for key := range keys {
		ordered = append(ordered, key)
	}
	slices.Sort(ordered)
	tValues, tState := table.FivenetUserPropertyValues, table.FivenetUserPropertyState
	for _, key := range ordered {
		value := model.FivenetUserPropertyValues{}
		if err := tValues.SELECT(tValues.UserID).
			FROM(tValues).
			WHERE(mysql.AND(tValues.UserID.EQ(mysql.Int32(userID)), tValues.PropertyKey.EQ(mysql.String(key)))).
			LIMIT(1).
			FOR(mysql.UPDATE()).
			QueryContext(ctx, tx, &value); err != nil &&
			!errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
		state := model.FivenetUserPropertyState{}
		if err := tState.SELECT(tState.AllColumns).
			FROM(tState).
			WHERE(mysql.AND(tState.UserID.EQ(mysql.Int32(userID)), tState.PropertyKey.EQ(mysql.String(key)))).
			LIMIT(1).
			FOR(mysql.UPDATE()).
			QueryContext(ctx, tx, &state); err != nil &&
			!errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}
	return keys, nil
}

// resolveIntegerMathMutations derives canonical scalar changes only after the
// value row is locked. This prevents a read/modify/write race across Sync
// clients and server processes.
func (s *Store) resolveIntegerMathMutations(
	ctx context.Context,
	tx qrm.DB,
	userID int32,
	mutations []*properties.PropertyMutation,
) ([]PropertyChange, error) {
	mathMutations := make([]*properties.PropertyMutation, 0, len(mutations))
	for _, mutation := range mutations {
		if mutation != nil && mutation.GetIntegerMath() != nil {
			mathMutations = append(mathMutations, mutation)
		}
	}
	if len(mathMutations) == 0 {
		return nil, nil
	}
	slices.SortFunc(mathMutations, func(left, right *properties.PropertyMutation) int {
		return strings.Compare(left.GetKey(), right.GetKey())
	})

	changes := make([]PropertyChange, 0, len(mathMutations))
	stored, err := s.GetPropertiesTx(ctx, tx, userID)
	if err != nil {
		return nil, err
	}
	byKey := make(map[string]*properties.Property, len(stored))
	for _, property := range stored {
		byKey[property.GetKey()] = property
	}
	for _, mutation := range mathMutations {
		definitions, err := s.ResolvePropertyDefinitions(ctx, []string{mutation.GetKey()})
		if err != nil {
			return nil, err
		}
		definition := definitions[mutation.GetKey()]
		if definition == nil {
			return nil, ErrInvalidProperty
		}

		current := int64(0)
		if property := byKey[mutation.GetKey()]; property != nil {
			if property.GetValue().GetIntegerValue() == 0 && property.GetValue().GetValue() == nil {
				return nil, ErrInvalidProperty
			}
			integer, ok := property.GetValue().GetValue().(*properties.PropertyValue_IntegerValue)
			if !ok {
				return nil, ErrInvalidProperty
			}
			current = integer.IntegerValue
		}
		value, err := propertyregistry.ApplyIntegerMath(
			definition,
			current,
			mutation.GetIntegerMath(),
		)
		if err != nil {
			return nil, ErrInvalidProperty
		}
		propertyValue := &properties.PropertyValue{}
		propertyValue.SetIntegerValue(value)
		changes = append(changes, PropertyChange{
			Property:  &properties.Property{Key: mutation.GetKey(), Value: propertyValue},
			Cause:     usersactivity.PropertyChangeCause_PROPERTY_CHANGE_CAUSE_MANUAL,
			Lifecycle: userMutationLifecycle(mutation),
		})
	}
	return changes, nil
}

// applyPropertyMutationStates persists lifecycle patches after their scalar
// mutations have been applied in the same transaction. State-only mutations
// are permitted only for an existing scalar property.
func (s *Store) applyPropertyMutationStates(
	ctx context.Context,
	tx qrm.DB,
	userID int32,
	mutations []*properties.PropertyMutation,
	sourceUserID *int32,
	reason string,
) ([]*usersactivity.UserActivity, error) {
	activities := make([]*usersactivity.UserActivity, 0, len(mutations))
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
				properties.PropertyScope_PROPERTY_SCOPE_USER,
				mutation.GetKey(),
			) {
			continue
		}
		definitions, err := s.ResolvePropertyDefinitions(ctx, []string{mutation.GetKey()})
		if err != nil {
			return nil, err
		}
		if err := propertyregistry.ValidateMutationState(
			definitions[mutation.GetKey()],
			mutation,
		); err != nil {
			return nil, fmt.Errorf("invalid state for property %q: %w", mutation.GetKey(), err)
		}
		if err := propertyregistry.ValidateMutationReason(
			definitions[mutation.GetKey()],
			mutation,
		); err != nil {
			return nil, fmt.Errorf("invalid reason for property %q: %w", mutation.GetKey(), err)
		}
		if mutation.GetState() == nil && mutation.Reason == nil && mutation.ClearReason == nil {
			continue
		}
		if mutation.GetClearValue() != nil {
			continue
		}
		// Value first, then state: this mirrors expiry cleanup and protects a
		// state-only patch from racing a scalar clear.
		tValues := table.FivenetUserPropertyValues
		locked := model.FivenetUserPropertyValues{}
		err = tValues.SELECT(tValues.UserID).FROM(tValues).WHERE(mysql.AND(
			tValues.UserID.EQ(mysql.Int32(userID)),
			tValues.PropertyKey.EQ(mysql.String(mutation.GetKey())),
		)).LIMIT(1).FOR(mysql.UPDATE()).QueryContext(ctx, tx, &locked)
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("cannot set state for missing property %q", mutation.GetKey())
		}
		if err != nil {
			return nil, err
		}
		values, err := s.GetPropertiesTx(ctx, tx, userID)
		if err != nil {
			return nil, err
		}
		{
			found := false
			for _, value := range values {
				if value.GetKey() == mutation.GetKey() {
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf(
					"cannot set state for missing property %q",
					mutation.GetKey(),
				)
			}
		}
		tState := table.FivenetUserPropertyState
		current := model.FivenetUserPropertyState{}
		err = tState.SELECT(tState.AllColumns).FROM(tState).WHERE(mysql.AND(
			tState.UserID.EQ(mysql.Int32(userID)),
			tState.PropertyKey.EQ(mysql.String(mutation.GetKey())),
		)).LIMIT(1).FOR(mysql.UPDATE()).QueryContext(ctx, tx, &current)
		if err != nil && !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
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
			if _, err := tState.INSERT(tState.UserID, tState.PropertyKey, tState.Reason, tState.ActiveAt, tState.ExpiresAt).
				VALUES(userID, mutation.GetKey(), stateReason, activeAt, expiresAt).
				ON_DUPLICATE_KEY_UPDATE(
					tState.Reason.SET(mysql.StringExp(mysql.RawString("VALUES(`reason`)"))),
					tState.ActiveAt.SET(mysql.TimestampExp(mysql.RawString("VALUES(`active_at`)"))),
					tState.ExpiresAt.SET(mysql.TimestampExp(mysql.RawString("VALUES(`expires_at`)"))),
				).
				ExecContext(ctx, tx); err != nil {
				return nil, err
			}
		}
		{
			activityMutation := mutation
			if mutation.Reason != nil && current.Reason != nil &&
				*mutation.Reason == *current.Reason {
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
			lifecycle := userMutationLifecycle(activityMutation)
			if lifecycle == nil {
				continue
			}
			for _, value := range values {
				if value.GetKey() == mutation.GetKey() {
					activities = append(
						activities,
						UserPropertyActivities(
							userID,
							sourceUserID,
							reason,
							values,
							[]PropertyChange{
								{
									Property:  value,
									Lifecycle: lifecycle,
									Cause:     usersactivity.PropertyChangeCause_PROPERTY_CHANGE_CAUSE_MANUAL,
								},
							},
						)...)
					break
				}
			}
		}
	}
	return activities, nil
}

// ApplyWantedMutationState persists the lifecycle reason for the managed
// wanted projection after its value/state row has been created in the same
// transaction. Shared public and Sync callers use this bridge.
func (s *Store) ApplyWantedMutationState(
	ctx context.Context,
	tx qrm.DB,
	userID int32,
	mutations []*properties.PropertyMutation,
) error {
	for _, mutation := range mutations {
		if mutation.GetKey() != "user.wanted" ||
			(mutation.Reason == nil && !mutation.GetClearReason()) {
			continue
		}
		tState := table.FivenetUserPropertyState
		reason := mysql.StringExp(mysql.RawString("NULL"))
		if mutation.Reason != nil {
			reason = mysql.String(*mutation.Reason)
		}
		if _, err := tState.UPDATE(tState.Reason).SET(reason).WHERE(mysql.AND(
			tState.UserID.EQ(mysql.Int32(userID)),
			tState.PropertyKey.EQ(mysql.String("user.wanted")),
		)).ExecContext(ctx, tx); err != nil {
			return fmt.Errorf("failed to update user.wanted state reason: %w", err)
		}
	}
	return nil
}

// UserWantedMutationActivities records lifecycle-only mutations which do not
// change the wanted scalar and therefore are invisible to HandleUserPropsChanges.
func UserWantedMutationActivities(
	previous, current *usersprops.UserProps,
	mutations []*properties.PropertyMutation,
	sourceUserID *int32,
	reason string,
) []*usersactivity.UserActivity {
	if previous == nil || current == nil {
		return nil
	}
	for _, mutation := range mutations {
		if mutation.GetKey() != "user.wanted" ||
			(mutation.GetState() == nil && mutation.Reason == nil && mutation.ClearReason == nil) ||
			mutation.GetSetValue() != nil ||
			mutation.GetClearValue() != nil {
			continue
		}
		return []*usersactivity.UserActivity{
			{
				SourceUserId: sourceUserID,
				TargetUserId: current.GetUserId(),
				Type:         usersactivity.UserActivityType_USER_ACTIVITY_TYPE_WANTED,
				Reason:       reason,
				Data: &usersactivity.UserActivityData{
					Data: &usersactivity.UserActivityData_WantedChange{
						WantedChange: &usersactivity.WantedChange{
							Wanted:     current.GetWanted(),
							WantedTill: current.GetWantedTill(),
						},
					},
				},
			},
		}
	}
	return nil
}

func userMutationLifecycle(
	mutation *properties.PropertyMutation,
) *usersactivity.PropertyLifecycleChange {
	lifecycle := propertyregistry.MutationLifecycleForMutation(mutation)
	if lifecycle == nil {
		return nil
	}
	return &usersactivity.PropertyLifecycleChange{
		ExpiresAt: lifecycle.ExpiresAt, ExpiryCleared: lifecycle.ExpiryCleared,
		StateReason: lifecycle.StateReason, StateReasonCleared: lifecycle.StateReasonCleared,
	}
}

// UserPropertyActivities builds activity entries for scalar property changes.
// It is kept in the citizens store package so persistence and audit semantics
// remain together for both service and store callers.
func UserPropertyActivities(
	targetUserID int32,
	sourceUserID *int32,
	reason string,
	previous []*properties.Property,
	changes []PropertyChange,
) []*usersactivity.UserActivity {
	old := make(map[string]*properties.Property, len(previous))
	for _, property := range previous {
		old[property.GetKey()] = property
	}
	activities := make([]*usersactivity.UserActivity, 0, len(changes))
	for _, change := range changes {
		before := old[change.Property.GetKey()]
		if change.Delete && before == nil {
			continue
		}
		if change.Lifecycle == nil &&
			change.Cause == usersactivity.PropertyChangeCause_PROPERTY_CHANGE_CAUSE_UNSPECIFIED &&
			!change.Delete &&
			before != nil &&
			proto.Equal(before.GetValue(), change.Property.GetValue()) {
			continue
		}

		activity := &usersactivity.UserActivity{
			SourceUserId: sourceUserID,
			TargetUserId: targetUserID,
			Type:         usersactivity.UserActivityType_USER_ACTIVITY_TYPE_PROPERTY,
			Key:          change.Property.GetKey(),
			Reason:       reason,
			Data: &usersactivity.UserActivityData{
				Data: &usersactivity.UserActivityData_PropertyChange{
					PropertyChange: &usersactivity.PropertyChange{
						Key:       change.Property.GetKey(),
						Lifecycle: change.Lifecycle,
						Cause:     change.Cause,
						OldValue: func() *properties.PropertyValue {
							if before == nil {
								return nil
							}
							return proto.Clone(before.GetValue()).(*properties.PropertyValue)
						}(),
						NewValue: func() *properties.PropertyValue {
							if change.Delete {
								return nil
							}
							return proto.Clone(change.Property.GetValue()).(*properties.PropertyValue)
						}(),
					},
				},
			},
		}
		activities = append(activities, activity)
	}

	return activities
}

func userPropertyFromModel(stored model.FivenetUserPropertyValues) (*properties.Property, error) {
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

func userPropertyValueColumns(
	property *properties.Property,
) (userPropertyValueType, any, any, any, any, any, any, error) {
	if property == nil || !isUserPropertyKey(property.GetKey()) || property.GetValue() == nil {
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

func isUserPropertyKey(key string) bool {
	return strings.HasPrefix(key, "user.") && key != "user.wanted" && len(key) > len("user.") &&
		!strings.HasSuffix(key, ".") && !strings.HasSuffix(key, "-") && !strings.HasSuffix(key, "_")
}

func userPropertyFilterCondition(
	filter *properties.PropertyFilter,
	userAlias string,
	index int,
) (mysql.BoolExpression, error) {
	if filter.GetKey() == "user.phone_number" {
		return propertyregistry.StringColumnFilterCondition(
			filter,
			"`"+userAlias+"`.phone_number",
		)
	}
	return propertyregistry.FilterCondition(
		filter,
		propertyregistry.FilterConfig{
			Scope: properties.PropertyScope_PROPERTY_SCOPE_USER, ScopeName: "user",
			ValuesTable: userPropertyValuesTable, ParentAlias: userAlias,
			ValueParentColumn: "user_id", ParentColumn: "id",
		},
		index,
	)
}

func (s *Store) IsSystemPropertyEnabled(ctx context.Context, key string) (bool, error) {
	enabled, err := s.propertyCache.SystemPropertyEnabled(
		ctx, s.db, properties.PropertyScope_PROPERTY_SCOPE_USER, key,
	)
	if errors.Is(err, propertyregistry.ErrSystemPropertyNotFound) {
		return false, ErrInvalidProperty
	}
	return enabled, err
}
