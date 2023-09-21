package jobs

import (
	"context"
	"errors"

	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	jobs "github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tConduct = table.FivenetJobsConduct.AS("conduct_entry")
)

func (s *Server) ConductListEntries(ctx context.Context, req *ConductListEntriesRequest) (*ConductListEntriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := tConduct.Job.EQ(jet.String(userInfo.Job))

	// Field Permission Check
	fieldsAttr, err := s.p.Attr(userInfo, JobsServicePerm, JobsServiceConductListEntriesPerm, JobsServiceConductListEntriesAccessPermField)
	if err != nil {
		return nil, ErrFailedQuery
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}

	if len(fields) == 0 {
		return nil, ErrFailedQuery
	} else if utils.InSlice(fields, "All") {
	} else if utils.InSlice(fields, "Own") {
		condition = condition.AND(tConduct.CreatorID.EQ(jet.Int32(userInfo.UserId)))
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
		return nil, ErrFailedQuery
	}

	pag, limit := req.Pagination.GetResponse()

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
		ORDER_BY(tConduct.CreatedAt.DESC(), tConduct.ID.DESC()).
		LIMIT(limit)

	resp := &ConductListEntriesResponse{
		Pagination: pag,
		Entries:    []*jobs.ConductEntry{},
	}
	if err := stmt.QueryContext(ctx, s.db, &resp.Entries); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, ErrFailedQuery
		}
	}

	resp.Pagination.Update(count.TotalCount, len(resp.Entries))

	return resp, nil
}

func (s *Server) ConductCreateEntry(ctx context.Context, req *ConductCreateEntryRequest) (*ConductCreateEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsService_ServiceDesc.ServiceName,
		Method:  "ConductCreateEntry",
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
		return nil, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	entry, err := s.getConductEntry(ctx, uint64(lastId))
	if err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	return &ConductCreateEntryResponse{
		Entry: entry,
	}, nil
}

func (s *Server) ConductUpdateEntry(ctx context.Context, req *ConductUpdateEntryRequest) (*ConductUpdateEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsService_ServiceDesc.ServiceName,
		Method:  "ConductUpdateEntry",
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
		return nil, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	entry, err := s.getConductEntry(ctx, uint64(lastId))
	if err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &ConductUpdateEntryResponse{
		Entry: entry,
	}, nil
}

func (s *Server) ConductDeleteEntry(ctx context.Context, req *ConductDeleteEntryRequest) (*ConductDeleteEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsService_ServiceDesc.ServiceName,
		Method:  "ConductDeleteEntry",
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
		return nil, err
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &ConductDeleteEntryResponse{}, nil
}
