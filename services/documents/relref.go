package documents

import (
	context "context"
	"fmt"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	documentsreferences "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/references"
	documentsrelations "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/relations"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications"
	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	pbdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorsdocuments "github.com/fivenet-app/fivenet/v2026/services/documents/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) GetDocumentReferences(
	ctx context.Context,
	req *pbdocuments.GetDocumentReferencesRequest,
) (*pbdocuments.GetDocumentReferencesResponse, error) {
	logging.InjectFields(ctx, logging.Fields{documentIDLogFieldKey, req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.canUserAccessDocument(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsdocuments.ErrFeedRefsViewDenied
	}

	dest, err := s.store.ListDocumentReferences(ctx, req.GetDocumentId(), userInfo.GetSuperuser())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	resp := &pbdocuments.GetDocumentReferencesResponse{References: dest}
	if len(dest) == 0 {
		return resp, nil
	}

	var docIds []int64
	for _, ref := range dest {
		docIds = append(docIds, ref.GetSourceDocumentId(), ref.GetTargetDocumentId())
	}

	ids, err := s.canUserAccessDocumentIDs(
		ctx,
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW,
		docIds...)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if len(ids) == 0 {
		return &pbdocuments.GetDocumentReferencesResponse{}, nil
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range dest {
		if dest[i].GetCreator() != nil {
			jobInfoFn(dest[i].GetCreator())
		}

		s.docCategories.Enrich(dest[i].GetSourceDocument())
		if dest[i].GetSourceDocument().GetCreator() != nil {
			jobInfoFn(dest[i].GetSourceDocument().GetCreator())
		}

		s.docCategories.Enrich(dest[i].GetTargetDocument())
	}

	return resp, nil
}

func (s *Server) GetDocumentRelations(
	ctx context.Context,
	req *pbdocuments.GetDocumentRelationsRequest,
) (*pbdocuments.GetDocumentRelationsResponse, error) {
	logging.InjectFields(ctx, logging.Fields{documentIDLogFieldKey, req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.canUserAccessDocument(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsdocuments.ErrFeedRelsViewDenied
	}

	relations, err := s.store.ListDocumentRelations(
		ctx,
		req.GetDocumentId(),
		userInfo.GetSuperuser(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range relations {
		if relations[i].GetSourceUser() != nil {
			jobInfoFn(relations[i].GetSourceUser())
		}
		if relations[i].GetTargetUser() != nil {
			jobInfoFn(relations[i].GetTargetUser())
		}
	}

	return &pbdocuments.GetDocumentRelationsResponse{
		Relations: relations,
	}, nil
}

func (s *Server) AddDocumentReference(
	ctx context.Context,
	req *pbdocuments.AddDocumentReferenceRequest,
) (*pbdocuments.AddDocumentReferenceResponse, error) {
	logging.InjectFields(ctx, logging.Fields{
		"fivenet.documents.source_document_id", req.GetReference().GetSourceDocumentId(),
		"fivenet.documents.target_document_id", req.GetReference().GetTargetDocumentId(),
	})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if req.GetReference().GetSourceDocumentId() == req.GetReference().GetTargetDocumentId() {
		return nil, errorsdocuments.ErrFeedRefSelf
	}

	// Check if user has access to both documents
	check, err := s.canUserAccessDocuments(
		ctx,
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT,
		req.GetReference().GetSourceDocumentId(),
		req.GetReference().GetTargetDocumentId(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsdocuments.ErrFeedRefAddDenied
	}

	req.Reference.CreatorId = &userInfo.UserId

	lastId, err := s.store.CreateDocumentReference(
		ctx,
		s.db,
		&documentsreferences.DocumentReference{
			SourceDocumentId: req.GetReference().GetSourceDocumentId(),
			TargetDocumentId: req.GetReference().GetTargetDocumentId(),
			Reference:        req.GetReference().GetReference(),
			CreatorId:        req.GetReference().CreatorId,
		},
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

	return &pbdocuments.AddDocumentReferenceResponse{Id: lastId}, nil
}

func (s *Server) RemoveDocumentReference(
	ctx context.Context,
	req *pbdocuments.RemoveDocumentReferenceRequest,
) (*pbdocuments.RemoveDocumentReferenceResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.reference_id", req.GetId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	ref, err := s.store.GetDocumentReference(ctx, req.GetId(), userInfo.GetSuperuser())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	check, err := s.canUserAccessDocuments(
		ctx,
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT,
		ref.GetSourceDocumentId(),
		ref.GetTargetDocumentId(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsdocuments.ErrFeedRefRemoveDenied
	}

	if err := s.store.DeleteDocumentReference(ctx, s.db, req.GetId()); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)

	return &pbdocuments.RemoveDocumentReferenceResponse{}, nil
}

func (s *Server) AddDocumentRelation(
	ctx context.Context,
	req *pbdocuments.AddDocumentRelationRequest,
) (*pbdocuments.AddDocumentRelationResponse, error) {
	logging.InjectFields(ctx, logging.Fields{
		documentIDLogFieldKey, req.GetRelation().GetDocumentId(),
		"fivenet.documents.source_user_id", req.GetRelation().GetSourceUserId(),
		"fivenet.documents.target_user_id", req.GetRelation().GetTargetUserId(),
	})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.canUserAccessDocument(
		ctx,
		req.GetRelation().GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsdocuments.ErrFeedRelAddDenied
	}

	req.Relation.SourceUserId = userInfo.GetUserId()

	lastId, created, err := s.store.CreateDocumentRelation(
		ctx,
		s.db,
		&documentsrelations.DocumentRelation{
			DocumentId:   req.GetRelation().GetDocumentId(),
			SourceUserId: req.GetRelation().GetSourceUserId(),
			Relation:     req.GetRelation().GetRelation(),
			TargetUserId: req.GetRelation().GetTargetUserId(),
		},
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if created {
		if err := s.addUserActivity(
			ctx,
			s.db,
			userInfo.GetUserId(),
			req.GetRelation().GetTargetUserId(),
			usersactivity.UserActivityType_USER_ACTIVITY_TYPE_DOCUMENT,
			"",
			&usersactivity.UserActivityData{
				Data: &usersactivity.UserActivityData_DocumentRelation{
					DocumentRelation: &usersactivity.CitizenDocumentRelation{
						Added:      true,
						DocumentId: req.GetRelation().GetDocumentId(),
						Relation:   req.GetRelation().GetRelation(),
					},
				},
			},
		); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		if req.GetRelation().
			GetRelation() ==
			documentsrelations.DocRelation_DOC_RELATION_MENTIONED {
			if err := s.notifyMentionedUser(
				ctx,
				req.GetRelation().GetDocumentId(),
				userInfo.GetUserId(),
				req.GetRelation().GetTargetUserId(),
			); err != nil {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}
		}
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

	return &pbdocuments.AddDocumentRelationResponse{Id: lastId}, nil
}

func (s *Server) RemoveDocumentRelation(
	ctx context.Context,
	req *pbdocuments.RemoveDocumentRelationRequest,
) (*pbdocuments.RemoveDocumentRelationResponse, error) {
	logging.InjectFields(ctx, logging.Fields{documentIDLogFieldKey, req.GetId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	rel, err := s.store.GetDocumentRelation(ctx, req.GetId(), userInfo.GetSuperuser())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	check, err := s.canUserAccessDocument(
		ctx,
		rel.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsdocuments.ErrFeedRelRemoveDenied
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	if err := s.store.DeleteDocumentRelation(ctx, tx, req.GetId()); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := s.addUserActivity(
		ctx,
		tx,
		userInfo.GetUserId(),
		rel.GetTargetUserId(),
		usersactivity.UserActivityType_USER_ACTIVITY_TYPE_DOCUMENT,
		"",
		&usersactivity.UserActivityData{
			Data: &usersactivity.UserActivityData_DocumentRelation{
				DocumentRelation: &usersactivity.CitizenDocumentRelation{
					Added:      false,
					DocumentId: rel.GetDocumentId(),
					// Relation:   rel.Relation,
				},
			},
		},
	); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)

	return &pbdocuments.RemoveDocumentRelationResponse{}, nil
}

func (s *Server) notifyMentionedUser(
	ctx context.Context,
	documentId int64,
	sourceUserId int32,
	targetUserId int32,
) error {
	userInfo, err := s.ui.GetUserInfo(ctx, targetUserId)
	if err != nil {
		return err
	}

	// Make sure target user has access to document
	check, err := s.canUserAccessDocument(
		ctx,
		documentId,
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil {
		return err
	}
	if !check {
		return nil
	}

	doc, err := s.getDocument(ctx, tDocument.ID.EQ(mysql.Int64(documentId)), userInfo, false)
	if err != nil {
		return err
	}
	if doc == nil {
		return nil
	}

	not := &notifications.Notification{
		UserId: targetUserId,
		Title: &common.I18NItem{
			Key: "notifications.documents.document_relation_mentioned.title",
		},
		Content: &common.I18NItem{
			Key:        "notifications.documents.document_relation_mentioned.content",
			Parameters: map[string]string{"title": doc.GetTitle()},
		},
		Type:     notifications.NotificationType_NOTIFICATION_TYPE_INFO,
		Category: notifications.NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT,
		Data: &notifications.Data{
			Link: &notifications.Link{
				To: fmt.Sprintf("/documents/%d", doc.GetId()),
			},
			CausedBy: &usershort.UserShort{
				UserId: sourceUserId,
			},
		},
	}
	if err := s.notifi.NotifyUser(ctx, not); err != nil {
		return err
	}

	return nil
}
