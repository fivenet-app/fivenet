package properties

import (
	"math"
	"testing"
	"time"

	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/stretchr/testify/require"
)

func TestValueColumnsAndStoredValueRoundTrip(t *testing.T) {
	t.Parallel()

	value := &propertypb.PropertyValue{}
	value.SetDateValue(&propertypb.LocalDate{Year: 2028, Month: 2, Day: 29})

	typeValue, _, _, _, _, _, date, err := ValueColumns(value)
	require.NoError(t, err)
	require.Equal(t, propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_DATE, typeValue)
	require.Equal(t, time.Date(2028, 2, 29, 0, 0, 0, 0, time.UTC), date)

	dateValue := date.(time.Time)
	roundTripped, err := DecodeStoredValue(StoredValue{Type: typeValue, Date: &dateValue})
	require.NoError(t, err)
	require.Equal(t, value, roundTripped)
}

func TestValueColumnsRejectsInvalidValues(t *testing.T) {
	t.Parallel()

	value := &propertypb.PropertyValue{}
	value.SetNumberValue(math.Inf(1))
	_, _, _, _, _, _, _, err := ValueColumns(value)
	require.ErrorIs(t, err, ErrInvalidValue)
}
