package documentsstore

import (
	"context"
	"database/sql"
	"errors"

	resourcesdatabase "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	resourcesdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents"
	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	documentsactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/activity"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type IStore interface {
	List(ctx context.Context, q ListQuery) ([]*resourcesdocuments.DocumentShort, error)
	Get(ctx context.Context, q GetQuery) (*resourcesdocuments.Document, error)
	UpdateDocumentOwner(
		ctx context.Context,
		tx qrm.DB,
		documentID int64,
		userInfo *userinfo.UserInfo,
		newOwner *usershort.UserShort,
	) error
}

type Store struct {
	db            *sql.DB
	subjectAccess *access.SubjectObjectAccess
}

type ListQuery struct {
	Search      string
	CategoryIDs []int64
	CreatorIDs  []int32
	From        *timestamp.Timestamp
	To          *timestamp.Timestamp
	Closed      *bool
	DocumentIDs []int64
	OnlyDrafts  *bool
	Sort        *resourcesdatabase.Sort
	Offset      int64
	Limit       int64

	IncludePhoneNumber bool
	UserInfo           *userinfo.UserInfo
}

type GetQuery struct {
	DocumentID         int64
	IncludePhoneNumber bool
	WithContent        bool
	UserInfo           *userinfo.UserInfo
}

func New(db *sql.DB) IStore {
	return &Store{
		db:            db,
		subjectAccess: access.NewDocumentsSubjectObjectAccess(db),
	}
}

func (s *Store) List(ctx context.Context, q ListQuery) ([]*resourcesdocuments.DocumentShort, error) {
	if q.UserInfo == nil {
		q.UserInfo = &userinfo.UserInfo{}
	}

	tDocumentShort := table.FivenetDocuments.AS("document_short")
	tCreator := table.FivenetUser.AS("creator")
	tDCategory := table.FivenetDocumentsCategories.AS("category")
	tWorkflowState := table.FivenetDocumentsWorkflowState.AS("workflow_state")
	tDMeta := table.FivenetDocumentsMeta.AS("meta")

	condition := mysql.Bool(true)
	if q.Search != "" {
		condition = dbutils.MATCH(tDocumentShort.Title, mysql.String(q.Search))
	}
	if len(q.CategoryIDs) > 0 {
		ids := make([]mysql.Expression, len(q.CategoryIDs))
		for i := range q.CategoryIDs {
			ids[i] = mysql.Int64(q.CategoryIDs[i])
		}
		condition = condition.AND(tDocumentShort.CategoryID.IN(ids...))
	}
	if len(q.CreatorIDs) > 0 {
		ids := make([]mysql.Expression, len(q.CreatorIDs))
		for i := range q.CreatorIDs {
			ids[i] = mysql.Int32(q.CreatorIDs[i])
		}
		condition = condition.AND(tDocumentShort.CreatorID.IN(ids...))
	}
	if q.From != nil {
		condition = condition.AND(tDocumentShort.CreatedAt.GT_EQ(mysql.TimestampT(q.From.AsTime())))
	}
	if q.To != nil {
		condition = condition.AND(tDocumentShort.CreatedAt.LT_EQ(mysql.TimestampT(q.To.AsTime())))
	}
	if q.Closed != nil {
		condition = condition.AND(tDocumentShort.Closed.EQ(mysql.Bool(*q.Closed)))
	}
	if len(q.DocumentIDs) > 0 {
		ids := make([]mysql.Expression, len(q.DocumentIDs))
		for i := range q.DocumentIDs {
			ids[i] = mysql.Int64(q.DocumentIDs[i])
		}
		condition = condition.AND(tDocumentShort.ID.IN(ids...))
	}
	if q.OnlyDrafts != nil {
		condition = condition.AND(tDocumentShort.Draft.EQ(mysql.Bool(*q.OnlyDrafts)))
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

	creatorSel := tDocumentShort.
		SELECT(
			tDocumentShort.ID.AS("id"),
			tDocumentShort.CreatedAt.AS("created_at"),
		).
		FROM(tDocumentShort).
		WHERE(mysql.AND(
			tDocumentShort.DeletedAt.IS_NULL(),
			tDocumentShort.CreatorID.EQ(mysql.Int32(q.UserInfo.GetUserId())),
			tDocumentShort.CreatorJob.EQ(mysql.String(q.UserInfo.GetJob())),
		))

	var existsAccess mysql.BoolExpression
	if !q.UserInfo.GetSuperuser() {
		existsAccess = s.subjectAccess.ACLAccessExistsCondition(
			tDocumentShort.ID,
			q.UserInfo,
			int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		)
	} else {
		existsAccess = mysql.Bool(true)
	}

	accessSel := tDocumentShort.
		SELECT(
			tDocumentShort.ID.AS("id"),
			tDocumentShort.CreatedAt.AS("created_at"),
		).
		FROM(tDocumentShort).
		WHERE(mysql.AND(
			tDocumentShort.DeletedAt.IS_NULL(),
			existsAccess,
		))

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

	if q.UserInfo.GetSuperuser() {
		columns = append(columns, tDocumentShort.DeletedAt)
	}
	if q.IncludePhoneNumber {
		columns = append(columns, tCreator.PhoneNumber)
	}

	orderBys := buildListOrder(q.Sort, tDocumentShort)
	docIDs := mysql.CTE("doc_ids")
	cteIDColumn := mysql.IntegerColumn("id").From(docIDs)

	innerStmt := mysql.
		SELECT(
			columns[0],
			columns[1:]...,
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
		WHERE(condition).
		ORDER_BY(orderBys...).
		OFFSET(q.Offset).
		LIMIT(q.Limit)

	stmt := mysql.WITH(
		docIDs.AS(pubSel.UNION(creatorSel).UNION(accessSel)),
	)(innerStmt)

	var docs []*resourcesdocuments.DocumentShort
	if err := stmt.QueryContext(ctx, s.db, &docs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return docs, nil
}

func (s *Store) Get(ctx context.Context, q GetQuery) (*resourcesdocuments.Document, error) {
	if q.UserInfo == nil {
		q.UserInfo = &userinfo.UserInfo{}
	}

	tDocument := table.FivenetDocuments.AS("document")
	tCreator := table.FivenetUser.AS("creator")
	tDCategory := table.FivenetDocumentsCategories.AS("category")
	tDPins := table.FivenetDocumentsPins.AS("pin")
	tWorkflowState := table.FivenetDocumentsWorkflowState.AS("workflow_state")
	tUserWorkflow := table.FivenetDocumentsWorkflowUsers.AS("workflow_user_state")
	tDMeta := table.FivenetDocumentsMeta.AS("meta")

	var wheres []mysql.BoolExpression
	if !q.UserInfo.GetSuperuser() {
		existsAccess := s.subjectAccess.ACLAccessExistsCondition(
			tDocument.ID,
			q.UserInfo,
			int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		)

		wheres = []mysql.BoolExpression{
			mysql.AND(
				tDocument.DeletedAt.IS_NULL(),
				mysql.OR(
					tDocument.Public.IS_TRUE(),
					tDocument.CreatorID.EQ(mysql.Int32(q.UserInfo.GetUserId())),
					existsAccess,
				),
			),
		}
	}

	wheres = append(wheres, tDocument.ID.EQ(mysql.Int64(q.DocumentID)))

	columns := dbutils.Columns{
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

func (s *Store) UpdateDocumentOwner(
	ctx context.Context,
	tx qrm.DB,
	documentID int64,
	userInfo *userinfo.UserInfo,
	newOwner *usershort.UserShort,
) error {
	if userInfo == nil {
		userInfo = &userinfo.UserInfo{}
	}

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

func buildListOrder(
	sort *resourcesdatabase.Sort,
	tDocumentShort *table.FivenetDocumentsTable,
) []mysql.OrderByClause {
	orderBys := []mysql.OrderByClause{}
	if sort != nil && len(sort.GetColumns()) > 0 {
		for _, sc := range sort.GetColumns() {
			var column mysql.Column
			switch sc.GetId() {
			case "title":
				column = tDocumentShort.Title
			case "createdAt":
				fallthrough
			default:
				column = tDocumentShort.CreatedAt
			}

			if sc.GetDesc() {
				orderBys = append(orderBys, column.DESC(), tDocumentShort.UpdatedAt.DESC())
			} else {
				orderBys = append(orderBys, column.ASC(), tDocumentShort.UpdatedAt.DESC())
			}
		}
	} else {
		orderBys = append(orderBys, tDocumentShort.UpdatedAt.DESC())
	}

	return orderBys
}
