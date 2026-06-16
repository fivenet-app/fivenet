package access

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type calculatedVisibilityBackend struct {
	access *SubjectObjectAccess
}

type visibilityRowsStatement interface {
	mysql.Statement
	AsTable(alias string) mysql.SelectTable
}

func (b *calculatedVisibilityBackend) VisibleIDsQuery(
	userInfo *userinfo.UserInfo,
	access int32,
	targetIDs ...int64,
) VisibilityQuery {
	ids := make([]mysql.Expression, 0, len(targetIDs))
	for _, targetID := range targetIDs {
		ids = append(ids, mysql.Int64(targetID))
	}

	condition := b.access.targetTableColumns.ID.IN(ids...)
	return b.VisibleIDsByConditionQuery(userInfo, access, condition)
}

func (b *calculatedVisibilityBackend) VisibleIDsStatement(
	userInfo *userinfo.UserInfo,
	access int32,
	targetIDs ...int64,
) mysql.Statement {
	query := b.VisibleIDsQuery(userInfo, access, targetIDs...)
	return b.idOnlyStatement(query)
}

func (b *calculatedVisibilityBackend) VisibleIDsByConditionStatement(
	userInfo *userinfo.UserInfo,
	access int32,
	condition mysql.BoolExpression,
) mysql.Statement {
	query := b.VisibleIDsByConditionQuery(userInfo, access, condition)
	return b.idOnlyStatement(query)
}

func (b *calculatedVisibilityBackend) ACLVisibleIDsByConditionStatement(
	userInfo *userinfo.UserInfo,
	access int32,
	condition mysql.BoolExpression,
) mysql.Statement {
	query := b.ACLVisibleIDsByConditionQuery(userInfo, access, condition)
	return b.idOnlyStatement(query)
}

func (b *calculatedVisibilityBackend) CountVisibleByConditionStatement(
	userInfo *userinfo.UserInfo,
	access int32,
	condition mysql.BoolExpression,
) mysql.Statement {
	if userInfo == nil {
		userInfo = &userinfo.UserInfo{}
	}
	condition = b.access.normalizeTargetCondition(condition)
	if userInfo.GetSuperuser() {
		return b.access.countTargetIDsStatement(condition)
	}

	query := b.VisibleIDsByConditionQuery(userInfo, access, condition)
	visibleID := mysql.IntegerColumn("id").From(query.Table)
	countSelect := mysql.SELECT(mysql.COUNT(visibleID).AS("exact_total")).FROM(query.Table)
	if len(query.CTEs) == 0 {
		return countSelect
	}
	return mysql.WITH(query.CTEs...)(countSelect)
}

func (b *calculatedVisibilityBackend) RefreshTargetVisibility(
	ctx context.Context,
	tx qrm.DB,
	targetID int64,
) error {
	return b.RefreshTargetVisibilityWithCreator(ctx, tx, targetID, 0, "")
}

func (b *calculatedVisibilityBackend) RefreshTargetVisibilityWithCreator(
	ctx context.Context,
	tx qrm.DB,
	targetID int64,
	creatorID int32,
	creatorJob string,
) error {
	if b.access.calculatedVisibilityPublicTable != nil {
		if _, err := b.access.calculatedVisibilityPublicTable.
			DELETE().
			WHERE(b.access.calculatedVisibilityPublicTargetID.EQ(mysql.Int64(targetID))).
			ExecContext(ctx, tx); err != nil {
			return err
		}
	}
	if b.access.calculatedVisibilityCreatorTable != nil {
		if _, err := b.access.calculatedVisibilityCreatorTable.
			DELETE().
			WHERE(b.access.calculatedVisibilityCreatorTargetID.EQ(mysql.Int64(targetID))).
			ExecContext(ctx, tx); err != nil {
			return err
		}
	}
	if _, err := b.access.calculatedVisibilitySubjectTable.
		DELETE().
		WHERE(b.access.calculatedVisibilitySubjectTargetID.EQ(mysql.Int64(targetID))).
		ExecContext(ctx, tx); err != nil {
		return err
	}

	row, err := b.loadTargetRow(ctx, tx, targetID)
	if err != nil || row == nil {
		return err
	}

	if row.Public && b.access.calculatedVisibilityPublicTable != nil {
		if err := b.insertPublicRow(ctx, tx, row.ID); err != nil {
			return err
		}
	}
	if creatorID != 0 || creatorJob != "" {
		row.CreatorID = &creatorID
		if creatorJob != "" {
			row.CreatorJob = &creatorJob
		}
	}
	if row.CreatorID == nil && row.CreatorJob == nil {
		creatorRow, err := b.loadTargetCreatorRow(ctx, tx, targetID)
		if err != nil {
			return err
		}
		if creatorRow != nil {
			row.CreatorID = creatorRow.CreatorID
			row.CreatorJob = creatorRow.CreatorJob
		}
	}
	if b.access.calculatedVisibilityCreatorTable != nil &&
		(row.CreatorID != nil || row.CreatorJob != nil) {
		creatorID := int64(0)
		if row.CreatorID != nil {
			creatorID = int64(*row.CreatorID)
		}
		creatorJob := ""
		if row.CreatorJob != nil {
			creatorJob = *row.CreatorJob
		}
		if err := b.insertCreatorRow(ctx, tx, row.ID, creatorID, creatorJob); err != nil {
			return err
		}
	}

	aclRows, err := b.loadACLRows(ctx, tx, targetID)
	if err != nil {
		return err
	}
	for _, aclRow := range aclRows {
		if err := b.insertSubjectRow(
			ctx,
			tx,
			row.ID,
			aclRow.SubjectID,
			aclRow.Access,
			aclRow.Effect,
		); err != nil {
			return err
		}
	}

	return nil
}

func (b *calculatedVisibilityBackend) VisibleIDsByConditionQuery(
	userInfo *userinfo.UserInfo,
	access int32,
	condition mysql.BoolExpression,
) VisibilityQuery {
	if userInfo == nil {
		userInfo = &userinfo.UserInfo{}
	}
	condition = b.access.normalizeTargetCondition(condition)
	if userInfo.GetSuperuser() {
		rows := b.visibleTargetRowsSelect(condition)
		return VisibilityQuery{Table: rows.AsTable("doc_ids"), Statement: rows}
	}

	var publicRows mysql.SelectStatement
	if b.access.visibility.HasRuleKind(VisibilityRulePublic) {
		publicRows = b.publicVisibleRowsSelect(condition)
	}
	var creatorRows mysql.SelectStatement
	if b.access.visibility.HasRuleKind(VisibilityRuleCreator) {
		creatorRows = b.creatorVisibleRowsSelect(userInfo, condition)
	}
	subjectCtes, subjectRows := b.subjectVisibleRowsComponents(userInfo, access, condition)

	combined := b.unionVisibleRows(publicRows, creatorRows, subjectRows)
	return VisibilityQuery{
		CTEs:      subjectCtes,
		Table:     combined.AsTable("doc_ids"),
		Statement: combined,
	}
}

func (b *calculatedVisibilityBackend) ACLVisibleIDsByConditionQuery(
	userInfo *userinfo.UserInfo,
	access int32,
	condition mysql.BoolExpression,
) VisibilityQuery {
	if userInfo == nil {
		userInfo = &userinfo.UserInfo{}
	}
	condition = b.access.normalizeTargetCondition(condition)
	if userInfo.GetSuperuser() {
		rows := b.visibleTargetRowsSelect(condition)
		return VisibilityQuery{Table: rows.AsTable("doc_ids"), Statement: rows}
	}

	var creatorRows mysql.SelectStatement
	if b.access.visibility.HasRuleKind(VisibilityRuleCreator) {
		creatorRows = b.creatorVisibleRowsSelect(userInfo, condition)
	}

	subjectCtes, subjectRows := b.subjectVisibleRowsComponents(userInfo, access, condition)
	if subjectRows == nil {
		rows := b.unionVisibleRows(func() mysql.SelectStatement {
			if b.access.visibility.HasRuleKind(VisibilityRulePublic) {
				return b.publicVisibleRowsSelect(condition)
			}
			return nil
		}(), creatorRows)
		return VisibilityQuery{CTEs: subjectCtes, Table: rows.AsTable("doc_ids"), Statement: rows}
	}

	combined := b.unionVisibleRows(subjectRows, creatorRows, func() mysql.SelectStatement {
		if b.access.visibility.HasRuleKind(VisibilityRulePublic) {
			return b.publicVisibleRowsSelect(condition)
		}
		return nil
	}())
	return VisibilityQuery{
		CTEs:      subjectCtes,
		Table:     combined.AsTable("doc_ids"),
		Statement: combined,
	}
}

func (b *calculatedVisibilityBackend) idOnlyStatement(query VisibilityQuery) mysql.Statement {
	idSelect := mysql.SELECT(mysql.IntegerColumn("id").From(query.Table).AS("id")).
		FROM(query.Table).
		DISTINCT()
	if len(query.CTEs) == 0 {
		return idSelect
	}
	return mysql.WITH(query.CTEs...)(idSelect)
}

func (b *calculatedVisibilityBackend) unionVisibleRows(
	rows ...mysql.SelectStatement,
) visibilityRowsStatement {
	filtered := make([]mysql.SelectStatement, 0, len(rows))
	for _, row := range rows {
		if row != nil {
			filtered = append(filtered, row)
		}
	}
	switch len(filtered) {
	case 0:
		return b.emptyVisibleRowsSelect()
	case 1:
		return filtered[0]
	default:
		result := mysql.UNION(filtered[0], filtered[1])
		for _, row := range filtered[2:] {
			result = mysql.UNION(result, row)
		}
		return result
	}
}

func (b *calculatedVisibilityBackend) emptyVisibleRowsSelect() mysql.SelectStatement {
	return mysql.SELECT(
		b.access.targetTableColumns.ID.AS("id"),
		b.access.targetTableColumns.CreatedAt.AS("created_at"),
		b.access.targetTableColumns.UpdatedAt.AS("updated_at"),
	).
		FROM(b.access.targetTable).
		WHERE(mysql.Bool(false)).
		DISTINCT()
}

func (b *calculatedVisibilityBackend) visibleTargetRowsSelect(
	condition mysql.BoolExpression,
) mysql.SelectStatement {
	return mysql.SELECT(
		b.access.targetTableColumns.ID.AS("id"),
		b.access.targetTableColumns.CreatedAt.AS("created_at"),
		b.access.targetTableColumns.UpdatedAt.AS("updated_at"),
	).
		FROM(b.access.targetTable).
		WHERE(condition).
		DISTINCT()
}

func (b *calculatedVisibilityBackend) publicVisibleRowsSelect(
	condition mysql.BoolExpression,
) mysql.SelectStatement {
	return mysql.SELECT(
		b.access.calculatedVisibilityPublicTargetID.AS("id"),
		b.access.targetTableColumns.CreatedAt.AS("created_at"),
		b.access.targetTableColumns.UpdatedAt.AS("updated_at"),
	).
		FROM(b.access.calculatedVisibilityPublicTable.
			INNER_JOIN(b.access.targetTable,
				b.access.calculatedVisibilityPublicTargetID.EQ(b.access.targetTableColumns.ID),
			),
		).
		WHERE(condition).
		DISTINCT()
}

func (b *calculatedVisibilityBackend) creatorVisibleRowsSelect(
	userInfo *userinfo.UserInfo,
	condition mysql.BoolExpression,
) mysql.SelectStatement {
	if userInfo == nil {
		return b.emptyVisibleRowsSelect()
	}

	creatorMatches := []mysql.BoolExpression{}
	if b.access.targetTableColumns.CreatorID != nil && userInfo.GetUserId() > 0 {
		creatorMatches = append(
			creatorMatches,
			b.access.calculatedVisibilityCreatorCreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
		)
	}
	if b.access.targetTableColumns.CreatorJob != nil && userInfo.GetJob() != "" {
		creatorMatches = append(
			creatorMatches,
			b.access.calculatedVisibilityCreatorCreatorJob.EQ(mysql.String(userInfo.GetJob())),
		)
	}
	if len(creatorMatches) == 0 {
		return b.emptyVisibleRowsSelect()
	}

	return mysql.SELECT(
		b.access.calculatedVisibilityCreatorTargetID.AS("id"),
		b.access.targetTableColumns.CreatedAt.AS("created_at"),
		b.access.targetTableColumns.UpdatedAt.AS("updated_at"),
	).
		FROM(b.access.calculatedVisibilityCreatorTable.
			INNER_JOIN(b.access.targetTable,
				b.access.calculatedVisibilityCreatorTargetID.EQ(b.access.targetTableColumns.ID),
			),
		).
		WHERE(mysql.AND(
			condition,
			mysql.AND(creatorMatches...),
		)).
		DISTINCT()
}

func (b *calculatedVisibilityBackend) subjectVisibleRowsComponents(
	userInfo *userinfo.UserInfo,
	access int32,
	condition mysql.BoolExpression,
) ([]mysql.CommonTableExpression, mysql.SelectStatement) {
	if userInfo == nil || userInfo.GetUserId() <= 0 {
		return nil, nil
	}

	actorSubjects := mysql.CTE("actor_subjects")
	tSubjectUsers := table.FivenetACLSubjectUsers.AS("asu")
	tSubjectQualis := table.FivenetACLSubjectQualifications.AS("asq")
	tQualiResults := table.FivenetQualificationsResults.AS("qr")
	tSubjectJobGrade := table.FivenetACLSubjectJobGradeScopes.AS("asjg")
	tUserJobs := table.FivenetUserJobs.AS("uj")

	actorSubjectsSelect := mysql.SELECT(tSubjectUsers.SubjectID.AS("subject_id")).
		FROM(tSubjectUsers).
		WHERE(tSubjectUsers.UserID.EQ(mysql.Int32(userInfo.GetUserId()))).
		UNION(
			mysql.SELECT(tSubjectQualis.SubjectID.AS("subject_id")).FROM(tSubjectQualis.
				INNER_JOIN(tQualiResults, mysql.AND(
					tQualiResults.QualificationID.EQ(tSubjectQualis.QualificationID),
					tQualiResults.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
					tQualiResults.DeletedAt.IS_NULL(),
					tQualiResults.Status.EQ(
						mysql.Int32(int32(qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL)),
					),
				),
				),
			),
		).
		UNION(
			mysql.SELECT(tSubjectJobGrade.SubjectID.AS("subject_id")).FROM(tSubjectJobGrade.
				INNER_JOIN(tUserJobs, mysql.AND(
					tUserJobs.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
					tUserJobs.Job.EQ(tSubjectJobGrade.Job),
					tUserJobs.Grade.GT_EQ(tSubjectJobGrade.MinimumGrade),
				)),
			),
		).
		UNION(
			mysql.SELECT(tSubjectJobGrade.SubjectID.AS("subject_id")).
				FROM(tSubjectJobGrade).
				WHERE(mysql.AND(
					tSubjectJobGrade.Job.EQ(mysql.String(userInfo.GetJob())),
					tSubjectJobGrade.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),
				)),
		)

	actorSubjectsClause := actorSubjects.AS(actorSubjectsSelect)
	actorSubjectID := mysql.IntegerColumn("subject_id").From(actorSubjectsClause)

	subjectRows := mysql.SELECT(
		b.access.calculatedVisibilitySubjectTargetID.AS("id"),
		b.access.targetTableColumns.CreatedAt.AS("created_at"),
		b.access.targetTableColumns.UpdatedAt.AS("updated_at"),
	).
		FROM(b.access.calculatedVisibilitySubjectTable.
			INNER_JOIN(b.access.targetTable,
				b.access.calculatedVisibilitySubjectTargetID.EQ(b.access.targetTableColumns.ID),
			).
			INNER_JOIN(actorSubjectsClause,
				actorSubjectID.EQ(b.access.calculatedVisibilitySubjectSubjectID),
			),
		).
		WHERE(mysql.AND(
			condition,
			b.access.calculatedVisibilitySubjectEffect.IS_TRUE(),
			b.access.calculatedVisibilitySubjectAccess.GT_EQ(mysql.Int32(access)),
		)).
		DISTINCT()

	return []mysql.CommonTableExpression{actorSubjectsClause}, subjectRows
}

func (b *calculatedVisibilityBackend) insertPublicRow(
	ctx context.Context,
	tx qrm.DB,
	targetID int64,
) error {
	stmt := b.access.calculatedVisibilityPublicTable.
		INSERT(b.access.calculatedVisibilityPublicTargetID).
		VALUES(targetID)
	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (b *calculatedVisibilityBackend) insertCreatorRow(
	ctx context.Context,
	tx qrm.DB,
	targetID int64,
	creatorID int64,
	creatorJob string,
) error {
	stmt := b.access.calculatedVisibilityCreatorTable.
		INSERT(
			b.access.calculatedVisibilityCreatorTargetID,
			b.access.calculatedVisibilityCreatorCreatorID,
			b.access.calculatedVisibilityCreatorCreatorJob,
		).
		VALUES(targetID, creatorID, creatorJob)
	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (b *calculatedVisibilityBackend) insertSubjectRow(
	ctx context.Context,
	tx qrm.DB,
	targetID int64,
	subjectID int64,
	access int32,
	effect bool,
) error {
	stmt := b.access.calculatedVisibilitySubjectTable.
		INSERT(
			b.access.calculatedVisibilitySubjectTargetID,
			b.access.calculatedVisibilitySubjectSubjectID,
			b.access.calculatedVisibilitySubjectAccess,
			b.access.calculatedVisibilitySubjectEffect,
		).
		VALUES(targetID, subjectID, access, effect)
	_, err := stmt.ExecContext(ctx, tx)
	return err
}

type calculatedVisibilityTargetRow struct {
	ID         int64   `alias:"id"`
	Public     bool    `alias:"public"`
	CreatorID  *int32  `alias:"creator_id"`
	CreatorJob *string `alias:"creator_job"`
}

func (b *calculatedVisibilityBackend) loadTargetRow(
	ctx context.Context,
	tx qrm.DB,
	targetID int64,
) (*calculatedVisibilityTargetRow, error) {
	columns := []mysql.Projection{
		b.access.targetTableColumns.ID.AS("id"),
	}
	if b.access.targetTableColumns.Public != nil {
		columns = append(columns, b.access.targetTableColumns.Public.AS("public"))
	}
	if b.access.targetTableColumns.CreatorID != nil {
		columns = append(columns,
			b.access.targetTableColumns.CreatorID.AS("creator_id"),
			b.access.targetTableColumns.CreatorJob.AS("creator_job"),
		)
	}

	if len(columns) == 0 {
		return nil, nil
	}

	stmt := mysql.SELECT(
		columns[0],
		columns[1:]...,
	).
		FROM(b.access.targetTable).
		WHERE(b.access.targetTableColumns.ID.EQ(mysql.Int64(targetID))).
		LIMIT(1)

	var dest []calculatedVisibilityTargetRow
	if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if len(dest) == 0 {
		return nil, nil
	}
	return &dest[0], nil
}

type calculatedVisibilityCreatorRow struct {
	CreatorID  *int32  `alias:"creator_id"`
	CreatorJob *string `alias:"creator_job"`
}

func (b *calculatedVisibilityBackend) loadTargetCreatorRow(
	ctx context.Context,
	tx qrm.DB,
	targetID int64,
) (*calculatedVisibilityCreatorRow, error) {
	stmt := mysql.SELECT(
		b.access.targetTableColumns.CreatorID.AS("creator_id"),
		b.access.targetTableColumns.CreatorJob.AS("creator_job"),
	).
		FROM(b.access.targetTable).
		WHERE(b.access.targetTableColumns.ID.EQ(mysql.Int64(targetID))).
		LIMIT(1)

	var dest []calculatedVisibilityCreatorRow
	if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if len(dest) == 0 {
		return nil, nil
	}
	return &dest[0], nil
}

type calculatedVisibilityACLRow struct {
	SubjectID int64 `alias:"subject_id"`
	Access    int32 `alias:"access"`
	Effect    bool  `alias:"effect"`
}

func (b *calculatedVisibilityBackend) loadACLRows(
	ctx context.Context,
	tx qrm.DB,
	targetID int64,
) ([]calculatedVisibilityACLRow, error) {
	stmt := b.access.accessTable.
		SELECT(
			b.access.accessColumns.SubjectID.AS("subject_id"),
			b.access.accessColumns.Access.AS("access"),
			b.access.accessColumns.Effect.AS("effect"),
		).
		FROM(b.access.accessTable).
		WHERE(b.access.accessColumns.TargetID.EQ(mysql.Int64(targetID))).
		ORDER_BY(b.access.accessColumns.ID.ASC())

	var dest []calculatedVisibilityACLRow
	if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return dest, nil
}
