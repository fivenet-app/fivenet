package properties

import (
	"testing"

	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalDateToTime(t *testing.T) {
	t.Parallel()

	value, err := LocalDateToTime(&propertypb.LocalDate{Year: 2028, Month: 2, Day: 29})
	require.NoError(t, err)
	assert.Equal(t, 2028, value.Year())
	assert.Equal(t, 2, int(value.Month()))
	assert.Equal(t, 29, value.Day())

	_, err = LocalDateToTime(&propertypb.LocalDate{Year: 2027, Month: 2, Day: 29})
	require.ErrorIs(t, err, ErrInvalidLocalDate)
}
