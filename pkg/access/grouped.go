package access

import (
	"context"
	"database/sql"
	"errors"

	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/utils/protoutils"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type Grouped[JobsU any, JobsT JobsAccessProtoMessage[JobsU, V], UsersU any, UsersT UsersAccessProtoMessage[UsersU, V], V protoutils.ProtoEnum] struct {
	db *sql.DB

	targetTable        jet.Table
	targetTableColumns *TargetTableColumns

	Jobs  *Jobs[JobsU, JobsT, V]
	Users *Users[UsersU, UsersT, V]
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

type GroupedAccessChanges[JobsU any, JobsT JobsAccessProtoMessage[JobsU, V], UsersU any, UsersT UsersAccessProtoMessage[UsersU, V], V protoutils.ProtoEnum] struct {
	Jobs  *AccessChangesJobs[JobsU, JobsT, V]
	Users *AccessChangesUsers[UsersU, UsersT, V]
}

func (a *GroupedAccessChanges[JobsU, JobsT, UsersU, UsersT, V]) IsEmpty() bool {
	return a.Jobs.IsEmpty() && a.Users.IsEmpty()
}

func NewGrouped[JobsU any, JobsT JobsAccessProtoMessage[JobsU, V], UsersU any, UsersT UsersAccessProtoMessage[UsersU, V], V protoutils.ProtoEnum](db *sql.DB, targetTable jet.Table, targetTableColumns *TargetTableColumns, jobs *Jobs[JobsU, JobsT, V], users *Users[UsersU, UsersT, V]) *Grouped[JobsU, JobsT, UsersU, UsersT, V] {
	return &Grouped[JobsU, JobsT, UsersU, UsersT, V]{
		db:                 db,
		targetTable:        targetTable,
		targetTableColumns: targetTableColumns,
		Jobs:               jobs,
		Users:              users,
	}
}

func (g *Grouped[JobsU, JobsT, UsersU, UsersT, V]) HandleAccessChanges(ctx context.Context, tx qrm.DB, targetId uint64, jobsIn []JobsT, usersIn []UsersT) (*GroupedAccessChanges[JobsU, JobsT, UsersU, UsersT, V], error) {
	changes := &GroupedAccessChanges[JobsU, JobsT, UsersU, UsersT, V]{
		Jobs:  &AccessChangesJobs[JobsU, JobsT, V]{},
		Users: &AccessChangesUsers[UsersU, UsersT, V]{},
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

	return changes, nil
}

func (g *Grouped[JobsU, JobsT, UsersU, UsersT, V]) CanUserAccessTarget(ctx context.Context, targetId uint64, userInfo *userinfo.UserInfo, access V) (bool, error) {
	out, err := g.CanUserAccessTargetIDs(ctx, userInfo, access, targetId)
	return len(out) > 0, err
}

func (g *Grouped[JobsU, JobsT, UsersU, UsersT, V]) CanUserAccessTargets(ctx context.Context, userInfo *userinfo.UserInfo, access V, targetIds ...uint64) (bool, error) {
	out, err := g.CanUserAccessTargetIDs(ctx, userInfo, access, targetIds...)
	return len(out) == len(targetIds), err
}

type canAccessIdsHelper struct {
	IDs []uint64 `alias:"id"`
}

func (g *Grouped[JobsU, JobsT, UsersU, UsersT, V]) CanUserAccessTargetIDs(ctx context.Context, userInfo *userinfo.UserInfo, access V, targetIds ...uint64) ([]uint64, error) {
	if len(targetIds) == 0 {
		return targetIds, nil
	}

	// Allow superusers access to any docs
	if userInfo.SuperUser {
		return targetIds, nil
	}

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

	dest := &canAccessIdsHelper{}
	if err := stmt.QueryContext(ctx, g.db, &dest.IDs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest.IDs, nil
}
