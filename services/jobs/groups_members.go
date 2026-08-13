package jobs

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"sort"
	"strconv"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	groupsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	colleaguehydrator "github.com/fivenet-app/fivenet/v2026/services/jobs/colleagues"
	errorsjobs "github.com/fivenet-app/fivenet/v2026/services/jobs/errors"
	jobsstore "github.com/fivenet-app/fivenet/v2026/stores/jobs"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

type groupMemberBundle struct {
	member    *jobsgroups.GroupResolvedMember
	sources   map[jobsgroups.GroupMemberSource]struct{}
	reasonKey map[string]struct{}
}

func newGroupMemberBundle(groupID int64, userID int32) *groupMemberBundle {
	return &groupMemberBundle{
		member: &jobsgroups.GroupResolvedMember{
			GroupId: groupID,
			UserId:  userID,
		},
		sources:   map[jobsgroups.GroupMemberSource]struct{}{},
		reasonKey: map[string]struct{}{},
	}
}

func (b *groupMemberBundle) addSource(source jobsgroups.GroupMemberSource) {
	if source == jobsgroups.GroupMemberSource_GROUP_MEMBER_SOURCE_UNSPECIFIED {
		return
	}
	b.sources[source] = struct{}{}
}

func (b *groupMemberBundle) hasSource(source jobsgroups.GroupMemberSource) bool {
	_, ok := b.sources[source]
	return ok
}

func (b *groupMemberBundle) addReason(
	source jobsgroups.GroupMemberSource,
	reasonType jobsgroups.GroupMembershipReasonType,
	detail *string,
	ruleID *int64,
) {
	key := strconv.Itoa(int(source)) + ":" + strconv.Itoa(int(reasonType))
	if detail != nil {
		key += ":" + *detail
	}
	if ruleID != nil {
		key += ":" + strconv.FormatInt(*ruleID, 10)
	}
	if _, ok := b.reasonKey[key]; ok {
		return
	}
	b.reasonKey[key] = struct{}{}
	b.member.Reasons = append(b.member.Reasons, &jobsgroups.GroupMembershipReason{
		Source: source,
		RuleId: ruleID,
		Type:   reasonType,
		Detail: detail,
	})
}

func (b *groupMemberBundle) finalize(includeReasons bool) {
	if len(b.sources) > 0 {
		b.member.Sources = make([]jobsgroups.GroupMemberSource, 0, len(b.sources))
		for source := range b.sources {
			b.member.Sources = append(b.member.Sources, source)
		}
		slices.SortFunc(b.member.Sources, func(a, b jobsgroups.GroupMemberSource) int {
			return int(a) - int(b)
		})
	}
	if !includeReasons {
		b.member.Reasons = nil
	}
}

func resolvedMemberMatchesSources(
	member *jobsgroups.GroupResolvedMember,
	sources []jobsgroups.GroupMemberSource,
) bool {
	if len(sources) == 0 {
		return true
	}

	hasFilter := false
	for _, source := range sources {
		if source == jobsgroups.GroupMemberSource_GROUP_MEMBER_SOURCE_UNSPECIFIED {
			continue
		}
		hasFilter = true
		if slices.Contains(member.GetSources(), source) {
			return true
		}
	}

	return !hasFilter
}

func compareResolvedMembersBySort(
	a *jobsgroups.GroupResolvedMember,
	b *jobsgroups.GroupResolvedMember,
	sort *database.Sort,
) int {
	if sort != nil {
		for _, column := range sort.GetColumns() {
			result, ok := compareResolvedMembersByColumn(a, b, column.GetId())
			if !ok || result == 0 {
				continue
			}
			if column.GetDesc() {
				return -result
			}
			return result
		}
	}

	if a.GetUserId() == b.GetUserId() {
		return compareInt64(a.GetGroupId(), b.GetGroupId())
	}
	return compareInt32(a.GetUserId(), b.GetUserId())
}

func compareResolvedMembersByColumn(
	a *jobsgroups.GroupResolvedMember,
	b *jobsgroups.GroupResolvedMember,
	column string,
) (int, bool) {
	switch column {
	case "user_id":
		return compareInt32(a.GetUserId(), b.GetUserId()), true
	case "group_id":
		return compareInt64(a.GetGroupId(), b.GetGroupId()), true
	case "is_member":
		return compareBool(a.GetIsMember(), b.GetIsMember()), true
	case "is_excluded":
		return compareBool(a.GetIsExcluded(), b.GetIsExcluded()), true
	case "is_leader":
		return compareBool(a.GetIsLeader(), b.GetIsLeader()), true
	case "sources_count":
		return compareInt64(int64(len(a.GetSources())), int64(len(b.GetSources()))), true
	default:
		return 0, false
	}
}

func compareInt32(a int32, b int32) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

func compareInt64(a int64, b int64) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

func compareBool(a bool, b bool) int {
	switch {
	case a == b:
		return 0
	case !a && b:
		return -1
	default:
		return 1
	}
}

func appendGroupColleagueTarget(
	targets []colleaguehydrator.Target,
	userID int32,
	set func(*jobscolleagues.Colleague),
) []colleaguehydrator.Target {
	if userID <= 0 || set == nil {
		return targets
	}

	return append(targets, colleaguehydrator.Target{
		UserID: userID,
		Set:    set,
	})
}

func appendGroupColleagueTargets(
	targets []colleaguehydrator.Target,
	groups []*jobsgroups.Group,
) []colleaguehydrator.Target {
	for _, group := range groups {
		targets = appendGroupColleagueTarget(
			targets,
			group.GetCreatedByUserId(),
			group.SetCreatedBy,
		)
	}

	return targets
}

func appendGroupRuleColleagueTargets(
	targets []colleaguehydrator.Target,
	rules []*jobsgroups.GroupRule,
) []colleaguehydrator.Target {
	for _, rule := range rules {
		targets = appendGroupColleagueTarget(
			targets,
			rule.GetCreatedByUserId(),
			rule.SetCreatedBy,
		)
	}

	return targets
}

func appendGroupActivityColleagueTargets(
	targets []colleaguehydrator.Target,
	activities []*jobsgroups.GroupActivity,
) []colleaguehydrator.Target {
	for _, activity := range activities {
		targets = appendGroupColleagueTarget(
			targets,
			activity.GetActorUserId(),
			activity.SetActorUser,
		)
		targets = appendGroupColleagueTarget(
			targets,
			activity.GetTargetUserId(),
			activity.SetTargetUser,
		)
	}

	return targets
}

func appendGroupLeaderColleagueTargets(
	targets []colleaguehydrator.Target,
	leaders []*jobsgroups.GroupLeader,
) []colleaguehydrator.Target {
	for _, leader := range leaders {
		targets = appendGroupColleagueTarget(targets, leader.GetUserId(), leader.SetColleague)
		targets = appendGroupColleagueTarget(
			targets,
			leader.GetCreatedByUserId(),
			leader.SetCreatedBy,
		)
	}

	return targets
}

func appendGroupManualMemberColleagueTargets(
	targets []colleaguehydrator.Target,
	members []*jobsgroups.GroupManualMember,
) []colleaguehydrator.Target {
	for _, member := range members {
		targets = appendGroupColleagueTarget(targets, member.GetUserId(), member.SetColleague)
		targets = appendGroupColleagueTarget(
			targets,
			member.GetCreatedByUserId(),
			member.SetCreatedBy,
		)
	}

	return targets
}

func appendGroupMemberExclusionColleagueTargets(
	targets []colleaguehydrator.Target,
	exclusions []*jobsgroups.GroupMemberExclusion,
) []colleaguehydrator.Target {
	for _, exclusion := range exclusions {
		targets = appendGroupColleagueTarget(targets, exclusion.GetUserId(), exclusion.SetColleague)
		targets = appendGroupColleagueTarget(
			targets,
			exclusion.GetCreatedByUserId(),
			exclusion.SetCreatedBy,
		)
	}

	return targets
}

func appendGroupResolvedMemberColleagueTargets(
	targets []colleaguehydrator.Target,
	members []*jobsgroups.GroupResolvedMember,
) []colleaguehydrator.Target {
	for _, member := range members {
		targets = appendGroupColleagueTarget(targets, member.GetUserId(), member.SetColleague)
	}

	return targets
}

func (s *Server) hydrateGroupColleagueTargets(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	targets []colleaguehydrator.Target,
) error {
	if len(targets) == 0 {
		return nil
	}

	if err := s.colleagueHydrator.HydrateTargets(
		ctx,
		s.db,
		userInfo,
		userInfo.GetJob(),
		targets,
		false,
	); err != nil {
		return errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	return nil
}

func (s *Server) hydrateGroupColleagues(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	groups ...*jobsgroups.Group,
) error {
	return s.hydrateGroupColleagueTargets(
		ctx,
		userInfo,
		appendGroupColleagueTargets(nil, groups),
	)
}

func (s *Server) getActiveGroupForJob(
	ctx context.Context,
	db qrm.DB,
	job string,
	groupID int64,
) (*jobsgroups.Group, error) {
	group, err := s.store.GetGroup(ctx, db, jobsstore.GroupQuery{
		Job:             job,
		IncludeArchived: false,
	}, groupID)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if group == nil {
		return nil, errorsjobs.ErrNotFoundOrNoPerms
	}

	return group, nil
}

func (s *Server) ensureGroupUserInJob(
	ctx context.Context,
	db qrm.DB,
	job string,
	userID int32,
) error {
	if userID <= 0 {
		return errorsjobs.ErrNotFoundOrNoPerms
	}

	ok, err := s.store.UserInJob(ctx, db, job, userID)
	if err != nil {
		return errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if !ok {
		return errorsjobs.ErrNotFoundOrNoPerms
	}

	return nil
}

func (s *Server) userMatchesGroupRules(
	ctx context.Context,
	db qrm.DB,
	group *jobsgroups.Group,
	userID int32,
) (bool, error) {
	ruleMatches, err := s.store.ListGroupRuleMemberMatches(ctx, db, group, "")
	if err != nil {
		return false, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	for _, match := range ruleMatches {
		if match.UserID == userID {
			return true, nil
		}
	}

	return false, nil
}

func (s *Server) recountAndGetGroup(
	ctx context.Context,
	db qrm.DB,
	job string,
	groupID int64,
) (*jobsgroups.Group, error) {
	if err := s.store.RecountGroupStats(ctx, db, groupID); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	group, err := s.store.GetGroup(ctx, db, jobsstore.GroupQuery{
		Job:             job,
		IncludeArchived: true,
	}, groupID)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if group == nil {
		return nil, errorsjobs.ErrFailedQuery
	}

	return group, nil
}

func (s *Server) resolveGroupMembers(
	ctx context.Context,
	group *jobsgroups.Group,
	search string,
	includeExcluded bool,
	includeLeaders bool,
	includeReasons bool,
	sources []jobsgroups.GroupMemberSource,
) ([]*jobsgroups.GroupResolvedMember, error) {
	groupID := group.GetId()
	manualMembers, err := s.store.ListGroupManualMembers(ctx, s.db, groupID, search)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	exclusions, err := s.store.ListGroupMemberExclusions(ctx, s.db, groupID, search)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	wantManual := len(sources) == 0 ||
		slices.Contains(sources, jobsgroups.GroupMemberSource_GROUP_MEMBER_SOURCE_MANUAL)
	wantRules := len(sources) == 0 ||
		slices.Contains(sources, jobsgroups.GroupMemberSource_GROUP_MEMBER_SOURCE_RULE)
	needsRulesForStrictManual := wantManual &&
		group.GetMembershipMode() == jobsgroups.GroupMembershipMode_GROUP_MEMBERSHIP_MODE_STRICT
	wantLeaders := includeLeaders ||
		slices.Contains(sources, jobsgroups.GroupMemberSource_GROUP_MEMBER_SOURCE_LEADER)
	var ruleMatches []*jobsstore.GroupRuleMemberMatch
	if wantRules || needsRulesForStrictManual {
		ruleMatches, err = s.store.ListGroupRuleMemberMatches(ctx, s.db, group, search)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}
	var leaders []*jobsgroups.GroupLeader
	if wantLeaders {
		leaders, err = s.store.ListGroupLeaders(ctx, s.db, groupID, search)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	bundles := map[int32]*groupMemberBundle{}
	getBundle := func(userID int32) *groupMemberBundle {
		bundle, ok := bundles[userID]
		if !ok {
			bundle = newGroupMemberBundle(groupID, userID)
			bundles[userID] = bundle
		}
		return bundle
	}

	if wantManual {
		for _, member := range manualMembers {
			bundle := getBundle(member.GetUserId())
			bundle.member.IsMember = true
			bundle.addSource(jobsgroups.GroupMemberSource_GROUP_MEMBER_SOURCE_MANUAL)
			if includeReasons {
				var detail *string
				if member.GetReason() != "" {
					reason := member.GetReason()
					detail = &reason
				}
				bundle.addReason(
					jobsgroups.GroupMemberSource_GROUP_MEMBER_SOURCE_MANUAL,
					jobsgroups.GroupMembershipReasonType_GROUP_MEMBERSHIP_REASON_TYPE_MANUAL,
					detail,
					nil,
				)
			}
		}
	}

	if wantRules {
		for _, match := range ruleMatches {
			bundle := getBundle(match.UserID)
			bundle.member.IsMember = true
			bundle.addSource(jobsgroups.GroupMemberSource_GROUP_MEMBER_SOURCE_RULE)
			if includeReasons {
				ruleID := match.RuleID
				bundle.addReason(
					jobsgroups.GroupMemberSource_GROUP_MEMBER_SOURCE_RULE,
					jobsgroups.GroupMembershipReasonType_GROUP_MEMBERSHIP_REASON_TYPE_RULE,
					nil,
					&ruleID,
				)
			}
		}
	}

	if group.GetMembershipMode() == jobsgroups.GroupMembershipMode_GROUP_MEMBERSHIP_MODE_STRICT {
		ruleMemberIDs := map[int32]struct{}{}
		for _, match := range ruleMatches {
			ruleMemberIDs[match.UserID] = struct{}{}
		}
		for _, bundle := range bundles {
			if !bundle.hasSource(jobsgroups.GroupMemberSource_GROUP_MEMBER_SOURCE_MANUAL) {
				continue
			}
			if _, ok := ruleMemberIDs[bundle.member.GetUserId()]; ok {
				continue
			}
			bundle.member.IsMember = false
		}
	}

	for _, exclusion := range exclusions {
		bundle := getBundle(exclusion.GetUserId())
		bundle.member.IsExcluded = true
		if includeReasons {
			var detail *string
			if exclusion.GetReason() != "" {
				reason := exclusion.GetReason()
				detail = &reason
			}
			bundle.addReason(
				jobsgroups.GroupMemberSource_GROUP_MEMBER_SOURCE_UNSPECIFIED,
				jobsgroups.GroupMembershipReasonType_GROUP_MEMBERSHIP_REASON_TYPE_EXCLUSION,
				detail,
				nil,
			)
		}
	}

	if wantLeaders {
		for _, leader := range leaders {
			bundle := getBundle(leader.GetUserId())
			bundle.member.IsLeader = true
			bundle.addSource(jobsgroups.GroupMemberSource_GROUP_MEMBER_SOURCE_LEADER)
			if includeReasons {
				bundle.addReason(
					jobsgroups.GroupMemberSource_GROUP_MEMBER_SOURCE_LEADER,
					jobsgroups.GroupMembershipReasonType_GROUP_MEMBERSHIP_REASON_TYPE_LEADER,
					nil,
					nil,
				)
			}
		}
	}

	resolved := make([]*jobsgroups.GroupResolvedMember, 0, len(bundles))
	for _, bundle := range bundles {
		if bundle.member.IsExcluded && !includeExcluded && !bundle.member.IsLeader {
			continue
		}
		if bundle.member.IsExcluded {
			bundle.member.IsMember = false
		}
		bundle.finalize(includeReasons)
		if !bundle.member.IsMember && !bundle.member.IsLeader && !bundle.member.IsExcluded {
			if len(sources) == 0 || !resolvedMemberMatchesSources(bundle.member, sources) {
				continue
			}
		}
		if !resolvedMemberMatchesSources(bundle.member, sources) {
			continue
		}
		resolved = append(resolved, bundle.member)
	}

	sort.SliceStable(resolved, func(i, j int) bool {
		if resolved[i].GetUserId() == resolved[j].GetUserId() {
			return resolved[i].GetGroupId() < resolved[j].GetGroupId()
		}
		return resolved[i].GetUserId() < resolved[j].GetUserId()
	})

	return resolved, nil
}

func (s *Server) ListGroupMembers(
	ctx context.Context,
	req *pbjobs.ListGroupMembersRequest,
) (*pbjobs.ListGroupMembersResponse, error) {
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

	members, err := s.resolveGroupMembers(
		ctx,
		group,
		req.GetSearch(),
		req.GetIncludeExcluded(),
		req.GetIncludeLeaders(),
		req.GetIncludeReasons(),
		req.GetSources(),
	)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(members, func(i, j int) bool {
		return compareResolvedMembersBySort(members[i], members[j], req.GetSort()) < 0
	})

	pag, limit := req.GetPagination().GetResponse(int64(len(members)))
	start := pag.GetOffset()
	end := start + limit
	if start > int64(len(members)) {
		start = int64(len(members))
	}
	if end > int64(len(members)) {
		end = int64(len(members))
	}

	resp := &pbjobs.ListGroupMembersResponse{
		Pagination: pag,
		Members:    members[start:end],
	}
	resp.Pagination.Update(len(resp.GetMembers()))

	targets := appendGroupResolvedMemberColleagueTargets(nil, resp.GetMembers())
	if err := s.hydrateGroupColleagueTargets(ctx, userInfo, targets); err != nil {
		return nil, err
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)
	return resp, nil
}

func paginateGroupItems[T any](
	pagination *database.PaginationRequest,
	items []T,
) (*database.PaginationResponse, []T) {
	pag, limit := pagination.GetResponse(int64(len(items)))
	start := pag.GetOffset()
	end := start + limit
	if start > int64(len(items)) {
		start = int64(len(items))
	}
	if end > int64(len(items)) {
		end = int64(len(items))
	}

	paged := items[start:end]
	pag.Update(len(paged))
	return pag, paged
}

func (s *Server) ListGroupManualMembers(
	ctx context.Context,
	req *pbjobs.ListGroupManualMembersRequest,
) (*pbjobs.ListGroupManualMembersResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.id", req.GetGroupId(),
	})

	if _, err := s.getActiveGroupForJob(
		ctx,
		s.db,
		userInfo.GetJob(),
		req.GetGroupId(),
	); err != nil {
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

	members, err := s.store.ListGroupManualMembers(ctx, s.db, req.GetGroupId(), req.GetSearch())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	pag, paged := paginateGroupItems(req.GetPagination(), members)

	resp := &pbjobs.ListGroupManualMembersResponse{
		Pagination:    pag,
		ManualMembers: paged,
	}
	targets := appendGroupManualMemberColleagueTargets(nil, resp.GetManualMembers())
	if err := s.hydrateGroupColleagueTargets(ctx, userInfo, targets); err != nil {
		return nil, err
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)
	return resp, nil
}

func (s *Server) ListGroupMemberExclusions(
	ctx context.Context,
	req *pbjobs.ListGroupMemberExclusionsRequest,
) (*pbjobs.ListGroupMemberExclusionsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.id", req.GetGroupId(),
	})

	if _, err := s.getActiveGroupForJob(
		ctx,
		s.db,
		userInfo.GetJob(),
		req.GetGroupId(),
	); err != nil {
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

	exclusions, err := s.store.ListGroupMemberExclusions(
		ctx,
		s.db,
		req.GetGroupId(),
		req.GetSearch(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	pag, paged := paginateGroupItems(req.GetPagination(), exclusions)

	resp := &pbjobs.ListGroupMemberExclusionsResponse{
		Pagination: pag,
		Exclusions: paged,
	}
	targets := appendGroupMemberExclusionColleagueTargets(nil, resp.GetExclusions())
	if err := s.hydrateGroupColleagueTargets(ctx, userInfo, targets); err != nil {
		return nil, err
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)
	return resp, nil
}

func (s *Server) ListGroupLeaders(
	ctx context.Context,
	req *pbjobs.ListGroupLeadersRequest,
) (*pbjobs.ListGroupLeadersResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.id", req.GetGroupId(),
	})

	if _, err := s.getActiveGroupForJob(
		ctx,
		s.db,
		userInfo.GetJob(),
		req.GetGroupId(),
	); err != nil {
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

	leaders, err := s.store.ListGroupLeaders(ctx, s.db, req.GetGroupId(), req.GetSearch())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	pag, paged := paginateGroupItems(req.GetPagination(), leaders)

	resp := &pbjobs.ListGroupLeadersResponse{
		Pagination: pag,
		Leaders:    paged,
	}
	targets := appendGroupLeaderColleagueTargets(nil, resp.GetLeaders())
	if err := s.hydrateGroupColleagueTargets(ctx, userInfo, targets); err != nil {
		return nil, err
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)
	return resp, nil
}

func (s *Server) AddGroupMember(
	ctx context.Context,
	req *pbjobs.AddGroupMemberRequest,
) (*pbjobs.AddGroupMemberResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.id", req.GetGroupId(),
		"fivenet.jobs.groups.user_id", req.GetUserId(),
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
	if err := s.ensureGroupUserInJob(ctx, tx, userInfo.GetJob(), req.GetUserId()); err != nil {
		return nil, err
	}
	if group.GetMembershipMode() == jobsgroups.GroupMembershipMode_GROUP_MEMBERSHIP_MODE_STRICT {
		matches, err := s.userMatchesGroupRules(ctx, tx, group, req.GetUserId())
		if err != nil {
			return nil, err
		}
		if !matches {
			return nil, errorsjobs.ErrGroupMemberRulesRequired
		}
	}

	var reason *string
	if req.HasReason() {
		reasonValue := req.GetReason()
		reason = &reasonValue
	}

	member, created, err := s.store.AddGroupManualMember(
		ctx,
		tx,
		req.GetGroupId(),
		req.GetUserId(),
		userInfo.GetUserId(),
		reason,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	activityType := jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_MEMBER_ADDED
	if !created {
		activityType = jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_UPDATED
	}
	if err := s.addGroupActivity(
		ctx,
		tx,
		userInfo.GetJob(),
		req.GetGroupId(),
		activityType,
		userInfo.GetUserId(),
		req.GetUserId(),
		0,
		reason,
		nil,
	); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	updated, err := s.recountAndGetGroup(ctx, tx, userInfo.GetJob(), req.GetGroupId())
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	grpc_audit.AddMeta(ctx, "jobs.group.id", strconv.FormatInt(req.GetGroupId(), 10))
	grpc_audit.AddMeta(ctx, "jobs.group.user_id", strconv.FormatInt(int64(req.GetUserId()), 10))
	if req.HasReason() {
		grpc_audit.AddMeta(ctx, "jobs.group.reason", req.GetReason())
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

	targets := appendGroupColleagueTargets(nil, []*jobsgroups.Group{updated})
	targets = appendGroupManualMemberColleagueTargets(
		targets,
		[]*jobsgroups.GroupManualMember{member},
	)
	if err := s.hydrateGroupColleagueTargets(ctx, userInfo, targets); err != nil {
		return nil, err
	}

	return &pbjobs.AddGroupMemberResponse{
		Member: member,
		Group:  updated,
	}, nil
}

func (s *Server) RemoveGroupMember(
	ctx context.Context,
	req *pbjobs.RemoveGroupMemberRequest,
) (*pbjobs.RemoveGroupMemberResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.id", req.GetGroupId(),
		"fivenet.jobs.groups.user_id", req.GetUserId(),
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

	if err := s.store.RemoveGroupManualMember(
		ctx,
		tx,
		req.GetGroupId(),
		req.GetUserId(),
	); err != nil {
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
		jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_MEMBER_REMOVED,
		userInfo.GetUserId(),
		req.GetUserId(),
		0,
		reason,
		nil,
	); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	updated, err := s.recountAndGetGroup(ctx, tx, userInfo.GetJob(), req.GetGroupId())
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	grpc_audit.AddMeta(ctx, "jobs.group.id", strconv.FormatInt(req.GetGroupId(), 10))
	grpc_audit.AddMeta(ctx, "jobs.group.user_id", strconv.FormatInt(int64(req.GetUserId()), 10))
	if req.HasReason() {
		grpc_audit.AddMeta(ctx, "jobs.group.reason", req.GetReason())
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)

	if err := s.hydrateGroupColleagues(ctx, userInfo, updated); err != nil {
		return nil, err
	}

	return &pbjobs.RemoveGroupMemberResponse{Group: updated}, nil
}

func (s *Server) ExcludeGroupMember(
	ctx context.Context,
	req *pbjobs.ExcludeGroupMemberRequest,
) (*pbjobs.ExcludeGroupMemberResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.id", req.GetGroupId(),
		"fivenet.jobs.groups.user_id", req.GetUserId(),
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
	if err := s.ensureGroupUserInJob(ctx, tx, userInfo.GetJob(), req.GetUserId()); err != nil {
		return nil, err
	}

	var reason *string
	if req.HasReason() {
		reasonValue := req.GetReason()
		reason = &reasonValue
	}

	exclusion, created, err := s.store.AddGroupMemberExclusion(
		ctx,
		tx,
		req.GetGroupId(),
		req.GetUserId(),
		req.GetReasonType(),
		userInfo.GetUserId(),
		reason,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	activityType := jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_MEMBER_EXCLUDED
	if !created {
		activityType = jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_UPDATED
	}
	if err := s.addGroupActivity(
		ctx,
		tx,
		userInfo.GetJob(),
		req.GetGroupId(),
		activityType,
		userInfo.GetUserId(),
		req.GetUserId(),
		0,
		reason,
		nil,
	); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	updated, err := s.recountAndGetGroup(ctx, tx, userInfo.GetJob(), req.GetGroupId())
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	grpc_audit.AddMeta(ctx, "jobs.group.id", strconv.FormatInt(req.GetGroupId(), 10))
	grpc_audit.AddMeta(ctx, "jobs.group.user_id", strconv.FormatInt(int64(req.GetUserId()), 10))
	if req.HasReason() {
		grpc_audit.AddMeta(ctx, "jobs.group.reason", req.GetReason())
	}
	grpc_audit.AddMeta(ctx, "jobs.group.reason_type", req.GetReasonType().String())
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

	targets := appendGroupColleagueTargets(nil, []*jobsgroups.Group{updated})
	targets = appendGroupMemberExclusionColleagueTargets(
		targets,
		[]*jobsgroups.GroupMemberExclusion{exclusion},
	)
	if err := s.hydrateGroupColleagueTargets(ctx, userInfo, targets); err != nil {
		return nil, err
	}

	return &pbjobs.ExcludeGroupMemberResponse{
		Exclusion: exclusion,
		Group:     updated,
	}, nil
}

func (s *Server) RemoveGroupMemberExclusion(
	ctx context.Context,
	req *pbjobs.RemoveGroupMemberExclusionRequest,
) (*pbjobs.RemoveGroupMemberExclusionResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.id", req.GetGroupId(),
		"fivenet.jobs.groups.user_id", req.GetUserId(),
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

	if err := s.store.RemoveGroupMemberExclusion(
		ctx,
		tx,
		req.GetGroupId(),
		req.GetUserId(),
	); err != nil {
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
		jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_MEMBER_EXCLUSION_REMOVED,
		userInfo.GetUserId(),
		req.GetUserId(),
		0,
		reason,
		nil,
	); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	updated, err := s.recountAndGetGroup(ctx, tx, userInfo.GetJob(), req.GetGroupId())
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	grpc_audit.AddMeta(ctx, "jobs.group.id", strconv.FormatInt(req.GetGroupId(), 10))
	grpc_audit.AddMeta(ctx, "jobs.group.user_id", strconv.FormatInt(int64(req.GetUserId()), 10))
	if req.HasReason() {
		grpc_audit.AddMeta(ctx, "jobs.group.reason", req.GetReason())
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)

	if err := s.hydrateGroupColleagues(ctx, userInfo, updated); err != nil {
		return nil, err
	}

	return &pbjobs.RemoveGroupMemberExclusionResponse{Group: updated}, nil
}

func (s *Server) AddGroupLeader(
	ctx context.Context,
	req *pbjobs.AddGroupLeaderRequest,
) (*pbjobs.AddGroupLeaderResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.id", req.GetGroupId(),
		"fivenet.jobs.groups.user_id", req.GetUserId(),
	})

	if err := s.ensureGroupAccess(
		ctx,
		userInfo,
		req.GetGroupId(),
		groupsaccess.AccessLevel_ACCESS_LEVEL_MANAGE,
	); err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	defer tx.Rollback()

	if _, err := s.getActiveGroupForJob(ctx, tx, userInfo.GetJob(), req.GetGroupId()); err != nil {
		return nil, err
	}
	if err := s.ensureGroupUserInJob(ctx, tx, userInfo.GetJob(), req.GetUserId()); err != nil {
		return nil, err
	}

	leader, created, err := s.store.AddGroupLeader(
		ctx,
		tx,
		req.GetGroupId(),
		req.GetUserId(),
		userInfo.GetUserId(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	var reason *string
	if req.HasReason() {
		reasonValue := req.GetReason()
		reason = &reasonValue
	}
	if created {
		if err := s.addGroupActivity(
			ctx,
			tx,
			userInfo.GetJob(),
			req.GetGroupId(),
			jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_LEADER_ADDED,
			userInfo.GetUserId(),
			req.GetUserId(),
			0,
			reason,
			nil,
		); err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	updated, err := s.recountAndGetGroup(ctx, tx, userInfo.GetJob(), req.GetGroupId())
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	grpc_audit.AddMeta(ctx, "jobs.group.id", strconv.FormatInt(req.GetGroupId(), 10))
	grpc_audit.AddMeta(ctx, "jobs.group.user_id", strconv.FormatInt(int64(req.GetUserId()), 10))
	if req.HasReason() {
		grpc_audit.AddMeta(ctx, "jobs.group.reason", req.GetReason())
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

	targets := appendGroupColleagueTargets(nil, []*jobsgroups.Group{updated})
	targets = appendGroupLeaderColleagueTargets(targets, []*jobsgroups.GroupLeader{leader})
	if err := s.hydrateGroupColleagueTargets(ctx, userInfo, targets); err != nil {
		return nil, err
	}

	return &pbjobs.AddGroupLeaderResponse{
		Leader: leader,
		Group:  updated,
	}, nil
}

func (s *Server) RemoveGroupLeader(
	ctx context.Context,
	req *pbjobs.RemoveGroupLeaderRequest,
) (*pbjobs.RemoveGroupLeaderResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.id", req.GetGroupId(),
		"fivenet.jobs.groups.user_id", req.GetUserId(),
	})

	if err := s.ensureGroupAccess(
		ctx,
		userInfo,
		req.GetGroupId(),
		groupsaccess.AccessLevel_ACCESS_LEVEL_MANAGE,
	); err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	defer tx.Rollback()

	if _, err := s.getActiveGroupForJob(ctx, tx, userInfo.GetJob(), req.GetGroupId()); err != nil {
		return nil, err
	}

	if err := s.store.RemoveGroupLeader(ctx, tx, req.GetGroupId(), req.GetUserId()); err != nil {
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
		jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_LEADER_REMOVED,
		userInfo.GetUserId(),
		req.GetUserId(),
		0,
		reason,
		nil,
	); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	updated, err := s.recountAndGetGroup(ctx, tx, userInfo.GetJob(), req.GetGroupId())
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	grpc_audit.AddMeta(ctx, "jobs.group.id", strconv.FormatInt(req.GetGroupId(), 10))
	grpc_audit.AddMeta(ctx, "jobs.group.user_id", strconv.FormatInt(int64(req.GetUserId()), 10))
	if req.HasReason() {
		grpc_audit.AddMeta(ctx, "jobs.group.reason", req.GetReason())
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)

	if err := s.hydrateGroupColleagues(ctx, userInfo, updated); err != nil {
		return nil, err
	}

	return &pbjobs.RemoveGroupLeaderResponse{Group: updated}, nil
}
