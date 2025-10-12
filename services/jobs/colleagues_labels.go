package jobs

import (
	context "context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	jobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	pbjobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/jobs"
	permsjobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/jobs/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2025/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorscitizens "github.com/fivenet-app/fivenet/v2025/services/citizens/errors"
	errorsjobs "github.com/fivenet-app/fivenet/v2025/services/jobs/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tJobLabels       = table.FivenetJobLabels.AS("label")
	tColleagueLabels = table.FivenetJobColleagueLabels
)

func (s *Server) GetColleagueLabels(
	ctx context.Context,
	req *pbjobs.GetColleagueLabelsRequest,
) (*pbjobs.GetColleagueLabelsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	resp := &pbjobs.GetColleagueLabelsResponse{
		Labels: []*jobs.Label{},
	}

	// Fields Permission Check
	fields, err := s.ps.AttrStringList(
		userInfo,
		permsjobs.JobsServicePerm,
		permsjobs.JobsServiceGetColleaguePerm,
		permsjobs.JobsServiceGetColleagueTypesPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if userInfo.GetSuperuser() {
		fields.Strings = []string{"Labels"}
	}

	if !fields.Contains("Labels") {
		// Fallback to checking if user has manage colleague labels permission
		if !s.ps.Can(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceManageLabelsPerm) {
			return nil, errorsjobs.ErrLabelsNoPerms
		}
	}

	condition := mysql.AND(
		tJobLabels.Job.EQ(mysql.String(userInfo.GetJob())),
		tJobLabels.DeletedAt.IS_NULL(),
	)

	if req.GetSearch() != "" {
		search := dbutils.PrepareForLikeSearch(req.GetSearch())
		if search != "" {
			condition = condition.AND(mysql.OR(
				tJobLabels.Name.LIKE(mysql.String(search)),
			))
		}
	}

	stmt := tJobLabels.
		SELECT(
			tJobLabels.ID,
			tJobLabels.Job,
			tJobLabels.DeletedAt,
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
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
	}

	return resp, nil
}

func (s *Server) ManageLabels(
	ctx context.Context,
	req *pbjobs.ManageLabelsRequest,
) (*pbjobs.ManageLabelsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	stmt := tJobLabels.
		SELECT(
			tJobLabels.ID,
			tJobLabels.DeletedAt,
			tJobLabels.Job,
			tJobLabels.Name,
			tJobLabels.Color,
			tJobLabels.Order,
		).
		FROM(tJobLabels).
		WHERE(mysql.AND(
			tJobLabels.Job.EQ(mysql.String(userInfo.GetJob())),
		))

	labels := []*jobs.Label{}
	if err := stmt.QueryContext(ctx, s.db, &labels); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
	}

	_, removed := utils.SlicesDifferenceFunc(labels, req.GetLabels(),
		func(in *jobs.Label) int64 {
			return in.GetId()
		})

	var i int32
	for _, label := range req.GetLabels() {
		label.Job = &userInfo.Job
		label.Order = i
		i++
	}

	tJobLabels := table.FivenetJobLabels
	if len(req.GetLabels()) > 0 {
		toCreate := []*jobs.Label{}
		toUpdate := []*jobs.Label{}

		for _, label := range req.GetLabels() {
			if label.GetId() == 0 {
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
					tJobLabels.Name.SET(mysql.StringExp(mysql.Raw("VALUES(`name`)"))),
					tJobLabels.Color.SET(mysql.StringExp(mysql.Raw("VALUES(`color`)"))),
					tJobLabels.Order.SET(mysql.IntExp(mysql.Raw("VALUES(`order`)"))),
					tJobLabels.DeletedAt.SET(mysql.TimestampExp(mysql.NULL)),
				)

			if _, err := insertStmt.ExecContext(ctx, s.db); err != nil {
				return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
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
						tJobLabels.Name.SET(mysql.String(label.GetName())),
						tJobLabels.Color.SET(mysql.String(label.GetColor())),
						tJobLabels.Order.SET(mysql.Int32(label.GetOrder())),
						tJobLabels.DeletedAt.SET(mysql.TimestampExp(mysql.NULL)),
					).
					WHERE(mysql.AND(
						tJobLabels.ID.EQ(mysql.Int64(label.GetId())),
						tJobLabels.Job.EQ(mysql.String(label.GetJob())),
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

		deleteStmt := tJobLabels.
			UPDATE().
			SET(
				tJobLabels.DeletedAt.SET(mysql.CURRENT_TIMESTAMP()),
			).
			WHERE(mysql.AND(
				tJobLabels.ID.IN(ids...),
				tJobLabels.Job.EQ(mysql.String(userInfo.GetJob())),
			)).
			LIMIT(int64(len(removed)))

		if _, err := deleteStmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
	}

	resp := &pbjobs.ManageLabelsResponse{
		Labels: []*jobs.Label{},
	}
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
	labels []*jobs.Label,
) (bool, error) {
	if len(labels) == 0 {
		return true, nil
	}

	idsExp := make([]mysql.Expression, len(labels))
	for i := range labels {
		idsExp[i] = mysql.Int64(labels[i].GetId())
	}

	stmt := tJobLabels.
		SELECT(
			mysql.COUNT(tJobLabels.ID).AS("data_count.total"),
		).
		FROM(tJobLabels).
		WHERE(mysql.AND(
			tJobLabels.Job.EQ(mysql.String(userInfo.GetJob())),
			tJobLabels.DeletedAt.IS_NULL(),
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

	return len(labels) == int(count.Total), nil
}

func (s *Server) getUserLabels(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	userId int32,
) (*jobs.Labels, error) {
	stmt := tColleagueLabels.
		SELECT(
			tJobLabels.ID,
			tJobLabels.Job,
			tJobLabels.Name,
			tJobLabels.Color,
		).
		FROM(
			tColleagueLabels.
				INNER_JOIN(tJobLabels,
					tJobLabels.ID.EQ(tColleagueLabels.LabelID),
				),
		).
		WHERE(mysql.AND(
			tColleagueLabels.UserID.EQ(mysql.Int32(userId)),
			tJobLabels.Job.EQ(mysql.String(userInfo.GetJob())),
			tJobLabels.DeletedAt.IS_NULL(),
		)).
		ORDER_BY(tJobLabels.Order.ASC())

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

func (s *Server) GetColleagueLabelsStats(
	ctx context.Context,
	req *pbjobs.GetColleagueLabelsStatsRequest,
) (*pbjobs.GetColleagueLabelsStatsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Types Permission Check
	fields, err := s.ps.AttrStringList(
		userInfo,
		permsjobs.JobsServicePerm,
		permsjobs.JobsServiceGetColleaguePerm,
		permsjobs.JobsServiceGetColleagueTypesPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if userInfo.GetSuperuser() {
		fields.Strings = []string{"Labels"}
	}

	if !fields.Contains("Labels") {
		return &pbjobs.GetColleagueLabelsStatsResponse{}, nil
	}

	tColleague := tables.User().AS("user")

	stmt := tColleagueLabels.
		SELECT(
			mysql.COUNT(tColleagueLabels.LabelID).AS("label_count.count"),
			tJobLabels.ID,
			tJobLabels.Job,
			tJobLabels.Name,
			tJobLabels.Color,
		).
		FROM(
			tColleagueLabels.
				INNER_JOIN(tJobLabels,
					tJobLabels.ID.EQ(tColleagueLabels.LabelID),
				).
				INNER_JOIN(tColleague,
					tColleague.ID.EQ(tColleagueLabels.UserID),
				),
		).
		WHERE(mysql.AND(
			tJobLabels.Job.EQ(mysql.String(userInfo.GetJob())),
			tJobLabels.DeletedAt.IS_NULL(),
			tColleagueLabels.Job.EQ(mysql.String(userInfo.GetJob())),
			tColleague.Job.EQ(mysql.String(userInfo.GetJob())),
		)).
		GROUP_BY(tJobLabels.ID).
		ORDER_BY(
			tJobLabels.Order.ASC(),
		)

	dest := []*jobs.LabelCount{}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
	}

	return &pbjobs.GetColleagueLabelsStatsResponse{
		Count: dest,
	}, nil
}
