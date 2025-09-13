package access

import (
	"context"
	"database/sql"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/qualifications"
	pbuserinfo "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/protoutils"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

// Grouped provides grouped access control for jobs, users, and qualifications.
// It holds references to the database, target table, and access logic for each group.
type Grouped[
	JobsU any,
	JobsT JobsAccessProtoMessage[JobsU, V],
	UsersU any,
	UsersT UsersAccessProtoMessage[UsersU, V],
	QualiU any,
	QualiT QualificationsAccessProtoMessage[QualiU, V],
	V protoutils.ProtoEnum,
] struct {
	// db is the SQL database connection.
	db *sql.DB

	// targetTable is the main table for access checks.
	targetTable jet.Table
	// targetTableColumns holds column references for the target table.
	targetTableColumns *TargetTableColumns

	// Jobs provides access logic for job-based permissions.
	Jobs *Jobs[JobsU, JobsT, V]
	// Users provides access logic for user-based permissions.
	Users *Users[UsersU, UsersT, V]
	// Qualifications provides access logic for qualification-based permissions.
	Qualifications *Qualifications[QualiU, QualiT, V]
}

// AccessChangesJobs holds the changes to be made for job-based access.
type AccessChangesJobs[JobsU any, JobsT JobsAccessProtoMessage[JobsU, V], V protoutils.ProtoEnum] struct {
	// ToCreate contains jobs to be created.
	ToCreate []JobsT
	// ToUpdate contains jobs to be updated.
	ToUpdate []JobsT
	// ToDelete contains jobs to be deleted.
	ToDelete []JobsT
}

// IsEmpty returns true if there are no job access changes to apply.
func (a *AccessChangesJobs[UsersU, UsersT, V]) IsEmpty() bool {
	return len(a.ToCreate) == 0 && len(a.ToUpdate) == 0 && len(a.ToDelete) == 0
}

// AccessChangesUsers holds the changes to be made for user-based access.
type AccessChangesUsers[UsersU any, UsersT UsersAccessProtoMessage[UsersU, V], V protoutils.ProtoEnum] struct {
	// ToCreate contains users to be created.
	ToCreate []UsersT
	// ToUpdate contains users to be updated.
	ToUpdate []UsersT
	// ToDelete contains users to be deleted.
	ToDelete []UsersT
}

// IsEmpty returns true if there are no user access changes to apply.
func (a *AccessChangesUsers[UsersU, UsersT, V]) IsEmpty() bool {
	return len(a.ToCreate) == 0 && len(a.ToUpdate) == 0 && len(a.ToDelete) == 0
}

// AccessChangesQualifications holds the changes to be made for qualification-based access.
type AccessChangesQualifications[QualiU any, QualiT QualificationsAccessProtoMessage[QualiU, V], V protoutils.ProtoEnum] struct {
	// ToCreate contains qualifications to be created.
	ToCreate []QualiT
	// ToUpdate contains qualifications to be updated.
	ToUpdate []QualiT
	// ToDelete contains qualifications to be deleted.
	ToDelete []QualiT
}

// IsEmpty returns true if there are no qualification access changes to apply.
func (a *AccessChangesQualifications[QualiU, QualiT, V]) IsEmpty() bool {
	return len(a.ToCreate) == 0 && len(a.ToUpdate) == 0 && len(a.ToDelete) == 0
}

// GroupedAccessChanges aggregates access changes for jobs, users, and qualifications.
type GroupedAccessChanges[JobsU any, JobsT JobsAccessProtoMessage[JobsU, V], UsersU any, UsersT UsersAccessProtoMessage[UsersU, V], QualiU any, QualiT QualificationsAccessProtoMessage[QualiU, V], V protoutils.ProtoEnum] struct {
	// Jobs holds job access changes.
	Jobs *AccessChangesJobs[JobsU, JobsT, V]
	// Users holds user access changes.
	Users *AccessChangesUsers[UsersU, UsersT, V]
	// Qualifications holds qualification access changes.
	Qualifications *AccessChangesQualifications[QualiU, QualiT, V]
}

// IsEmpty returns true if there are no grouped access changes to apply.
func (a *GroupedAccessChanges[JobsU, JobsT, UsersU, UsersT, QualiU, QualiT, V]) IsEmpty() bool {
	return a.Jobs.IsEmpty() && a.Users.IsEmpty() && a.Qualifications.IsEmpty()
}

// NewGrouped creates a new Grouped instance for access control.
func NewGrouped[
	JobsU any,
	JobsT JobsAccessProtoMessage[JobsU, V],
	UsersU any,
	UsersT UsersAccessProtoMessage[UsersU, V],
	QualiU any,
	QualiT QualificationsAccessProtoMessage[QualiU, V],
	V protoutils.ProtoEnum,
](db *sql.DB, targetTable jet.Table, targetTableColumns *TargetTableColumns, jobs *Jobs[JobsU, JobsT, V], users *Users[UsersU, UsersT, V], qualis *Qualifications[QualiU, QualiT, V]) *Grouped[JobsU, JobsT, UsersU, UsersT, QualiU, QualiT, V] {
	return &Grouped[JobsU, JobsT, UsersU, UsersT, QualiU, QualiT, V]{
		db:                 db,
		targetTable:        targetTable,
		targetTableColumns: targetTableColumns,
		Jobs:               jobs,
		Users:              users,
		Qualifications:     qualis,
	}
}

// HandleAccessChanges processes and categorizes access changes for jobs, users, and qualifications.
func (g *Grouped[JobsU, JobsT, UsersU, UsersT, QualiU, QualiT, V]) HandleAccessChanges(
	ctx context.Context,
	tx qrm.DB,
	targetId int64,
	jobsIn []JobsT,
	usersIn []UsersT,
	qualisIn []QualiT,
) (*GroupedAccessChanges[JobsU, JobsT, UsersU, UsersT, QualiU, QualiT, V], error) {
	changes := &GroupedAccessChanges[JobsU, JobsT, UsersU, UsersT, QualiU, QualiT, V]{
		Jobs:           &AccessChangesJobs[JobsU, JobsT, V]{},
		Users:          &AccessChangesUsers[UsersU, UsersT, V]{},
		Qualifications: &AccessChangesQualifications[QualiU, QualiT, V]{},
	}
	var err error

	if g.Jobs != nil {
		if changes.Jobs.ToCreate, changes.Jobs.ToUpdate, changes.Jobs.ToDelete, err = g.Jobs.HandleAccessChanges(ctx, tx, targetId, jobsIn); err != nil {
			return nil, err
		}
	}

	if g.Users != nil {
		if changes.Users.ToCreate, changes.Users.ToUpdate, changes.Users.ToDelete, err = g.Users.HandleAccessChanges(ctx, tx, targetId, usersIn); err != nil {
			return nil, err
		}
	}

	if g.Qualifications != nil {
		if changes.Qualifications.ToCreate, changes.Qualifications.ToUpdate, changes.Qualifications.ToDelete, err = g.Qualifications.HandleAccessChanges(ctx, tx, targetId, qualisIn); err != nil {
			return nil, err
		}
	}

	return changes, nil
}

// CanUserAccessTarget checks if a user can access a specific target based on access rights.
func (g *Grouped[JobsU, JobsT, UsersU, UsersT, QualiU, QualiT, V]) CanUserAccessTarget(
	ctx context.Context,
	targetId int64,
	userInfo *pbuserinfo.UserInfo,
	access V,
) (bool, error) {
	out, err := g.CanUserAccessTargetIDs(ctx, userInfo, access, targetId)
	return len(out) > 0, err
}

// CanUserAccessTargets checks if a user can access multiple targets based on access rights.
func (g *Grouped[JobsU, JobsT, UsersU, UsersT, QualiU, QualiT, V]) CanUserAccessTargets(
	ctx context.Context,
	userInfo *pbuserinfo.UserInfo,
	access V,
	targetIds ...int64,
) (bool, error) {
	out, err := g.CanUserAccessTargetIDs(ctx, userInfo, access, targetIds...)
	return len(out) == len(targetIds), err
}

type canAccessIdsHelper struct {
	IDs []int64 `alias:"id"`
}

// CanUserAccessTargetIDs retrieves target IDs that a user can access based on access rights.
func (g *Grouped[JobsU, JobsT, UsersU, UsersT, QualiU, QualiT, V]) CanUserAccessTargetIDs(
	ctx context.Context,
	userInfo *pbuserinfo.UserInfo,
	access V,
	targetIds ...int64,
) ([]int64, error) {
	if len(targetIds) == 0 {
		return targetIds, nil
	}

	// Allow superusers access to any docs
	if userInfo.GetSuperuser() {
		return targetIds, nil
	}

	stmt := g.getAccessQuery(userInfo, targetIds, access)

	dest := &canAccessIdsHelper{}
	if err := stmt.QueryContext(ctx, g.db, &dest.IDs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest.IDs, nil
}

// getAccessQuery constructs a query to check user access for given target IDs.
func (g *Grouped[JobsU, JobsT, UsersU, UsersT, QualiU, QualiT, V]) getAccessQuery(
	userInfo *pbuserinfo.UserInfo,
	targetIds []int64,
	access V,
) jet.SelectStatement {
	ids := make([]jet.Expression, len(targetIds))
	for i := range targetIds {
		ids[i] = jet.Int64(targetIds[i])
	}

	accessCheckConditions := make([]jet.BoolExpression, 0, 4)
	orderBys := []jet.OrderByClause{g.targetTableColumns.ID.DESC()}

	// Creator-based access (keeps your nil-logic intact)
	creatorCond := jet.Bool(false)
	if g.targetTableColumns.CreatorJob != nil {
		creatorCond = g.targetTableColumns.CreatorJob.EQ(jet.String(userInfo.GetJob()))
	}
	if g.targetTableColumns.CreatorID != nil {
		if g.targetTableColumns.CreatorJob == nil {
			creatorCond = g.targetTableColumns.CreatorID.EQ(jet.Int32(userInfo.GetUserId()))
		} else {
			creatorCond = creatorCond.AND(
				g.targetTableColumns.CreatorID.EQ(jet.Int32(userInfo.GetUserId())),
			)
		}
	}
	accessCheckConditions = append(accessCheckConditions, creatorCond)

	// Direct user access via EXISTS (if Users access table is configured)
	if g.Users != nil {
		userAccessExists := jet.EXISTS(
			jet.
				SELECT(jet.Int(1)).
				FROM(g.Users.table).
				WHERE(
					g.Users.columns.TargetID.EQ(g.targetTableColumns.ID).
						AND(g.Users.columns.Access.GT_EQ(jet.Int32(int32(access.Number())))).
						AND(g.Users.columns.UserId.EQ(jet.Int32(userInfo.GetUserId()))),
				),
		)
		accessCheckConditions = append(accessCheckConditions, userAccessExists)
	}

	// Job + grade access via EXISTS (if Jobs access table is configured)
	if g.Jobs != nil {
		jobAccessExists := jet.EXISTS(
			jet.
				SELECT(jet.Int(1)).
				FROM(g.Jobs.table).
				WHERE(
					g.Jobs.columns.TargetID.EQ(g.targetTableColumns.ID).
						AND(g.Jobs.columns.Access.GT_EQ(jet.Int32(int32(access.Number())))).
						AND(
							jet.AND(
								g.Jobs.columns.Job.EQ(jet.String(userInfo.GetJob())),
								g.Jobs.columns.MinimumGrade.LT_EQ(
									jet.Int32(userInfo.GetJobGrade()),
								),
							),
						),
				),
		)
		accessCheckConditions = append(accessCheckConditions, jobAccessExists)
	}

	// Qualification-based access via EXISTS (if Qualifications are configured)
	if g.Qualifications != nil {
		qualExists := jet.EXISTS(
			jet.
				SELECT(jet.Int(1)).
				FROM(g.Qualifications.table.
					INNER_JOIN(tQualiResults,
						tQualiResults.QualificationID.EQ(g.Qualifications.columns.QualificationId),
					),
				).
				WHERE(
					jet.AND(
						g.Qualifications.columns.QualificationId.IS_NOT_NULL(),
						tQualiResults.DeletedAt.IS_NULL(),
						tQualiResults.QualificationID.EQ(g.Qualifications.columns.QualificationId),
						tQualiResults.UserID.EQ(jet.Int32(userInfo.GetUserId())),
						tQualiResults.Status.EQ(
							jet.Int32(int32(qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL)),
						),
					),
				),
		)
		accessCheckConditions = append(accessCheckConditions, qualExists)
	}

	stmt := g.targetTable.
		SELECT(
			g.targetTableColumns.ID.AS("id"),
		).
		FROM(g.targetTable).
		WHERE(
			jet.AND(
				g.targetTableColumns.ID.IN(ids...),
				g.targetTableColumns.DeletedAt.IS_NULL(),
				jet.OR(accessCheckConditions...),
			),
		).
		ORDER_BY(orderBys...)

	return stmt
}
