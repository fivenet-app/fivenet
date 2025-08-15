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
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	permsdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/access"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2025/services/documents/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

const (
	CommentsDefaultPageSize = 8
	CommentsMaxLength       = 2048
)

var tDComments = table.FivenetDocumentsComments

func (s *Server) GetComments(
	ctx context.Context,
	req *pbdocuments.GetCommentsRequest,
) (*pbdocuments.GetCommentsResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsdocuments.ErrCommentViewDenied
	}

	tDComments := tDComments.AS("comment")
	var condition jet.BoolExpression
	if userInfo.GetSuperuser() {
		condition = jet.AND(
			tDComments.DocumentID.EQ(jet.Int64(req.GetDocumentId())),
		)
	} else {
		condition = jet.AND(
			tDComments.DocumentID.EQ(jet.Int64(req.GetDocumentId())),
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

	pag, limit := req.GetPagination().GetResponseWithPageSize(count.Total, CommentsDefaultPageSize)
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
	if userInfo.GetSuperuser() {
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
			req.GetPagination().GetOffset(),
		).
		ORDER_BY(
			tDComments.CreatedAt.DESC(),
		).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Comments); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	resp.GetPagination().Update(len(resp.GetComments()))

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

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "PostComment",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetComment().GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_COMMENT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsdocuments.ErrCommentPostDenied
	}

	if len(req.GetComment().GetContent().GetRawContent()) > CommentsMaxLength {
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
			req.GetComment().GetDocumentId(),
			req.GetComment().GetContent(),
			userInfo.GetUserId(),
			userInfo.GetJob(),
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
		DocumentId:   req.GetComment().GetDocumentId(),
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_COMMENT_ADDED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.GetJob(),
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := s.notifyUsersNewComment(ctx, req.GetComment().GetDocumentId(), userInfo.GetUserId()); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_CREATED

	comment, err := s.getComment(ctx, lastId, userInfo)
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
		logging.Fields{"fivenet.documents.id", req.GetComment().GetDocumentId()},
	)
	logging.InjectFields(
		ctx,
		logging.Fields{"fivenet.documents.comment_id", req.GetComment().GetId()},
	)

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "EditComment",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetComment().GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_COMMENT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsdocuments.ErrCommentEditDenied
	}

	comment, err := s.getComment(ctx, req.GetComment().GetId(), userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !userInfo.GetSuperuser() && comment.GetCreatorId() != userInfo.GetUserId() {
		return nil, errorsdocuments.ErrCommentEditDenied
	}

	if len(req.GetComment().GetContent().GetRawContent()) > CommentsMaxLength {
		return nil, errorsdocuments.ErrCommentPostDenied
	}

	stmt := tDComments.
		UPDATE(
			tDComments.Comment,
		).
		SET(
			tDComments.Comment.SET(jet.String(req.GetComment().GetContent().GetRawContent())),
		).
		WHERE(jet.AND(
			tDComments.ID.EQ(jet.Int64(req.GetComment().GetId())),
			tDComments.DeletedAt.IS_NULL(),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	comment.Content = req.GetComment().GetContent()

	if _, err := addDocumentActivity(ctx, s.db, &documents.DocActivity{
		DocumentId:   req.GetComment().GetDocumentId(),
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_COMMENT_UPDATED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.GetJob(),
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbdocuments.EditCommentResponse{
		Comment: comment,
	}, nil
}

func (s *Server) getComment(
	ctx context.Context,
	id int64,
	userInfo *userinfo.UserInfo,
) (*documents.Comment, error) {
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
			tUserProps.AvatarFileID.AS("creator.avatar_file_id"),
			tAvatar.FilePath.AS("creator.avatar"),
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
			tDComments.ID.EQ(jet.Int64(id)),
		).
		LIMIT(1)

	comment := &documents.Comment{}
	if err := stmt.QueryContext(ctx, s.db, comment); err != nil {
		return nil, err
	}

	if comment.GetCreator() != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, comment.GetCreator())
	}

	return comment, nil
}

func (s *Server) DeleteComment(
	ctx context.Context,
	req *pbdocuments.DeleteCommentRequest,
) (*pbdocuments.DeleteCommentResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.GetCommentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "DeleteComment",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	comment, err := s.getComment(ctx, req.GetCommentId(), userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if comment.GetCreatorJob() == "" {
		comment.CreatorJob = userInfo.GetJob()
	}

	check, err := s.access.CanUserAccessTarget(
		ctx,
		comment.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_COMMENT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsdocuments.ErrCommentDeleteDenied
	}

	// Field Permission Check
	fields, err := s.ps.AttrStringList(
		userInfo,
		permsdocuments.DocumentsServicePerm,
		permsdocuments.DocumentsServiceDeleteCommentPerm,
		permsdocuments.DocumentsServiceDeleteCommentAccessPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !access.CheckIfHasOwnJobAccess(
		fields,
		userInfo,
		comment.GetCreatorJob(),
		comment.GetCreator(),
	) {
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
			tDComments.ID.EQ(jet.Int64(req.GetCommentId())),
			tDComments.DeletedAt.IS_NULL(),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if _, err := addDocumentActivity(ctx, s.db, &documents.DocActivity{
		DocumentId:   comment.GetDocumentId(),
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_COMMENT_DELETED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.GetJob(),
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return &pbdocuments.DeleteCommentResponse{}, nil
}

func (s *Server) notifyUsersNewComment(
	ctx context.Context,
	documentId int64,
	sourceUserId int32,
) error {
	userInfo, err := s.ui.GetUserInfoWithoutAccountId(ctx, sourceUserId)
	if err != nil {
		return err
	}

	doc, err := s.getDocument(ctx, tDocument.ID.EQ(jet.Int64(documentId)), userInfo, false)
	if err != nil {
		return err
	}
	if doc == nil || doc.GetDeletedAt() != nil {
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
			tDComments.DocumentID.EQ(jet.Int64(doc.GetId())),
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
	if doc.CreatorId != nil && sourceUserId != doc.GetCreatorId() &&
		!slices.Contains(targetUserIds, doc.GetCreatorId()) {
		userInfo, err := s.ui.GetUserInfoWithoutAccountId(ctx, sourceUserId)
		if err != nil {
			return err
		}

		check, err := s.access.CanUserAccessTarget(
			ctx,
			doc.GetId(),
			userInfo,
			documents.AccessLevel_ACCESS_LEVEL_VIEW,
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
		userInfo, err := s.ui.GetUserInfoWithoutAccountId(ctx, sourceUserId)
		if err != nil {
			return err
		}

		check, err := s.access.CanUserAccessTarget(
			ctx,
			doc.GetId(),
			userInfo,
			documents.AccessLevel_ACCESS_LEVEL_VIEW,
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
				CausedBy: &users.UserShort{
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
