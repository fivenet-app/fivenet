package citizens

import (
	context "context"

	pbaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	citizenslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/citizens/labels"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbcitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens"
	permscitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorscitizens "github.com/fivenet-app/fivenet/v2026/services/citizens/errors"
)

var (
	tCitizensLabelsJob = table.FivenetUserLabelsJob
	tCitizenLabels     = table.FivenetUserLabels
)

var labelSubjectAccessOptions = access.SubjectAccessOptions{
	BlockedAccess: -1,
}

func (s *Server) ListLabels(
	ctx context.Context,
	req *pbcitizens.ListLabelsRequest,
) (*pbcitizens.ListLabelsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Either user can see labels of citizens
	fields, err := permscitizens.CitizensService.ListCitizens.FieldsTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	var canCreateLabel bool
	if !fields.Contains(permscitizens.CitizensServiceListCitizensFieldsPermValueUserPropsLabels) {
		canCreateLabel = s.ps.Can(userInfo, permscitizens.LabelsService.CreateOrUpdateLabel.Perm)
		if !canCreateLabel {
			return nil, errorscitizens.ErrLabelNotFound
		}
	}

	labels, err := s.store.ListLabels(
		ctx,
		s.db,
		userInfo,
		req.GetSearch(),
		req.GetOwnJobOnly(),
		canCreateLabel,
		int32(min(req.GetMinAccess(), citizenslabels.AccessLevel_ACCESS_LEVEL_VIEW)),
		userInfo.GetJobAdmin(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if err := s.fillLabelAccess(ctx, labels.GetList()...); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	return &pbcitizens.ListLabelsResponse{Labels: labels.GetList()}, nil
}

func (s *Server) GetLabel(
	ctx context.Context,
	req *pbcitizens.GetLabelRequest,
) (*pbcitizens.GetLabelResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.labelsAccess.CanUserAccessTarget(
		ctx,
		req.GetId(),
		userInfo,
		int32(citizenslabels.AccessLevel_ACCESS_LEVEL_VIEW),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	if !check {
		return nil, errorscitizens.ErrLabelNotFound
	}

	label, err := s.store.GetLabel(
		ctx,
		s.db,
		userInfo.GetJob(),
		req.GetId(),
		userInfo.GetJobAdmin(),
	)
	if err != nil {
		return nil, err
	}

	if label == nil || label.GetId() == 0 {
		return nil, errorscitizens.ErrLabelNotFound
	}

	if !userInfo.GetJobAdmin() {
		if label.GetDeletedAt() != nil {
			return nil, errorscitizens.ErrLabelNotFound
		}

		label.DeletedAt = nil
	}

	if err := s.fillLabelAccess(ctx, label); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	return &pbcitizens.GetLabelResponse{
		Label: label,
	}, nil
}

func (s *Server) CreateOrUpdateLabel(
	ctx context.Context,
	req *pbcitizens.CreateOrUpdateLabelRequest,
) (*pbcitizens.CreateOrUpdateLabelResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	label := req.GetLabel()
	label.Job = &userInfo.Job

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	defer tx.Rollback()

	if req.GetLabel().GetId() > 0 {
		if err := s.store.UpdateLabel(ctx, tx, label, userInfo.GetJob()); err != nil {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	} else {
		sortOrder, err := s.store.NextLabelSortOrder(ctx, tx, userInfo.GetJob())
		if err != nil {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
		label.SortOrder = sortOrder

		lastId, err := s.store.InsertLabel(ctx, tx, label)
		if err != nil {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}

		label.SetId(lastId)

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)
	}

	fallbackAccess := &pbaccess.Access{
		Jobs: []*pbaccess.JobAccess{
			{
				TargetId:     label.GetId(),
				Job:          userInfo.GetJob(),
				MinimumGrade: userInfo.GetJobGrade(),
				Access:       int32(citizenslabels.AccessLevel_ACCESS_LEVEL_REMOVE),
			},
		},
	}
	normalizedAccess, err := access.NormalizeAccess(
		label.GetAccess(),
		nil,
		fallbackAccess,
		15,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if _, err := s.labelsAccess.ReplaceTargetAccess(
		ctx,
		tx,
		s.labelsAccessResolver,
		label.GetId(),
		normalizedAccess,
		labelSubjectAccessOptions,
	); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	label, err = s.store.GetLabel(ctx, s.db, label.GetJob(), label.GetId(), userInfo.GetJobAdmin())
	if err != nil {
		return nil, err
	}
	if !userInfo.GetJobAdmin() {
		label.DeletedAt = nil
	}

	// Retrieve labels access
	if err := s.fillLabelAccess(ctx, label); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	return &pbcitizens.CreateOrUpdateLabelResponse{
		Label: label,
	}, nil
}

func (s *Server) DeleteLabel(
	ctx context.Context,
	req *pbcitizens.DeleteLabelRequest,
) (*pbcitizens.DeleteLabelResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	label, err := s.store.GetLabel(
		ctx,
		s.db,
		userInfo.GetJob(),
		req.GetId(),
		userInfo.GetJobAdmin(),
	)
	if err != nil {
		return nil, err
	}
	if label == nil || label.GetJob() != userInfo.GetJob() {
		return nil, errorscitizens.ErrFailedQuery
	}

	var deletedAtTime *timestamp.Timestamp
	if label.GetDeletedAt() == nil || !userInfo.GetJobAdmin() {
		deletedAtTime = timestamp.Now()
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	} else {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_RESTORED)
	}

	if err := s.store.DeleteLabel(
		ctx,
		s.db,
		userInfo.GetJob(),
		req.GetId(),
		deletedAtTime,
	); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	return &pbcitizens.DeleteLabelResponse{}, nil
}

func (s *Server) fillLabelAccess(ctx context.Context, labels ...*citizenslabels.Label) error {
	for _, label := range labels {
		if label == nil || label.GetId() == 0 {
			continue
		}
		access, err := s.labelsAccess.ListTargetAccess(
			ctx,
			s.db,
			label.GetId(),
			labelSubjectAccessOptions,
		)
		if err != nil {
			return err
		}
		label.Access = &pbaccess.Access{
			Jobs: access.GetJobs(),
		}
	}
	return nil
}

func (s *Server) ReorderLabels(
	ctx context.Context,
	req *pbcitizens.ReorderLabelsRequest,
) (*pbcitizens.ReorderLabelsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	labelIds := utils.RemoveSliceDuplicates(req.GetLabelIds())
	if err := s.store.ReorderLabels(ctx, userInfo.GetJob(), labelIds); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return &pbcitizens.ReorderLabelsResponse{}, nil
}
