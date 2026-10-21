package citizensstore

import (
	"context"
	"errors"
	"fmt"
	"strings"

	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
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
	Property *properties.Property
	Delete   bool
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

	return result, nil
}

func (s *Store) GetPropertiesTx(
	ctx context.Context,
	tx qrm.DB,
	userID int32,
) ([]*properties.Property, error) {
	return s.getProperties(ctx, tx, userID)
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
	return values, nil
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
	if len(changes) == 0 {
		return nil, nil
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
		if !change.Delete && before != nil &&
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
						Key: change.Property.GetKey(),
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
