package jobs

import (
	"testing"

	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
	groupsshort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups/short"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	permsjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs/perms"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/stretchr/testify/require"
)

func TestAttachColleagueGroupsPreservesExistingProps(t *testing.T) {
	t.Parallel()

	labels := &jobscolleagues.ColleagueProps{UserId: 7, Job: "police"}
	colleagues := []*jobscolleagues.Colleague{{UserId: 7, Props: labels}}
	groups := map[int32][]*groupsshort.GroupMemberShort{
		7: {{Id: 42, Job: "police", Name: "Traffic"}},
	}

	attachColleagueGroups(colleagues, groups)

	if colleagues[0].GetProps() != labels {
		t.Fatal("expected existing colleague props to be preserved")
	}
	if len(colleagues[0].GetProps().GetGroups()) != 1 {
		t.Fatalf("expected one group, got %d", len(colleagues[0].GetProps().GetGroups()))
	}
}

func TestGetConditionForColleagueAccessUsesJobGradeForRankAccess(t *testing.T) {
	t.Parallel()

	s := &Server{}
	tActivity := tColleagueActivity.AS("colleague_activity")
	tTargetUserJobs := table.FivenetUserJobs.AS("target_user_jobs")

	condition := s.getConditionForColleagueAccess(
		tActivity,
		tTargetUserJobs,
		[]permsjobs.ColleaguesServiceGetColleagueAccessPermValue{
			permsjobs.ColleaguesServiceGetColleagueAccessPermValueLowerRank,
		},
		&userinfo.UserInfo{
			UserId:   12,
			Job:      "police",
			JobGrade: 5,
		},
	)

	stmt := tActivity.
		SELECT(tActivity.ID).
		FROM(tActivity).
		WHERE(condition)
	sql, args := stmt.Sql()

	require.Contains(t, sql, "target_user_jobs.grade < ?")
	require.NotContains(t, sql, "target_user.id < ?")
	require.Equal(t, []any{int32(5)}, args)
}

func TestGetConditionForColleagueAccessUsesTargetUserIDForOwnAccess(t *testing.T) {
	t.Parallel()

	s := &Server{}
	tActivity := tColleagueActivity.AS("colleague_activity")
	tTargetUserJobs := table.FivenetUserJobs.AS("target_user_jobs")

	condition := s.getConditionForColleagueAccess(
		tActivity,
		tTargetUserJobs,
		[]permsjobs.ColleaguesServiceGetColleagueAccessPermValue{
			permsjobs.ColleaguesServiceGetColleagueAccessPermValueOwn,
		},
		&userinfo.UserInfo{
			UserId:   12,
			Job:      "police",
			JobGrade: 5,
		},
	)

	stmt := tActivity.
		SELECT(tActivity.ID).
		FROM(tActivity).
		WHERE(condition)
	sql, args := stmt.Sql()

	require.Contains(t, sql, "colleague_activity.target_user_id = ?")
	require.Equal(t, []any{int32(12)}, args)
}
