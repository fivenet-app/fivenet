package properties

import (
	"errors"
	"math"
	"slices"
	"time"
	"unicode/utf8"

	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
)

var ErrInvalidValue = errors.New("invalid property value")

// StoredValue is the database-neutral representation of a scalar property
// value. Exactly one value field must be populated according to Type.
type StoredValue struct {
	Type      propertypb.PropertyValueType
	String    *string
	Integer   *int64
	Number    *float64
	Bool      *bool
	Timestamp *time.Time
	Date      *time.Time
}

// ValidateValueConstraints checks a scalar value against its definition after
// ValueColumns has established its type and basic representation validity.
func ValidateValueConstraints(
	definition *propertypb.PropertyDefinition,
	value *propertypb.PropertyValue,
) error {
	if definition == nil || value == nil {
		return ErrInvalidValue
	}
	valueType, stringValue, integerValue, numberValue, _, timestampValue, dateValue, err := ValueColumns(
		value,
	)
	if err != nil || valueType != definition.GetValueType() {
		return ErrInvalidValue
	}
	constraints := definition.GetValueConstraints()
	if constraints == nil {
		return nil
	}
	// Definitions normally pass this check when they are persisted. Rechecking
	// here also protects write paths from manually corrupted legacy rows.
	if !IsValueConstraintsValid(definition.GetValueType(), constraints) {
		return ErrInvalidValue
	}
	switch valueType {
	case propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_STRING:
		text := stringValue.(string)
		length, ok := utils.ToUint32Checked(utf8.RuneCountInString(text))
		if !ok || constraints.StringMinLength != nil && length < constraints.GetStringMinLength() ||
			constraints.StringMaxLength != nil && length > constraints.GetStringMaxLength() ||
			len(constraints.GetStringAllowedValues()) > 0 &&
				constraints.GetStringAllowedValuesMode() != propertypb.PropertyStringAllowedValuesMode_PROPERTY_STRING_ALLOWED_VALUES_MODE_SUGGEST &&
				!slices.Contains(constraints.GetStringAllowedValues(), text) {
			return ErrInvalidValue
		}
	case propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_INTEGER:
		number := integerValue.(int64)
		if constraints.IntegerMin != nil && number < constraints.GetIntegerMin() ||
			constraints.IntegerMax != nil && number > constraints.GetIntegerMax() {
			return ErrInvalidValue
		}
	case propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_NUMBER:
		number := numberValue.(float64)
		if constraints.NumberMin != nil && number < constraints.GetNumberMin() ||
			constraints.NumberMax != nil && number > constraints.GetNumberMax() {
			return ErrInvalidValue
		}
	case propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_TIMESTAMP:
		timestamp := timestampValue.(time.Time)
		if constraints.TimestampMin != nil &&
			timestamp.Before(constraints.GetTimestampMin().AsTime()) ||
			constraints.TimestampMax != nil &&
				timestamp.After(constraints.GetTimestampMax().AsTime()) {
			return ErrInvalidValue
		}
	case propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_DATE:
		date := dateValue.(time.Time)
		if constraints.DateMin != nil {
			minDate, err := LocalDateToTime(constraints.GetDateMin())
			if err != nil || date.Before(minDate) {
				return ErrInvalidValue
			}
		}
		if constraints.DateMax != nil {
			maxDate, err := LocalDateToTime(constraints.GetDateMax())
			if err != nil || date.After(maxDate) {
				return ErrInvalidValue
			}
		}
	}
	return nil
}

// ValueColumns returns the values used by the scalar property table columns.
// All columns other than the one selected by Type are nil.
func ValueColumns(
	value *propertypb.PropertyValue,
) (propertypb.PropertyValueType, any, any, any, any, any, any, error) {
	if value == nil {
		return 0, nil, nil, nil, nil, nil, nil, ErrInvalidValue
	}

	switch value := value.GetValue().(type) {
	case *propertypb.PropertyValue_StringValue:
		return propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_STRING, value.StringValue, nil, nil, nil, nil, nil, nil
	case *propertypb.PropertyValue_IntegerValue:
		return propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_INTEGER, nil, value.IntegerValue, nil, nil, nil, nil, nil
	case *propertypb.PropertyValue_NumberValue:
		if math.IsNaN(value.NumberValue) || math.IsInf(value.NumberValue, 0) {
			return 0, nil, nil, nil, nil, nil, nil, ErrInvalidValue
		}
		return propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_NUMBER, nil, nil, value.NumberValue, nil, nil, nil, nil
	case *propertypb.PropertyValue_BoolValue:
		return propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_BOOL, nil, nil, nil, value.BoolValue, nil, nil, nil
	case *propertypb.PropertyValue_TimestampValue:
		if value.TimestampValue == nil || value.TimestampValue.GetTimestamp() == nil {
			return 0, nil, nil, nil, nil, nil, nil, ErrInvalidValue
		}
		return propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_TIMESTAMP, nil, nil, nil, nil, value.TimestampValue.AsTime(), nil, nil
	case *propertypb.PropertyValue_DateValue:
		date, err := LocalDateToTime(value.DateValue)
		if err != nil {
			return 0, nil, nil, nil, nil, nil, nil, ErrInvalidValue
		}
		return propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_DATE, nil, nil, nil, nil, nil, date, nil
	default:
		return 0, nil, nil, nil, nil, nil, nil, ErrInvalidValue
	}
}

func DecodeStoredValue(value StoredValue) (*propertypb.PropertyValue, error) {
	result := &propertypb.PropertyValue{}
	switch value.Type {
	case propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_STRING:
		if value.String == nil {
			return nil, ErrInvalidValue
		}
		result.SetStringValue(*value.String)
	case propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_INTEGER:
		if value.Integer == nil {
			return nil, ErrInvalidValue
		}
		result.SetIntegerValue(*value.Integer)
	case propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_NUMBER:
		if value.Number == nil || math.IsNaN(*value.Number) || math.IsInf(*value.Number, 0) {
			return nil, ErrInvalidValue
		}
		result.SetNumberValue(*value.Number)
	case propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_BOOL:
		if value.Bool == nil {
			return nil, ErrInvalidValue
		}
		result.SetBoolValue(*value.Bool)
	case propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_TIMESTAMP:
		if value.Timestamp == nil {
			return nil, ErrInvalidValue
		}
		result.SetTimestampValue(timestamp.New(*value.Timestamp))
	case propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_DATE:
		if value.Date == nil {
			return nil, ErrInvalidValue
		}
		result.SetDateValue(LocalDateFromTime(*value.Date))
	default:
		return nil, ErrInvalidValue
	}
	return result, nil
}
