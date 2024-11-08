package docstore

import (
	context "context"
	"errors"
	"fmt"
	"slices"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common"
	database "github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	errorsdocstore "github.com/fivenet-app/fivenet/gen/go/proto/services/docstore/errors"
	permsdocstore "github.com/fivenet-app/fivenet/gen/go/proto/services/docstore/perms"
	"github.com/fivenet-app/fivenet/pkg/access"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	CommentsDefaultPageSize = 8
)

var tDComments = table.FivenetDocumentsComments

func (s *Server) GetComments(ctx context.Context, req *GetCommentsRequest) (*GetCommentsResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.id", int64(req.DocumentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	if !check {
		return nil, errorsdocstore.ErrCommentViewDenied
	}

	tDComments := tDComments.AS("comment")
	var condition jet.BoolExpression
	if userInfo.SuperUser {
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
			jet.COUNT(tDComments.ID).AS("datacount.totalcount"),
		).
		FROM(
			tDComments,
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, CommentsDefaultPageSize)
	resp := &GetCommentsResponse{
		Pagination: pag,
		Comments:   []*documents.Comment{},
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	columns := jet.ProjectionList{
		tDComments.ID,
		tDComments.DocumentID,
		tDComments.CreatedAt,
		tDComments.UpdatedAt,
		tDComments.Comment,
		tDComments.CreatorID,
		tCreator.ID,
		tCreator.Job,
		tCreator.JobGrade,
		tCreator.Firstname,
		tCreator.Lastname,
		tUserProps.Avatar.AS("creator.avatar"),
	}
	if userInfo.SuperUser {
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
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	resp.Pagination.Update(len(resp.Comments))

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := 0; i < len(resp.Comments); i++ {
		if resp.Comments[i].Creator != nil {
			jobInfoFn(resp.Comments[i].Creator)
		}
	}

	return resp, nil
}

func (s *Server) PostComment(ctx context.Context, req *PostCommentRequest) (*PostCommentResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.id", int64(req.Comment.DocumentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "PostComment",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.Comment.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_COMMENT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsdocstore.ErrCommentPostDenied
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
			req.Comment.Comment,
			userInfo.UserId,
			userInfo.Job,
		)

	result, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	if _, err := s.addDocumentActivity(ctx, s.db, &documents.DocActivity{
		DocumentId:   req.Comment.DocumentId,
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_COMMENT_ADDED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.Job,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	if err := s.notifyUsersNewComment(ctx, req.Comment.DocumentId, userInfo.UserId); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	comment, err := s.getComment(ctx, uint64(lastId), userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	comment.CreatorJob = ""

	return &PostCommentResponse{
		Comment: comment,
	}, nil
}

func (s *Server) EditComment(ctx context.Context, req *EditCommentRequest) (*EditCommentResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.id", int64(req.Comment.DocumentId)))
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.comment_id", int64(req.Comment.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "EditComment",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.Comment.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_COMMENT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsdocstore.ErrCommentEditDenied
	}

	comment, err := s.getComment(ctx, req.Comment.Id, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	if !userInfo.SuperUser && *comment.CreatorId != userInfo.UserId {
		return nil, errorsdocstore.ErrCommentEditDenied
	}

	stmt := tDComments.
		UPDATE(
			tDComments.Comment,
		).
		SET(
			tDComments.Comment.SET(jet.String(req.Comment.Comment)),
		).
		WHERE(
			jet.AND(
				tDComments.ID.EQ(jet.Uint64(req.Comment.Id)),
				tDComments.DeletedAt.IS_NULL(),
			),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	comment.Comment = req.Comment.Comment

	if _, err := s.addDocumentActivity(ctx, s.db, &documents.DocActivity{
		DocumentId:   req.Comment.DocumentId,
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_COMMENT_UPDATED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.Job,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &EditCommentResponse{
		Comment: comment,
	}, nil
}

func (s *Server) getComment(ctx context.Context, id uint64, userInfo *userinfo.UserInfo) (*documents.Comment, error) {
	comment := &documents.Comment{}

	tDComments := tDComments.AS("comment")
	stmt := tDComments.
		SELECT(
			tDComments.ID,
			tDComments.CreatedAt,
			tDComments.UpdatedAt,
			tDComments.DocumentID,
			tDComments.Comment,
			tDComments.CreatorID,
			tDComments.CreatorJob,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tUserProps.Avatar.AS("creator.avatar"),
		).
		FROM(
			tDComments.
				LEFT_JOIN(tCreator,
					tDComments.CreatorID.EQ(tCreator.ID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tCreator.ID),
				),
		).
		WHERE(
			tDComments.ID.EQ(jet.Uint64(id)),
		).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, s.db, comment); err != nil {
		return nil, err
	}

	if comment.Creator != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, comment.Creator)
	}

	return comment, nil
}

func (s *Server) DeleteComment(ctx context.Context, req *DeleteCommentRequest) (*DeleteCommentResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.comment_id", int64(req.CommentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "DeleteComment",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	comment, err := s.getComment(ctx, req.CommentId, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	if comment.CreatorJob == "" {
		comment.CreatorJob = userInfo.Job
	}

	check, err := s.access.CanUserAccessTarget(ctx, comment.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_COMMENT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsdocstore.ErrCommentDeleteDenied
	}

	// Field Permission Check
	fieldsAttr, err := s.ps.Attr(userInfo, permsdocstore.DocStoreServicePerm, permsdocstore.DocStoreServiceDeleteCommentPerm, permsdocstore.DocStoreServiceDeleteCommentAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}
	if !access.CheckIfHasAccess(fields, userInfo, comment.CreatorJob, comment.Creator) {
		return nil, errorsdocstore.ErrCommentDeleteDenied
	}

	stmt := tDComments.
		UPDATE(
			tDComments.DeletedAt,
		).
		SET(
			tDComments.DeletedAt.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(
			jet.AND(
				tDComments.ID.EQ(jet.Uint64(req.CommentId)),
				tDComments.DeletedAt.IS_NULL(),
			),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	if _, err := s.addDocumentActivity(ctx, s.db, &documents.DocActivity{
		DocumentId:   uint64(comment.DocumentId),
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_COMMENT_DELETED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.Job,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeleteCommentResponse{}, nil
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

	// Get the last 3 commentors to send them a notification
	stmt := tDComments.
		SELECT(
			tDComments.CreatorID,
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
				Key: "notifications.document_comment_added.title",
			},
			Content: &common.TranslateItem{
				Key:        "notifications.document_comment_added.content",
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
