package demo

import (
	"strings"
	"testing"

	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestDemo(seed uint64) *Demo {
	cfg := &config.Config{}
	cfg.Demo.Seed = seed
	cfg.Demo.TargetJob = PoliceJob

	d := &Demo{cfg: cfg}
	d.initRandomizers()

	return d
}

func TestDemoIdentifierDeterministic(t *testing.T) {
	t.Parallel()
	d := newTestDemo(42)

	license := stableLicenseToken("demochar", 7)
	got := d.charIdentifier(1, license)
	want := "char1:860d10f4cb5bb61a609ca73e14b0255dc366f2c2e9494bcfb7a34b9da9"

	assert.Equal(t, want, got, "expected identifier %q, got %q", want, got)
}

func TestBuildFakeUserProfileDeterministic(t *testing.T) {
	t.Parallel()
	licenses := []string{"dmv", "drive", "weapon"}

	d1 := newTestDemo(42)
	p1 := d1.buildFakeUserProfile(1, "demochar", licenses)

	d2 := newTestDemo(42)
	p2 := d2.buildFakeUserProfile(1, "demochar", licenses)

	assert.Equal(
		t,
		p2.Identifier,
		p1.Identifier,
		"identifier mismatch: %q vs %q",
		p1.Identifier,
		p2.Identifier,
	)
	assert.Equal(
		t,
		p2.Firstname,
		p1.Firstname,
		"name mismatch: %q %q vs %q %q",
		p1.Firstname,
		p1.Lastname,
		p2.Firstname,
		p2.Lastname,
	)
	assert.Equal(
		t,
		p2.Lastname,
		p1.Lastname,
		"name mismatch: %q %q vs %q %q",
		p1.Firstname,
		p1.Lastname,
		p2.Firstname,
		p2.Lastname,
	)
	assert.Equal(
		t,
		p2.PrimaryJob,
		p1.PrimaryJob,
		"primary job mismatch: %s/%d vs %s/%d",
		p1.PrimaryJob,
		p1.PrimaryJobGrade,
		p2.PrimaryJob,
		p2.PrimaryJobGrade,
	)
	assert.Equal(
		t,
		p2.PrimaryJobGrade,
		p1.PrimaryJobGrade,
		"primary job mismatch: %s/%d vs %s/%d",
		p1.PrimaryJob,
		p1.PrimaryJobGrade,
		p2.PrimaryJob,
		p2.PrimaryJobGrade,
	)
	assert.Equal(
		t,
		p2.PhoneNumber,
		p1.PhoneNumber,
		"phone mismatch: %q vs %q",
		p1.PhoneNumber,
		p2.PhoneNumber,
	)
	assert.Len(
		t,
		p1.Jobs, len(p2.Jobs),
		"profile sizes mismatch: jobs %d/%d licenses %d/%d",
		len(p1.Jobs),
		len(p2.Jobs),
		len(p1.Licenses),
		len(p2.Licenses),
	)
	assert.Len(
		t,
		p1.Licenses, len(p2.Licenses),
		"profile sizes mismatch: jobs %d/%d licenses %d/%d",
		len(p1.Jobs),
		len(p2.Jobs),
		len(p1.Licenses),
		len(p2.Licenses),
	)
}

func TestPickUserJobsFromConfiguredPool(t *testing.T) {
	t.Parallel()
	d := newTestDemo(1337)

	pool := map[string]map[int32]struct{}{}
	for _, job := range demoSeedJobs {
		pool[job.Name] = map[int32]struct{}{}
	}
	for _, grade := range demoSeedJobGrades {
		if _, ok := pool[grade.JobName]; ok {
			pool[grade.JobName][grade.Grade] = struct{}{}
		}
	}

	for range 250 {
		jobs := d.pickUserJobs()
		require.NotEmpty(t, jobs, "expected at least one job")
		assert.True(t, jobs[0].IsPrimary, "expected first job to be primary")

		for _, job := range jobs {
			grades, ok := pool[job.Job]
			require.True(t, ok, "job %q not in demo pool", job.Job)
			_, ok = grades[job.Grade]
			assert.True(t, ok, "job grade %d for %q not in demo seed grades", job.Grade, job.Job)
		}
	}
}

func TestBuildTargetJobUserProfileUsesTargetJob(t *testing.T) {
	t.Parallel()
	d := newTestDemo(99)
	d.cfg.Demo.TargetJob = "ambulance"

	profile := d.buildTargetJobUserProfile(3, []string{"drive"})
	assert.Equal(
		t,
		"ambulance",
		profile.PrimaryJob,
		"expected primary job ambulance, got %q",
		profile.PrimaryJob,
	)
	require.Len(
		t,
		profile.Jobs,
		1,
		"expected exactly one primary ambulance job, got %+v",
		profile.Jobs,
	)
	assert.Equal(
		t,
		"ambulance",
		profile.Jobs[0].Job,
		"expected exactly one primary ambulance job, got %+v",
		profile.Jobs,
	)
	assert.True(
		t,
		profile.Jobs[0].IsPrimary,
		"expected exactly one primary ambulance job, got %+v",
		profile.Jobs,
	)
	assert.True(
		t,
		strings.HasPrefix(profile.Identifier, "char1:"),
		"expected char1 identifier, got %q",
		profile.Identifier,
	)
}
