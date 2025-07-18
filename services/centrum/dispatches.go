package centrum

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	centrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	pbcentrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/users"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	jsoniter "github.com/json-iterator/go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var (
	tDispatch        = table.FivenetCentrumDispatches.AS("dispatch")
	tDispatchStatus  = table.FivenetCentrumDispatchesStatus.AS("dispatchstatus")
	tDispatchHeatmap = table.FivenetCentrumDispatchesHeatmaps.AS("heatmap")
)

func (s *Server) ListDispatches(ctx context.Context, req *pbcentrum.ListDispatchesRequest) (*pbcentrum.ListDispatchesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "ListDispatches",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	jobs := []string{userInfo.Job}
	// TODO add other jobs based on user access
	jobsOut, _ := json.Marshal(jobs)

	condition := jet.BoolExp(dbutils.JSON_CONTAINS(tDispatch.Jobs, jet.StringExp(jet.String(string(jobsOut))))).
		AND(
			tDispatchStatus.ID.IS_NULL().OR(
				tDispatchStatus.ID.EQ(
					jet.RawInt("SELECT MAX(`dispatchstatus`.`id`) FROM `fivenet_centrum_dispatches_status` AS `dispatchstatus` WHERE `dispatchstatus`.`dispatch_id` = `dispatch`.`id`"),
				),
			),
		)

	if len(req.Status) > 0 {
		statuses := make([]jet.Expression, len(req.Status))
		for i := range req.Status {
			statuses[i] = jet.Int16(int16(*req.Status[i].Enum()))
		}

		condition = condition.AND(tDispatchStatus.Status.IN(statuses...))
	}
	if len(req.NotStatus) > 0 {
		statuses := make([]jet.Expression, len(req.NotStatus))
		for i := range req.NotStatus {
			statuses[i] = jet.Int16(int16(*req.NotStatus[i].Enum()))
		}

		condition = condition.AND(tDispatchStatus.Status.NOT_IN(statuses...))
	}

	if len(req.Ids) > 0 {
		ids := make([]jet.Expression, len(req.Ids))
		for i := range req.Ids {
			ids[i] = jet.Uint64(req.Ids[i])
		}

		condition = condition.AND(tDispatch.ID.IN(ids...))
	}

	if req.Postal != nil && *req.Postal != "" {
		condition = condition.AND(tDispatch.Postal.EQ(jet.String(*req.Postal)))
	}

	countStmt := tDispatch.
		SELECT(
			jet.COUNT(tDispatch.ID).AS("data_count.total"),
		).
		FROM(
			tDispatch.
				LEFT_JOIN(tDispatchStatus,
					tDispatchStatus.DispatchID.EQ(tDispatch.ID),
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.Total, 20)
	resp := &pbcentrum.ListDispatchesResponse{
		Pagination: pag,
	}
	if count.Total <= 0 {
		return resp, nil
	}

	tUsers := tables.User().AS("colleague")

	stmt := tDispatch.
		SELECT(
			tDispatch.ID,
			tDispatch.CreatedAt,
			tDispatch.UpdatedAt,
			tDispatch.Jobs,
			tDispatch.Message,
			tDispatch.Description,
			tDispatch.Attributes,
			tDispatch.References,
			tDispatch.X,
			tDispatch.Y,
			tDispatch.Postal,
			tDispatch.Anon,
			tDispatch.CreatorID,
			tDispatchStatus.ID,
			tDispatchStatus.CreatedAt,
			tDispatchStatus.DispatchID,
			tDispatchStatus.UnitID,
			tDispatchStatus.Status,
			tDispatchStatus.Reason,
			tDispatchStatus.Code,
			tDispatchStatus.UserID,
			tDispatchStatus.X,
			tDispatchStatus.Y,
			tDispatchStatus.Postal,
			tDispatchStatus.CreatorJob,
			tUsers.ID,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Sex,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
		).
		FROM(tDispatch.
			LEFT_JOIN(tDispatchStatus,
				tDispatchStatus.DispatchID.EQ(tDispatch.ID),
			).
			LEFT_JOIN(tUsers,
				tUsers.ID.EQ(tDispatchStatus.UserID),
			)).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		ORDER_BY(
			tDispatch.ID.DESC(),
		).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Dispatches); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	resp.Pagination.Update(len(resp.Dispatches))

	publicJobs := s.appCfg.Get().JobInfo.PublicJobs
	for i := range resp.Dispatches {
		var err error
		resp.Dispatches[i].Units, err = s.dispatches.LoadDispatchAssignments(ctx, resp.Dispatches[i].Id)
		if err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		if resp.Dispatches[i].CreatorId != nil {
			resp.Dispatches[i].Creator, err = users.RetrieveUserById(ctx, s.db, *resp.Dispatches[i].CreatorId)
			if err != nil {
				return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}

			if resp.Dispatches[i].Creator != nil {
				// Clear dispatch creator's job info if it isn't a public job
				if !slices.Contains(publicJobs, resp.Dispatches[i].Creator.Job) {
					resp.Dispatches[i].Creator.Job = ""
				}
				resp.Dispatches[i].Creator.JobGrade = 0
			}
		}

		// Ensure dispatch has a valid job list (fallback to deprecated Jobs field for old dispatches)
		if resp.Dispatches[i].Jobs == nil || len(resp.Dispatches[i].Jobs.GetJobs()) == 0 {
			resp.Dispatches[i].Jobs = &centrum.JobList{
				Jobs: []*centrum.Job{
					{
						Name: resp.Dispatches[i].Job,
					},
				},
			}
			resp.Dispatches[i].Job = ""
		}
		for _, job := range resp.Dispatches[i].Jobs.GetJobs() {
			s.enricher.EnrichJobName(job)
		}
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_VIEWED

	return resp, nil
}

func (s *Server) GetDispatch(ctx context.Context, req *pbcentrum.GetDispatchRequest) (*pbcentrum.GetDispatchResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.centrum.dispatch_id", int64(req.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "GetDispatch",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	jobs, _, err := s.settings.GetJobAccessList(ctx, userInfo.Job, userInfo.JobGrade)
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}
	jobsOut, _ := json.Marshal(jobs)

	condition := tDispatchStatus.ID.IS_NULL().OR(
		tDispatchStatus.ID.EQ(
			jet.RawInt("SELECT MAX(`dispatchstatus`.`id`) FROM `fivenet_centrum_dispatches_status` AS `dispatchstatus` WHERE `dispatchstatus`.`dispatch_id` = `dispatch`.`id`"),
		),
	).
		AND(tDispatch.ID.EQ(jet.Uint64(req.Id))).
		AND(jet.BoolExp(dbutils.JSON_CONTAINS(tDispatch.Jobs, jet.StringExp(jet.String(string(jobsOut))))))

	resp := &pbcentrum.GetDispatchResponse{
		Dispatch: &centrum.Dispatch{},
	}

	tUsers := tables.User().AS("colleague")

	stmt := tDispatch.
		SELECT(
			tDispatch.ID,
			tDispatch.CreatedAt,
			tDispatch.UpdatedAt,
			tDispatch.Jobs,
			tDispatch.Message,
			tDispatch.Description,
			tDispatch.Attributes,
			tDispatch.References,
			tDispatch.X,
			tDispatch.Y,
			tDispatch.Postal,
			tDispatch.Anon,
			tDispatch.CreatorID,
			tDispatchStatus.ID,
			tDispatchStatus.CreatedAt,
			tDispatchStatus.DispatchID,
			tDispatchStatus.UnitID,
			tDispatchStatus.Status,
			tDispatchStatus.Reason,
			tDispatchStatus.Code,
			tDispatchStatus.UserID,
			tDispatchStatus.X,
			tDispatchStatus.Y,
			tDispatchStatus.Postal,
			tDispatchStatus.CreatorJob,
			tUsers.ID,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Sex,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
		).
		FROM(tDispatch.
			LEFT_JOIN(tDispatchStatus,
				tDispatchStatus.DispatchID.EQ(tDispatch.ID),
			).
			LEFT_JOIN(tUsers,
				tUsers.ID.EQ(tDispatchStatus.UserID),
			)).
		WHERE(condition).
		ORDER_BY(
			tDispatch.ID.DESC(),
		).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, s.db, resp.Dispatch); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	if resp.Dispatch == nil || resp.Dispatch.Id <= 0 {
		return &pbcentrum.GetDispatchResponse{}, nil
	}

	resp.Dispatch.Units, err = s.dispatches.LoadDispatchAssignments(ctx, resp.Dispatch.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	if resp.Dispatch.CreatorId != nil && *resp.Dispatch.CreatorId > 0 {
		creator, err := users.RetrieveUserById(ctx, s.db, *resp.Dispatch.CreatorId)
		if err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		if creator != nil {
			resp.Dispatch.Creator = creator
			// Clear dispatch creator's job info if not a visible job
			if !slices.Contains(s.appCfg.Get().JobInfo.PublicJobs, creator.Job) {
				resp.Dispatch.Creator.Job = ""
			}
			resp.Dispatch.Creator.JobGrade = 0
		}
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_VIEWED

	return resp, nil
}

func (s *Server) CreateDispatch(ctx context.Context, req *pbcentrum.CreateDispatchRequest) (*pbcentrum.CreateDispatchResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "CreateDispatch",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	// Make sure jobs and creator id are set
	if len(req.Dispatch.Jobs.GetJobs()) > 0 {
		for _, job := range req.Dispatch.Jobs.GetJobStrings() {
			if !s.jobs.Has(job) {
				return nil, errorscentrum.ErrFailedQuery
			}
		}
	} else {
		req.Dispatch.Jobs = &centrum.JobList{
			Jobs: []*centrum.Job{
				{
					Name: userInfo.Job,
				},
			},
		}
	}
	req.Dispatch.CreatorId = &userInfo.UserId
	req.Dispatch.CreatedAt = timestamp.Now()

	dsp, err := s.dispatches.Create(ctx, req.Dispatch)
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_CREATED

	return &pbcentrum.CreateDispatchResponse{
		Dispatch: dsp,
	}, nil
}

func (s *Server) UpdateDispatch(ctx context.Context, req *pbcentrum.UpdateDispatchRequest) (*pbcentrum.UpdateDispatchResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.centrum.dispatch_id", int64(req.Dispatch.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "UpdateDispatch",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	dsp, err := s.dispatches.Update(ctx, &userInfo.UserId, req.Dispatch)
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbcentrum.UpdateDispatchResponse{
		Dispatch: dsp,
	}, nil
}

func (s *Server) TakeDispatch(ctx context.Context, req *pbcentrum.TakeDispatchRequest) (*pbcentrum.TakeDispatchResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64Slice("fivenet.centrum.dispatch_ids", utils.SliceUint64ToInt64(req.DispatchIds)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "TakeDispatch",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	unitMapping, err := s.tracker.GetUserMapping(userInfo.UserId)
	if err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	if unitMapping == nil || unitMapping.UnitId == nil || *unitMapping.UnitId <= 0 {
		return nil, errorscentrum.ErrNotOnDuty
	}

	if err := s.dispatches.TakeDispatch(ctx, userInfo.Job, userInfo.UserId, *unitMapping.UnitId, req.Resp, req.DispatchIds); err != nil {
		return nil, err
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbcentrum.TakeDispatchResponse{}, nil
}

func (s *Server) UpdateDispatchStatus(ctx context.Context, req *pbcentrum.UpdateDispatchStatusRequest) (*pbcentrum.UpdateDispatchStatusResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "UpdateDispatchStatus",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	dsp, err := s.dispatches.Get(ctx, req.DispatchId)
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	if !s.helpers.CheckIfUserIsPartOfDispatch(ctx, userInfo, dsp, true) && !userInfo.Superuser {
		return nil, errorscentrum.ErrNotPartOfDispatch
	}

	var statusUnitId *uint64
	userMapping, err := s.tracker.GetUserMapping(userInfo.UserId)
	if err != nil {
		if !s.helpers.CheckIfUserIsDispatcher(ctx, userInfo.Job, userInfo.UserId) {
			return nil, errorscentrum.ErrNotPartOfDispatch
		}
	} else {
		statusUnitId = userMapping.UnitId
	}

	if _, err := s.dispatches.UpdateStatus(ctx, dsp.Id, &centrum.DispatchStatus{
		CreatedAt:  timestamp.Now(),
		DispatchId: dsp.Id,
		UnitId:     statusUnitId,
		Status:     req.Status,
		Code:       req.Code,
		Reason:     req.Reason,
		UserId:     &userInfo.UserId,
		CreatorJob: &userInfo.Job,
	}); err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	if (req.Status == centrum.StatusDispatch_STATUS_DISPATCH_EN_ROUTE ||
		req.Status == centrum.StatusDispatch_STATUS_DISPATCH_ON_SCENE ||
		req.Status == centrum.StatusDispatch_STATUS_DISPATCH_NEED_ASSISTANCE) && statusUnitId != nil {
		if unit, err := s.units.Get(ctx, *statusUnitId); err == nil {
			// Set unit to busy when unit accepts a dispatch
			if unit.Status == nil || unit.Status.Status != centrum.StatusUnit_STATUS_UNIT_BUSY {
				if _, err := s.units.UpdateStatus(ctx, *statusUnitId, &centrum.UnitStatus{
					CreatedAt:  timestamp.Now(),
					UnitId:     unit.Id,
					Status:     centrum.StatusUnit_STATUS_UNIT_BUSY,
					UserId:     &userInfo.UserId,
					CreatorId:  &userInfo.UserId,
					CreatorJob: &userInfo.Job,
				}); err != nil {
					return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
				}
			}
		}
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbcentrum.UpdateDispatchStatusResponse{}, nil
}

func (s *Server) AssignDispatch(ctx context.Context, req *pbcentrum.AssignDispatchRequest) (*pbcentrum.AssignDispatchResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.centrum.dispatch_id", int64(req.DispatchId)))
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64Slice("fivenet.centrum.units.to_add", utils.SliceUint64ToInt64(req.ToAdd)))
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64Slice("fivenet.centrum.units.to_remove", utils.SliceUint64ToInt64(req.ToRemove)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "AssignDispatch",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	dsp, err := s.dispatches.Get(ctx, req.DispatchId)
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	if !slices.Contains(dsp.Jobs.GetJobStrings(), userInfo.Job) {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	expiresAt := time.Time{}
	if req.Forced == nil || !*req.Forced {
		expiresAt = s.settings.DispatchAssignmentExpirationTime()
	}

	if err := s.dispatches.UpdateAssignments(ctx, &userInfo.UserId, dsp.Id, req.ToAdd, req.ToRemove, expiresAt); err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbcentrum.AssignDispatchResponse{}, nil
}

func (s *Server) ListDispatchActivity(ctx context.Context, req *pbcentrum.ListDispatchActivityRequest) (*pbcentrum.ListDispatchActivityResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	countStmt := tDispatchStatus.
		SELECT(
			jet.COUNT(jet.DISTINCT(tDispatchStatus.ID)).AS("data_count.total"),
		).
		FROM(
			tDispatchStatus.
				INNER_JOIN(tDispatch,
					tDispatch.ID.EQ(tDispatchStatus.DispatchID),
				),
		).
		WHERE(jet.AND(
			tDispatchStatus.DispatchID.EQ(jet.Uint64(req.Id)),
		))

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.Total, 10)
	resp := &pbcentrum.ListDispatchActivityResponse{
		Pagination: pag,
	}
	if count.Total <= 0 {
		return resp, nil
	}

	tUsers := tables.User().AS("colleague")
	tAvatar := table.FivenetFiles.AS("avatar")

	stmt := tDispatchStatus.
		SELECT(
			tDispatchStatus.ID,
			tDispatchStatus.CreatedAt,
			tDispatchStatus.UnitID,
			tDispatchStatus.Status,
			tDispatchStatus.Reason,
			tDispatchStatus.Code,
			tDispatchStatus.UserID,
			tDispatchStatus.X,
			tDispatchStatus.Y,
			tDispatchStatus.Postal,
			tDispatchStatus.CreatorJob,
			tUsers.ID,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Sex,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
			tColleagueProps.UserID,
			tColleagueProps.Job,
			tColleagueProps.NamePrefix,
			tColleagueProps.NameSuffix,
			tUserProps.AvatarFileID.AS("colleague.avatar_file_id"),
			tAvatar.FilePath.AS("colleague.avatar"),
		).
		FROM(
			tDispatchStatus.
				LEFT_JOIN(tUsers,
					tUsers.ID.EQ(tDispatchStatus.UserID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tDispatchStatus.UserID).
						AND(tUsers.Job.EQ(jet.String(userInfo.Job))),
				).
				LEFT_JOIN(tColleagueProps,
					tColleagueProps.UserID.EQ(tUsers.ID).
						AND(tColleagueProps.Job.EQ(tUsers.Job)),
				).
				LEFT_JOIN(tAvatar,
					tAvatar.ID.EQ(tUserProps.AvatarFileID),
				),
		).
		WHERE(
			tDispatchStatus.DispatchID.EQ(jet.Uint64(req.Id)),
		).
		ORDER_BY(tDispatchStatus.ID.DESC()).
		OFFSET(req.Pagination.Offset).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Activity); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.Activity {
		if resp.Activity[i].UnitId != nil && *resp.Activity[i].UnitId > 0 {
			var err error
			resp.Activity[i].Unit, err = s.units.Get(ctx, *resp.Activity[i].UnitId)
			if err != nil {
				return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}
		}

		if resp.Activity[i].User != nil {
			jobInfoFn(resp.Activity[i].User)
		}
	}

	resp.Pagination.Update(len(resp.Activity))

	return resp, nil
}

func (s *Server) DeleteDispatch(ctx context.Context, req *pbcentrum.DeleteDispatchRequest) (*pbcentrum.DeleteDispatchResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "DeleteDispatch",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	dsp, err := s.dispatches.Get(ctx, req.Id)
	if err != nil {
		if !errors.Is(err, jetstream.ErrKeyNotFound) {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	if !userInfo.Superuser {
		if dsp == nil || !slices.Contains(dsp.Jobs.GetJobStrings(), userInfo.Job) {
			return nil, errorscentrum.ErrNotPartOfDispatch
		}
	}

	if err := s.dispatches.Delete(ctx, req.Id, true); err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return &pbcentrum.DeleteDispatchResponse{}, nil
}
