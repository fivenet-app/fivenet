package mailerstore

import (
	"context"
	"database/sql"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/content"
	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	maileremails "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/emails"
	mailermessages "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/messages"
	mailersettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/settings"
	mailertemplates "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/templates"
	mailerthreads "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/threads"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type IStore interface {
	CountThreads(ctx context.Context, db qrm.DB, in ThreadListQuery) (int64, error)
	ListThreads(ctx context.Context, db qrm.DB, in ThreadListQuery) ([]*mailerthreads.Thread, error)
	GetThread(
		ctx context.Context,
		db qrm.DB,
		threadID int64,
		emailID int64,
		includeDeleted bool,
	) (*mailerthreads.Thread, error)
	UpdateThreadTime(ctx context.Context, db qrm.DB, threadID int64) error
	AddThreadRecipients(
		ctx context.Context,
		db qrm.DB,
		threadID int64,
		recipients []*mailerthreads.ThreadRecipientEmail,
	) error
	ListThreadRecipients(
		ctx context.Context,
		db qrm.DB,
		threadID int64,
		includeDeleted bool,
	) ([]*mailerthreads.ThreadRecipientEmail, error)
	GetThreadState(
		ctx context.Context,
		db qrm.DB,
		threadID int64,
		emailID int64,
	) (*mailerthreads.ThreadState, error)
	SetThreadState(ctx context.Context, db qrm.DB, state *mailerthreads.ThreadState) error
	SetUnreadState(
		ctx context.Context,
		db qrm.DB,
		threadID int64,
		senderID int64,
		emailIDs []int64,
	) error
	CountThreadMessages(
		ctx context.Context,
		db qrm.DB,
		threadID int64,
		includeDeleted bool,
	) (int64, error)
	ListThreadMessages(
		ctx context.Context,
		db qrm.DB,
		in MessageListQuery,
		includeDeleted bool,
	) ([]*mailermessages.Message, error)
	DeleteThread(
		ctx context.Context,
		q qrm.DB,
		threadID int64,
		deletedAt *timestamp.Timestamp,
	) error
	GetMessage(
		ctx context.Context,
		db qrm.DB,
		messageID int64,
		includeDeleted bool,
	) (*mailermessages.Message, error)
	CreateMessage(ctx context.Context, db qrm.DB, msg *mailermessages.Message) (int64, error)
	DeleteMessage(
		ctx context.Context,
		q qrm.DB,
		threadID int64,
		messageID int64,
		deletedAt *timestamp.Timestamp,
	) error
	CountEmails(ctx context.Context, db qrm.DB, condition mysql.BoolExpression) (int64, error)
	ListUserEmails(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		pag *database.PaginationRequest,
		includeDisabled bool,
		includeDeleted bool,
	) ([]*maileremails.Email, error)
	ListEmails(
		ctx context.Context,
		db qrm.DB,
		userInfo *userinfo.UserInfo,
		pag *database.PaginationRequest,
		all bool,
	) (*database.PaginationResponse, []*maileremails.Email, error)
	GetEmailByCondition(
		ctx context.Context,
		db qrm.DB,
		condition mysql.BoolExpression,
	) (*maileremails.Email, error)
	GetEmail(
		ctx context.Context,
		db qrm.DB,
		emailID int64,
		includeDeleted bool,
	) (*maileremails.Email, error)
	CreateEmail(
		ctx context.Context,
		q qrm.DB,
		email *maileremails.Email,
		creatorID int32,
	) (int64, error)
	DeleteEmail(
		ctx context.Context,
		q qrm.DB,
		emailID int64,
		deletedAt *timestamp.Timestamp,
	) error
	GetUserShort(ctx context.Context, db qrm.DB, userID int32) (*usershort.UserShort, error)
	ListRecipientsByEmails(
		ctx context.Context,
		db qrm.DB,
		recipients []string,
	) ([]*mailerthreads.ThreadRecipientEmail, error)
	GetEmailSettings(
		ctx context.Context,
		db qrm.DB,
		emailID int64,
	) (*mailersettings.EmailSettings, error)
	UpsertEmailSettingsSignature(
		ctx context.Context,
		db qrm.DB,
		emailID int64,
		signature *content.Content,
	) error
	AddBlockedEmails(ctx context.Context, db qrm.DB, emailID int64, blockedEmails []string) error
	DeleteBlockedEmails(ctx context.Context, db qrm.DB, emailID int64, blockedEmails []string) error
	ListTemplates(
		ctx context.Context,
		db qrm.DB,
		emailID int64,
		limit int64,
	) ([]*mailertemplates.Template, error)
	GetTemplate(
		ctx context.Context,
		db qrm.DB,
		id int64,
		emailID *int64,
	) (*mailertemplates.Template, error)
	CountTemplatesByCreatorJob(ctx context.Context, db qrm.DB, job string) (int64, error)
	CreateTemplate(
		ctx context.Context,
		db qrm.DB,
		template *mailertemplates.Template,
		creatorID int32,
	) (int64, error)
	UpdateTemplate(ctx context.Context, db qrm.DB, template *mailertemplates.Template) error
	DeleteTemplate(ctx context.Context, db qrm.DB, id int64) error
}

type Store struct {
	db            *sql.DB
	subjectAccess *access.SubjectObjectAccess
}

func New(db *sql.DB) IStore {
	return &Store{
		db:            db,
		subjectAccess: access.NewMailerEmailsSubjectObjectAccess(db),
	}
}
