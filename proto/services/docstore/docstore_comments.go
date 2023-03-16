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
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	ok, err := s.checkIfUserHasAccessToDoc(ctx, userID, job, jobGrade, documents.DOC_ACCESS_VIEW)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to view document comments!")
	}

	condition := adc.DocumentID.EQ(jet.Uint64(req.DocumentID))
	countStmt := adc.SELECT(
		jet.COUNT(ad.ID).AS("total_count"),
	).
		FROM(
			adc,
		).
		WHERE(condition)
	var count struct{ TotalCount int64 }
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, err
	}

	stmt := adc.SELECT(
		adc.ID,
		adc.Comment,
		adc.CreatorID,
		u.ID,
		u.Identifier,
		u.Job,
		u.JobGrade,
		u.Firstname,
		u.Lastname,
	).
		FROM(
			adc.
				LEFT_JOIN(u,
					adc.CreatorID.EQ(u.ID),
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
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	check, err := s.checkIfUserHasAccessToDoc(ctx, userID, job, jobGrade, documents.DOC_ACCESS_VIEW)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to post a comment on this document!")
	}

	// Clean comment from
	req.Comment.Comment = htmlsanitizer.StripTags(req.Comment.Comment)

	stmt := adc.INSERT(
		adc.DocumentID,
		adc.Comment,
		adc.CreatorID,
	).
		VALUES(
			req.Comment.DocumentId,
			req.Comment.Comment,
			userID,
		)

	resp := &PostDocumentCommentResponse{}
	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	return resp, nil
}
func (s *Server) EditDocumentComment(ctx context.Context, req *EditDocumentCommentRequest) (*EditDocumentCommentResponse, error) {
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	check, err := s.checkIfUserHasAccessToDoc(ctx, userID, job, jobGrade, documents.DOC_ACCESS_VIEW)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to edit this comment!")
	}

	stmt := adc.UPDATE().
		SET(
			adc.Comment.SET(jet.String(req.Comment.Comment)),
		).
		WHERE(
			adc.ID.EQ(jet.Uint64(req.Comment.Id)),
		)

	resp := &EditDocumentCommentResponse{}
	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	return resp, nil
}
