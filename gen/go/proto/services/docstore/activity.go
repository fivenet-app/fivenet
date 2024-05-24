package docstore

import (
	"context"
	"errors"
	"regexp"
	"strings"

	database "github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/documents"
	errorsdocstore "github.com/fivenet-app/fivenet/gen/go/proto/services/docstore/errors"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	ActivityDefaultPageSize = 10
)

var (
	tDocActivity = table.FivenetDocumentsActivity
)

func (s *Server) addDocumentActivity(ctx context.Context, tx qrm.DB, activitiy *documents.DocActivity) (uint64, error) {
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

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		if !dbutils.IsDuplicateError(err) {
			return 0, err
		}
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(lastId), nil
}

func (s *Server) ListDocumentActivity(ctx context.Context, req *ListDocumentActivityRequest) (*ListDocumentActivityResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.id", int64(req.DocumentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	if !check {
		return nil, errorsdocstore.ErrDocViewDenied
	}

	tDocActivity := table.FivenetDocumentsActivity.AS("doc_activity")

	condition := tDocActivity.DocumentID.EQ(jet.Uint64(req.DocumentId))
	if len(req.ActivityTypes) > 0 {
		ids := make([]jet.Expression, len(req.ActivityTypes))
		for i := 0; i < len(req.ActivityTypes); i++ {
			ids[i] = jet.Int16(int16(*req.ActivityTypes[i].Enum()))
		}
		condition = condition.AND(tDocActivity.ActivityType.IN(ids...))
	}

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
		if !errors.Is(err, qrm.ErrNoRows) {
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
			}
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, ActivityDefaultPageSize)
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
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	resp.Pagination.Update(len(resp.Activity))

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
	// If no "htmldiff" change markers are found, return empty string
	if !strings.Contains(out, "htmldiff") {
		return "", nil
	}

	return out, nil
}

// generateDocumentDiff Only generates diff if the old and new contents are not equal, using a simple "string comparison"
func (s *Server) generateDocumentDiff(oldDoc *documents.Document, newDoc *documents.Document) (*documents.DocUpdated, error) {
	diff := &documents.DocUpdated{}

	if !strings.EqualFold(oldDoc.Title, newDoc.Title) {
		titleDiff, err := s.generateDiff(oldDoc.Title, newDoc.Title)
		if err != nil {
			return nil, err
		}
		if titleDiff != "" {
			diff.TitleDiff = &titleDiff
		}
	}

	if !strings.EqualFold(oldDoc.Content, newDoc.Content) {
		contentDiff, err := s.generateDiff(oldDoc.Content, newDoc.Content)
		if err != nil {
			return nil, err
		}
		if contentDiff != "" {
			diff.ContentDiff = &contentDiff
		}
	}

	if !strings.EqualFold(oldDoc.State, newDoc.State) {
		stateDiff, err := s.generateDiff(oldDoc.State, newDoc.State)
		if err != nil {
			return nil, err
		}
		if stateDiff != "" {
			diff.StateDiff = &stateDiff
		}
	}

	return diff, nil
}
