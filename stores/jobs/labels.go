package jobsstore

import (
	"context"
	"errors"
	"math"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	jobslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/labels"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
)

func (s *Store) GetColleagueLabels(
	ctx context.Context,
	db qrm.DB,
	job string,
	search string,
	includeDeleted bool,
) ([]*jobslabels.Label, error) {
	condition := mysql.AND(
		tJobLabels.Job.EQ(mysql.String(job)),
		mysql.OR(
			mysql.Bool(includeDeleted),
			tJobLabels.DeletedAt.IS_NULL(),
		),
	)

	if search = dbutils.PrepareForLikeSearch(search); search != "" {
		condition = condition.AND(tJobLabels.Name.LIKE(mysql.String(search)))
	}

	stmt := tJobLabels.
		SELECT(
			tJobLabels.ID,
			tJobLabels.Job,
			tJobLabels.DeletedAt,
			tJobLabels.Name,
			tJobLabels.Color,
			tJobLabels.Icon,
			tJobLabels.SortOrder,
		).
		FROM(tJobLabels).
		WHERE(condition).
		ORDER_BY(
			tJobLabels.SortOrder.ASC(),
			tJobLabels.SortKey.ASC(),
		)

	labels := []*jobslabels.Label{}
	if err := stmt.QueryContext(ctx, db, &labels); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return labels, nil
}

func (s *Store) GetUsersLabels(
	ctx context.Context,
	db qrm.DB,
	job string,
	userIds []int32,
	includeDeleted bool,
) ([]*UserLabels, error) {
	if len(userIds) == 0 {
		return []*UserLabels{}, nil
	}

	labels := make([]*UserLabels, 0, len(userIds))
	for _, userId := range userIds {
		userLabels, err := s.GetUserLabels(ctx, db, job, userId, includeDeleted)
		if err != nil {
			return nil, err
		}
		if len(userLabels.GetList()) == 0 {
			continue
		}

		labels = append(labels, &UserLabels{UserId: userId, Labels: userLabels})
	}

	return labels, nil
}

func (s *Store) GetLabel(
	ctx context.Context,
	db qrm.DB,
	job string,
	labelId int64,
	includeDeleted bool,
) (*jobslabels.Label, error) {
	stmt := tJobLabels.
		SELECT(
			tJobLabels.ID,
			tJobLabels.DeletedAt,
			tJobLabels.Job,
			tJobLabels.Name,
			tJobLabels.Color,
			tJobLabels.Icon,
			tJobLabels.SortOrder,
		).
		FROM(tJobLabels).
		WHERE(mysql.AND(
			tJobLabels.ID.EQ(mysql.Int64(labelId)),
			tJobLabels.Job.EQ(mysql.String(job)),
			mysql.OR(
				mysql.Bool(includeDeleted),
				tJobLabels.DeletedAt.IS_NULL(),
			),
		)).
		LIMIT(1)

	label := &jobslabels.Label{}
	if err := stmt.QueryContext(ctx, db, label); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if label.GetId() == 0 {
		return nil, nil
	}

	return label, nil
}

func (s *Store) NextLabelSortOrder(
	ctx context.Context,
	db qrm.Queryable,
	job string,
) (int32, error) {
	stmt := tJobLabels.
		SELECT(
			mysql.COALESCE(mysql.MAX(tJobLabels.SortOrder), mysql.Int32(-1)).
				AS("sort_order"),
		).
		FROM(tJobLabels).
		WHERE(mysql.AND(
			tJobLabels.Job.EQ(mysql.String(job)),
			tJobLabels.DeletedAt.IS_NULL(),
		))

	var dest struct {
		SortOrder int32 `alias:"sort_order"`
	}
	if err := stmt.QueryContext(ctx, db, &dest); err != nil {
		return 0, err
	}

	return dest.SortOrder + 1, nil
}

func (s *Store) UpdateLabel(
	ctx context.Context,
	db qrm.DB,
	label *jobslabels.Label,
	job string,
) error {
	tJobLabels := table.FivenetJobLabels
	stmt := tJobLabels.
		UPDATE(
			tJobLabels.Name,
			tJobLabels.Color,
			tJobLabels.Icon,
		).
		SET(
			tJobLabels.Name.SET(mysql.String(label.GetName())),
			tJobLabels.Color.SET(mysql.String(label.GetColor())),
			tJobLabels.Icon.SET(dbutils.StringEmpty(label.GetIcon())),
		).
		WHERE(mysql.AND(
			tJobLabels.ID.EQ(mysql.Int64(label.GetId())),
			tJobLabels.Job.EQ(mysql.String(job)),
		))

	_, err := stmt.ExecContext(ctx, db)
	return err
}

func (s *Store) InsertLabel(
	ctx context.Context,
	db qrm.DB,
	label *jobslabels.Label,
) (int64, error) {
	tJobLabels := table.FivenetJobLabels
	stmt := tJobLabels.
		INSERT(
			tJobLabels.Job,
			tJobLabels.SortOrder,
			tJobLabels.Name,
			tJobLabels.Color,
			tJobLabels.Icon,
		).
		VALUES(
			label.Job,
			label.GetSortOrder(),
			label.GetName(),
			label.GetColor(),
			dbutils.StringEmpty(label.GetIcon()),
		)

	result, err := stmt.ExecContext(ctx, db)
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
	db qrm.DB,
	job string,
	labelId int64,
	deletedAt *timestamp.Timestamp,
) error {
	tJobLabels := table.FivenetJobLabels
	stmt := tJobLabels.
		UPDATE(
			tJobLabels.DeletedAt,
		).
		SET(
			tJobLabels.DeletedAt.SET(dbutils.TimestampToMySQL(deletedAt)),
		).
		WHERE(mysql.AND(
			tJobLabels.ID.EQ(mysql.Int64(labelId)),
			tJobLabels.Job.EQ(mysql.String(job)),
		)).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, db)
	return err
}

func (s *Store) ReorderLabels(
	ctx context.Context,
	job string,
	labelIds []int64,
) error {
	tJobLabels := table.FivenetJobLabels
	stmt := tJobLabels.
		SELECT(tJobLabels.ID).
		FROM(tJobLabels).
		WHERE(mysql.AND(
			tJobLabels.Job.EQ(mysql.String(job)),
			tJobLabels.DeletedAt.IS_NULL(),
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

		if _, err := tJobLabels.
			UPDATE().
			SET(
				tJobLabels.SortOrder.SET(mysql.Int32(int32(idx))),
			).
			WHERE(mysql.AND(
				tJobLabels.ID.EQ(mysql.Int64(labelID)),
				tJobLabels.Job.EQ(mysql.String(job)),
				tJobLabels.DeletedAt.IS_NULL(),
			)).
			LIMIT(1).
			ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Store) GetColleagueLabelsStats(
	ctx context.Context,
	db qrm.DB,
	job string,
) ([]*jobslabels.LabelCount, error) {
	tColleague := table.FivenetUser.AS("user")

	stmt := tColleagueLabels.
		SELECT(
			mysql.COUNT(tColleagueLabels.LabelID).AS("label_count.count"),
			tJobLabels.ID,
			tJobLabels.Job,
			tJobLabels.Name,
			tJobLabels.Color,
			tJobLabels.Icon,
		).
		FROM(
			tColleagueLabels.
				INNER_JOIN(tJobLabels,
					tJobLabels.ID.EQ(tColleagueLabels.LabelID),
				).
				INNER_JOIN(tColleague,
					tColleague.ID.EQ(tColleagueLabels.UserID),
				).
				INNER_JOIN(tUserJobs,
					mysql.AND(
						tUserJobs.UserID.EQ(tColleagueLabels.UserID),
						tUserJobs.Job.EQ(mysql.String(job)),
					),
				),
		).
		WHERE(mysql.AND(
			tJobLabels.Job.EQ(mysql.String(job)),
			tJobLabels.DeletedAt.IS_NULL(),
			tColleagueLabels.Job.EQ(mysql.String(job)),
			tColleague.Job.EQ(mysql.String(job)),
		)).
		GROUP_BY(tJobLabels.ID).
		ORDER_BY(
			tJobLabels.SortOrder.ASC(),
			tJobLabels.SortKey.ASC(),
		)

	dest := []*jobslabels.LabelCount{}
	if err := stmt.QueryContext(ctx, db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}

func (s *Store) ValidateLabels(
	ctx context.Context,
	db qrm.DB,
	job string,
	labels []*jobslabels.Label,
) (bool, error) {
	if len(labels) == 0 {
		return true, nil
	}

	idsExp := make([]mysql.Expression, len(labels))
	for i := range labels {
		idsExp[i] = mysql.Int64(labels[i].GetId())
	}

	stmt := tJobLabels.
		SELECT(mysql.COUNT(tJobLabels.ID).AS("data_count.total")).
		FROM(tJobLabels).
		WHERE(mysql.AND(
			tJobLabels.Job.EQ(mysql.String(job)),
			tJobLabels.DeletedAt.IS_NULL(),
			tJobLabels.ID.IN(idsExp...),
		)).
		LIMIT(10)

	var count database.DataCount
	if err := stmt.QueryContext(ctx, db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return false, err
		}
	}

	return len(labels) == int(count.Total), nil
}
