package jobs

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/notifications"
	pbjobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/jobs"
	permsjobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/jobs/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsjobs "github.com/fivenet-app/fivenet/v2025/services/jobs/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

var tConduct = table.FivenetJobConduct.AS("conduct_entry")

func (s *Server) ListConductEntries(ctx context.Context, req *pbjobs.ListConductEntriesRequest) (*pbjobs.ListConductEntriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := tConduct.Job.EQ(jet.String(userInfo.Job))

	// Field Permission Check
	fields, err := s.ps.AttrStringList(userInfo, permsjobs.ConductServicePerm, permsjobs.ConductServiceListConductEntriesPerm, permsjobs.ConductServiceListConductEntriesAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	// "All" is a pass, but if no fields or "Own" is given, return user's created conduct entries
	if fields.Contains("All") {
	} else if fields.Len() == 0 || fields.Contains("Own") {
		condition = condition.AND(tConduct.CreatorID.EQ(jet.Int32(userInfo.UserId)))
	} else {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	if len(req.Ids) > 0 {
		ids := make([]jet.Expression, len(req.Ids))
		for i := range req.Ids {
			ids[i] = jet.Uint64(req.Ids[i])
		}

		condition = condition.AND(tConduct.ID.IN(ids...))
	}
	if len(req.Types) > 0 {
		ts := make([]jet.Expression, len(req.Types))
		for i := range req.Types {
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
		for i := range req.UserIds {
			ids[i] = jet.Int32(req.UserIds[i])
		}

		condition = condition.AND(
			tConduct.TargetUserID.IN(ids...),
		)
	}

	countStmt := tConduct.
		SELECT(jet.COUNT(tConduct.ID).AS("data_count.total")).
		FROM(tConduct).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponse(count.Total)
	resp := &pbjobs.ListConductEntriesResponse{
		Pagination: pag,
	}
	if count.Total <= 0 {
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

	tColleague := tables.User().AS("target_user")
	tUserUserProps := tUserProps.AS("target_user_props")
	tColleagueAvatar := tAvatar.AS("target_user_avatar")
	tCreator := tColleague.AS("creator")
	tCreatorUserProps := tUserProps.AS("creator_props")
	tCreatorAvatar := tAvatar.AS("creator_avatar")

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
			tColleague.ID,
			tColleague.Job,
			tColleague.JobGrade,
			tColleague.Firstname,
			tColleague.Lastname,
			tColleague.Dateofbirth,
			tColleague.PhoneNumber,
			tUserUserProps.AvatarFileID.AS("target_user.avatar_file_id"),
			tColleagueAvatar.FilePath.AS("target_user.avatar"),
			tColleagueProps.UserID,
			tColleagueProps.Job,
			tColleagueProps.AbsenceBegin,
			tColleagueProps.AbsenceEnd,
			tColleagueProps.NamePrefix,
			tColleagueProps.NameSuffix,
			tConduct.CreatorID,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tCreator.PhoneNumber,
			tCreatorUserProps.AvatarFileID.AS("creator.avatar_file_id"),
			tCreatorAvatar.FilePath.AS("creator.avatar"),
		).
		FROM(
			tConduct.
				INNER_JOIN(tColleague,
					tColleague.ID.EQ(tConduct.TargetUserID),
				).
				LEFT_JOIN(tUserUserProps,
					tUserUserProps.UserID.EQ(tConduct.TargetUserID),
				).
				LEFT_JOIN(tColleagueProps,
					tColleagueProps.UserID.EQ(tConduct.TargetUserID).
						AND(tColleague.Job.EQ(jet.String(userInfo.Job))),
				).
				LEFT_JOIN(tCreator,
					tCreator.ID.EQ(tConduct.CreatorID),
				).
				LEFT_JOIN(tCreatorUserProps,
					tCreatorUserProps.UserID.EQ(tConduct.CreatorID),
				).
				LEFT_JOIN(tColleagueAvatar,
					tColleagueAvatar.ID.EQ(tUserUserProps.AvatarFileID),
				).
				LEFT_JOIN(tCreatorAvatar,
					tCreatorAvatar.ID.EQ(tCreatorUserProps.AvatarFileID),
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
	for i := range resp.Entries {
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

func (s *Server) CreateConductEntry(ctx context.Context, req *pbjobs.CreateConductEntryRequest) (*pbjobs.CreateConductEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbjobs.ConductService_ServiceDesc.ServiceName,
		Method:  "CreateConductEntry",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	req.Entry.Job = userInfo.Job

	tConduct := table.FivenetJobConduct
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

	auditEntry.State = audit.EventType_EVENT_TYPE_CREATED

	return &pbjobs.CreateConductEntryResponse{
		Entry: entry,
	}, nil
}

func (s *Server) UpdateConductEntry(ctx context.Context, req *pbjobs.UpdateConductEntryRequest) (*pbjobs.UpdateConductEntryResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.jobs.conduct_id", req.Entry.Id})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbjobs.ConductService_ServiceDesc.ServiceName,
		Method:  "UpdateConductEntry",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	entry, err := s.getConductEntry(ctx, req.Entry.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	if entry == nil || entry.Job != userInfo.Job {
		return nil, errorsjobs.ErrFailedQuery
	}

	if req.Entry.Type <= 0 {
		req.Entry.Type = entry.Type
	}
	if req.Entry.TargetUserId == 0 {
		req.Entry.TargetUserId = entry.TargetUserId
	}
	req.Entry.Job = userInfo.Job

	tConduct := table.FivenetJobConduct
	stmt := tConduct.
		UPDATE(
			tConduct.Type,
			tConduct.Message,
			tConduct.ExpiresAt,
			tConduct.TargetUserID,
		).
		SET(
			req.Entry.Type,
			req.Entry.Message,
			req.Entry.ExpiresAt,
			req.Entry.TargetUserId,
		).
		WHERE(jet.AND(
			tConduct.Job.EQ(jet.String(req.Entry.Job)),
			tConduct.ID.EQ(jet.Uint64(req.Entry.Id)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	entry, err = s.getConductEntry(ctx, entry.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	s.notifi.SendObjectEvent(ctx, &notifications.ObjectEvent{
		Type:      notifications.ObjectType_OBJECT_TYPE_JOBS_CONDUCT,
		Id:        &entry.Id,
		EventType: notifications.ObjectEventType_OBJECT_EVENT_TYPE_UPDATED,

		UserId: &userInfo.UserId,
		Job:    &userInfo.Job,
	})

	return &pbjobs.UpdateConductEntryResponse{
		Entry: entry,
	}, nil
}

func (s *Server) DeleteConductEntry(ctx context.Context, req *pbjobs.DeleteConductEntryRequest) (*pbjobs.DeleteConductEntryResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.jobs.conduct_id", req.Id})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbjobs.ConductService_ServiceDesc.ServiceName,
		Method:  "DeleteConductEntry",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	entry, err := s.getConductEntry(ctx, req.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	deletedAtTime := jet.CURRENT_TIMESTAMP()
	if entry != nil && entry.DeletedAt != nil && userInfo.Superuser {
		deletedAtTime = jet.TimestampExp(jet.NULL)
	}

	tConduct := table.FivenetJobConduct
	stmt := tConduct.
		UPDATE(
			tConduct.DeletedAt,
		).
		SET(
			tConduct.DeletedAt.SET(deletedAtTime),
		).
		WHERE(jet.AND(
			tConduct.Job.EQ(jet.String(userInfo.Job)),
			tConduct.ID.EQ(jet.Uint64(req.Id)),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return &pbjobs.DeleteConductEntryResponse{}, nil
}
