package documents

import (
	context "context"
	"errors"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents"
	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	documentsactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/activity"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	permscitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2026/services/documents/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) listDocumentsQuery(
	where mysql.BoolExpression,
	onlyColumns mysql.ProjectionList,
	additionalColumns mysql.ProjectionList,
	userInfo *userinfo.UserInfo,
	customizeFn func(stmt mysql.SelectStatement) mysql.SelectStatement,
) mysql.Statement {
	tCreator := table.FivenetUser.AS("creator")

	wheres := []mysql.BoolExpression{}
	if where != nil {
		wheres = append(wheres, where)
	}

	pubSel := tDocumentShort.
		SELECT(
			tDocumentShort.ID.AS("id"),
			tDocumentShort.CreatedAt.AS("created_at"),
		).
		FROM(tDocumentShort).
		WHERE(mysql.AND(
			tDocumentShort.DeletedAt.IS_NULL(),
			tDocumentShort.Public.EQ(mysql.Bool(true)),
		))

	// branch 2: created by this user + job
	creatorSel := tDocumentShort.
		SELECT(
			tDocumentShort.ID.AS("id"),
			tDocumentShort.CreatedAt.AS("created_at"),
		).
		FROM(tDocumentShort).
		WHERE(mysql.AND(
			tDocumentShort.DeletedAt.IS_NULL(),
			tDocumentShort.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
			tDocumentShort.CreatorJob.EQ(mysql.String(userInfo.GetJob())),
		))

	var existsAccess mysql.BoolExpression
	if !userInfo.GetSuperuser() {
		existsAccess = mysql.EXISTS(
			mysql.
				SELECT(mysql.Int(1)).
				FROM(tDAccess).
				WHERE(mysql.AND(
					tDAccess.TargetID.EQ(tDocumentShort.ID),
					mysql.OR(
						// Direct user access
						tDAccess.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
						// or job + grade access
						mysql.AND(
							tDAccess.Job.EQ(mysql.String(userInfo.GetJob())),
							tDAccess.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),
						),
					),
					tDAccess.Access.GT_EQ(
						mysql.Int32(int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW)),
					),
				)),
		)
	} else {
		existsAccess = mysql.Bool(true)
	}

	accessSel := tDocumentShort.SELECT(
		tDocumentShort.ID.AS("id"),
		tDocumentShort.CreatedAt.AS("created_at"),
	).
		FROM(tDocumentShort).
		WHERE(mysql.AND(
			tDocumentShort.DeletedAt.IS_NULL(),
			existsAccess,
		))

	var columns mysql.ProjectionList
	if onlyColumns != nil {
		columns = append(columns, onlyColumns...)
	} else {
		columns = append(columns,
			tDocumentShort.ID,
			tDocumentShort.CreatedAt,
			tDocumentShort.UpdatedAt,
			tDocumentShort.DeletedAt,
			tDocumentShort.CategoryID,
			tDCategory.ID,
			tDCategory.Name,
			tDCategory.Description,
			tDCategory.Job,
			tDCategory.Color,
			tDCategory.Icon,
			tDocumentShort.Title,
			tDocumentShort.ContentType,
			// tDocumentShort.Summary.AS("document_short.content"), // Summary is currently unused
			tDocumentShort.CreatorID,
			tDocumentShort.TemplateID,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tDocumentShort.CreatorJob,
			tWorkflowState.DocumentID,
			tWorkflowState.AutoCloseTime,
			tWorkflowState.NextReminderTime,
			tDocumentShort.State.AS("meta.state"),
			tDocumentShort.Closed.AS("meta.closed"),
			tDocumentShort.Draft.AS("meta.draft"),
			tDocumentShort.Public.AS("meta.public"),
			tDMeta.DocumentID,
			tDMeta.Approved,
			tDMeta.ApRequiredTotal,
			tDMeta.ApCollectedApproved,
			tDMeta.ApRequiredRemaining,
			tDMeta.ApDeclinedCount,
			tDMeta.ApPendingCount,
			tDMeta.ApAnyDeclined,
			tDMeta.ApPoliciesActive,
			tDMeta.CommentCount,
		)

		if userInfo.GetSuperuser() {
			columns = append(columns, tDocumentShort.DeletedAt)
		}

		// Field Permission Check
		fields, _ := s.ps.AttrStringList(userInfo, permscitizens.CitizensServicePerm, permscitizens.CitizensServiceListCitizensPerm, permscitizens.CitizensServiceListCitizensFieldsPermField)
		if fields.Contains("PhoneNumber") {
			columns = append(columns, tCreator.PhoneNumber)
		}

		if additionalColumns != nil {
			columns = append(columns, additionalColumns...)
		}
	}

	// Union
	docIDs := mysql.CTE("doc_ids")

	cteIDColumn := mysql.IntegerColumn("id").From(docIDs)

	innerStmt := mysql.
		SELECT(
			columns[0],
			columns[1:],
		).
		FROM(docIDs.
			INNER_JOIN(tDocumentShort, tDocumentShort.ID.EQ(cteIDColumn)).
			LEFT_JOIN(tDCategory,
				mysql.AND(
					tDocumentShort.CategoryID.EQ(tDCategory.ID),
					tDCategory.DeletedAt.IS_NULL(),
				),
			).
			LEFT_JOIN(tCreator,
				tDocumentShort.CreatorID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tWorkflowState,
				tWorkflowState.DocumentID.EQ(tDocumentShort.ID),
			).
			LEFT_JOIN(tDMeta,
				tDMeta.DocumentID.EQ(tDocumentShort.ID),
			),
		).
		WHERE(mysql.AND(
			wheres...,
		))

	if customizeFn != nil {
		innerStmt = customizeFn(innerStmt)
	}

	return mysql.WITH(
		docIDs.AS(
			pubSel.UNION(creatorSel).UNION(accessSel), // UNION (distinct) to dedupe
		),
	)(innerStmt)
}

func (s *Server) getDocumentQuery(
	where mysql.BoolExpression,
	onlyColumns mysql.ProjectionList,
	userInfo *userinfo.UserInfo,
	withContent bool,
) mysql.SelectStatement {
	tCreator := table.FivenetUser.AS("creator")

	var wheres []mysql.BoolExpression
	if !userInfo.GetSuperuser() {
		accessExists := mysql.EXISTS(
			mysql.
				SELECT(mysql.Int(1)).
				FROM(tDAccess).
				WHERE(mysql.AND(
					tDAccess.TargetID.EQ(tDocument.ID),
					mysql.OR(
						// Direct user access
						tDAccess.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
						// or job + grade access
						mysql.AND(
							tDAccess.Job.EQ(mysql.String(userInfo.GetJob())),
							tDAccess.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),
						),
					),
					tDAccess.Access.GT_EQ(
						mysql.Int32(int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW)),
					),
				)),
		)

		wheres = []mysql.BoolExpression{
			mysql.AND(
				tDocument.DeletedAt.IS_NULL(),
				mysql.OR(
					tDocument.Public.IS_TRUE(),
					tDocument.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
					accessExists,
				),
			),
		}
	}

	if where != nil {
		wheres = append(wheres, where)
	}

	var columns mysql.ProjectionList
	if onlyColumns != nil {
		columns = append(columns, onlyColumns)
	} else {
		columns = append(columns,
			tDocument.ID,
			tDocument.CreatedAt,
			tDocument.UpdatedAt,
			tDocument.DeletedAt,
			tDocument.CategoryID,
			tDCategory.ID,
			tDCategory.Name,
			tDCategory.Description,
			tDCategory.Job,
			tDCategory.Color,
			tDCategory.Icon,
			tDocument.Title,
			tDocument.ContentType,
			tDocument.CreatorID,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tDocument.CreatorJob,
			tDocument.State.AS("meta.state"),
			tDocument.Closed.AS("meta.closed"),
			tDocument.Draft.AS("meta.draft"),
			tDocument.Public.AS("meta.public"),
			tDocument.TemplateID,
			tDPins.State,
			tDPins.Job,
			tDPins.UserID,
			tWorkflowState.DocumentID,
			tWorkflowState.AutoCloseTime,
			tWorkflowState.NextReminderTime,
			tUserWorkflow.DocumentID,
			tUserWorkflow.UserID,
			tUserWorkflow.ManualReminderTime,
			tUserWorkflow.ManualReminderMessage,
			tDMeta.DocumentID,
			tDMeta.Approved,
			tDMeta.ApRequiredTotal,
			tDMeta.ApCollectedApproved,
			tDMeta.ApRequiredRemaining,
			tDMeta.ApDeclinedCount,
			tDMeta.ApPendingCount,
			tDMeta.ApAnyDeclined,
			tDMeta.ApPoliciesActive,
			tDMeta.CommentCount,
		)

		if withContent {
			columns = append(columns,
				tDocument.Data,
				tDocument.ContentJSON,
			)
		}

		if userInfo.GetSuperuser() {
			columns = append(columns, tDocument.DeletedAt)
		}

		// Field Permission Check
		fields, _ := s.ps.AttrStringList(userInfo, permscitizens.CitizensServicePerm, permscitizens.CitizensServiceListCitizensPerm, permscitizens.CitizensServiceListCitizensFieldsPermField)
		if fields.Contains("PhoneNumber") {
			columns = append(columns, tCreator.PhoneNumber)
		}
	}

	return tDocument.
		SELECT(
			columns[0],
			columns[1:]...,
		).
		FROM(tDocument.
			LEFT_JOIN(tDCategory,
				mysql.AND(
					tDocument.CategoryID.EQ(tDCategory.ID),
					tDCategory.DeletedAt.IS_NULL(),
				),
			).
			LEFT_JOIN(tCreator,
				tDocument.CreatorID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tDPins,
				tDPins.DocumentID.EQ(tDocument.ID),
			).
			LEFT_JOIN(tWorkflowState,
				tWorkflowState.DocumentID.EQ(tDocument.ID),
			).
			LEFT_JOIN(tUserWorkflow,
				mysql.AND(
					tUserWorkflow.DocumentID.EQ(tDocument.ID),
					tUserWorkflow.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
				),
			).
			LEFT_JOIN(tDMeta,
				tDMeta.DocumentID.EQ(tDocument.ID),
			),
		).
		WHERE(mysql.AND(
			wheres...,
		)).
		ORDER_BY(
			tDocument.CreatedAt.DESC(),
			tDocument.UpdatedAt.DESC(),
		)
}

func (s *Server) getDocumentMeta(
	ctx context.Context,
	tx qrm.DB,
	documentId int64,
) (*documents.DocumentMeta, error) {
	tDMeta := tDMeta.AS("document_meta")

	stmt := tDMeta.
		SELECT(
			tDMeta.DocumentID,
			tDMeta.Approved,
			tDMeta.ApRequiredTotal,
			tDMeta.ApCollectedApproved,
			tDMeta.ApRequiredRemaining,
			tDMeta.ApDeclinedCount,
			tDMeta.ApPendingCount,
			tDMeta.ApAnyDeclined,
			tDMeta.ApPoliciesActive,
			tDMeta.CommentCount,
		).
		FROM(tDMeta).
		WHERE(
			tDMeta.DocumentID.EQ(mysql.Int64(documentId)),
		).
		LIMIT(1)

	dest := &documents.DocumentMeta{}
	if err := stmt.QueryContext(ctx, tx, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	if dest.DocumentId == 0 {
		// No meta found, return "default" values
		dest.DocumentId = documentId
	}

	return dest, nil
}

func (s *Server) updateDocumentOwner(
	ctx context.Context,
	tx qrm.DB,
	documentId int64,
	userInfo *userinfo.UserInfo,
	newOwner *usershort.UserShort,
) error {
	stmt := tDocument.
		UPDATE(
			tDocument.CreatorID,
		).
		SET(
			newOwner.GetUserId(),
		).
		WHERE(
			tDocument.ID.EQ(mysql.Int64(documentId)),
		).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if _, err := addDocumentActivity(ctx, tx, &documentsactivity.DocActivity{
		DocumentId:   documentId,
		ActivityType: documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_OWNER_CHANGED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.GetJob(),
		Data: &documentsactivity.DocActivityData{
			Data: &documentsactivity.DocActivityData_OwnerChanged{
				OwnerChanged: &documentsactivity.DocOwnerChanged{
					NewOwnerId: newOwner.GetUserId(),
					NewOwner:   newOwner,
				},
			},
		},
	}); err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return nil
}
