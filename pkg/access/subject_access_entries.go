package access

import (
	"context"
	"errors"
	"slices"
	"strconv"

	resourcesaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/access"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
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
	tQualiSuccess := table.FivenetQualificationsResultSuccessMap.AS(
		"subject_acl_qualification_results_exists",
	)
	tSubjectJobGrade := table.FivenetACLSubjectJobGradeScopes.AS("subject_acl_job_grade_exists")
	tUserJobs := table.FivenetUserJobs.AS("subject_acl_user_jobs_exists")

	actorSubjectsSelect := mysql.
		SELECT(
			tSubjectUsers.SubjectID.AS("subject_id"),
			mysql.Int32(SubjectSpecificityUser).
				AS("specificity"),
			mysql.Int32(-1).AS("grade_specificity"),
		).
		FROM(tSubjectUsers).
		WHERE(tSubjectUsers.UserID.EQ(mysql.Int32(userInfo.GetUserId()))).
		UNION(
			mysql.
				SELECT(
					tSubjectQualis.SubjectID.AS("subject_id"),
					mysql.Int32(SubjectSpecificityQualification).
						AS("specificity"),
					mysql.Int32(-1).AS("grade_specificity"),
				).
				FROM(tSubjectQualis.
					INNER_JOIN(tQualiSuccess, mysql.AND(
						tQualiSuccess.QualificationID.EQ(tSubjectQualis.QualificationID),
						tQualiSuccess.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
					),
					),
				),
		).
		UNION_ALL(
			mysql.
				SELECT(
					tSubjectJobGrade.SubjectID.AS("subject_id"),
					mysql.Int32(SubjectSpecificityJobGrade).AS("specificity"),
					tSubjectJobGrade.MinimumGrade.AS("grade_specificity"),
				).
				FROM(tSubjectJobGrade.
					INNER_JOIN(tUserJobs, mysql.AND(
						tUserJobs.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
						tUserJobs.Job.EQ(tSubjectJobGrade.Job),
						tUserJobs.Grade.GT_EQ(tSubjectJobGrade.MinimumGrade),
						tUserJobs.IsPrimary.IS_TRUE(),
					)),
				),
		).
		UNION_ALL(
			mysql.
				SELECT(
					tSubjectJobGrade.SubjectID.AS("subject_id"),
					mysql.Int32(SubjectSpecificityJobGrade).AS("specificity"),
					tSubjectJobGrade.MinimumGrade.AS("grade_specificity"),
				).
				FROM(tSubjectJobGrade).
				WHERE(mysql.AND(
					tSubjectJobGrade.Job.EQ(mysql.String(userInfo.GetJob())),
					tSubjectJobGrade.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),
				)),
		)

	actorSubjectsTable := actorSubjectsSelect.AsTable("acl_access_actor_subjects")
	actorSubjectID := mysql.IntegerColumn("subject_id").From(actorSubjectsTable)
	actorSpecificity := mysql.IntegerColumn("specificity").From(actorSubjectsTable)
	actorGradeSpecificity := mysql.IntegerColumn("grade_specificity").From(actorSubjectsTable)

	matchingACLSelect := mysql.
		SELECT(
			columns.TargetID.AS("target_id"),
			columns.Effect.AS("effect"),
			columns.Access.AS("access"),
			actorSpecificity.AS("specificity"),
			mysql.COALESCE(actorGradeSpecificity, mysql.Int32(-1)).AS("grade_specificity"),
			mysql.ROW_NUMBER().OVER(
				mysql.PARTITION_BY(columns.TargetID).
					ORDER_BY(
						actorSpecificity.DESC(),
						mysql.COALESCE(actorGradeSpecificity, mysql.Int32(-1)).DESC(),
						columns.Effect.ASC(),
					),
			).AS("access_rank"),
		).
		FROM(tAccess.
			INNER_JOIN(actorSubjectsTable,
				actorSubjectID.EQ(columns.SubjectID),
			),
		).
		WHERE(
			// FIXME previously `columns.TargetID.EQ(targetID),` was taken into account but not anymore
			columns.Access.GT_EQ(mysql.Int32(access)),
		).
		DISTINCT()

	matchingACLTable := matchingACLSelect.AsTable("acl_access_matching")
	matchingTargetID := mysql.IntegerColumn("target_id").From(matchingACLTable)
	matchingEffect := mysql.BoolColumn("effect").From(matchingACLTable)
	matchingRank := mysql.IntegerColumn("access_rank").From(matchingACLTable)

	winningACLSelect := mysql.
		SELECT(
			matchingTargetID.AS("target_id"),
			matchingEffect.AS("effect"),
			matchingRank.AS("access_rank"),
		).
		FROM(matchingACLTable).
		WHERE(matchingRank.EQ(mysql.Int(1)))

	winningACLTable := winningACLSelect.AsTable("acl_access_winning")
	winningTargetID := mysql.IntegerColumn("target_id").From(winningACLTable)
	winningEffect := mysql.BoolColumn("effect").From(winningACLTable)

	return mysql.EXISTS(
		mysql.
			SELECT(mysql.Int(1)).
			FROM(winningACLTable).
			WHERE(mysql.AND(
				winningTargetID.EQ(targetID),
				winningEffect.IS_TRUE(),
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
	access := opts.BlockedAccess
	if _, maxAccess, ok := accessRangeBounds(opts.DeniedAccessLevels); ok {
		access = maxAccess
	}

	if err := a.CreateEntry(ctx, tx, targetID, subjectID, access, AccessEffectDeny); err != nil {
		return err
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
		func(entry *resourcesaccess.JobAccess) int64 {
			return entry.GetId()
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
		func(entry *resourcesaccess.UserAccess) int64 {
			return entry.GetId()
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
			return entry.GetId()
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

	blockedJobIndexes := map[subjectJobAccessKeyValue]int{}
	blockedJobAccesses := map[subjectJobAccessKeyValue]int32{}
	blockedUserIndexes := map[int32]int{}
	blockedUserAccesses := map[int32]int32{}
	blockedQualificationIndexes := map[int64]int{}
	blockedQualificationAccesses := map[int64]int32{}
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
				Id:             row.ID,
				TargetId:       row.TargetID,
				Job:            *row.ACLJob,
				MinimumGrade:   *row.ACLMinimumGrade,
				Access:         row.Access,
				RequiredAccess: requiredAccessPtr(row.Access),
			}

			if !row.Effect {
				key := subjectJobAccessKeyValue{
					job:          entry.GetJob(),
					minimumGrade: entry.GetMinimumGrade(),
				}
				if idx, ok := blockedJobIndexes[key]; ok {
					if row.Access <= blockedJobAccesses[key] {
						continue
					}
					blockedJobAccesses[key] = row.Access
					out.Jobs[idx].Id = row.ID
					continue
				}
				blockedJobIndexes[key] = len(out.Jobs)
				blockedJobAccesses[key] = row.Access
				entry.Access = opts.BlockedAccess
				entry.RequiredAccess = requiredAccessPtr(deniedAccessFloor(opts, row.Access))
				out.Jobs = append(out.Jobs, entry)
				continue
			}

			allowedJobs = append(allowedJobs, entry)

		case SubjectTypeUser:
			if row.SubjectUserID == nil {
				continue
			}

			entry := &resourcesaccess.UserAccess{
				Id:             row.ID,
				TargetId:       row.TargetID,
				UserId:         *row.SubjectUserID,
				Access:         row.Access,
				RequiredAccess: requiredAccessPtr(row.Access),
				User:           subjectAccessUserShort(row),
			}

			if !row.Effect {
				if idx, ok := blockedUserIndexes[entry.GetUserId()]; ok {
					if row.Access <= blockedUserAccesses[entry.GetUserId()] {
						continue
					}
					blockedUserAccesses[entry.GetUserId()] = row.Access
					out.Users[idx].Id = row.ID
					continue
				}
				blockedUserIndexes[entry.GetUserId()] = len(out.Users)
				blockedUserAccesses[entry.GetUserId()] = row.Access
				entry.Access = opts.BlockedAccess
				entry.RequiredAccess = requiredAccessPtr(deniedAccessFloor(opts, row.Access))
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
				RequiredAccess:  requiredAccessPtr(row.Access),
			}

			if !row.Effect {
				if idx, ok := blockedQualificationIndexes[entry.GetQualificationId()]; ok {
					if row.Access <= blockedQualificationAccesses[entry.GetQualificationId()] {
						continue
					}
					blockedQualificationAccesses[entry.GetQualificationId()] = row.Access
					out.Qualifications[idx].Id = row.ID
					continue
				}
				blockedQualificationIndexes[entry.GetQualificationId()] = len(out.Qualifications)
				blockedQualificationAccesses[entry.GetQualificationId()] = row.Access
				entry.Access = opts.BlockedAccess
				entry.RequiredAccess = requiredAccessPtr(deniedAccessFloor(opts, row.Access))
				out.Qualifications = append(out.Qualifications, entry)
				continue
			}

			allowedQualifications = append(allowedQualifications, entry)
		}
	}

	out.Jobs = append(out.Jobs, allowedJobs...)
	out.Users = append(out.Users, allowedUsers...)
	out.Qualifications = append(out.Qualifications, allowedQualifications...)

	return out
}

func subjectAccessUserShort(row subjectAccessRow) *usershort.UserShort {
	if row.SubjectUserID == nil {
		return nil
	}

	user := &usershort.UserShort{
		UserId:      *row.SubjectUserID,
		Job:         utils.Deref(row.UserJob),
		JobGrade:    utils.Deref(row.UserJobGrade),
		Firstname:   utils.Deref(row.UserFirstname),
		Lastname:    utils.Deref(row.UserLastname),
		Dateofbirth: utils.Deref(row.UserDateofbirth),
	}
	if row.UserPhoneNumber != nil {
		user.PhoneNumber = row.UserPhoneNumber
	}

	return user
}

type subjectJobAccessKeyValue struct {
	job          string
	minimumGrade int32
}

func requiredAccessPtr(access int32) *int32 {
	v := access
	return &v
}

func deniedAccessFloor(opts SubjectAccessOptions, access int32) int32 {
	if floor, _, ok := accessRangeBounds(opts.DeniedAccessLevels); ok {
		return floor
	}
	return access
}

func accessRangeBounds(levels []int32) (int32, int32, bool) {
	if len(levels) == 0 {
		return 0, 0, false
	}

	minLevel := levels[0]
	maxLevel := levels[0]
	for _, level := range levels[1:] {
		if level < minLevel {
			minLevel = level
		}
		if level > maxLevel {
			maxLevel = level
		}
	}

	return minLevel, maxLevel, true
}

func subjectJobAccessKey(job string, minimumGrade int32) string {
	return job + ":" + strconv.FormatInt(int64(minimumGrade), 10)
}
