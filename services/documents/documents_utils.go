package documents

import (
	context "context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	permscitizens "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	errorsdocuments "github.com/fivenet-app/fivenet/v2025/services/documents/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) listDocumentsQuery(
	where jet.BoolExpression,
	onlyColumns jet.ProjectionList,
	additionalColumns jet.ProjectionList,
	userInfo *userinfo.UserInfo,
	customizeFn func(stmt jet.SelectStatement) jet.SelectStatement,
) jet.Statement {
	tCreator := tables.User().AS("creator")

	wheres := []jet.BoolExpression{}
	if where != nil {
		wheres = append(wheres, where)
	}

	if !userInfo.GetSuperuser() {
	}

	pubSel := tDocumentShort.
		SELECT(
			tDocumentShort.ID.AS("id"),
			tDocumentShort.CreatedAt.AS("created_at"),
		).
		FROM(tDocumentShort).
		WHERE(jet.AND(
			tDocumentShort.DeletedAt.IS_NULL(),
			tDocumentShort.Public.EQ(jet.Bool(true)),
		))

	// branch 2: created by this user + job
	creatorSel := tDocumentShort.
		SELECT(
			tDocumentShort.ID.AS("id"),
			tDocumentShort.CreatedAt.AS("created_at"),
		).
		FROM(tDocumentShort).
		WHERE(tDocumentShort.DeletedAt.IS_NULL().
			AND(tDocumentShort.CreatorID.EQ(jet.Int32(userInfo.GetUserId()))).
			AND(tDocumentShort.CreatorJob.EQ(jet.String(userInfo.GetJob()))))

	existsAccess := jet.EXISTS(
		jet.
			SELECT(jet.Int(1)).
			FROM(tDAccess).
			WHERE(
				jet.AND(
					tDAccess.TargetID.EQ(tDocumentShort.ID),
					jet.OR(
						// Direct user access
						tDAccess.UserID.EQ(jet.Int32(userInfo.GetUserId())),
						// or job + grade access
						jet.AND(
							tDAccess.Job.EQ(jet.String(userInfo.GetJob())),
							tDAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.GetJobGrade())),
						),
					),
					tDAccess.Access.GT_EQ(
						jet.Int32(int32(documents.AccessLevel_ACCESS_LEVEL_VIEW)),
					),
				),
			),
	)

	accessSel := tDocumentShort.SELECT(
		tDocumentShort.ID.AS("id"),
		tDocumentShort.CreatedAt.AS("created_at"),
	).
		FROM(tDocumentShort).
		WHERE(jet.AND(
			tDocumentShort.DeletedAt.IS_NULL(),
			existsAccess,
		))

	var columns jet.ProjectionList
	if onlyColumns != nil {
		columns = onlyColumns
	} else {
		columns = jet.ProjectionList{
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
			// tDocumentShort.Summary.AS("document_short.content"), // Summary is unused at the moment
			tDocumentShort.CreatorID,
			tDocumentShort.TemplateID,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tDocumentShort.CreatorJob,
			tDocumentShort.State,
			tDocumentShort.Closed,
			tDocumentShort.Draft,
			tDocumentShort.Public,
			tDocumentShort.TemplateID,
			tDWorkflow.DocumentID,
			tDWorkflow.AutoCloseTime,
			tDWorkflow.NextReminderTime,
		}

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
	docIDs := jet.CTE("doc_ids")

	cteIDColumn := jet.IntegerColumn("id").From(docIDs)

	innerStmt := jet.
		SELECT(
			columns[0],
			columns[1:],
		).
		FROM(
			docIDs.
				INNER_JOIN(tDocumentShort, tDocumentShort.ID.EQ(cteIDColumn)).
				LEFT_JOIN(tDCategory,
					tDocumentShort.CategoryID.EQ(tDCategory.ID).
						AND(tDCategory.DeletedAt.IS_NULL()),
				).
				LEFT_JOIN(tCreator,
					tDocumentShort.CreatorID.EQ(tCreator.ID),
				).
				LEFT_JOIN(tDWorkflow,
					tDWorkflow.DocumentID.EQ(tDocumentShort.ID),
				),
		).
		WHERE(jet.AND(
			wheres...,
		))

	if customizeFn != nil {
		innerStmt = customizeFn(innerStmt)
	}

	return jet.WITH(
		docIDs.AS(
			pubSel.UNION(creatorSel).UNION(accessSel), // UNION (distinct) to dedupe
		),
	)(innerStmt)
}

func (s *Server) getDocumentQuery(
	where jet.BoolExpression,
	onlyColumns jet.ProjectionList,
	userInfo *userinfo.UserInfo,
	withContent bool,
) jet.SelectStatement {
	tCreator := tables.User().AS("creator")

	var wheres []jet.BoolExpression
	if !userInfo.GetSuperuser() {
		accessExists := jet.EXISTS(
			jet.
				SELECT(jet.Int(1)).
				FROM(tDAccess).
				WHERE(
					jet.AND(
						tDAccess.TargetID.EQ(tDocument.ID),
						jet.OR(
							// Direct user access
							tDAccess.UserID.EQ(jet.Int32(userInfo.GetUserId())),
							// or job + grade access
							jet.AND(
								tDAccess.Job.EQ(jet.String(userInfo.GetJob())),
								tDAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.GetJobGrade())),
							),
						),
						tDAccess.Access.GT_EQ(
							jet.Int32(int32(documents.AccessLevel_ACCESS_LEVEL_VIEW)),
						),
					),
				),
		)

		wheres = []jet.BoolExpression{
			jet.AND(
				tDocument.DeletedAt.IS_NULL(),
				jet.OR(
					tDocument.Public.IS_TRUE(),
					tDocument.CreatorID.EQ(jet.Int32(userInfo.GetUserId())),
					accessExists,
				),
			),
		}
	}

	if where != nil {
		wheres = append(wheres, where)
	}

	var columns []jet.Projection
	if onlyColumns != nil {
		columns = onlyColumns
	} else {
		columns = jet.ProjectionList{
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
			tDocument.State,
			tDocument.Closed,
			tDocument.Draft,
			tDocument.Public,
			tDocument.TemplateID,
			tDPins.State,
			tDPins.Job,
			tDPins.UserID,
			tDWorkflow.DocumentID,
			tDWorkflow.AutoCloseTime,
			tDWorkflow.NextReminderTime,
			tUserWorkflow.DocumentID,
			tUserWorkflow.UserID,
			tUserWorkflow.ManualReminderTime,
			tUserWorkflow.ManualReminderMessage,
		}

		if withContent {
			columns = append(columns,
				tDocument.Data,
				tDocument.Content,
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
				tDocument.CategoryID.EQ(tDCategory.ID).
					AND(tDCategory.DeletedAt.IS_NULL()),
			).
			LEFT_JOIN(tCreator,
				tDocument.CreatorID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tDPins,
				tDPins.DocumentID.EQ(tDocument.ID),
			).
			LEFT_JOIN(tDWorkflow,
				tDWorkflow.DocumentID.EQ(tDocument.ID),
			).
			LEFT_JOIN(tUserWorkflow,
				tUserWorkflow.DocumentID.EQ(tDocument.ID).
					AND(tUserWorkflow.UserID.EQ(jet.Int32(userInfo.GetUserId()))),
			),
		).
		WHERE(jet.AND(
			wheres...,
		)).
		ORDER_BY(
			tDocument.CreatedAt.DESC(),
			tDocument.UpdatedAt.DESC(),
		)
}

func (s *Server) updateDocumentOwner(
	ctx context.Context,
	tx qrm.DB,
	documentId int64,
	userInfo *userinfo.UserInfo,
	newOwner *users.UserShort,
) error {
	stmt := tDocument.
		UPDATE(
			tDocument.CreatorID,
		).
		SET(
			newOwner.GetUserId(),
		).
		WHERE(
			tDocument.ID.EQ(jet.Int64(documentId)),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if _, err := addDocumentActivity(ctx, tx, &documents.DocActivity{
		DocumentId:   documentId,
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_OWNER_CHANGED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.GetJob(),
		Data: &documents.DocActivityData{
			Data: &documents.DocActivityData_OwnerChanged{
				OwnerChanged: &documents.DocOwnerChanged{
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
