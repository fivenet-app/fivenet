package citizens

import (
	context "context"
	"database/sql"
	"errors"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	citizenslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/citizens/labels"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbcitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorscitizens "github.com/fivenet-app/fivenet/v2026/services/citizens/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tCitizensLabelsJob       = table.FivenetUserLabelsJob.AS("label")
	tCitizenLabels           = table.FivenetUserLabels
	tCitizensLabelsJobAccess = table.FivenetUserLabelsJobJobAccess
)

func (s *Server) ListLabels(
	ctx context.Context,
	req *pbcitizens.ListLabelsRequest,
) (*pbcitizens.ListLabelsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := mysql.Bool(true)

	if !userInfo.GetSuperuser() {
		condition = condition.AND(tCitizensLabelsJob.DeletedAt.IS_NULL())
	}

	if search := dbutils.PrepareForLikeSearch(req.GetSearch()); search != "" {
		condition = condition.AND(tCitizensLabelsJob.Name.LIKE(mysql.String(search)))
	}

	if !userInfo.GetSuperuser() {
		jobAccessExists := mysql.EXISTS(
			mysql.
				SELECT(mysql.Int(1)).
				FROM(tCitizensLabelsJobAccess).
				WHERE(mysql.AND(
					tCitizensLabelsJobAccess.LabelID.EQ(tCitizensLabelsJob.ID),
					tCitizensLabelsJobAccess.Access.GT_EQ(
						mysql.Int32(int32(req.GetMinAccess())),
					),
					tCitizensLabelsJobAccess.Job.EQ(mysql.String(userInfo.GetJob())),
					tCitizensLabelsJobAccess.MinimumGrade.LT_EQ(
						mysql.Int32(userInfo.GetJobGrade()),
					),
				)),
		)

		condition = mysql.AND(
			condition,
			jobAccessExists,
		)
	}

	columns := mysql.ProjectionList{
		tCitizensLabelsJob.ID,
		tCitizensLabelsJob.CreatedAt,
		tCitizensLabelsJob.Name,
		tCitizensLabelsJob.Color,
		tCitizensLabelsJob.Icon,
		tCitizensLabelsJob.Settings,
	}
	if userInfo.GetSuperuser() {
		columns = append(columns, tCitizensLabelsJob.DeletedAt)
	}

	stmt := tCitizensLabelsJob.
		SELECT(
			columns[0],
			columns[1:],
		).
		FROM(tCitizensLabelsJob).
		WHERE(condition).
		ORDER_BY(
			tCitizensLabelsJob.SortKey.ASC(),
		).
		LIMIT(15)

	resp := &pbcitizens.ListLabelsResponse{
		Labels: []*citizenslabels.Label{},
	}
	if err := stmt.QueryContext(ctx, s.db, &resp.Labels); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
	}

	return resp, nil
}

func (s *Server) GetLabel(
	ctx context.Context,
	req *pbcitizens.GetLabelRequest,
) (*pbcitizens.GetLabelResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	label, err := s.getLabel(ctx, s.db, userInfo.GetJob(), req.GetId())
	if err != nil {
		return nil, err
	}

	if label.GetId() == 0 {
		return nil, errorscitizens.ErrLabelNotFound
	}

	if !userInfo.GetSuperuser() {
		if label.GetDeletedAt() != nil {
			return nil, errorscitizens.ErrLabelNotFound
		}

		label.DeletedAt = nil
	}

	return &pbcitizens.GetLabelResponse{
		Label: label,
	}, nil
}

func (s *Server) getLabel(
	ctx context.Context,
	tx *sql.DB,
	job string,
	labelId int64,
) (*citizenslabels.Label, error) {
	stmt := tCitizensLabelsJob.
		SELECT(
			tCitizensLabelsJob.ID,
			tCitizensLabelsJob.CreatedAt,
			tCitizensLabelsJob.UpdatedAt,
			tCitizensLabelsJob.DeletedAt,
			tCitizensLabelsJob.Job,
			tCitizensLabelsJob.Name,
			tCitizensLabelsJob.Color,
			tCitizensLabelsJob.Icon,
			tCitizensLabelsJob.Settings,
		).
		FROM(tCitizensLabelsJob).
		WHERE(mysql.AND(
			tCitizensLabelsJob.ID.EQ(mysql.Int64(labelId)),
			tCitizensLabelsJob.Job.EQ(mysql.String(job)),
		)).
		LIMIT(1)

	label := &citizenslabels.Label{}
	if err := stmt.QueryContext(ctx, tx, label); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
	}

	if label.GetId() == 0 {
		return nil, errorscitizens.ErrLabelNotFound
	}

	return label, nil
}

func (s *Server) CreateOrUpdateLabel(
	ctx context.Context,
	req *pbcitizens.CreateOrUpdateLabelRequest,
) (*pbcitizens.CreateOrUpdateLabelResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	label := req.GetLabel()
	label.Job = &userInfo.Job

	tCitizensLabelsJob := table.FivenetUserLabelsJob

	if req.GetLabel().GetId() > 0 {
		stmt := tCitizensLabelsJob.
			UPDATE(
				tCitizensLabelsJob.Name,
				tCitizensLabelsJob.Color,
				tCitizensLabelsJob.Icon,
				tCitizensLabelsJob.Settings,
			).
			SET(
				label.Name,
				label.Color,
				label.Icon,
				label.Settings,
			).
			WHERE(mysql.AND(
				tCitizensLabelsJob.ID.EQ(mysql.Int64(label.GetId())),
				tCitizensLabelsJob.Job.EQ(mysql.String(userInfo.GetJob())),
			)).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
	} else {
		stmt := tCitizensLabelsJob.
			INSERT(
				tCitizensLabelsJob.Job,
				tCitizensLabelsJob.Name,
				tCitizensLabelsJob.Color,
				tCitizensLabelsJob.Icon,
				tCitizensLabelsJob.Settings,
			).
			VALUES(
				label.Job,
				label.Name,
				label.Color,
				label.Icon,
				label.Settings,
			)

		result, err := stmt.ExecContext(ctx, s.db)
		if err != nil {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}

		lastInsertId, err := result.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}

		label.Id = lastInsertId
	}

	label, err := s.getLabel(ctx, s.db, label.GetJob(), label.GetId())
	if err != nil {
		return nil, err
	}
	if !userInfo.GetSuperuser() {
		label.DeletedAt = nil
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return &pbcitizens.CreateOrUpdateLabelResponse{
		Label: label,
	}, nil
}

func (s *Server) DeleteLabel(
	ctx context.Context,
	req *pbcitizens.DeleteLabelRequest,
) (*pbcitizens.DeleteLabelResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	label, err := s.getLabel(ctx, s.db, userInfo.GetJob(), req.GetId())
	if err != nil {
		return nil, err
	}
	if label == nil || label.GetJob() != userInfo.GetJob() {
		return nil, errorscitizens.ErrFailedQuery
	}

	deletedAtTime := mysql.CURRENT_TIMESTAMP()
	if label.GetDeletedAt() != nil && userInfo.GetSuperuser() {
		deletedAtTime = mysql.TimestampExp(mysql.NULL)
	}

	stmt := tCitizensLabelsJob.
		UPDATE(
			tCitizensLabelsJob.DeletedAt,
		).
		SET(
			tCitizensLabelsJob.DeletedAt.SET(deletedAtTime),
		).
		WHERE(mysql.AND(
			tCitizensLabelsJob.ID.EQ(mysql.Int64(req.GetId())),
			tCitizensLabelsJob.Job.EQ(mysql.String(userInfo.GetJob())),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)

	return &pbcitizens.DeleteLabelResponse{}, nil
}

func (s *Server) getUserLabels(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	userId int32,
) (*citizenslabels.Labels, error) {
	jobAccessExists := mysql.EXISTS(
		mysql.
			SELECT(mysql.Int(1)).
			FROM(tCitizensLabelsJobAccess).
			WHERE(mysql.AND(
				tCitizensLabelsJobAccess.LabelID.EQ(tCitizensLabelsJob.ID),
				tCitizensLabelsJobAccess.Access.GT_EQ(
					mysql.Int32(int32(citizenslabels.AccessLevel_ACCESS_LEVEL_VIEW)),
				),
				tCitizensLabelsJobAccess.Job.EQ(mysql.String(userInfo.GetJob())),
				tCitizensLabelsJobAccess.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),
			)),
	)

	stmt := tCitizenLabels.
		SELECT(
			tCitizensLabelsJob.ID,
			tCitizensLabelsJob.Job,
			tCitizensLabelsJob.Name,
			tCitizensLabelsJob.Color,
			tCitizensLabelsJob.Icon,
			tCitizensLabelsJob.Settings,
		).
		FROM(
			tCitizenLabels.
				INNER_JOIN(tCitizensLabelsJob,
					tCitizensLabelsJob.ID.EQ(tCitizenLabels.LabelID),
				),
		).
		WHERE(mysql.AND(
			tCitizensLabelsJob.DeletedAt.IS_NULL(),
			tCitizenLabels.UserID.EQ(mysql.Int32(userId)),
			jobAccessExists,
		)).
		ORDER_BY(tCitizensLabelsJob.SortKey.ASC(), tCitizensLabelsJob.ID.DESC()).
		LIMIT(25)

	list := &citizenslabels.Labels{
		List: []*citizenslabels.Label{},
	}
	if err := stmt.QueryContext(ctx, s.db, &list.List); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return list, nil
}

func (s *Server) validateLabels(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	labels []*citizenslabels.Label,
) (bool, error) {
	if len(labels) == 0 {
		return true, nil
	}

	idsExp := make([]mysql.Expression, len(labels))
	for i := range labels {
		idsExp[i] = mysql.Int64(labels[i].GetId())
	}

	stmt := tCitizensLabelsJob.
		SELECT(
			mysql.COUNT(tCitizensLabelsJob.ID).AS("data_count.total"),
		).
		FROM(tCitizensLabelsJob).
		WHERE(mysql.AND(
			tCitizensLabelsJob.Job.EQ(mysql.String(userInfo.GetJob())),
			tCitizensLabelsJob.ID.IN(idsExp...),
		)).
		LIMIT(20)

	var count database.DataCount
	if err := stmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return false, err
		}
	}

	return len(labels) == int(count.Total), nil
}

func (s *Server) AddCitizenLabels(
	ctx context.Context,
	req *pbcitizens.AddCitizenLabelsRequest,
) (*pbcitizens.AddCitizenLabelsResponse, error) {
	// TODO

	return nil, nil
}

func (s *Server) RemoveCitizenLabels(
	ctx context.Context,
	req *pbcitizens.RemoveCitizenLabelsRequest,
) (*pbcitizens.RemoveCitizenLabelsResponse, error) {
	// TODO

	return nil, nil
}
