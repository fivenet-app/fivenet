package citizenstore

import (
	context "context"
	"errors"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	users "github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	errorscitizenstore "github.com/fivenet-app/fivenet/gen/go/proto/services/citizenstore/errors"
	permscompletor "github.com/fivenet-app/fivenet/gen/go/proto/services/completor/perms"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/utils"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) ManageCitizenAttributes(ctx context.Context, req *ManageCitizenAttributesRequest) (*ManageCitizenAttributesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CitizenStoreService_ServiceDesc.ServiceName,
		Method:  "ManageCitizenAttributes",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	resp := &ManageCitizenAttributesResponse{
		Attributes: []*users.CitizenAttribute{},
	}

	stmt := tJobCitizenAttributes.
		SELECT(
			tJobCitizenAttributes.ID,
			tJobCitizenAttributes.Job,
			tJobCitizenAttributes.Name,
			tJobCitizenAttributes.Color,
		).
		FROM(tJobCitizenAttributes).
		WHERE(
			tJobCitizenAttributes.Job.EQ(jet.String(userInfo.Job)),
		)

	if err := stmt.QueryContext(ctx, s.db, &resp.Attributes); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscitizenstore.ErrFailedQuery)
		}
	}

	_, removed := utils.SlicesDifferenceFunc(resp.Attributes, req.Attributes,
		func(in *users.CitizenAttribute) string {
			return in.Name
		})

	for i := 0; i < len(req.Attributes); i++ {
		req.Attributes[i].Job = &userInfo.Job
	}

	tJobCitizenAttributes := table.FivenetJobCitizenAttributes
	insertStmt := tJobCitizenAttributes.
		INSERT(
			tJobCitizenAttributes.Job,
			tJobCitizenAttributes.Name,
			tJobCitizenAttributes.Color,
		).
		MODELS(req.Attributes).
		ON_DUPLICATE_KEY_UPDATE(
			tJobCitizenAttributes.Job.SET(jet.StringExp(jet.Raw("VALUES(`job`)"))),
			tJobCitizenAttributes.Name.SET(jet.StringExp(jet.Raw("VALUES(`name`)"))),
			tJobCitizenAttributes.Color.SET(jet.StringExp(jet.Raw("VALUES(`color`)"))),
		)

	if _, err := insertStmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorscitizenstore.ErrFailedQuery)
	}

	if len(removed) > 0 {
		ids := make([]jet.Expression, len(removed))

		for i := range removed {
			ids[i] = jet.Uint64(removed[i].Id)
		}

		deleteStmt := tJobCitizenAttributes.
			DELETE().
			WHERE(jet.AND(
				tJobCitizenAttributes.ID.IN(ids...),
				tJobCitizenAttributes.Job.EQ(jet.String(userInfo.Job)),
			)).
			LIMIT(int64(len(removed)))

		if _, err := deleteStmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorscitizenstore.ErrFailedQuery)
		}
	}

	resp.Attributes = []*users.CitizenAttribute{}
	if err := stmt.QueryContext(ctx, s.db, &resp.Attributes); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscitizenstore.ErrFailedQuery)
		}
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return resp, nil
}

func (s *Server) validateCitizenAttributes(ctx context.Context, userInfo *userinfo.UserInfo, attributes []*users.CitizenAttribute) (bool, error) {
	if len(attributes) == 0 {
		return true, nil
	}

	jobsAttr, err := s.p.Attr(userInfo, permscompletor.CompletorServicePerm, permscompletor.CompletorServiceCompleteCitizenAttributesPerm, permscompletor.CompletorServiceCompleteCitizenAttributesJobsPermField)
	if err != nil {
		return false, errswrap.NewError(err, errorscitizenstore.ErrFailedQuery)
	}
	var jobs perms.StringList
	if jobsAttr != nil {
		jobs = jobsAttr.([]string)
	}

	if len(jobs) == 0 {
		jobs = append(jobs, userInfo.Job)
	}

	jobsExp := make([]jet.Expression, len(jobs))
	for i := 0; i < len(jobs); i++ {
		jobsExp[i] = jet.String(jobs[i])
	}

	idsExp := make([]jet.Expression, len(attributes))
	for i := 0; i < len(attributes); i++ {
		idsExp[i] = jet.Uint64(attributes[i].Id)
	}

	stmt := tJobCitizenAttributes.
		SELECT(
			jet.COUNT(tJobCitizenAttributes.ID).AS("datacount.totalcount"),
		).
		FROM(tJobCitizenAttributes).
		WHERE(jet.AND(
			tJobCitizenAttributes.Job.IN(jobsExp...),
			tJobCitizenAttributes.ID.IN(idsExp...),
		)).
		ORDER_BY(
			tJobCitizenAttributes.Name.DESC(),
		).
		LIMIT(10)

	dest := database.DataCount{}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return false, err
		}
	}

	return len(attributes) == int(dest.TotalCount), nil
}

func (s *Server) getUserAttributes(ctx context.Context, userInfo *userinfo.UserInfo, userId int32) (*users.CitizenAttributes, error) {
	jobsAttr, err := s.p.Attr(userInfo, permscompletor.CompletorServicePerm, permscompletor.CompletorServiceCompleteCitizenAttributesPerm, permscompletor.CompletorServiceCompleteCitizenAttributesJobsPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizenstore.ErrFailedQuery)
	}
	var jobs perms.StringList
	if jobsAttr != nil {
		jobs = jobsAttr.([]string)
	}

	if len(jobs) == 0 {
		jobs = append(jobs, userInfo.Job)
	}

	jobsExp := make([]jet.Expression, len(jobs))
	for i := 0; i < len(jobs); i++ {
		jobsExp[i] = jet.String(jobs[i])
	}

	stmt := tUserCitizenAttributes.
		SELECT(
			tJobCitizenAttributes.ID,
			tJobCitizenAttributes.Job,
			tJobCitizenAttributes.Name,
			tJobCitizenAttributes.Color,
		).
		FROM(
			tUserCitizenAttributes.
				INNER_JOIN(tJobCitizenAttributes,
					tJobCitizenAttributes.ID.EQ(tUserCitizenAttributes.AttributeID),
				),
		).
		WHERE(jet.AND(
			tUserCitizenAttributes.UserID.EQ(jet.Int32(userId)),
			tJobCitizenAttributes.Job.IN(jobsExp...),
		))

	list := &users.CitizenAttributes{}
	if err := stmt.QueryContext(ctx, s.db, list); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return list, nil
}

func (s *Server) updateCitizenAttributes(ctx context.Context, tx qrm.DB, userId int32, added []*users.CitizenAttribute, removed []*users.CitizenAttribute) error {
	tUserCitizenAttributes := table.FivenetUserCitizenAttributes

	if len(added) > 0 {
		addedAttributes := make([]*model.FivenetUserCitizenAttributes, len(added))
		for i, attribute := range added {
			addedAttributes[i] = &model.FivenetUserCitizenAttributes{
				UserID:      userId,
				AttributeID: attribute.Id,
			}
		}

		stmt := tUserCitizenAttributes.
			INSERT(
				tUserCitizenAttributes.UserID,
				tUserCitizenAttributes.AttributeID,
			).
			MODELS(addedAttributes)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return err
			}
		}
	}

	if len(removed) > 0 {
		ids := make([]jet.Expression, len(removed))

		for i := range removed {
			ids[i] = jet.Uint64(removed[i].Id)
		}

		stmt := tUserCitizenAttributes.
			DELETE().
			WHERE(jet.AND(
				tUserCitizenAttributes.UserID.EQ(jet.Int32(userId)),
				tUserCitizenAttributes.AttributeID.IN(ids...),
			)).
			LIMIT(int64(len(removed)))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}
