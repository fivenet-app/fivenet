package docstore

import (
	context "context"

	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	"github.com/galexrt/fivenet/gen/go/proto/resources/documents"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	errorsdocstore "github.com/galexrt/fivenet/gen/go/proto/services/docstore/errors"
	permsdocstore "github.com/galexrt/fivenet/gen/go/proto/services/docstore/perms"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/errswrap"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	CommentsDefaultPageLimit = 7
)

var (
	tDComments = table.FivenetDocumentsComments
)

func (s *Server) GetComments(ctx context.Context, req *GetCommentsRequest) (*GetCommentsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.id", int64(req.DocumentId)))

	ok, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(errorsdocstore.ErrFailedQuery, err)
	}
	if !ok {
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
		return nil, errswrap.NewError(errorsdocstore.ErrFailedQuery, err)
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(CommentsDefaultPageLimit)
	resp := &GetCommentsResponse{
		Pagination: pag,
		Comments:   []*documents.Comment{},
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	columns := []jet.Projection{
		tDComments.ID,
		tDComments.DocumentID,
		tDComments.CreatedAt,
		tDComments.UpdatedAt,
		tDComments.Comment,
		tDComments.CreatorID,
		tCreator.ID,
		tCreator.Identifier,
		tCreator.Job,
		tCreator.JobGrade,
		tCreator.Firstname,
		tCreator.Lastname,
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
		return nil, errswrap.NewError(errorsdocstore.ErrFailedQuery, err)
	}

	resp.Pagination.Update(count.TotalCount, len(resp.Comments))

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := 0; i < len(resp.Comments); i++ {
		if resp.Comments[i].Creator != nil {
			jobInfoFn(resp.Comments[i].Creator)
		}
	}

	return resp, nil
}

func (s *Server) PostComment(ctx context.Context, req *PostCommentRequest) (*PostCommentResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "PostComment",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.id", int64(req.Comment.DocumentId)))

	check, err := s.checkIfUserHasAccessToDoc(ctx, req.Comment.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_COMMENT)
	if err != nil {
		return nil, errswrap.NewError(errorsdocstore.ErrFailedQuery, err)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsdocstore.ErrCommentPostDenied
	}

	stmt := tDComments.
		INSERT(
			tDComments.DocumentID,
			tDComments.Comment,
			tDComments.CreatorID,
		).
		VALUES(
			req.Comment.DocumentId,
			req.Comment.Comment,
			userInfo.UserId,
		)

	result, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, errswrap.NewError(errorsdocstore.ErrFailedQuery, err)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, errswrap.NewError(errorsdocstore.ErrFailedQuery, err)
	}

	if _, err := s.addDocumentActivity(ctx, s.db, &documents.DocActivity{
		DocumentId:   req.Comment.DocumentId,
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_COMMENT_ADDED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.Job,
	}); err != nil {
		return nil, errswrap.NewError(errorsdocstore.ErrFailedQuery, err)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	return &PostCommentResponse{
		Id: uint64(lastId),
	}, nil
}

func (s *Server) EditComment(ctx context.Context, req *EditCommentRequest) (*EditCommentResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "EditComment",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.id", int64(req.Comment.DocumentId)), attribute.Int64("fivenet.docstore.comment_id", int64(req.Comment.Id)))

	check, err := s.checkIfUserHasAccessToDoc(ctx, req.Comment.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_COMMENT)
	if err != nil {
		return nil, errswrap.NewError(errorsdocstore.ErrFailedQuery, err)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsdocstore.ErrCommentEditDenied
	}

	comment, err := s.getComment(ctx, req.Comment.Id)
	if err != nil {
		return nil, errswrap.NewError(errorsdocstore.ErrFailedQuery, err)
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
		return nil, errswrap.NewError(errorsdocstore.ErrFailedQuery, err)
	}

	comment.Comment = req.Comment.Comment

	if _, err := s.addDocumentActivity(ctx, s.db, &documents.DocActivity{
		DocumentId:   req.Comment.DocumentId,
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_COMMENT_UPDATED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.Job,
	}); err != nil {
		return nil, errswrap.NewError(errorsdocstore.ErrFailedQuery, err)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &EditCommentResponse{
		Comment: comment,
	}, nil
}

func (s *Server) getComment(ctx context.Context, id uint64) (*documents.Comment, error) {
	comment := &documents.Comment{}

	dComments := tDComments.AS("comment")
	stmt := dComments.
		SELECT(
			dComments.ID,
			dComments.CreatedAt,
			dComments.UpdatedAt,
			dComments.Comment,
			dComments.CreatorID,
		).
		FROM(
			dComments,
		).
		WHERE(
			dComments.ID.EQ(jet.Uint64(id)),
		).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, s.db, comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *Server) DeleteComment(ctx context.Context, req *DeleteCommentRequest) (*DeleteCommentResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "DeleteComment",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.comment_id", int64(req.CommentId)))

	comment, err := s.getComment(ctx, req.CommentId)
	if err != nil {
		return nil, errswrap.NewError(errorsdocstore.ErrFailedQuery, err)
	}

	check, err := s.checkIfUserHasAccessToDoc(ctx, comment.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_COMMENT)
	if err != nil {
		return nil, errswrap.NewError(errorsdocstore.ErrFailedQuery, err)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsdocstore.ErrCommentDeleteDenied
	}

	// Field Permission Check
	fieldsAttr, err := s.ps.Attr(userInfo, permsdocstore.DocStoreServicePerm, permsdocstore.DocStoreServiceDeleteCommentPerm, permsdocstore.DocStoreServiceDeleteCommentAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(errorsdocstore.ErrFailedQuery, err)
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}
	if !s.checkIfHasAccess(fields, userInfo, comment.Creator.Job, comment.Creator) {
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
		return nil, errswrap.NewError(errorsdocstore.ErrFailedQuery, err)
	}

	if _, err := s.addDocumentActivity(ctx, s.db, &documents.DocActivity{
		DocumentId:   uint64(comment.DocumentId),
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_COMMENT_DELETED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.Job,
	}); err != nil {
		return nil, errswrap.NewError(errorsdocstore.ErrFailedQuery, err)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeleteCommentResponse{}, nil
}
