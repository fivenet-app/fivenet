package usersel

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"

	jobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	groupsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups/access"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	jobsstore "github.com/fivenet-app/fivenet/v2026/stores/jobs"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
)

const MaxResolvedUsers = 5000

type Resolver struct {
	store       jobsstore.IGroupsQuery
	groupAccess access.JobGroupsAccess
}

type IResolver interface {
	Resolve(
		ctx context.Context,
		db qrm.DB,
		userInfo *pbuserinfo.UserInfo,
		selector *jobs.UserSelector,
		opts ResolveOpts,
	) ([]int32, error)
}

type Params struct {
	fx.In

	Store       jobsstore.IGroupsQuery
	GroupAccess *access.JobGroupsObjectAccess
}

func New(p Params) IResolver {
	return &Resolver{
		store:       p.Store,
		groupAccess: p.GroupAccess,
	}
}

func NewWithAccess(
	store jobsstore.IGroupsQuery,
	groupAccess access.JobGroupsAccess,
) IResolver {
	return &Resolver{
		store:       store,
		groupAccess: groupAccess,
	}
}

type ResolveOpts struct {
	MaxResolvedUsers int
}

func (r *Resolver) Resolve(
	ctx context.Context,
	db qrm.DB,
	userInfo *pbuserinfo.UserInfo,
	selector *jobs.UserSelector,
	opts ResolveOpts,
) ([]int32, error) {
	if selector == nil {
		selector = &jobs.UserSelector{}
	}
	if opts.MaxResolvedUsers <= 0 {
		opts.MaxResolvedUsers = MaxResolvedUsers
	}

	hasExplicit := len(selector.GetUserIds()) > 0
	hasGroups := selector.GetGroups() != nil && len(selector.GetGroups().GetGroupIds()) > 0

	if !hasExplicit && !hasGroups {
		return nil, nil
	}

	job := userInfo.GetJob()
	seen := map[int32]struct{}{}

	for _, uid := range selector.GetUserIds() {
		seen[uid] = struct{}{}
	}

	if groupSel := selector.GetGroups(); groupSel != nil && len(groupSel.GetGroupIds()) > 0 {
		filteredGroupSel, err := r.filterAccessibleGroups(ctx, userInfo, groupSel)
		if err != nil {
			return nil, err
		}
		if len(filteredGroupSel.GetGroupIds()) > 0 {
			groupQuery := jobsstore.GroupQuery{
				Job:             job,
				IncludeArchived: false,
			}

			for _, gid := range filteredGroupSel.GetGroupIds() {
				group, err := r.store.GetGroup(ctx, db, groupQuery, gid)
				if err != nil {
					if errors.Is(err, sql.ErrNoRows) {
						continue
					}
					return nil, fmt.Errorf("failed to resolve group %d. %w", gid, err)
				}

				members, err := r.resolveGroupMembers(ctx, db, group, filteredGroupSel)
				if err != nil {
					return nil, fmt.Errorf("failed to resolve group %d members. %w", gid, err)
				}
				for _, uid := range members {
					seen[uid] = struct{}{}
				}
			}
		}
	}

	result := make([]int32, 0, len(seen))
	for uid := range seen {
		result = append(result, uid)
	}

	if opts.MaxResolvedUsers > 0 && len(result) > opts.MaxResolvedUsers {
		return nil, fmt.Errorf(
			"too many resolved users: %d (max %d)",
			len(result),
			opts.MaxResolvedUsers,
		)
	}

	slices.Sort(result)
	return result, nil
}

func (r *Resolver) filterAccessibleGroups(
	ctx context.Context,
	userInfo *pbuserinfo.UserInfo,
	groupSel *jobs.GroupUserSelector,
) (*jobs.GroupUserSelector, error) {
	if groupSel == nil || len(groupSel.GetGroupIds()) == 0 {
		return groupSel, nil
	}

	allowedGroupIDs, err := r.groupAccess.CanUserAccessTargetIDs(
		ctx,
		userInfo,
		int32(groupsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		groupSel.GetGroupIds()...,
	)
	if err != nil {
		return nil, err
	}

	allowed := make(map[int64]struct{}, len(allowedGroupIDs))
	for _, groupID := range allowedGroupIDs {
		allowed[groupID] = struct{}{}
	}

	filteredGroupIDs := make([]int64, 0, len(groupSel.GetGroupIds()))
	for _, groupID := range groupSel.GetGroupIds() {
		if _, ok := allowed[groupID]; ok {
			filteredGroupIDs = append(filteredGroupIDs, groupID)
		}
	}

	return &jobs.GroupUserSelector{
		GroupIds:        filteredGroupIDs,
		IncludeLeaders:  groupSel.GetIncludeLeaders(),
		IncludeExcluded: groupSel.GetIncludeExcluded(),
	}, nil
}

func (r *Resolver) resolveGroupMembers(
	ctx context.Context,
	db qrm.DB,
	group *jobsgroups.Group,
	groupSel *jobs.GroupUserSelector,
) ([]int32, error) {
	manualMembers, err := r.store.ListGroupManualMembers(ctx, db, group.GetId(), "")
	if err != nil {
		return nil, err
	}
	ruleMatches, err := r.store.ListGroupRuleMemberMatches(ctx, db, group, "")
	if err != nil {
		return nil, err
	}

	excluded := map[int32]struct{}{}
	if !groupSel.GetIncludeExcluded() {
		exclusions, err := r.store.ListGroupMemberExclusions(ctx, db, group.GetId(), "")
		if err != nil {
			return nil, err
		}
		for _, exclusion := range exclusions {
			excluded[exclusion.GetUserId()] = struct{}{}
		}
	}

	members := map[int32]struct{}{}
	ruleMemberIDs := map[int32]struct{}{}
	for _, match := range ruleMatches {
		ruleMemberIDs[match.UserID] = struct{}{}
		if _, ok := excluded[match.UserID]; !ok {
			members[match.UserID] = struct{}{}
		}
	}
	for _, member := range manualMembers {
		if _, ok := excluded[member.GetUserId()]; ok {
			continue
		}
		if group.GetMembershipMode() == jobsgroups.GroupMembershipMode_GROUP_MEMBERSHIP_MODE_STRICT {
			if _, ok := ruleMemberIDs[member.GetUserId()]; !ok {
				continue
			}
		}
		members[member.GetUserId()] = struct{}{}
	}

	if groupSel.GetIncludeLeaders() {
		leaders, err := r.store.ListGroupLeaders(ctx, db, group.GetId(), "")
		if err != nil {
			return nil, err
		}
		for _, leader := range leaders {
			members[leader.GetUserId()] = struct{}{}
		}
	}

	result := make([]int32, 0, len(members))
	for uid := range members {
		result = append(result, uid)
	}

	return result, nil
}

func HasSelection(selector *jobs.UserSelector) bool {
	if selector == nil {
		return false
	}
	if len(selector.GetUserIds()) > 0 {
		return true
	}
	return selector.GetGroups() != nil && len(selector.GetGroups().GetGroupIds()) > 0
}

func GroupsOnly(selector *jobs.UserSelector) *jobs.UserSelector {
	if selector == nil {
		return &jobs.UserSelector{}
	}

	groupSel := selector.GetGroups()
	if groupSel == nil {
		return &jobs.UserSelector{}
	}

	return &jobs.UserSelector{
		Groups: &jobs.GroupUserSelector{
			GroupIds:        slices.Clone(groupSel.GetGroupIds()),
			IncludeLeaders:  groupSel.GetIncludeLeaders(),
			IncludeExcluded: groupSel.GetIncludeExcluded(),
		},
	}
}
