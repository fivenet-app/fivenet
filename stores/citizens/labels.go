package citizensstore

import (
	"context"
	"errors"
	"math"

	citizenslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/citizens/labels"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tCitizensLabelsJob = table.FivenetUserLabelsJob
	tCitizenLabels     = table.FivenetUserLabels
)

func (s *Store) ListLabels(
	ctx context.Context,
	q qrm.Queryable,
	userInfo *userinfo.UserInfo,
	search string,
	ownJobOnly bool,
	canCreateLabel bool,
	minAccess int32,
	includeDeleted bool,
) (*citizenslabels.Labels, error) {
	tLabel := tCitizensLabelsJob.AS("label")
	visibilityCondition := buildLabelVisibilityCondition(tLabel, search)

	columns := mysql.ProjectionList{
		tLabel.ID,
		tLabel.CreatedAt,
		tLabel.SortOrder,
		tLabel.Name,
		tLabel.Color,
		tLabel.Icon,
		tLabel.Settings,
	}
	if includeDeleted {
		columns = append(columns, tLabel.DeletedAt)
	}

	var (
		stmt mysql.Statement
		ctes []mysql.CommonTableExpression
	)

	if userInfo.GetSuperuser() || (ownJobOnly && canCreateLabel) {
		visibilityCondition = visibilityCondition.AND(
			tLabel.Job.EQ(mysql.String(userInfo.GetJob())),
		)

		stmt = tLabel.
			SELECT(columns[0], columns[1:]...).
			FROM(tLabel).
			WHERE(visibilityCondition).
			GROUP_BY(
				tLabel.ID,
				tLabel.CreatedAt,
				tLabel.Name,
				tLabel.Color,
				tLabel.Icon,
				tLabel.Settings,
				tLabel.SortKey,
			).
			ORDER_BY(
				tLabel.SortOrder.ASC(),
				tLabel.SortKey.ASC(),
			).
			LIMIT(20)
	} else {
		visibleIDs := s.labelsAccess.VisibleIDsByConditionQuery(
			userInfo,
			minAccess,
			false,
			buildLabelVisibilityCondition(tCitizensLabelsJob, search),
		)
		ctes = visibleIDs.CTEs
		visibleLabelID := mysql.IntegerColumn("id").From(visibleIDs.Table)

		stmt = tLabel.
			SELECT(columns[0], columns[1:]...).
			FROM(
				visibleIDs.Table.
					INNER_JOIN(tLabel, tLabel.ID.EQ(visibleLabelID)),
			).
			GROUP_BY(
				tLabel.ID,
				tLabel.CreatedAt,
				tLabel.Name,
				tLabel.Color,
				tLabel.Icon,
				tLabel.Settings,
				tLabel.SortKey,
			).
			ORDER_BY(
				tLabel.SortOrder.ASC(),
				tLabel.SortKey.ASC(),
			).
			LIMIT(20)
	}

	if len(ctes) > 0 {
		stmt = mysql.WITH(ctes...)(stmt)
	}

	resp := &citizenslabels.Labels{
		List: []*citizenslabels.Label{},
	}
	if err := stmt.QueryContext(ctx, q, &resp.List); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return resp, nil
}

func buildLabelVisibilityCondition(
	label *table.FivenetUserLabelsJobTable,
	search string,
) mysql.BoolExpression {
	condition := label.DeletedAt.IS_NULL()

	if search = dbutils.PrepareForLikeSearch(search); search != "" {
		condition = condition.AND(label.Name.LIKE(mysql.String(search)))
	}

	return condition
}

func (s *Store) NextLabelSortOrder(
	ctx context.Context,
	q qrm.Queryable,
	job string,
) (int32, error) {
	stmt := tCitizensLabelsJob.
		SELECT(
			mysql.COALESCE(mysql.MAX(tCitizensLabelsJob.SortOrder), mysql.Int32(-1)).
				AS("sort_order"),
		).
		FROM(tCitizensLabelsJob).
		WHERE(mysql.AND(
			tCitizensLabelsJob.Job.EQ(mysql.String(job)),
			tCitizensLabelsJob.DeletedAt.IS_NULL(),
		))

	var dest struct {
		SortOrder int32 `alias:"sort_order"`
	}
	if err := stmt.QueryContext(ctx, q, &dest); err != nil {
		return 0, err
	}

	return dest.SortOrder + 1, nil
}

func (s *Store) GetLabel(
	ctx context.Context,
	q qrm.Queryable,
	job string,
	labelId int64,
	includeDeleted bool,
) (*citizenslabels.Label, error) {
	tCitizensLabelsJob := tCitizensLabelsJob.AS("label")

	stmt := tCitizensLabelsJob.
		SELECT(
			tCitizensLabelsJob.ID,
			tCitizensLabelsJob.CreatedAt,
			tCitizensLabelsJob.UpdatedAt,
			tCitizensLabelsJob.DeletedAt,
			tCitizensLabelsJob.Job,
			tCitizensLabelsJob.SortOrder,
			tCitizensLabelsJob.Name,
			tCitizensLabelsJob.Color,
			tCitizensLabelsJob.Icon,
			tCitizensLabelsJob.Settings,
		).
		FROM(tCitizensLabelsJob).
		WHERE(mysql.AND(
			tCitizensLabelsJob.ID.EQ(mysql.Int64(labelId)),
			tCitizensLabelsJob.Job.EQ(mysql.String(job)),
			mysql.OR(
				mysql.Bool(includeDeleted),
				tCitizensLabelsJob.DeletedAt.IS_NULL(),
			),
		)).
		LIMIT(1)

	label := &citizenslabels.Label{}
	if err := stmt.QueryContext(ctx, q, label); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if label.GetId() == 0 {
		return nil, nil
	}

	return label, nil
}

func (s *Store) UpdateLabel(
	ctx context.Context,
	tx qrm.DB,
	label *citizenslabels.Label,
	job string,
) error {
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
			tCitizensLabelsJob.Job.EQ(mysql.String(job)),
		)).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) InsertLabel(
	ctx context.Context,
	tx qrm.DB,
	label *citizenslabels.Label,
) (int64, error) {
	stmt := tCitizensLabelsJob.
		INSERT(
			tCitizensLabelsJob.Job,
			tCitizensLabelsJob.SortOrder,
			tCitizensLabelsJob.Name,
			tCitizensLabelsJob.Color,
			tCitizensLabelsJob.Icon,
			tCitizensLabelsJob.Settings,
		).
		VALUES(
			label.Job,
			label.SortOrder,
			label.Name,
			label.Color,
			label.Icon,
			label.Settings,
		)

	result, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return 0, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastInsertId, nil
}

func (s *Store) DeleteLabel(
	ctx context.Context,
	tx qrm.DB,
	job string,
	labelId int64,
	deletedAt *timestamp.Timestamp,
) error {
	stmt := tCitizensLabelsJob.
		UPDATE(
			tCitizensLabelsJob.DeletedAt,
		).
		SET(
			tCitizensLabelsJob.DeletedAt.SET(dbutils.TimestampToMySQL(deletedAt)),
		).
		WHERE(mysql.AND(
			tCitizensLabelsJob.ID.EQ(mysql.Int64(labelId)),
			tCitizensLabelsJob.Job.EQ(mysql.String(job)),
		)).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) ReorderLabels(
	ctx context.Context,
	job string,
	labelIds []int64,
) error {
	stmt := tCitizensLabelsJob.
		SELECT(tCitizensLabelsJob.ID).
		FROM(tCitizensLabelsJob).
		WHERE(mysql.AND(
			tCitizensLabelsJob.Job.EQ(mysql.String(job)),
			tCitizensLabelsJob.DeletedAt.IS_NULL(),
		)).
		LIMIT(int64(len(labelIds)))

	var dest []int64
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	existing := make(map[int64]struct{}, len(labelIds))
	for _, labelID := range dest {
		existing[labelID] = struct{}{}
	}

	if len(existing) != len(labelIds) {
		return errors.New("invalid labels")
	}

	for _, labelID := range labelIds {
		if _, ok := existing[labelID]; !ok {
			return errors.New("invalid labels")
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for idx, labelID := range labelIds {
		if idx > math.MaxInt32 {
			return errors.New("invalid labels")
		}

		if _, err := tCitizensLabelsJob.
			UPDATE().
			SET(
				tCitizensLabelsJob.SortOrder.SET(mysql.Int32(int32(idx))),
			).
			WHERE(mysql.AND(
				tCitizensLabelsJob.ID.EQ(mysql.Int64(labelID)),
				tCitizensLabelsJob.Job.EQ(mysql.String(job)),
				tCitizensLabelsJob.DeletedAt.IS_NULL(),
			)).
			LIMIT(1).
			ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Store) GetUserLabels(
	ctx context.Context,
	q qrm.Queryable,
	condition mysql.BoolExpression,
) (*citizenslabels.Labels, error) {
	tCitizensLabelsJob := tCitizensLabelsJob.AS("label")

	stmt := tCitizenLabels.
		SELECT(
			tCitizensLabelsJob.ID,
			tCitizensLabelsJob.Job,
			tCitizensLabelsJob.Name,
			tCitizensLabelsJob.Color,
			tCitizensLabelsJob.Icon,
			tCitizensLabelsJob.Settings,
			tCitizenLabels.ExpiresAt.AS("label.expiresAt"),
		).
		FROM(
			tCitizenLabels.
				INNER_JOIN(tCitizensLabelsJob,
					tCitizensLabelsJob.ID.EQ(tCitizenLabels.LabelID),
				),
		).
		WHERE(condition).
		ORDER_BY(
			tCitizensLabelsJob.SortOrder.ASC(),
			tCitizensLabelsJob.SortKey.ASC(),
			tCitizensLabelsJob.ID.DESC(),
		).
		LIMIT(25)

	list := &citizenslabels.Labels{
		List: []*citizenslabels.Label{},
	}
	if err := stmt.QueryContext(ctx, q, &list.List); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return list, nil
}

func (s *Store) GetUserLabelsForUser(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	userId int32,
) (*citizenslabels.Labels, error) {
	tCitizensLabelsJob := table.FivenetUserLabelsJob.AS("label")

	includeDeleted := userInfo.GetSuperuser()
	condition := mysql.AND(
		mysql.OR(
			mysql.Bool(includeDeleted),
			tCitizensLabelsJob.DeletedAt.IS_NULL(),
		),
		tCitizenLabels.UserID.EQ(mysql.Int32(userId)),
	)
	if !userInfo.GetSuperuser() {
		jobAccessExists := s.labelsAccess.ACLAccessExistsCondition(
			tCitizensLabelsJob.ID,
			userInfo,
			int32(citizenslabels.AccessLevel_ACCESS_LEVEL_VIEW),
		)

		condition = condition.AND(jobAccessExists)
	}

	return s.GetUserLabels(ctx, s.db, condition)
}

func (s *Store) ValidateLabels(
	ctx context.Context,
	userJob string,
	labels []*citizenslabels.Label,
) (bool, error) {
	if len(labels) == 0 {
		return true, nil
	}

	idsExp := make([]mysql.Expression, len(labels))
	for i := range labels {
		// Remove access and settings info from passed in labels
		labels[i].Access = nil
		labels[i].Settings = nil

		idsExp[i] = mysql.Int64(labels[i].GetId())
	}

	stmt := tCitizensLabelsJob.
		SELECT(mysql.COUNT(tCitizensLabelsJob.ID).AS("data_count.total")).
		FROM(tCitizensLabelsJob).
		WHERE(mysql.AND(
			tCitizensLabelsJob.Job.EQ(mysql.String(userJob)),
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
