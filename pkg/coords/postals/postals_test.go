package postals

import (
	"os"
	"testing"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	// Create a temporary file to simulate the postals file
	tmpFile, err := os.CreateTemp(t.TempDir(), "postals_test_*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Write sample postal data to the temporary file
	sampleData := `[{"Code":"12345","Latitude":40.7128,"Longitude":-74.0060},{"Code":"67890","Latitude":34.0522,"Longitude":-118.2437}]`
	_, err = tmpFile.WriteString(sampleData)
	require.NoError(t, err)
	tmpFile.Close()

	// Create a mock config pointing to the temporary file
	cfg := &config.Config{
		PostalsFile: tmpFile.Name(),
	}

	// Call the New function
	postals, err := New(cfg)
	require.NoError(t, err)
	assert.NotNil(t, postals)

	// Verify the postalCodesMap is populated correctly
	postal, ok := postalCodesMap["12345"]
	assert.True(t, ok)
	assert.NotNil(t, postal)
	assert.Equal(t, "12345", *postal.Code)

	// Make sure non-existing code doesn't exist
	postal, ok = postalCodesMap["789123"]
	assert.False(t, ok)
	assert.Nil(t, postal)
}

func TestByCode(t *testing.T) {
	code := "67890"
	// Prepare the postalCodesMap for testing
	postalCodesMap["67890"] = &Postal{
		Code: &code,
		X:    34.0522,
		Y:    -118.2437,
	}

	// Test existing code
	postal, ok := ByCode("67890")
	assert.True(t, ok)
	assert.NotNil(t, postal)
	assert.Equal(t, "67890", *postal.Code)

	// Test non-existing code
	postal, ok = ByCode("99999")
	assert.False(t, ok)
	assert.Nil(t, postal)
}
