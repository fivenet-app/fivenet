package properties

import (
	"testing"
	"time"

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

func TestLocalDateFromTime(t *testing.T) {
	t.Parallel()

	assert.Equal(t, &propertypb.LocalDate{Year: 2028, Month: 2, Day: 29},
		LocalDateFromTime(time.Date(2028, time.February, 29, 12, 0, 0, 0, time.FixedZone("UTC+2", 2*60*60))))
	assert.Nil(t, LocalDateFromTime(time.Date(10000, time.January, 1, 0, 0, 0, 0, time.UTC)))
}
