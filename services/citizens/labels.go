package citizens

import (
	context "context"

	accessProto "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	citizenslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/citizens/labels"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbcitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens"
	permscitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorscitizens "github.com/fivenet-app/fivenet/v2026/services/citizens/errors"
	"github.com/go-jet/jet/v2/mysql"
)

var tCitizensLabelsJob = table.FivenetUserLabelsJob.AS("label")
var tCitizenLabels = table.FivenetUserLabels

var labelSubjectAccessOptions = access.SubjectAccessOptions{BlockedAccess: -1}

func labelJobAccess(jobs []*citizenslabels.JobAccess) *accessProto.Access {
	return &accessProto.Access{Jobs: jobs}
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

	condition := mysql.Bool(true)

	if !userInfo.GetSuperuser() {
		condition = condition.AND(tCitizensLabelsJob.DeletedAt.IS_NULL())
	}

	if search := dbutils.PrepareForLikeSearch(req.GetSearch()); search != "" {
		condition = condition.AND(tCitizensLabelsJob.Name.LIKE(mysql.String(search)))
	}

	// When an user can create/update labels, the user is allowed to be returned all of their job's labels.
	if req.GetOwnJobOnly() && canCreateLabel {
		condition = condition.AND(tCitizensLabelsJob.Job.EQ(mysql.String(userInfo.GetJob())))
	} else if !userInfo.GetSuperuser() {
		minAccess := min(req.GetMinAccess(), citizenslabels.AccessLevel_ACCESS_LEVEL_VIEW)

		jobAccessExists := s.labelsAccess.ACLAccessExistsCondition(
			tCitizensLabelsJob.ID,
			userInfo,
			int32(minAccess),
		)

		condition = mysql.AND(
			condition,
			jobAccessExists,
		)
	}

	resp, err := s.store.ListLabels(ctx, s.db, condition, userInfo.GetSuperuser())
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if err := s.fillLabelAccess(ctx, resp.GetList()...); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	return &pbcitizens.ListLabelsResponse{Labels: resp.GetList()}, nil
}

func (s *Server) GetLabel(
	ctx context.Context,
	req *pbcitizens.GetLabelRequest,
) (*pbcitizens.GetLabelResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	label, err := s.store.GetLabel(ctx, s.db, userInfo.GetJob(), req.GetId())
	if err != nil {
		return nil, err
	}

	if label == nil || label.GetId() == 0 {
		return nil, errorscitizens.ErrLabelNotFound
	}

	if !userInfo.GetSuperuser() {
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

		lastInsertId, err := s.store.InsertLabel(ctx, tx, label)
		if err != nil {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}

		label.Id = lastInsertId

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)
	}

	access := label.GetAccess()
	if access == nil || len(access.GetJobs()) == 0 {
		access = &citizenslabels.LabelAccess{
			Jobs: []*citizenslabels.JobAccess{
				{
					TargetId:     label.GetId(),
					Job:          userInfo.GetJob(),
					MinimumGrade: userInfo.GetJobGrade(),
					Access:       int32(citizenslabels.AccessLevel_ACCESS_LEVEL_REMOVE),
				},
			},
		}
	}
	if _, err := s.labelsAccess.ReplaceTargetAccess(
		ctx,
		tx,
		s.labelsAccessResolver,
		label.GetId(),
		labelJobAccess(access.GetJobs()),
		labelSubjectAccessOptions,
	); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	label, err = s.store.GetLabel(ctx, s.db, label.GetJob(), label.GetId())
	if err != nil {
		return nil, err
	}
	if !userInfo.GetSuperuser() {
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

	label, err := s.store.GetLabel(ctx, s.db, userInfo.GetJob(), req.GetId())
	if err != nil {
		return nil, err
	}
	if label == nil || label.GetJob() != userInfo.GetJob() {
		return nil, errorscitizens.ErrFailedQuery
	}

	var deletedAtTime *timestamp.Timestamp
	if label.GetDeletedAt() == nil || !userInfo.GetSuperuser() {
		deletedAtTime = timestamp.Now()
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	} else {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_RESTORED)
	}

	if err := s.store.DeleteLabel(ctx, s.db, userInfo.GetJob(), req.GetId(), deletedAtTime); err != nil {
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
		label.Access = &citizenslabels.LabelAccess{Jobs: access.GetJobs()}
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
