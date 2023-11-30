package docstore

import (
	"context"
	"fmt"

	"github.com/galexrt/fivenet/gen/go/proto/resources/common"
	"github.com/galexrt/fivenet/gen/go/proto/resources/documents"
	"github.com/galexrt/fivenet/gen/go/proto/resources/notifications"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/notifi"
	"github.com/galexrt/fivenet/query/fivenet/model"
	jet "github.com/go-jet/jet/v2/mysql"
)

func (s *Server) RequestDocumentAction(ctx context.Context, req *RequestDocumentActionRequest) (*RequestDocumentActionResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "RequestDocumentAction",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	doc, err := s.getDocument(ctx, tDocument.ID.EQ(jet.Uint64(req.DocumentId)), userInfo)
	if err != nil {
		return nil, ErrFailedQuery
	}

	resp := &RequestDocumentActionResponse{}

	if err := s.AddDocumentActivity(ctx, s.db, &documents.DocActivity{
		DocumentId:   doc.Id,
		CreatorId:    &userInfo.UserId,
		ActivityType: req.RequestType,
		CreatorJob:   userInfo.Job,
		Reason:       req.Reason,
		Data:         &documents.DocActivityData{},
	}); err != nil {
		return nil, ErrFailedQuery
	}

	// If the document has no creator anymore, nothing we can do here
	if doc.CreatorId == nil {
		return resp, nil
	}

	// TODO check if such a request was already made in the last 6 hours

	if err := s.notifyUser(ctx, doc, userInfo.UserId, int32(*doc.CreatorId)); err != nil {
		return nil, ErrFailedQuery
	}

	return resp, nil
}

// TODO add logic to accept/decline a document request action

func (s *Server) notifyUser(ctx context.Context, doc *documents.Document, sourceUserId int32, targetUserId int32) error {
	userInfo, err := s.ui.GetUserInfoWithoutAccountId(ctx, targetUserId)
	if err != nil {
		return err
	}

	if doc.Creator != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, doc.Creator)
	}

	nType := string(notifi.InfoType)
	not := &notifications.Notification{
		UserId: targetUserId,
		Title: &common.TranslateItem{
			Key: "notifications.notifi.document_relation_mentioned.title",
		},
		Content: &common.TranslateItem{
			Key:        "notifications.notifi.document_relation_mentioned.content",
			Parameters: []string{doc.Title},
		},
		Type:     &nType,
		Category: notifications.NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT,
		Data: &notifications.Data{
			Link: &notifications.Link{
				To: fmt.Sprintf("/documents/%d", doc.Id),
			},
			CausedBy: &users.UserShort{
				UserId: sourceUserId,
			},
		},
	}
	if err := s.notif.NotifyUser(ctx, not); err != nil {
		return err
	}

	return nil
}
