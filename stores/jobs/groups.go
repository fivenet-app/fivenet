package jobsstore

import (
	"context"
	"errors"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var tJobGroups = table.FivenetJobGroups.AS("job_group")

func buildGroupStates(q GroupsQuery) []jobsgroups.GroupState {
	if len(q.States) > 0 {
		return q.States
	}

	states := []jobsgroups.GroupState{
		jobsgroups.GroupState_GROUP_STATE_ACTIVE,
	}
	if q.IncludeInactive {
		states = append(states, jobsgroups.GroupState_GROUP_STATE_INACTIVE)
	}
	if q.IncludeArchived {
		states = append(states, jobsgroups.GroupState_GROUP_STATE_ARCHIVED)
	}
	return states
}

func buildGroupCondition(q GroupsQuery) mysql.BoolExpression {
	condition := mysql.AND(
		tJobGroups.JobID.EQ(mysql.Int64(q.JobID)),
	)

	states := buildGroupStates(q)
	if len(states) > 0 {
		expressions := make([]mysql.Expression, len(states))
		for i := range states {
			expressions[i] = mysql.Int32(int32(states[i]))
		}
		condition = condition.AND(tJobGroups.State.IN(expressions...))
	}

	if search := dbutils.PrepareForLikeSearch(q.Search); search != "" {
		like := mysql.String(search)
		condition = condition.AND(mysql.OR(
			tJobGroups.Name.LIKE(like),
			tJobGroups.ShortName.LIKE(like),
			tJobGroups.Description.LIKE(like),
		))
	}

	return condition
}

func buildGroupOrderBy(sort *database.Sort) []mysql.OrderByClause {
	orderBys := []mysql.OrderByClause{}
	if sort != nil && len(sort.GetColumns()) > 0 {
		for _, sc := range sort.GetColumns() {
			var columns []mysql.Column
			switch sc.GetId() {
			case "name":
				columns = append(columns, tJobGroups.Name, tJobGroups.SortOrder, tJobGroups.ID)
			case "state":
				columns = append(columns, tJobGroups.State, tJobGroups.SortOrder, tJobGroups.Name)
			case "sort_order":
				columns = append(columns, tJobGroups.SortOrder, tJobGroups.Name, tJobGroups.ID)
			case "updated_at":
				columns = append(
					columns,
					tJobGroups.UpdatedAt,
					tJobGroups.SortOrder,
					tJobGroups.Name,
				)
			case "created_at":
				columns = append(
					columns,
					tJobGroups.CreatedAt,
					tJobGroups.SortOrder,
					tJobGroups.Name,
				)
			case "members_count":
				columns = append(
					columns,
					tJobGroups.MembersCount,
					tJobGroups.SortOrder,
					tJobGroups.Name,
				)
			case "leaders_count":
				columns = append(
					columns,
					tJobGroups.LeadersCount,
					tJobGroups.SortOrder,
					tJobGroups.Name,
				)
			case "rules_count":
				columns = append(
					columns,
					tJobGroups.RulesCount,
					tJobGroups.SortOrder,
					tJobGroups.Name,
				)
			case "exclusions_count":
				columns = append(
					columns,
					tJobGroups.ExclusionsCount,
					tJobGroups.SortOrder,
					tJobGroups.Name,
				)
			case "id":
				fallthrough
			default:
				columns = append(columns, tJobGroups.SortOrder, tJobGroups.Name, tJobGroups.ID)
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

	if len(orderBys) == 0 {
		orderBys = append(orderBys,
			tJobGroups.SortOrder.ASC(),
			tJobGroups.Name.ASC(),
			tJobGroups.ID.ASC(),
		)
	}

	return orderBys
}

func normalizeGroupForInsert(group *jobsgroups.Group) {
	if group == nil {
		return
	}
	if group.GetType() == jobsgroups.GroupType_GROUP_TYPE_UNSPECIFIED {
		group.Type = jobsgroups.GroupType_GROUP_TYPE_MANUAL
	}
	if group.GetState() == jobsgroups.GroupState_GROUP_STATE_UNSPECIFIED {
		group.State = jobsgroups.GroupState_GROUP_STATE_ACTIVE
	}
	if group.GetMembershipMode() == jobsgroups.GroupMembershipMode_GROUP_MEMBERSHIP_MODE_UNSPECIFIED {
		group.MembershipMode = jobsgroups.GroupMembershipMode_GROUP_MEMBERSHIP_MODE_FLEXIBLE
	}
}

func (s *Store) CountGroups(ctx context.Context, db qrm.DB, q GroupsQuery) (int64, error) {
	countStmt := tJobGroups.
		SELECT(mysql.COUNT(tJobGroups.ID).AS("data_count.total")).
		FROM(tJobGroups).
		WHERE(buildGroupCondition(q))

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return count.Total, nil
}

func (s *Store) ListGroups(
	ctx context.Context,
	db qrm.DB,
	q GroupsQuery,
) ([]*jobsgroups.Group, error) {
	stmt := tJobGroups.
		SELECT(
			tJobGroups.ID,
			tJobGroups.JobID,
			tJobGroups.Name,
			tJobGroups.Description,
			tJobGroups.ShortName,
			tJobGroups.LogoFileID,
			tJobGroups.Color,
			tJobGroups.Type,
			tJobGroups.State,
			tJobGroups.MembershipMode,
			tJobGroups.SortOrder,
			tJobGroups.MembersCount,
			tJobGroups.LeadersCount,
			tJobGroups.RulesCount,
			tJobGroups.ExclusionsCount,
			tJobGroups.CreatedByUserID,
			tJobGroups.UpdatedByUserID,
			tJobGroups.CreatedAt,
			tJobGroups.UpdatedAt,
			tJobGroups.ArchivedAt,
		).
		FROM(tJobGroups).
		WHERE(buildGroupCondition(q)).
		OFFSET(q.Offset).
		ORDER_BY(buildGroupOrderBy(q.Sort)...).
		LIMIT(q.Limit)

	groups := []*jobsgroups.Group{}
	if err := stmt.QueryContext(ctx, db, &groups); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return groups, nil
}

func (s *Store) GetGroup(
	ctx context.Context,
	db qrm.DB,
	q GroupQuery,
	id int64,
) (*jobsgroups.Group, error) {
	condition := mysql.AND(
		tJobGroups.JobID.EQ(mysql.Int64(q.JobID)),
		tJobGroups.ID.EQ(mysql.Int64(id)),
	)
	if !q.IncludeArchived {
		condition = condition.AND(
			tJobGroups.State.NOT_EQ(mysql.Int32(int32(jobsgroups.GroupState_GROUP_STATE_ARCHIVED))),
		)
	}

	stmt := tJobGroups.
		SELECT(
			tJobGroups.ID,
			tJobGroups.JobID,
			tJobGroups.Name,
			tJobGroups.Description,
			tJobGroups.ShortName,
			tJobGroups.LogoFileID,
			tJobGroups.Color,
			tJobGroups.Type,
			tJobGroups.State,
			tJobGroups.MembershipMode,
			tJobGroups.SortOrder,
			tJobGroups.MembersCount,
			tJobGroups.LeadersCount,
			tJobGroups.RulesCount,
			tJobGroups.ExclusionsCount,
			tJobGroups.CreatedByUserID,
			tJobGroups.UpdatedByUserID,
			tJobGroups.CreatedAt,
			tJobGroups.UpdatedAt,
			tJobGroups.ArchivedAt,
		).
		FROM(tJobGroups).
		WHERE(condition).
		LIMIT(1)

	dest := &jobsgroups.Group{}
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

func (s *Store) CreateGroup(
	ctx context.Context,
	db qrm.DB,
	group *jobsgroups.Group,
) (int64, error) {
	normalizeGroupForInsert(group)

	insertStmt := tJobGroups.
		INSERT(
			tJobGroups.JobID,
			tJobGroups.Name,
			tJobGroups.Description,
			tJobGroups.ShortName,
			tJobGroups.LogoFileID,
			tJobGroups.Color,
			tJobGroups.Type,
			tJobGroups.State,
			tJobGroups.MembershipMode,
			tJobGroups.SortOrder,
			tJobGroups.CreatedByUserID,
			tJobGroups.UpdatedByUserID,
		).
		VALUES(
			group.GetJobId(),
			group.GetName(),
			dbutils.StringPP(group.Description),
			dbutils.StringPP(group.ShortName),
			dbutils.StringPP(group.LogoFileId),
			dbutils.StringPP(group.Color),
			int32(group.GetType()),
			int32(group.GetState()),
			int32(group.GetMembershipMode()),
			group.GetSortOrder(),
			group.GetCreatedByUserId(),
			dbutils.Int64P(group.GetUpdatedByUserId()),
		)

	res, err := insertStmt.ExecContext(ctx, db)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Store) UpdateGroup(ctx context.Context, db qrm.DB, group *jobsgroups.Group) error {
	updateStmt := tJobGroups.
		UPDATE(
			tJobGroups.Name,
			tJobGroups.Description,
			tJobGroups.ShortName,
			tJobGroups.LogoFileID,
			tJobGroups.Color,
			tJobGroups.Type,
			tJobGroups.State,
			tJobGroups.MembershipMode,
			tJobGroups.SortOrder,
			tJobGroups.UpdatedByUserID,
			tJobGroups.UpdatedAt,
		).
		SET(
			tJobGroups.Name.SET(mysql.String(group.GetName())),
			tJobGroups.Description.SET(dbutils.StringPP(group.Description)),
			tJobGroups.ShortName.SET(dbutils.StringPP(group.ShortName)),
			tJobGroups.LogoFileID.SET(dbutils.StringPP(group.LogoFileId)),
			tJobGroups.Color.SET(dbutils.StringPP(group.Color)),
			tJobGroups.Type.SET(mysql.Int32(int32(group.GetType()))),
			tJobGroups.State.SET(mysql.Int32(int32(group.GetState()))),
			tJobGroups.MembershipMode.SET(mysql.Int32(int32(group.GetMembershipMode()))),
			tJobGroups.SortOrder.SET(mysql.Int32(group.GetSortOrder())),
			tJobGroups.UpdatedByUserID.SET(dbutils.Int64P(group.GetUpdatedByUserId())),
			tJobGroups.UpdatedAt.SET(mysql.CURRENT_TIMESTAMP()),
		).
		WHERE(mysql.AND(
			tJobGroups.ID.EQ(mysql.Int64(group.GetId())),
			tJobGroups.JobID.EQ(mysql.Int64(group.GetJobId())),
		)).
		LIMIT(1)

	_, err := updateStmt.ExecContext(ctx, db)
	return err
}

func (s *Store) ArchiveGroup(
	ctx context.Context,
	db qrm.DB,
	jobID int64,
	id int64,
	updatedByUserID int64,
) error {
	updateStmt := tJobGroups.
		UPDATE(
			tJobGroups.State,
			tJobGroups.ArchivedAt,
			tJobGroups.UpdatedByUserID,
			tJobGroups.UpdatedAt,
		).
		SET(
			tJobGroups.State.SET(mysql.Int32(int32(jobsgroups.GroupState_GROUP_STATE_ARCHIVED))),
			tJobGroups.ArchivedAt.SET(mysql.CURRENT_TIMESTAMP()),
			tJobGroups.UpdatedByUserID.SET(dbutils.Int64P(updatedByUserID)),
			tJobGroups.UpdatedAt.SET(mysql.CURRENT_TIMESTAMP()),
		).
		WHERE(mysql.AND(
			tJobGroups.ID.EQ(mysql.Int64(id)),
			tJobGroups.JobID.EQ(mysql.Int64(jobID)),
		)).
		LIMIT(1)

	_, err := updateStmt.ExecContext(ctx, db)
	return err
}

func (s *Store) RestoreGroup(
	ctx context.Context,
	db qrm.DB,
	jobID int64,
	id int64,
	updatedByUserID int64,
) error {
	updateStmt := tJobGroups.
		UPDATE(
			tJobGroups.State,
			tJobGroups.ArchivedAt,
			tJobGroups.UpdatedByUserID,
			tJobGroups.UpdatedAt,
		).
		SET(
			tJobGroups.State.SET(mysql.Int32(int32(jobsgroups.GroupState_GROUP_STATE_ACTIVE))),
			tJobGroups.ArchivedAt.SET(mysql.TimestampExp(mysql.NULL)),
			tJobGroups.UpdatedByUserID.SET(dbutils.Int64P(updatedByUserID)),
			tJobGroups.UpdatedAt.SET(mysql.CURRENT_TIMESTAMP()),
		).
		WHERE(mysql.AND(
			tJobGroups.ID.EQ(mysql.Int64(id)),
			tJobGroups.JobID.EQ(mysql.Int64(jobID)),
		)).
		LIMIT(1)

	_, err := updateStmt.ExecContext(ctx, db)
	return err
}
