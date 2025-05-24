package jobs

import (
	"context"
	"errors"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/protobuf/proto"
)

func GetColleagueProps(ctx context.Context, tx qrm.DB, job string, userId int32, fields []string) (*ColleagueProps, error) {
	tColleagueProps := table.FivenetJobColleagueProps.AS("colleague_props")

	columns := []jet.Projection{
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
		WHERE(jet.AND(
			tColleagueProps.UserID.EQ(jet.Int32(userId)),
			tColleagueProps.Job.EQ(jet.String(job)),
		)).
		LIMIT(1)

	dest := &ColleagueProps{
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

func GetUserLabels(ctx context.Context, tx qrm.DB, job string, userId int32) (*Labels, error) {
	tJobLabels := table.FivenetJobLabels.AS("label")
	tUserLabels := table.FivenetJobColleagueLabels

	stmt := tUserLabels.
		SELECT(
			tJobLabels.ID,
			tJobLabels.Job,
			tJobLabels.Name,
			tJobLabels.Color,
		).
		FROM(
			tUserLabels.
				INNER_JOIN(tJobLabels,
					tJobLabels.ID.EQ(tUserLabels.LabelID),
				),
		).
		WHERE(jet.AND(
			tUserLabels.UserID.EQ(jet.Int32(userId)),
			tJobLabels.Job.EQ(jet.String(job)),
			tJobLabels.DeletedAt.IS_NULL(),
		)).
		ORDER_BY(
			tJobLabels.Order.ASC(),
		)

	list := &Labels{
		List: []*Label{},
	}
	if err := stmt.QueryContext(ctx, tx, &list.List); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return list, nil
}

func (x *ColleagueProps) HandleChanges(ctx context.Context, tx qrm.DB, in *ColleagueProps, job string, sourceUserId *int32, reason string) ([]*ColleagueActivity, error) {
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

	tColleagueProps := table.FivenetJobColleagueProps

	updateSets := []jet.ColumnAssigment{
		tColleagueProps.AbsenceBegin.SET(jet.DateExp(jet.Raw("VALUES(`absence_begin`)"))),
		tColleagueProps.AbsenceEnd.SET(jet.DateExp(jet.Raw("VALUES(`absence_end`)"))),
	}

	// Generate the update sets
	if in.Note != nil {
		// Set empty note to null
		if *in.Note == "" {
			updateSets = append(updateSets, tColleagueProps.Note.SET(jet.StringExp(jet.NULL)))
		} else {
			updateSets = append(updateSets, tColleagueProps.Note.SET(jet.String(*in.Note)))
		}
	} else {
		in.Note = x.Note
	}

	if in.Labels == nil {
		in.Labels = x.Labels
	}

	if in.NamePrefix != nil || in.NameSuffix != nil {
		if in.NamePrefix != nil {
			*in.NamePrefix = strings.TrimSpace(*in.NamePrefix) // Trim spaces
			updateSets = append(updateSets, tColleagueProps.NamePrefix.SET(jet.String(*in.NamePrefix)))
		} else {
			in.NamePrefix = x.NamePrefix
		}
		if in.NameSuffix != nil {
			*in.NameSuffix = strings.TrimSpace(*in.NameSuffix) // Trim spaces
			updateSets = append(updateSets, tColleagueProps.NameSuffix.SET(jet.String(*in.NameSuffix)))
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

	activities := []*ColleagueActivity{}

	// Create user activity entries
	if in.Labels != nil && !proto.Equal(in.Labels, x.Labels) {
		added, removed := utils.SlicesDifferenceFunc(x.Labels.List, in.Labels.List,
			func(in *Label) uint64 {
				return in.Id
			})

		if err := x.updateLabels(ctx, tx, x.UserId, job, added, removed); err != nil {
			return nil, err
		}

		activities = append(activities, &ColleagueActivity{
			Job:          job,
			SourceUserId: sourceUserId,
			TargetUserId: x.UserId,
			ActivityType: ColleagueActivityType_COLLEAGUE_ACTIVITY_TYPE_LABELS,
			Reason:       reason,
			Data: &ColleagueActivityData{
				Data: &ColleagueActivityData_LabelsChange{
					LabelsChange: &LabelsChange{
						Added:   added,
						Removed: removed,
					},
				},
			},
		})
	}

	// Compare absence dates if any were set
	if (in.AbsenceBegin == nil && in.AbsenceEnd == nil && x.AbsenceBegin != nil && x.AbsenceEnd != nil) ||
		(in.AbsenceBegin != nil && (x.AbsenceBegin == nil || in.AbsenceBegin.AsTime().Compare(x.AbsenceBegin.AsTime()) != 0) ||
			in.AbsenceEnd != nil && (x.AbsenceEnd == nil || in.AbsenceEnd.AsTime().Compare(x.AbsenceEnd.AsTime()) != 0)) {
		activities = append(activities, &ColleagueActivity{
			Job:          job,
			SourceUserId: sourceUserId,
			TargetUserId: x.UserId,
			ActivityType: ColleagueActivityType_COLLEAGUE_ACTIVITY_TYPE_ABSENCE_DATE,
			Reason:       reason,
			Data: &ColleagueActivityData{
				Data: &ColleagueActivityData_AbsenceDate{
					AbsenceDate: &AbsenceDateChange{
						AbsenceBegin: in.AbsenceBegin,
						AbsenceEnd:   in.AbsenceEnd,
					},
				},
			},
		})
	}

	if in.Note != nil && (x.Note == nil || *in.Note != *x.Note) {
		activities = append(activities, &ColleagueActivity{
			Job:          job,
			SourceUserId: sourceUserId,
			TargetUserId: x.UserId,
			ActivityType: ColleagueActivityType_COLLEAGUE_ACTIVITY_TYPE_NOTE,
			Reason:       reason,
		})
	}

	if in.NamePrefix != nil && (x.NamePrefix == nil || *in.NamePrefix != *x.NamePrefix) ||
		in.NameSuffix != nil && (x.NameSuffix == nil || *in.NameSuffix != *x.NameSuffix) {
		activities = append(activities, &ColleagueActivity{
			Job:          job,
			SourceUserId: sourceUserId,
			TargetUserId: x.UserId,
			ActivityType: ColleagueActivityType_COLLEAGUE_ACTIVITY_TYPE_NAME,
			Reason:       reason,
			Data: &ColleagueActivityData{
				Data: &ColleagueActivityData_NameChange{
					NameChange: &NameChange{
						Prefix: in.NamePrefix,
						Suffix: in.NameSuffix,
					},
				},
			},
		})
	}

	return activities, nil
}

func (x *ColleagueProps) updateLabels(ctx context.Context, tx qrm.DB, userId int32, job string, added []*Label, removed []*Label) error {
	tUserLabels := table.FivenetJobColleagueLabels

	if len(added) > 0 {
		addedLabels := make([]*model.FivenetJobColleagueLabels, len(added))
		for i, label := range added {
			addedLabels[i] = &model.FivenetJobColleagueLabels{
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
