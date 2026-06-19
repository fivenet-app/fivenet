package jobs

import (
	"context"
	"strconv"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorsjobs "github.com/fivenet-app/fivenet/v2026/services/jobs/errors"
	jobsstore "github.com/fivenet-app/fivenet/v2026/stores/jobs"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) resolveCurrentJobID(userInfo *userinfo.UserInfo) (int64, error) {
	job := s.enricher.GetJobByName(userInfo.GetJob())
	if job == nil {
		return 0, errorsjobs.ErrFailedQuery
	}

	return job.GetId(), nil
}

func (s *Server) ListGroups(
	ctx context.Context,
	req *pbjobs.ListGroupsRequest,
) (*pbjobs.ListGroupsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	jobID, err := s.resolveCurrentJobID(userInfo)
	if err != nil {
		return nil, err
	}

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.job_id", jobID,
	})

	count, err := s.store.CountGroups(ctx, s.db, jobsstore.GroupsQuery{
		JobID:           jobID,
		States:          req.GetStates(),
		Search:          req.GetSearch(),
		IncludeCounts:   req.GetIncludeCounts(),
		IncludeInactive: req.GetIncludeInactive(),
		IncludeArchived: req.GetIncludeArchived(),
		Sort:            req.GetSort(),
		Offset:          req.GetPagination().GetOffset(),
	})
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
		JobID:           jobID,
		States:          req.GetStates(),
		Search:          req.GetSearch(),
		IncludeCounts:   req.GetIncludeCounts(),
		IncludeInactive: req.GetIncludeInactive(),
		IncludeArchived: req.GetIncludeArchived(),
		Sort:            req.GetSort(),
		Offset:          req.GetPagination().GetOffset(),
		Limit:           limit,
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	resp.Groups = groups

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)
	return resp, nil
}

func (s *Server) GetGroup(
	ctx context.Context,
	req *pbjobs.GetGroupRequest,
) (*pbjobs.GetGroupResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	jobID, err := s.resolveCurrentJobID(userInfo)
	if err != nil {
		return nil, err
	}

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.id", req.GetId(),
		"fivenet.jobs.groups.job_id", jobID,
	})

	group, err := s.store.GetGroup(ctx, s.db, jobsstore.GroupQuery{
		JobID:           jobID,
		IncludeArchived: req.GetIncludeArchived(),
	}, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if group == nil {
		return nil, errorsjobs.ErrNotFoundOrNoPerms
	}

	resp := &pbjobs.GetGroupResponse{Group: group}
	if req.GetIncludeRules() {
		resp.Rules = []*jobsgroups.GroupRule{}
	}
	if req.GetIncludeLeaders() {
		resp.Leaders = []*jobsgroups.GroupLeader{}
	}
	if req.GetIncludeManualMembers() {
		resp.ManualMembers = []*jobsgroups.GroupManualMember{}
	}
	if req.GetIncludeExclusions() {
		resp.Exclusions = []*jobsgroups.GroupMemberExclusion{}
	}
	if req.GetIncludeResolvedMembers() {
		resp.ResolvedMembers = []*jobsgroups.GroupResolvedMember{}
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)
	return resp, nil
}

func (s *Server) CreateGroup(
	ctx context.Context,
	req *pbjobs.CreateGroupRequest,
) (*pbjobs.CreateGroupResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	jobID, err := s.resolveCurrentJobID(userInfo)
	if err != nil {
		return nil, err
	}

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.job_id", jobID,
		"fivenet.jobs.groups.name", req.GetName(),
	})

	group := &jobsgroups.Group{
		JobId:           jobID,
		Name:            req.GetName(),
		Type:            req.GetType(),
		MembershipMode:  req.GetMembershipMode(),
		SortOrder:       req.GetSortOrder(),
		CreatedByUserId: int64(userInfo.GetUserId()),
		UpdatedByUserId: int64(userInfo.GetUserId()),
	}
	if req.HasDescription() {
		description := req.GetDescription()
		group.Description = &description
	}
	if req.HasShortName() {
		shortName := req.GetShortName()
		group.ShortName = &shortName
	}
	if req.HasLogoFileId() {
		logoFileID := req.GetLogoFileId()
		group.LogoFileId = &logoFileID
	}
	if req.HasColor() {
		color := req.GetColor()
		group.Color = &color
	}

	id, err := s.store.CreateGroup(ctx, s.db, group)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	grpc_audit.AddMeta(ctx, "jobs.group.id", strconv.FormatInt(id, 10))
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

	created, err := s.store.GetGroup(ctx, s.db, jobsstore.GroupQuery{
		JobID:           jobID,
		IncludeArchived: true,
	}, id)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if created == nil {
		return nil, errorsjobs.ErrFailedQuery
	}

	return &pbjobs.CreateGroupResponse{Group: created}, nil
}

func (s *Server) UpdateGroup(
	ctx context.Context,
	req *pbjobs.UpdateGroupRequest,
) (*pbjobs.UpdateGroupResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	jobID, err := s.resolveCurrentJobID(userInfo)
	if err != nil {
		return nil, err
	}

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.id", req.GetId(),
		"fivenet.jobs.groups.job_id", jobID,
	})

	group, err := s.store.GetGroup(ctx, s.db, jobsstore.GroupQuery{
		JobID:           jobID,
		IncludeArchived: false,
	}, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if group == nil {
		return nil, errorsjobs.ErrNotFoundOrNoPerms
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
	if req.HasLogoFileId() {
		logoFileID := req.GetLogoFileId()
		group.LogoFileId = &logoFileID
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
	if req.HasSortOrder() {
		group.SortOrder = req.GetSortOrder()
	}
	if req.HasState() {
		if req.GetState() == jobsgroups.GroupState_GROUP_STATE_ARCHIVED {
			return nil, errorsjobs.ErrNotFoundOrNoPerms
		}
		if req.GetState() != jobsgroups.GroupState_GROUP_STATE_UNSPECIFIED {
			group.State = req.GetState()
		}
	}

	group.UpdatedByUserId = int64(userInfo.GetUserId())

	if err := s.store.UpdateGroup(ctx, s.db, group); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	updated, err := s.store.GetGroup(ctx, s.db, jobsstore.GroupQuery{
		JobID:           jobID,
		IncludeArchived: true,
	}, group.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if updated == nil {
		return nil, errorsjobs.ErrFailedQuery
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	return &pbjobs.UpdateGroupResponse{Group: updated}, nil
}

func (s *Server) ArchiveGroup(
	ctx context.Context,
	req *pbjobs.ArchiveGroupRequest,
) (*pbjobs.ArchiveGroupResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	jobID, err := s.resolveCurrentJobID(userInfo)
	if err != nil {
		return nil, err
	}

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.id", req.GetId(),
		"fivenet.jobs.groups.job_id", jobID,
	})

	if err := s.store.ArchiveGroup(
		ctx,
		s.db,
		jobID,
		req.GetId(),
		int64(userInfo.GetUserId()),
	); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	group, err := s.store.GetGroup(ctx, s.db, jobsstore.GroupQuery{
		JobID:           jobID,
		IncludeArchived: true,
	}, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if group == nil {
		return nil, errorsjobs.ErrNotFoundOrNoPerms
	}

	if req.HasReason() {
		grpc_audit.AddMeta(ctx, "jobs.group.reason", req.GetReason())
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	return &pbjobs.ArchiveGroupResponse{Group: group}, nil
}

func (s *Server) RestoreGroup(
	ctx context.Context,
	req *pbjobs.RestoreGroupRequest,
) (*pbjobs.RestoreGroupResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	jobID, err := s.resolveCurrentJobID(userInfo)
	if err != nil {
		return nil, err
	}

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.id", req.GetId(),
		"fivenet.jobs.groups.job_id", jobID,
	})

	if err := s.store.RestoreGroup(
		ctx,
		s.db,
		jobID,
		req.GetId(),
		int64(userInfo.GetUserId()),
	); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	group, err := s.store.GetGroup(ctx, s.db, jobsstore.GroupQuery{
		JobID:           jobID,
		IncludeArchived: true,
	}, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if group == nil {
		return nil, errorsjobs.ErrNotFoundOrNoPerms
	}

	if req.HasReason() {
		grpc_audit.AddMeta(ctx, "jobs.group.reason", req.GetReason())
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_RESTORED)
	return &pbjobs.RestoreGroupResponse{Group: group}, nil
}
