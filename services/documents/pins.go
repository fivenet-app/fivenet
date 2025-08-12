package documents

import (
	"context"
	"errors"
	"slices"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2025/services/documents/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

var tDPins = table.FivenetDocumentsPins.AS("pin")

func (s *Server) ListDocumentPins(
	ctx context.Context,
	req *pbdocuments.ListDocumentPinsRequest,
) (*pbdocuments.ListDocumentPinsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tDPins := table.FivenetDocumentsPins.AS("document_pin")

	var idCondition jet.BoolExpression
	if req.Personal != nil && req.GetPersonal() {
		idCondition = tDPins.UserID.EQ(jet.Int32(userInfo.GetUserId()))
	} else {
		idCondition = jet.OR(
			tDPins.Job.EQ(jet.String(userInfo.GetJob())),
			tDPins.UserID.EQ(jet.Int32(userInfo.GetUserId())),
		)
	}

	countStmt := tDPins.
		SELECT(
			jet.COUNT(tDPins.DocumentID).AS("data_count.total"),
		).
		FROM(tDPins).
		WHERE(idCondition).
		LIMIT(50)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count.Total, 50)
	resp := &pbdocuments.ListDocumentPinsResponse{
		Pagination: pag,
	}
	if count.Total <= 0 {
		return resp, nil
	}

	// Select document IDs from pins
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

	docPins := []*documents.DocumentPin{}
	if err := idStmt.QueryContext(ctx, s.db, &docPins); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	docIdsExpr := make([]jet.Expression, len(docPins))
	for k, pin := range docPins {
		docIdsExpr[k] = jet.Uint64(pin.GetDocumentId())
	}
	condition := tDocumentShort.ID.IN(docIdsExpr...)

	// Select the documents with the list of pin doc ids
	stmt := s.listDocumentsQuery(condition, nil, nil, userInfo).
		ORDER_BY(
			tDocumentShort.CreatedAt.DESC(),
			tDocumentShort.UpdatedAt.DESC(),
		).
		OFFSET(req.GetPagination().GetOffset()).
		GROUP_BY(tDocumentShort.ID).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Documents); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetDocuments() {
		if resp.GetDocuments()[i].GetCreator() != nil {
			jobInfoFn(resp.GetDocuments()[i].GetCreator())
		}
	}

	for i := range docPins {
		idx := slices.IndexFunc(resp.GetDocuments(), func(doc *documents.DocumentShort) bool {
			return doc.GetId() == docPins[i].GetDocumentId()
		})
		if idx > -1 {
			if resp.GetDocuments()[idx].GetPin() != nil {
				if docPins[i].Job != nil {
					resp.Documents[idx].Pin.Job = docPins[i].Job
				}
				if docPins[i].UserId != nil {
					resp.Documents[idx].Pin.UserId = docPins[i].UserId
				}
			} else {
				resp.Documents[idx].Pin = docPins[i]
			}
		} else {
			// If the document is not found in the response, add it as a placeholder to the response
			resp.Documents = append(resp.Documents, &documents.DocumentShort{
				Id:  docPins[i].GetDocumentId(),
				Pin: docPins[i],
			})
		}
	}

	resp.GetPagination().Update(len(resp.GetDocuments()))

	return resp, nil
}

func (s *Server) ToggleDocumentPin(
	ctx context.Context,
	req *pbdocuments.ToggleDocumentPinRequest,
) (*pbdocuments.ToggleDocumentPinResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "ToggleDocumentPin",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	// Adding a pin requires view access to the document, but removing a pin does not
	if req.GetState() {
		check, err := s.access.CanUserAccessTarget(
			ctx,
			req.GetDocumentId(),
			userInfo,
			documents.AccessLevel_ACCESS_LEVEL_VIEW,
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
		}
		if !check && !userInfo.GetSuperuser() {
			return nil, errorsdocuments.ErrDocViewDenied
		}
	}

	tDPins := table.FivenetDocumentsPins

	if req.GetState() {
		job := jet.NULL
		userId := jet.NULL
		if req.Personal != nil && req.GetPersonal() {
			userId = jet.Int32(userInfo.GetUserId())
		} else {
			job = jet.String(userInfo.GetJob())
		}

		stmt := tDPins.
			INSERT(
				tDPins.DocumentID,
				tDPins.Job,
				tDPins.UserID,
				tDPins.CreatorID,
			).
			VALUES(
				req.GetDocumentId(),
				job,
				userId,
				userInfo.GetUserId(),
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}
		}
	} else {
		condition := tDPins.DocumentID.EQ(jet.Uint64(req.GetDocumentId()))
		if req.Personal != nil && req.GetPersonal() {
			condition = condition.AND(jet.AND(
				tDPins.UserID.EQ(jet.Int32(userInfo.GetUserId())),
				tDPins.Job.IS_NULL(),
			))
		} else {
			condition = condition.AND(jet.AND(
				tDPins.Job.EQ(jet.String(userInfo.GetJob())),
				tDPins.UserID.IS_NULL(),
			))
		}

		stmt := tDPins.
			DELETE().
			WHERE(condition).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	pin, err := s.getDocumentPin(ctx, req.GetDocumentId(), userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.ToggleDocumentPinResponse{
		Pin: pin,
	}, nil
}

func (s *Server) getDocumentPin(
	ctx context.Context,
	documentId uint64,
	userInfo *userinfo.UserInfo,
) (*documents.DocumentPin, error) {
	tDPins := table.FivenetDocumentsPins.AS("document_pin")

	condition := jet.AND(
		tDPins.DocumentID.EQ(jet.Uint64(documentId)),
		jet.OR(
			jet.AND(
				tDPins.Job.EQ(jet.String(userInfo.GetJob())),
				tDPins.UserID.IS_NULL(),
			),
			jet.AND(
				tDPins.UserID.EQ(jet.Int32(userInfo.GetUserId())),
				tDPins.Job.IS_NULL(),
			),
		),
	)

	stmt := tDPins.
		SELECT(
			tDPins.DocumentID,
			tDPins.Job,
			tDPins.UserID,
			tDPins.CreatedAt,
			tDPins.State,
			tDPins.CreatorID,
		).
		WHERE(condition).
		ORDER_BY(tDPins.UserID.DESC()).
		LIMIT(2)

	pins := []*documents.DocumentPin{}
	if err := stmt.QueryContext(ctx, s.db, &pins); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if len(pins) == 0 {
		return nil, nil
	}

	pin := &documents.DocumentPin{
		DocumentId: documentId,
	}
	for _, p := range pins {
		if p.Job != nil {
			pin.Job = p.Job
		}
		if p.UserId != nil {
			pin.UserId = p.UserId
		}
		pin.State = p.GetState()
		pin.CreatedAt = p.GetCreatedAt()
		pin.CreatorId = p.GetCreatorId()
	}

	return pin, nil
}
