package docstore

import (
	context "context"

	"github.com/galexrt/fivenet/pkg/auth"
	"github.com/galexrt/fivenet/pkg/htmlsanitizer"
	database "github.com/galexrt/fivenet/proto/resources/common/database"
	"github.com/galexrt/fivenet/proto/resources/documents"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	CommentsDefaultPageLimit = 5
)

var (
	dComments = table.FivenetDocumentsComments
)

func (s *Server) GetDocumentComments(ctx context.Context, req *GetDocumentCommentsRequest) (*GetDocumentCommentsResponse, error) {
	userId, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	ok, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userId, job, jobGrade, true, documents.DOC_ACCESS_VIEW)
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

	resp := &GetDocumentCommentsResponse{
		Pagination: database.EmptyPaginationResponseWithPageSize(req.Pagination.Offset, CommentsDefaultPageLimit),
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
		LIMIT(CommentsDefaultPageLimit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Comments); err != nil {
		return nil, err
	}

	database.PaginationHelperWithPageSize(resp.Pagination,
		count.TotalCount,
		req.Pagination.Offset,
		len(resp.Comments),
		CommentsDefaultPageLimit,
	)

	for i := 0; i < len(resp.Comments); i++ {
		if resp.Comments[i].Creator != nil {
			s.c.EnrichJobInfo(resp.Comments[i].Creator)
		}
	}

	return resp, nil
}

func (s *Server) PostDocumentComment(ctx context.Context, req *PostDocumentCommentRequest) (*PostDocumentCommentResponse, error) {
	userId, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	check, err := s.checkIfUserHasAccessToDoc(ctx, req.Comment.DocumentId, userId, job, jobGrade, false, documents.DOC_ACCESS_COMMENT)
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
			userId,
		)

	result, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &PostDocumentCommentResponse{
		Id: uint64(lastId),
	}, nil
}
func (s *Server) EditDocumentComment(ctx context.Context, req *EditDocumentCommentRequest) (*EditDocumentCommentResponse, error) {
	userId, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	check, err := s.checkIfUserHasAccessToDoc(ctx, req.Comment.DocumentId, userId, job, jobGrade, false, documents.DOC_ACCESS_COMMENT)
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
	if comment.CreatorId != userId {
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
		LIMIT(1)

	if err := stmt.QueryContext(ctx, s.db, comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *Server) DeleteDocumentComment(ctx context.Context, req *DeleteDocumentCommentRequest) (*DeleteDocumentCommentResponse, error) {
	userId := auth.GetUserIDFromContext(ctx)

	comment, err := s.getDocumentComment(ctx, req.CommentId)
	if err != nil {
		return nil, err
	}
	if comment.CreatorId != userId {
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

	return &DeleteDocumentCommentResponse{}, nil
}
