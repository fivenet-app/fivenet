package properties

import (
	"errors"
	"math"
	"time"

	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
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
