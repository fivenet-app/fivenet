package perms

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultPermGuard(t *testing.T) {
	t.Parallel()

	guard, err := DefaultPermGuard("settings.ConfigService", "GetAppConfig")
	require.NoError(t, err)
	assert.Equal(t, BuildGuard("settings", "ConfigService", "GetAppConfig"), guard)
}

func TestDefaultPermGuardRejectsMalformedValues(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		category string
		name     string
	}{
		{category: "", name: "GetAppConfig"},
		{category: "ConfigService", name: "GetAppConfig"},
		{category: "settings.", name: "GetAppConfig"},
		{category: ".ConfigService", name: "GetAppConfig"},
		{category: "settings.ConfigService", name: ""},
	} {
		_, err := DefaultPermGuard(test.category, test.name)
		assert.Error(t, err, "%q/%q should be rejected", test.category, test.name)
	}
}
