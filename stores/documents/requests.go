package documentsstore

import (
	"context"
	"errors"

	documentsrequests "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/requests"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) AddDocumentReq(
	ctx context.Context,
	tx qrm.DB,
	request *documentsrequests.DocRequest,
) (int64, error) {
	tDocRequest := table.FivenetDocumentsRequests
	stmt := tDocRequest.
		INSERT(
			tDocRequest.DocumentID,
			tDocRequest.RequestType,
			tDocRequest.CreatorID,
			tDocRequest.CreatorJob,
			tDocRequest.Reason,
			tDocRequest.Data,
			tDocRequest.Accepted,
		).
		VALUES(
			request.GetDocumentId(),
			request.GetRequestType(),
			request.GetCreatorId(),
			request.GetCreatorJob(),
			request.Reason,
			request.GetData(),
			request.Accepted,
		)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		if !dbutils.IsDuplicateError(err) {
			return 0, err
		}
	}

	return res.LastInsertId()
}

func (s *Store) UpdateDocumentReq(
	ctx context.Context,
	tx qrm.DB,
	id int64,
	request *documentsrequests.DocRequest,
) error {
	tDocRequest := table.FivenetDocumentsRequests
	stmt := tDocRequest.
		UPDATE(
			tDocRequest.DocumentID,
			tDocRequest.RequestType,
			tDocRequest.CreatorID,
			tDocRequest.CreatorJob,
			tDocRequest.Reason,
			tDocRequest.Data,
			tDocRequest.Accepted,
		).
		SET(
			request.GetDocumentId(),
			request.GetRequestType(),
			request.GetCreatorId(),
			request.GetCreatorJob(),
			request.Reason,
			request.GetData(),
			request.Accepted,
		).
		WHERE(tDocRequest.ID.EQ(mysql.Int64(id))).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	if err != nil && !dbutils.IsDuplicateError(err) {
		return err
	}
	return nil
}

func (s *Store) GetDocumentReq(
	ctx context.Context,
	db qrm.DB,
	condition mysql.BoolExpression,
) (*documentsrequests.DocRequest, error) {
	tDocRequest := table.FivenetDocumentsRequests.AS("doc_request")
	tCreator := table.FivenetUser.AS("creator")

	stmt := tDocRequest.
		SELECT(
			tDocRequest.ID,
			tDocRequest.CreatedAt,
			tDocRequest.UpdatedAt,
			tDocRequest.DocumentID,
			tDocRequest.RequestType,
			tDocRequest.CreatorID,
			tDocRequest.CreatorJob,
			tDocRequest.Reason,
			tDocRequest.Data,
			tCreator.ID,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Job,
			tCreator.Dateofbirth,
			tCreator.PhoneNumber,
		).
		FROM(
			tDocRequest.
				INNER_JOIN(tCreator,
					tCreator.ID.EQ(tDocRequest.CreatorID),
				),
		).
		WHERE(condition).
		LIMIT(1)

	var dest documentsrequests.DocRequest
	if err := stmt.QueryContext(ctx, db, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &dest, nil
}

func (s *Store) DeleteDocumentReq(ctx context.Context, tx qrm.DB, id int64) error {
	tDocRequest := table.FivenetDocumentsRequests
	stmt := tDocRequest.
		DELETE().
		WHERE(tDocRequest.ID.EQ(mysql.Int64(id))).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}
