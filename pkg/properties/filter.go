package properties

import (
	"errors"
	"fmt"
	"math"
	"strings"

	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/go-jet/jet/v2/mysql"
)

var ErrInvalidFilter = errors.New("invalid property filter")

type FilterConfig struct {
	Scope             propertypb.PropertyScope
	ScopeName         string
	ValuesTable       string
	ParentAlias       string
	ValueParentColumn string
	ParentColumn      string
}

// StringColumnFilterCondition compiles a logical property filter against a
// dedicated string column. The column expression is supplied by trusted store
// code; filter values remain parameterized.
func StringColumnFilterCondition(
	filter *propertypb.PropertyFilter,
	column string,
) (mysql.BoolExpression, error) {
	if filter == nil || filter.GetValue() == nil {
		return nil, ErrInvalidFilter
	}
	if _, ok := filter.GetValue().GetValue().(*propertypb.PropertyValue_StringValue); !ok {
		return nil, ErrInvalidFilter
	}

	var operator string
	value := any(filter.GetValue().GetStringValue())
	var escape string
	switch filter.GetOperation() {
	case propertypb.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_EQUALS:
		operator = "="
	case propertypb.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_NOT_EQUALS:
		operator = "<>"
	case propertypb.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_GREATER_THAN:
		operator = ">"
	case propertypb.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_GREATER_THAN_OR_EQUAL:
		operator = ">="
	case propertypb.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_LESS_THAN:
		operator = "<"
	case propertypb.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_LESS_THAN_OR_EQUAL:
		operator = "<="
	case propertypb.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_CONTAINS:
		operator = "LIKE"
		escape = " ESCAPE '\\\\'"
		value = "%" + strings.ReplaceAll(
			strings.ReplaceAll(filter.GetValue().GetStringValue(), "%", "\\%"),
			"_",
			"\\_",
		) + "%"
	default:
		return nil, ErrInvalidFilter
	}

	return mysql.RawBool(
		column+" "+operator+" $value"+escape,
		mysql.RawArgs{"$value": value},
	), nil
}

// FilterCondition compiles a definition-aware correlated EXISTS condition for
// a scalar property value table. System definitions come from the registry;
// custom definitions are checked in the database. System overrides only need
// a disabled-row check because generated system definitions default enabled.
func FilterCondition(
	filter *propertypb.PropertyFilter,
	config FilterConfig,
	index int,
) (mysql.BoolExpression, error) {
	valueType, column, value, err := filterValue(filter.GetValue())
	if err != nil ||
		filter.GetOperation() == propertypb.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_UNSPECIFIED {
		return nil, ErrInvalidFilter
	}
	var operator string
	var escape string
	switch filter.GetOperation() {
	case propertypb.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_EQUALS:
		operator = "="
	case propertypb.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_NOT_EQUALS:
		operator = "<>"
	case propertypb.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_GREATER_THAN:
		operator = ">"
	case propertypb.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_GREATER_THAN_OR_EQUAL:
		operator = ">="
	case propertypb.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_LESS_THAN:
		operator = "<"
	case propertypb.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_LESS_THAN_OR_EQUAL:
		operator = "<="
	case propertypb.PropertyFilterOperation_PROPERTY_FILTER_OPERATION_CONTAINS:
		if valueType != propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_STRING {
			return nil, ErrInvalidFilter
		}
		operator, escape = "LIKE", " ESCAPE '\\\\'"
		value = "%" + strings.ReplaceAll(
			strings.ReplaceAll(fmt.Sprint(value), "%", "\\%"),
			"_",
			"\\_",
		) + "%"
	default:
		return nil, ErrInvalidFilter
	}

	definitionCheck := "EXISTS (SELECT 1 FROM fivenet_property_definitions AS property_definition WHERE property_definition.scope = $scope AND property_definition.property_key = $key AND property_definition.deleted_at IS NULL AND property_definition.enabled IS TRUE AND property_definition.searchable IS TRUE)"
	if definition, ok := SystemPropertyDefinition(config.Scope, filter.GetKey()); ok {
		if !definition.GetSearchable() || definition.GetValueType() != valueType {
			return nil, ErrInvalidFilter
		}
		definitionCheck = "NOT EXISTS (SELECT 1 FROM fivenet_property_system_overrides AS property_override WHERE property_override.scope = $scope AND property_override.property_key = $key AND property_override.enabled IS FALSE)"
	}
	alias := fmt.Sprintf("property_filter_%d", index)
	parentColumn := config.ParentColumn
	if parentColumn == "" {
		parentColumn = "id"
	}
	return mysql.RawBool(
		"EXISTS (SELECT 1 FROM "+config.ValuesTable+" AS "+alias+" WHERE "+alias+"."+config.ValueParentColumn+" = `"+config.ParentAlias+"`."+parentColumn+" AND "+alias+".property_key = $key AND "+alias+".value_type = $type AND "+alias+"."+column+" "+operator+" $value"+escape+" AND "+definitionCheck+")",
		mysql.RawArgs{
			"$key":   filter.GetKey(),
			"$scope": config.ScopeName,
			"$type":  int32(valueType),
			"$value": value,
		},
	), nil
}

func filterValue(
	value *propertypb.PropertyValue,
) (propertypb.PropertyValueType, string, any, error) {
	if value == nil {
		return 0, "", nil, ErrInvalidFilter
	}
	switch v := value.GetValue().(type) {
	case *propertypb.PropertyValue_StringValue:
		return propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_STRING, "value_string", v.StringValue, nil
	case *propertypb.PropertyValue_IntegerValue:
		return propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_INTEGER, "value_integer", v.IntegerValue, nil
	case *propertypb.PropertyValue_NumberValue:
		if math.IsNaN(v.NumberValue) || math.IsInf(v.NumberValue, 0) {
			return 0, "", nil, ErrInvalidFilter
		}
		return propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_NUMBER, "value_number", v.NumberValue, nil
	case *propertypb.PropertyValue_BoolValue:
		return propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_BOOL, "value_bool", v.BoolValue, nil
	case *propertypb.PropertyValue_TimestampValue:
		if v.TimestampValue == nil || v.TimestampValue.GetTimestamp() == nil {
			return 0, "", nil, ErrInvalidFilter
		}
		return propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_TIMESTAMP, "value_timestamp", v.TimestampValue.AsTime(), nil
	case *propertypb.PropertyValue_DateValue:
		date, err := LocalDateToTime(v.DateValue)
		if err != nil {
			return 0, "", nil, ErrInvalidFilter
		}
		return propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_DATE, "value_date", date, nil
	default:
		return 0, "", nil, ErrInvalidFilter
	}
}
