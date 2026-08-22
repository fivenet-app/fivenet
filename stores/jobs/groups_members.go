package jobsstore

import (
	"context"
	"database/sql"
	"errors"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	groupspolicy "github.com/fivenet-app/fivenet/v2026/stores/jobs/groupspolicy"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func groupManualMemberColumns(
	tMembers *table.FivenetJobGroupManualMembersTable,
) mysql.ProjectionList {
	return mysql.ProjectionList{
		tMembers.GroupID.AS("group_manual_member.group_id"),
		tMembers.UserID.AS("group_manual_member.user_id"),
		tMembers.Reason.AS("group_manual_member.reason"),
		tMembers.CreatedByUserID.AS("group_manual_member.created_by_user_id"),
		tMembers.CreatedAt.AS("group_manual_member.created_at"),
	}
}

func groupMemberExclusionColumns(
	tExclusions *table.FivenetJobGroupMemberExclusionsTable,
) mysql.ProjectionList {
	return mysql.ProjectionList{
		tExclusions.GroupID.AS("group_member_exclusion.group_id"),
		tExclusions.UserID.AS("group_member_exclusion.user_id"),
		tExclusions.ReasonType.AS("group_member_exclusion.reason_type"),
		tExclusions.Reason.AS("group_member_exclusion.reason"),
		tExclusions.CreatedByUserID.AS("group_member_exclusion.created_by_user_id"),
		tExclusions.CreatedAt.AS("group_member_exclusion.created_at"),
	}
}

func groupLeaderColumns(tLeaders *table.FivenetJobGroupLeadersTable) mysql.ProjectionList {
	return mysql.ProjectionList{
		tLeaders.GroupID.AS("group_leader.group_id"),
		tLeaders.UserID.AS("group_leader.user_id"),
		tLeaders.CreatedByUserID.AS("group_leader.created_by_user_id"),
		tLeaders.CreatedAt.AS("group_leader.created_at"),
	}
}

func groupMemberSearchCondition(search string, user *table.FivenetUserTable) mysql.BoolExpression {
	search = dbutils.PrepareForLikeSearch(search)
	if search == "" {
		return nil
	}

	like := mysql.String(search)
	return mysql.OR(
		mysql.CONCAT(user.Firstname, mysql.String(" "), user.Lastname).LIKE(like),
		user.Identifier.LIKE(like),
	)
}

func groupMemberCondition(
	groupID int64,
	groupIDColumn mysql.ColumnInteger,
	search string,
	user *table.FivenetUserTable,
) mysql.BoolExpression {
	condition := groupIDColumn.EQ(mysql.Int64(groupID))
	if searchCondition := groupMemberSearchCondition(search, user); searchCondition != nil {
		condition = condition.AND(searchCondition)
	}

	return condition
}

func activeGroupMemberCondition(
	groupID int64,
	groupIDColumn mysql.ColumnInteger,
	search string,
	user *table.FivenetUserTable,
	userJobs *table.FivenetUserJobsTable,
	groupJob mysql.StringExpression,
) mysql.BoolExpression {
	return groupMemberCondition(groupID, groupIDColumn, search, user).
		AND(userJobs.Job.EQ(groupJob))
}

func (s *Store) UserInJob(ctx context.Context, db qrm.DB, job string, userID int32) (bool, error) {
	tUserJobs := table.FivenetUserJobs
	stmt := tUserJobs.
		SELECT(mysql.Int(1).AS("found")).
		FROM(tUserJobs).
		WHERE(mysql.AND(
			tUserJobs.UserID.EQ(mysql.Int32(userID)),
			tUserJobs.Job.EQ(mysql.String(job)),
		)).
		LIMIT(1)

	dest := struct{ Found int }{}
	if err := stmt.QueryContext(ctx, db, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return dest.Found == 1, nil
}

func (s *Store) RecountGroupStats(ctx context.Context, db qrm.DB, groupID int64) error {
	groupInfo := struct {
		Type           jobsgroups.GroupType
		Job            string
		MembershipMode jobsgroups.GroupMembershipMode
	}{}
	if err := tJobGroups.
		SELECT(
			tJobGroups.Type.AS("type"),
			tJobGroups.Job.AS("job"),
			tJobGroups.MembershipMode.AS("membership_mode"),
		).
		FROM(tJobGroups).
		WHERE(tJobGroups.ID.EQ(mysql.Int64(groupID))).
		LIMIT(1).
		QueryContext(ctx, db, &groupInfo); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
		return sql.ErrNoRows
	}
	group := &jobsgroups.Group{
		Type:           groupInfo.Type,
		Id:             groupID,
		Job:            groupInfo.Job,
		MembershipMode: groupInfo.MembershipMode,
	}

	var manualMembers []*jobsgroups.GroupManualMember
	var ruleMatches []*GroupRuleMemberMatch
	var exclusions []*jobsgroups.GroupMemberExclusion
	var err error
	if groupspolicy.AllowsManualMembers(group.GetType()) {
		manualMembers, err = s.ListGroupManualMembers(ctx, db, GroupItemsQuery{GroupID: groupID})
		if err != nil {
			return err
		}
	}
	if groupspolicy.AllowsRules(group.GetType()) ||
		groupspolicy.RequiresManualMembersMatchRules(group.GetType(), group.GetMembershipMode()) {
		ruleMatches, err = s.ListGroupRuleMemberMatches(ctx, db, group, "")
		if err != nil {
			return err
		}
	}
	if groupspolicy.AllowsExclusions(group.GetType()) {
		exclusions, err = s.ListGroupMemberExclusions(ctx, db, GroupItemsQuery{GroupID: groupID})
		if err != nil {
			return err
		}
	}
	leadersCount, err := s.CountGroupLeaders(ctx, db, GroupItemsQuery{GroupID: groupID})
	if err != nil {
		return err
	}
	var rulesCount int64
	if groupspolicy.AllowsRules(group.GetType()) {
		rulesCount, err = s.CountGroupRules(ctx, db, GroupItemsQuery{GroupID: groupID})
		if err != nil {
			return err
		}
	}

	excluded := map[int32]struct{}{}
	for _, exclusion := range exclusions {
		excluded[exclusion.GetUserId()] = struct{}{}
	}

	members := map[int32]struct{}{}
	ruleMemberIDs := map[int32]struct{}{}
	if groupspolicy.AllowsRules(group.GetType()) {
		for _, match := range ruleMatches {
			ruleMemberIDs[match.UserID] = struct{}{}
			if _, ok := excluded[match.UserID]; !ok {
				members[match.UserID] = struct{}{}
			}
		}
	}
	if groupspolicy.AllowsManualMembers(group.GetType()) {
		for _, member := range manualMembers {
			if _, ok := excluded[member.GetUserId()]; ok {
				continue
			}
			if groupspolicy.RequiresManualMembersMatchRules(
				group.GetType(),
				group.GetMembershipMode(),
			) {
				if _, ok := ruleMemberIDs[member.GetUserId()]; !ok {
					continue
				}
			}
			members[member.GetUserId()] = struct{}{}
		}
	}

	_, err = tJobGroups.
		UPDATE().
		SET(
			tJobGroups.MembersCount.SET(mysql.Int64(int64(len(members)))),
			tJobGroups.LeadersCount.SET(mysql.Int64(leadersCount)),
			tJobGroups.RulesCount.SET(mysql.Int64(rulesCount)),
			tJobGroups.ExclusionsCount.SET(mysql.Int64(int64(len(exclusions)))),
		).
		WHERE(tJobGroups.ID.EQ(mysql.Int64(groupID))).
		ExecContext(ctx, db)
	return err
}

func (s *Store) CountGroupLeaders(
	ctx context.Context,
	db qrm.DB,
	q GroupItemsQuery,
) (int64, error) {
	tLeaders := table.FivenetJobGroupLeaders.AS("gl")
	tUser := table.FivenetUser.AS("u")
	tUserJobs := table.FivenetUserJobs.AS("uj")
	tJobGroups := table.FivenetJobGroups.AS("g")
	var count database.DataCount
	if err := tLeaders.
		SELECT(mysql.COUNT(tLeaders.UserID).AS("data_count.total")).
		FROM(tLeaders.
			INNER_JOIN(tUser, tUser.ID.EQ(tLeaders.UserID)).
			INNER_JOIN(tJobGroups, tJobGroups.ID.EQ(tLeaders.GroupID)).
			INNER_JOIN(tUserJobs, tUserJobs.UserID.EQ(tUser.ID)),
		).
		WHERE(activeGroupMemberCondition(q.GroupID, tLeaders.GroupID, q.Search, tUser, tUserJobs, tJobGroups.Job)).
		QueryContext(ctx, db, &count); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return count.Total, nil
}

func (s *Store) CountGroupManualMembers(
	ctx context.Context,
	db qrm.DB,
	q GroupItemsQuery,
) (int64, error) {
	tMembers := table.FivenetJobGroupManualMembers.AS("mm")
	tUser := table.FivenetUser.AS("u")
	tUserJobs := table.FivenetUserJobs.AS("uj")
	tJobGroups := table.FivenetJobGroups.AS("g")
	var count database.DataCount
	if err := tMembers.
		SELECT(mysql.COUNT(tMembers.UserID).AS("data_count.total")).
		FROM(tMembers.
			INNER_JOIN(tUser, tUser.ID.EQ(tMembers.UserID)).
			INNER_JOIN(tJobGroups, tJobGroups.ID.EQ(tMembers.GroupID)).
			INNER_JOIN(tUserJobs, tUserJobs.UserID.EQ(tUser.ID)),
		).
		WHERE(activeGroupMemberCondition(q.GroupID, tMembers.GroupID, q.Search, tUser, tUserJobs, tJobGroups.Job)).
		QueryContext(ctx, db, &count); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return count.Total, nil
}

func (s *Store) CountGroupMemberExclusions(
	ctx context.Context,
	db qrm.DB,
	q GroupItemsQuery,
) (int64, error) {
	tExclusions := table.FivenetJobGroupMemberExclusions.AS("me")
	tUser := table.FivenetUser.AS("u")
	var count database.DataCount
	if err := tExclusions.
		SELECT(mysql.COUNT(tExclusions.UserID).AS("data_count.total")).
		FROM(tExclusions.INNER_JOIN(tUser, tUser.ID.EQ(tExclusions.UserID))).
		WHERE(groupMemberCondition(q.GroupID, tExclusions.GroupID, q.Search, tUser)).
		QueryContext(ctx, db, &count); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return count.Total, nil
}

func (s *Store) CountGroupRules(ctx context.Context, db qrm.DB, q GroupItemsQuery) (int64, error) {
	tRules := table.FivenetJobGroupRules
	var count database.DataCount
	if err := tRules.
		SELECT(mysql.COUNT(tRules.ID).AS("data_count.total")).
		FROM(tRules).
		WHERE(tRules.GroupID.EQ(mysql.Int64(q.GroupID))).
		QueryContext(ctx, db, &count); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return count.Total, nil
}

func (s *Store) ListGroupManualMembers(
	ctx context.Context,
	db qrm.DB,
	q GroupItemsQuery,
) ([]*jobsgroups.GroupManualMember, error) {
	tMembers := table.FivenetJobGroupManualMembers.AS("mm")
	tUser := table.FivenetUser.AS("u")
	tUserJobs := table.FivenetUserJobs.AS("uj")
	tJobGroups := table.FivenetJobGroups.AS("g")
	columns := groupManualMemberColumns(tMembers)
	stmt := tMembers.
		SELECT(columns[0], columns[1:]...).
		FROM(tMembers.
			INNER_JOIN(tUser, tUser.ID.EQ(tMembers.UserID)).
			INNER_JOIN(tJobGroups, tJobGroups.ID.EQ(tMembers.GroupID)).
			INNER_JOIN(tUserJobs, tUserJobs.UserID.EQ(tUser.ID)),
		).
		WHERE(activeGroupMemberCondition(q.GroupID, tMembers.GroupID, q.Search, tUser, tUserJobs, tJobGroups.Job)).
		ORDER_BY(tMembers.CreatedAt.ASC(), tMembers.UserID.ASC())
	if q.Offset > 0 {
		stmt = stmt.OFFSET(q.Offset)
	}
	if q.Limit > 0 {
		stmt = stmt.LIMIT(q.Limit)
	}

	members := []*jobsgroups.GroupManualMember{}
	if err := stmt.QueryContext(ctx, db, &members); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return []*jobsgroups.GroupManualMember{}, nil
		}
		return nil, err
	}

	return members, nil
}

func (s *Store) ListGroupMemberExclusions(
	ctx context.Context,
	db qrm.DB,
	q GroupItemsQuery,
) ([]*jobsgroups.GroupMemberExclusion, error) {
	tExclusions := table.FivenetJobGroupMemberExclusions.AS("me")
	tUser := table.FivenetUser.AS("u")
	columns := groupMemberExclusionColumns(tExclusions)
	stmt := tExclusions.
		SELECT(columns[0], columns[1:]...).
		FROM(tExclusions.INNER_JOIN(tUser, tUser.ID.EQ(tExclusions.UserID))).
		WHERE(groupMemberCondition(q.GroupID, tExclusions.GroupID, q.Search, tUser)).
		ORDER_BY(tExclusions.CreatedAt.ASC(), tExclusions.UserID.ASC())
	if q.Offset > 0 {
		stmt = stmt.OFFSET(q.Offset)
	}
	if q.Limit > 0 {
		stmt = stmt.LIMIT(q.Limit)
	}

	exclusions := []*jobsgroups.GroupMemberExclusion{}
	if err := stmt.QueryContext(ctx, db, &exclusions); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return []*jobsgroups.GroupMemberExclusion{}, nil
		}
		return nil, err
	}

	return exclusions, nil
}

func (s *Store) ListGroupLeaders(
	ctx context.Context,
	db qrm.DB,
	q GroupItemsQuery,
) ([]*jobsgroups.GroupLeader, error) {
	tLeaders := table.FivenetJobGroupLeaders.AS("gl")
	tUser := table.FivenetUser.AS("u")
	tUserJobs := table.FivenetUserJobs.AS("uj")
	tJobGroups := table.FivenetJobGroups.AS("g")
	columns := groupLeaderColumns(tLeaders)
	stmt := tLeaders.
		SELECT(columns[0], columns[1:]...).
		FROM(tLeaders.
			INNER_JOIN(tUser, tUser.ID.EQ(tLeaders.UserID)).
			INNER_JOIN(tJobGroups, tJobGroups.ID.EQ(tLeaders.GroupID)).
			INNER_JOIN(tUserJobs, tUserJobs.UserID.EQ(tUser.ID)),
		).
		WHERE(activeGroupMemberCondition(q.GroupID, tLeaders.GroupID, q.Search, tUser, tUserJobs, tJobGroups.Job)).
		ORDER_BY(tLeaders.CreatedAt.ASC(), tLeaders.UserID.ASC())
	if q.Offset > 0 {
		stmt = stmt.OFFSET(q.Offset)
	}
	if q.Limit > 0 {
		stmt = stmt.LIMIT(q.Limit)
	}

	leaders := []*jobsgroups.GroupLeader{}
	if err := stmt.QueryContext(ctx, db, &leaders); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return []*jobsgroups.GroupLeader{}, nil
		}
		return nil, err
	}

	return leaders, nil
}

func (s *Store) getGroupManualMember(
	ctx context.Context,
	db qrm.DB,
	groupID int64,
	userID int32,
) (*jobsgroups.GroupManualMember, error) {
	tMembers := table.FivenetJobGroupManualMembers
	columns := groupManualMemberColumns(tMembers)
	member := &jobsgroups.GroupManualMember{}
	if err := tMembers.
		SELECT(columns[0], columns[1:]...).
		FROM(tMembers).
		WHERE(mysql.AND(
			tMembers.GroupID.EQ(mysql.Int64(groupID)),
			tMembers.UserID.EQ(mysql.Int32(userID)),
		)).
		LIMIT(1).
		QueryContext(ctx, db, member); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	if member.GetGroupId() == 0 {
		return nil, sql.ErrNoRows
	}

	return member, nil
}

func (s *Store) getGroupMemberExclusion(
	ctx context.Context,
	db qrm.DB,
	groupID int64,
	userID int32,
) (*jobsgroups.GroupMemberExclusion, error) {
	tExclusions := table.FivenetJobGroupMemberExclusions
	columns := groupMemberExclusionColumns(tExclusions)
	exclusion := &jobsgroups.GroupMemberExclusion{}
	if err := tExclusions.
		SELECT(columns[0], columns[1:]...).
		FROM(tExclusions).
		WHERE(mysql.AND(
			tExclusions.GroupID.EQ(mysql.Int64(groupID)),
			tExclusions.UserID.EQ(mysql.Int32(userID)),
		)).
		LIMIT(1).
		QueryContext(ctx, db, exclusion); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	if exclusion.GetGroupId() == 0 {
		return nil, sql.ErrNoRows
	}

	return exclusion, nil
}

func (s *Store) getGroupLeader(
	ctx context.Context,
	db qrm.DB,
	groupID int64,
	userID int32,
) (*jobsgroups.GroupLeader, error) {
	tLeaders := table.FivenetJobGroupLeaders
	columns := groupLeaderColumns(tLeaders)
	leader := &jobsgroups.GroupLeader{}
	if err := tLeaders.
		SELECT(columns[0], columns[1:]...).
		FROM(tLeaders).
		WHERE(mysql.AND(
			tLeaders.GroupID.EQ(mysql.Int64(groupID)),
			tLeaders.UserID.EQ(mysql.Int32(userID)),
		)).
		LIMIT(1).
		QueryContext(ctx, db, leader); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	if leader.GetGroupId() == 0 {
		return nil, sql.ErrNoRows
	}

	return leader, nil
}

func (s *Store) AddGroupManualMember(
	ctx context.Context,
	db qrm.DB,
	groupID int64,
	userID int32,
	createdByUserID int32,
	reason *string,
) (*jobsgroups.GroupManualMember, bool, error) {
	created := false
	if _, err := s.getGroupManualMember(ctx, db, groupID, userID); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, false, err
		}
		created = true
	}

	var reasonValue any
	if reason != nil {
		reasonValue = *reason
	}

	tMembers := table.FivenetJobGroupManualMembers
	_, err := tMembers.
		INSERT(
			tMembers.GroupID,
			tMembers.UserID,
			tMembers.Reason,
			tMembers.CreatedByUserID,
			tMembers.CreatedAt,
		).
		VALUES(groupID, userID, reasonValue, createdByUserID, mysql.CURRENT_TIMESTAMP()).
		ON_DUPLICATE_KEY_UPDATE(
			tMembers.Reason.SET(mysql.RawString("VALUES(`reason`)")),
		).
		ExecContext(ctx, db)
	if err != nil {
		return nil, false, err
	}

	member, err := s.getGroupManualMember(ctx, db, groupID, userID)
	return member, created, err
}

func (s *Store) RemoveGroupManualMember(
	ctx context.Context,
	db qrm.DB,
	groupID int64,
	userID int32,
) error {
	tMembers := table.FivenetJobGroupManualMembers
	result, err := tMembers.
		DELETE().
		WHERE(mysql.AND(
			tMembers.GroupID.EQ(mysql.Int64(groupID)),
			tMembers.UserID.EQ(mysql.Int32(userID)),
		)).
		ExecContext(ctx, db)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *Store) AddGroupMemberExclusion(
	ctx context.Context,
	db qrm.DB,
	groupID int64,
	userID int32,
	reasonType jobsgroups.GroupExclusionReason,
	createdByUserID int32,
	reason *string,
) (*jobsgroups.GroupMemberExclusion, bool, error) {
	created := false
	if _, err := s.getGroupMemberExclusion(ctx, db, groupID, userID); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, false, err
		}
		created = true
	}

	if reasonType == jobsgroups.GroupExclusionReason_GROUP_EXCLUSION_REASON_UNSPECIFIED {
		reasonType = jobsgroups.GroupExclusionReason_GROUP_EXCLUSION_REASON_MANUAL
	}
	var reasonValue any
	if reason != nil {
		reasonValue = *reason
	}

	tExclusions := table.FivenetJobGroupMemberExclusions
	_, err := tExclusions.
		INSERT(
			tExclusions.GroupID,
			tExclusions.UserID,
			tExclusions.ReasonType,
			tExclusions.Reason,
			tExclusions.CreatedByUserID,
			tExclusions.CreatedAt,
		).
		VALUES(groupID, userID, int32(reasonType), reasonValue, createdByUserID, mysql.CURRENT_TIMESTAMP()).
		ON_DUPLICATE_KEY_UPDATE(
			tExclusions.ReasonType.SET(mysql.RawInt("VALUES(`reason_type`)")),
			tExclusions.Reason.SET(mysql.RawString("VALUES(`reason`)")),
		).
		ExecContext(ctx, db)
	if err != nil {
		return nil, false, err
	}

	exclusion, err := s.getGroupMemberExclusion(ctx, db, groupID, userID)
	return exclusion, created, err
}

func (s *Store) RemoveGroupMemberExclusion(
	ctx context.Context,
	db qrm.DB,
	groupID int64,
	userID int32,
) error {
	tExclusions := table.FivenetJobGroupMemberExclusions
	result, err := tExclusions.
		DELETE().
		WHERE(mysql.AND(
			tExclusions.GroupID.EQ(mysql.Int64(groupID)),
			tExclusions.UserID.EQ(mysql.Int32(userID)),
		)).
		ExecContext(ctx, db)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *Store) AddGroupLeader(
	ctx context.Context,
	db qrm.DB,
	groupID int64,
	userID int32,
	createdByUserID int32,
) (*jobsgroups.GroupLeader, bool, error) {
	created := false
	if _, err := s.getGroupLeader(ctx, db, groupID, userID); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, false, err
		}
		created = true
	}

	tLeaders := table.FivenetJobGroupLeaders
	_, err := tLeaders.
		INSERT(
			tLeaders.GroupID,
			tLeaders.UserID,
			tLeaders.CreatedByUserID,
			tLeaders.CreatedAt,
		).
		VALUES(groupID, userID, createdByUserID, mysql.CURRENT_TIMESTAMP()).
		ON_DUPLICATE_KEY_UPDATE(
			tLeaders.UserID.SET(mysql.RawInt("VALUES(`user_id`)")),
		).
		ExecContext(ctx, db)
	if err != nil {
		return nil, false, err
	}

	leader, err := s.getGroupLeader(ctx, db, groupID, userID)
	return leader, created, err
}

func (s *Store) RemoveGroupLeader(
	ctx context.Context,
	db qrm.DB,
	groupID int64,
	userID int32,
) error {
	tLeaders := table.FivenetJobGroupLeaders
	result, err := tLeaders.
		DELETE().
		WHERE(mysql.AND(
			tLeaders.GroupID.EQ(mysql.Int64(groupID)),
			tLeaders.UserID.EQ(mysql.Int32(userID)),
		)).
		ExecContext(ctx, db)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
