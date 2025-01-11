package citizenstore

import (
	context "context"
	"errors"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	users "github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	errorscitizenstore "github.com/fivenet-app/fivenet/gen/go/proto/services/citizenstore/errors"
	permscompletor "github.com/fivenet-app/fivenet/gen/go/proto/services/completor/perms"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/utils"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tJobCitizenLabels  = table.FivenetJobCitizenLabels.AS("citizen_label")
	tUserCitizenLabels = table.FivenetUserCitizenLabels
)

func (s *Server) ManageCitizenLabels(ctx context.Context, req *ManageCitizenLabelsRequest) (*ManageCitizenLabelsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CitizenStoreService_ServiceDesc.ServiceName,
		Method:  "ManageCitizenLabels",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	resp := &ManageCitizenLabelsResponse{
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

	for i := 0; i < len(req.Labels); i++ {
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

	jobsAttr, err := s.ps.Attr(userInfo, permscompletor.CompletorServicePerm, permscompletor.CompletorServiceCompleteCitizenLabelsPerm, permscompletor.CompletorServiceCompleteCitizenLabelsJobsPermField)
	if err != nil {
		return false, errswrap.NewError(err, errorscitizenstore.ErrFailedQuery)
	}
	var jobs perms.StringList
	if jobsAttr != nil {
		jobs = jobsAttr.([]string)
	}

	if len(jobs) == 0 {
		jobs = append(jobs, userInfo.Job)
	}

	jobsExp := make([]jet.Expression, len(jobs))
	for i := 0; i < len(jobs); i++ {
		jobsExp[i] = jet.String(jobs[i])
	}

	idsExp := make([]jet.Expression, len(attributes))
	for i := 0; i < len(attributes); i++ {
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
		ORDER_BY(
			tJobCitizenLabels.Name.DESC(),
		).
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
	jobsAttr, err := s.ps.Attr(userInfo, permscompletor.CompletorServicePerm, permscompletor.CompletorServiceCompleteCitizenLabelsPerm, permscompletor.CompletorServiceCompleteCitizenLabelsJobsPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizenstore.ErrFailedQuery)
	}
	var jobs perms.StringList
	if jobsAttr != nil {
		jobs = jobsAttr.([]string)
	}

	if len(jobs) == 0 {
		jobs = append(jobs, userInfo.Job)
	}

	jobsExp := make([]jet.Expression, len(jobs))
	for i := 0; i < len(jobs); i++ {
		jobsExp[i] = jet.String(jobs[i])
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
		))

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
