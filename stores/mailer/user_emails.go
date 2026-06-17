package mailerstore

import (
	"context"
	"errors"

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
	tEmailAlias := table.FivenetMailerEmails.AS("email")

	condition := mysql.Bool(includeDeleted).OR(tEmailAlias.DeletedAt.IS_NULL())
	if !includeDisabled {
		condition = condition.AND(tEmailAlias.Deactivated.IS_FALSE())
	}

	if !userInfo.GetSuperuser() {
		acl := s.subjectAccess.ACLAccessExistsCondition(
			tEmailAlias.ID,
			userInfo,
			int32(maileraccess.AccessLevel_ACCESS_LEVEL_READ),
		)

		condition = condition.AND(mysql.OR(
			tEmailAlias.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
			acl,
		))
	}

	stmt := tEmailAlias.
		SELECT(
			tEmailAlias.ID,
			tEmailAlias.CreatedAt,
			tEmailAlias.UpdatedAt,
			tEmailAlias.DeletedAt,
			tEmailAlias.Deactivated,
			tEmailAlias.Job,
			tEmailAlias.UserID,
			tEmailAlias.Email,
			tEmailAlias.EmailChanged,
			tEmailAlias.Label,
		).
		FROM(tEmailAlias).
		WHERE(condition).
		ORDER_BY(tEmailAlias.Job.ASC(), tEmailAlias.Label.ASC())

	if pag != nil {
		stmt = stmt.OFFSET(pag.GetOffset())
	}

	emails := []*maileremails.Email{}
	if err := stmt.QueryContext(ctx, s.dbOr(db), &emails); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return emails, nil
}
