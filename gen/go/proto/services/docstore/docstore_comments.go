package docstore

import (
	context "context"

	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	"github.com/galexrt/fivenet/gen/go/proto/resources/documents"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/htmlsanitizer"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	CommentsDefaultPageLimit = 7
)

var (
	dComments = table.FivenetDocumentsComments
)

func (s *Server) GetDocumentComments(ctx context.Context, req *GetDocumentCommentsRequest) (*GetDocumentCommentsResponse, error) {
	userInfo := auth.GetUserInfoFromContext(ctx)
	ok, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userInfo, true, documents.ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to view document comments!")
	}

	dComments := dComments.AS("documentcomment")
	condition := jet.AND(
		dComments.DocumentID.EQ(jet.Uint64(req.DocumentId)),
		dComments.DeletedAt.IS_NULL(),
	)

	countStmt := dComments.
		SELECT(
			jet.COUNT(dComments.ID).AS("datacount.totalcount"),
		).
		FROM(
			dComments,
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, err
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(CommentsDefaultPageLimit)
	resp := &GetDocumentCommentsResponse{
		Pagination: pag,
		Comments:   []*documents.DocumentComment{},
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := dComments.
		SELECT(
			dComments.ID,
			dComments.DocumentID,
			dComments.CreatedAt,
			dComments.UpdatedAt,
			dComments.Comment,
			dComments.CreatorID,
			uCreator.ID,
			uCreator.Identifier,
			uCreator.Job,
			uCreator.JobGrade,
			uCreator.Firstname,
			uCreator.Lastname,
		).
		FROM(
			dComments.
				LEFT_JOIN(uCreator,
					dComments.CreatorID.EQ(uCreator.ID),
				),
		).
		WHERE(condition).
		OFFSET(
			req.Pagination.Offset,
		).
		ORDER_BY(
			dComments.CreatedAt.DESC(),
		).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Comments); err != nil {
		return nil, err
	}

	resp.Pagination.Update(count.TotalCount, len(resp.Comments))

	for i := 0; i < len(resp.Comments); i++ {
		if resp.Comments[i].Creator != nil {
			s.c.EnrichJobInfo(resp.Comments[i].Creator)
		}
	}

	return resp, nil
}

func (s *Server) PostDocumentComment(ctx context.Context, req *PostDocumentCommentRequest) (*PostDocumentCommentResponse, error) {
	userInfo := auth.GetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "PostDocumentComment",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.AddEntryWithData(auditEntry, req)

	check, err := s.checkIfUserHasAccessToDoc(ctx, req.Comment.DocumentId, userInfo, false, documents.ACCESS_LEVEL_COMMENT)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to post a comment on this document!")
	}

	// Clean comment from
	req.Comment.Comment = htmlsanitizer.StripTags(req.Comment.Comment)

	stmt := dComments.
		INSERT(
			dComments.DocumentID,
			dComments.Comment,
			dComments.CreatorID,
		).
		VALUES(
			req.Comment.DocumentId,
			req.Comment.Comment,
			userInfo.UserId,
		)

	result, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EVENT_TYPE_CREATED)

	return &PostDocumentCommentResponse{
		Id: uint64(lastId),
	}, nil
}

func (s *Server) EditDocumentComment(ctx context.Context, req *EditDocumentCommentRequest) (*EditDocumentCommentResponse, error) {
	userInfo := auth.GetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "EditDocumentComment",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.AddEntryWithData(auditEntry, req)

	check, err := s.checkIfUserHasAccessToDoc(ctx, req.Comment.DocumentId, userInfo, false, documents.ACCESS_LEVEL_COMMENT)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to edit this comment!")
	}

	comment, err := s.getDocumentComment(ctx, req.Comment.Id)
	if err != nil {
		return nil, err
	}
	if comment.CreatorId != userInfo.UserId {
		return nil, status.Error(codes.PermissionDenied, "You can't edit others document comments!")
	}

	// Clean comment from
	req.Comment.Comment = htmlsanitizer.StripTags(req.Comment.Comment)

	stmt := dComments.
		UPDATE(
			dComments.Comment,
		).
		SET(
			dComments.Comment.SET(jet.String(req.Comment.Comment)),
		).
		WHERE(
			jet.AND(
				dComments.ID.EQ(jet.Uint64(req.Comment.Id)),
				dComments.DeletedAt.IS_NULL(),
			),
		)

	resp := &EditDocumentCommentResponse{}
	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EVENT_TYPE_UPDATED)

	return resp, nil
}

func (s *Server) getDocumentComment(ctx context.Context, id uint64) (*documents.DocumentComment, error) {
	comment := &documents.DocumentComment{}

	dComments := dComments.AS("documentcomment")
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

func (s *Server) DeleteDocumentComment(ctx context.Context, req *DeleteDocumentCommentRequest) (*DeleteDocumentCommentResponse, error) {
	userInfo := auth.GetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "DeleteDocumentComment",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.AddEntryWithData(auditEntry, req)

	comment, err := s.getDocumentComment(ctx, req.CommentId)
	if err != nil {
		return nil, err
	}
	// If the requestor is not the creator nor a superuser
	if comment.CreatorId != userInfo.UserId && !userInfo.SuperUser {
		return nil, status.Error(codes.PermissionDenied, "You can't delete others document comments!")
	}

	stmt := dComments.
		UPDATE(
			dComments.DeletedAt,
		).
		SET(
			dComments.DeletedAt.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(
			jet.AND(
				dComments.ID.EQ(jet.Uint64(req.CommentId)),
				dComments.DeletedAt.IS_NULL(),
			),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EVENT_TYPE_DELETED)

	return &DeleteDocumentCommentResponse{}, nil
}
