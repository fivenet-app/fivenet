package documentsstore

import (
	"context"
	"errors"
	"slices"

	resourcesdatabase "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	resourcesdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents"
	documentspins "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/pins"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) ListDocumentPins(
	ctx context.Context,
	q ListDocumentPinsQuery,
) (*resourcesdatabase.PaginationResponse, []*resourcesdocuments.DocumentShort, error) {
	if q.Pagination == nil {
		q.Pagination = &resourcesdatabase.PaginationRequest{}
	}
	if q.UserInfo == nil {
		q.UserInfo = &userinfo.UserInfo{}
	}

	tDPins := table.FivenetDocumentsPins.AS("document_pin")

	var idCondition mysql.BoolExpression
	if q.Personal {
		idCondition = tDPins.UserID.EQ(mysql.Int32(q.UserInfo.GetUserId()))
	} else {
		idCondition = mysql.OR(
			tDPins.Job.EQ(mysql.String(q.UserInfo.GetJob())),
			tDPins.UserID.EQ(mysql.Int32(q.UserInfo.GetUserId())),
		)
	}

	countStmt := tDPins.
		SELECT(mysql.COUNT(tDPins.DocumentID).AS("data_count.total")).
		FROM(tDPins).
		WHERE(idCondition).
		LIMIT(50)

	var count resourcesdatabase.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, nil, err
		}
	}

	pag, limit := q.Pagination.GetResponseWithPageSize(count.Total, 50)
	if count.Total <= 0 {
		return pag, []*resourcesdocuments.DocumentShort{}, nil
	}

	idStmt := tDPins.
		SELECT(
			tDPins.DocumentID,
			tDPins.Job,
			tDPins.UserID,
			tDPins.CreatedAt,
			tDPins.State,
			tDPins.CreatorID,
		).
		FROM(tDPins).
		WHERE(idCondition).
		LIMIT(50)

	docPins := []*documentspins.DocumentPin{}
	if err := idStmt.QueryContext(ctx, s.db, &docPins); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, nil, err
		}
	}

	docIds := make([]int64, len(docPins))
	for i, pin := range docPins {
		docIds[i] = pin.GetDocumentId()
	}

	docs, err := s.List(ctx, ListQuery{
		DocumentIDs: docIds,
		UserInfo:    q.UserInfo,
		Offset:      q.Pagination.GetOffset(),
		Limit:       limit,
	})
	if err != nil {
		return nil, nil, err
	}

	for _, pin := range docPins {
		idx := slices.IndexFunc(docs, func(doc *resourcesdocuments.DocumentShort) bool {
			return doc.GetId() == pin.GetDocumentId()
		})
		if idx >= 0 {
			if docs[idx].Pin != nil {
				if pin.Job != nil {
					docs[idx].Pin.Job = pin.Job
				}
				if pin.UserId != nil {
					docs[idx].Pin.UserId = pin.UserId
				}
			} else {
				docs[idx].Pin = pin
			}
			continue
		}

		docs = append(docs, &resourcesdocuments.DocumentShort{Id: pin.GetDocumentId(), Pin: pin})
	}

	return pag, docs, nil
}

func (s *Store) CreateDocumentPin(
	ctx context.Context,
	tx qrm.DB,
	documentID int64,
	userInfo *userinfo.UserInfo,
	personal bool,
) error {
	tDPins := table.FivenetDocumentsPins
	job := mysql.NULL
	userId := mysql.NULL
	if personal {
		userId = mysql.Int32(userInfo.GetUserId())
	} else {
		job = mysql.String(userInfo.GetJob())
	}

	stmt := tDPins.
		INSERT(
			tDPins.DocumentID,
			tDPins.Job,
			tDPins.UserID,
			tDPins.CreatorID,
		).
		VALUES(
			documentID,
			job,
			userId,
			userInfo.GetUserId(),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		if !dbutils.IsDuplicateError(err) {
			return err
		}
	}

	return nil
}

func (s *Store) DeleteDocumentPin(
	ctx context.Context,
	tx qrm.DB,
	documentID int64,
	userInfo *userinfo.UserInfo,
	personal bool,
) error {
	tDPins := table.FivenetDocumentsPins
	condition := tDPins.DocumentID.EQ(mysql.Int64(documentID))
	if personal {
		condition = condition.AND(mysql.AND(
			tDPins.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
			tDPins.Job.IS_NULL(),
		))
	} else {
		condition = condition.AND(mysql.AND(
			tDPins.Job.EQ(mysql.String(userInfo.GetJob())),
			tDPins.UserID.IS_NULL(),
		))
	}

	stmt := tDPins.
		DELETE().
		WHERE(condition).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}
