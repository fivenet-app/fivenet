package documentsstore

import (
	"context"
	"errors"

	resourcesaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/access"
	resourcesdatabase "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	resourcesdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents"
	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	documentsactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/activity"
	documentspins "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/pins"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2026/services/documents/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) List(
	ctx context.Context,
	q ListQuery,
) ([]*resourcesdocuments.DocumentShort, error) {
	tDocumentShort := table.FivenetDocuments.AS("document_short")
	tDocumentPage := table.FivenetDocuments.AS("document_page")
	tCreator := table.FivenetUser.AS("creator")
	tDCategory := table.FivenetDocumentsCategories.AS("category")
	tWorkflowState := table.FivenetDocumentsWorkflowState.AS("workflow_state")
	tDMeta := table.FivenetDocumentsMeta.AS("meta")
	superuser := q.UserInfo != nil && q.UserInfo.GetSuperuser()

	columns := dbutils.Columns{
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
	}
	if superuser {
		columns = append(columns, tDocumentShort.DeletedAt)
	}
	if q.IncludePhoneNumber {
		columns = append(columns, tCreator.PhoneNumber)
	}

	var (
		stmt mysql.Statement
		ctes []mysql.CommonTableExpression
	)
	if superuser {
		stmt = mysql.
			SELECT(
				columns[0],
				columns[1:]...,
			).
			FROM(tDocumentShort.
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
			WHERE(buildListDocumentsCondition(tDocumentShort, q)).
			ORDER_BY(buildListOrder(
				q.Sort,
				tDocumentShort.Title,
				tDocumentShort.CreatedAt,
				tDocumentShort.UpdatedAt,
			)...).
			OFFSET(q.Offset).
			LIMIT(q.Limit)
	} else {
		visibleIDs := s.subjectAccess.VisibleIDsByConditionQuery(
			q.UserInfo,
			int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
			false,
			buildListDocumentsCondition(tDocumentPage, q),
		)
		ctes = visibleIDs.CTEs
		visibleDocID := mysql.IntegerColumn("id").From(visibleIDs.Table)
		docPage := mysql.
			SELECT(
				visibleDocID.AS("id"),
				tDocumentPage.CreatedAt.AS("created_at"),
				tDocumentPage.UpdatedAt.AS("updated_at"),
			).
			FROM(visibleIDs.Table.
				INNER_JOIN(tDocumentPage,
					tDocumentPage.ID.EQ(visibleDocID),
				),
			).
			ORDER_BY(buildListOrder(
				q.Sort,
				tDocumentPage.Title,
				tDocumentPage.CreatedAt,
				tDocumentPage.UpdatedAt,
			)...).
			OFFSET(q.Offset).
			LIMIT(q.Limit).
			AsTable("doc_page")

		docPageID := mysql.IntegerColumn("id").From(docPage)

		stmt = mysql.
			SELECT(
				columns[0],
				columns[1:]...,
			).
			FROM(docPage.
				INNER_JOIN(tDocumentShort, tDocumentShort.ID.EQ(docPageID)).
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
			ORDER_BY(buildListOrder(
				q.Sort,
				tDocumentShort.Title,
				tDocumentShort.CreatedAt,
				tDocumentShort.UpdatedAt,
			)...)
	}

	if len(ctes) > 0 {
		stmt = mysql.WITH(ctes...)(stmt)
	}
	var docs []*resourcesdocuments.DocumentShort
	if err := stmt.QueryContext(ctx, s.db, &docs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return docs, nil
}

func (s *Store) Get(ctx context.Context, q GetQuery) (*resourcesdocuments.Document, error) {
	tDocument := table.FivenetDocuments.AS("document")
	tCreator := table.FivenetUser.AS("creator")
	tDCategory := table.FivenetDocumentsCategories.AS("category")
	tDPins := table.FivenetDocumentsPins.AS("pin")
	tWorkflowState := table.FivenetDocumentsWorkflowState.AS("workflow_state")
	tUserWorkflow := table.FivenetDocumentsWorkflowUsers.AS("workflow_user_state")
	tDMeta := table.FivenetDocumentsMeta.AS("meta")

	if !q.UserInfo.GetSuperuser() {
		visible, err := s.subjectAccess.CanUserAccessTarget(
			ctx,
			q.DocumentID,
			q.UserInfo,
			int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		)
		if err != nil {
			return nil, err
		}
		if !visible {
			return nil, nil
		}
	}

	wheres := []mysql.BoolExpression{
		tDocument.ID.EQ(mysql.Int64(q.DocumentID)),
	}
	if !q.UserInfo.GetSuperuser() {
		wheres = append(wheres, tDocument.DeletedAt.IS_NULL())
	}

	columns := dbutils.Columns{
		tDocument.ID,
		tDocument.CreatedAt,
		tDocument.UpdatedAt,
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
	}

	if q.WithContent {
		columns = append(columns,
			tDocument.Data,
			tDocument.ContentJSON,
		)
	}
	if q.UserInfo.GetSuperuser() {
		columns = append(columns, tDocument.DeletedAt)
	}
	if q.IncludePhoneNumber {
		columns = append(columns, tCreator.PhoneNumber)
	}

	stmt := tDocument.
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
					tUserWorkflow.UserID.EQ(mysql.Int32(q.UserInfo.GetUserId())),
				),
			).
			LEFT_JOIN(tDMeta,
				tDMeta.DocumentID.EQ(tDocument.ID),
			),
		).
		WHERE(mysql.AND(wheres...)).
		ORDER_BY(
			tDocument.CreatedAt.DESC(),
			tDocument.UpdatedAt.DESC(),
		).
		LIMIT(1)

	var doc resourcesdocuments.Document
	if err := stmt.QueryContext(ctx, s.db, &doc); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if doc.GetId() <= 0 {
		return nil, nil
	}
	if doc.Meta == nil {
		doc.Meta = &resourcesdocuments.DocumentMeta{DocumentId: doc.GetId()}
	}

	return &doc, nil
}

func (s *Store) GetDocumentAccess(
	ctx context.Context,
	documentID int64,
) (*documentsaccess.DocumentAccess, error) {
	return s.subjectAccess.ListTargetAccess(ctx, s.db, documentID, documentSubjectAccessOptions)
}

func (s *Store) GetDocumentMeta(
	ctx context.Context,
	db qrm.DB,
	documentID int64,
) (*resourcesdocuments.DocumentMeta, error) {
	tDMeta := table.FivenetDocumentsMeta.AS("document_meta")

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
		WHERE(tDMeta.DocumentID.EQ(mysql.Int64(documentID))).
		LIMIT(1)

	dest := &resourcesdocuments.DocumentMeta{}
	if err := stmt.QueryContext(ctx, db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.DocumentId == 0 {
		dest.DocumentId = documentID
	}

	return dest, nil
}

func (s *Store) GetDocumentRequiredAccess(
	ctx context.Context,
	documentID int64,
	userInfo *userinfo.UserInfo,
) (*resourcesaccess.Access, error) {
	doc, err := s.Get(ctx, GetQuery{
		DocumentID: documentID,
		UserInfo:   userInfo,
	})
	if err != nil {
		return nil, err
	}
	if doc.GetTemplateId() <= 0 {
		return nil, nil
	}

	tmpl, err := s.GetTemplate(ctx, doc.GetTemplateId(), false)
	if err != nil {
		return nil, err
	}
	if tmpl == nil || tmpl.GetContentAccess() == nil || tmpl.GetContentAccess().IsEmpty() {
		return nil, nil
	}

	return tmpl.GetContentAccess(), nil
}

func (s *Store) GetDocumentPin(
	ctx context.Context,
	documentID int64,
	userInfo *userinfo.UserInfo,
) (*documentspins.DocumentPin, error) {
	tDPins := table.FivenetDocumentsPins.AS("document_pin")

	condition := mysql.AND(
		tDPins.DocumentID.EQ(mysql.Int64(documentID)),
		mysql.OR(
			mysql.AND(
				tDPins.Job.EQ(mysql.String(userInfo.GetJob())),
				tDPins.UserID.IS_NULL(),
			),
			mysql.AND(
				tDPins.Job.IS_NULL(),
				tDPins.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
			),
		),
	)

	stmt := tDPins.
		SELECT(
			tDPins.DocumentID,
			tDPins.Job,
			tDPins.UserID,
			tDPins.CreatedAt,
			tDPins.State,
			tDPins.CreatorID,
		).
		WHERE(condition).
		ORDER_BY(tDPins.UserID.DESC()).
		LIMIT(2)

	pins := []*documentspins.DocumentPin{}
	if err := stmt.QueryContext(ctx, s.db, &pins); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if len(pins) == 0 {
		return nil, nil
	}

	pin := &documentspins.DocumentPin{DocumentId: documentID}
	for _, p := range pins {
		if p.Job != nil {
			pin.Job = p.Job
		}
		if p.UserId != nil {
			pin.UserId = p.UserId
		}
		pin.State = p.GetState()
		pin.CreatedAt = p.GetCreatedAt()
		pin.CreatorId = p.GetCreatorId()
	}

	return pin, nil
}

func (s *Store) UpdateDocumentAccess(
	ctx context.Context,
	tx qrm.DB,
	documentID int64,
	userInfo *userinfo.UserInfo,
	docAccess *documentsaccess.DocumentAccess,
	addActivity bool,
) error {
	if docAccess == nil {
		docAccess = &documentsaccess.DocumentAccess{}
	}

	if documentsaccess.DocumentAccessHasDuplicates(docAccess) {
		return errorsdocuments.ErrDocAccessDuplicate
	}

	changes, err := s.subjectAccess.ReplaceTargetAccess(
		ctx,
		tx,
		s.subjectResolver,
		documentID,
		docAccess,
		documentSubjectAccessOptions,
	)
	if err != nil {
		if dbutils.IsDuplicateError(err) {
			return errorsdocuments.ErrDocAccessDuplicate
		}
		return err
	}

	if addActivity && !changes.IsEmpty() {
		if _, err := addDocumentActivity(ctx, tx, &documentsactivity.DocActivity{
			DocumentId:   documentID,
			ActivityType: documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_ACCESS_UPDATED,
			CreatorId:    &userInfo.UserId,
			CreatorJob:   userInfo.GetJob(),
			Data: &documentsactivity.DocActivityData{
				Data: &documentsactivity.DocActivityData_AccessUpdated{
					AccessUpdated: &documentsactivity.DocAccessUpdated{
						Jobs: &documentsactivity.DocAccessJobsDiff{
							ToCreate: changes.Jobs.ToCreate,
							ToUpdate: changes.Jobs.ToUpdate,
							ToDelete: changes.Jobs.ToDelete,
						},
						Users: &documentsactivity.DocAccessUsersDiff{
							ToCreate: changes.Users.ToCreate,
							ToUpdate: changes.Users.ToUpdate,
							ToDelete: changes.Users.ToDelete,
						},
					},
				},
			},
		}); err != nil {
			return err
		}
	}

	return nil
}

func addDocumentActivity(
	ctx context.Context,
	tx qrm.DB,
	activity *documentsactivity.DocActivity,
) (int64, error) {
	stmt := table.FivenetDocumentsActivity.
		INSERT(
			table.FivenetDocumentsActivity.DocumentID,
			table.FivenetDocumentsActivity.ActivityType,
			table.FivenetDocumentsActivity.CreatorID,
			table.FivenetDocumentsActivity.CreatorJob,
			table.FivenetDocumentsActivity.Reason,
			table.FivenetDocumentsActivity.Data,
		).
		VALUES(
			activity.GetDocumentId(),
			activity.GetActivityType(),
			activity.CreatorId,
			activity.GetCreatorJob(),
			activity.Reason,
			activity.GetData(),
		)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		if !dbutils.IsDuplicateError(err) {
			return 0, err
		}
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (s *Store) UpdateDocumentOwner(
	ctx context.Context,
	tx qrm.DB,
	documentID int64,
	userInfo *userinfo.UserInfo,
	newOwner *usershort.UserShort,
) error {
	tDocument := table.FivenetDocuments.AS("document")
	stmt := tDocument.
		UPDATE(
			tDocument.CreatorID,
		).
		SET(
			newOwner.GetUserId(),
		).
		WHERE(
			tDocument.ID.EQ(mysql.Int64(documentID)),
		).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	if err := s.subjectAccess.RefreshTargetVisibilityWithCreator(
		ctx,
		tx,
		documentID,
		newOwner.GetUserId(),
		newOwner.GetJob(),
	); err != nil {
		return err
	}

	creatorVisibilityStmt := table.FivenetDocumentsVisibilityCreator.
		INSERT(
			table.FivenetDocumentsVisibilityCreator.TargetID,
			table.FivenetDocumentsVisibilityCreator.CreatorID,
			table.FivenetDocumentsVisibilityCreator.CreatorJob,
		).
		VALUES(
			documentID,
			newOwner.GetUserId(),
			newOwner.GetJob(),
		).
		ON_DUPLICATE_KEY_UPDATE(
			table.FivenetDocumentsVisibilityCreator.TargetID.SET(
				mysql.RawInt("VALUES(`target_id`)"),
			),
			table.FivenetDocumentsVisibilityCreator.CreatorID.SET(
				mysql.RawInt("VALUES(`creator_id`)"),
			),
			table.FivenetDocumentsVisibilityCreator.CreatorJob.SET(
				mysql.RawString("VALUES(`creator_job`)"),
			),
		)
	if _, err := creatorVisibilityStmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	tDocActivity := table.FivenetDocumentsActivity
	if _, err := tDocActivity.
		INSERT(
			tDocActivity.DocumentID,
			tDocActivity.ActivityType,
			tDocActivity.CreatorID,
			tDocActivity.CreatorJob,
			tDocActivity.Data,
		).
		VALUES(
			documentID,
			documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_OWNER_CHANGED,
			&userInfo.UserId,
			userInfo.GetJob(),
			&documentsactivity.DocActivityData{
				Data: &documentsactivity.DocActivityData_OwnerChanged{
					OwnerChanged: &documentsactivity.DocOwnerChanged{
						NewOwnerId: newOwner.GetUserId(),
						NewOwner:   newOwner,
					},
				},
			},
		).
		ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}

func buildListDocumentsCondition(
	document *table.FivenetDocumentsTable,
	q ListQuery,
) mysql.BoolExpression {
	condition := mysql.Bool(true)
	if q.Search != "" {
		search := mysql.String(dbutils.PrepareForLikeSearch(q.Search))
		searchCategory := table.FivenetDocumentsCategories.AS("search_category")
		condition = condition.AND(mysql.OR(
			document.Title.LIKE(search),
			document.ContentText.LIKE(search),
			document.CreatorJob.LIKE(search),
			mysql.EXISTS(
				searchCategory.
					SELECT(mysql.Int(1)).
					FROM(searchCategory).
					WHERE(mysql.AND(
						searchCategory.ID.EQ(document.CategoryID),
						searchCategory.DeletedAt.IS_NULL(),
						mysql.OR(
							searchCategory.Name.LIKE(search),
							searchCategory.Description.LIKE(search),
						),
					)),
			),
		))
	}
	if len(q.CategoryIDs) > 0 {
		ids := make([]mysql.Expression, len(q.CategoryIDs))
		for i := range q.CategoryIDs {
			ids[i] = mysql.Int64(q.CategoryIDs[i])
		}
		condition = condition.AND(document.CategoryID.IN(ids...))
	}
	if len(q.CreatorIDs) > 0 {
		ids := make([]mysql.Expression, len(q.CreatorIDs))
		for i := range q.CreatorIDs {
			ids[i] = mysql.Int32(q.CreatorIDs[i])
		}
		condition = condition.AND(document.CreatorID.IN(ids...))
	}
	if q.From != nil {
		condition = condition.AND(document.CreatedAt.GT_EQ(mysql.TimestampT(q.From.AsTime())))
	}
	if q.To != nil {
		condition = condition.AND(document.CreatedAt.LT_EQ(mysql.TimestampT(q.To.AsTime())))
	}
	if q.Closed != nil {
		condition = condition.AND(document.Closed.EQ(mysql.Bool(*q.Closed)))
	}
	if len(q.DocumentIDs) > 0 {
		ids := make([]mysql.Expression, len(q.DocumentIDs))
		for i := range q.DocumentIDs {
			ids[i] = mysql.Int64(q.DocumentIDs[i])
		}
		condition = condition.AND(document.ID.IN(ids...))
	}
	if q.OnlyDrafts != nil {
		condition = condition.AND(document.Draft.EQ(mysql.Bool(*q.OnlyDrafts)))
	}

	return condition
}

func buildListOrder(
	sort *resourcesdatabase.Sort,
	title mysql.ColumnString,
	createdAt mysql.ColumnTimestamp,
	updatedAt mysql.ColumnTimestamp,
) []mysql.OrderByClause {
	orderBys := []mysql.OrderByClause{}
	if sort != nil && len(sort.GetColumns()) > 0 {
		for _, sc := range sort.GetColumns() {
			var column mysql.Column
			switch sc.GetId() {
			case "title":
				column = title
			case "createdAt":
				fallthrough
			default:
				column = createdAt
			}

			if sc.GetDesc() {
				orderBys = append(orderBys, column.DESC(), updatedAt.DESC())
			} else {
				orderBys = append(orderBys, column.ASC(), updatedAt.DESC())
			}
		}
	} else {
		orderBys = append(orderBys, updatedAt.DESC())
	}

	return orderBys
}
