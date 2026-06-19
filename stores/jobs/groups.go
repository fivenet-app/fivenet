package jobsstore

import (
	"context"
	"errors"
	"strings"

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
		tJobGroups.Job.EQ(mysql.String(q.Job)),
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

	if len(q.IDs) > 0 {
		// TODO improve logic by sorting the ids first
		expressions := make([]mysql.Expression, len(q.IDs))
		for i := range q.IDs {
			expressions = append(expressions, mysql.Int64(q.IDs[i]))
		}
		condition = condition.AND(tJobGroups.ID.IN(expressions...))
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
				columns = append(columns, tJobGroups.Name, tJobGroups.SortRank, tJobGroups.ID)
			case "state":
				columns = append(columns, tJobGroups.State, tJobGroups.SortRank, tJobGroups.Name)
			case "sort_rank":
				columns = append(columns, tJobGroups.SortRank, tJobGroups.Name, tJobGroups.ID)
			case "updated_at":
				columns = append(
					columns,
					tJobGroups.UpdatedAt,
					tJobGroups.SortRank,
					tJobGroups.Name,
				)
			case "created_at":
				columns = append(
					columns,
					tJobGroups.CreatedAt,
					tJobGroups.SortRank,
					tJobGroups.Name,
				)
			case "members_count":
				columns = append(
					columns,
					tJobGroups.MembersCount,
					tJobGroups.SortRank,
					tJobGroups.Name,
				)
			case "leaders_count":
				columns = append(
					columns,
					tJobGroups.LeadersCount,
					tJobGroups.SortRank,
					tJobGroups.Name,
				)
			case "rules_count":
				columns = append(
					columns,
					tJobGroups.RulesCount,
					tJobGroups.SortRank,
					tJobGroups.Name,
				)
			case "exclusions_count":
				columns = append(
					columns,
					tJobGroups.ExclusionsCount,
					tJobGroups.SortRank,
					tJobGroups.Name,
				)
			case "id":
				fallthrough
			default:
				columns = append(columns, tJobGroups.SortRank, tJobGroups.Name, tJobGroups.ID)
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
		orderBys = append(
			orderBys,
			tJobGroups.SortRank.ASC(),
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
	tLogoFile := table.FivenetFiles.AS("logo_file")
	stmt := tJobGroups.
		SELECT(
			tJobGroups.ID.AS("group.id"),
			tJobGroups.Job.AS("group.job"),
			tJobGroups.Name.AS("group.name"),
			tJobGroups.Description.AS("group.description"),
			tJobGroups.ShortName.AS("group.short_name"),
			tJobGroups.LogoFileID.AS("group.logo_file_id"),
			tJobGroups.Color.AS("group.color"),
			tJobGroups.Type.AS("group.type"),
			tJobGroups.State.AS("group.state"),
			tJobGroups.MembershipMode.AS("group.membership_mode"),
			tJobGroups.SortRank.AS("group.sort_rank"),
			tJobGroups.MembersCount.AS("group.members_count"),
			tJobGroups.LeadersCount.AS("group.leaders_count"),
			tJobGroups.RulesCount.AS("group.rules_count"),
			tJobGroups.ExclusionsCount.AS("group.exclusions_count"),
			tJobGroups.CreatedByUserID.AS("group.created_by_user_id"),
			tJobGroups.UpdatedByUserID.AS("group.updated_by_user_id"),
			tJobGroups.CreatedAt.AS("group.created_at"),
			tJobGroups.UpdatedAt.AS("group.updated_at"),
			tJobGroups.DeletedAt.AS("group.deleted_at"),
			tLogoFile.FilePath.AS("logo_file.file_path"),
		).
		FROM(
			tJobGroups.
				LEFT_JOIN(tLogoFile,
					mysql.AND(
						tLogoFile.ID.EQ(tJobGroups.LogoFileID),
						tLogoFile.DeletedAt.IS_NULL(),
					),
				),
		).
		WHERE(buildGroupCondition(q)).
		ORDER_BY(buildGroupOrderBy(q.Sort)...).
		LIMIT(q.Limit).
		OFFSET(q.Offset)

	dest := []*jobsgroups.Group{}
	if err := stmt.QueryContext(ctx, db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}

func (s *Store) GetGroup(
	ctx context.Context,
	db qrm.DB,
	q GroupQuery,
	id int64,
) (*jobsgroups.Group, error) {
	tLogoFile := table.FivenetFiles.AS("logo_file")
	condition := mysql.AND(
		tJobGroups.Job.EQ(mysql.String(q.Job)),
		tJobGroups.ID.EQ(mysql.Int64(id)),
	)
	if !q.IncludeArchived {
		condition = condition.AND(
			tJobGroups.State.NOT_EQ(mysql.Int32(int32(jobsgroups.GroupState_GROUP_STATE_ARCHIVED))),
		)
	}

	stmt := tJobGroups.
		SELECT(
			tJobGroups.ID.AS("group.id"),
			tJobGroups.Job.AS("group.job"),
			tJobGroups.Name.AS("group.name"),
			tJobGroups.Description.AS("group.description"),
			tJobGroups.ShortName.AS("group.short_name"),
			tJobGroups.LogoFileID.AS("group.logo_file_id"),
			tJobGroups.Color.AS("group.color"),
			tJobGroups.Type.AS("group.type"),
			tJobGroups.State.AS("group.state"),
			tJobGroups.MembershipMode.AS("group.membership_mode"),
			tJobGroups.SortRank.AS("group.sort_rank"),
			tJobGroups.MembersCount.AS("group.members_count"),
			tJobGroups.LeadersCount.AS("group.leaders_count"),
			tJobGroups.RulesCount.AS("group.rules_count"),
			tJobGroups.ExclusionsCount.AS("group.exclusions_count"),
			tJobGroups.CreatedByUserID.AS("group.created_by_user_id"),
			tJobGroups.UpdatedByUserID.AS("group.updated_by_user_id"),
			tJobGroups.CreatedAt.AS("group.created_at"),
			tJobGroups.UpdatedAt.AS("group.updated_at"),
			tJobGroups.DeletedAt.AS("group.deleted_at"),
			tLogoFile.ID.AS("logo_file.id"),
			tLogoFile.FilePath.AS("logo_file.file_path"),
			tLogoFile.ByteSize.AS("logo_file.byte_size"),
			tLogoFile.ContentType.AS("logo_file.content_type"),
			tLogoFile.CreatedAt.AS("logo_file.created_at"),
		).
		FROM(
			tJobGroups.
				LEFT_JOIN(tLogoFile,
					mysql.AND(
						tLogoFile.ID.EQ(tJobGroups.LogoFileID),
						tLogoFile.DeletedAt.IS_NULL(),
					),
				),
		).
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

	sortRank := strings.TrimSpace(group.GetSortRank())
	if sortRank == "" {
		var err error
		sortRank, err = s.nextGroupRank(ctx, db, group.GetJob(), 0)
		if err != nil {
			return 0, err
		}
	}

	tJobGroups := table.FivenetJobGroups
	insertStmt := tJobGroups.
		INSERT(
			tJobGroups.Job,
			tJobGroups.Name,
			tJobGroups.Description,
			tJobGroups.ShortName,
			tJobGroups.LogoFileID,
			tJobGroups.Color,
			tJobGroups.Type,
			tJobGroups.State,
			tJobGroups.MembershipMode,
			tJobGroups.SortRank,
			tJobGroups.CreatedByUserID,
			tJobGroups.UpdatedByUserID,
		).
		VALUES(
			group.GetJob(),
			group.GetName(),
			dbutils.StringEmpty(group.GetDescription()),
			dbutils.StringEmpty(group.GetShortName()),
			dbutils.Int64P(group.GetLogoFileId()),
			dbutils.StringEmpty(group.GetColor()),
			int32(group.GetType()),
			int32(group.GetState()),
			int32(group.GetMembershipMode()),
			sortRank,
			group.GetCreatedByUserId(),
			group.GetUpdatedByUserId(),
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
	tJobGroups := table.FivenetJobGroups
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
			tJobGroups.SortRank,
			tJobGroups.UpdatedByUserID,
			tJobGroups.UpdatedAt,
		).
		SET(
			tJobGroups.Name.SET(mysql.String(group.GetName())),
			tJobGroups.Description.SET(dbutils.StringPP(group.Description)),
			tJobGroups.ShortName.SET(dbutils.StringPP(group.ShortName)),
			tJobGroups.LogoFileID.SET(dbutils.Int64P(group.GetLogoFileId())),
			tJobGroups.Color.SET(dbutils.StringPP(group.Color)),
			tJobGroups.Type.SET(mysql.Int32(int32(group.GetType()))),
			tJobGroups.State.SET(mysql.Int32(int32(group.GetState()))),
			tJobGroups.MembershipMode.SET(mysql.Int32(int32(group.GetMembershipMode()))),
			tJobGroups.SortRank.SET(mysql.String(group.GetSortRank())),
			tJobGroups.UpdatedByUserID.SET(dbutils.Int32P(group.GetUpdatedByUserId())),
			tJobGroups.UpdatedAt.SET(mysql.CURRENT_TIMESTAMP()),
		).
		WHERE(mysql.AND(
			tJobGroups.ID.EQ(mysql.Int64(group.GetId())),
			tJobGroups.Job.EQ(mysql.String(group.GetJob())),
		)).
		LIMIT(1)

	_, err := updateStmt.ExecContext(ctx, db)
	return err
}

func (s *Store) ArchiveGroup(
	ctx context.Context,
	db qrm.DB,
	job string,
	id int64,
	updatedByUserID int32,
) error {
	tJobGroups := table.FivenetJobGroups
	updateStmt := tJobGroups.
		UPDATE(
			tJobGroups.State,
			tJobGroups.DeletedAt,
			tJobGroups.UpdatedByUserID,
			tJobGroups.UpdatedAt,
		).
		SET(
			tJobGroups.State.SET(mysql.Int32(int32(jobsgroups.GroupState_GROUP_STATE_ARCHIVED))),
			tJobGroups.DeletedAt.SET(mysql.CURRENT_TIMESTAMP()),
			tJobGroups.UpdatedByUserID.SET(dbutils.Int32P(updatedByUserID)),
			tJobGroups.UpdatedAt.SET(mysql.CURRENT_TIMESTAMP()),
		).
		WHERE(mysql.AND(
			tJobGroups.ID.EQ(mysql.Int64(id)),
			tJobGroups.Job.EQ(mysql.String(job)),
		)).
		LIMIT(1)

	_, err := updateStmt.ExecContext(ctx, db)
	return err
}

func (s *Store) RestoreGroup(
	ctx context.Context,
	db qrm.DB,
	job string,
	id int64,
	updatedByUserID int32,
) error {
	tJobGroups := table.FivenetJobGroups
	updateStmt := tJobGroups.
		UPDATE(
			tJobGroups.State,
			tJobGroups.DeletedAt,
			tJobGroups.UpdatedByUserID,
			tJobGroups.UpdatedAt,
		).
		SET(
			tJobGroups.State.SET(mysql.Int32(int32(jobsgroups.GroupState_GROUP_STATE_ACTIVE))),
			tJobGroups.DeletedAt.SET(mysql.TimestampExp(mysql.NULL)),
			tJobGroups.UpdatedByUserID.SET(dbutils.Int32P(updatedByUserID)),
			tJobGroups.UpdatedAt.SET(mysql.CURRENT_TIMESTAMP()),
		).
		WHERE(mysql.AND(
			tJobGroups.ID.EQ(mysql.Int64(id)),
			tJobGroups.Job.EQ(mysql.String(job)),
		)).
		LIMIT(1)

	_, err := updateStmt.ExecContext(ctx, db)
	return err
}
