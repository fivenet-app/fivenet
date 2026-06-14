package completor

import (
	"context"
	"errors"

	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type CitizensQuery struct {
	Search      string
	CurrentJob  bool
	UserJob     string
	UserIDs     []int32
	UserIDsOnly bool
}

func (s *Store) CompleteCitizens(
	ctx context.Context,
	q CitizensQuery,
) ([]*usershort.UserShort, error) {
	tUsers := table.FivenetUser.AS("user_short")

	orderBys := []mysql.OrderByClause{}
	condition := s.customDB.Conditions.User.GetFilter(tUsers.Alias())

	if q.CurrentJob {
		condition = condition.AND(tUsers.Job.EQ(mysql.String(q.UserJob)))
	}

	if search := dbutils.PrepareForLikeSearch(q.Search); search != "" {
		condition = condition.AND(
			mysql.CONCAT(tUsers.Firstname, mysql.String(" "), tUsers.Lastname).
				LIKE(mysql.String(search)),
		)
	}

	if len(q.UserIDs) > 0 {
		userIDs := []mysql.Expression{}
		for _, v := range q.UserIDs {
			userIDs = append(userIDs, mysql.Int32(v))
		}

		if q.UserIDsOnly {
			condition = condition.OR(tUsers.ID.IN(userIDs...))
		}

		// Make sure to sort by the user IDs if provided
		orderBys = append(orderBys, tUsers.ID.IN(userIDs...).DESC())
	}

	orderBys = append(orderBys, tUsers.Lastname.ASC())

	columns := mysql.ProjectionList{
		tUsers.ID,
		tUsers.Firstname,
		tUsers.Lastname,
		tUsers.Dateofbirth,
	}
	if q.CurrentJob {
		columns = append(columns, tUsers.Job, tUsers.JobGrade)
	}

	stmt := tUsers.
		SELECT(
			columns[0],
			columns[1:]...,
		).
		OPTIMIZER_HINTS(mysql.OptimizerHint("idx_users_firstname_lastname_fulltext")).
		FROM(tUsers).
		WHERE(condition).
		ORDER_BY(orderBys...).
		LIMIT(15)

	var dest []*usershort.UserShort
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}
