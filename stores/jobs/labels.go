package jobsstore

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	jobslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/labels"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) GetColleagueLabels(
	ctx context.Context,
	db qrm.DB,
	job string,
	search string,
) ([]*jobslabels.Label, error) {
	condition := mysql.AND(
		tJobLabels.Job.EQ(mysql.String(job)),
		tJobLabels.DeletedAt.IS_NULL(),
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
) ([]*UserLabels, error) {
	if len(userIds) == 0 {
		return []*UserLabels{}, nil
	}

	labels := make([]*UserLabels, 0, len(userIds))
	for _, userId := range userIds {
		userLabels, err := s.GetUserLabels(ctx, db, job, userId)
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

func (s *Store) ManageLabels(
	ctx context.Context,
	db qrm.DB,
	job string,
	labels []*jobslabels.Label,
) ([]*jobslabels.Label, error) {
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
			tJobLabels.Job.EQ(mysql.String(job)),
		))

	current := []*jobslabels.Label{}
	if err := stmt.QueryContext(ctx, db, &current); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	_, removed := utils.SlicesDifferenceFunc(current, labels, func(in *jobslabels.Label) int64 {
		return in.GetId()
	})

	var i int32
	for _, label := range labels {
		label.Job = &job
		label.SortOrder = i
		i++
	}

	if len(labels) > 0 {
		toCreate := []*jobslabels.Label{}
		toUpdate := []*jobslabels.Label{}
		for _, label := range labels {
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
					tJobLabels.Icon,
					tJobLabels.SortOrder,
					tJobLabels.DeletedAt,
				).
				MODELS(toCreate).
				ON_DUPLICATE_KEY_UPDATE(
					tJobLabels.Name.SET(mysql.RawString("VALUES(`name`)")),
					tJobLabels.Color.SET(mysql.RawString("VALUES(`color`)")),
					tJobLabels.Icon.SET(mysql.RawString("VALUES(`icon`)")),
					tJobLabels.SortOrder.SET(mysql.RawInt("VALUES(`sort_order`)")),
					tJobLabels.DeletedAt.SET(mysql.TimestampExp(mysql.NULL)),
				)

			if _, err := insertStmt.ExecContext(ctx, db); err != nil {
				return nil, err
			}
		}

		if len(toUpdate) > 0 {
			for _, label := range toUpdate {
				updateStmt := tJobLabels.
					UPDATE(
						tJobLabels.Name,
						tJobLabels.Color,
						tJobLabels.Icon,
						tJobLabels.SortOrder,
						tJobLabels.DeletedAt,
					).
					SET(
						tJobLabels.Name.SET(mysql.String(label.GetName())),
						tJobLabels.Color.SET(mysql.String(label.GetColor())),
						tJobLabels.Icon.SET(dbutils.StringEmpty(label.GetIcon())),
						tJobLabels.SortOrder.SET(mysql.Int32(label.GetSortOrder())),
						tJobLabels.DeletedAt.SET(mysql.TimestampExp(mysql.NULL)),
					).
					WHERE(mysql.AND(
						tJobLabels.ID.EQ(mysql.Int64(label.GetId())),
						tJobLabels.Job.EQ(mysql.String(label.GetJob())),
					)).
					LIMIT(1)

				if _, err := updateStmt.ExecContext(ctx, db); err != nil {
					return nil, err
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
			SET(tJobLabels.DeletedAt.SET(mysql.CURRENT_TIMESTAMP())).
			WHERE(mysql.AND(
				tJobLabels.ID.IN(ids...),
				tJobLabels.Job.EQ(mysql.String(job)),
			)).
			LIMIT(int64(len(removed)))

		if _, err := deleteStmt.ExecContext(ctx, db); err != nil {
			return nil, err
		}
	}

	updated := []*jobslabels.Label{}
	if err := stmt.QueryContext(ctx, db, &updated); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return updated, nil
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
