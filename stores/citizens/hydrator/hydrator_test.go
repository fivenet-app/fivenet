package citizenshydrator

import (
	"testing"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
)

type testUserAwareEnricher struct {
	mstlystcdata.DummyEnricher

	calls int
}

func (e *testUserAwareEnricher) EnrichJobInfoSafe(
	userInfo *userinfo.UserInfo,
	usrs ...common.IJobInfo,
) {
	for _, usr := range usrs {
		e.calls++
		usr.SetJob("unemployed")
		usr.SetJobGrade(0)
		usr.SetJobLabel("N/A")
		usr.SetJobGradeLabel("N/A")
	}
}

func (e *testUserAwareEnricher) EnrichJobInfoSafeFunc(
	userInfo *userinfo.UserInfo,
) func(usr common.IJobInfo) {
	return func(usr common.IJobInfo) {
		e.calls++
		usr.SetJob("unemployed")
		usr.SetJobGrade(0)
		usr.SetJobLabel("N/A")
		usr.SetJobGradeLabel("N/A")
	}
}

func TestEnrichShortsSafeRedactsJobInfo(t *testing.T) {
	t.Parallel()

	enricher := &testUserAwareEnricher{}
	h := &Hydrator{
		enricher: enricher,
	}

	user := &usershort.UserShort{
		UserId:   42,
		Job:      "police",
		JobGrade: 4,
	}

	h.enrichShortsSafe(nil, user, nil)

	if enricher.calls != 1 {
		t.Fatalf("expected 1 enrichment call, got %d", enricher.calls)
	}
	if got := user.GetJob(); got != "unemployed" {
		t.Fatalf("expected job to be redacted, got %q", got)
	}
	if got := user.GetJobGrade(); got != 0 {
		t.Fatalf("expected job grade to be redacted, got %d", got)
	}
	if got := user.GetJobLabel(); got != "N/A" {
		t.Fatalf("expected job label to be redacted, got %q", got)
	}
	if got := user.GetJobGradeLabel(); got != "N/A" {
		t.Fatalf("expected job grade label to be redacted, got %q", got)
	}
}
