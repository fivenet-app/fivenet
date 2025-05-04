package wiki

import (
	"context"
	"errors"
	"strings"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/content"
	database "github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/wiki"
	pbwiki "github.com/fivenet-app/fivenet/gen/go/proto/services/wiki"
	"github.com/fivenet-app/fivenet/pkg/dbutils"
	"github.com/fivenet-app/fivenet/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	errorswiki "github.com/fivenet-app/fivenet/services/wiki/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tPActivity = table.FivenetWikiPageActivity

func (s *Server) ListPageActivity(ctx context.Context, req *pbwiki.ListPageActivityRequest) (*pbwiki.ListPageActivityResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.wiki.id", int64(req.PageId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(ctx, req.PageId, userInfo, wiki.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	if !check {
		return nil, errorswiki.ErrPageDenied
	}

	tPActivity := table.FivenetWikiPageActivity.AS("page_activity")
	condition := tPActivity.PageID.EQ(jet.Uint64(req.PageId))

	countStmt := tPActivity.
		SELECT(
			jet.COUNT(tPActivity.ID).AS("datacount.totalcount"),
		).
		FROM(
			tPActivity,
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
			}
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, 10)
	resp := &pbwiki.ListPageActivityResponse{
		Pagination: pag,
		Activity:   []*wiki.PageActivity{},
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	tCreator := tables.Users().AS("creator")

	stmt := tPActivity.
		SELECT(
			tPActivity.ID,
			tPActivity.CreatedAt,
			tPActivity.PageID,
			tPActivity.ActivityType,
			tPActivity.CreatorID,
			tPActivity.CreatorJob,
			tPActivity.Reason,
			tPActivity.Data,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
		).
		FROM(
			tPActivity.
				LEFT_JOIN(tCreator,
					tCreator.ID.EQ(tPActivity.CreatorID),
				),
		).
		WHERE(condition).
		OFFSET(
			req.Pagination.Offset,
		).
		ORDER_BY(
			tPActivity.ID.DESC(),
		).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Activity); err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	resp.Pagination.Update(len(resp.Activity))

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.Activity {
		if resp.Activity[i].Creator != nil {
			jobInfoFn(resp.Activity[i].Creator)
		}
	}

	return resp, nil
}

func (s *Server) addPageActivity(ctx context.Context, tx qrm.DB, activitiy *wiki.PageActivity) (uint64, error) {
	stmt := tPActivity.
		INSERT(
			tPActivity.PageID,
			tPActivity.ActivityType,
			tPActivity.CreatorID,
			tPActivity.CreatorJob,
			tPActivity.Reason,
			tPActivity.Data,
		).
		VALUES(
			activitiy.PageId,
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

// generatePageDiff Generates diff if the old and new contents are not equal, using a simple "string comparison"
func (s *Server) generatePageDiff(old *wiki.Page, new *wiki.Page) (*wiki.PageUpdated, error) {
	diff := &wiki.PageUpdated{}

	if !strings.EqualFold(old.Meta.Title, new.Meta.Title) {
		titleDiff, err := s.htmlDiff.FancyDiff(old.Meta.Title, new.Meta.Title)
		if err != nil {
			return nil, err
		}
		if titleDiff != "" {
			diff.TitleDiff = &titleDiff
		}
	}

	if !strings.EqualFold(old.Meta.Description, new.Meta.Description) {
		descriptionDiff, err := s.htmlDiff.FancyDiff(old.Meta.Description, new.Meta.Description)
		if err != nil {
			return nil, err
		}
		if descriptionDiff != "" {
			diff.DescriptionDiff = &descriptionDiff
		}
	}

	newRawContent, err := content.PrettyHTML(*new.Content.RawContent)
	if err != nil {
		return nil, err
	}
	if d := s.htmlDiff.PatchDiff(*old.Content.RawContent, newRawContent); d != "" {
		diff.ContentDiff = &d
	}

	return diff, nil
}
