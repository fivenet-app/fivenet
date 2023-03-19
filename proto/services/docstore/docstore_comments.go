package docstore

import (
	context "context"

	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/htmlsanitizer"
	"github.com/galexrt/arpanet/proto/resources/documents"
	jet "github.com/go-jet/jet/v2/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	condition := jet.AND(
		dComments.DocumentID.EQ(jet.Uint64(req.DocumentId)),
		dComments.DeletedAt.IS_NULL(),
	)
	countStmt := dComments.
		SELECT(
			jet.COUNT(docs.ID).AS("total_count"),
		).
		FROM(
			dComments,
		).
		WHERE(condition)
	var count struct{ TotalCount int64 }
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, err
	}

	stmt := dComments.
		SELECT(
			dComments.ID,
			dComments.Comment,
			dComments.CreatorID,
			u.ID,
			u.Identifier,
			u.Job,
			u.JobGrade,
			u.Firstname,
			u.Lastname,
		).
		FROM(
			dComments.
				LEFT_JOIN(u,
					dComments.CreatorID.EQ(u.ID),
				),
		).
		WHERE(condition)

	resp := &GetDocumentCommentsResponse{}
	if err := stmt.QueryContext(ctx, s.db, resp.Comments); err != nil {
		return nil, err
	}

	resp.TotalCount = count.TotalCount
	if req.Offset >= resp.TotalCount {
		resp.Offset = 0
	} else {
		resp.Offset = req.Offset
	}
	resp.End = resp.Offset + int64(len(resp.Comments))

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

	resp := &PostDocumentCommentResponse{}
	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	return resp, nil
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

	stmt := dComments.
		UPDATE().
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
