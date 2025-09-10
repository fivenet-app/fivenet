package citizens

import (
	"context"
	"database/sql"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	pb "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/citizens"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type UserActivityStore struct {
	db     *sql.DB
	sorter *database.Builder
}

func New(db *sql.DB) *UserActivityStore {
	return &UserActivityStore{
		db:     db,
		sorter: newUserActivitySorter(),
	}
}

func newUserActivitySorter() *database.Builder {
	tUserActivity := table.FivenetUserActivity.AS("user_activity")

	return database.New(
		database.SpecMap{
			"createdAt": database.Column{Col: tUserActivity.CreatedAt, NullsLast: true},
		},
		[]mysql.OrderByClause{
			tUserActivity.CreatedAt.DESC().NULLS_LAST(),
		},
		[]mysql.OrderByClause{
			tUserActivity.ID.DESC(),
		},
		3,
	)
}

func (s *UserActivityStore) List(
	ctx context.Context,
	req *pb.ListUserActivityRequest,
) ([]*users.UserActivity, error) {
	tUserActivity := table.FivenetUserActivity.AS("user_activity")

	orderBys := s.sorter.Build(req.GetSort())

	stmt := mysql.
		SELECT(
			tUserActivity.AllColumns,
		).
		FROM(tUserActivity).
		ORDER_BY(orderBys...)

	var activities []*users.UserActivity
	if err := stmt.QueryContext(ctx, s.db, &activities); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return activities, nil
}
