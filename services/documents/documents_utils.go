package documents

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	permscitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
)

func (s *Server) getDocumentQuery(
	where mysql.BoolExpression,
	onlyColumns mysql.ProjectionList,
	userInfo *userinfo.UserInfo,
	withContent bool,
) mysql.SelectStatement {
	tCreator := table.FivenetUser.AS("creator")

	var wheres []mysql.BoolExpression
	if !userInfo.GetSuperuser() {
		wheres = []mysql.BoolExpression{
			mysql.AND(
				tDocument.DeletedAt.IS_NULL(),
				mysql.OR(
					tDocument.Public.IS_TRUE(),
					tDocument.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
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
		fields, _ := permscitizens.CitizensService.ListCitizens.FieldsTyped.Get(s.ps, userInfo)
		if fields.Contains(permscitizens.CitizensServiceListCitizensFieldsPermValuePhoneNumber) {
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
