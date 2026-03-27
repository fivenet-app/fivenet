package colleaguesstore

import (
	"context"
	"errors"

	jobslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/labels"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func GetUserLabels(
	ctx context.Context,
	tx qrm.DB,
	job string,
	userId int32,
) (*jobslabels.Labels, error) {
	tJobLabels := table.FivenetJobLabels.AS("label")
	tUserLabels := table.FivenetJobColleagueLabels

	stmt := tUserLabels.
		SELECT(
			tJobLabels.ID,
			tJobLabels.Job,
			tJobLabels.Name,
			tJobLabels.Color,
			tJobLabels.Icon,
		).
		FROM(
			tUserLabels.
				INNER_JOIN(tJobLabels,
					tJobLabels.ID.EQ(tUserLabels.LabelID),
				),
		).
		WHERE(mysql.AND(
			tUserLabels.UserID.EQ(mysql.Int32(userId)),
			tJobLabels.Job.EQ(mysql.String(job)),
			tJobLabels.DeletedAt.IS_NULL(),
		)).
		ORDER_BY(
			tJobLabels.Order.ASC(),
		)

	list := &jobslabels.Labels{
		List: []*jobslabels.Label{},
	}
	if err := stmt.QueryContext(ctx, tx, &list.List); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return list, nil
}
