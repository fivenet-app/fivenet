package documents

import (
	"context"

	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	pbdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	errorsdocuments "github.com/fivenet-app/fivenet/v2026/services/documents/errors"
	documentsstore "github.com/fivenet-app/fivenet/v2026/stores/documents"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) ListDocumentPins(
	ctx context.Context,
	req *pbdocuments.ListDocumentPinsRequest,
) (*pbdocuments.ListDocumentPinsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	pag, docs, err := s.store.ListDocumentPins(ctx, documentsstore.ListDocumentPinsQuery{
		Personal:   req.GetPersonal(),
		Pagination: req.GetPagination(),
		UserInfo:   userInfo,
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	resp := &pbdocuments.ListDocumentPinsResponse{
		Pagination: pag,
	}
	if pag.GetTotalCount() <= 0 {
		return resp, nil
	}
	resp.Documents = docs

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetDocuments() {
		if resp.GetDocuments()[i].GetCreator() != nil {
			jobInfoFn(resp.GetDocuments()[i].GetCreator())
		}
	}

	return resp, nil
}

func (s *Server) ToggleDocumentPin(
	ctx context.Context,
	req *pbdocuments.ToggleDocumentPinRequest,
) (*pbdocuments.ToggleDocumentPinResponse, error) {
	logging.InjectFields(ctx, logging.Fields{documentIDLogFieldKey, req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Adding a pin requires view access to the document, but removing a pin does not
	if req.GetState() {
		check, err := s.canUserAccessDocument(
			ctx,
			req.GetDocumentId(),
			userInfo,
			documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW,
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
		}
		if !check && !userInfo.GetJobAdmin() {
			return nil, errorsdocuments.ErrDocViewDenied
		}
	}

	if req.GetState() {
		if err := s.store.CreateDocumentPin(
			ctx,
			s.db,
			req.GetDocumentId(),
			userInfo,
			req.GetPersonal(),
		); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	} else {
		if err := s.store.DeleteDocumentPin(
			ctx,
			s.db,
			req.GetDocumentId(),
			userInfo,
			req.GetPersonal(),
		); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	pin, err := s.store.GetDocumentPin(ctx, req.GetDocumentId(), userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.ToggleDocumentPinResponse{
		Pin: pin,
	}, nil
}
