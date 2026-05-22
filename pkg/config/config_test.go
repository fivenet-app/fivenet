package config

import (
	"testing"

	"github.com/creasty/defaults"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDemoDefaults(t *testing.T) {
	t.Parallel()
	cfg := &Config{}
	require.NoError(t, defaults.Set(cfg), "failed to set defaults")

	assert.True(t, cfg.Demo.Features.Dispatches, "expected demo.features.dispatches default to be true")
	assert.True(t, cfg.Demo.Features.Locations, "expected demo.features.locations default to be true")
	assert.True(t, cfg.Demo.Features.Timeclock, "expected demo.features.timeclock default to be true")
	assert.False(t, cfg.Demo.Features.Users, "expected demo.features.users default to be false")

	assert.EqualValues(t, 50, cfg.Demo.FakeUsers.Count, "expected demo.fakeUsers.count default 50, got %d", cfg.Demo.FakeUsers.Count)
}
