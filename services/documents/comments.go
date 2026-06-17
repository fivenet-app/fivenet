package documents

import (
	context "context"
	"errors"
	"fmt"
	"slices"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	documentsactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/activity"
	documentscomment "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/comment"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	pbdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents"
	permsdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2026/services/documents/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

const (
	CommentsDefaultPageSize = 8
	CommentsMaxLength       = 2048
)

func (s *Server) GetComments(
	ctx context.Context,
	req *pbdocuments.GetCommentsRequest,
) (*pbdocuments.GetCommentsResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.canUserAccessDocument(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsdocuments.ErrCommentViewDenied
	}

	count, err := s.store.CountComments(ctx, s.db, req.GetDocumentId(), userInfo.GetSuperuser())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(int64(count), CommentsDefaultPageSize)
	resp := &pbdocuments.GetCommentsResponse{
		Pagination: pag,
		Comments:   []*documentscomment.Comment{},
	}
	if count <= 0 {
		return resp, nil
	}
	resp.Comments, err = s.store.ListComments(
		ctx,
		req.GetDocumentId(),
		userInfo,
		req.GetPagination().GetOffset(),
		limit,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetComments() {
		if resp.GetComments()[i].GetCreator() != nil {
			jobInfoFn(resp.GetComments()[i].GetCreator())
		}
	}

	return resp, nil
}

func (s *Server) PostComment(
	ctx context.Context,
	req *pbdocuments.PostCommentRequest,
) (*pbdocuments.PostCommentResponse, error) {
	logging.InjectFields(
		ctx,
		logging.Fields{"fivenet.documents.id", req.GetComment().GetDocumentId()},
	)

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.canUserAccessDocument(
		ctx,
		req.GetComment().GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_COMMENT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsdocuments.ErrCommentPostDenied
	}

	// Check comment length
	extracted := req.GetComment().GetContent().Extract()
	if len(extracted.GetText()) > CommentsMaxLength {
		return nil, errorsdocuments.ErrCommentPostDenied
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	lastId, err := s.store.CreateComment(ctx, tx, req.GetComment(), userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := s.notifyUsersNewComment(
		ctx,
		tx,
		req.GetComment().GetDocumentId(),
		userInfo.GetUserId(),
	); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := s.store.UpdateCommentsCount(ctx, tx, req.GetComment().GetDocumentId()); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

	comment, err := s.store.GetComment(ctx, lastId, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	comment.CreatorJob = ""

	return &pbdocuments.PostCommentResponse{
		Comment: comment,
	}, nil
}

func (s *Server) EditComment(
	ctx context.Context,
	req *pbdocuments.EditCommentRequest,
) (*pbdocuments.EditCommentResponse, error) {
	logging.InjectFields(
		ctx,
		logging.Fields{
			"fivenet.documents.id", req.GetComment().GetDocumentId(),
			"fivenet.documents.comment_id", req.GetComment().GetId(),
		},
	)

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.canUserAccessDocument(
		ctx,
		req.GetComment().GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_COMMENT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsdocuments.ErrCommentEditDenied
	}

	comment, err := s.store.GetComment(ctx, req.GetComment().GetId(), userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !userInfo.GetSuperuser() && comment.GetCreatorId() != userInfo.GetUserId() {
		return nil, errorsdocuments.ErrCommentEditDenied
	}

	// Check comment length
	extracted := req.GetComment().GetContent().Extract()
	if len(extracted.GetText()) > CommentsMaxLength {
		return nil, errorsdocuments.ErrCommentPostDenied
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	if err := s.store.UpdateComment(ctx, tx, req.GetComment(), userInfo); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	comment.Content = req.GetComment().GetContent()

	if err := s.store.UpdateCommentsCount(ctx, tx, comment.GetDocumentId()); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return &pbdocuments.EditCommentResponse{
		Comment: comment,
	}, nil
}

func (s *Server) DeleteComment(
	ctx context.Context,
	req *pbdocuments.DeleteCommentRequest,
) (*pbdocuments.DeleteCommentResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.GetCommentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	comment, err := s.store.GetComment(ctx, req.GetCommentId(), userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if comment.GetCreatorJob() == "" {
		comment.CreatorJob = userInfo.GetJob()
	}

	check, err := s.canUserAccessDocument(
		ctx,
		comment.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_COMMENT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsdocuments.ErrCommentDeleteDenied
	}

	// Field Permission Check
	fields, err := permsdocuments.CommentsService.DeleteComment.AccessTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !access.CheckIfHasOwnJobAccess(
		fields.StringList(),
		userInfo,
		comment.GetCreatorJob(),
		comment.GetCreator(),
	) {
		return nil, errorsdocuments.ErrCommentDeleteDenied
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	activityType := documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_COMMENT_DELETED
	var deletedAtTime *timestamp.Timestamp
	if comment.GetDeletedAt() == nil || !userInfo.GetSuperuser() {
		deletedAtTime = timestamp.Now()
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	} else {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_RESTORED)
		activityType = documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_COMMENT_RESTORED
	}

	if err := s.store.DeleteComment(
		ctx,
		tx,
		comment,
		userInfo,
		deletedAtTime,
		activityType,
	); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := s.store.UpdateCommentsCount(ctx, tx, comment.GetDocumentId()); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.DeleteCommentResponse{}, nil
}

func (s *Server) notifyUsersNewComment(
	ctx context.Context,
	tx qrm.DB,
	documentId int64,
	sourceUserId int32,
) error {
	userInfo, err := s.ui.GetUserInfo(ctx, sourceUserId)
	if err != nil {
		return err
	}

	doc, err := s.getDocument(ctx, tDocument.ID.EQ(mysql.Int64(documentId)), userInfo, false)
	if err != nil {
		return err
	}
	if doc == nil || doc.GetDeletedAt() != nil {
		return nil
	}

	lastPerCreator := table.FivenetDocumentsComments.
		SELECT(
			table.FivenetDocumentsComments.CreatorID.AS("creator_id"),
			mysql.MAX(table.FivenetDocumentsComments.CreatedAt).AS("last_at"),
		).
		FROM(table.FivenetDocumentsComments).
		WHERE(mysql.AND(
			table.FivenetDocumentsComments.DocumentID.EQ(mysql.Int64(documentId)),
			table.FivenetDocumentsComments.CreatorID.NOT_EQ(mysql.Int32(sourceUserId)),
		)).
		GROUP_BY(table.FivenetDocumentsComments.CreatorID).
		AsTable("dc")

	// Get the last 3 commentors to send them a notification
	stmt := lastPerCreator.
		SELECT(
			mysql.RawInt("creator_id"),
		).
		FROM(lastPerCreator).
		ORDER_BY(
			mysql.RawTimestamp("last_at").DESC(),
		).
		LIMIT(3)

	var targetUserIds []int32
	if err := stmt.QueryContext(ctx, tx, &targetUserIds); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	// If we have a document creator, make sure to inform the creator if necessary
	if doc.CreatorId != nil && sourceUserId != doc.GetCreatorId() &&
		!slices.Contains(targetUserIds, doc.GetCreatorId()) {
		userInfo, err := s.ui.GetUserInfo(ctx, sourceUserId)
		if err != nil {
			return err
		}

		check, err := s.canUserAccessDocument(
			ctx,
			doc.GetId(),
			userInfo,
			documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW,
		)
		if err != nil {
			return err
		}
		if check {
			targetUserIds = append(targetUserIds, doc.GetCreatorId())
		}
	}

	for _, targetUserId := range targetUserIds {
		// Don't send notifications to source user
		if targetUserId == sourceUserId {
			continue
		}

		// Make sure user has access to document
		userInfo, err := s.ui.GetUserInfo(ctx, sourceUserId)
		if err != nil {
			return err
		}

		check, err := s.canUserAccessDocument(
			ctx,
			doc.GetId(),
			userInfo,
			documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW,
		)
		if err != nil {
			return err
		}
		if !check {
			continue
		}

		not := &notifications.Notification{
			UserId: targetUserId,
			Title: &common.I18NItem{
				Key: "notifications.documents.document_comment_added.title",
			},
			Content: &common.I18NItem{
				Key:        "notifications.documents.document_comment_added.content",
				Parameters: map[string]string{"title": doc.GetTitle()},
			},
			Type:     notifications.NotificationType_NOTIFICATION_TYPE_INFO,
			Category: notifications.NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT,
			Data: &notifications.Data{
				Link: &notifications.Link{
					To: fmt.Sprintf("/documents/%d#comments", doc.GetId()),
				},
				CausedBy: &usershort.UserShort{
					UserId: sourceUserId,
				},
			},
		}
		if err := s.notifi.NotifyUser(ctx, not); err != nil {
			return err
		}
	}

	return nil
}
