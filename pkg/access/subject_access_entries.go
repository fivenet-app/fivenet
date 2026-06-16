package access

import (
	"context"
	"errors"
	"fmt"
	"slices"

	resourcesaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type SubjectAccessOptions struct {
	BlockedAccess      int32
	DeniedAccessLevels []int32
}

type SubjectAccessEntryChanges[T any] struct {
	ToCreate []T
	ToUpdate []T
	ToDelete []T
}

func (c *SubjectAccessEntryChanges[T]) IsEmpty() bool {
	return len(c.ToCreate) == 0 && len(c.ToUpdate) == 0 && len(c.ToDelete) == 0
}

type SubjectAccessChanges struct {
	Jobs           *SubjectAccessEntryChanges[*resourcesaccess.JobAccess]
	Users          *SubjectAccessEntryChanges[*resourcesaccess.UserAccess]
	Qualifications *SubjectAccessEntryChanges[*resourcesaccess.QualificationAccess]
}

func (c *SubjectAccessChanges) IsEmpty() bool {
	return c.Jobs.IsEmpty() && c.Users.IsEmpty() && c.Qualifications.IsEmpty()
}

func (a *SubjectObjectAccess) ListTargetAccess(
	ctx context.Context,
	tx qrm.DB,
	targetID int64,
	opts SubjectAccessOptions,
) (*resourcesaccess.Access, error) {
	tAccess := a.accessTable
	tSubjects := table.FivenetACLSubjects.AS("acl_subject")
	tJobGrade := table.FivenetACLSubjectJobGradeScopes.AS("subject_job_grade")
	tQualifications := table.FivenetACLSubjectQualifications.AS("subject_qualification")
	tUserSubjects := table.FivenetACLSubjectUsers.AS("subject_user")
	tUsers := table.FivenetUser.AS("user_short")

	columns := a.accessColumns

	stmt := tAccess.
		SELECT(
			columns.ID.AS("subject_access_row.id"),
			columns.TargetID.AS("subject_access_row.target_id"),
			columns.Access.AS("subject_access_row.access"),
			columns.Effect.AS("subject_access_row.effect"),
			tSubjects.SubjectType.AS("subject_access_row.subject_type"),
			tJobGrade.Job.AS("subject_access_row.acl_job"),
			tJobGrade.MinimumGrade.AS("subject_access_row.acl_minimum_grade"),
			tQualifications.QualificationID.AS("subject_access_row.acl_qualification_id"),
			tUserSubjects.UserID.AS("subject_access_row.subject_user_id"),
			tUsers.Job.AS("subject_access_row.user_job"),
			tUsers.JobGrade.AS("subject_access_row.user_job_grade"),
			tUsers.Firstname.AS("subject_access_row.user_firstname"),
			tUsers.Lastname.AS("subject_access_row.user_lastname"),
			tUsers.Dateofbirth.AS("subject_access_row.user_dateofbirth"),
			tUsers.PhoneNumber.AS("subject_access_row.user_phone_number"),
		).
		FROM(tAccess.
			INNER_JOIN(tSubjects,
				tSubjects.ID.EQ(columns.SubjectID),
			).
			LEFT_JOIN(tJobGrade,
				tJobGrade.SubjectID.EQ(columns.SubjectID),
			).
			LEFT_JOIN(tQualifications,
				tQualifications.SubjectID.EQ(columns.SubjectID),
			).
			LEFT_JOIN(tUserSubjects,
				tUserSubjects.SubjectID.EQ(columns.SubjectID),
			).
			LEFT_JOIN(tUsers,
				tUsers.ID.EQ(tUserSubjects.UserID),
			),
		).
		WHERE(columns.TargetID.EQ(mysql.Int64(targetID))).
		ORDER_BY(columns.ID.ASC())

	var rows []subjectAccessRow
	if err := stmt.QueryContext(ctx, tx, &rows); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return subjectAccessRowsToProto(rows, opts), nil
}

func (a *SubjectObjectAccess) ReplaceTargetAccess(
	ctx context.Context,
	tx qrm.DB,
	resolver *SubjectResolver,
	targetID int64,
	in *resourcesaccess.Access,
	opts SubjectAccessOptions,
) (*SubjectAccessChanges, error) {
	current, err := a.ListTargetAccess(ctx, tx, targetID, opts)
	if err != nil {
		return nil, err
	}

	changes := compareSubjectAccess(current, in)

	if err := a.ClearTarget(ctx, tx, targetID); err != nil {
		return nil, err
	}

	if err := a.createTargetAccess(ctx, tx, resolver, targetID, in, opts); err != nil {
		return nil, err
	}

	if err := a.RefreshTargetVisibility(ctx, tx, targetID); err != nil {
		return nil, err
	}

	return changes, nil
}

func (a *SubjectObjectAccess) CreateUserAccess(
	ctx context.Context,
	tx qrm.DB,
	resolver *SubjectResolver,
	targetID int64,
	userID int32,
	access int32,
) error {
	subjectID, err := resolver.EnsureUserSubject(ctx, tx, userID)
	if err != nil {
		return err
	}

	if err := a.CreateEntry(ctx, tx, targetID, subjectID, access, AccessEffectAllow); err != nil {
		return err
	}

	return a.RefreshTargetVisibility(ctx, tx, targetID)
}

func (a *SubjectObjectAccess) ACLAccessExistsCondition(
	targetID mysql.IntegerExpression,
	userInfo interface {
		GetUserId() int32
		GetJob() string
		GetJobGrade() int32
	},
	access int32,
) mysql.BoolExpression {
	tAccess := a.accessTable
	columns := a.accessColumns
	tSubjectUsers := table.FivenetACLSubjectUsers.AS("subject_acl_user_exists")
	tSubjectQualis := table.FivenetACLSubjectQualifications.AS("subject_acl_qualification_exists")
	tQualiResults := table.FivenetQualificationsResults.AS(
		"subject_acl_qualification_results_exists",
	)
	tSubjectJobGrade := table.FivenetACLSubjectJobGradeScopes.AS("subject_acl_job_grade_exists")
	tUserJobs := table.FivenetUserJobs.AS("subject_acl_user_jobs_exists")

	return mysql.EXISTS(
		mysql.SELECT(mysql.Int(1)).
			FROM(tAccess).
			WHERE(mysql.AND(
				columns.TargetID.EQ(targetID),
				columns.Effect.IS_TRUE(),
				columns.Access.GT_EQ(mysql.Int32(access)),
				mysql.OR(
					columns.SubjectID.IN(
						tSubjectUsers.
							SELECT(tSubjectUsers.SubjectID).
							FROM(tSubjectUsers).
							WHERE(tSubjectUsers.UserID.EQ(mysql.Int32(userInfo.GetUserId()))),
					),
					columns.SubjectID.IN(
						tSubjectQualis.
							SELECT(tSubjectQualis.SubjectID).
							FROM(tSubjectQualis.
								INNER_JOIN(tQualiResults,
									mysql.AND(
										tQualiResults.QualificationID.EQ(
											tSubjectQualis.QualificationID,
										),
										tQualiResults.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
										tQualiResults.DeletedAt.IS_NULL(),
										tQualiResults.Status.EQ(
											mysql.Int32(
												int32(
													qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL,
												),
											),
										),
									),
								),
							),
					),
					columns.SubjectID.IN(
						tSubjectJobGrade.
							SELECT(tSubjectJobGrade.SubjectID).
							FROM(tSubjectJobGrade.
								INNER_JOIN(tUserJobs,
									mysql.AND(
										tUserJobs.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
										tUserJobs.Job.EQ(tSubjectJobGrade.Job),
										tUserJobs.Grade.GT_EQ(tSubjectJobGrade.MinimumGrade),
									),
								),
							),
					),
					columns.SubjectID.IN(
						tSubjectJobGrade.
							SELECT(tSubjectJobGrade.SubjectID).
							FROM(tSubjectJobGrade).
							WHERE(mysql.AND(
								tSubjectJobGrade.Job.EQ(mysql.String(userInfo.GetJob())),
								tSubjectJobGrade.MinimumGrade.LT_EQ(
									mysql.Int32(userInfo.GetJobGrade()),
								),
							)),
					),
				),
			)),
	)
}

func (a *SubjectObjectAccess) createTargetAccess(
	ctx context.Context,
	tx qrm.DB,
	resolver *SubjectResolver,
	targetID int64,
	in *resourcesaccess.Access,
	opts SubjectAccessOptions,
) error {
	if in == nil {
		return nil
	}

	for _, jobAccess := range in.GetJobs() {
		subjectID, err := resolver.EnsureJobGradeSubject(
			ctx,
			tx,
			jobAccess.GetJob(),
			jobAccess.GetMinimumGrade(),
		)
		if err != nil {
			return err
		}

		if jobAccess.GetAccess() == opts.BlockedAccess {
			if err := a.createDeniedAccess(ctx, tx, targetID, subjectID, opts); err != nil {
				return err
			}
			continue
		}

		if err := a.CreateEntry(
			ctx,
			tx,
			targetID,
			subjectID,
			jobAccess.GetAccess(),
			AccessEffectAllow,
		); err != nil {
			return err
		}
	}

	for _, userAccess := range in.GetUsers() {
		subjectID, err := resolver.EnsureUserSubject(ctx, tx, userAccess.GetUserId())
		if err != nil {
			return err
		}

		if userAccess.GetAccess() == opts.BlockedAccess {
			if err := a.createDeniedAccess(ctx, tx, targetID, subjectID, opts); err != nil {
				return err
			}
			continue
		}

		if err := a.CreateEntry(
			ctx,
			tx,
			targetID,
			subjectID,
			userAccess.GetAccess(),
			AccessEffectAllow,
		); err != nil {
			return err
		}
	}

	for _, qualiAccess := range in.GetQualifications() {
		subjectID, err := resolver.EnsureQualificationSubject(
			ctx,
			tx,
			qualiAccess.GetQualificationId(),
		)
		if err != nil {
			return err
		}

		if qualiAccess.GetAccess() == opts.BlockedAccess {
			if err := a.createDeniedAccess(ctx, tx, targetID, subjectID, opts); err != nil {
				return err
			}
			continue
		}

		if err := a.CreateEntry(
			ctx,
			tx,
			targetID,
			subjectID,
			qualiAccess.GetAccess(),
			AccessEffectAllow,
		); err != nil {
			return err
		}
	}

	return nil
}

func (a *SubjectObjectAccess) createDeniedAccess(
	ctx context.Context,
	tx qrm.DB,
	targetID int64,
	subjectID int64,
	opts SubjectAccessOptions,
) error {
	for _, level := range opts.DeniedAccessLevels {
		if err := a.CreateEntry(ctx, tx, targetID, subjectID, level, AccessEffectDeny); err != nil {
			return err
		}
	}

	return nil
}

func compareSubjectAccess(
	current *resourcesaccess.Access,
	in *resourcesaccess.Access,
) *SubjectAccessChanges {
	changes := &SubjectAccessChanges{
		Jobs:           &SubjectAccessEntryChanges[*resourcesaccess.JobAccess]{},
		Users:          &SubjectAccessEntryChanges[*resourcesaccess.UserAccess]{},
		Qualifications: &SubjectAccessEntryChanges[*resourcesaccess.QualificationAccess]{},
	}
	if current == nil {
		current = &resourcesaccess.Access{}
	}
	if in == nil {
		in = &resourcesaccess.Access{}
	}

	changes.Jobs.ToCreate, changes.Jobs.ToUpdate, changes.Jobs.ToDelete = compareSubjectJobAccess(
		current.GetJobs(),
		in.GetJobs(),
	)
	changes.Users.ToCreate, changes.Users.ToUpdate, changes.Users.ToDelete = compareSubjectUserAccess(
		current.GetUsers(),
		in.GetUsers(),
	)
	changes.Qualifications.ToCreate, changes.Qualifications.ToUpdate, changes.Qualifications.ToDelete = compareSubjectQualificationAccess(
		current.GetQualifications(),
		in.GetQualifications(),
	)

	return changes
}

func compareSubjectJobAccess(
	current []*resourcesaccess.JobAccess,
	in []*resourcesaccess.JobAccess,
) ([]*resourcesaccess.JobAccess, []*resourcesaccess.JobAccess, []*resourcesaccess.JobAccess) {
	return compareSubjectAccessEntries(
		current,
		in,
		func(entry *resourcesaccess.JobAccess) subjectAccessKey {
			return subjectAccessKey{
				job:          entry.GetJob(),
				minimumGrade: entry.GetMinimumGrade(),
			}
		},
		func(entry *resourcesaccess.JobAccess) int64 {
			return entry.GetId()
		},
		func(entry *resourcesaccess.JobAccess) int32 {
			return entry.GetAccess()
		},
		func(entry *resourcesaccess.JobAccess, access int32) {
			entry.SetAccess(access)
		},
	)
}

func compareSubjectUserAccess(
	current []*resourcesaccess.UserAccess,
	in []*resourcesaccess.UserAccess,
) ([]*resourcesaccess.UserAccess, []*resourcesaccess.UserAccess, []*resourcesaccess.UserAccess) {
	return compareSubjectAccessEntries(
		current,
		in,
		func(entry *resourcesaccess.UserAccess) int32 {
			return entry.GetUserId()
		},
		func(entry *resourcesaccess.UserAccess) int64 {
			return entry.GetId()
		},
		func(entry *resourcesaccess.UserAccess) int32 {
			return entry.GetAccess()
		},
		func(entry *resourcesaccess.UserAccess, access int32) {
			entry.SetAccess(access)
		},
	)
}

func compareSubjectQualificationAccess(
	current []*resourcesaccess.QualificationAccess,
	in []*resourcesaccess.QualificationAccess,
) ([]*resourcesaccess.QualificationAccess, []*resourcesaccess.QualificationAccess, []*resourcesaccess.QualificationAccess) {
	return compareSubjectAccessEntries(
		current,
		in,
		func(entry *resourcesaccess.QualificationAccess) int64 {
			return entry.GetQualificationId()
		},
		func(entry *resourcesaccess.QualificationAccess) int64 {
			return entry.GetId()
		},
		func(entry *resourcesaccess.QualificationAccess) int32 {
			return entry.GetAccess()
		},
		func(entry *resourcesaccess.QualificationAccess, access int32) {
			entry.SetAccess(access)
		},
	)
}

type subjectAccessKey struct {
	job          string
	minimumGrade int32
}

func compareSubjectAccessEntries[T any, K comparable](
	current []T,
	in []T,
	keyFn func(T) K,
	idFn func(T) int64,
	accessFn func(T) int32,
	setAccessFn func(T, int32),
) ([]T, []T, []T) {
	toCreate := make([]T, 0, len(in))
	toUpdate := make([]T, 0, len(current))
	toDelete := make([]T, 0, len(current))

	if len(current) == 0 {
		return in, toUpdate, toDelete
	}

	slices.SortFunc(current, func(a, b T) int {
		switch {
		case idFn(a) < idFn(b):
			return -1
		case idFn(a) > idFn(b):
			return 1
		default:
			return 0
		}
	})

	inputByKey := make(map[K]T, len(in))
	for _, inputEntry := range in {
		inputByKey[keyFn(inputEntry)] = inputEntry
	}

	for _, currentEntry := range current {
		inputEntry, found := inputByKey[keyFn(currentEntry)]
		if !found {
			toDelete = append(toDelete, currentEntry)
			continue
		}

		if accessFn(currentEntry) != accessFn(inputEntry) {
			setAccessFn(currentEntry, accessFn(inputEntry))
			toUpdate = append(toUpdate, currentEntry)
		}
	}

	currentKeys := make(map[K]struct{}, len(current))
	for _, currentEntry := range current {
		currentKeys[keyFn(currentEntry)] = struct{}{}
	}

	for _, inputEntry := range in {
		if _, ok := currentKeys[keyFn(inputEntry)]; !ok {
			toCreate = append(toCreate, inputEntry)
		}
	}

	return toCreate, toUpdate, toDelete
}

type subjectAccessRow struct {
	ID          int64
	TargetID    int64
	Access      int32
	Effect      bool
	SubjectType int16

	ACLJob             *string
	ACLMinimumGrade    *int32
	ACLQualificationID *int64

	SubjectUserID   *int32
	UserJob         *string
	UserJobGrade    *int32
	UserFirstname   *string
	UserLastname    *string
	UserDateofbirth *string
	UserPhoneNumber *string
}

func subjectAccessRowsToProto(
	rows []subjectAccessRow,
	opts SubjectAccessOptions,
) *resourcesaccess.Access {
	out := &resourcesaccess.Access{
		Jobs:           make([]*resourcesaccess.JobAccess, 0, len(rows)),
		Users:          make([]*resourcesaccess.UserAccess, 0, len(rows)),
		Qualifications: make([]*resourcesaccess.QualificationAccess, 0, len(rows)),
	}

	blockedJobKeys := map[string]struct{}{}
	blockedUserIDs := map[int32]struct{}{}
	blockedQualificationIDs := map[int64]struct{}{}
	allowedJobs := make([]*resourcesaccess.JobAccess, 0, len(rows))
	allowedUsers := make([]*resourcesaccess.UserAccess, 0, len(rows))
	allowedQualifications := make([]*resourcesaccess.QualificationAccess, 0, len(rows))

	for _, row := range rows {
		switch SubjectType(row.SubjectType) {
		case SubjectTypeJobGrade:
			if row.ACLJob == nil || row.ACLMinimumGrade == nil {
				continue
			}

			entry := &resourcesaccess.JobAccess{
				Id:           row.ID,
				TargetId:     row.TargetID,
				Job:          *row.ACLJob,
				MinimumGrade: *row.ACLMinimumGrade,
				Access:       row.Access,
			}

			if !row.Effect {
				key := subjectJobAccessKey(entry.GetJob(), entry.GetMinimumGrade())
				if _, ok := blockedJobKeys[key]; ok {
					continue
				}
				blockedJobKeys[key] = struct{}{}
				entry.Access = opts.BlockedAccess
				out.Jobs = append(out.Jobs, entry)
				continue
			}

			allowedJobs = append(allowedJobs, entry)

		case SubjectTypeUser:
			if row.SubjectUserID == nil {
				continue
			}

			entry := &resourcesaccess.UserAccess{
				Id:       row.ID,
				TargetId: row.TargetID,
				UserId:   *row.SubjectUserID,
				Access:   row.Access,
				User:     subjectAccessUserShort(row),
			}

			if !row.Effect {
				if _, ok := blockedUserIDs[entry.GetUserId()]; ok {
					continue
				}
				blockedUserIDs[entry.GetUserId()] = struct{}{}
				entry.Access = opts.BlockedAccess
				out.Users = append(out.Users, entry)
				continue
			}

			allowedUsers = append(allowedUsers, entry)
		case SubjectTypeQualification:
			if row.ACLQualificationID == nil {
				continue
			}

			entry := &resourcesaccess.QualificationAccess{
				Id:              row.ID,
				TargetId:        row.TargetID,
				QualificationId: *row.ACLQualificationID,
				Access:          row.Access,
			}

			if !row.Effect {
				if _, ok := blockedQualificationIDs[entry.GetQualificationId()]; ok {
					continue
				}
				blockedQualificationIDs[entry.GetQualificationId()] = struct{}{}
				entry.Access = opts.BlockedAccess
				out.Qualifications = append(out.Qualifications, entry)
				continue
			}

			allowedQualifications = append(allowedQualifications, entry)
		}
	}

	out.Jobs = append(out.Jobs, filterAllowedSubjectJobAccess(allowedJobs, blockedJobKeys)...)
	out.Users = append(out.Users, filterAllowedSubjectUserAccess(allowedUsers, blockedUserIDs)...)
	out.Qualifications = append(
		out.Qualifications,
		filterAllowedSubjectQualificationAccess(allowedQualifications, blockedQualificationIDs)...)

	return out
}

func filterAllowedSubjectQualificationAccess(
	rows []*resourcesaccess.QualificationAccess,
	blocked map[int64]struct{},
) []*resourcesaccess.QualificationAccess {
	out := make([]*resourcesaccess.QualificationAccess, 0, len(rows))
	for _, row := range rows {
		if _, ok := blocked[row.GetQualificationId()]; ok {
			continue
		}
		out = append(out, row)
	}

	return out
}

func subjectAccessUserShort(row subjectAccessRow) *usershort.UserShort {
	if row.SubjectUserID == nil {
		return nil
	}

	user := &usershort.UserShort{
		UserId:      *row.SubjectUserID,
		Job:         derefString(row.UserJob),
		JobGrade:    derefInt32(row.UserJobGrade),
		Firstname:   derefString(row.UserFirstname),
		Lastname:    derefString(row.UserLastname),
		Dateofbirth: derefString(row.UserDateofbirth),
	}
	if row.UserPhoneNumber != nil {
		user.PhoneNumber = row.UserPhoneNumber
	}

	return user
}

func derefString(in *string) string {
	if in == nil {
		return ""
	}
	return *in
}

func derefInt32(in *int32) int32 {
	if in == nil {
		return 0
	}
	return *in
}

func subjectJobAccessKey(job string, minimumGrade int32) string {
	return fmt.Sprintf("%s:%d", job, minimumGrade)
}

func filterAllowedSubjectJobAccess(
	rows []*resourcesaccess.JobAccess,
	blocked map[string]struct{},
) []*resourcesaccess.JobAccess {
	out := make([]*resourcesaccess.JobAccess, 0, len(rows))
	for _, row := range rows {
		if _, ok := blocked[subjectJobAccessKey(row.GetJob(), row.GetMinimumGrade())]; ok {
			continue
		}
		out = append(out, row)
	}

	return out
}

func filterAllowedSubjectUserAccess(
	rows []*resourcesaccess.UserAccess,
	blocked map[int32]struct{},
) []*resourcesaccess.UserAccess {
	out := make([]*resourcesaccess.UserAccess, 0, len(rows))
	for _, row := range rows {
		if _, ok := blocked[row.GetUserId()]; ok {
			continue
		}
		out = append(out, row)
	}

	return out
}
