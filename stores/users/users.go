package users

import (
	"context"
	"database/sql"

	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
)

type UserStore struct {
	db *sql.DB
}

// New creates a new UserStore instance.
func New(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) GetShort(ctx context.Context, userId int32) (*usershort.UserShort, error) {
	tUsers := table.FivenetUser.AS("user_short")

	stmt := tUsers.
		SELECT(
			tUsers.ID.AS("user_short.userid"),
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Dateofbirth,
		).
		FROM(tUsers).
		WHERE(tUsers.ID.EQ(mysql.Int32(userId))).
		LIMIT(1)

	user := &usershort.UserShort{}
	if err := stmt.QueryContext(ctx, s.db, user); err != nil {
		return nil, err
	}

	return user, nil
}
