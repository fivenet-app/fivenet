package citizens

import (
	context "context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	users "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	pbcitizens "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/citizens"
	permscompletor "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/completor/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorscitizens "github.com/fivenet-app/fivenet/v2025/services/citizens/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tCitizensLabelsJob = table.FivenetUserLabelsJob.AS("label")
	tUserLabels        = table.FivenetUserLabels
)

func (s *Server) ManageLabels(
	ctx context.Context,
	req *pbcitizens.ManageLabelsRequest,
) (*pbcitizens.ManageLabelsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcitizens.CitizensService_ServiceDesc.ServiceName,
		Method:  "ManageLabels",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	resp := &pbcitizens.ManageLabelsResponse{
		Labels: []*users.Label{},
	}

	stmt := tCitizensLabelsJob.
		SELECT(
			tCitizensLabelsJob.ID,
			tCitizensLabelsJob.Job,
			tCitizensLabelsJob.Name,
			tCitizensLabelsJob.Color,
		).
		FROM(tCitizensLabelsJob).
		WHERE(
			tCitizensLabelsJob.Job.EQ(jet.String(userInfo.GetJob())),
		)

	if err := stmt.QueryContext(ctx, s.db, &resp.Labels); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
	}

	_, removed := utils.SlicesDifferenceFunc(resp.GetLabels(), req.GetLabels(),
		func(in *users.Label) string {
			return in.GetName()
		})

	for i := range req.GetLabels() {
		req.Labels[i].Job = &userInfo.Job
	}

	tCitizensLabelsJob := table.FivenetUserLabelsJob

	if len(req.GetLabels()) > 0 {
		toCreate := []*users.Label{}
		toUpdate := []*users.Label{}

		for _, attribute := range req.GetLabels() {
			if attribute.GetId() == 0 {
				toCreate = append(toCreate, attribute)
			} else {
				toUpdate = append(toUpdate, attribute)
			}
		}

		if len(toCreate) > 0 {
			insertStmt := tCitizensLabelsJob.
				INSERT(
					tCitizensLabelsJob.Job,
					tCitizensLabelsJob.Name,
					tCitizensLabelsJob.Color,
				).
				MODELS(toCreate).
				ON_DUPLICATE_KEY_UPDATE(
					tCitizensLabelsJob.Name.SET(jet.StringExp(jet.Raw("VALUES(`name`)"))),
					tCitizensLabelsJob.Color.SET(jet.StringExp(jet.Raw("VALUES(`color`)"))),
				)

			if _, err := insertStmt.ExecContext(ctx, s.db); err != nil {
				return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
			}
		}

		if len(toUpdate) > 0 {
			for _, attribute := range toUpdate {
				updateStmt := tCitizensLabelsJob.
					UPDATE(
						tCitizensLabelsJob.Name,
						tCitizensLabelsJob.Color,
					).
					SET(
						tCitizensLabelsJob.Name.SET(jet.String(attribute.GetName())),
						tCitizensLabelsJob.Color.SET(jet.String(attribute.GetColor())),
					).
					WHERE(jet.AND(
						tCitizensLabelsJob.ID.EQ(jet.Int64(attribute.GetId())),
						tCitizensLabelsJob.Job.EQ(jet.String(attribute.GetJob())),
					))

				if _, err := updateStmt.ExecContext(ctx, s.db); err != nil {
					return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
				}
			}
		}
	}

	if len(removed) > 0 {
		ids := make([]jet.Expression, len(removed))

		for i := range removed {
			ids[i] = jet.Int64(removed[i].GetId())
		}

		deleteStmt := tCitizensLabelsJob.
			DELETE().
			WHERE(jet.AND(
				tCitizensLabelsJob.ID.IN(ids...),
				tCitizensLabelsJob.Job.EQ(jet.String(userInfo.GetJob())),
			)).
			LIMIT(int64(len(removed)))

		if _, err := deleteStmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
	}

	resp.Labels = []*users.Label{}
	if err := stmt.QueryContext(ctx, s.db, &resp.Labels); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return resp, nil
}

func (s *Server) validateLabels(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	attributes []*users.Label,
) (bool, error) {
	if len(attributes) == 0 {
		return true, nil
	}

	jobs, err := s.ps.AttrStringList(
		userInfo,
		permscompletor.CompletorServicePerm,
		permscompletor.CompletorServiceCompleteCitizenLabelsPerm,
		permscompletor.CompletorServiceCompleteCitizenLabelsJobsPermField,
	)
	if err != nil {
		return false, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if jobs.Len() == 0 {
		jobs.Strings = append(jobs.Strings, userInfo.GetJob())
	}

	jobsExp := make([]jet.Expression, len(jobs.GetStrings()))
	for i := range jobs.GetStrings() {
		jobsExp[i] = jet.String(jobs.GetStrings()[i])
	}

	idsExp := make([]jet.Expression, len(attributes))
	for i := range attributes {
		idsExp[i] = jet.Int64(attributes[i].GetId())
	}

	stmt := tCitizensLabelsJob.
		SELECT(
			jet.COUNT(tCitizensLabelsJob.ID).AS("data_count.total"),
		).
		FROM(tCitizensLabelsJob).
		WHERE(jet.AND(
			tCitizensLabelsJob.Job.IN(jobsExp...),
			tCitizensLabelsJob.ID.IN(idsExp...),
		)).
		LIMIT(10)

	var count database.DataCount
	if err := stmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return false, err
		}
	}

	return len(attributes) == int(count.Total), nil
}

func (s *Server) getUserLabels(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	userId int32,
) (*users.Labels, error) {
	jobs, err := s.ps.AttrStringList(
		userInfo,
		permscompletor.CompletorServicePerm,
		permscompletor.CompletorServiceCompleteCitizenLabelsPerm,
		permscompletor.CompletorServiceCompleteCitizenLabelsJobsPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if jobs.Len() == 0 {
		jobs.Strings = append(jobs.Strings, userInfo.GetJob())
	}

	jobsExp := make([]jet.Expression, jobs.Len())
	for i := range jobs.GetStrings() {
		jobsExp[i] = jet.String(jobs.GetStrings()[i])
	}

	stmt := tUserLabels.
		SELECT(
			tCitizensLabelsJob.ID,
			tCitizensLabelsJob.Job,
			tCitizensLabelsJob.Name,
			tCitizensLabelsJob.Color,
		).
		FROM(
			tUserLabels.
				INNER_JOIN(tCitizensLabelsJob,
					tCitizensLabelsJob.ID.EQ(tUserLabels.LabelID),
				),
		).
		WHERE(jet.AND(
			tUserLabels.UserID.EQ(jet.Int32(userId)),
			tCitizensLabelsJob.Job.IN(jobsExp...),
		)).
		ORDER_BY(
			tCitizensLabelsJob.SortKey.ASC(),
		)

	list := &users.Labels{
		List: []*users.Label{},
	}
	if err := stmt.QueryContext(ctx, s.db, &list.List); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return list, nil
}
