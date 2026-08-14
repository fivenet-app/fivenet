package jobs

import (
	"context"
	"errors"
	"testing"

	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	qualificationsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/access"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	errorsjobs "github.com/fivenet-app/fivenet/v2026/services/jobs/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type stubQualificationAccess struct {
	allowedIDs []int64
	calls      int
	access     int32
	targetIDs  []int64
	err        error
}

func (s *stubQualificationAccess) CanUserAccessTargetIDs(
	_ context.Context,
	_ *pbuserinfo.UserInfo,
	access int32,
	targetIDs ...int64,
) ([]int64, error) {
	s.calls++
	s.access = access
	s.targetIDs = append([]int64(nil), targetIDs...)

	return s.allowedIDs, s.err
}

func TestEnrichGroupRuleGradeLabels(t *testing.T) {
	t.Parallel()

	server := &Server{enricher: mstlystcdata.NewDummyUserAwareEnricher()}
	minimumRule := &jobsgroups.GroupRule{}
	minimumRule.SetGrade(&jobsgroups.GroupGradeRule{
		Type:  jobsgroups.GroupGradeRuleType_GROUP_GRADE_RULE_TYPE_MINIMUM,
		Grade: new(int32(3)),
	})
	rangeRule := &jobsgroups.GroupRule{}
	rangeRule.SetGrade(&jobsgroups.GroupGradeRule{
		Type:     jobsgroups.GroupGradeRuleType_GROUP_GRADE_RULE_TYPE_RANGE,
		MinGrade: new(int32(2)),
		MaxGrade: new(int32(5)),
	})

	server.enrichGroupRuleGradeLabels("police", minimumRule, rangeRule)

	assert.Equal(t, "Rank 3", minimumRule.GetGrade().GetGradeLabel())
	assert.Equal(t, "Rank 2", rangeRule.GetGrade().GetMinGradeLabel())
	assert.Equal(t, "Rank 5", rangeRule.GetGrade().GetMaxGradeLabel())
}

func TestEnsureGroupRuleQualificationAccessAllowed(t *testing.T) {
	t.Parallel()

	qualificationAccess := &stubQualificationAccess{allowedIDs: []int64{10, 11}}
	server := &Server{qualificationAccess: qualificationAccess}
	rule := &jobsgroups.GroupRule{}
	rule.SetQualification(&jobsgroups.GroupQualificationRule{
		Type:             jobsgroups.GroupQualificationRuleType_GROUP_QUALIFICATION_RULE_TYPE_ALL,
		QualificationIds: []int64{10, 11},
	})

	err := server.ensureGroupRuleQualificationAccess(
		t.Context(),
		&pbuserinfo.UserInfo{UserId: 42, Job: "police"},
		rule,
	)

	require.NoError(t, err)
	assert.Equal(t, 1, qualificationAccess.calls)
	assert.Equal(
		t,
		int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		qualificationAccess.access,
	)
	assert.Equal(t, []int64{10, 11}, qualificationAccess.targetIDs)
}

func TestEnsureGroupRuleQualificationAccessDenied(t *testing.T) {
	t.Parallel()

	qualificationAccess := &stubQualificationAccess{allowedIDs: []int64{10}}
	server := &Server{qualificationAccess: qualificationAccess}
	rule := &jobsgroups.GroupRule{}
	rule.SetQualification(&jobsgroups.GroupQualificationRule{
		Type:             jobsgroups.GroupQualificationRuleType_GROUP_QUALIFICATION_RULE_TYPE_ANY,
		QualificationIds: []int64{10, 11},
	})

	err := server.ensureGroupRuleQualificationAccess(
		t.Context(),
		&pbuserinfo.UserInfo{UserId: 42, Job: "police"},
		rule,
	)

	require.ErrorIs(t, err, errorsjobs.ErrNotFoundOrNoPerms)
}

func TestEnsureGroupRuleQualificationAccessWrapsAccessError(t *testing.T) {
	t.Parallel()

	accessErr := errors.New("lookup failed")
	server := &Server{
		qualificationAccess: &stubQualificationAccess{err: accessErr},
	}
	rule := &jobsgroups.GroupRule{}
	rule.SetQualification(&jobsgroups.GroupQualificationRule{
		Type:             jobsgroups.GroupQualificationRuleType_GROUP_QUALIFICATION_RULE_TYPE_ANY,
		QualificationIds: []int64{10},
	})

	err := server.ensureGroupRuleQualificationAccess(
		t.Context(),
		&pbuserinfo.UserInfo{UserId: 42, Job: "police"},
		rule,
	)

	require.ErrorIs(t, err, accessErr)
	assert.Equal(t, codes.Internal, status.Code(err))
}

func TestEnsureGroupRuleQualificationAccessSkipsGradeRule(t *testing.T) {
	t.Parallel()

	qualificationAccess := &stubQualificationAccess{}
	server := &Server{qualificationAccess: qualificationAccess}
	rule := &jobsgroups.GroupRule{}
	rule.SetGrade(&jobsgroups.GroupGradeRule{
		Type:  jobsgroups.GroupGradeRuleType_GROUP_GRADE_RULE_TYPE_EXACT,
		Grade: new(int32),
	})

	err := server.ensureGroupRuleQualificationAccess(
		t.Context(),
		&pbuserinfo.UserInfo{UserId: 42, Job: "police"},
		rule,
	)

	require.NoError(t, err)
	assert.Zero(t, qualificationAccess.calls)
}
