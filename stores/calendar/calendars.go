package calendarstore

import (
	"context"
	"errors"

	calendarresource "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar"
	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) CountCalendars(ctx context.Context, q ListQuery) (int64, error) {
	tCreator := table.FivenetUser.AS("creator")
	stmt := s.countCalendarsStmt(q, tCreator)

	var count database.DataCount
	if err := stmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return count.Total, nil
}

func (s *Store) ListCalendars(
	ctx context.Context,
	q ListQuery,
	offset, limit int64,
) ([]*calendarresource.Calendar, error) {
	tCreator := table.FivenetUser.AS("creator")
	tAvatar := table.FivenetFiles.AS("profile_picture")
	var userID int32
	if q.UserInfo != nil {
		userID = q.UserInfo.GetUserId()
	}

	stmt := s.listCalendarsStmt(q, userID, tCreator, tAvatar, offset, limit)

	var calendars []*calendarresource.Calendar
	if err := stmt.QueryContext(ctx, s.db, &calendars); err != nil {
		return nil, err
	}

	if !q.OnlyPublic {
		// FIXME temporary solution to have the access field populated for clients to decide user access
		// E.g., batching the look up(s) would be fine as well instead of doing one lookup per calendar
		for _, calendar := range calendars {
			access, err := s.ListTargetAccess(
				ctx,
				calendar.GetId(),
				birthdaySyncSubjectAccessOptions,
			)
			if err != nil {
				return nil, err
			}
			calendar.SetAccess(access)
		}
	}

	return calendars, nil
}

func (s *Store) GetCalendar(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	condition mysql.BoolExpression,
) (*calendarresource.Calendar, error) {
	stmt := s.getCalendarStmt(userInfo, condition)

	dest := &calendarresource.Calendar{}
	if err := stmt.QueryContext(ctx, s.db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.GetId() == 0 {
		return nil, nil
	}

	return dest, nil
}

func (s *Store) GetAccessibleCalendar(
	ctx context.Context,
	calendarID int64,
	userInfo *userinfo.UserInfo,
	accessLevel calendaraccess.AccessLevel,
	publicOk bool,
) (*calendarresource.Calendar, error) {
	check, err := s.CheckIfUserHasAccessToCalendar(ctx, calendarID, userInfo, accessLevel, publicOk)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, nil
	}

	return s.GetCalendar(ctx, userInfo, tCalendar.ID.EQ(mysql.Int64(calendarID)))
}

func (s *Store) countCalendarsStmt(
	q ListQuery,
	tCreator *table.FivenetUserTable,
) mysql.SelectStatement {
	condition, _ := s.listConditions(q)

	return tCalendar.
		SELECT(
			mysql.COUNT(mysql.DISTINCT(tCalendar.ID)).AS("data_count.total"),
		).
		FROM(tCalendar.
			LEFT_JOIN(tCreator,
				tCalendar.CreatorID.EQ(tCreator.ID),
			),
		).
		WHERE(condition)
}

func (s *Store) listConditions(
	q ListQuery,
) (mysql.BoolExpression, []mysql.OrderByClause) {
	includeDeleted := q.UserInfo != nil && q.UserInfo.GetJobAdmin()
	condition := mysql.Bool(includeDeleted).OR(tCalendar.DeletedAt.IS_NULL())
	if q.OnlyPublic {
		return condition.AND(tCalendar.Public.IS_TRUE()),
			[]mysql.OrderByClause{
				tCalendar.Name.ASC(),
			}
	}

	if q.UserInfo == nil {
		return condition, []mysql.OrderByClause{tCalendar.Name.ASC()}
	}

	subsCondition := tCalendar.ID.IN(tCalendarSubs.
		SELECT(
			tCalendarSubs.CalendarID,
		).
		FROM(tCalendarSubs).
		WHERE(mysql.AND(
			tCalendarSubs.UserID.EQ(mysql.Int32(q.UserInfo.GetUserId())),
		)),
	)

	minAccessLevel := calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW
	if q.MinAccessLevel != nil {
		minAccessLevel = *q.MinAccessLevel
		subsCondition = mysql.Bool(false)
	}

	var accessExists mysql.BoolExpression
	if !q.UserInfo.GetJobAdmin() {
		accessExists = s.access.ACLAccessExistsCondition(
			tCalendar.ID,
			q.UserInfo,
			int32(minAccessLevel),
		)
	} else {
		accessExists = mysql.OR(
			tCalendar.SystemKind.IS_NULL(),
			tCalendar.SystemKind.NOT_EQ(
				mysql.Int32(
					int32(calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS),
				),
			),
			s.birthdayCalendarVisible(tCalendar.ID, minAccessLevel, q.UserInfo),
		)
	}

	orderBys := []mysql.OrderByClause{tCalendar.Name.ASC()}
	creatorPrivateCondition := mysql.AND(
		tCalendar.Job.IS_NULL(),
		tCalendar.CreatorID.EQ(mysql.Int32(q.UserInfo.GetUserId())),
	)
	condition = mysql.AND(
		condition,
		mysql.OR(
			subsCondition,
			creatorPrivateCondition,
			accessExists,
			s.birthdayCalendarVisible(tCalendar.ID, minAccessLevel, q.UserInfo),
		),
	)

	if q.After != nil {
		condition = condition.AND(
			tCalendar.UpdatedAt.GT_EQ(mysql.TimestampT(q.After.AsTime())),
		)
	}

	if len(q.CalendarIDs) > 0 {
		calendarIDs := []mysql.Expression{}
		for _, v := range q.CalendarIDs {
			calendarIDs = append(calendarIDs, mysql.Int64(v))
		}

		orderBys = append(orderBys, tCalendar.ID.IN(calendarIDs...).DESC())
	}

	return condition, orderBys
}

func (s *Store) birthdayCalendarVisible(
	calendarID mysql.IntegerExpression,
	accessLevel calendaraccess.AccessLevel,
	userInfo *userinfo.UserInfo,
) mysql.BoolExpression {
	return mysql.AND(
		tCalendar.SystemKind.EQ(
			mysql.Int32(
				int32(calendarresource.CalendarSystemKind_CALENDAR_SYSTEM_KIND_JOB_BIRTHDAYS),
			),
		),
		mysql.OR(
			mysql.Bool(userInfo.GetJobAdmin()),
			tCalendar.Job.EQ(mysql.String(userInfo.GetJob())),
		),
		s.access.ACLAccessExistsCondition(calendarID, userInfo, int32(accessLevel)),
	)
}

func (s *Store) listCalendarsStmt(
	q ListQuery,
	userID int32,
	tCreator *table.FivenetUserTable,
	tAvatar *table.FivenetFilesTable,
	offset, limit int64,
) mysql.SelectStatement {
	condition, orderBys := s.listConditions(q)

	selectColumns := []mysql.Projection{
		tCalendar.ID,
		tCalendar.CreatedAt,
		tCalendar.UpdatedAt,
		tCalendar.DeletedAt,
		tCalendar.Job,
		tCalendar.DiscordSettings,
		tCalendar.Name,
		tCalendar.Description,
		tCalendar.Public,
		tCalendar.Closed,
		tCalendar.Color,
		tCalendar.Icon,
		tCalendar.CreatorID,
		tCreator.ID,
		tCreator.Job,
		tCreator.JobGrade,
		tCreator.Firstname,
		tCreator.Lastname,
		tCreator.Dateofbirth,
		tCreator.PhoneNumber,
		tUserProps.AvatarFileID.AS("creator.profile_picture_file_id"),
		tAvatar.FilePath.AS("creator.profile_picture"),
		tCalendarSubs.CalendarID,
		tCalendarSubs.UserID,
		tCalendarSubs.CreatedAt,
		tCalendarSubs.Confirmed,
		tCalendarSubs.Muted,
	}
	selectColumns = append(selectColumns, tCalendar.SystemKind)

	return tCalendar.
		SELECT(selectColumns[0], selectColumns[1:]...).
		FROM(tCalendar.
			LEFT_JOIN(tCreator,
				tCalendar.CreatorID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tCalendar.CreatorID),
			).
			LEFT_JOIN(tCalendarSubs,
				mysql.AND(
					tCalendarSubs.CalendarID.EQ(tCalendar.ID),
					tCalendarSubs.UserID.EQ(mysql.Int32(userID)),
				),
			).
			LEFT_JOIN(tAvatar,
				tAvatar.ID.EQ(tUserProps.AvatarFileID),
			),
		).
		WHERE(condition).
		ORDER_BY(orderBys...).
		OFFSET(offset).
		LIMIT(limit)
}

func (s *Store) getCalendarStmt(
	userInfo *userinfo.UserInfo,
	condition mysql.BoolExpression,
) mysql.SelectStatement {
	tCreator := table.FivenetUser.AS("creator")
	tAvatar := table.FivenetFiles.AS("profile_picture")
	var userID int32
	if userInfo != nil {
		userID = userInfo.GetUserId()
	}
	includeDeleted := userInfo != nil && userInfo.GetJobAdmin()

	columns := []mysql.Projection{
		tCalendar.ID,
		tCalendar.CreatedAt,
		tCalendar.UpdatedAt,
		tCalendar.DeletedAt,
		tCalendar.Job,
		tCalendar.DiscordSettings,
		tCalendar.Name,
		tCalendar.Description,
		tCalendar.Public,
		tCalendar.Closed,
		tCalendar.Color,
		tCalendar.Icon,
		tCalendar.CreatorID,
		tCalendar.CreatorJob,
		tCreator.ID,
		tCreator.Job,
		tCreator.JobGrade,
		tCreator.Firstname,
		tCreator.Lastname,
		tCreator.Dateofbirth,
		tCreator.PhoneNumber,
		tUserProps.AvatarFileID.AS("creator.profile_picture_file_id"),
		tAvatar.FilePath.AS("creator.profile_picture"),
		tCalendarSubs.CalendarID,
		tCalendarSubs.UserID,
		tCalendarSubs.CreatedAt,
		tCalendarSubs.Confirmed,
		tCalendarSubs.Muted,
	}
	columns = append(columns, tCalendar.SystemKind)

	return tCalendar.
		SELECT(columns[0], columns[1:]...).
		FROM(tCalendar.
			LEFT_JOIN(tCreator,
				tCalendar.CreatorID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tCalendar.CreatorID),
			).
			LEFT_JOIN(tCalendarSubs,
				mysql.AND(
					tCalendarSubs.CalendarID.EQ(tCalendar.ID),
					tCalendarSubs.UserID.EQ(mysql.Int32(userID)),
				),
			).
			LEFT_JOIN(tAvatar,
				tAvatar.ID.EQ(tUserProps.AvatarFileID),
			),
		).
		WHERE(mysql.AND(
			mysql.OR(
				mysql.Bool(includeDeleted),
				tCalendar.DeletedAt.IS_NULL(),
			),
			condition,
		)).
		LIMIT(1)
}

func (s *Store) CreateCalendar(
	ctx context.Context,
	tx qrm.DB,
	cal *calendarresource.Calendar,
	userInfo *userinfo.UserInfo,
) (int64, error) {
	tCalendar := table.FivenetCalendar
	systemKind := mysql.IntExp(mysql.NULL)
	if cal.HasSystemKind() {
		systemKind = mysql.IntExp(mysql.Int32(int32(cal.GetSystemKind())))
	}
	var existing struct{ ID int64 }
	if cal.HasSystemKind() {
		stmt := tCalendar.
			SELECT(tCalendar.ID.AS("id")).
			FROM(tCalendar).
			WHERE(mysql.AND(
				tCalendar.Job.EQ(mysql.String(cal.GetJob())),
				tCalendar.SystemKind.EQ(mysql.Int32(int32(cal.GetSystemKind()))),
			)).
			LIMIT(1)

		if err := stmt.QueryContext(
			ctx,
			tx,
			&existing,
		); err != nil &&
			!errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	if existing.ID > 0 {
		if _, err := tCalendar.
			UPDATE(
				tCalendar.SystemKind,
				tCalendar.DiscordSettings,
				tCalendar.Name,
				tCalendar.Description,
				tCalendar.Public,
				tCalendar.Closed,
				tCalendar.Color,
				tCalendar.Icon,
			).
			SET(
				systemKind,
				cal.GetDiscordSettings(),
				cal.GetName(),
				cal.GetDescription(),
				cal.GetPublic(),
				cal.GetClosed(),
				cal.GetColor(),
				dbutils.StringEmpty(cal.GetIcon()),
			).
			WHERE(tCalendar.ID.EQ(mysql.Int64(existing.ID))).
			LIMIT(1).
			ExecContext(ctx, tx); err != nil {
			return 0, err
		}
		return existing.ID, nil
	}

	res, err := tCalendar.
		INSERT(
			tCalendar.Job,
			tCalendar.SystemKind,
			tCalendar.DiscordSettings,
			tCalendar.Name,
			tCalendar.Description,
			tCalendar.Public,
			tCalendar.Closed,
			tCalendar.Color,
			tCalendar.Icon,
			tCalendar.CreatorID,
			tCalendar.CreatorJob,
		).
		VALUES(
			cal.GetJob(),
			systemKind,
			cal.GetDiscordSettings(),
			cal.GetName(),
			cal.GetDescription(),
			cal.GetPublic(),
			cal.GetClosed(),
			cal.GetColor(),
			dbutils.StringEmpty(cal.GetIcon()),
			userInfo.GetUserId(),
			userInfo.GetJob(),
		).
		ExecContext(ctx, tx)
	if err != nil {
		return 0, err
	}

	if cal.GetId() > 0 {
		return cal.GetId(), nil
	}
	return res.LastInsertId()
}

func (s *Store) UpdateCalendar(
	ctx context.Context,
	tx qrm.DB,
	cal *calendarresource.Calendar,
) error {
	tCalendar := table.FivenetCalendar
	stmt := tCalendar.
		UPDATE(
			tCalendar.DiscordSettings,
			tCalendar.Name,
			tCalendar.Description,
			tCalendar.Public,
			tCalendar.Closed,
			tCalendar.Color,
			tCalendar.Icon,
		).
		SET(
			cal.GetDiscordSettings(),
			cal.GetName(),
			cal.GetDescription(),
			cal.GetPublic(),
			cal.GetClosed(),
			cal.GetColor(),
			dbutils.StringEmpty(cal.GetIcon()),
		).
		WHERE(mysql.AND(
			tCalendar.ID.EQ(mysql.Int64(cal.GetId())),
		)).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) DeleteCalendar(
	ctx context.Context,
	tx qrm.DB,
	calendarID int64,
	deletedAt *timestamp.Timestamp,
) error {
	tCalendar := table.FivenetCalendar
	stmt := tCalendar.
		UPDATE(
			tCalendar.DeletedAt,
		).
		SET(
			tCalendar.DeletedAt.SET(dbutils.TimestampToMySQL(deletedAt)),
		).
		WHERE(tCalendar.ID.EQ(mysql.Int64(calendarID))).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}
