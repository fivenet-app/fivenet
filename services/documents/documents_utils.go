package documents

import (
	context "context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	permscitizens "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	errorsdocuments "github.com/fivenet-app/fivenet/v2025/services/documents/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) listDocumentsQuery(where jet.BoolExpression, onlyColumns jet.ProjectionList, additionalColumns jet.ProjectionList, userInfo *userinfo.UserInfo) jet.SelectStatement {
	tCreator := tables.User().AS("creator")

	wheres := []jet.BoolExpression{}
	if !userInfo.Superuser {
		wheres = []jet.BoolExpression{
			jet.AND(
				tDocumentShort.DeletedAt.IS_NULL(),
				jet.OR(
					tDocumentShort.Public.IS_TRUE(),
					jet.AND(
						tDocumentShort.CreatorID.EQ(jet.Int32(userInfo.UserId)),
						tDocumentShort.CreatorJob.EQ(jet.String(userInfo.Job)),
					),
					jet.OR(
						tDAccess.UserID.EQ(jet.Int32(userInfo.UserId)),
						jet.AND(
							tDAccess.Job.EQ(jet.String(userInfo.Job)),
							tDAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade)),
						),
					),
				),
			),
		}
	}

	if where != nil {
		wheres = append(wheres, where)
	}

	var q jet.SelectStatement
	if onlyColumns != nil {
		q = tDocumentShort.
			SELECT(
				onlyColumns,
			)
	} else {
		columns := jet.ProjectionList{
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
			tDocumentShort.Summary.AS("document_short.content"),
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
			tDocumentShort.Public,
			tDocumentShort.TemplateID,
			tWorkflow.DocumentID,
			tWorkflow.AutoCloseTime,
			tWorkflow.NextReminderTime,
		}

		if userInfo.Superuser {
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

		q = tDocumentShort.SELECT(columns[0], columns[1:])
	}

	var tables jet.ReadableTable
	if !userInfo.Superuser {
		tables = tDocumentShort.
			LEFT_JOIN(tDAccess,
				tDAccess.TargetID.EQ(tDocumentShort.ID).
					AND(tDAccess.Access.GT_EQ(jet.Int32(int32(documents.AccessLevel_ACCESS_LEVEL_VIEW)))),
			).
			LEFT_JOIN(tDCategory,
				tDocumentShort.CategoryID.EQ(tDCategory.ID).
					AND(tDCategory.DeletedAt.IS_NULL()),
			).
			LEFT_JOIN(tCreator,
				tDocumentShort.CreatorID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tWorkflow,
				tWorkflow.DocumentID.EQ(tDocumentShort.ID),
			)
	} else {
		tables = tDocumentShort.
			LEFT_JOIN(tDCategory,
				tDocumentShort.CategoryID.EQ(tDCategory.ID).
					AND(tDCategory.DeletedAt.IS_NULL()),
			).
			LEFT_JOIN(tCreator,
				tDocumentShort.CreatorID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tWorkflow,
				tWorkflow.DocumentID.EQ(tDocumentShort.ID),
			)
	}

	return q.
		FROM(tables).
		WHERE(jet.AND(
			wheres...,
		))
}

func (s *Server) getDocumentQuery(where jet.BoolExpression, onlyColumns jet.ProjectionList, userInfo *userinfo.UserInfo, withContent bool) jet.SelectStatement {
	tCreator := tables.User().AS("creator")

	var wheres []jet.BoolExpression
	if !userInfo.Superuser {
		wheres = []jet.BoolExpression{
			jet.AND(
				tDocument.DeletedAt.IS_NULL(),
				jet.OR(
					jet.OR(
						tDocument.Public.IS_TRUE(),
						tDocument.CreatorID.EQ(jet.Int32(userInfo.UserId)),
					),
					jet.OR(
						tDAccess.UserID.EQ(jet.Int32(userInfo.UserId)),
						jet.AND(
							tDAccess.Job.EQ(jet.String(userInfo.Job)),
							tDAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade)),
						),
					),
				),
			),
		}
	}

	if where != nil {
		wheres = append(wheres, where)
	}

	var q jet.SelectStatement
	if onlyColumns != nil {
		q = tDocument.
			SELECT(
				onlyColumns,
			)
	} else {
		columns := jet.ProjectionList{
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
			tWorkflow.DocumentID,
			tWorkflow.AutoCloseTime,
			tWorkflow.NextReminderTime,
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

		if userInfo.Superuser {
			columns = append(columns, tDocument.DeletedAt)
		}

		// Field Permission Check
		fields, _ := s.ps.AttrStringList(userInfo, permscitizens.CitizensServicePerm, permscitizens.CitizensServiceListCitizensPerm, permscitizens.CitizensServiceListCitizensFieldsPermField)
		if fields.Contains("PhoneNumber") {
			columns = append(columns, tCreator.PhoneNumber)
		}

		q = tDocument.SELECT(columns[0], columns[1:])
	}

	var tables jet.ReadableTable
	if !userInfo.Superuser {
		tables = tDocument.
			LEFT_JOIN(tDAccess,
				tDAccess.TargetID.EQ(tDocument.ID).
					AND(tDAccess.Access.GT_EQ(jet.Int32(int32(documents.AccessLevel_ACCESS_LEVEL_VIEW)))),
			).
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
			LEFT_JOIN(tWorkflow,
				tWorkflow.DocumentID.EQ(tDocument.ID),
			).
			LEFT_JOIN(tUserWorkflow,
				tUserWorkflow.DocumentID.EQ(tDocument.ID).
					AND(tUserWorkflow.UserID.EQ(jet.Int32(userInfo.UserId))),
			)
	} else {
		tables = tDocument.
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
			LEFT_JOIN(tWorkflow,
				tWorkflow.DocumentID.EQ(tDocument.ID),
			).
			LEFT_JOIN(tUserWorkflow,
				tUserWorkflow.DocumentID.EQ(tDocument.ID).
					AND(tUserWorkflow.UserID.EQ(jet.Int32(userInfo.UserId))),
			)
	}

	return q.
		FROM(tables).
		WHERE(jet.AND(
			wheres...,
		)).
		ORDER_BY(
			tDocument.CreatedAt.DESC(),
			tDocument.UpdatedAt.DESC(),
		)
}

func (s *Server) updateDocumentOwner(ctx context.Context, tx qrm.DB, documentId uint64, userInfo *userinfo.UserInfo, newOwner *users.UserShort) error {
	stmt := tDocument.
		UPDATE(
			tDocument.CreatorID,
		).
		SET(
			newOwner.UserId,
		).
		WHERE(
			tDocument.ID.EQ(jet.Uint64(documentId)),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if _, err := addDocumentActivity(ctx, tx, &documents.DocActivity{
		DocumentId:   documentId,
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_OWNER_CHANGED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.Job,
		Data: &documents.DocActivityData{
			Data: &documents.DocActivityData_OwnerChanged{
				OwnerChanged: &documents.DocOwnerChanged{
					NewOwnerId: newOwner.UserId,
					NewOwner:   newOwner,
				},
			},
		},
	}); err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return nil
}
