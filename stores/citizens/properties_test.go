package citizensstore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreApplyPropertyChangesPersistsAndReturnsTypedValue(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	store := New(Params{DB: db}).(*Store)

	value := &properties.PropertyValue{}
	value.SetStringValue("high-risk")
	property := &properties.Property{Key: "user.risk_level", Value: value}
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO fivenet_user_property_values")).
		WithArgs(int32(42), "user.risk_level", uint8(userPropertyValueTypeString), "high-risk", nil, nil, nil, nil, nil, uint8(userPropertyValueTypeString)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_user_property_values")).
		WithArgs(int32(42), 1, "user").
		WillReturnRows(sqlmock.NewRows([]string{
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
		}).AddRow(int32(42), "user.risk_level", userPropertyValueTypeString, "high-risk", nil, nil, nil, nil, nil, time.Now()))

	persisted, err := store.ApplyPropertyChanges(
		t.Context(),
		db,
		42,
		[]PropertyChange{{Property: property}},
	)
	require.NoError(t, err)
	require.Len(t, persisted, 1)
	assert.Equal(t, "high-risk", persisted[0].GetValue().GetStringValue())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListCustomPropertyKeysExcludesDisabledAndDeleted(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	store := New(Params{DB: db}).(*Store)

	mock.ExpectQuery("SELECT .*property_key.*FROM fivenet_property_definitions.*").
		WithArgs("user").
		WillReturnRows(sqlmock.NewRows([]string{"property_key"}).AddRow("user.risk_level"))

	keys, err := store.ListCustomPropertyKeys(t.Context())
	require.NoError(t, err)
	assert.Equal(t, []string{"user.risk_level"}, keys)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListPropertiesBatchesUsers(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	store := New(Params{DB: db}).(*Store)

	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_user_property_values")).
		WithArgs(int32(42), int32(43), 1, "user").
		WillReturnRows(sqlmock.NewRows([]string{
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
		}).
			AddRow(int32(42), "user.risk_level", userPropertyValueTypeString, "high", nil, nil, nil, nil, nil, time.Now()).
			AddRow(int32(43), "user.risk_level", userPropertyValueTypeString, "low", nil, nil, nil, nil, nil, time.Now()),
		)

	values, err := store.ListProperties(t.Context(), []int32{42, 43})
	require.NoError(t, err)
	require.Len(t, values[42], 1)
	require.Len(t, values[43], 1)
	assert.Equal(t, "high", values[42][0].GetValue().GetStringValue())
	assert.Equal(t, "low", values[43][0].GetValue().GetStringValue())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPropertyFilterConditionEscapesContainsValue(t *testing.T) {
	t.Parallel()
	value := &properties.PropertyValue{}
	value.SetStringValue("50%_off")
	condition, err := userPropertyFilterCondition(&properties.PropertyFilter{
		Key:       "user.note",
		Operation: properties.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_CONTAINS,
		Value:     value,
	}, "user", 0)
	require.NoError(t, err)
	user := table.FivenetUser.AS("user")
	stmt := user.SELECT(user.ID).FROM(user).WHERE(mysql.AND(condition))
	sql, args := stmt.Sql()
	assert.Contains(t, sql, "fivenet_user_property_values AS property_filter_0")
	assert.Contains(t, sql, "LIKE ? ESCAPE '\\\\'")
	assert.Contains(t, args, "%50\\%\\_off%")
}

func TestUserPropertyFilterConditionUsesPhoneColumn(t *testing.T) {
	t.Parallel()
	value := &properties.PropertyValue{}
	value.SetStringValue("50%_off")

	condition, err := userPropertyFilterCondition(&properties.PropertyFilter{
		Key:       "user.phone_number",
		Operation: properties.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_CONTAINS,
		Value:     value,
	}, "user", 0)
	require.NoError(t, err)
	user := table.FivenetUser.AS("user")
	stmt := user.SELECT(user.ID).FROM(user).WHERE(mysql.AND(condition))
	sql, args := stmt.Sql()
	assert.Contains(t, sql, "`user`.phone_number LIKE ? ESCAPE '\\\\'")
	assert.Contains(t, args, "%50\\%\\_off%")
}

func TestListCitizensOptionsForPropertiesRequiresEnabledDefinition(t *testing.T) {
	t.Parallel()

	propertyList := &permissionsattributes.PropertyList{Keys: []string{
		"user.dateofbirth", "user.height", "user.phone_number", "user.visum",
	}}
	definitions := map[string]*properties.PropertyDefinition{
		"user.dateofbirth":  {Enabled: true},
		"user.height":       {Enabled: false},
		"user.phone_number": {Enabled: true},
		"user.visum":        {Enabled: true},
	}

	options := ListCitizensOptionsForProperties(propertyList, definitions)
	assert.True(t, options.IncludeDateOfBirth)
	assert.False(t, options.IncludeHeight)
	assert.True(t, options.IncludePhoneNumber)
	assert.True(t, options.IncludeVisum)
	assert.False(t, options.IncludeSex)
}
