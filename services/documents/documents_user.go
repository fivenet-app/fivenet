package documents

import (
	context "context"

	documentsrelations "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/relations"
	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
	pbdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2026/services/documents/errors"
	documentsstore "github.com/fivenet-app/fivenet/v2026/stores/documents"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

var tUserActivity = table.FivenetUserActivity

func (s *Server) ListUserDocuments(
	ctx context.Context,
	req *pbdocuments.ListUserDocumentsRequest,
) (*pbdocuments.ListUserDocumentsResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.user_id", req.GetUserId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)
	count, relations, err := s.store.ListUserDocuments(ctx, documentsstore.ListUserDocumentsQuery{
		UserID:         req.GetUserId(),
		IncludeCreated: req.GetIncludeCreated(),
		Closed:         req.Closed,
		Relations:      req.GetRelations(),
		Sort:           req.GetSort(),
		Pagination:     req.GetPagination(),
		UserInfo:       userInfo,
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	pag, _ := req.GetPagination().GetResponseWithPageSize(count.Total, 20)

	resp := &pbdocuments.ListUserDocumentsResponse{
		Pagination: pag,
		Relations:  []*documentsrelations.DocumentRelation{},
	}
	if count.Total <= 0 {
		return resp, nil
	}
	resp.Relations = relations

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetRelations() {
		if resp.GetRelations()[i].GetSourceUser() != nil {
			jobInfoFn(resp.GetRelations()[i].GetSourceUser())
		}

		if doc := resp.GetRelations()[i].GetDocument(); doc != nil &&
			doc.GetCreator() != nil {
			jobInfoFn(doc.GetCreator())
		}
	}

	return resp, nil
}

func (s *Server) addUserActivity(
	ctx context.Context,
	tx qrm.DB,
	userId int32,
	targetUserId int32,
	aType usersactivity.UserActivityType,
	reason string,
	data *usersactivity.UserActivityData,
) error {
	reasonField := mysql.NULL
	if reason != "" {
		reasonField = mysql.String(reason)
	}

	stmt := tUserActivity.
		INSERT(
			tUserActivity.SourceUserID,
			tUserActivity.TargetUserID,
			tUserActivity.Type,
			tUserActivity.Reason,
			tUserActivity.Data,
		).
		VALUES(
			userId,
			targetUserId,
			int32(aType),
			reasonField,
			data,
		)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}
