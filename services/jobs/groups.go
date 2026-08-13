package jobs

import (
	"context"
	"strconv"

	resourcesaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	groupsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups/access"
	pbjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	colleaguehydrator "github.com/fivenet-app/fivenet/v2026/services/jobs/colleagues"
	errorsjobs "github.com/fivenet-app/fivenet/v2026/services/jobs/errors"
	jobsstore "github.com/fivenet-app/fivenet/v2026/stores/jobs"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) ensureHighestJobGradeAccess(
	job string,
	in *resourcesaccess.Access,
) (*resourcesaccess.Access, error) {
	if in == nil || s.enricher == nil {
		return in, nil
	}

	jobData := s.enricher.GetJobByName(job)
	if jobData == nil || len(jobData.GetGrades()) == 0 {
		return in, nil
	}

	highestGrade := jobData.GetGrades()[len(jobData.GetGrades())-1].GetGrade()
	if highestGrade <= 0 {
		return in, nil
	}

	return access.NormalizeAccess(
		in,
		&resourcesaccess.Access{
			Jobs: []*resourcesaccess.JobAccess{
				{
					Job:          job,
					MinimumGrade: highestGrade,
					Access:       int32(groupsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
					Required:     new(true),
				},
			},
		},
		nil,
		15,
	)
}

func (s *Server) ListGroups(
	ctx context.Context,
	req *pbjobs.ListGroupsRequest,
) (*pbjobs.ListGroupsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	count, err := s.store.CountGroups(ctx, s.db, jobsstore.GroupsQuery{
		Job:             userInfo.GetJob(),
		States:          req.GetStates(),
		Search:          req.GetSearch(),
		IncludeCounts:   req.GetIncludeCounts(),
		IncludeInactive: req.GetIncludeInactive(),
		IncludeArchived: req.GetIncludeArchived(),
		Sort:            req.GetSort(),
		Offset:          req.GetPagination().GetOffset(),
	}, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count, defaultPageSize)
	resp := &pbjobs.ListGroupsResponse{
		Pagination: pag,
		Groups:     []*jobsgroups.Group{},
	}
	if count <= 0 {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)
		return resp, nil
	}

	groups, err := s.store.ListGroups(ctx, s.db, jobsstore.GroupsQuery{
		Job:             userInfo.GetJob(),
		States:          req.GetStates(),
		Search:          req.GetSearch(),
		IncludeCounts:   req.GetIncludeCounts(),
		IncludeInactive: req.GetIncludeInactive(),
		IncludeArchived: req.GetIncludeArchived(),
		Sort:            req.GetSort(),
		Offset:          req.GetPagination().GetOffset(),
		Limit:           limit,
	}, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	resp.Groups = groups
	if err := s.hydrateGroupColleagueTargets(
		ctx,
		userInfo,
		appendGroupColleagueTargets(nil, resp.GetGroups()),
	); err != nil {
		return nil, err
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)
	return resp, nil
}

func (s *Server) GetGroup(
	ctx context.Context,
	req *pbjobs.GetGroupRequest,
) (*pbjobs.GetGroupResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.id", req.GetId(),
	})

	group, err := s.store.GetGroup(ctx, s.db, jobsstore.GroupQuery{
		Job:             userInfo.GetJob(),
		IncludeArchived: req.GetIncludeArchived(),
	}, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if group == nil {
		return nil, errorsjobs.ErrNotFoundOrNoPerms
	}
	if err := s.ensureGroupAccess(
		ctx,
		userInfo,
		group.GetId(),
		groupsaccess.AccessLevel_ACCESS_LEVEL_VIEW,
	); err != nil {
		return nil, err
	}

	resp := &pbjobs.GetGroupResponse{Group: group}
	if req.GetIncludeRules() {
		resp.Rules, err = s.store.ListGroupRules(ctx, s.db, group.GetId())
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		s.enrichGroupRuleGradeLabels(group.GetJob(), resp.Rules...)
	}
	if req.GetIncludeLeaders() {
		resp.Leaders, err = s.store.ListGroupLeaders(ctx, s.db, group.GetId(), "")
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}
	if req.GetIncludeManualMembers() {
		resp.ManualMembers, err = s.store.ListGroupManualMembers(ctx, s.db, group.GetId(), "")
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}
	if req.GetIncludeExclusions() {
		resp.Exclusions, err = s.store.ListGroupMemberExclusions(ctx, s.db, group.GetId(), "")
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}
	if req.GetIncludeResolvedMembers() {
		resp.ResolvedMembers, err = s.resolveGroupMembers(
			ctx,
			group,
			"",
			req.GetIncludeExclusions(),
			req.GetIncludeLeaders(),
			true,
			nil,
		)
		if err != nil {
			return nil, err
		}
	}
	if s.groupAccess != nil {
		resp.Access, err = s.groupAccess.ListTargetAccess(
			ctx,
			s.db,
			group.GetId(),
			access.SubjectAccessOptions{BlockedAccess: -1},
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	targets := []colleaguehydrator.Target{}
	targets = appendGroupColleagueTargets(targets, []*jobsgroups.Group{resp.GetGroup()})
	targets = appendGroupRuleColleagueTargets(targets, resp.GetRules())
	targets = appendGroupLeaderColleagueTargets(targets, resp.GetLeaders())
	targets = appendGroupManualMemberColleagueTargets(targets, resp.GetManualMembers())
	targets = appendGroupMemberExclusionColleagueTargets(targets, resp.GetExclusions())
	targets = appendGroupResolvedMemberColleagueTargets(targets, resp.GetResolvedMembers())
	if err := s.hydrateGroupColleagueTargets(ctx, userInfo, targets); err != nil {
		return nil, err
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)
	return resp, nil
}

func (s *Server) CreateGroup(
	ctx context.Context,
	req *pbjobs.CreateGroupRequest,
) (*pbjobs.CreateGroupResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	job := userInfo.GetJob()
	if req.GetJob() != "" && req.GetJob() != job {
		return nil, errorsjobs.ErrNotFoundOrNoPerms
	}

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.name", req.GetName(),
	})

	group := &jobsgroups.Group{
		Job:             job,
		Name:            req.GetName(),
		Type:            req.GetType(),
		MembershipMode:  req.GetMembershipMode(),
		CreatedByUserId: int32Ptr(userInfo.GetUserId()),
		UpdatedByUserId: int32Ptr(userInfo.GetUserId()),
	}
	if req.HasSortRank() {
		group.SortRank = req.GetSortRank()
	}
	if req.HasDescription() {
		description := req.GetDescription()
		group.Description = &description
	}
	if req.HasShortName() {
		shortName := req.GetShortName()
		group.ShortName = &shortName
	}
	if req.HasColor() {
		color := req.GetColor()
		group.Color = &color
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	defer tx.Rollback()

	id, err := s.store.CreateGroup(ctx, tx, group)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	group.Id = id
	if req.HasAccess() && s.groupAccess != nil {
		accessPayload, err := s.ensureHighestJobGradeAccess(job, req.GetAccess())
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		if _, err := s.groupAccess.ReplaceTargetAccess(
			ctx,
			tx,
			s.groupAccessResolver,
			id,
			accessPayload,
			access.SubjectAccessOptions{BlockedAccess: -1},
		); err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}
	if err := s.addGroupActivity(
		ctx,
		tx,
		job,
		id,
		jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_CREATED,
		userInfo.GetUserId(),
		0,
		0,
		nil,
		nil,
	); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	seenLeaders := map[int32]struct{}{}
	for _, userID := range req.GetLeaderUserIds() {
		if _, ok := seenLeaders[userID]; ok {
			continue
		}
		seenLeaders[userID] = struct{}{}
		if err := s.ensureGroupUserInJob(ctx, tx, job, userID); err != nil {
			return nil, err
		}
		if _, _, err := s.store.AddGroupLeader(
			ctx,
			tx,
			id,
			userID,
			userInfo.GetUserId(),
		); err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		if err := s.addGroupActivity(
			ctx,
			tx,
			job,
			id,
			jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_LEADER_ADDED,
			userInfo.GetUserId(),
			userID,
			0,
			nil,
			nil,
		); err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	for _, input := range req.GetRules() {
		rule, err := groupRuleFromInput(id, 0, userInfo.GetUserId(), input, true)
		if err != nil {
			return nil, err
		}
		createdRule, err := s.store.CreateGroupRule(ctx, tx, rule)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		if err := s.addGroupActivity(
			ctx,
			tx,
			job,
			id,
			jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_RULE_ADDED,
			userInfo.GetUserId(),
			0,
			createdRule.GetId(),
			nil,
			groupActivityRuleData(createdRule),
		); err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	seenMembers := map[int32]struct{}{}
	for _, userID := range req.GetManualMemberUserIds() {
		if _, ok := seenMembers[userID]; ok {
			continue
		}
		seenMembers[userID] = struct{}{}
		if err := s.ensureGroupUserInJob(ctx, tx, job, userID); err != nil {
			return nil, err
		}
		if group.GetMembershipMode() == jobsgroups.GroupMembershipMode_GROUP_MEMBERSHIP_MODE_STRICT {
			matches, err := s.userMatchesGroupRules(ctx, tx, group, userID)
			if err != nil {
				return nil, err
			}
			if !matches {
				return nil, errorsjobs.ErrGroupMemberRulesRequired
			}
		}
		if _, _, err := s.store.AddGroupManualMember(
			ctx,
			tx,
			id,
			userID,
			userInfo.GetUserId(),
			nil,
		); err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		if err := s.addGroupActivity(
			ctx,
			tx,
			job,
			id,
			jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_MEMBER_ADDED,
			userInfo.GetUserId(),
			userID,
			0,
			nil,
			nil,
		); err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	created, err := s.recountAndGetGroup(ctx, tx, job, id)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	grpc_audit.AddMeta(ctx, "jobs.group.id", strconv.FormatInt(id, 10))
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

	if err := s.hydrateGroupColleagues(ctx, userInfo, created); err != nil {
		return nil, err
	}

	return &pbjobs.CreateGroupResponse{Group: created}, nil
}

func (s *Server) UpdateGroup(
	ctx context.Context,
	req *pbjobs.UpdateGroupRequest,
) (*pbjobs.UpdateGroupResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.id", req.GetId(),
	})

	group, err := s.store.GetGroup(ctx, s.db, jobsstore.GroupQuery{
		Job:             userInfo.GetJob(),
		IncludeArchived: false,
	}, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if group == nil {
		return nil, errorsjobs.ErrNotFoundOrNoPerms
	}
	if err := s.ensureGroupAccess(
		ctx,
		userInfo,
		group.GetId(),
		groupsaccess.AccessLevel_ACCESS_LEVEL_EDIT,
	); err != nil {
		return nil, err
	}

	if req.HasName() {
		group.Name = req.GetName()
	}
	if req.HasDescription() {
		description := req.GetDescription()
		group.Description = &description
	}
	if req.HasShortName() {
		shortName := req.GetShortName()
		group.ShortName = &shortName
	}
	if req.HasColor() {
		color := req.GetColor()
		group.Color = &color
	}
	if req.HasType() && req.GetType() != jobsgroups.GroupType_GROUP_TYPE_UNSPECIFIED {
		group.Type = req.GetType()
	}
	if req.HasMembershipMode() &&
		req.GetMembershipMode() != jobsgroups.GroupMembershipMode_GROUP_MEMBERSHIP_MODE_UNSPECIFIED {
		group.MembershipMode = req.GetMembershipMode()
	}
	if req.HasSortRank() {
		group.SortRank = req.GetSortRank()
	}
	if req.HasState() {
		if req.GetState() == jobsgroups.GroupState_GROUP_STATE_ARCHIVED {
			return nil, errorsjobs.ErrNotFoundOrNoPerms
		}
		if req.GetState() != jobsgroups.GroupState_GROUP_STATE_UNSPECIFIED {
			group.State = req.GetState()
		}
	}

	group.UpdatedByUserId = int32Ptr(userInfo.GetUserId())

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	defer tx.Rollback()

	if err := s.store.UpdateGroup(ctx, tx, group); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if req.HasAccess() && s.groupAccess != nil {
		accessPayload, err := s.ensureHighestJobGradeAccess(group.GetJob(), req.GetAccess())
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		if _, err := s.groupAccess.ReplaceTargetAccess(
			ctx,
			tx,
			s.groupAccessResolver,
			group.GetId(),
			accessPayload,
			access.SubjectAccessOptions{BlockedAccess: -1},
		); err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}
	if err := s.addGroupActivity(
		ctx,
		tx,
		userInfo.GetJob(),
		group.GetId(),
		jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_UPDATED,
		userInfo.GetUserId(),
		0,
		0,
		nil,
		nil,
	); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	updated, err := s.store.GetGroup(ctx, tx, jobsstore.GroupQuery{
		Job:             userInfo.GetJob(),
		IncludeArchived: true,
	}, group.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if updated == nil {
		return nil, errorsjobs.ErrFailedQuery
	}
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	if err := s.hydrateGroupColleagues(ctx, userInfo, updated); err != nil {
		return nil, err
	}

	return &pbjobs.UpdateGroupResponse{Group: updated}, nil
}

func (s *Server) ArchiveGroup(
	ctx context.Context,
	req *pbjobs.ArchiveGroupRequest,
) (*pbjobs.ArchiveGroupResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.id", req.GetId(),
	})

	if err := s.ensureGroupAccess(
		ctx,
		userInfo,
		req.GetId(),
		groupsaccess.AccessLevel_ACCESS_LEVEL_MANAGE,
	); err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	defer tx.Rollback()

	if err := s.store.ArchiveGroup(
		ctx,
		tx,
		userInfo.GetJob(),
		req.GetId(),
		userInfo.GetUserId(),
	); err != nil {
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
		req.GetId(),
		jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_ARCHIVED,
		userInfo.GetUserId(),
		0,
		0,
		reason,
		nil,
	); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	group, err := s.store.GetGroup(ctx, tx, jobsstore.GroupQuery{
		Job:             userInfo.GetJob(),
		IncludeArchived: true,
	}, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if group == nil {
		return nil, errorsjobs.ErrNotFoundOrNoPerms
	}
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	if req.HasReason() {
		grpc_audit.AddMeta(ctx, "jobs.group.reason", req.GetReason())
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	if err := s.hydrateGroupColleagues(ctx, userInfo, group); err != nil {
		return nil, err
	}

	return &pbjobs.ArchiveGroupResponse{Group: group}, nil
}

func (s *Server) RestoreGroup(
	ctx context.Context,
	req *pbjobs.RestoreGroupRequest,
) (*pbjobs.RestoreGroupResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.id", req.GetId(),
	})

	if err := s.ensureGroupAccess(
		ctx,
		userInfo,
		req.GetId(),
		groupsaccess.AccessLevel_ACCESS_LEVEL_MANAGE,
	); err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	defer tx.Rollback()

	if err := s.store.RestoreGroup(
		ctx,
		tx,
		userInfo.GetJob(),
		req.GetId(),
		userInfo.GetUserId(),
	); err != nil {
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
		req.GetId(),
		jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_RESTORED,
		userInfo.GetUserId(),
		0,
		0,
		reason,
		nil,
	); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	group, err := s.store.GetGroup(ctx, tx, jobsstore.GroupQuery{
		Job:             userInfo.GetJob(),
		IncludeArchived: true,
	}, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if group == nil {
		return nil, errorsjobs.ErrNotFoundOrNoPerms
	}
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	if req.HasReason() {
		grpc_audit.AddMeta(ctx, "jobs.group.reason", req.GetReason())
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_RESTORED)
	if err := s.hydrateGroupColleagues(ctx, userInfo, group); err != nil {
		return nil, err
	}

	return &pbjobs.RestoreGroupResponse{Group: group}, nil
}
