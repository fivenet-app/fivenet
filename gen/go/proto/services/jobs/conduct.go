package jobs

import (
	"context"
	"errors"
	"slices"

	database "github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	errorsjobs "github.com/fivenet-app/fivenet/gen/go/proto/services/jobs/errors"
	permsjobs "github.com/fivenet-app/fivenet/gen/go/proto/services/jobs/perms"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tConduct = table.FivenetJobsConduct.AS("conduct_entry")

func (s *Server) ListConductEntries(ctx context.Context, req *ListConductEntriesRequest) (*ListConductEntriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := tConduct.Job.EQ(jet.String(userInfo.Job))

	// Field Permission Check
	fieldsAttr, err := s.ps.Attr(userInfo, permsjobs.JobsConductServicePerm, permsjobs.JobsConductServiceListConductEntriesPerm, permsjobs.JobsConductServiceListConductEntriesAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
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
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	if len(req.Ids) > 0 {
		ids := make([]jet.Expression, len(req.Ids))
		for i := 0; i < len(req.Ids); i++ {
			ids[i] = jet.Uint64(req.Ids[i])
		}

		condition = condition.AND(tConduct.ID.IN(ids...))
	}
	if len(req.Types) > 0 {
		ts := make([]jet.Expression, len(req.Types))
		for i := 0; i < len(req.Types); i++ {
			ts[i] = jet.Int16(int16(req.Types[i].Number()))
		}

		condition = condition.AND(tConduct.Type.IN(ts...))
	}
	if len(req.Ids) == 0 && (req.ShowExpired == nil || !*req.ShowExpired) {
		condition = condition.AND(jet.OR(
			tConduct.ExpiresAt.IS_NULL(),
			tConduct.ExpiresAt.GT_EQ(
				jet.CURRENT_DATE(),
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
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponse(count.TotalCount)
	resp := &ListConductEntriesResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	// Convert proto sort to db sorting
	orderBys := []jet.OrderByClause{}
	if req.Sort != nil {
		var columns []jet.Column
		switch req.Sort.Column {
		case "type":
			columns = append(columns, tConduct.Type, tConduct.ID)
		case "id":
			fallthrough
		default:
			columns = append(columns, tConduct.ID)
		}

		for _, column := range columns {
			if req.Sort.Direction == database.AscSortDirection {
				orderBys = append(orderBys, column.ASC())
			} else {
				orderBys = append(orderBys, column.DESC())
			}
		}
	} else {
		orderBys = append(orderBys, tConduct.ID.DESC())
	}

	tUser := tUser.AS("target_user")
	tUserUserProps := tUserProps.AS("target_user_props")
	tCreator := tUser.AS("creator")
	tCreatorUserProps := tUserProps.AS("creator_props")
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
			tUser.Job,
			tUser.JobGrade,
			tUser.Firstname,
			tUser.Lastname,
			tUser.Dateofbirth,
			tUser.PhoneNumber,
			tUserUserProps.Avatar.AS("target_user.avatar"),
			tJobsUserProps.UserID,
			tJobsUserProps.Job,
			tJobsUserProps.AbsenceBegin,
			tJobsUserProps.AbsenceEnd,
			tConduct.CreatorID,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tCreator.PhoneNumber,
			tCreatorUserProps.Avatar.AS("creator.avatar"),
		).
		FROM(
			tConduct.
				INNER_JOIN(tUser,
					tUser.ID.EQ(tConduct.TargetUserID),
				).
				LEFT_JOIN(tUserUserProps,
					tUserUserProps.UserID.EQ(tConduct.TargetUserID),
				).
				LEFT_JOIN(tJobsUserProps,
					tJobsUserProps.UserID.EQ(tConduct.TargetUserID).
						AND(tUser.Job.EQ(jet.String(userInfo.Job))),
				).
				LEFT_JOIN(tCreator,
					tCreator.ID.EQ(tConduct.CreatorID),
				).
				LEFT_JOIN(tCreatorUserProps,
					tCreatorUserProps.UserID.EQ(tConduct.CreatorID),
				),
		).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		ORDER_BY(orderBys...).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Entries); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
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

	resp.Pagination.Update(len(resp.Entries))

	return resp, nil
}

func (s *Server) CreateConductEntry(ctx context.Context, req *CreateConductEntryRequest) (*CreateConductEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsConductService_ServiceDesc.ServiceName,
		Method:  "CreateConductEntry",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	req.Entry.Job = userInfo.Job

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
			req.Entry.ExpiresAt,
			req.Entry.TargetUserId,
			userInfo.UserId,
		)

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	req.Entry.Id = uint64(lastId)

	entry, err := s.getConductEntry(ctx, req.Entry.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	return &CreateConductEntryResponse{
		Entry: entry,
	}, nil
}

func (s *Server) UpdateConductEntry(ctx context.Context, req *UpdateConductEntryRequest) (*UpdateConductEntryResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.jobs.conduct.id", int64(req.Entry.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsConductService_ServiceDesc.ServiceName,
		Method:  "UpdateConductEntry",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	entry, err := s.getConductEntry(ctx, req.Entry.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	if req.Entry.Type == 0 {
		req.Entry.Type = entry.Type
	}
	if req.Entry.ExpiresAt == nil {
		req.Entry.ExpiresAt = entry.ExpiresAt
	}
	if req.Entry.TargetUser == nil {
		req.Entry.TargetUser = entry.TargetUser
	}
	if req.Entry.TargetUserId == 0 {
		req.Entry.TargetUserId = entry.TargetUserId
	}

	req.Entry.Job = userInfo.Job

	tConduct := table.FivenetJobsConduct
	stmt := tConduct.
		UPDATE(
			tConduct.Type,
			tConduct.Message,
			tConduct.ExpiresAt,
			tConduct.TargetUserID,
		).
		SET(
			int16(entry.Type),
			req.Entry.Message,
			req.Entry.ExpiresAt,
			req.Entry.TargetUserId,
		).
		WHERE(jet.AND(
			tConduct.Job.EQ(jet.String(userInfo.Job)),
			tConduct.ID.EQ(jet.Uint64(req.Entry.Id)),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	entry, err = s.getConductEntry(ctx, entry.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &UpdateConductEntryResponse{
		Entry: entry,
	}, nil
}

func (s *Server) DeleteConductEntry(ctx context.Context, req *DeleteConductEntryRequest) (*DeleteConductEntryResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.jobs.conduct.id", int64(req.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsConductService_ServiceDesc.ServiceName,
		Method:  "DeleteConductEntry",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	stmt := tConduct.
		DELETE().
		WHERE(jet.AND(
			tConduct.Job.EQ(jet.String(userInfo.Job)),
			tConduct.ID.EQ(jet.Uint64(req.Id)),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeleteConductEntryResponse{}, nil
}
