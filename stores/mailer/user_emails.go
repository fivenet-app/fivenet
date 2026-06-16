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
) ([]*maileremails.Email, error) {
	if userInfo == nil {
		userInfo = &userinfo.UserInfo{}
	}

	tEmail := table.FivenetMailerEmails
	tEmailAlias := table.FivenetMailerEmails.AS("email")

	condition := mysql.Bool(true)
	baseCondition := tEmail.DeletedAt.IS_NULL()
	if !includeDisabled {
		baseCondition = baseCondition.AND(tEmail.Deactivated.IS_FALSE())
	}

	if !userInfo.GetSuperuser() {
		acl := s.subjectAccess.ACLAccessExistsCondition(
			tEmail.ID,
			userInfo,
			int32(maileraccess.AccessLevel_ACCESS_LEVEL_READ),
		)

		condition = condition.AND(baseCondition.AND(mysql.OR(
			tEmail.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
			acl,
		)))
	} else {
		condition = condition.AND(baseCondition)
	}

	visibleStmt := s.subjectAccess.VisibleIDsByConditionStatement(
		userInfo,
		int32(maileraccess.AccessLevel_ACCESS_LEVEL_READ),
		condition,
	)
	var visibleIDs []int64
	if err := visibleStmt.QueryContext(ctx, s.dbOr(db), &visibleIDs); err != nil {
		return nil, err
	}

	if len(visibleIDs) == 0 {
		return []*maileremails.Email{}, nil
	}

	ids := make([]mysql.Expression, len(visibleIDs))
	for i := range visibleIDs {
		ids[i] = mysql.Int64(visibleIDs[i])
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
		WHERE(tEmailAlias.ID.IN(ids...)).
		ORDER_BY(tEmailAlias.Job.ASC(), tEmailAlias.Label.ASC())

	if pag != nil {
		stmt = stmt.OFFSET(pag.GetOffset())
		stmt = stmt.LIMIT(int64(len(visibleIDs)))
	}

	emails := []*maileremails.Email{}
	if err := stmt.QueryContext(ctx, s.dbOr(db), &emails); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return emails, nil
}
