package livemap

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	livemapmarkers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/livemap/markers"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	pblivemap "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/livemap"
	permslivemap "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/livemap/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorslivemap "github.com/fivenet-app/fivenet/v2026/services/livemap/errors"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) CreateOrUpdateMarker(
	ctx context.Context,
	req *pblivemap.CreateOrUpdateMarkerRequest,
) (*pblivemap.CreateOrUpdateMarkerResponse, error) {
	if req.GetMarker() != nil && req.GetMarker().GetId() > 0 {
		logging.InjectFields(
			ctx,
			logging.Fields{"fivenet.livemap.marker_id", req.GetMarker().GetId()},
		)
	}

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if req.Marker.Postal == nil || *req.Marker.Postal == "" {
		if postal, ok := s.postals.Closest(req.Marker.X, req.Marker.Y); postal != nil && ok {
			req.Marker.Postal = postal.Code
		}
	}

	if req.GetMarker().GetId() <= 0 {
		id, err := s.store.CreateMarker(
			ctx,
			req.GetMarker(),
			userInfo.GetUserId(),
			userInfo.GetJob(),
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
		}

		req.Marker.Id = id
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)
	} else {
		fields, err := permslivemap.LivemapService.CreateOrUpdateMarker.AccessTyped.Get(
			s.ps,
			userInfo,
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
		}

		marker, err := s.store.GetMarker(ctx, req.GetMarker().GetId())
		if err != nil {
			return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
		}

		if !access.CheckIfHasOwnJobAccess(
			fields.StringList(),
			userInfo,
			marker.GetCreator().GetJob(),
			marker.GetCreator(),
		) {
			return nil, errorslivemap.ErrMarkerDenied
		}

		if err := s.store.UpdateMarker(ctx, req.GetMarker(), marker.GetJob()); err != nil {
			return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
		}

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	}

	marker, err := s.store.GetMarker(ctx, req.GetMarker().GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
	}
	s.enricher.EnrichJobName(marker)

	if err := s.sendUpdateEvent(
		ctx,
		MarkerTopic,
		MarkerUpdate,
		marker.GetJob(),
		marker,
	); err != nil {
		return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
	}

	return &pblivemap.CreateOrUpdateMarkerResponse{Marker: marker}, nil
}

func (s *Server) DeleteMarker(
	ctx context.Context,
	req *pblivemap.DeleteMarkerRequest,
) (*pblivemap.DeleteMarkerResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.livemap.marker_id", req.GetId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	fields, err := permslivemap.LivemapService.DeleteMarker.AccessTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
	}

	marker, err := s.store.GetMarker(ctx, req.GetId())
	if err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
		}

		return &pblivemap.DeleteMarkerResponse{}, nil
	}
	s.enricher.EnrichJobName(marker)

	if !access.CheckIfHasOwnJobAccess(
		fields.StringList(),
		userInfo,
		marker.GetCreator().GetJob(),
		marker.GetCreator(),
	) {
		return nil, errorslivemap.ErrMarkerDenied
	}

	if err := s.store.SoftDeleteMarker(ctx, req.GetId()); err != nil {
		return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
	}

	if err := s.sendUpdateEvent(
		ctx,
		MarkerTopic,
		MarkerDelete,
		marker.GetJob(),
		marker,
	); err != nil {
		return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
	}

	if markers, ok := s.markersDeletedCache.Load(marker.GetJob()); ok {
		s.markersDeletedCache.Store(marker.GetJob(), append(markers, marker.GetId()))
	}

	if markers, ok := s.markersCache.Load(marker.GetJob()); ok {
		s.markersCache.Store(
			marker.GetJob(),
			slices.DeleteFunc(markers, func(m *livemapmarkers.MarkerMarker) bool {
				return m.GetId() == marker.GetId()
			}),
		)
	}

	return &pblivemap.DeleteMarkerResponse{}, nil
}

func (s *Server) getMarkerMarkers(
	jobs *permissionsattributes.StringList,
) ([]*livemapmarkers.MarkerMarker, []int64) {
	updated := []*livemapmarkers.MarkerMarker{}
	deleted := []int64{}

	for _, job := range jobs.GetStrings() {
		markers, _ := s.markersCache.Load(job)

		for _, marker := range markers {
			if marker.GetExpiresAt() == nil || time.Since(marker.GetExpiresAt().AsTime()) < 0 {
				updated = append(updated, marker)
			} else {
				deleted = append(deleted, marker.GetId())
			}
		}

		deletedMarkers, _ := s.markersDeletedCache.Load(job)
		deleted = append(deleted, deletedMarkers...)
	}

	return updated, deleted
}

func (s *Server) refreshMarkers(ctx context.Context) error {
	dest, err := s.store.ListActiveMarkers(ctx)
	if err != nil {
		return err
	}

	markers := map[string][]*livemapmarkers.MarkerMarker{}
	for _, m := range dest {
		s.enricher.EnrichJobName(m)

		if _, ok := markers[m.GetJob()]; !ok {
			markers[m.GetJob()] = []*livemapmarkers.MarkerMarker{}
		}

		markers[m.GetJob()] = append(markers[m.GetJob()], m)
	}

	for job, ms := range markers {
		if len(ms) == 0 {
			s.markersCache.Delete(job)
		} else {
			s.markersCache.Store(job, ms)
		}
	}

	for key := range s.markersCache.All() {
		if _, ok := markers[key]; !ok {
			s.markersCache.Delete(key)
		}
	}

	return s.refreshDeletedMarkers(ctx)
}

func (s *Server) refreshDeletedMarkers(ctx context.Context) error {
	deletedMarkers := map[string][]int64{}

	dest, err := s.store.ListDeletedMarkers(ctx)
	if err != nil {
		return err
	}

	for _, m := range dest {
		if _, ok := deletedMarkers[m.GetJob()]; !ok {
			deletedMarkers[m.GetJob()] = []int64{}
		}

		deletedMarkers[m.GetJob()] = append(deletedMarkers[m.GetJob()], m.GetId())
	}

	for job, ms := range deletedMarkers {
		if len(ms) == 0 {
			s.markersDeletedCache.Delete(job)
		} else {
			s.markersDeletedCache.Store(job, ms)
		}
	}

	for key := range s.markersDeletedCache.All() {
		if _, ok := deletedMarkers[key]; !ok {
			s.markersDeletedCache.Delete(key)
		}
	}

	return nil
}
