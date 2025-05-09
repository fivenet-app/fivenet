package citizenstore

import (
	context "context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/rector"
	users "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	pbcitizenstore "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/citizenstore"
	permscompletor "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/completor/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorscitizenstore "github.com/fivenet-app/fivenet/v2025/services/citizenstore/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tJobCitizenLabels  = table.FivenetJobCitizenLabels.AS("citizen_label")
	tUserCitizenLabels = table.FivenetUserCitizenLabels
)

func (s *Server) ManageCitizenLabels(ctx context.Context, req *pbcitizenstore.ManageCitizenLabelsRequest) (*pbcitizenstore.ManageCitizenLabelsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbcitizenstore.CitizenStoreService_ServiceDesc.ServiceName,
		Method:  "ManageCitizenLabels",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	resp := &pbcitizenstore.ManageCitizenLabelsResponse{
		Labels: []*users.CitizenLabel{},
	}

	stmt := tJobCitizenLabels.
		SELECT(
			tJobCitizenLabels.ID,
			tJobCitizenLabels.Job,
			tJobCitizenLabels.Name,
			tJobCitizenLabels.Color,
		).
		FROM(tJobCitizenLabels).
		WHERE(
			tJobCitizenLabels.Job.EQ(jet.String(userInfo.Job)),
		)

	if err := stmt.QueryContext(ctx, s.db, &resp.Labels); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscitizenstore.ErrFailedQuery)
		}
	}

	_, removed := utils.SlicesDifferenceFunc(resp.Labels, req.Labels,
		func(in *users.CitizenLabel) string {
			return in.Name
		})

	for i := range req.Labels {
		req.Labels[i].Job = &userInfo.Job
	}

	tJobCitizenLabels := table.FivenetJobCitizenLabels

	if len(req.Labels) > 0 {
		toCreate := []*users.CitizenLabel{}
		toUpdate := []*users.CitizenLabel{}

		for _, attribute := range req.Labels {
			if attribute.Id == 0 {
				toCreate = append(toCreate, attribute)
			} else {
				toUpdate = append(toUpdate, attribute)
			}
		}

		if len(toCreate) > 0 {
			insertStmt := tJobCitizenLabels.
				INSERT(
					tJobCitizenLabels.Job,
					tJobCitizenLabels.Name,
					tJobCitizenLabels.Color,
				).
				MODELS(toCreate).
				ON_DUPLICATE_KEY_UPDATE(
					tJobCitizenLabels.Name.SET(jet.StringExp(jet.Raw("VALUES(`name`)"))),
					tJobCitizenLabels.Color.SET(jet.StringExp(jet.Raw("VALUES(`color`)"))),
				)

			if _, err := insertStmt.ExecContext(ctx, s.db); err != nil {
				return nil, errswrap.NewError(err, errorscitizenstore.ErrFailedQuery)
			}
		}

		if len(toUpdate) > 0 {
			for _, attribute := range toUpdate {
				updateStmt := tJobCitizenLabels.
					UPDATE(
						tJobCitizenLabels.Name,
						tJobCitizenLabels.Color,
					).
					SET(
						tJobCitizenLabels.Name.SET(jet.String(attribute.Name)),
						tJobCitizenLabels.Color.SET(jet.String(attribute.Color)),
					).
					WHERE(jet.AND(
						tJobCitizenLabels.ID.EQ(jet.Uint64(attribute.Id)),
						tJobCitizenLabels.Job.EQ(jet.String(*attribute.Job)),
					))

				if _, err := updateStmt.ExecContext(ctx, s.db); err != nil {
					return nil, errswrap.NewError(err, errorscitizenstore.ErrFailedQuery)
				}
			}
		}
	}

	if len(removed) > 0 {
		ids := make([]jet.Expression, len(removed))

		for i := range removed {
			ids[i] = jet.Uint64(removed[i].Id)
		}

		deleteStmt := tJobCitizenLabels.
			DELETE().
			WHERE(jet.AND(
				tJobCitizenLabels.ID.IN(ids...),
				tJobCitizenLabels.Job.EQ(jet.String(userInfo.Job)),
			)).
			LIMIT(int64(len(removed)))

		if _, err := deleteStmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorscitizenstore.ErrFailedQuery)
		}
	}

	resp.Labels = []*users.CitizenLabel{}
	if err := stmt.QueryContext(ctx, s.db, &resp.Labels); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscitizenstore.ErrFailedQuery)
		}
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return resp, nil
}

func (s *Server) validateCitizenLabels(ctx context.Context, userInfo *userinfo.UserInfo, attributes []*users.CitizenLabel) (bool, error) {
	if len(attributes) == 0 {
		return true, nil
	}

	jobs, err := s.ps.AttrStringList(userInfo, permscompletor.CompletorServicePerm, permscompletor.CompletorServiceCompleteCitizenLabelsPerm, permscompletor.CompletorServiceCompleteCitizenLabelsJobsPermField)
	if err != nil {
		return false, errswrap.NewError(err, errorscitizenstore.ErrFailedQuery)
	}

	if jobs.Len() == 0 {
		jobs.Strings = append(jobs.Strings, userInfo.Job)
	}

	jobsExp := make([]jet.Expression, len(jobs.Strings))
	for i := range jobs.Strings {
		jobsExp[i] = jet.String(jobs.Strings[i])
	}

	idsExp := make([]jet.Expression, len(attributes))
	for i := range attributes {
		idsExp[i] = jet.Uint64(attributes[i].Id)
	}

	stmt := tJobCitizenLabels.
		SELECT(
			jet.COUNT(tJobCitizenLabels.ID).AS("datacount.totalcount"),
		).
		FROM(tJobCitizenLabels).
		WHERE(jet.AND(
			tJobCitizenLabels.Job.IN(jobsExp...),
			tJobCitizenLabels.ID.IN(idsExp...),
		)).
		LIMIT(10)

	var count database.DataCount
	if err := stmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return false, err
		}
	}

	return len(attributes) == int(count.TotalCount), nil
}

func (s *Server) getUserLabels(ctx context.Context, userInfo *userinfo.UserInfo, userId int32) (*users.CitizenLabels, error) {
	jobs, err := s.ps.AttrStringList(userInfo, permscompletor.CompletorServicePerm, permscompletor.CompletorServiceCompleteCitizenLabelsPerm, permscompletor.CompletorServiceCompleteCitizenLabelsJobsPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizenstore.ErrFailedQuery)
	}

	if jobs.Len() == 0 {
		jobs.Strings = append(jobs.Strings, userInfo.Job)
	}

	jobsExp := make([]jet.Expression, jobs.Len())
	for i := range jobs.Strings {
		jobsExp[i] = jet.String(jobs.Strings[i])
	}

	stmt := tUserCitizenLabels.
		SELECT(
			tJobCitizenLabels.ID,
			tJobCitizenLabels.Job,
			tJobCitizenLabels.Name,
			tJobCitizenLabels.Color,
		).
		FROM(
			tUserCitizenLabels.
				INNER_JOIN(tJobCitizenLabels,
					tJobCitizenLabels.ID.EQ(tUserCitizenLabels.AttributeID),
				),
		).
		WHERE(jet.AND(
			tUserCitizenLabels.UserID.EQ(jet.Int32(userId)),
			tJobCitizenLabels.Job.IN(jobsExp...),
		)).
		ORDER_BY(
			tJobCitizenLabels.SortKey.ASC(),
		)

	list := &users.CitizenLabels{
		List: []*users.CitizenLabel{},
	}
	if err := stmt.QueryContext(ctx, s.db, &list.List); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return list, nil
}
