package documentsstore

import (
	"context"
	"database/sql"

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
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
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
	GetDocumentMeta(
		ctx context.Context,
		db qrm.DB,
		documentID int64,
	) (*resourcesdocuments.DocumentMeta, error)
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
	UpsertWorkflowState(
		ctx context.Context,
		tx qrm.DB,
		documentID int64,
		workflow *documentsworkflow.Workflow,
	) error
	UpsertWorkflowUserState(
		ctx context.Context,
		tx qrm.DB,
		state *documentsworkflow.WorkflowUserState,
	) error
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
	db                 *sql.DB
	subjectAccess      *access.SubjectObjectAccess
	subjectResolver    *access.SubjectResolver
	templateAccess     *access.SubjectObjectAccess
	userDocumentSorter *resourcesdatabase.SorterBuilder
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
		userDocumentSorter: resourcesdatabase.New(
			resourcesdatabase.SpecMap{
				"createdAt": resourcesdatabase.Column{
					Col: table.FivenetDocuments.AS("document").CreatedAt,
				},
			},
			[]mysql.OrderByClause{table.FivenetDocuments.AS("document").CreatedAt.DESC()},
			nil,
			"createdAt",
			3,
		),
	}
}
