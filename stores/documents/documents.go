package documentsstore

import (
	"context"
	"database/sql"
	"errors"

	resourcesaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/access"
	resourcesdatabase "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	resourcesdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents"
	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	documentsactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/activity"
	documentsapproval "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/approval"
	documentscategory "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/category"
	documentscomment "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/comment"
	documentspins "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/pins"
	documentsreferences "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/references"
	documentsrelations "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/relations"
	documentsrequests "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/requests"
	documentsstamps "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/stamps"
	documentstemplates "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/templates"
	documentsworkflow "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/workflow"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	pbdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2026/services/documents/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type IStore interface {
	List(ctx context.Context, q ListQuery) ([]*resourcesdocuments.DocumentShort, error)
	Get(ctx context.Context, q GetQuery) (*resourcesdocuments.Document, error)
	ListTemplates(
		ctx context.Context,
		userInfo *userinfo.UserInfo,
	) ([]*documentstemplates.TemplateShort, error)
	GetTemplate(ctx context.Context, templateID int64) (*documentstemplates.Template, error)
	GetDocumentAccess(
		ctx context.Context,
		documentID int64,
	) (*documentsaccess.DocumentAccess, error)
	GetDocumentMeta(ctx context.Context, db qrm.DB, documentID int64) (*resourcesdocuments.DocumentMeta, error)
	GetDocumentRequiredAccess(
		ctx context.Context,
		documentID int64,
		userInfo *userinfo.UserInfo,
	) (*resourcesaccess.Access, error)
	GetDocumentPin(
		ctx context.Context,
		documentID int64,
		userInfo *userinfo.UserInfo,
	) (*documentspins.DocumentPin, error)
	ListDocumentPins(
		ctx context.Context,
		q ListDocumentPinsQuery,
	) (*resourcesdatabase.PaginationResponse, []*resourcesdocuments.DocumentShort, error)
	CreateDocumentPin(
		ctx context.Context,
		tx qrm.DB,
		documentID int64,
		userInfo *userinfo.UserInfo,
		personal bool,
	) error
	DeleteDocumentPin(
		ctx context.Context,
		tx qrm.DB,
		documentID int64,
		userInfo *userinfo.UserInfo,
		personal bool,
	) error
	ListCategories(
		ctx context.Context,
		userInfo *userinfo.UserInfo,
	) ([]*documentscategory.Category, error)
	GetCategory(
		ctx context.Context,
		id int64,
		userInfo *userinfo.UserInfo,
	) (*documentscategory.Category, error)
	CreateCategory(
		ctx context.Context,
		tx qrm.DB,
		category *documentscategory.Category,
		userInfo *userinfo.UserInfo,
	) (int64, error)
	UpdateCategory(
		ctx context.Context,
		tx qrm.DB,
		category *documentscategory.Category,
		userInfo *userinfo.UserInfo,
	) error
	DeleteCategory(
		ctx context.Context,
		tx qrm.DB,
		id int64,
		userInfo *userinfo.UserInfo,
		deletedAt *timestamp.Timestamp,
	) error
	UpsertWorkflowState(ctx context.Context, tx qrm.DB, documentID int64, workflow *documentsworkflow.Workflow) error
	UpsertWorkflowUserState(ctx context.Context, tx qrm.DB, state *documentsworkflow.WorkflowUserState) error
	DeleteWorkflowUserState(ctx context.Context, tx qrm.DB, documentID int64, userID int32) error
	ListDocumentActivity(
		ctx context.Context,
		q ListDocumentActivityQuery,
	) (resourcesdatabase.DataCount, []*documentsactivity.DocActivity, error)
	ListDocumentReqs(
		ctx context.Context,
		q ListDocumentReqsQuery,
	) (resourcesdatabase.DataCount, []*documentsrequests.DocRequest, error)
	AddDocumentReq(
		ctx context.Context,
		tx qrm.DB,
		request *documentsrequests.DocRequest,
	) (int64, error)
	UpdateDocumentReq(
		ctx context.Context,
		tx qrm.DB,
		id int64,
		request *documentsrequests.DocRequest,
	) error
	GetDocumentReq(
		ctx context.Context,
		db qrm.DB,
		condition mysql.BoolExpression,
	) (*documentsrequests.DocRequest, error)
	DeleteDocumentReq(ctx context.Context, tx qrm.DB, id int64) error
	ListApprovalTasksInbox(
		ctx context.Context,
		q ListApprovalTasksInboxQuery,
	) (resourcesdatabase.DataCount, []*documentsapproval.ApprovalTask, error)
	GetApprovalPolicy(
		ctx context.Context,
		db qrm.DB,
		condition mysql.BoolExpression,
	) (*documentsapproval.ApprovalPolicy, error)
	CreateApprovalPolicy(
		ctx context.Context,
		tx qrm.DB,
		documentID int64,
		pol *documentsapproval.ApprovalPolicy,
	) error
	GetApprovalTask(
		ctx context.Context,
		db qrm.DB,
		taskID int64,
	) (*documentsapproval.ApprovalTask, error)
	CreateApprovalTasks(
		ctx context.Context,
		tx qrm.DB,
		userInfo *userinfo.UserInfo,
		documentID int64,
		snapDate *timestamp.Timestamp,
		seeds []*pbdocuments.ApprovalTaskSeed,
	) (int32, int32, error)
	DeleteApprovalTasks(
		ctx context.Context,
		tx qrm.DB,
		documentID int64,
		snapDate *timestamp.Timestamp,
		deleteAllPending bool,
		taskIDs []int64,
		pendingCount int32,
	) error
	ExpireApprovalTasks(ctx context.Context, tx qrm.DB) (int64, error)
	ResetApprovalProgress(ctx context.Context, tx qrm.DB, documentID int64) error
	RecomputeApprovalPolicyTx(
		ctx context.Context,
		tx qrm.DB,
		documentID int64,
		snapDate *timestamp.Timestamp,
	) error
	ListApprovals(
		ctx context.Context,
		q ListApprovalsQuery,
	) (resourcesdatabase.DataCount, []*documentsapproval.Approval, error)
	GetApproval(
		ctx context.Context,
		db qrm.DB,
		approvalID int64,
	) (*documentsapproval.Approval, error)
	CheckJobStampCount(ctx context.Context, job string) (int64, error)
	ListUsableStamps(
		ctx context.Context,
		q ListUsableStampsQuery,
	) (*resourcesdatabase.PaginationResponse, []*documentsstamps.Stamp, error)
	GetStamp(ctx context.Context, stampID int64) (*documentsstamps.Stamp, error)
	CreateStamp(ctx context.Context, tx qrm.DB, stamp *documentsstamps.Stamp) (int64, error)
	UpdateStamp(ctx context.Context, tx qrm.DB, stamp *documentsstamps.Stamp) error
	DeleteStamp(ctx context.Context, tx qrm.DB, stampID int64) error
	ListUserDocuments(
		ctx context.Context,
		q ListUserDocumentsQuery,
	) (resourcesdatabase.DataCount, []*documentsrelations.DocumentRelation, error)
	GetDocumentReference(
		ctx context.Context,
		id int64,
	) (*documentsreferences.DocumentReference, error)
	ListDocumentReferences(
		ctx context.Context,
		documentID int64,
	) ([]*documentsreferences.DocumentReference, error)
	CreateDocumentReference(
		ctx context.Context,
		tx qrm.DB,
		ref *documentsreferences.DocumentReference,
	) (int64, error)
	DeleteDocumentReference(ctx context.Context, tx qrm.DB, id int64) error
	GetDocumentRelation(ctx context.Context, id int64) (*documentsrelations.DocumentRelation, error)
	ListDocumentRelations(
		ctx context.Context,
		documentID int64,
	) ([]*documentsrelations.DocumentRelation, error)
	CreateDocumentRelation(
		ctx context.Context,
		tx qrm.DB,
		rel *documentsrelations.DocumentRelation,
	) (int64, bool, error)
	DeleteDocumentRelation(ctx context.Context, tx qrm.DB, id int64) error
	CountComments(
		ctx context.Context,
		tx qrm.DB,
		documentID int64,
		includeDeleted bool,
	) (int32, error)
	ListComments(
		ctx context.Context,
		documentID int64,
		userInfo *userinfo.UserInfo,
		offset int64,
		limit int64,
	) ([]*documentscomment.Comment, error)
	GetComment(
		ctx context.Context,
		id int64,
		userInfo *userinfo.UserInfo,
	) (*documentscomment.Comment, error)
	CreateComment(
		ctx context.Context,
		tx qrm.DB,
		comment *documentscomment.Comment,
		userInfo *userinfo.UserInfo,
	) (int64, error)
	UpdateComment(
		ctx context.Context,
		tx qrm.DB,
		comment *documentscomment.Comment,
		userInfo *userinfo.UserInfo,
	) error
	DeleteComment(
		ctx context.Context,
		tx qrm.DB,
		comment *documentscomment.Comment,
		userInfo *userinfo.UserInfo,
		deletedAt *timestamp.Timestamp,
		activityType documentsactivity.DocActivityType,
	) error
	UpdateCommentsCount(ctx context.Context, tx qrm.DB, documentID int64) error
	UpdateDocumentAccess(
		ctx context.Context,
		tx qrm.DB,
		documentID int64,
		userInfo *userinfo.UserInfo,
		docAccess *documentsaccess.DocumentAccess,
		addActivity bool,
	) error
	CreateTemplate(
		ctx context.Context,
		tx qrm.DB,
		tmpl *documentstemplates.Template,
		creatorJob string,
		categoryID *int64,
	) (int64, error)
	UpdateTemplate(
		ctx context.Context,
		tx qrm.DB,
		tmpl *documentstemplates.Template,
		categoryID *int64,
	) error
	DeleteTemplate(ctx context.Context, tx qrm.DB, templateID int64, creatorJob string) error
	UpdateDocumentOwner(
		ctx context.Context,
		tx qrm.DB,
		documentID int64,
		userInfo *userinfo.UserInfo,
		newOwner *usershort.UserShort,
	) error
}

type Store struct {
	db              *sql.DB
	subjectAccess   *access.SubjectObjectAccess
	subjectResolver *access.SubjectResolver
	templateAccess  *access.SubjectObjectAccess
}

var documentSubjectAccessOptions = access.SubjectAccessOptions{
	BlockedAccess: int32(documentsaccess.AccessLevel_ACCESS_LEVEL_BLOCKED),
	DeniedAccessLevels: []int32{
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_COMMENT),
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_STATUS),
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_ACCESS),
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
	},
}

const documentAccessEntryLimit = 20

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

type ListUserDocumentsQuery struct {
	UserID         int32
	IncludeCreated bool
	Closed         *bool
	Relations      []documentsrelations.DocRelation
	Sort           *resourcesdatabase.Sort
	Pagination     *resourcesdatabase.PaginationRequest
	UserInfo       *userinfo.UserInfo
}

type ListDocumentPinsQuery struct {
	Personal   bool
	Pagination *resourcesdatabase.PaginationRequest
	UserInfo   *userinfo.UserInfo
}

type ListDocumentActivityQuery struct {
	DocumentID    int64
	ActivityTypes []documentsactivity.DocActivityType
	Pagination    *resourcesdatabase.PaginationRequest
	UserInfo      *userinfo.UserInfo
}

type ListDocumentReqsQuery struct {
	DocumentID int64
	Pagination *resourcesdatabase.PaginationRequest
	UserInfo   *userinfo.UserInfo
}

type ListApprovalTasksInboxQuery struct {
	Pagination      *resourcesdatabase.PaginationRequest
	UserInfo        *userinfo.UserInfo
	Statuses        []documentsapproval.ApprovalTaskStatus
	NotAlreadyActed bool
	OnlyDrafts      *bool
}

type ListUsableStampsQuery struct {
	Pagination *resourcesdatabase.PaginationRequest
	UserInfo   *userinfo.UserInfo
}

type ListApprovalsQuery struct {
	DocumentID   int64
	SnapshotDate *timestamp.Timestamp
	Status       documentsapproval.ApprovalStatus
	UserID       int32
	TaskID       int64
	Pagination   *resourcesdatabase.PaginationRequest
	UserInfo     *userinfo.UserInfo
}

func New(db *sql.DB) IStore {
	return &Store{
		db:              db,
		subjectAccess:   access.NewDocumentsSubjectObjectAccess(db),
		subjectResolver: access.NewSubjectResolver(db),
		templateAccess:  access.NewDocumentTemplatesSubjectObjectAccess(db),
	}
}

func (s *Store) List(
	ctx context.Context,
	q ListQuery,
) ([]*resourcesdocuments.DocumentShort, error) {
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

func (s *Store) GetDocumentAccess(
	ctx context.Context,
	documentID int64,
) (*documentsaccess.DocumentAccess, error) {
	return s.subjectAccess.ListTargetAccess(ctx, s.db, documentID, documentSubjectAccessOptions)
}

func (s *Store) GetDocumentMeta(ctx context.Context, db qrm.DB, documentID int64) (*resourcesdocuments.DocumentMeta, error) {
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

	tmpl, err := s.GetTemplate(ctx, doc.GetTemplateId())
	if err != nil {
		return nil, err
	}
	if tmpl.GetContentAccess() == nil || tmpl.GetContentAccess().IsEmpty() {
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

	requiredAccess, err := s.GetDocumentRequiredAccess(ctx, documentID, userInfo)
	if err != nil {
		return err
	}

	fallbackAccess := &resourcesaccess.Access{
		Jobs: []*resourcesaccess.JobAccess{
			{
				Job:          userInfo.GetJob(),
				MinimumGrade: userInfo.GetJobGrade(),
				Access:       int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
			},
		},
	}

	docAccess, err = access.NormalizeAccess(
		docAccess,
		requiredAccess,
		fallbackAccess,
		documentAccessEntryLimit,
	)
	if err != nil {
		if isAccessEntryLimitError(err) {
			return errorsdocuments.ErrDocRequiredAccessTemplate
		}
		return err
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

func isAccessEntryLimitError(err error) bool {
	var limitErr *access.AccessEntryLimitError
	return errors.As(err, &limitErr)
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
