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
	grpc_audit "github.com/fivenet-app/fivenet/v2025/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorscitizens "github.com/fivenet-app/fivenet/v2025/services/citizens/errors"
	"github.com/go-jet/jet/v2/mysql"
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
			tCitizensLabelsJob.Job.EQ(mysql.String(userInfo.GetJob())),
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
					tCitizensLabelsJob.Name.SET(mysql.StringExp(mysql.Raw("VALUES(`name`)"))),
					tCitizensLabelsJob.Color.SET(mysql.StringExp(mysql.Raw("VALUES(`color`)"))),
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
						tCitizensLabelsJob.Name.SET(mysql.String(attribute.GetName())),
						tCitizensLabelsJob.Color.SET(mysql.String(attribute.GetColor())),
					).
					WHERE(mysql.AND(
						tCitizensLabelsJob.ID.EQ(mysql.Int64(attribute.GetId())),
						tCitizensLabelsJob.Job.EQ(mysql.String(attribute.GetJob())),
					))

				if _, err := updateStmt.ExecContext(ctx, s.db); err != nil {
					return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
				}
			}
		}
	}

	if len(removed) > 0 {
		ids := make([]mysql.Expression, len(removed))

		for i := range removed {
			ids[i] = mysql.Int64(removed[i].GetId())
		}

		deleteStmt := tCitizensLabelsJob.
			DELETE().
			WHERE(mysql.AND(
				tCitizensLabelsJob.ID.IN(ids...),
				tCitizensLabelsJob.Job.EQ(mysql.String(userInfo.GetJob())),
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

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

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

	jobsExp := make([]mysql.Expression, len(jobs.GetStrings()))
	for i := range jobs.GetStrings() {
		jobsExp[i] = mysql.String(jobs.GetStrings()[i])
	}

	idsExp := make([]mysql.Expression, len(attributes))
	for i := range attributes {
		idsExp[i] = mysql.Int64(attributes[i].GetId())
	}

	stmt := tCitizensLabelsJob.
		SELECT(
			mysql.COUNT(tCitizensLabelsJob.ID).AS("data_count.total"),
		).
		FROM(tCitizensLabelsJob).
		WHERE(mysql.AND(
			tCitizensLabelsJob.Job.IN(jobsExp...),
			tCitizensLabelsJob.ID.IN(idsExp...),
		)).
		LIMIT(25)

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

	jobsExp := make([]mysql.Expression, jobs.Len())
	for i := range jobs.GetStrings() {
		jobsExp[i] = mysql.String(jobs.GetStrings()[i])
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
		WHERE(mysql.AND(
			tUserLabels.UserID.EQ(mysql.Int32(userId)),
			tCitizensLabelsJob.Job.IN(jobsExp...),
		)).
		ORDER_BY(tCitizensLabelsJob.SortKey.ASC())

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
