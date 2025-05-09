package jobs

import (
	context "context"
	"errors"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	jobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/rector"
	pbjobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/jobs"
	permsjobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/jobs/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorscitizenstore "github.com/fivenet-app/fivenet/v2025/services/citizenstore/errors"
	errorsjobs "github.com/fivenet-app/fivenet/v2025/services/jobs/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tJobLabels  = table.FivenetJobsLabels.AS("label")
	tUserLabels = table.FivenetJobsLabelsUsers
)

func (s *Server) GetColleagueLabels(ctx context.Context, req *pbjobs.GetColleagueLabelsRequest) (*pbjobs.GetColleagueLabelsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	resp := &pbjobs.GetColleagueLabelsResponse{
		Labels: []*jobs.Label{},
	}

	// Fields Permission Check
	fields, err := s.ps.AttrStringList(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceGetColleaguePerm, permsjobs.JobsServiceGetColleagueTypesPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if userInfo.SuperUser {
		fields.Strings = []string{"Labels"}
	}

	if !fields.Contains("Labels") {
		// Fallback to checking if user has manage colleague labels permission
		if !s.ps.Can(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceManageColleagueLabelsPerm) {
			return nil, errorsjobs.ErrLabelsNoPerms
		}
	}

	condition := tJobLabels.Job.EQ(jet.String(userInfo.Job))

	if req.Search != nil && *req.Search != "" {
		*req.Search = strings.TrimSpace(*req.Search)
		*req.Search = strings.ReplaceAll(*req.Search, "%", "")
		*req.Search = strings.ReplaceAll(*req.Search, " ", "%")
		*req.Search = "%" + *req.Search + "%"
		condition = condition.AND(jet.OR(
			tJobLabels.Name.LIKE(jet.String(*req.Search)),
		))
	}

	stmt := tJobLabels.
		SELECT(
			tJobLabels.ID,
			tJobLabels.Job,
			tJobLabels.Name,
			tJobLabels.Color,
			tJobLabels.Order,
		).
		FROM(tJobLabels).
		WHERE(condition).
		ORDER_BY(
			tJobLabels.Order.ASC(),
			tJobLabels.SortKey.ASC(),
		)

	if err := stmt.QueryContext(ctx, s.db, &resp.Labels); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscitizenstore.ErrFailedQuery)
		}
	}

	return resp, nil
}

func (s *Server) ManageColleagueLabels(ctx context.Context, req *pbjobs.ManageColleagueLabelsRequest) (*pbjobs.ManageColleagueLabelsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbjobs.JobsService_ServiceDesc.ServiceName,
		Method:  "ManageColleagueLabels",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	stmt := tJobLabels.
		SELECT(
			tJobLabels.ID,
			tJobLabels.Job,
			tJobLabels.Name,
			tJobLabels.Color,
			tJobLabels.Order,
		).
		FROM(tJobLabels).
		WHERE(
			tJobLabels.Job.EQ(jet.String(userInfo.Job)),
		)

	labels := []*jobs.Label{}
	if err := stmt.QueryContext(ctx, s.db, &labels); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscitizenstore.ErrFailedQuery)
		}
	}

	_, removed := utils.SlicesDifferenceFunc(labels, req.Labels,
		func(in *jobs.Label) uint64 {
			return in.Id
		})

	for i := range req.Labels {
		req.Labels[i].Job = &userInfo.Job
		req.Labels[i].Order = int32(i)
	}

	tJobLabels := table.FivenetJobsLabels
	if len(req.Labels) > 0 {
		toCreate := []*jobs.Label{}
		toUpdate := []*jobs.Label{}

		for _, label := range req.Labels {
			if label.Id == 0 {
				toCreate = append(toCreate, label)
			} else {
				toUpdate = append(toUpdate, label)
			}
		}

		if len(toCreate) > 0 {
			insertStmt := tJobLabels.
				INSERT(
					tJobLabels.Job,
					tJobLabels.Name,
					tJobLabels.Color,
					tJobLabels.Order,
				).
				MODELS(toCreate).
				ON_DUPLICATE_KEY_UPDATE(
					tJobLabels.Name.SET(jet.StringExp(jet.Raw("VALUES(`name`)"))),
					tJobLabels.Color.SET(jet.StringExp(jet.Raw("VALUES(`color`)"))),
					tJobLabels.Order.SET(jet.IntExp(jet.Raw("VALUES(`order`)"))),
				)

			if _, err := insertStmt.ExecContext(ctx, s.db); err != nil {
				return nil, errswrap.NewError(err, errorscitizenstore.ErrFailedQuery)
			}
		}

		if len(toUpdate) > 0 {
			for _, label := range toUpdate {
				updateStmt := tJobLabels.
					UPDATE(
						tJobLabels.Name,
						tJobLabels.Color,
						tJobLabels.Order,
					).
					SET(
						tJobLabels.Name.SET(jet.String(label.Name)),
						tJobLabels.Color.SET(jet.String(label.Color)),
						tJobLabels.Order.SET(jet.Int32(label.Order)),
					).
					WHERE(jet.AND(
						tJobLabels.ID.EQ(jet.Uint64(label.Id)),
						tJobLabels.Job.EQ(jet.String(*label.Job)),
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

		deleteStmt := tJobLabels.
			DELETE().
			WHERE(jet.AND(
				tJobLabels.ID.IN(ids...),
				tJobLabels.Job.EQ(jet.String(userInfo.Job)),
			)).
			LIMIT(int64(len(removed)))

		if _, err := deleteStmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorscitizenstore.ErrFailedQuery)
		}
	}

	resp := &pbjobs.ManageColleagueLabelsResponse{
		Labels: []*jobs.Label{},
	}
	if err := stmt.QueryContext(ctx, s.db, &resp.Labels); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscitizenstore.ErrFailedQuery)
		}
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return resp, nil
}

func (s *Server) validateLabels(ctx context.Context, userInfo *userinfo.UserInfo, labels []*jobs.Label) (bool, error) {
	if len(labels) == 0 {
		return true, nil
	}

	idsExp := make([]jet.Expression, len(labels))
	for i := range labels {
		idsExp[i] = jet.Uint64(labels[i].Id)
	}

	stmt := tJobLabels.
		SELECT(
			jet.COUNT(tJobLabels.ID).AS("datacount.totalcount"),
		).
		FROM(tJobLabels).
		WHERE(jet.AND(
			tJobLabels.Job.EQ(jet.String(userInfo.Job)),
			tJobLabels.ID.IN(idsExp...),
		)).
		ORDER_BY(
			tJobLabels.Name.DESC(),
		).
		LIMIT(10)

	var count database.DataCount
	if err := stmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return false, err
		}
	}

	return len(labels) == int(count.TotalCount), nil
}

func (s *Server) getUserLabels(ctx context.Context, userInfo *userinfo.UserInfo, userId int32) (*jobs.Labels, error) {
	stmt := tUserLabels.
		SELECT(
			tJobLabels.ID,
			tJobLabels.Job,
			tJobLabels.Name,
			tJobLabels.Color,
		).
		FROM(
			tUserLabels.
				INNER_JOIN(tJobLabels,
					tJobLabels.ID.EQ(tUserLabels.LabelID),
				),
		).
		WHERE(jet.AND(
			tUserLabels.UserID.EQ(jet.Int32(userId)),
			tJobLabels.Job.EQ(jet.String(userInfo.Job)),
		)).
		ORDER_BY(
			tJobLabels.Order.ASC(),
		)

	list := &jobs.Labels{
		List: []*jobs.Label{},
	}
	if err := stmt.QueryContext(ctx, s.db, &list.List); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return list, nil
}

func (s *Server) GetColleagueLabelsStats(ctx context.Context, req *pbjobs.GetColleagueLabelsStatsRequest) (*pbjobs.GetColleagueLabelsStatsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Types Permission Check
	fields, err := s.ps.AttrStringList(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceGetColleaguePerm, permsjobs.JobsServiceGetColleagueTypesPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if userInfo.SuperUser {
		fields.Strings = []string{"Labels"}
	}

	if !fields.Contains("Labels") {
		return &pbjobs.GetColleagueLabelsStatsResponse{}, nil
	}

	tUser := tables.Users().AS("user")

	stmt := tUserLabels.
		SELECT(
			jet.COUNT(tUserLabels.LabelID).AS("label_count.count"),
			tJobLabels.ID,
			tJobLabels.Job,
			tJobLabels.Name,
			tJobLabels.Color,
		).
		FROM(
			tUserLabels.
				INNER_JOIN(tJobLabels,
					tJobLabels.ID.EQ(tUserLabels.LabelID),
				).
				INNER_JOIN(tUser,
					tUser.ID.EQ(tUserLabels.UserID),
				),
		).
		WHERE(jet.AND(
			tUserLabels.Job.EQ(jet.String(userInfo.Job)),
			tUser.Job.EQ(jet.String(userInfo.Job)),
		)).
		GROUP_BY(tJobLabels.ID).
		ORDER_BY(
			tJobLabels.Order.ASC(),
		)

	dest := []*jobs.LabelCount{}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscitizenstore.ErrFailedQuery)
		}
	}

	return &pbjobs.GetColleagueLabelsStatsResponse{
		Count: dest,
	}, nil
}
