package documents

import timestamp "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"

func (x *ApprovalTask) GetJobGrade() int32 {
	return x.GetMinimumGrade()
}

func (x *ApprovalTask) SetJobGrade(grade int32) {
	x.MinimumGrade = &grade
}

func (x *ApprovalTask) SetJobGradeLabel(label string) {
	x.JobGradeLabel = &label
}

func (x *ApprovalTask) SetJob(job string) {
	x.Job = &job
}

func (x *ApprovalTask) SetJobLabel(label string) {
	x.JobLabel = &label
}

func (x *Approval) GetJob() string {
	return *x.UserJob
}

func (x *Approval) GetJobGrade() int32 {
	return x.GetUserGrade()
}

func (x *Approval) SetJob(job string) {
	x.UserJob = &job
}

func (x *Approval) SetJobLabel(label string) {
	x.UserJobLabel = &label
}

func (x *Approval) SetJobGrade(grade int32) {
	x.UserGrade = &grade
}

func (x *Approval) SetJobGradeLabel(label string) {
	x.UserGradeLabel = &label
}

func (x *ApprovalPolicy) Default() {
	if x.SnapshotDate == nil {
		x.SnapshotDate = timestamp.Now()
	}

	if x.RuleKind == ApprovalRuleKind_APPROVAL_RULE_KIND_UNSPECIFIED {
		x.RuleKind = ApprovalRuleKind_APPROVAL_RULE_KIND_REQUIRE_ALL
	}

	if x.RuleKind == ApprovalRuleKind_APPROVAL_RULE_KIND_QUORUM_ANY && x.GetRequiredCount() == 0 {
		requiredCount := int32(0)
		x.RequiredCount = &requiredCount
	}

	if x.OnEditBehavior == OnEditBehavior_ON_EDIT_BEHAVIOR_UNSPECIFIED {
		x.OnEditBehavior = OnEditBehavior_ON_EDIT_BEHAVIOR_KEEP_PROGRESS
	}
}
