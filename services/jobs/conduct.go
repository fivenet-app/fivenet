package jobs

import (
	"context"
	"strconv"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	jobsconduct "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/conduct"
	notificationsclientview "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications/clientview"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs"
	permsjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorsjobs "github.com/fivenet-app/fivenet/v2026/services/jobs/errors"
	jobsstore "github.com/fivenet-app/fivenet/v2026/stores/jobs"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) ListConductEntries(
	ctx context.Context,
	req *pbjobs.ListConductEntriesRequest,
) (*pbjobs.ListConductEntriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Field Permission Check
	fields, err := permsjobs.ConductService.ListConductEntries.AccessTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	allAccess := fields.Contains(permsjobs.ConductServiceListConductEntriesAccessPermValueAll)
	ownOnly := !allAccess &&
		(fields.Len() == 0 || fields.Contains(permsjobs.ConductServiceListConductEntriesAccessPermValueOwn))
	if !allAccess && !ownOnly {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	q := jobsstore.ConductQuery{
		Job:         userInfo.GetJob(),
		Sort:        req.GetSort(),
		Offset:      req.GetPagination().GetOffset(),
		Limit:       0,
		Types:       req.GetTypes(),
		ShowExpired: req.GetShowExpired(),
		ShowDrafts:  req.GetShowDrafts(),
		UserIDs:     req.GetUserIds(),
		IDs:         req.GetIds(),
		CreatorID:   userInfo.GetUserId(),
		OwnOnly:     ownOnly,
		AllAccess:   allAccess,
	}

	total, err := s.store.CountConductEntries(ctx, s.db, q)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	pag, limit := req.GetPagination().GetResponse(total)
	resp := &pbjobs.ListConductEntriesResponse{
		Pagination: pag,
		Entries:    []*jobsconduct.ConductEntry{},
	}
	if total <= 0 {
		return resp, nil
	}

	q.Limit = limit
	resp.Entries, err = s.store.ListConductEntries(ctx, s.db, q)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetEntries() {
		if resp.GetEntries()[i].GetTargetUser() != nil {
			jobInfoFn(resp.GetEntries()[i].GetTargetUser())
		}
		if resp.GetEntries()[i].GetCreator() != nil {
			jobInfoFn(resp.GetEntries()[i].GetCreator())
		}

		if resp.GetEntries()[i].GetMessage() != nil &&
			resp.GetEntries()[i].GetMessage().GetContent() != nil {
			rawHtml, err := resp.GetEntries()[i].GetMessage().GetContent().ToHTML()
			if err != nil {
				return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
			}
			//nolint:staticcheck,nolintlint // legacy preview
			resp.GetEntries()[i].GetMessage().RawHtml = &rawHtml
		}
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)
	return resp, nil
}

func (s *Server) GetConductEntry(
	ctx context.Context,
	req *pbjobs.GetConductEntryRequest,
) (*pbjobs.GetConductEntryResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.jobs.conduct.id", req.GetId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)
	resp := &pbjobs.GetConductEntryResponse{Entry: &jobsconduct.ConductEntry{}}

	entry, err := s.store.GetConductEntry(ctx, s.db, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if entry == nil {
		return nil, errorsjobs.ErrFailedQuery
	}

	files, err := s.fHandler.ListFilesForParentID(ctx, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	entry.Files = files

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	if entry.GetTargetUser() != nil {
		jobInfoFn(entry.GetTargetUser())
	}
	if entry.GetCreator() != nil {
		jobInfoFn(entry.GetCreator())
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)
	resp.Entry = entry
	return resp, nil
}

func (s *Server) CreateConductEntry(
	ctx context.Context,
	req *pbjobs.CreateConductEntryRequest,
) (*pbjobs.CreateConductEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	req.Entry.Job = userInfo.GetJob()
	req.Entry.CreatorId = userInfo.GetUserId()

	lastId, err := s.store.CreateConductEntry(ctx, s.db, req.GetEntry())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	req.Entry.Id = lastId

	entry, err := s.store.GetConductEntry(ctx, s.db, lastId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	grpc_audit.AddMeta(ctx, "conduct.id", strconv.Itoa(int(lastId)))
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

	return &pbjobs.CreateConductEntryResponse{Entry: entry}, nil
}

func (s *Server) UpdateConductEntry(
	ctx context.Context,
	req *pbjobs.UpdateConductEntryRequest,
) (*pbjobs.UpdateConductEntryResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.jobs.conduct_id", req.GetEntry().GetId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	entry, err := s.store.GetConductEntry(ctx, s.db, req.GetEntry().GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if entry == nil || entry.GetJob() != userInfo.GetJob() {
		return nil, errorsjobs.ErrFailedQuery
	}

	if req.GetEntry().GetDraft() && !entry.GetDraft() && !userInfo.GetSuperuser() {
		req.Entry.Draft = entry.GetDraft()
	}
	if req.GetEntry().GetType() <= 0 {
		req.Entry.Type = entry.GetType()
	}
	if req.GetEntry().GetTargetUserId() == 0 {
		req.Entry.TargetUserId = entry.GetTargetUserId()
	}
	req.Entry.Job = userInfo.GetJob()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	defer tx.Rollback()

	if err := s.store.UpdateConductEntry(ctx, tx, req.GetEntry()); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	if _, _, err := s.fHandler.HandleFileChangesForParent(
		ctx,
		tx,
		req.GetEntry().GetId(),
		req.GetEntry().GetFiles(),
	); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	entry, err = s.store.GetConductEntry(ctx, s.db, entry.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	s.notifi.SendObjectEvent(ctx, &notificationsclientview.ObjectEvent{
		Type:      notificationsclientview.ObjectType_OBJECT_TYPE_JOBS_CONDUCT,
		Id:        &entry.Id,
		EventType: notificationsclientview.ObjectEventType_OBJECT_EVENT_TYPE_UPDATED,
		UserId:    &userInfo.UserId,
		Job:       &userInfo.Job,
	})

	return &pbjobs.UpdateConductEntryResponse{Entry: entry}, nil
}

func (s *Server) DeleteConductEntry(
	ctx context.Context,
	req *pbjobs.DeleteConductEntryRequest,
) (*pbjobs.DeleteConductEntryResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.jobs.conduct_id", req.GetId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)
	entry, err := s.store.GetConductEntry(ctx, s.db, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	var deletedAtTime *timestamp.Timestamp
	if entry == nil || entry.GetDeletedAt() == nil || !userInfo.GetSuperuser() {
		deletedAtTime = timestamp.Now()
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	} else {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_RESTORED)
	}

	if err := s.store.DeleteConductEntry(
		ctx,
		s.db,
		userInfo.GetJob(),
		req.GetId(),
		deletedAtTime,
	); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	return &pbjobs.DeleteConductEntryResponse{}, nil
}
