package jobsstore

import (
	"context"
	"errors"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	jobsconduct "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/conduct"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) CountConductEntries(ctx context.Context, db qrm.DB, q ConductQuery) (int64, error) {
	condition := mysql.AND(tConduct.Job.EQ(mysql.String(q.Job)))
	if !q.AllAccess {
		condition = condition.AND(tConduct.CreatorID.EQ(mysql.Int32(q.CreatorID)))
	}
	if !q.ShowDrafts {
		condition = condition.AND(tConduct.Draft.EQ(mysql.Bool(false)))
	}
	if len(q.IDs) > 0 {
		ids := make([]mysql.Expression, len(q.IDs))
		for i := range q.IDs {
			ids[i] = mysql.Int64(q.IDs[i])
		}
		condition = condition.AND(tConduct.ID.IN(ids...))
	}
	if len(q.Types) > 0 {
		ts := make([]mysql.Expression, len(q.Types))
		for i := range q.Types {
			ts[i] = mysql.Int32(int32(q.Types[i].Number()))
		}
		condition = condition.AND(tConduct.Type.IN(ts...))
	}
	if len(q.IDs) == 0 && !q.ShowExpired {
		condition = condition.AND(
			mysql.OR(tConduct.ExpiresAt.IS_NULL(), tConduct.ExpiresAt.GT_EQ(mysql.CURRENT_DATE())),
		)
	}
	if len(q.UserIDs) > 0 {
		ids := make([]mysql.Expression, len(q.UserIDs))
		for i := range q.UserIDs {
			ids[i] = mysql.Int32(q.UserIDs[i])
		}
		condition = condition.AND(tConduct.TargetUserID.IN(ids...))
	}

	countStmt := tConduct.
		SELECT(mysql.COUNT(tConduct.ID).AS("data_count.total")).
		FROM(tConduct).
		WHERE(condition)
	var count database.DataCount
	if err := countStmt.QueryContext(ctx, db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}
	return count.Total, nil
}

func (s *Store) ListConductEntries(
	ctx context.Context,
	db qrm.DB,
	q ConductQuery,
) ([]*jobsconduct.ConductEntry, error) {
	condition := mysql.AND(tConduct.Job.EQ(mysql.String(q.Job)))
	if !q.AllAccess {
		condition = condition.AND(tConduct.CreatorID.EQ(mysql.Int32(q.CreatorID)))
	}
	if !q.ShowDrafts {
		condition = condition.AND(tConduct.Draft.EQ(mysql.Bool(false)))
	}
	if len(q.IDs) > 0 {
		ids := make([]mysql.Expression, len(q.IDs))
		for i := range q.IDs {
			ids[i] = mysql.Int64(q.IDs[i])
		}
		condition = condition.AND(tConduct.ID.IN(ids...))
	}
	if len(q.Types) > 0 {
		ts := make([]mysql.Expression, len(q.Types))
		for i := range q.Types {
			ts[i] = mysql.Int32(int32(q.Types[i].Number()))
		}
		condition = condition.AND(tConduct.Type.IN(ts...))
	}
	if len(q.IDs) == 0 && !q.ShowExpired {
		condition = condition.AND(
			mysql.OR(tConduct.ExpiresAt.IS_NULL(), tConduct.ExpiresAt.GT_EQ(mysql.CURRENT_DATE())),
		)
	}
	if len(q.UserIDs) > 0 {
		ids := make([]mysql.Expression, len(q.UserIDs))
		for i := range q.UserIDs {
			ids[i] = mysql.Int32(q.UserIDs[i])
		}
		condition = condition.AND(tConduct.TargetUserID.IN(ids...))
	}

	orderBys := []mysql.OrderByClause{}
	if q.Sort != nil && len(q.Sort.GetColumns()) > 0 {
		for _, sc := range q.Sort.GetColumns() {
			var columns []mysql.Column
			switch sc.GetId() {
			case "type":
				columns = append(columns, tConduct.Type, tConduct.ID)
			case "id":
				fallthrough
			default:
				columns = append(columns, tConduct.ID)
			}
			for _, column := range columns {
				if sc.GetDesc() {
					orderBys = append(orderBys, column.DESC())
				} else {
					orderBys = append(orderBys, column.ASC())
				}
			}
		}
	} else {
		orderBys = append(orderBys, tConduct.ID.DESC())
	}

	tColleague := table.FivenetUser.AS("target_user")
	tUserUserProps := table.FivenetUserProps.AS("target_user_props")
	tColleagueAvatar := table.FivenetFiles.AS("target_user_profile_picture")
	tCreator := table.FivenetUser.AS("creator")
	tCreatorUserProps := table.FivenetUserProps.AS("creator_props")
	tCreatorAvatar := table.FivenetFiles.AS("creator_profile_picture")

	columns := mysql.ProjectionList{
		tConduct.CreatedAt,
		tConduct.UpdatedAt,
		tConduct.Job,
		tConduct.Type,
		tConduct.Draft,
		tConduct.Message,
		tConduct.ExpiresAt,
		tConduct.TargetUserID,
		tColleague.ID,
		tColleague.Job,
		tColleague.JobGrade,
		tColleague.Firstname,
		tColleague.Lastname,
		tColleague.Dateofbirth,
		tColleague.PhoneNumber,
		tUserUserProps.AvatarFileID.AS("target_user.profile_picture_file_id"),
		tColleagueAvatar.FilePath.AS("target_user.profile_picture"),
		tColleagueProps.UserID,
		tColleagueProps.Job,
		tColleagueProps.AbsenceBegin,
		tColleagueProps.AbsenceEnd,
		tColleagueProps.NamePrefix,
		tColleagueProps.NameSuffix,
		tConduct.CreatorID,
		tCreator.ID,
		tCreator.Job,
		tCreator.JobGrade,
		tCreator.Firstname,
		tCreator.Lastname,
		tCreator.Dateofbirth,
		tCreator.PhoneNumber,
		tCreatorUserProps.AvatarFileID.AS("creator.profile_picture_file_id"),
		tCreatorAvatar.FilePath.AS("creator.profile_picture"),
	}
	if q.AllAccess {
		columns = append(columns, tConduct.DeletedAt)
	}

	stmt := tConduct.
		SELECT(tConduct.ID, columns...).
		FROM(tConduct.
			LEFT_JOIN(tColleague, tColleague.ID.EQ(tConduct.TargetUserID)).
			LEFT_JOIN(tUserUserProps, tUserUserProps.UserID.EQ(tConduct.TargetUserID)).
			LEFT_JOIN(tColleagueProps, mysql.AND(tColleagueProps.UserID.EQ(tConduct.TargetUserID), tColleague.Job.EQ(mysql.String(q.Job)))).
			LEFT_JOIN(tCreator, tCreator.ID.EQ(tConduct.CreatorID)).
			LEFT_JOIN(tCreatorUserProps, tCreatorUserProps.UserID.EQ(tConduct.CreatorID)).
			LEFT_JOIN(tColleagueAvatar, tColleagueAvatar.ID.EQ(tUserUserProps.AvatarFileID)).
			LEFT_JOIN(tCreatorAvatar, tCreatorAvatar.ID.EQ(tCreatorUserProps.AvatarFileID)),
		).
		WHERE(condition).
		OFFSET(q.Offset).
		ORDER_BY(orderBys...).
		LIMIT(q.Limit)

	entries := []*jobsconduct.ConductEntry{}
	if err := stmt.QueryContext(ctx, db, &entries); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return entries, nil
}

func (s *Store) GetConductEntry(
	ctx context.Context,
	db qrm.DB,
	id int64,
) (*jobsconduct.ConductEntry, error) {
	tColleague := table.FivenetUser.AS("target_user")
	tCreator := tColleague.AS("creator")

	stmt := tConduct.
		SELECT(
			tConduct.ID,
			tConduct.CreatedAt,
			tConduct.UpdatedAt,
			tConduct.DeletedAt,
			tConduct.Job,
			tConduct.Type,
			tConduct.Draft,
			tConduct.Message,
			tConduct.ExpiresAt,
			tConduct.TargetUserID,
			tColleague.ID,
			tColleague.Firstname,
			tColleague.Lastname,
			tColleague.Dateofbirth,
			tColleague.PhoneNumber,
			tConduct.CreatorID,
			tCreator.ID,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tCreator.PhoneNumber,
		).
		FROM(tConduct.
			LEFT_JOIN(tColleague, tColleague.ID.EQ(tConduct.TargetUserID)).
			LEFT_JOIN(tCreator, tCreator.ID.EQ(tConduct.CreatorID)),
		).
		WHERE(tConduct.ID.EQ(mysql.Int64(id))).
		LIMIT(1)

	dest := &jobsconduct.ConductEntry{}
	if err := stmt.QueryContext(ctx, db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.GetId() == 0 {
		return nil, nil
	}

	return dest, nil
}

func (s *Store) CreateConductEntry(
	ctx context.Context,
	db qrm.DB,
	entry *jobsconduct.ConductEntry,
) (int64, error) {
	tConduct := table.FivenetJobConduct
	stmt := tConduct.
		INSERT(
			tConduct.Job,
			tConduct.Type,
			tConduct.Draft,
			tConduct.Message,
			tConduct.ExpiresAt,
			tConduct.TargetUserID,
			tConduct.CreatorID,
		).
		VALUES(
			entry.GetJob(),
			entry.GetType(),
			entry.GetDraft(),
			entry.GetMessage(),
			entry.GetExpiresAt(),
			dbutils.Int32P(entry.GetTargetUserId()),
			entry.GetCreatorId(),
		)

	res, err := stmt.ExecContext(ctx, db)
	if err != nil {
		return 0, err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

func (s *Store) UpdateConductEntry(
	ctx context.Context,
	db qrm.DB,
	entry *jobsconduct.ConductEntry,
) error {
	tConduct := table.FivenetJobConduct
	stmt := tConduct.
		UPDATE(
			tConduct.Type,
			tConduct.Draft,
			tConduct.Message,
			tConduct.ExpiresAt,
			tConduct.TargetUserID,
		).
		SET(
			entry.GetType(),
			entry.GetDraft(),
			entry.GetMessage(),
			entry.GetExpiresAt(),
			entry.GetTargetUserId(),
		).
		WHERE(mysql.AND(
			tConduct.Job.EQ(mysql.String(entry.GetJob())),
			tConduct.ID.EQ(mysql.Int64(entry.GetId())),
		)).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, db)
	return err
}

func (s *Store) DeleteConductEntry(
	ctx context.Context,
	db qrm.DB,
	job string,
	id int64,
	deletedAt *timestamp.Timestamp,
) error {
	tConduct := table.FivenetJobConduct
	stmt := tConduct.
		UPDATE(tConduct.DeletedAt).
		SET(tConduct.DeletedAt.SET(dbutils.TimestampToMySQL(deletedAt))).
		WHERE(mysql.AND(
			tConduct.Job.EQ(mysql.String(job)),
			tConduct.ID.EQ(mysql.Int64(id)),
		)).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, db)
	return err
}
