package jobsstore

import (
	"context"
	"errors"
	"strings"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
	colleaguesactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues/activity"
	jobslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/labels"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/protobuf/proto"
)

func (s *Store) CreateColleagueActivity(
	ctx context.Context,
	db qrm.DB,
	activities ...*colleaguesactivity.ColleagueActivity,
) error {
	if len(activities) == 0 {
		return nil
	}

	tColleagueActivity := table.FivenetJobColleagueActivity
	stmt := tColleagueActivity.
		INSERT(
			tColleagueActivity.Job,
			tColleagueActivity.SourceUserID,
			tColleagueActivity.TargetUserID,
			tColleagueActivity.ActivityType,
			tColleagueActivity.Reason,
			tColleagueActivity.Data,
		)

	for _, activity := range activities {
		stmt = stmt.VALUES(
			activity.GetJob(),
			activity.SourceUserId,
			activity.TargetUserId,
			activity.GetActivityType(),
			activity.GetReason(),
			activity.Data,
		)
	}

	_, err := stmt.ExecContext(ctx, db)
	return err
}

func (s *Store) CountColleagues(
	ctx context.Context,
	db qrm.DB,
	q ListColleaguesQuery,
) (int64, error) {
	tColleague := table.FivenetUser.AS("colleague")
	condition := mysql.AND(
		tUserJobs.Job.EQ(mysql.String(q.Job)),
		mysql.Bool(true),
	)
	if q.Where != nil {
		condition = condition.AND(q.Where)
	}

	if len(q.UserIDs) > 0 && q.UserOnly {
		userIds := make([]mysql.Expression, len(q.UserIDs))
		for i := range q.UserIDs {
			userIds[i] = mysql.Int32(q.UserIDs[i])
		}
		condition = condition.AND(tColleague.ID.IN(userIds...))
	} else if search := dbutils.PrepareForLikeSearch(q.Search); search != "" {
		condition = condition.AND(
			mysql.CONCAT(tColleague.Firstname, mysql.String(" "), tColleague.Lastname).
				LIKE(mysql.String(search)),
		)
	}

	if q.Absent {
		condition = condition.AND(mysql.AND(
			tColleagueProps.AbsenceBegin.IS_NOT_NULL(),
			tColleagueProps.AbsenceEnd.IS_NOT_NULL(),
			tColleagueProps.AbsenceBegin.LT_EQ(mysql.CURRENT_DATE()),
			tColleagueProps.AbsenceEnd.GT_EQ(mysql.CURRENT_DATE()),
		))
	}
	if q.NamePrefix != "" {
		namePrefix := dbutils.PrepareForLikeSearch(q.NamePrefix)
		if namePrefix != "" {
			condition = condition.AND(
				mysql.AND(
					tColleagueProps.NamePrefix.IS_NOT_NULL(),
					tColleagueProps.NamePrefix.LIKE(mysql.String(namePrefix)),
				),
			)
		}
	}
	if q.NameSuffix != "" {
		nameSuffix := dbutils.PrepareForLikeSearch(q.NameSuffix)
		if nameSuffix != "" {
			condition = condition.AND(
				mysql.AND(
					tColleagueProps.NameSuffix.IS_NOT_NULL(),
					tColleagueProps.NameSuffix.LIKE(mysql.String(nameSuffix)),
				),
			)
		}
	}
	if len(q.LabelIDs) > 0 {
		labelIDExprs := make([]mysql.Expression, len(q.LabelIDs))
		for i := range q.LabelIDs {
			labelIDExprs[i] = mysql.Int64(q.LabelIDs[i])
		}

		condition = condition.AND(mysql.EXISTS(
			mysql.
				SELECT(mysql.Int(1)).
				FROM(tColleagueLabels).
				WHERE(mysql.AND(
					tColleagueLabels.UserID.EQ(tColleague.ID),
					tColleagueLabels.Job.EQ(mysql.String(q.Job)),
					tColleagueLabels.LabelID.IN(labelIDExprs...),
				)),
		))
	}

	countStmt := tColleague.
		SELECT(mysql.COUNT(mysql.DISTINCT(tColleague.ID)).AS("data_count.total")).
		FROM(
			tColleague.
				INNER_JOIN(tUserJobs, tUserJobs.UserID.EQ(tColleague.ID)).
				LEFT_JOIN(tColleagueProps, mysql.AND(tColleagueProps.UserID.EQ(tColleague.ID), tColleagueProps.Job.EQ(mysql.String(q.Job)))),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return count.Total, nil
}

func (s *Store) ListColleagues(
	ctx context.Context,
	db qrm.DB,
	q ListColleaguesQuery,
) ([]*jobscolleagues.Colleague, error) {
	tColleague := table.FivenetUser.AS("colleague")
	condition := mysql.AND(
		tUserJobs.Job.EQ(mysql.String(q.Job)),
		mysql.Bool(true),
	)
	if q.Where != nil {
		condition = condition.AND(q.Where)
	}

	if len(q.UserIDs) > 0 && q.UserOnly {
		userIds := make([]mysql.Expression, len(q.UserIDs))
		for i := range q.UserIDs {
			userIds[i] = mysql.Int32(q.UserIDs[i])
		}
		condition = condition.AND(tColleague.ID.IN(userIds...))
	} else if search := dbutils.PrepareForLikeSearch(q.Search); search != "" {
		condition = condition.AND(
			mysql.CONCAT(tColleague.Firstname, mysql.String(" "), tColleague.Lastname).
				LIKE(mysql.String(search)),
		)
	}

	if q.Absent {
		condition = condition.AND(mysql.AND(
			tColleagueProps.AbsenceBegin.IS_NOT_NULL(),
			tColleagueProps.AbsenceEnd.IS_NOT_NULL(),
			tColleagueProps.AbsenceBegin.LT_EQ(mysql.CURRENT_DATE()),
			tColleagueProps.AbsenceEnd.GT_EQ(mysql.CURRENT_DATE()),
		))
	}
	if q.NamePrefix != "" {
		namePrefix := dbutils.PrepareForLikeSearch(q.NamePrefix)
		if namePrefix != "" {
			condition = condition.AND(
				mysql.AND(
					tColleagueProps.NamePrefix.IS_NOT_NULL(),
					tColleagueProps.NamePrefix.LIKE(mysql.String(namePrefix)),
				),
			)
		}
	}
	if q.NameSuffix != "" {
		nameSuffix := dbutils.PrepareForLikeSearch(q.NameSuffix)
		if nameSuffix != "" {
			condition = condition.AND(
				mysql.AND(
					tColleagueProps.NameSuffix.IS_NOT_NULL(),
					tColleagueProps.NameSuffix.LIKE(mysql.String(nameSuffix)),
				),
			)
		}
	}
	if len(q.LabelIDs) > 0 {
		labelIDExprs := make([]mysql.Expression, len(q.LabelIDs))
		for i := range q.LabelIDs {
			labelIDExprs[i] = mysql.Int64(q.LabelIDs[i])
		}

		condition = condition.AND(mysql.EXISTS(
			mysql.
				SELECT(mysql.Int(1)).
				FROM(tColleagueLabels).
				WHERE(mysql.AND(
					tColleagueLabels.UserID.EQ(tColleague.ID),
					tColleagueLabels.Job.EQ(mysql.String(q.Job)),
					tColleagueLabels.LabelID.IN(labelIDExprs...),
				)),
		))
	}

	orderBys := []mysql.OrderByClause{
		tUserJobs.Grade.ASC(),
		tColleague.Firstname.ASC(),
		tColleague.Lastname.ASC(),
	}
	if q.Sort != nil && len(q.Sort.GetColumns()) > 0 {
		orderBys = []mysql.OrderByClause{}
		for _, sc := range q.Sort.GetColumns() {
			var columns []mysql.Column
			switch sc.GetId() {
			case nameColumn:
				columns = append(columns, tColleague.Firstname, tColleague.Lastname)
			case rankColumn:
				fallthrough
			default:
				columns = append(columns, tUserJobs.Grade)
			}
			for _, column := range columns {
				if sc.GetDesc() {
					orderBys = append(orderBys, column.DESC())
				} else {
					orderBys = append(orderBys, column.ASC())
				}
			}
		}
	}

	stmt := tColleague.
		SELECT(
			tColleague.ID,
			tUserJobs.Job.AS("colleague.job"),
			tUserJobs.Grade.AS("colleague.job_grade"),
			tColleague.Firstname,
			tColleague.Lastname,
			tColleague.Dateofbirth,
			tColleague.PhoneNumber,
			tUserProps.AvatarFileID.AS("colleague.profile_picture_file_id"),
			tAvatar.FilePath.AS("colleague.profile_picture"),
			tUserProps.Email.AS("colleague.email"),
			tColleagueProps.UserID,
			tColleagueProps.Job,
			tColleagueProps.AbsenceBegin,
			tColleagueProps.AbsenceEnd,
			tColleagueProps.NamePrefix,
			tColleagueProps.NameSuffix,
		).
		FROM(
			tColleague.
				INNER_JOIN(tUserJobs, tUserJobs.UserID.EQ(tColleague.ID)).
				LEFT_JOIN(tUserProps, tUserProps.UserID.EQ(tColleague.ID)).
				LEFT_JOIN(tColleagueProps, mysql.AND(tColleagueProps.UserID.EQ(tColleague.ID), tColleagueProps.Job.EQ(mysql.String(q.Job)))).
				LEFT_JOIN(tAvatar, tAvatar.ID.EQ(tUserProps.AvatarFileID)),
		).
		WHERE(condition).
		OFFSET(q.Offset).
		ORDER_BY(orderBys...).
		LIMIT(q.Limit)

	colleagues := []*jobscolleagues.Colleague{}
	if err := stmt.QueryContext(ctx, db, &colleagues); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return colleagues, nil
}

func (s *Store) GetColleague(
	ctx context.Context,
	db qrm.DB,
	job string,
	userId int32,
	withColumns mysql.ProjectionList,
	includeDeleted bool,
) (*jobscolleagues.Colleague, error) {
	tColleague := table.FivenetUser.AS("colleague")
	columns := mysql.ProjectionList{
		tColleague.Firstname,
		tColleague.Lastname,
		tUserJobs.Job.AS("colleague.job"),
		tUserJobs.Grade.AS("colleague.job_grade"),
		tColleague.Dateofbirth,
		tColleague.PhoneNumber,
		tUserProps.AvatarFileID.AS("colleague.profile_picture_file_id"),
		tAvatar.FilePath.AS("colleague.profile_picture"),
		tUserProps.Email.AS("colleague.email"),
		tColleagueProps.UserID,
		tColleagueProps.Job,
		tColleagueProps.AbsenceBegin,
		tColleagueProps.AbsenceEnd,
		tColleagueProps.NamePrefix,
		tColleagueProps.NameSuffix,
	}
	columns = append(columns, withColumns...)

	stmt := tColleague.
		SELECT(tColleague.ID, columns...).
		FROM(
			tColleague.
				INNER_JOIN(tUserJobs, mysql.AND(tUserJobs.UserID.EQ(tColleague.ID), tUserJobs.Job.EQ(mysql.String(job)))).
				LEFT_JOIN(tUserProps, tUserProps.UserID.EQ(tColleague.ID)).
				LEFT_JOIN(tColleagueProps, mysql.AND(tColleagueProps.UserID.EQ(tColleague.ID), tColleagueProps.Job.EQ(mysql.String(job)))).
				LEFT_JOIN(tAvatar, tAvatar.ID.EQ(tUserProps.AvatarFileID)),
		).
		WHERE(tColleague.ID.EQ(mysql.Int32(userId))).
		LIMIT(1)

	dest := &jobscolleagues.Colleague{}
	if err := stmt.QueryContext(ctx, db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}
	if dest.GetUserId() == 0 {
		return nil, nil
	}
	if dest.GetProps() == nil {
		dest.Props = &jobscolleagues.ColleagueProps{
			UserId: dest.GetUserId(),
			Job:    job,
		}
	}

	labels, err := s.GetUserLabels(ctx, db, job, userId, includeDeleted)
	if err != nil {
		return nil, err
	}
	dest.Props.Labels = labels

	return dest, nil
}

func (s *Store) CountColleagueActivity(ctx context.Context, db qrm.DB, q ListQuery) (int64, error) {
	tActivity := tColleagueActivity.AS("colleague_activity")
	tTargetColleague := table.FivenetUser.AS("target_user")

	condition := mysql.AND(tActivity.Job.EQ(mysql.String(q.Job)))
	if q.Where != nil {
		condition = condition.AND(q.Where)
	}

	countStmt := tActivity.
		SELECT(mysql.COUNT(mysql.DISTINCT(tActivity.ID)).AS("data_count.total")).
		FROM(
			tActivity.
				INNER_JOIN(tTargetColleague, tTargetColleague.ID.EQ(tActivity.TargetUserID)),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return count.Total, nil
}

func (s *Store) ListColleagueActivity(
	ctx context.Context,
	db qrm.DB,
	q ListQuery,
) ([]*colleaguesactivity.ColleagueActivity, error) {
	tActivity := tColleagueActivity.AS("colleague_activity")

	condition := mysql.AND(tActivity.Job.EQ(mysql.String(q.Job)))
	if q.Where != nil {
		condition = condition.AND(q.Where)
	}

	orderBys := []mysql.OrderByClause{tActivity.CreatedAt.DESC()}
	if q.Sort != nil && len(q.Sort.GetColumns()) > 0 {
		orderBys = []mysql.OrderByClause{}
		for _, sc := range q.Sort.GetColumns() {
			column := tActivity.CreatedAt
			if sc.GetDesc() {
				orderBys = append(orderBys, column.DESC())
			} else {
				orderBys = append(orderBys, column.ASC())
			}
		}
	}

	tTargetColleague := table.FivenetUser.AS("target_user")
	tTargetUserJobs := table.FivenetUserJobs.AS("target_user_jobs")
	tSourceUser := tTargetColleague.AS("source_user")
	tSourceUserJobs := table.FivenetUserJobs.AS("source_user_jobs")
	tTargetUserProps := tUserProps.AS("target_user_props")
	tTargetUserAvatar := tAvatar.AS("target_user_profile_picture")
	tTargetColleagueProps := tColleagueProps.AS("fivenet_colleague_props")
	tSourceUserProps := tUserProps.AS("source_user_props")
	tSourceUserAvatar := tAvatar.AS("source_user_profile_picture")

	stmt := tActivity.
		SELECT(
			tActivity.ID,
			tActivity.CreatedAt,
			tActivity.Job,
			tActivity.SourceUserID,
			tActivity.TargetUserID,
			tActivity.ActivityType,
			tActivity.Reason,
			tActivity.Data,
			tTargetColleague.ID,
			tTargetUserJobs.Job.AS("target_user.job"),
			tTargetUserJobs.Grade.AS("target_user.job_grade"),
			tTargetColleague.Job,
			tTargetColleague.JobGrade,
			tTargetColleague.Firstname,
			tTargetColleague.Lastname,
			tTargetColleague.Dateofbirth,
			tTargetColleague.PhoneNumber,
			tTargetUserProps.AvatarFileID.AS("target_user.profile_picture_file_id"),
			tTargetUserAvatar.FilePath.AS("target_user.profile_picture"),
			tTargetColleagueProps.UserID,
			tTargetColleagueProps.Job,
			tTargetColleagueProps.AbsenceBegin,
			tTargetColleagueProps.AbsenceEnd,
			tTargetColleagueProps.NamePrefix,
			tTargetColleagueProps.NameSuffix,
			tSourceUser.ID,
			tSourceUserJobs.Job.AS("source_user.job"),
			tSourceUserJobs.Grade.AS("source_user.job_grade"),
			tSourceUser.Firstname,
			tSourceUser.Lastname,
			tSourceUser.Dateofbirth,
			tSourceUser.PhoneNumber,
			tSourceUserProps.AvatarFileID.AS("source_user.profile_picture_file_id"),
			tSourceUserAvatar.FilePath.AS("source_user.profile_picture"),
		).
		FROM(
			tActivity.
				INNER_JOIN(tTargetColleague, tTargetColleague.ID.EQ(tActivity.TargetUserID)).
				LEFT_JOIN(tTargetUserJobs, mysql.AND(tTargetUserJobs.UserID.EQ(tTargetColleague.ID), tTargetUserJobs.Job.EQ(mysql.String(q.Job)))).
				LEFT_JOIN(tTargetUserProps, tTargetUserProps.UserID.EQ(tTargetColleague.ID)).
				LEFT_JOIN(tTargetUserAvatar, tTargetUserAvatar.ID.EQ(tTargetUserProps.AvatarFileID)).
				LEFT_JOIN(tTargetColleagueProps, mysql.AND(tTargetColleagueProps.UserID.EQ(tTargetColleague.ID), tTargetColleague.Job.EQ(mysql.String(q.Job)))).
				LEFT_JOIN(tSourceUser, tSourceUser.ID.EQ(tActivity.SourceUserID)).
				LEFT_JOIN(tSourceUserJobs, mysql.AND(tSourceUserJobs.UserID.EQ(tSourceUser.ID), tSourceUserJobs.Job.EQ(mysql.String(q.Job)))).
				LEFT_JOIN(tSourceUserProps, tSourceUserProps.UserID.EQ(tSourceUser.ID)).
				LEFT_JOIN(tSourceUserAvatar, tSourceUserAvatar.ID.EQ(tSourceUserProps.AvatarFileID)),
		).
		WHERE(condition).
		OFFSET(q.Offset).
		ORDER_BY(orderBys...).
		LIMIT(q.Limit)

	activity := []*colleaguesactivity.ColleagueActivity{}
	if err := stmt.QueryContext(ctx, db, &activity); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return activity, nil
}

func (s *Store) GetColleagueProps(
	ctx context.Context,
	db qrm.DB,
	job string,
	userId int32,
	fields []string,
	includeDeleted bool,
) (*jobscolleagues.ColleagueProps, error) {
	columns := mysql.ProjectionList{
		tColleagueProps.Job,
		tColleagueProps.AbsenceBegin,
		tColleagueProps.AbsenceEnd,
		tColleagueProps.NamePrefix,
		tColleagueProps.NameSuffix,
	}

	if fields == nil {
		fields = append(fields, "Note")
	}

	for _, field := range fields {
		switch field {
		case "Note":
			columns = append(columns, tColleagueProps.Note)
		}
	}

	stmt := tColleagueProps.
		SELECT(
			tColleagueProps.UserID,
			columns...,
		).
		FROM(tColleagueProps).
		WHERE(mysql.AND(
			tColleagueProps.UserID.EQ(mysql.Int32(userId)),
			tColleagueProps.Job.EQ(mysql.String(job)),
		)).
		LIMIT(1)

	dest := &jobscolleagues.ColleagueProps{UserId: userId}
	if err := stmt.QueryContext(ctx, db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	labels, err := s.GetUserLabels(ctx, db, job, userId, includeDeleted)
	if err != nil {
		return nil, err
	}
	dest.Labels = labels

	return dest, nil
}

func (s *Store) HandleColleaguePropsChanges(
	ctx context.Context,
	db qrm.DB,
	x *jobscolleagues.ColleagueProps,
	in *jobscolleagues.ColleagueProps,
	job string,
	sourceUserId *int32,
	reason string,
) ([]*colleaguesactivity.ColleagueActivity, error) {
	tColleagueProps := table.FivenetJobColleagueProps

	updateSets := []mysql.ColumnAssigment{}

	absenceBegin := mysql.DateExp(mysql.NULL)
	absenceEnd := mysql.DateExp(mysql.NULL)
	if in.GetAbsenceBegin() != nil && in.GetAbsenceEnd() != nil {
		if in.GetAbsenceBegin().GetTimestamp() == nil {
			in.AbsenceBegin = nil
		} else {
			absenceBegin = mysql.DateT(in.GetAbsenceBegin().AsTime())
		}

		if in.GetAbsenceEnd().GetTimestamp() == nil {
			in.AbsenceEnd = nil
		} else {
			absenceEnd = mysql.DateT(in.GetAbsenceEnd().AsTime())
		}
	} else {
		in.AbsenceBegin = x.GetAbsenceBegin()
		in.AbsenceEnd = x.GetAbsenceEnd()
	}

	updateSets = append(updateSets,
		tColleagueProps.AbsenceBegin.SET(mysql.DateExp(mysql.Raw("VALUES(`absence_begin`)"))),
		tColleagueProps.AbsenceEnd.SET(mysql.DateExp(mysql.Raw("VALUES(`absence_end`)"))),
	)

	if in.Note != nil {
		if in.GetNote() == "" {
			updateSets = append(updateSets, tColleagueProps.Note.SET(mysql.StringExp(mysql.NULL)))
		} else {
			updateSets = append(updateSets, tColleagueProps.Note.SET(mysql.String(in.GetNote())))
		}
	} else {
		in.Note = x.Note
	}

	if in.GetLabels() == nil {
		in.Labels = x.GetLabels()
	}

	if in.NamePrefix != nil || in.NameSuffix != nil {
		if in.NamePrefix != nil {
			*in.NamePrefix = strings.TrimSpace(in.GetNamePrefix())
			updateSets = append(
				updateSets,
				tColleagueProps.NamePrefix.SET(mysql.String(in.GetNamePrefix())),
			)
		} else {
			in.NamePrefix = x.NamePrefix
		}
		if in.NameSuffix != nil {
			*in.NameSuffix = strings.TrimSpace(in.GetNameSuffix())
			updateSets = append(
				updateSets,
				tColleagueProps.NameSuffix.SET(mysql.String(in.GetNameSuffix())),
			)
		} else {
			in.NameSuffix = x.NameSuffix
		}
	} else {
		in.NamePrefix = x.NamePrefix
		in.NameSuffix = x.NameSuffix
	}

	stmt := tColleagueProps.
		INSERT(
			tColleagueProps.UserID,
			tColleagueProps.Job,
			tColleagueProps.AbsenceBegin,
			tColleagueProps.AbsenceEnd,
			tColleagueProps.Note,
			tColleagueProps.NamePrefix,
			tColleagueProps.NameSuffix,
		).
		VALUES(
			x.GetUserId(),
			job,
			absenceBegin,
			absenceEnd,
			in.Note,
			in.NamePrefix,
			in.NameSuffix,
		).
		ON_DUPLICATE_KEY_UPDATE(updateSets...)

	if _, err := stmt.ExecContext(ctx, db); err != nil {
		return nil, err
	}

	activities := []*colleaguesactivity.ColleagueActivity{}

	if in.GetLabels() != nil && !proto.Equal(in.GetLabels(), x.GetLabels()) {
		added, removed := utils.SlicesDifferenceFunc(
			x.GetLabels().GetList(),
			in.GetLabels().GetList(),
			func(in *jobslabels.Label) int64 { return in.GetId() },
		)

		if err := updateColleaguesLabels(ctx, db, x.GetUserId(), job, added, removed); err != nil {
			return nil, err
		}

		activities = append(activities, &colleaguesactivity.ColleagueActivity{
			Job:          job,
			SourceUserId: sourceUserId,
			TargetUserId: x.GetUserId(),
			ActivityType: colleaguesactivity.ColleagueActivityType_COLLEAGUE_ACTIVITY_TYPE_LABELS,
			Reason:       reason,
			Data: &colleaguesactivity.ColleagueActivityData{
				Data: &colleaguesactivity.ColleagueActivityData_LabelsChange{
					LabelsChange: &colleaguesactivity.LabelsChange{Added: added, Removed: removed},
				},
			},
		})
	}

	if (in.GetAbsenceBegin() == nil && in.GetAbsenceEnd() == nil && x.GetAbsenceBegin() != nil && x.GetAbsenceEnd() != nil) ||
		(in.GetAbsenceBegin() != nil && (x.GetAbsenceBegin() == nil || in.GetAbsenceBegin().AsTime().Compare(x.GetAbsenceBegin().AsTime()) != 0) ||
			in.GetAbsenceEnd() != nil && (x.GetAbsenceEnd() == nil || in.GetAbsenceEnd().AsTime().Compare(x.GetAbsenceEnd().AsTime()) != 0)) {
		activities = append(activities, &colleaguesactivity.ColleagueActivity{
			Job:          job,
			SourceUserId: sourceUserId,
			TargetUserId: x.GetUserId(),
			ActivityType: colleaguesactivity.ColleagueActivityType_COLLEAGUE_ACTIVITY_TYPE_ABSENCE_DATE,
			Reason:       reason,
			Data: &colleaguesactivity.ColleagueActivityData{
				Data: &colleaguesactivity.ColleagueActivityData_AbsenceDate{
					AbsenceDate: &colleaguesactivity.AbsenceDateChange{
						AbsenceBegin: in.GetAbsenceBegin(),
						AbsenceEnd:   in.GetAbsenceEnd(),
					},
				},
			},
		})
	}

	if in.Note != nil && (x.Note == nil || in.GetNote() != x.GetNote()) {
		activities = append(activities, &colleaguesactivity.ColleagueActivity{
			Job:          job,
			SourceUserId: sourceUserId,
			TargetUserId: x.GetUserId(),
			ActivityType: colleaguesactivity.ColleagueActivityType_COLLEAGUE_ACTIVITY_TYPE_NOTE,
			Reason:       reason,
		})
	}

	if in.NamePrefix != nil && (x.NamePrefix == nil || in.GetNamePrefix() != x.GetNamePrefix()) ||
		in.NameSuffix != nil && (x.NameSuffix == nil || in.GetNameSuffix() != x.GetNameSuffix()) {
		activities = append(activities, &colleaguesactivity.ColleagueActivity{
			Job:          job,
			SourceUserId: sourceUserId,
			TargetUserId: x.GetUserId(),
			ActivityType: colleaguesactivity.ColleagueActivityType_COLLEAGUE_ACTIVITY_TYPE_NAME,
			Reason:       reason,
			Data: &colleaguesactivity.ColleagueActivityData{
				Data: &colleaguesactivity.ColleagueActivityData_NameChange{
					NameChange: &colleaguesactivity.NameChange{
						Prefix: in.NamePrefix,
						Suffix: in.NameSuffix,
					},
				},
			},
		})
	}

	return activities, nil
}

func updateColleaguesLabels(
	ctx context.Context,
	db qrm.DB,
	userId int32,
	job string,
	added []*jobslabels.Label,
	removed []*jobslabels.Label,
) error {
	tColleagueLabels := table.FivenetJobColleagueLabels
	if len(added) > 0 {
		addedLabels := make([]*model.FivenetJobColleagueLabels, len(added))
		for i, label := range added {
			addedLabels[i] = &model.FivenetJobColleagueLabels{
				UserID:  userId,
				Job:     job,
				LabelID: label.GetId(),
			}
		}

		stmt := tColleagueLabels.
			INSERT(
				tColleagueLabels.UserID,
				tColleagueLabels.Job,
				tColleagueLabels.LabelID,
			).
			MODELS(addedLabels)

		if _, err := stmt.ExecContext(ctx, db); err != nil && !dbutils.IsDuplicateError(err) {
			return err
		}
	}

	if len(removed) > 0 {
		ids := make([]mysql.Expression, len(removed))
		for i := range removed {
			ids[i] = mysql.Int64(removed[i].GetId())
		}

		stmt := tColleagueLabels.
			DELETE().
			WHERE(mysql.AND(
				tColleagueLabels.UserID.EQ(mysql.Int32(userId)),
				tColleagueLabels.LabelID.IN(ids...),
			)).
			LIMIT(int64(len(removed)))

		if _, err := stmt.ExecContext(ctx, db); err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) GetUserLabels(
	ctx context.Context,
	db qrm.DB,
	job string,
	userId int32,
	includeDeleted bool,
) (*jobslabels.Labels, error) {
	stmt := tColleagueLabels.
		SELECT(
			tJobLabels.ID,
			tJobLabels.Job,
			tJobLabels.Name,
			tJobLabels.Color,
			tJobLabels.Icon,
			tJobLabels.SortOrder,
		).
		FROM(
			tColleagueLabels.
				INNER_JOIN(tJobLabels,
					tJobLabels.ID.EQ(tColleagueLabels.LabelID),
				),
		).
		WHERE(mysql.AND(
			tColleagueLabels.UserID.EQ(mysql.Int32(userId)),
			tJobLabels.Job.EQ(mysql.String(job)),
			mysql.OR(
				mysql.Bool(includeDeleted),
				tJobLabels.DeletedAt.IS_NULL(),
			),
		)).
		ORDER_BY(
			tJobLabels.SortOrder.ASC(),
			tJobLabels.SortKey.ASC(),
		)

	list := &jobslabels.Labels{List: []*jobslabels.Label{}}
	if err := stmt.QueryContext(ctx, db, &list.List); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return list, nil
}

type UserLabels struct {
	UserId int32 `alias:"userId" sql:"primary_key"`
	Labels *jobslabels.Labels
}
