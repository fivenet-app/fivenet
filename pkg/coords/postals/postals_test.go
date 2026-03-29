package postals

import (
	"os"
	"testing"

	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/quadtree"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLoadsAndQueriesPostalCodes(t *testing.T) {
	t.Parallel()

	// Arrange a temporary postals file with a couple of known entries.
	tmpFile, err := os.CreateTemp(t.TempDir(), "postals_test_*.json")
	require.NoError(t, err)

	sampleData := `[
		{"code":"12345","x":40.7128,"y":-74.0060},
		{"code":"67890","x":34.0522,"y":-118.2437}
	]`
	_, err = tmpFile.WriteString(sampleData)
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	cfg := &config.Config{
		PostalsFile: tmpFile.Name(),
	}

	// Load the postals store from disk through the public constructor.
	postals, err := New(cfg)
	require.NoError(t, err)
	require.NotNil(t, postals)

	byCode := func(code string) quadtree.FilterFunc {
		return func(point orb.Pointer) bool {
			postal, ok := point.(*Postal)
			return ok && postal.Code != nil && *postal.Code == code
		}
	}

	target := orb.Point{40.7128, -74.0060}

	// Query the loaded store through the embedded coords API.
	require.True(t, postals.Has(target, byCode("12345")))

	postal := postals.Get(target, byCode("12345"))
	require.NotNil(t, postal)
	require.NotNil(t, postal.Code)
	assert.Equal(t, "12345", *postal.Code)
	assert.Equal(t, target, postal.Point())

	postal, ok := postals.ByCode("12345")
	require.True(t, ok)
	require.NotNil(t, postal)
	require.NotNil(t, postal.Code)
	assert.Equal(t, "12345", *postal.Code)

	// Verify missing postal codes are not reported as present.
	missing := orb.Point{999, 999}
	assert.False(t, postals.Has(missing, byCode("99999")))

	postal, ok = postals.ByCode("99999")
	assert.False(t, ok)
	assert.Nil(t, postal)

	// Confirm spatial nearest-neighbor lookup still works on the same store.
	closest, ok := postals.Closest(34.0522, -118.2437)
	require.True(t, ok)
	require.NotNil(t, closest)
	require.NotNil(t, closest.Code)
	assert.Equal(t, "67890", *closest.Code)
	assert.Equal(t, orb.Point{34.0522, -118.2437}, closest.Point())
}
