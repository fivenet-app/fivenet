package jobs

import (
	"context"

	"github.com/fivenet-app/fivenet/pkg/utils"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/protobuf/proto"
)

func (x *JobsUserProps) HandleChanges(ctx context.Context, tx qrm.DB, in *JobsUserProps, job string, sourceUserId *int32, reason string) ([]*JobsUserActivity, error) {
	absenceBegin := jet.DateExp(jet.NULL)
	absenceEnd := jet.DateExp(jet.NULL)
	if in.AbsenceBegin != nil && in.AbsenceEnd != nil {
		if in.AbsenceBegin.Timestamp == nil {
			in.AbsenceBegin = nil
		} else {
			absenceBegin = jet.DateT(in.AbsenceBegin.AsTime())
		}

		if in.AbsenceEnd.Timestamp == nil {
			in.AbsenceEnd = nil
		} else {
			absenceEnd = jet.DateT(in.AbsenceEnd.AsTime())
		}
	} else {
		in.AbsenceBegin = x.AbsenceBegin
		in.AbsenceEnd = x.AbsenceEnd
	}

	tJobsUserProps := table.FivenetJobsUserProps

	updateSets := []jet.ColumnAssigment{
		tJobsUserProps.AbsenceBegin.SET(jet.DateExp(jet.Raw("VALUES(`absence_begin`)"))),
		tJobsUserProps.AbsenceEnd.SET(jet.DateExp(jet.Raw("VALUES(`absence_end`)"))),
	}

	// Generate the update sets
	if in.Note != nil {
		updateSets = append(updateSets, tJobsUserProps.Note.SET(jet.String(*in.Note)))
	} else {
		in.Note = x.Note
	}

	if in.Labels == nil {
		in.Labels = x.Labels
	}

	if in.NamePrefix != nil || in.NameSuffix != nil {
		if in.NamePrefix != nil {
			updateSets = append(updateSets, tJobsUserProps.NamePrefix.SET(jet.String(*in.NamePrefix)))
		} else {
			in.NamePrefix = x.NamePrefix
		}
		if in.NameSuffix != nil {
			updateSets = append(updateSets, tJobsUserProps.NameSuffix.SET(jet.String(*in.NameSuffix)))
		} else {
			in.NameSuffix = x.NameSuffix
		}
	} else {
		in.NamePrefix = x.NamePrefix
		in.NameSuffix = x.NameSuffix
	}

	stmt := tJobsUserProps.
		INSERT(
			tJobsUserProps.UserID,
			tJobsUserProps.Job,
			tJobsUserProps.AbsenceBegin,
			tJobsUserProps.AbsenceEnd,
			tJobsUserProps.Note,
			tJobsUserProps.NamePrefix,
			tJobsUserProps.NameSuffix,
		).
		VALUES(
			x.UserId,
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

	activities := []*JobsUserActivity{}

	// Create user activity entries
	if in.Labels != nil && !proto.Equal(in.Labels, x.Labels) {
		added, removed := utils.SlicesDifferenceFunc(x.Labels.List, in.Labels.List,
			func(in *Label) uint64 {
				return in.Id
			})

		if err := x.updateLabels(ctx, tx, x.UserId, job, added, removed); err != nil {
			return nil, err
		}

		activities = append(activities, &JobsUserActivity{
			Job:          job,
			SourceUserId: sourceUserId,
			TargetUserId: x.UserId,
			ActivityType: JobsUserActivityType_JOBS_USER_ACTIVITY_TYPE_LABELS,
			Reason:       reason,
			Data: &JobsUserActivityData{
				Data: &JobsUserActivityData_LabelsChange{
					LabelsChange: &ColleagueLabelsChange{
						Added:   added,
						Removed: removed,
					},
				},
			},
		})
	}

	// Compare absence dates if any were set
	if in.AbsenceBegin != nil && (x.AbsenceBegin == nil || in.AbsenceBegin.AsTime().Compare(x.AbsenceBegin.AsTime()) != 0) ||
		in.AbsenceEnd != nil && (x.AbsenceEnd == nil || in.AbsenceEnd.AsTime().Compare(x.AbsenceEnd.AsTime()) != 0) {
		activities = append(activities, &JobsUserActivity{
			Job:          job,
			SourceUserId: sourceUserId,
			TargetUserId: x.UserId,
			ActivityType: JobsUserActivityType_JOBS_USER_ACTIVITY_TYPE_ABSENCE_DATE,
			Reason:       reason,
			Data: &JobsUserActivityData{
				Data: &JobsUserActivityData_AbsenceDate{
					AbsenceDate: &ColleagueAbsenceDate{
						AbsenceBegin: in.AbsenceBegin,
						AbsenceEnd:   in.AbsenceEnd,
					},
				},
			},
		})
	}

	if (in.Note == nil && x.Note != nil) || (in.Note != nil && x.Note == nil) {
		activities = append(activities, &JobsUserActivity{
			Job:          job,
			SourceUserId: sourceUserId,
			TargetUserId: x.UserId,
			ActivityType: JobsUserActivityType_JOBS_USER_ACTIVITY_TYPE_NOTE,
			Reason:       reason,
		})
	}

	if in.NamePrefix != nil && (x.NamePrefix == nil || in.NamePrefix != x.NamePrefix) ||
		in.NameSuffix != nil && (x.NameSuffix == nil || in.NameSuffix != x.NameSuffix) {
		activities = append(activities, &JobsUserActivity{
			Job:          job,
			SourceUserId: sourceUserId,
			TargetUserId: x.UserId,
			ActivityType: JobsUserActivityType_JOBS_USER_ACTIVITY_TYPE_NAME,
			Reason:       reason,
			Data: &JobsUserActivityData{
				Data: &JobsUserActivityData_NameChange{
					NameChange: &ColleagueNameChange{
						Prefix: in.NamePrefix,
						Suffix: in.NameSuffix,
					},
				},
			},
		})
	}

	return activities, nil
}

func (x *JobsUserProps) updateLabels(ctx context.Context, tx qrm.DB, userId int32, job string, added []*Label, removed []*Label) error {
	tUserLabels := table.FivenetJobsLabelsUsers

	if len(added) > 0 {
		addedLabels := make([]*model.FivenetJobsLabelsUsers, len(added))
		for i, label := range added {
			addedLabels[i] = &model.FivenetJobsLabelsUsers{
				UserID:  userId,
				Job:     job,
				LabelID: label.Id,
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
		ids := make([]jet.Expression, len(removed))

		for i := range removed {
			ids[i] = jet.Uint64(removed[i].Id)
		}

		stmt := tUserLabels.
			DELETE().
			WHERE(jet.AND(
				tUserLabels.UserID.EQ(jet.Int32(userId)),
				tUserLabels.LabelID.IN(ids...),
			)).
			LIMIT(int64(len(removed)))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}
