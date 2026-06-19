package jobs

import (
	"testing"

	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/stretchr/testify/assert"
)

func TestEnrichGroupRuleGradeLabels(t *testing.T) {
	t.Parallel()

	server := &Server{enricher: mstlystcdata.NewDummyUserAwareEnricher()}
	minimumRule := &jobsgroups.GroupRule{}
	minimumRule.SetGrade(&jobsgroups.GroupGradeRule{
		Type:  jobsgroups.GroupGradeRuleType_GROUP_GRADE_RULE_TYPE_MINIMUM,
		Grade: int32Ptr(3),
	})
	rangeRule := &jobsgroups.GroupRule{}
	rangeRule.SetGrade(&jobsgroups.GroupGradeRule{
		Type:     jobsgroups.GroupGradeRuleType_GROUP_GRADE_RULE_TYPE_RANGE,
		MinGrade: int32Ptr(2),
		MaxGrade: int32Ptr(5),
	})

	server.enrichGroupRuleGradeLabels("police", minimumRule, rangeRule)

	assert.Equal(t, "Rank 3", minimumRule.GetGrade().GetGradeLabel())
	assert.Equal(t, "Rank 2", rangeRule.GetGrade().GetMinGradeLabel())
	assert.Equal(t, "Rank 5", rangeRule.GetGrade().GetMaxGradeLabel())
}
