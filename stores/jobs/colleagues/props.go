package colleaguesstore

import (
	"context"
	"errors"
	"strings"

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

func GetColleagueProps(
	ctx context.Context,
	tx qrm.DB,
	job string,
	userId int32,
	fields []string,
) (*jobscolleagues.ColleagueProps, error) {
	tColleagueProps := table.FivenetJobColleagueProps.AS("colleague_props")

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

	dest := &jobscolleagues.ColleagueProps{
		UserId: userId,
	}
	if err := stmt.QueryContext(ctx, tx, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	labels, err := GetUserLabels(ctx, tx, job, userId)
	if err != nil {
		return nil, err
	}
	dest.Labels = labels

	return dest, nil
}

func HandleColleaguesPropsChanges(
	ctx context.Context,
	tx qrm.DB,
	x *jobscolleagues.ColleagueProps,
	in *jobscolleagues.ColleagueProps,
	job string,
	sourceUserId *int32,
	reason string,
) ([]*colleaguesactivity.ColleagueActivity, error) {
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

	tColleagueProps := table.FivenetJobColleagueProps

	updateSets := []mysql.ColumnAssigment{
		tColleagueProps.AbsenceBegin.SET(mysql.DateExp(mysql.Raw("VALUES(`absence_begin`)"))),
		tColleagueProps.AbsenceEnd.SET(mysql.DateExp(mysql.Raw("VALUES(`absence_end`)"))),
	}

	// Generate the update sets
	if in.Note != nil {
		// Set empty note to null
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
			*in.NamePrefix = strings.TrimSpace(in.GetNamePrefix()) // Trim spaces
			updateSets = append(
				updateSets,
				tColleagueProps.NamePrefix.SET(mysql.String(in.GetNamePrefix())),
			)
		} else {
			in.NamePrefix = x.NamePrefix
		}
		if in.NameSuffix != nil {
			*in.NameSuffix = strings.TrimSpace(in.GetNameSuffix()) // Trim spaces
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
		ON_DUPLICATE_KEY_UPDATE(
			updateSets...,
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, err
	}

	activities := []*colleaguesactivity.ColleagueActivity{}

	// Create user activity entries
	if in.GetLabels() != nil && !proto.Equal(in.GetLabels(), x.GetLabels()) {
		added, removed := utils.SlicesDifferenceFunc(
			x.GetLabels().GetList(),
			in.GetLabels().GetList(),
			func(in *jobslabels.Label) int64 {
				return in.GetId()
			},
		)

		if err := updateColleaguesLabels(ctx, tx, x.GetUserId(), job, added, removed); err != nil {
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
					LabelsChange: &colleaguesactivity.LabelsChange{
						Added:   added,
						Removed: removed,
					},
				},
			},
		})
	}

	// Compare absence dates if any were set
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
	tx qrm.DB,
	userId int32,
	job string,
	added []*jobslabels.Label,
	removed []*jobslabels.Label,
) error {
	tUserLabels := table.FivenetJobColleagueLabels

	if len(added) > 0 {
		addedLabels := make([]*model.FivenetJobColleagueLabels, len(added))
		for i, label := range added {
			addedLabels[i] = &model.FivenetJobColleagueLabels{
				UserID:  userId,
				Job:     job,
				LabelID: label.GetId(),
			}
		}

		stmt := tUserLabels.
			INSERT(
				tUserLabels.UserID,
				tUserLabels.Job,
				tUserLabels.LabelID,
			).
			MODELS(addedLabels)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return err
			}
		}
	}

	if len(removed) > 0 {
		ids := make([]mysql.Expression, len(removed))

		for i := range removed {
			ids[i] = mysql.Int64(removed[i].GetId())
		}

		stmt := tUserLabels.
			DELETE().
			WHERE(mysql.AND(
				tUserLabels.UserID.EQ(mysql.Int32(userId)),
				tUserLabels.LabelID.IN(ids...),
			)).
			LIMIT(int64(len(removed)))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}
