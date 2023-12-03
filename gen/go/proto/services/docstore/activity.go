package docstore

import (
	"context"
	"regexp"
	"strings"

	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	"github.com/galexrt/fivenet/gen/go/proto/resources/documents"
	errorsdocstore "github.com/galexrt/fivenet/gen/go/proto/services/docstore/errors"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/utils/dbutils"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	ActivityDefaultPageLimit = 10
)

var (
	tDocActivity = table.FivenetDocumentsActivity
)

func (s *Server) AddDocumentActivity(ctx context.Context, tx qrm.DB, activitiy *documents.DocActivity) error {
	stmt := tDocActivity.
		INSERT(
			tDocActivity.DocumentID,
			tDocActivity.ActivityType,
			tDocActivity.CreatorID,
			tDocActivity.CreatorJob,
			tDocActivity.Reason,
			tDocActivity.Data,
		).
		VALUES(
			activitiy.DocumentId,
			activitiy.ActivityType,
			activitiy.CreatorId,
			activitiy.CreatorJob,
			activitiy.Reason,
			activitiy.Data,
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		if !dbutils.IsDuplicateError(err) {
			return err
		}
	}

	return nil
}

func (s *Server) ListDocumentActivity(ctx context.Context, req *ListDocumentActivityRequest) (*ListDocumentActivityResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.id", int64(req.DocumentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	ok, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errorsdocstore.ErrFailedQuery
	}
	if !ok {
		return nil, errorsdocstore.ErrDocViewDenied
	}

	tDocActivity := table.FivenetDocumentsActivity.AS("doc_activity")

	condition := tDocActivity.DocumentID.EQ(jet.Uint64(req.DocumentId))

	countStmt := tDocActivity.
		SELECT(
			jet.COUNT(tDocActivity.ID).AS("datacount.totalcount"),
		).
		FROM(
			tDocActivity,
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, err
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(ActivityDefaultPageLimit)
	resp := &ListDocumentActivityResponse{
		Pagination: pag,
		Activity:   []*documents.DocActivity{},
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := tDocActivity.
		SELECT(
			tDocActivity.ID,
			tDocActivity.CreatedAt,
			tDocActivity.DocumentID,
			tDocActivity.ActivityType,
			tDocActivity.CreatorID,
			tDocActivity.CreatorJob,
			tDocActivity.Reason,
			tDocActivity.Data,
			tCreator.ID,
			tCreator.Identifier,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
		).
		FROM(
			tDocActivity.
				LEFT_JOIN(tCreator,
					tCreator.ID.EQ(tDocActivity.CreatorID),
				),
		).
		WHERE(condition).
		OFFSET(
			req.Pagination.Offset,
		).
		ORDER_BY(
			tDocActivity.ID.DESC(),
		).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Activity); err != nil {
		return nil, err
	}

	resp.Pagination.Update(count.TotalCount, len(resp.Activity))

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := 0; i < len(resp.Activity); i++ {
		if resp.Activity[i].Creator != nil {
			jobInfoFn(resp.Activity[i].Creator)
		}
	}

	return resp, nil
}

var brFixer = regexp.MustCompile(`(?m)(<br>)+([ \n]*)(<br>)+`)

func (s *Server) generateDiff(oldContent string, newContent string) (string, error) {
	oldContent = brFixer.ReplaceAllString(oldContent, "<br>")
	newContent = brFixer.ReplaceAllString(newContent, "<br>")
	res, err := s.htmlDiff.HTMLdiff([]string{oldContent, newContent})
	if err != nil {
		// Fallback to the new content
		return newContent, nil
	}

	out := res[0]
	// If no change markers found, return empty string
	if !strings.Contains(out, "bg-") {
		return "", nil
	}

	return out, nil
}

func (s *Server) generateDocumentDiff(old *documents.Document, newDoc *documents.Document) (*documents.DocUpdated, error) {
	diff := &documents.DocUpdated{}

	titleDiff, err := s.generateDiff(old.Title, newDoc.Title)
	if err != nil {
		return nil, err
	}
	if titleDiff != "" {
		diff.TitleDiff = &titleDiff
	}

	contentDiff, err := s.generateDiff(old.Content, newDoc.Content)
	if err != nil {
		return nil, err
	}
	if contentDiff != "" {
		diff.ContentDiff = &contentDiff
	}

	stateDiff, err := s.generateDiff(old.State, newDoc.State)
	if err != nil {
		return nil, err
	}
	if stateDiff != "" {
		diff.StateDiff = &stateDiff
	}

	return diff, nil
}
