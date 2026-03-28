package documents

import (
	"context"
	"errors"
	"strings"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents"
	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	documentsactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/activity"
	pbdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/textdiff"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2026/services/documents/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

const (
	ActivityDefaultPageSize = 10
)

var tDocActivity = table.FivenetDocumentsActivity

func (s *Server) ListDocumentActivity(
	ctx context.Context,
	req *pbdocuments.ListDocumentActivityRequest,
) (*pbdocuments.ListDocumentActivityResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check {
		return nil, errorsdocuments.ErrDocViewDenied
	}

	tDocActivity := table.FivenetDocumentsActivity.AS("doc_activity")

	condition := tDocActivity.DocumentID.EQ(mysql.Int64(req.GetDocumentId()))
	if len(req.GetActivityTypes()) > 0 {
		ids := make([]mysql.Expression, len(req.GetActivityTypes()))
		for i := range req.GetActivityTypes() {
			ids[i] = mysql.Int32(int32(*req.GetActivityTypes()[i].Enum()))
		}
		condition = condition.AND(tDocActivity.ActivityType.IN(ids...))
	}

	countStmt := tDocActivity.
		SELECT(
			mysql.COUNT(tDocActivity.ID).AS("data_count.total"),
		).
		FROM(
			tDocActivity,
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}
		}
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count.Total, ActivityDefaultPageSize)
	resp := &pbdocuments.ListDocumentActivityResponse{
		Pagination: pag,
		Activity:   []*documentsactivity.DocActivity{},
	}
	if count.Total <= 0 {
		return resp, nil
	}

	tCreator := table.FivenetUser.AS("creator")

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
			req.GetPagination().GetOffset(),
		).
		ORDER_BY(
			tDocActivity.ID.DESC(),
		).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Activity); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetActivity() {
		if resp.GetActivity()[i].GetCreator() != nil {
			jobInfoFn(resp.GetActivity()[i].GetCreator())
		}
	}

	return resp, nil
}

func addDocumentActivity(
	ctx context.Context,
	tx qrm.DB,
	activitiy *documentsactivity.DocActivity,
) (int64, error) {
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
			activitiy.GetDocumentId(),
			activitiy.GetActivityType(),
			activitiy.CreatorId,
			activitiy.GetCreatorJob(),
			activitiy.Reason,
			activitiy.GetData(),
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

	return lastId, nil
}

// generateDocumentDiff Generates diff if the old and new contents are not equal, using a simple "string comparison".
func (s *Server) generateDocumentDiff(
	old *documents.Document,
	new *documents.Document,
) (*documentsactivity.DocUpdated, error) {
	diff := &documentsactivity.DocUpdated{}

	if !strings.EqualFold(old.GetTitle(), new.GetTitle()) {
		if titleDiff := textdiff.DiffText(old.GetTitle(), new.GetTitle()); titleDiff.HasChanges() {
			diff.TitleCdiff = titleDiff
		}
	}

	if !strings.EqualFold(old.GetMeta().GetState(), new.GetMeta().GetState()) {
		if stateDiff := textdiff.DiffText(
			old.GetMeta().GetState(),
			new.GetMeta().GetState(),
		); stateDiff.HasChanges() {
			diff.StateCdiff = stateDiff
		}
	}

	if d := textdiff.DiffText(
		old.GetContent().Extract().GetText(),
		new.GetContent().Extract().GetText(),
	); d.HasChanges() {
		diff.ContentCdiff = d
	}

	return diff, nil
}
