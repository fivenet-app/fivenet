package mailerstore

import (
	"context"
	"errors"
	"fmt"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	maileraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/access"
	maileremails "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/emails"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) ListUserEmails(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	pag *database.PaginationRequest,
	includeDisabled bool,
	includeDeleted bool,
) ([]*maileremails.Email, error) {
	condition := mysql.Bool(true)
	if !includeDisabled {
		condition = condition.AND(table.FivenetMailerEmails.Deactivated.IS_FALSE())
	}

	visibleIDs := s.subjectAccess.VisibleIDsByConditionQuery(
		userInfo,
		int32(maileraccess.AccessLevel_ACCESS_LEVEL_READ),
		includeDeleted,
		condition,
	)
	ctes := visibleIDs.CTEs
	visibleEmailID := mysql.IntegerColumn("id").From(visibleIDs.Table)

	stmt := tEmails.
		SELECT(
			tEmails.ID,
			tEmails.CreatedAt,
			tEmails.UpdatedAt,
			tEmails.DeletedAt,
			tEmails.Deactivated,
			tEmails.Job,
			tEmails.UserID,
			tEmails.Email,
			tEmails.EmailChanged,
			tEmails.Label,
		).
		FROM(
			visibleIDs.Table.
				INNER_JOIN(tEmails,
					tEmails.ID.EQ(visibleEmailID),
				),
		).
		ORDER_BY(tEmails.Job.ASC(), tEmails.Label.ASC())

	if pag != nil {
		stmt = stmt.OFFSET(pag.GetOffset())
	}

	var finalStmt mysql.Statement = stmt
	if len(ctes) > 0 {
		finalStmt = mysql.WITH(ctes...)(stmt)
	}
	fmt.Println(finalStmt.DebugSql())

	emails := []*maileremails.Email{}
	if err := finalStmt.QueryContext(ctx, s.dbOr(db), &emails); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return emails, nil
}
