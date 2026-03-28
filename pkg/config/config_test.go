package config

import (
	"testing"

	"github.com/creasty/defaults"
)

func TestDemoDefaults(t *testing.T) {
	t.Parallel()
	cfg := &Config{}
	if err := defaults.Set(cfg); err != nil {
		t.Fatalf("failed to set defaults: %v", err)
	}

	if cfg.Demo.Features.Dispatches != true {
		t.Fatal("expected demo.features.dispatches default to be true")
	}
	if cfg.Demo.Features.Locations != true {
		t.Fatal("expected demo.features.locations default to be true")
	}
	if cfg.Demo.Features.Timeclock != true {
		t.Fatal("expected demo.features.timeclock default to be true")
	}
	if cfg.Demo.Features.Users != false {
		t.Fatal("expected demo.features.users default to be false")
	}

	if cfg.Demo.FakeUsers.Count != 50 {
		t.Fatalf("expected demo.fakeUsers.count default 50, got %d", cfg.Demo.FakeUsers.Count)
	}
}
