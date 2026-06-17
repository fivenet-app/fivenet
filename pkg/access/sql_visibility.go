package access

import (
	"context"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type sqlVisibilityBackend struct {
	access *SubjectObjectAccess
}

func (b *sqlVisibilityBackend) VisibleIDsQuery(
	userInfo *userinfo.UserInfo,
	access int32,
	includeDeleted bool,
	targetIDs ...int64,
) VisibilityQuery {
	ids := make([]mysql.Expression, 0, len(targetIDs))
	for _, targetID := range targetIDs {
		ids = append(ids, mysql.Int64(targetID))
	}

	condition := b.access.targetTableColumns.ID.IN(ids...)
	return b.VisibleIDsByConditionQuery(userInfo, access, includeDeleted, condition)
}

func (b *sqlVisibilityBackend) VisibleIDsStatement(
	userInfo *userinfo.UserInfo,
	access int32,
	includeDeleted bool,
	targetIDs ...int64,
) mysql.Statement {
	query := b.VisibleIDsQuery(userInfo, access, includeDeleted, targetIDs...)
	if len(query.CTEs) == 0 {
		return query.Statement
	}
	return mysql.WITH(query.CTEs...)(query.Statement)
}

func (b *sqlVisibilityBackend) VisibleIDsByConditionStatement(
	userInfo *userinfo.UserInfo,
	access int32,
	includeDeleted bool,
	condition mysql.BoolExpression,
) mysql.Statement {
	query := b.VisibleIDsByConditionQuery(userInfo, access, includeDeleted, condition)
	if len(query.CTEs) == 0 {
		return query.Statement
	}
	return mysql.WITH(query.CTEs...)(query.Statement)
}

func (b *sqlVisibilityBackend) VisibleIDsByConditionQuery(
	userInfo *userinfo.UserInfo,
	access int32,
	includeDeleted bool,
	condition mysql.BoolExpression,
) VisibilityQuery {
	condition = b.access.normalizeTargetCondition(condition, includeDeleted)
	if userInfo != nil && userInfo.GetSuperuser() {
		return VisibilityQuery{
			Table:     b.access.visibleTargetIDsStatement(condition).AsTable("doc_ids"),
			Statement: b.access.visibleTargetIDsStatement(condition),
		}
	}

	implicitStmt := b.access.implicitVisibleIDsStatement(userInfo, condition)
	aclCtes, aclSelect := b.access.aclVisibleIDsComponents(userInfo, access, condition, false)
	if implicitStmt == nil {
		return VisibilityQuery{
			CTEs:      aclCtes,
			Table:     aclSelect.AsTable("doc_ids"),
			Statement: aclSelect,
		}
	}

	combined := implicitStmt.UNION(aclSelect)
	return VisibilityQuery{
		CTEs:      aclCtes,
		Table:     combined.AsTable("doc_ids"),
		Statement: combined,
	}
}

func (b *sqlVisibilityBackend) ACLVisibleIDsByConditionQuery(
	userInfo *userinfo.UserInfo,
	access int32,
	includeDeleted bool,
	condition mysql.BoolExpression,
) VisibilityQuery {
	condition = b.access.normalizeTargetCondition(condition, includeDeleted)
	if userInfo != nil && userInfo.GetSuperuser() {
		return VisibilityQuery{
			Table:     b.access.visibleTargetIDsStatement(condition).AsTable("doc_ids"),
			Statement: b.access.visibleTargetIDsStatement(condition),
		}
	}

	aclCtes, aclSelect := b.access.aclVisibleIDsComponents(userInfo, access, condition, false)
	return VisibilityQuery{
		CTEs:      aclCtes,
		Table:     aclSelect.AsTable("doc_ids"),
		Statement: aclSelect,
	}
}

func (b *sqlVisibilityBackend) ACLVisibleIDsByConditionStatement(
	userInfo *userinfo.UserInfo,
	access int32,
	includeDeleted bool,
	condition mysql.BoolExpression,
) mysql.Statement {
	query := b.ACLVisibleIDsByConditionQuery(userInfo, access, includeDeleted, condition)
	if len(query.CTEs) == 0 {
		return query.Statement
	}
	return mysql.WITH(query.CTEs...)(query.Statement)
}

func (b *sqlVisibilityBackend) CountVisibleByConditionStatement(
	userInfo *userinfo.UserInfo,
	access int32,
	includeDeleted bool,
	condition mysql.BoolExpression,
) mysql.Statement {
	condition = b.access.normalizeTargetCondition(condition, includeDeleted)
	if userInfo != nil && userInfo.GetSuperuser() {
		return b.access.countTargetIDsStatement(condition)
	}

	query := b.VisibleIDsByConditionQuery(userInfo, access, includeDeleted, condition)
	visibleID := mysql.IntegerColumn("id").From(query.Table)
	countSelect := mysql.
		SELECT(mysql.COUNT(visibleID).AS("exact_total")).FROM(query.Table)
	if len(query.CTEs) == 0 {
		return countSelect
	}
	return mysql.WITH(
		query.CTEs...,
	)(countSelect)
}

func (b *sqlVisibilityBackend) RefreshTargetVisibility(
	ctx context.Context,
	tx qrm.DB,
	targetID int64,
) error {
	return nil
}
