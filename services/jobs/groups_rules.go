package jobs

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	groupsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups/access"
	pbjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorsjobs "github.com/fivenet-app/fivenet/v2026/services/jobs/errors"
	jobspolicy "github.com/fivenet-app/fivenet/v2026/stores/jobs/jobspolicy"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func groupRuleFromInput(
	groupID int64,
	ruleID int64,
	createdByUserID int32,
	input *pbjobs.GroupRuleInput,
	defaultEnabled bool,
) (*jobsgroups.GroupRule, error) {
	if input == nil || !input.HasRule() {
		return nil, status.Error(codes.InvalidArgument, "job group rule is required")
	}

	enabled := defaultEnabled
	if input.HasEnabled() {
		enabled = input.GetEnabled()
	}

	rule := &jobsgroups.GroupRule{
		Id:              ruleID,
		GroupId:         groupID,
		Enabled:         enabled,
		CreatedByUserId: &createdByUserID,
	}

	switch input.WhichRule() {
	case pbjobs.GroupRuleInput_Grade_case:
		grade := input.GetGrade()
		if err := validateGradeRule(grade); err != nil {
			return nil, err
		}
		rule.Type = jobsgroups.GroupRuleType_GROUP_RULE_TYPE_GRADE
		rule.SetGrade(grade)
	case pbjobs.GroupRuleInput_Qualification_case:
		qualification := input.GetQualification()
		if err := validateQualificationRule(qualification); err != nil {
			return nil, err
		}
		rule.Type = jobsgroups.GroupRuleType_GROUP_RULE_TYPE_QUALIFICATION
		rule.SetQualification(qualification)
	default:
		return nil, status.Error(codes.InvalidArgument, "unsupported job group rule type")
	}

	return rule, nil
}

func validateGradeRule(rule *jobsgroups.GroupGradeRule) error {
	if rule == nil {
		return status.Error(codes.InvalidArgument, "grade rule is required")
	}

	switch rule.GetType() {
	case jobsgroups.GroupGradeRuleType_GROUP_GRADE_RULE_TYPE_MINIMUM,
		jobsgroups.GroupGradeRuleType_GROUP_GRADE_RULE_TYPE_EXACT:
		if !rule.HasGrade() {
			return status.Error(codes.InvalidArgument, "grade rule grade is required")
		}
	case jobsgroups.GroupGradeRuleType_GROUP_GRADE_RULE_TYPE_RANGE:
		if !rule.HasMinGrade() || !rule.HasMaxGrade() {
			return status.Error(codes.InvalidArgument, "grade rule min and max grade are required")
		}
		if rule.GetMinGrade() > rule.GetMaxGrade() {
			return status.Error(codes.InvalidArgument, "grade rule min grade must be <= max grade")
		}
	default:
		return status.Error(codes.InvalidArgument, "unsupported grade rule type")
	}

	return nil
}

func validateQualificationRule(rule *jobsgroups.GroupQualificationRule) error {
	if rule == nil {
		return status.Error(codes.InvalidArgument, "qualification rule is required")
	}

	switch rule.GetType() {
	case jobsgroups.GroupQualificationRuleType_GROUP_QUALIFICATION_RULE_TYPE_ALL,
		jobsgroups.GroupQualificationRuleType_GROUP_QUALIFICATION_RULE_TYPE_ANY:
	default:
		return status.Error(codes.InvalidArgument, "unsupported qualification rule type")
	}

	seen := map[int64]struct{}{}
	for _, qualificationID := range rule.GetQualificationIds() {
		if qualificationID <= 0 {
			return status.Error(codes.InvalidArgument, "qualification IDs must be positive")
		}
		if _, ok := seen[qualificationID]; ok {
			return status.Error(codes.InvalidArgument, "qualification IDs must be unique")
		}
		seen[qualificationID] = struct{}{}
	}
	if len(seen) == 0 {
		return status.Error(codes.InvalidArgument, "at least one qualification ID is required")
	}

	return nil
}

func (s *Server) enrichGroupRuleGradeLabels(job string, rules ...*jobsgroups.GroupRule) {
	if s.enricher == nil {
		return
	}

	for _, rule := range rules {
		grade := rule.GetGrade()
		if grade == nil {
			continue
		}

		if grade.HasGrade() {
			if _, jobGrade := s.enricher.GetJobGrade(job, grade.GetGrade()); jobGrade != nil {
				grade.SetGradeLabel(jobGrade.GetLabel())
			}
		}
		if grade.HasMinGrade() {
			if _, jobGrade := s.enricher.GetJobGrade(job, grade.GetMinGrade()); jobGrade != nil {
				grade.SetMinGradeLabel(jobGrade.GetLabel())
			}
		}
		if grade.HasMaxGrade() {
			if _, jobGrade := s.enricher.GetJobGrade(job, grade.GetMaxGrade()); jobGrade != nil {
				grade.SetMaxGradeLabel(jobGrade.GetLabel())
			}
		}
	}
}

func (s *Server) CreateGroupRule(
	ctx context.Context,
	req *pbjobs.CreateGroupRuleRequest,
) (*pbjobs.CreateGroupRuleResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.id", req.GetGroupId(),
	})

	rule, err := groupRuleFromInput(
		req.GetGroupId(),
		0,
		userInfo.GetUserId(),
		req.GetRule(),
		true,
	)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	defer tx.Rollback()

	group, err := s.getActiveGroupForJob(ctx, tx, userInfo.GetJob(), req.GetGroupId())
	if err != nil {
		return nil, err
	}
	if err := s.ensureGroupAccess(
		ctx,
		userInfo,
		req.GetGroupId(),
		groupsaccess.AccessLevel_ACCESS_LEVEL_EDIT,
	); err != nil {
		return nil, err
	}
	if err := validateGroupPolicyAllowedMutation(group, jobspolicy.MutationRuleAdd); err != nil {
		return nil, err
	}

	created, err := s.store.CreateGroupRule(ctx, tx, rule)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	var reason *string
	if req.HasReason() {
		reasonValue := req.GetReason()
		reason = &reasonValue
	}
	if err := s.addGroupActivity(
		ctx,
		tx,
		userInfo.GetJob(),
		req.GetGroupId(),
		jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_RULE_ADDED,
		userInfo.GetUserId(),
		0,
		created.GetId(),
		reason,
		groupActivityRuleData(created),
	); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	group, err = s.recountAndGetGroup(ctx, tx, userInfo.GetJob(), req.GetGroupId())
	if err != nil {
		return nil, err
	}
	if err := s.validateGroupPolicyAgainstExistingData(ctx, tx, group); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	grpc_audit.AddMeta(ctx, "jobs.group.id", strconv.FormatInt(req.GetGroupId(), 10))
	grpc_audit.AddMeta(ctx, "jobs.group.rule.id", strconv.FormatInt(created.GetId(), 10))
	if req.HasReason() {
		grpc_audit.AddMeta(ctx, "jobs.group.reason", req.GetReason())
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

	s.enrichGroupRuleGradeLabels(userInfo.GetJob(), created)
	targets := appendGroupColleagueTargets(nil, []*jobsgroups.Group{group})
	targets = appendGroupRuleColleagueTargets(targets, []*jobsgroups.GroupRule{created})
	if err := s.hydrateGroupColleagueTargets(ctx, userInfo, targets); err != nil {
		return nil, err
	}

	return &pbjobs.CreateGroupRuleResponse{Rule: created, Group: group}, nil
}

func (s *Server) ListGroupRules(
	ctx context.Context,
	req *pbjobs.ListGroupRulesRequest,
) (*pbjobs.ListGroupRulesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.id", req.GetGroupId(),
	})

	group, err := s.getActiveGroupForJob(ctx, s.db, userInfo.GetJob(), req.GetGroupId())
	if err != nil {
		return nil, err
	}
	if err := s.ensureGroupAccess(
		ctx,
		userInfo,
		req.GetGroupId(),
		groupsaccess.AccessLevel_ACCESS_LEVEL_VIEW,
	); err != nil {
		return nil, err
	}

	rules, err := s.store.ListGroupRules(ctx, s.db, req.GetGroupId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	pag, paged := paginateGroupItems(req.GetPagination(), rules)
	s.enrichGroupRuleGradeLabels(group.GetJob(), paged...)

	resp := &pbjobs.ListGroupRulesResponse{
		Pagination: pag,
		Rules:      paged,
	}
	targets := appendGroupRuleColleagueTargets(nil, resp.GetRules())
	if err := s.hydrateGroupColleagueTargets(ctx, userInfo, targets); err != nil {
		return nil, err
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)
	return resp, nil
}

func (s *Server) UpdateGroupRule(
	ctx context.Context,
	req *pbjobs.UpdateGroupRuleRequest,
) (*pbjobs.UpdateGroupRuleResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.id", req.GetGroupId(),
		"fivenet.jobs.groups.rule.id", req.GetRuleId(),
	})

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	defer tx.Rollback()

	group, err := s.getActiveGroupForJob(ctx, tx, userInfo.GetJob(), req.GetGroupId())
	if err != nil {
		return nil, err
	}
	if err := s.ensureGroupAccess(
		ctx,
		userInfo,
		req.GetGroupId(),
		groupsaccess.AccessLevel_ACCESS_LEVEL_EDIT,
	); err != nil {
		return nil, err
	}
	if err := validateGroupPolicyAllowedMutation(group, jobspolicy.MutationRuleUpdate); err != nil {
		return nil, err
	}

	existing, err := s.store.GetGroupRule(ctx, tx, req.GetGroupId(), req.GetRuleId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if existing == nil {
		return nil, errorsjobs.ErrNotFoundOrNoPerms
	}

	rule, err := groupRuleFromInput(
		req.GetGroupId(),
		req.GetRuleId(),
		existing.GetCreatedByUserId(),
		req.GetRule(),
		existing.GetEnabled(),
	)
	if err != nil {
		return nil, err
	}

	updated, err := s.store.UpdateGroupRule(ctx, tx, rule, userInfo.GetUserId())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorsjobs.ErrNotFoundOrNoPerms
		}
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	var reason *string
	if req.HasReason() {
		reasonValue := req.GetReason()
		reason = &reasonValue
	}
	if err := s.addGroupActivity(
		ctx,
		tx,
		userInfo.GetJob(),
		req.GetGroupId(),
		jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_RULE_UPDATED,
		userInfo.GetUserId(),
		0,
		updated.GetId(),
		reason,
		groupActivityRuleData(updated),
	); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	group, err = s.recountAndGetGroup(ctx, tx, userInfo.GetJob(), req.GetGroupId())
	if err != nil {
		return nil, err
	}
	if err := s.validateGroupPolicyAgainstExistingData(ctx, tx, group); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	grpc_audit.AddMeta(ctx, "jobs.group.id", strconv.FormatInt(req.GetGroupId(), 10))
	grpc_audit.AddMeta(ctx, "jobs.group.rule.id", strconv.FormatInt(req.GetRuleId(), 10))
	if req.HasReason() {
		grpc_audit.AddMeta(ctx, "jobs.group.reason", req.GetReason())
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	s.enrichGroupRuleGradeLabels(userInfo.GetJob(), updated)
	targets := appendGroupColleagueTargets(nil, []*jobsgroups.Group{group})
	targets = appendGroupRuleColleagueTargets(targets, []*jobsgroups.GroupRule{updated})
	if err := s.hydrateGroupColleagueTargets(ctx, userInfo, targets); err != nil {
		return nil, err
	}

	return &pbjobs.UpdateGroupRuleResponse{Rule: updated, Group: group}, nil
}

func (s *Server) DeleteGroupRule(
	ctx context.Context,
	req *pbjobs.DeleteGroupRuleRequest,
) (*pbjobs.DeleteGroupRuleResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.id", req.GetGroupId(),
		"fivenet.jobs.groups.rule.id", req.GetRuleId(),
	})

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	defer tx.Rollback()

	if _, err := s.getActiveGroupForJob(ctx, tx, userInfo.GetJob(), req.GetGroupId()); err != nil {
		return nil, err
	}
	if err := s.ensureGroupAccess(
		ctx,
		userInfo,
		req.GetGroupId(),
		groupsaccess.AccessLevel_ACCESS_LEVEL_EDIT,
	); err != nil {
		return nil, err
	}

	existing, err := s.store.GetGroupRule(ctx, tx, req.GetGroupId(), req.GetRuleId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if existing == nil {
		return nil, errorsjobs.ErrNotFoundOrNoPerms
	}

	if err := s.store.DeleteGroupRule(ctx, tx, req.GetGroupId(), req.GetRuleId()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorsjobs.ErrNotFoundOrNoPerms
		}
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	var reason *string
	if req.HasReason() {
		reasonValue := req.GetReason()
		reason = &reasonValue
	}
	if err := s.addGroupActivity(
		ctx,
		tx,
		userInfo.GetJob(),
		req.GetGroupId(),
		jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_RULE_REMOVED,
		userInfo.GetUserId(),
		0,
		req.GetRuleId(),
		reason,
		groupActivityRuleData(existing),
	); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	group, err := s.recountAndGetGroup(ctx, tx, userInfo.GetJob(), req.GetGroupId())
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	grpc_audit.AddMeta(ctx, "jobs.group.id", strconv.FormatInt(req.GetGroupId(), 10))
	grpc_audit.AddMeta(ctx, "jobs.group.rule.id", strconv.FormatInt(req.GetRuleId(), 10))
	if req.HasReason() {
		grpc_audit.AddMeta(ctx, "jobs.group.reason", req.GetReason())
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)

	if err := s.hydrateGroupColleagues(ctx, userInfo, group); err != nil {
		return nil, err
	}

	return &pbjobs.DeleteGroupRuleResponse{Group: group}, nil
}
