package userinfo

import (
	"fmt"
	"testing"
	"time"

	accounts "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	pbtimestamp "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type labelEnricher struct{}

func (labelEnricher) EnrichJobInfo(user common.IJobInfo) {
	user.SetJobLabel(user.GetJob())
	user.SetJobGradeLabel(fmt.Sprintf("Rank %d", user.GetJobGrade()))
}

func (labelEnricher) EnrichJobInfoNoFallback(common.IJobInfo) {}

func (labelEnricher) EnrichJobName(common.IJobName) {}

func (labelEnricher) GetJobByName(string) *jobs.Job { return nil }

func (labelEnricher) GetJobGrade(string, int32) (*jobs.Job, *jobs.JobGrade) {
	return nil, nil
}

func TestBuildUserInfoChangedEvent(t *testing.T) {
	t.Parallel()

	evt := BuildUserInfoChangedEvent(
		42,
		77,
		nil,
		"police",
		3,
		labelEnricher{},
	)

	require.NotNil(t, evt)
	require.NotNil(t, evt.GetChangedAt())

	assert.Equal(t, int64(42), evt.GetAccountId())
	assert.Equal(t, int32(77), evt.GetUserId())
	assert.Same(t, evt.GetChangedAt(), evt.ChangedAt)
	assert.Equal(t, "police", evt.GetNewJob())
	assert.Equal(t, "police", evt.GetNewJobLabel())
	assert.Equal(t, int32(3), evt.GetNewJobGrade())
	assert.Equal(t, "Rank 3", evt.GetNewJobGradeLabel())
}

func TestBuildUserInfoChangedEventUsesProvidedTimestamp(t *testing.T) {
	t.Parallel()

	changedAt := pbtimestamp.New(time.Unix(123, 456))
	evt := BuildUserInfoChangedEvent(1, 2, changedAt, "ems", 4, nil)

	require.NotNil(t, evt)
	assert.Same(t, changedAt, evt.GetChangedAt())
	assert.Equal(t, "ems", evt.GetNewJob())
	assert.Equal(t, int32(4), evt.GetNewJobGrade())
}

func TestBuildAccountGroupsChangedEvent(t *testing.T) {
	t.Parallel()

	groups := &accounts.AccountGroups{Groups: []string{"supporter", "donator"}}
	evt := BuildAccountGroupsChangedEvent(42, nil, groups, true, false)

	require.NotNil(t, evt)
	require.NotNil(t, evt.GetChangedAt())
	require.NotNil(t, evt.GetNewGroups())

	assert.Equal(t, int64(42), evt.GetAccountId())
	assert.Equal(t, []string{"supporter", "donator"}, evt.GetNewGroups().GetGroups())
	assert.True(t, evt.GetCanBeSuperuser())
	assert.False(t, evt.GetCanBeConfigAdmin())

	groups.Groups[0] = "modified"
	assert.Equal(t, []string{"supporter", "donator"}, evt.GetNewGroups().GetGroups())
}
