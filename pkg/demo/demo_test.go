package demo

import (
	"strings"
	"testing"

	"github.com/fivenet-app/fivenet/v2026/pkg/config"
)

func newTestDemo(seed uint64) *Demo {
	cfg := &config.Config{}
	cfg.Demo.Seed = seed
	cfg.Demo.TargetJob = "police"

	d := &Demo{cfg: cfg}
	d.initRandomizers()

	return d
}

func TestDemoIdentifierDeterministic(t *testing.T) {
	d := newTestDemo(42)

	license := stableLicenseToken("demochar", 7)
	got := d.charIdentifier(1, license)
	want := "char1:c0c9c6cac4c35df2a9277c5f7e8c7b0e492f8e9c"

	if got != want {
		t.Fatalf("expected identifier %q, got %q", want, got)
	}
}

func TestBuildFakeUserProfileDeterministic(t *testing.T) {
	licenses := []string{"dmv", "drive", "weapon"}

	d1 := newTestDemo(42)
	p1 := d1.buildFakeUserProfile(1, "demochar", licenses)

	d2 := newTestDemo(42)
	p2 := d2.buildFakeUserProfile(1, "demochar", licenses)

	if p1.Identifier != p2.Identifier {
		t.Fatalf("identifier mismatch: %q vs %q", p1.Identifier, p2.Identifier)
	}
	if p1.Firstname != p2.Firstname || p1.Lastname != p2.Lastname {
		t.Fatalf(
			"name mismatch: %q %q vs %q %q",
			p1.Firstname,
			p1.Lastname,
			p2.Firstname,
			p2.Lastname,
		)
	}
	if p1.PrimaryJob != p2.PrimaryJob || p1.PrimaryJobGrade != p2.PrimaryJobGrade {
		t.Fatalf(
			"primary job mismatch: %s/%d vs %s/%d",
			p1.PrimaryJob,
			p1.PrimaryJobGrade,
			p2.PrimaryJob,
			p2.PrimaryJobGrade,
		)
	}
	if p1.PhoneNumber != p2.PhoneNumber {
		t.Fatalf("phone mismatch: %q vs %q", p1.PhoneNumber, p2.PhoneNumber)
	}
	if len(p1.Jobs) != len(p2.Jobs) || len(p1.Licenses) != len(p2.Licenses) {
		t.Fatalf(
			"profile sizes mismatch: jobs %d/%d licenses %d/%d",
			len(p1.Jobs),
			len(p2.Jobs),
			len(p1.Licenses),
			len(p2.Licenses),
		)
	}
}

func TestPickUserJobsFromConfiguredPool(t *testing.T) {
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

	for i := 0; i < 250; i++ {
		jobs := d.pickUserJobs()
		if len(jobs) == 0 {
			t.Fatal("expected at least one job")
		}
		if !jobs[0].IsPrimary {
			t.Fatal("expected first job to be primary")
		}

		for _, job := range jobs {
			grades, ok := pool[job.Job]
			if !ok {
				t.Fatalf("job %q not in demo pool", job.Job)
			}
			if _, ok := grades[job.Grade]; !ok {
				t.Fatalf("job grade %d for %q not in demo seed grades", job.Grade, job.Job)
			}
		}
	}
}

func TestBuildTargetJobUserProfileUsesTargetJob(t *testing.T) {
	d := newTestDemo(99)
	d.cfg.Demo.TargetJob = "ambulance"

	profile := d.buildTargetJobUserProfile(3, []string{"drive"})
	if profile.PrimaryJob != "ambulance" {
		t.Fatalf("expected primary job ambulance, got %q", profile.PrimaryJob)
	}
	if len(profile.Jobs) != 1 || profile.Jobs[0].Job != "ambulance" || !profile.Jobs[0].IsPrimary {
		t.Fatalf("expected exactly one primary ambulance job, got %+v", profile.Jobs)
	}
	if !strings.HasPrefix(profile.Identifier, "char1:") {
		t.Fatalf("expected char1 identifier, got %q", profile.Identifier)
	}
}
