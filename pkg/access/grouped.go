package access

import (
	"context"
	"database/sql"
	"errors"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/utils/protoutils"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type Grouped[
	JobsU any,
	JobsT JobsAccessProtoMessage[JobsU, V],
	UsersU any,
	UsersT UsersAccessProtoMessage[UsersU, V],
	QualiU any,
	QualiT QualificationsAccessProtoMessage[QualiU, V],
	V protoutils.ProtoEnum,
] struct {
	db *sql.DB

	targetTable        jet.Table
	targetTableColumns *TargetTableColumns

	Jobs           *Jobs[JobsU, JobsT, V]
	Users          *Users[UsersU, UsersT, V]
	Qualifications *Qualifications[QualiU, QualiT, V]
}

type AccessChangesJobs[JobsU any, JobsT JobsAccessProtoMessage[JobsU, V], V protoutils.ProtoEnum] struct {
	ToCreate []JobsT
	ToUpdate []JobsT
	ToDelete []JobsT
}

func (a *AccessChangesJobs[UsersU, UsersT, V]) IsEmpty() bool {
	return len(a.ToCreate) == 0 && len(a.ToUpdate) == 0 && len(a.ToDelete) == 0
}

type AccessChangesUsers[UsersU any, UsersT UsersAccessProtoMessage[UsersU, V], V protoutils.ProtoEnum] struct {
	ToCreate []UsersT
	ToUpdate []UsersT
	ToDelete []UsersT
}

func (a *AccessChangesUsers[UsersU, UsersT, V]) IsEmpty() bool {
	return len(a.ToCreate) == 0 && len(a.ToUpdate) == 0 && len(a.ToDelete) == 0
}

type AccessChangesQualifications[QualiU any, QualiT QualificationsAccessProtoMessage[QualiU, V], V protoutils.ProtoEnum] struct {
	ToCreate []QualiT
	ToUpdate []QualiT
	ToDelete []QualiT
}

func (a *AccessChangesQualifications[QualiU, QualiT, V]) IsEmpty() bool {
	return len(a.ToCreate) == 0 && len(a.ToUpdate) == 0 && len(a.ToDelete) == 0
}

type GroupedAccessChanges[JobsU any, JobsT JobsAccessProtoMessage[JobsU, V], UsersU any, UsersT UsersAccessProtoMessage[UsersU, V], QualiU any, QualiT QualificationsAccessProtoMessage[QualiU, V], V protoutils.ProtoEnum] struct {
	Jobs           *AccessChangesJobs[JobsU, JobsT, V]
	Users          *AccessChangesUsers[UsersU, UsersT, V]
	Qualifications *AccessChangesQualifications[QualiU, QualiT, V]
}

func (a *GroupedAccessChanges[JobsU, JobsT, UsersU, UsersT, QualiU, QualiT, V]) IsEmpty() bool {
	return a.Jobs.IsEmpty() && a.Users.IsEmpty() && a.Qualifications.IsEmpty()
}

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

func (g *Grouped[JobsU, JobsT, UsersU, UsersT, QualiU, QualiT, V]) HandleAccessChanges(ctx context.Context, tx qrm.DB, targetId uint64, jobsIn []JobsT, usersIn []UsersT, qualisIn []QualiT) (*GroupedAccessChanges[JobsU, JobsT, UsersU, UsersT, QualiU, QualiT, V], error) {
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

func (g *Grouped[JobsU, JobsT, UsersU, UsersT, QualiU, QualiT, V]) CanUserAccessTarget(ctx context.Context, targetId uint64, userInfo *userinfo.UserInfo, access V) (bool, error) {
	out, err := g.CanUserAccessTargetIDs(ctx, userInfo, access, targetId)
	return len(out) > 0, err
}

func (g *Grouped[JobsU, JobsT, UsersU, UsersT, QualiU, QualiT, V]) CanUserAccessTargets(ctx context.Context, userInfo *userinfo.UserInfo, access V, targetIds ...uint64) (bool, error) {
	out, err := g.CanUserAccessTargetIDs(ctx, userInfo, access, targetIds...)
	return len(out) == len(targetIds), err
}

type canAccessIdsHelper struct {
	IDs []uint64 `alias:"id"`
}

func (g *Grouped[JobsU, JobsT, UsersU, UsersT, QualiU, QualiT, V]) CanUserAccessTargetIDs(ctx context.Context, userInfo *userinfo.UserInfo, access V, targetIds ...uint64) ([]uint64, error) {
	if len(targetIds) == 0 {
		return targetIds, nil
	}

	// Allow superusers access to any docs
	if userInfo.SuperUser {
		return targetIds, nil
	}

	stmt := g.GetAccessQuery(userInfo, targetIds, access)

	dest := &canAccessIdsHelper{}
	if err := stmt.QueryContext(ctx, g.db, &dest.IDs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest.IDs, nil
}

func (g *Grouped[JobsU, JobsT, UsersU, UsersT, QualiU, QualiT, V]) GetAccessQuery(userInfo *userinfo.UserInfo, targetIds []uint64, access V) jet.SelectStatement {
	ids := make([]jet.Expression, len(targetIds))
	for i := 0; i < len(targetIds); i++ {
		ids[i] = jet.Uint64(targetIds[i])
	}

	var from jet.ReadableTable
	from = g.targetTable

	accessCheckConditions := []jet.BoolExpression{}
	accessCheckCondition := jet.Bool(false)
	if g.targetTableColumns.CreatorJob != nil {
		accessCheckCondition = g.targetTableColumns.CreatorJob.EQ(jet.String(userInfo.Job))
	}
	if g.targetTableColumns.CreatorID != nil {
		if g.targetTableColumns.CreatorJob == nil {
			accessCheckCondition = g.targetTableColumns.CreatorID.EQ(jet.Int32(userInfo.UserId))
		} else {
			accessCheckCondition = accessCheckCondition.AND(g.targetTableColumns.CreatorID.EQ(jet.Int32(userInfo.UserId)))
		}
	}
	accessCheckConditions = append(accessCheckConditions, accessCheckCondition)

	orderBys := []jet.OrderByClause{g.targetTableColumns.ID.DESC()}

	if g.Jobs != nil {
		from = from.
			LEFT_JOIN(g.Jobs.table,
				g.Jobs.columns.TargetID.EQ(g.targetTableColumns.ID).
					AND(g.Jobs.columns.Job.EQ(jet.String(userInfo.Job))).
					AND(g.Jobs.columns.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
			)

		condition := []jet.BoolExpression{
			g.Jobs.columns.Access.IS_NOT_NULL(),
			g.Jobs.columns.Access.GT_EQ(jet.Int32(int32(access.Number()))),
		}
		if g.Users != nil {
			condition = append(condition, g.Users.columns.Access.IS_NULL())
		}
		accessCheckConditions = append(accessCheckConditions, jet.AND(condition...))

		orderBys = append(orderBys, g.Jobs.columns.MinimumGrade)
	}

	if g.Users != nil {
		from = from.
			LEFT_JOIN(g.Users.table,
				g.Users.columns.TargetID.EQ(g.targetTableColumns.ID).
					AND(g.Users.columns.UserId.EQ(jet.Int32(userInfo.UserId))),
			)

		condition := []jet.BoolExpression{
			g.Users.columns.Access.IS_NOT_NULL(),
			g.Users.columns.Access.GT_EQ(jet.Int32(int32(access.Number()))),
		}
		accessCheckConditions = append(accessCheckConditions, jet.AND(condition...))
	}

	if g.Qualifications != nil {
		from = from.
			LEFT_JOIN(g.Qualifications.table,
				g.Qualifications.columns.TargetID.EQ(g.targetTableColumns.ID),
			).
			LEFT_JOIN(tQualiResults,
				tQualiResults.QualificationID.EQ(g.Qualifications.columns.QualificationId).
					AND(tQualiResults.UserID.EQ(jet.Int32(userInfo.UserId))),
			)

		condition := []jet.BoolExpression{
			g.Qualifications.columns.Access.IS_NOT_NULL(),
			g.Qualifications.columns.Access.GT_EQ(jet.Int32(int32(access.Number()))),
			tQualiResults.DeletedAt.IS_NULL(),
			tQualiResults.QualificationID.EQ(g.Qualifications.columns.QualificationId),
			tQualiResults.Status.EQ(jet.Int32(int32(qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL.Number()))),
		}
		accessCheckConditions = append(accessCheckConditions, jet.AND(condition...))
	}

	stmt := g.targetTable.
		SELECT(
			g.targetTableColumns.ID.AS("id"),
		).
		FROM(from).
		WHERE(jet.AND(
			g.targetTableColumns.ID.IN(ids...),
			g.targetTableColumns.DeletedAt.IS_NULL(),
			jet.OR(accessCheckConditions...),
		)).
		GROUP_BY(g.targetTableColumns.ID).
		ORDER_BY(orderBys...)

	return stmt
}
