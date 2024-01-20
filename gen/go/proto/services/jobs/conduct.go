package jobs

import (
	"context"
	"errors"
	"slices"

	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	errorsjobs "github.com/galexrt/fivenet/gen/go/proto/services/jobs/errors"
	permsjobs "github.com/galexrt/fivenet/gen/go/proto/services/jobs/perms"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/errswrap"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tConduct = table.FivenetJobsConduct.AS("conduct_entry")
)

func (s *Server) ListConductEntries(ctx context.Context, req *ListConductEntriesRequest) (*ListConductEntriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := tConduct.Job.EQ(jet.String(userInfo.Job))

	// Field Permission Check
	fieldsAttr, err := s.p.Attr(userInfo, permsjobs.JobsConductServicePerm, permsjobs.JobsConductServiceListConductEntriesPerm, permsjobs.JobsConductServiceListConductEntriesAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}

	// "All" is a pass, but if no fields or "Own" is given, return user's created conduct entries
	if slices.Contains(fields, "All") {
	} else if len(fields) == 0 || slices.Contains(fields, "Own") {
		condition = condition.AND(tConduct.CreatorID.EQ(jet.Int32(userInfo.UserId)))
	} else {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	if len(req.Types) > 0 {
		ts := make([]jet.Expression, len(req.Types))
		for i := 0; i < len(req.Types); i++ {
			ts[i] = jet.Int16(int16(req.Types[i].Number()))
		}

		condition = condition.AND(
			tConduct.Type.IN(ts...),
		)
	}
	if req.ShowExpired == nil || !*req.ShowExpired {
		condition = condition.AND(jet.OR(
			tConduct.ExpiresAt.IS_NULL(),
			tConduct.ExpiresAt.GT_EQ(
				jet.CURRENT_TIMESTAMP(),
			),
		))
	}
	if len(req.UserIds) > 0 {
		ids := make([]jet.Expression, len(req.UserIds))
		for i := 0; i < len(req.UserIds); i++ {
			ids[i] = jet.Int32(req.UserIds[i])
		}

		condition = condition.AND(
			tConduct.TargetUserID.IN(ids...),
		)
	}

	countStmt := tConduct.
		SELECT(jet.COUNT(tConduct.ID).AS("datacount.totalcount")).
		FROM(tConduct).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	pag, limit := req.Pagination.GetResponse()
	resp := &ListConductEntriesResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	tUser := tUser.AS("target_user")
	tCreator := tUser.AS("creator")
	stmt := tConduct.
		SELECT(
			tConduct.ID,
			tConduct.CreatedAt,
			tConduct.UpdatedAt,
			tConduct.Job,
			tConduct.Type,
			tConduct.Message,
			tConduct.ExpiresAt,
			tConduct.TargetUserID,
			tUser.ID,
			tUser.Identifier,
			tUser.Firstname,
			tUser.Lastname,
			tUser.Dateofbirth,
			tUser.PhoneNumber,
			tConduct.CreatorID,
			tCreator.ID,
			tCreator.Identifier,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tCreator.PhoneNumber,
		).
		FROM(
			tConduct.
				INNER_JOIN(tUser,
					tUser.ID.EQ(tConduct.TargetUserID),
				).
				LEFT_JOIN(tCreator,
					tCreator.ID.EQ(tConduct.CreatorID),
				),
		).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		ORDER_BY(tConduct.CreatedAt.DESC(), tConduct.ID.DESC()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Entries); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := 0; i < len(resp.Entries); i++ {
		if resp.Entries[i].TargetUser != nil {
			jobInfoFn(resp.Entries[i].TargetUser)
		}
		if resp.Entries[i].Creator != nil {
			jobInfoFn(resp.Entries[i].Creator)
		}
	}

	resp.Pagination.Update(count.TotalCount, len(resp.Entries))

	return resp, nil
}

func (s *Server) CreateConductEntry(ctx context.Context, req *CreateConductEntryRequest) (*CreateConductEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsService_ServiceDesc.ServiceName,
		Method:  "CreateConductEntry",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	expiresAt := jet.NULL
	if req.Entry.ExpiresAt != nil {
		expiresAt = jet.TimestampT(req.Entry.ExpiresAt.AsTime())
	}

	tConduct := table.FivenetJobsConduct
	stmt := tConduct.
		INSERT(
			tConduct.Job,
			tConduct.Type,
			tConduct.Message,
			tConduct.ExpiresAt,
			tConduct.TargetUserID,
			tConduct.CreatorID,
		).
		VALUES(
			userInfo.Job,
			req.Entry.Type,
			req.Entry.Message,
			expiresAt,
			req.Entry.TargetUserId,
			userInfo.UserId,
		)

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	entry, err := s.getConductEntry(ctx, uint64(lastId))
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	return &CreateConductEntryResponse{
		Entry: entry,
	}, nil
}

func (s *Server) UpdateConductEntry(ctx context.Context, req *UpdateConductEntryRequest) (*UpdateConductEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsService_ServiceDesc.ServiceName,
		Method:  "UpdateConductEntry",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	tConduct := table.FivenetJobsConduct
	stmt := tConduct.
		UPDATE(
			tConduct.Job,
			tConduct.Type,
			tConduct.Message,
			tConduct.ExpiresAt,
			tConduct.TargetUserID,
		).
		SET(
			tConduct.Type.SET(jet.Int16(int16(req.Entry.Type))),
			tConduct.Message.SET(jet.String(req.Entry.Message)),
			tConduct.ExpiresAt.SET(jet.TimestampT(req.Entry.ExpiresAt.AsTime())),
			tConduct.TargetUserID.SET(jet.Int32(req.Entry.TargetUserId)),
		).
		WHERE(jet.AND(
			tConduct.Job.EQ(jet.String(userInfo.Job)),
			tConduct.ID.EQ(jet.Uint64(req.Entry.Id)),
		))

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	entry, err := s.getConductEntry(ctx, uint64(lastId))
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &UpdateConductEntryResponse{
		Entry: entry,
	}, nil
}

func (s *Server) DeleteConductEntry(ctx context.Context, req *DeleteConductEntryRequest) (*DeleteConductEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsService_ServiceDesc.ServiceName,
		Method:  "DeleteConductEntry",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	stmt := tConduct.
		DELETE().
		WHERE(jet.AND(
			tConduct.Job.EQ(jet.String(userInfo.Job)),
			tConduct.ID.EQ(jet.Uint64(req.Id)),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeleteConductEntryResponse{}, nil
}
