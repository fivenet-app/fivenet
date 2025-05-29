package documents

import (
	context "context"
	"errors"
	"fmt"
	"slices"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	permsdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/access"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2025/services/documents/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	CommentsDefaultPageSize = 8
	CommentsMaxLength       = 2048
)

var tDComments = table.FivenetDocumentsComments

func (s *Server) GetComments(ctx context.Context, req *pbdocuments.GetCommentsRequest) (*pbdocuments.GetCommentsResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.documents.id", int64(req.DocumentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && !userInfo.Superuser {
		return nil, errorsdocuments.ErrCommentViewDenied
	}

	tDComments := tDComments.AS("comment")
	var condition jet.BoolExpression
	if userInfo.Superuser {
		condition = jet.AND(
			tDComments.DocumentID.EQ(jet.Uint64(req.DocumentId)),
		)
	} else {
		condition = jet.AND(
			tDComments.DocumentID.EQ(jet.Uint64(req.DocumentId)),
			tDComments.DeletedAt.IS_NULL(),
		)
	}

	countStmt := tDComments.
		SELECT(
			jet.COUNT(tDComments.ID).AS("data_count.total"),
		).
		FROM(
			tDComments,
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.Total, CommentsDefaultPageSize)
	resp := &pbdocuments.GetCommentsResponse{
		Pagination: pag,
		Comments:   []*documents.Comment{},
	}
	if count.Total <= 0 {
		return resp, nil
	}

	tCreator := tables.User().AS("creator")
	tAvatar := table.FivenetFiles.AS("avatar")

	columns := jet.ProjectionList{
		tDComments.ID,
		tDComments.DocumentID,
		tDComments.CreatedAt,
		tDComments.UpdatedAt,
		tDComments.Comment.AS("comment.content"),
		tDComments.CreatorID,
		tCreator.ID,
		tCreator.Job,
		tCreator.JobGrade,
		tCreator.Firstname,
		tCreator.Lastname,
		tUserProps.AvatarFileID.AS("creator.avatar_file_id"),
		tAvatar.FilePath.AS("creator.avatar"),
	}
	if userInfo.Superuser {
		columns = append(columns, tDComments.DeletedAt)
	}

	stmt := tDComments.
		SELECT(
			columns[0],
			columns[1:]...,
		).
		FROM(
			tDComments.
				LEFT_JOIN(tCreator,
					tDComments.CreatorID.EQ(tCreator.ID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tCreator.ID),
				).
				LEFT_JOIN(tAvatar,
					tAvatar.ID.EQ(tUserProps.AvatarFileID),
				),
		).
		WHERE(condition).
		OFFSET(
			req.Pagination.Offset,
		).
		ORDER_BY(
			tDComments.CreatedAt.DESC(),
		).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Comments); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	resp.Pagination.Update(len(resp.Comments))

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.Comments {
		if resp.Comments[i].Creator != nil {
			jobInfoFn(resp.Comments[i].Creator)
		}
	}

	return resp, nil
}

func (s *Server) PostComment(ctx context.Context, req *pbdocuments.PostCommentRequest) (*pbdocuments.PostCommentResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.documents.id", int64(req.Comment.DocumentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "PostComment",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.Comment.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_COMMENT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && !userInfo.Superuser {
		return nil, errorsdocuments.ErrCommentPostDenied
	}

	if len(*req.Comment.Content.RawContent) > CommentsMaxLength {
		return nil, errorsdocuments.ErrCommentPostDenied
	}

	stmt := tDComments.
		INSERT(
			tDComments.DocumentID,
			tDComments.Comment,
			tDComments.CreatorID,
			tDComments.CreatorJob,
		).
		VALUES(
			req.Comment.DocumentId,
			req.Comment.Content,
			userInfo.UserId,
			userInfo.Job,
		)

	result, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if _, err := addDocumentActivity(ctx, s.db, &documents.DocActivity{
		DocumentId:   req.Comment.DocumentId,
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_COMMENT_ADDED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.Job,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := s.notifyUsersNewComment(ctx, req.Comment.DocumentId, userInfo.UserId); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_CREATED

	comment, err := s.getComment(ctx, uint64(lastId), userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	comment.CreatorJob = ""

	return &pbdocuments.PostCommentResponse{
		Comment: comment,
	}, nil
}

func (s *Server) EditComment(ctx context.Context, req *pbdocuments.EditCommentRequest) (*pbdocuments.EditCommentResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.documents.id", int64(req.Comment.DocumentId)))
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.documents.comment_id", int64(req.Comment.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "EditComment",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.Comment.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_COMMENT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && !userInfo.Superuser {
		return nil, errorsdocuments.ErrCommentEditDenied
	}

	comment, err := s.getComment(ctx, req.Comment.Id, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !userInfo.Superuser && *comment.CreatorId != userInfo.UserId {
		return nil, errorsdocuments.ErrCommentEditDenied
	}

	if len(*req.Comment.Content.RawContent) > CommentsMaxLength {
		return nil, errorsdocuments.ErrCommentPostDenied
	}

	stmt := tDComments.
		UPDATE(
			tDComments.Comment,
		).
		SET(
			tDComments.Comment.SET(jet.String(*req.Comment.Content.RawContent)),
		).
		WHERE(jet.AND(
			tDComments.ID.EQ(jet.Uint64(req.Comment.Id)),
			tDComments.DeletedAt.IS_NULL(),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	comment.Content = req.Comment.Content

	if _, err := addDocumentActivity(ctx, s.db, &documents.DocActivity{
		DocumentId:   req.Comment.DocumentId,
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_COMMENT_UPDATED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.Job,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbdocuments.EditCommentResponse{
		Comment: comment,
	}, nil
}

func (s *Server) getComment(ctx context.Context, id uint64, userInfo *userinfo.UserInfo) (*documents.Comment, error) {
	tDComments := tDComments.AS("comment")
	tCreator := tables.User().AS("creator")
	tAvatar := table.FivenetFiles.AS("avatar")

	stmt := tDComments.
		SELECT(
			tDComments.ID,
			tDComments.CreatedAt,
			tDComments.UpdatedAt,
			tDComments.DocumentID,
			tDComments.Comment.AS("comment.content"),
			tDComments.CreatorID,
			tDComments.CreatorJob,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tUserProps.AvatarFileID.AS("colleague.avatar_file_id"),
			tAvatar.FilePath.AS("colleague.avatar"),
		).
		FROM(
			tDComments.
				LEFT_JOIN(tCreator,
					tDComments.CreatorID.EQ(tCreator.ID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tCreator.ID),
				).
				LEFT_JOIN(tAvatar,
					tAvatar.ID.EQ(tUserProps.AvatarFileID),
				),
		).
		WHERE(
			tDComments.ID.EQ(jet.Uint64(id)),
		).
		LIMIT(1)

	comment := &documents.Comment{}
	if err := stmt.QueryContext(ctx, s.db, comment); err != nil {
		return nil, err
	}

	if comment.Creator != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, comment.Creator)
	}

	return comment, nil
}

func (s *Server) DeleteComment(ctx context.Context, req *pbdocuments.DeleteCommentRequest) (*pbdocuments.DeleteCommentResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.documents.comment_id", int64(req.CommentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "DeleteComment",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	comment, err := s.getComment(ctx, req.CommentId, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if comment.CreatorJob == "" {
		comment.CreatorJob = userInfo.Job
	}

	check, err := s.access.CanUserAccessTarget(ctx, comment.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_COMMENT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && !userInfo.Superuser {
		return nil, errorsdocuments.ErrCommentDeleteDenied
	}

	// Field Permission Check
	fields, err := s.ps.AttrStringList(userInfo, permsdocuments.DocumentsServicePerm, permsdocuments.DocumentsServiceDeleteCommentPerm, permsdocuments.DocumentsServiceDeleteCommentAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !access.CheckIfHasAccess(fields, userInfo, comment.CreatorJob, comment.Creator) {
		return nil, errorsdocuments.ErrCommentDeleteDenied
	}

	stmt := tDComments.
		UPDATE(
			tDComments.DeletedAt,
		).
		SET(
			tDComments.DeletedAt.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(jet.AND(
			tDComments.ID.EQ(jet.Uint64(req.CommentId)),
			tDComments.DeletedAt.IS_NULL(),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if _, err := addDocumentActivity(ctx, s.db, &documents.DocActivity{
		DocumentId:   uint64(comment.DocumentId),
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_COMMENT_DELETED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.Job,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return &pbdocuments.DeleteCommentResponse{}, nil
}

func (s *Server) notifyUsersNewComment(ctx context.Context, documentId uint64, sourceUserId int32) error {
	userInfo, err := s.ui.GetUserInfoWithoutAccountId(ctx, sourceUserId)
	if err != nil {
		return err
	}

	doc, err := s.getDocument(ctx, tDocument.ID.EQ(jet.Uint64(documentId)), userInfo, false)
	if err != nil {
		return err
	}
	if doc == nil || doc.DeletedAt != nil {
		return nil
	}

	tCreator := tables.User().AS("creator")

	// Get the last 3 commentors to send them a notification
	stmt := tDComments.
		SELECT(
			jet.DISTINCT(tDComments.CreatorID),
		).
		FROM(
			tDComments.
				LEFT_JOIN(tCreator,
					tDComments.CreatorID.EQ(tCreator.ID),
				),
		).
		WHERE(jet.AND(
			tDComments.DocumentID.EQ(jet.Uint64(doc.Id)),
			tDComments.CreatorID.NOT_EQ(jet.Int32(sourceUserId)),
		)).
		GROUP_BY(tDComments.CreatorID).
		ORDER_BY(
			tDComments.CreatedAt.DESC(),
		).
		LIMIT(3)

	var targetUserIds []int32
	if err := stmt.QueryContext(ctx, s.db, &targetUserIds); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	// If we have a document creator, make sure to inform the creator if necessary
	if doc.CreatorId != nil && sourceUserId != *doc.CreatorId && !slices.Contains(targetUserIds, *doc.CreatorId) {
		userInfo, err := s.ui.GetUserInfoWithoutAccountId(ctx, sourceUserId)
		if err != nil {
			return err
		}

		check, err := s.access.CanUserAccessTarget(ctx, doc.Id, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
		if err != nil {
			return err
		}
		if check {
			targetUserIds = append(targetUserIds, *doc.CreatorId)
		}
	}

	for _, targetUserId := range targetUserIds {
		// Don't send notifications to source user
		if targetUserId == sourceUserId {
			continue
		}

		// Make sure user has access to document
		userInfo, err := s.ui.GetUserInfoWithoutAccountId(ctx, sourceUserId)
		if err != nil {
			return err
		}

		check, err := s.access.CanUserAccessTarget(ctx, doc.Id, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
		if err != nil {
			return err
		}
		if !check {
			continue
		}

		not := &notifications.Notification{
			UserId: targetUserId,
			Title: &common.TranslateItem{
				Key: "notifications.documents.document_comment_added.title",
			},
			Content: &common.TranslateItem{
				Key:        "notifications.documents.document_comment_added.content",
				Parameters: map[string]string{"title": doc.Title},
			},
			Type:     notifications.NotificationType_NOTIFICATION_TYPE_INFO,
			Category: notifications.NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT,
			Data: &notifications.Data{
				Link: &notifications.Link{
					To: fmt.Sprintf("/documents/%d#comments", doc.Id),
				},
				CausedBy: &users.UserShort{
					UserId: sourceUserId,
				},
			},
		}
		if err := s.notif.NotifyUser(ctx, not); err != nil {
			return err
		}
	}

	return nil
}
