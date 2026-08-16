package jobs

import (
	context "context"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	jobslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/labels"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs"
	permsjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	errorscitizens "github.com/fivenet-app/fivenet/v2026/services/citizens/errors"
	errorsjobs "github.com/fivenet-app/fivenet/v2026/services/jobs/errors"
)

func (s *Server) GetColleagueLabels(
	ctx context.Context,
	req *pbjobs.GetColleagueLabelsRequest,
) (*pbjobs.GetColleagueLabelsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	resp := &pbjobs.GetColleagueLabelsResponse{
		Labels: []*jobslabels.Label{},
	}

	// Fields Permission Check
	fields, err := permsjobs.ColleaguesService.GetColleague.TypesTyped.Get(s.perms, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if userInfo.GetJobAdmin() {
		fields.Set(permsjobs.ColleaguesServiceGetColleagueTypesPermValueLabels)
	}
	if !fields.Contains(permsjobs.ColleaguesServiceGetColleagueTypesPermValueLabels) {
		// Fallback to checking if user has manage colleague labels permission
		if !s.perms.Can(userInfo, permsjobs.ColleaguesService.CreateOrUpdateLabel.Perm) {
			return nil, errorsjobs.ErrLabelsNoPerms
		}
	}

	labels, err := s.store.GetColleagueLabels(
		ctx,
		s.db,
		userInfo.GetJob(),
		req.GetSearch(),
		userInfo.GetJobAdmin(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	resp.Labels = labels

	return resp, nil
}

func (s *Server) CreateOrUpdateLabel(
	ctx context.Context,
	req *pbjobs.CreateOrUpdateLabelRequest,
) (*pbjobs.CreateOrUpdateLabelResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	label := req.GetLabel()
	if label == nil {
		return nil, errorsjobs.ErrFailedQuery
	}
	label.Job = &userInfo.Job

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	defer tx.Rollback()

	if label.GetId() > 0 {
		existing, err := s.store.GetLabel(ctx, tx, userInfo.GetJob(), label.GetId(), false)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		if existing == nil {
			return nil, errorsjobs.ErrLabelNotFound
		}

		if err := s.store.UpdateLabel(ctx, tx, label, userInfo.GetJob()); err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	} else {
		sortOrder, err := s.store.NextLabelSortOrder(ctx, tx, userInfo.GetJob())
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		label.SortOrder = sortOrder

		lastId, err := s.store.InsertLabel(ctx, tx, label)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}

		label.SetId(lastId)

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	return &pbjobs.CreateOrUpdateLabelResponse{Label: label}, nil
}

func (s *Server) DeleteLabel(
	ctx context.Context,
	req *pbjobs.DeleteLabelRequest,
) (*pbjobs.DeleteLabelResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	label, err := s.store.GetLabel(ctx, s.db, userInfo.GetJob(), req.GetId(), false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if label == nil || label.GetJob() != userInfo.GetJob() {
		return nil, errorsjobs.ErrLabelNotFound
	}

	if err := s.store.DeleteLabel(
		ctx,
		s.db,
		userInfo.GetJob(),
		req.GetId(),
		timestamp.Now(),
	); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)

	return &pbjobs.DeleteLabelResponse{}, nil
}

func (s *Server) ReorderLabels(
	ctx context.Context,
	req *pbjobs.ReorderLabelsRequest,
) (*pbjobs.ReorderLabelsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	labelIds := utils.SliceDedup(req.GetLabelIds())
	if err := s.store.ReorderLabels(ctx, userInfo.GetJob(), labelIds); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return &pbjobs.ReorderLabelsResponse{}, nil
}

func (s *Server) GetColleagueLabelsStats(
	ctx context.Context,
	req *pbjobs.GetColleagueLabelsStatsRequest,
) (*pbjobs.GetColleagueLabelsStatsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Types Permission Check
	fields, err := permsjobs.ColleaguesService.GetColleague.TypesTyped.Get(s.perms, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if userInfo.GetJobAdmin() {
		fields.Set(permsjobs.ColleaguesServiceGetColleagueTypesPermValueLabels)
	}
	if !fields.Contains(permsjobs.ColleaguesServiceGetColleagueTypesPermValueLabels) {
		return &pbjobs.GetColleagueLabelsStatsResponse{}, nil
	}

	dest, err := s.store.GetColleagueLabelsStats(ctx, s.db, userInfo.GetJob())
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	return &pbjobs.GetColleagueLabelsStatsResponse{
		Count: dest,
	}, nil
}
