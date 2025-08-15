package centrum

import (
	"context"
	"encoding/json"
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
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/nats-io/nats.go/jetstream"
)

var (
	tDispatch        = table.FivenetCentrumDispatches.AS("dispatch")
	tDispatchStatus  = table.FivenetCentrumDispatchesStatus.AS("dispatchstatus")
	tDispatchHeatmap = table.FivenetCentrumDispatchesHeatmaps.AS("heatmap")
)

func (s *Server) ListDispatches(
	ctx context.Context,
	req *pbcentrum.ListDispatchesRequest,
) (*pbcentrum.ListDispatchesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "ListDispatches",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	jobs, _, err := s.settings.GetJobAccessList(ctx, userInfo.GetJob(), userInfo.GetJobGrade())
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}
	jobsOut, _ := json.Marshal(jobs)

	condition := jet.BoolExp(dbutils.JSON_CONTAINS(tDispatch.Jobs, jet.StringExp(jet.String(string(jobsOut))))).
		AND(
			tDispatchStatus.ID.IS_NULL().OR(
				tDispatchStatus.ID.EQ(
					jet.RawInt("SELECT MAX(`dispatchstatus`.`id`) FROM `fivenet_centrum_dispatches_status` AS `dispatchstatus` WHERE `dispatchstatus`.`dispatch_id` = `dispatch`.`id`"),
				),
			),
		)

	if len(req.GetStatus()) > 0 {
		statuses := make([]jet.Expression, len(req.GetStatus()))
		for i := range req.GetStatus() {
			statuses[i] = jet.Int32(int32(*req.GetStatus()[i].Enum()))
		}

		condition = condition.AND(tDispatchStatus.Status.IN(statuses...))
	}
	if len(req.GetNotStatus()) > 0 {
		statuses := make([]jet.Expression, len(req.GetNotStatus()))
		for i := range req.GetNotStatus() {
			statuses[i] = jet.Int32(int32(*req.GetNotStatus()[i].Enum()))
		}

		condition = condition.AND(tDispatchStatus.Status.NOT_IN(statuses...))
	}

	if len(req.GetIds()) > 0 {
		ids := make([]jet.Expression, len(req.GetIds()))
		for i := range req.GetIds() {
			ids[i] = jet.Int64(req.GetIds()[i])
		}

		condition = condition.AND(tDispatch.ID.IN(ids...))
	}

	if req.Postal != nil && req.GetPostal() != "" {
		condition = condition.AND(tDispatch.Postal.EQ(jet.String(req.GetPostal())))
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

	pag, limit := req.GetPagination().GetResponseWithPageSize(count.Total, 20)
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
		OFFSET(req.GetPagination().GetOffset()).
		ORDER_BY(
			tDispatch.ID.DESC(),
		).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Dispatches); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	dsps := resp.GetDispatches()
	resp.GetPagination().Update(len(dsps))

	publicJobs := s.appCfg.Get().JobInfo.GetPublicJobs()
	for i := range dsps {
		var err error
		resp.Dispatches[i].Units, err = s.dispatches.LoadDispatchAssignments(
			ctx,
			dsps[i].GetId(),
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		if resp.Dispatches[i].CreatorId != nil {
			resp.Dispatches[i].Creator, err = users.RetrieveUserById(
				ctx,
				s.db,
				dsps[i].GetCreatorId(),
			)
			if err != nil {
				return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}

			if dsps[i].GetCreator() != nil {
				// Clear dispatch creator's job info if it isn't a public job
				if !slices.Contains(publicJobs, dsps[i].GetCreator().GetJob()) {
					resp.Dispatches[i].Creator.Job = ""
				}
				resp.Dispatches[i].Creator.JobGrade = 0
			}
		}

		// Ensure dispatch has a valid job list (fallback to deprecated Jobs field for old dispatches)
		if dsps[i].GetJobs() == nil ||
			len(dsps[i].GetJobs().GetJobs()) == 0 {
			resp.Dispatches[i].Jobs = &centrum.JobList{
				Jobs: []*centrum.Job{
					{
						//nolint:staticcheck // This is a fallback for old dispatches.
						Name: dsps[i].GetJob(),
					},
				},
			}
			//nolint:staticcheck // Clear old job info. This is a fallback for old dispatches.
			resp.Dispatches[i].Job = ""
		}
		for _, job := range dsps[i].GetJobs().GetJobs() {
			s.enricher.EnrichJobName(job)
		}
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_VIEWED

	return resp, nil
}

func (s *Server) GetDispatch(
	ctx context.Context,
	req *pbcentrum.GetDispatchRequest,
) (*pbcentrum.GetDispatchResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.centrum.dispatch_id", req.GetId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "GetDispatch",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	jobs, _, err := s.settings.GetJobAccessList(ctx, userInfo.GetJob(), userInfo.GetJobGrade())
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}
	jobsOut, _ := json.Marshal(jobs)

	condition := tDispatchStatus.ID.IS_NULL().OR(
		tDispatchStatus.ID.EQ(
			jet.RawInt("SELECT MAX(`dispatchstatus`.`id`) FROM `fivenet_centrum_dispatches_status` AS `dispatchstatus` WHERE `dispatchstatus`.`dispatch_id` = `dispatch`.`id`"),
		),
	).
		AND(tDispatch.ID.EQ(jet.Int64(req.GetId()))).
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

	if err := stmt.QueryContext(ctx, s.db, resp.GetDispatch()); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	if resp.GetDispatch() == nil || resp.GetDispatch().GetId() <= 0 {
		return &pbcentrum.GetDispatchResponse{}, nil
	}

	resp.Dispatch.Units, err = s.dispatches.LoadDispatchAssignments(ctx, resp.GetDispatch().GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	if resp.Dispatch.CreatorId != nil && resp.GetDispatch().GetCreatorId() > 0 {
		creator, err := users.RetrieveUserById(ctx, s.db, resp.GetDispatch().GetCreatorId())
		if err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		if creator != nil {
			resp.Dispatch.Creator = creator
			// Clear dispatch creator's job info if not a visible job
			if !slices.Contains(s.appCfg.Get().JobInfo.GetPublicJobs(), creator.GetJob()) {
				resp.Dispatch.Creator.Job = ""
			}
			resp.Dispatch.Creator.JobGrade = 0
		}
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_VIEWED

	return resp, nil
}

func (s *Server) CreateDispatch(
	ctx context.Context,
	req *pbcentrum.CreateDispatchRequest,
) (*pbcentrum.CreateDispatchResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "CreateDispatch",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	// Make sure jobs and creator id are set
	if len(req.GetDispatch().GetJobs().GetJobs()) > 0 {
		for _, job := range req.GetDispatch().GetJobs().GetJobStrings() {
			if !s.jobs.Has(job) {
				return nil, errorscentrum.ErrFailedQuery
			}
		}
	} else {
		req.Dispatch.Jobs = &centrum.JobList{
			Jobs: []*centrum.Job{
				{
					Name: userInfo.GetJob(),
				},
			},
		}
	}
	req.Dispatch.CreatorId = &userInfo.UserId
	req.Dispatch.CreatedAt = timestamp.Now()

	dsp, err := s.dispatches.Create(ctx, req.GetDispatch())
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_CREATED

	return &pbcentrum.CreateDispatchResponse{
		Dispatch: dsp,
	}, nil
}

func (s *Server) UpdateDispatch(
	ctx context.Context,
	req *pbcentrum.UpdateDispatchRequest,
) (*pbcentrum.UpdateDispatchResponse, error) {
	logging.InjectFields(
		ctx,
		logging.Fields{"fivenet.centrum.dispatch_id", req.GetDispatch().GetId()},
	)

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "UpdateDispatch",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	dsp, err := s.dispatches.Update(ctx, &userInfo.UserId, req.GetDispatch())
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbcentrum.UpdateDispatchResponse{
		Dispatch: dsp,
	}, nil
}

func (s *Server) TakeDispatch(
	ctx context.Context,
	req *pbcentrum.TakeDispatchRequest,
) (*pbcentrum.TakeDispatchResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.centrum.dispatch_ids", req.GetDispatchIds()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "TakeDispatch",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	unitMapping, err := s.tracker.GetUserMapping(userInfo.GetUserId())
	if err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	if unitMapping == nil || unitMapping.UnitId == nil || unitMapping.GetUnitId() <= 0 {
		return nil, errorscentrum.ErrNotOnDuty
	}

	if err := s.dispatches.TakeDispatch(ctx, userInfo.GetJob(), userInfo.GetUserId(), unitMapping.GetUnitId(), req.GetResp(), req.GetDispatchIds()); err != nil {
		return nil, err
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbcentrum.TakeDispatchResponse{}, nil
}

func (s *Server) UpdateDispatchStatus(
	ctx context.Context,
	req *pbcentrum.UpdateDispatchStatusRequest,
) (*pbcentrum.UpdateDispatchStatusResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "UpdateDispatchStatus",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	dsp, err := s.dispatches.Get(ctx, req.GetDispatchId())
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	if !s.helpers.CheckIfUserIsPartOfDispatch(ctx, userInfo, dsp, true) &&
		!userInfo.GetSuperuser() {
		return nil, errorscentrum.ErrNotPartOfDispatch
	}

	var statusUnitId *int64
	userMapping, err := s.tracker.GetUserMapping(userInfo.GetUserId())
	if err != nil {
		if !s.helpers.CheckIfUserIsDispatcher(ctx, userInfo.GetJob(), userInfo.GetUserId()) {
			return nil, errorscentrum.ErrNotPartOfDispatch
		}
	} else {
		statusUnitId = userMapping.UnitId
	}

	if _, err := s.dispatches.UpdateStatus(ctx, dsp.GetId(), &centrum.DispatchStatus{
		CreatedAt:  timestamp.Now(),
		DispatchId: dsp.GetId(),
		UnitId:     statusUnitId,
		Status:     req.GetStatus(),
		Code:       req.Code,
		Reason:     req.Reason,
		UserId:     &userInfo.UserId,
		CreatorJob: &userInfo.Job,
	}); err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	if (req.GetStatus() == centrum.StatusDispatch_STATUS_DISPATCH_EN_ROUTE ||
		req.GetStatus() == centrum.StatusDispatch_STATUS_DISPATCH_ON_SCENE ||
		req.GetStatus() == centrum.StatusDispatch_STATUS_DISPATCH_NEED_ASSISTANCE) && statusUnitId != nil {
		if unit, err := s.units.Get(ctx, *statusUnitId); err == nil {
			// Set unit to busy when unit accepts a dispatch
			if unit.GetStatus() == nil ||
				unit.GetStatus().GetStatus() != centrum.StatusUnit_STATUS_UNIT_BUSY {
				if _, err := s.units.UpdateStatus(ctx, *statusUnitId, &centrum.UnitStatus{
					CreatedAt:  timestamp.Now(),
					UnitId:     unit.GetId(),
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

func (s *Server) AssignDispatch(
	ctx context.Context,
	req *pbcentrum.AssignDispatchRequest,
) (*pbcentrum.AssignDispatchResponse, error) {
	logging.InjectFields(ctx, logging.Fields{
		"fivenet.centrum.dispatch_id", req.GetDispatchId(),
		"fivenet.centrum.units.to_add", req.GetToAdd(),
		"fivenet.centrum.units.to_remove", req.GetToRemove(),
	})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "AssignDispatch",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	dsp, err := s.dispatches.Get(ctx, req.GetDispatchId())
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	if !slices.Contains(dsp.GetJobs().GetJobStrings(), userInfo.GetJob()) {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	expiresAt := time.Time{}
	if req.Forced == nil || !req.GetForced() {
		expiresAt = s.settings.DispatchAssignmentExpirationTime()
	}

	if err := s.dispatches.UpdateAssignments(ctx, &userInfo.UserId, dsp.GetId(), req.GetToAdd(), req.GetToRemove(), expiresAt); err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbcentrum.AssignDispatchResponse{}, nil
}

func (s *Server) ListDispatchActivity(
	ctx context.Context,
	req *pbcentrum.ListDispatchActivityRequest,
) (*pbcentrum.ListDispatchActivityResponse, error) {
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
			tDispatchStatus.DispatchID.EQ(jet.Int64(req.GetId())),
		))

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count.Total, 10)
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
						AND(tUsers.Job.EQ(jet.String(userInfo.GetJob()))),
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
			tDispatchStatus.DispatchID.EQ(jet.Int64(req.GetId())),
		).
		ORDER_BY(tDispatchStatus.ID.DESC()).
		OFFSET(req.GetPagination().GetOffset()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Activity); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetActivity() {
		if resp.Activity[i].UnitId != nil && resp.GetActivity()[i].GetUnitId() > 0 {
			var err error
			resp.Activity[i].Unit, err = s.units.Get(ctx, resp.GetActivity()[i].GetUnitId())
			if err != nil {
				return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}
		}

		if resp.GetActivity()[i].GetUser() != nil {
			jobInfoFn(resp.GetActivity()[i].GetUser())
		}
	}

	resp.GetPagination().Update(len(resp.GetActivity()))

	return resp, nil
}

func (s *Server) DeleteDispatch(
	ctx context.Context,
	req *pbcentrum.DeleteDispatchRequest,
) (*pbcentrum.DeleteDispatchResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "DeleteDispatch",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	dsp, err := s.dispatches.Get(ctx, req.GetId())
	if err != nil {
		if !errors.Is(err, jetstream.ErrKeyNotFound) {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	if !userInfo.GetSuperuser() {
		if dsp == nil || !slices.Contains(dsp.GetJobs().GetJobStrings(), userInfo.GetJob()) {
			return nil, errorscentrum.ErrNotPartOfDispatch
		}
	}

	if err := s.dispatches.Delete(ctx, req.GetId(), true); err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return &pbcentrum.DeleteDispatchResponse{}, nil
}
