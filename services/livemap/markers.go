package livemap

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	livemapmarkers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/livemap/markers"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
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
	reqMarker := req.GetMarker()
	if reqMarker != nil && reqMarker.GetId() > 0 {
		logging.InjectFields(
			ctx,
			logging.Fields{"fivenet.livemap.marker_id", reqMarker.GetId()},
		)
	}

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if reqMarker.Postal == nil || reqMarker.GetPostal() == "" {
		if postal, ok := s.postals.Closest(
			reqMarker.GetX(),
			reqMarker.GetY(),
		); postal != nil &&
			ok {
			reqMarker.Postal = postal.Code
		}
	}

	if reqMarker.GetId() <= 0 {
		if err := s.resolveMarkerPublic(nil, reqMarker, userInfo); err != nil {
			return nil, err
		}

		id, err := s.store.CreateMarker(
			ctx,
			reqMarker,
			func() *int32 {
				creatorID := userInfo.GetUserId()
				return &creatorID
			}(),
			userInfo.GetJob(),
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
		}

		reqMarker.SetId(id)
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)
	} else {
		fields, err := permslivemap.LivemapService.CreateOrUpdateMarker.AccessTyped.Get(
			s.ps,
			userInfo,
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
		}

		checkMarker, err := s.store.GetMarker(ctx, reqMarker.GetId())
		if err != nil {
			return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
		}

		if err := s.requirePublicMarkerMutationAccess(checkMarker, userInfo); err != nil {
			return nil, err
		}

		if err := s.resolveMarkerPublic(checkMarker, reqMarker, userInfo); err != nil {
			return nil, err
		}

		if !access.CheckIfHasOwnJobAccess(
			fields.StringList(),
			userInfo,
			checkMarker.GetCreator().GetJob(),
			checkMarker.GetCreator(),
		) {
			return nil, errorslivemap.ErrMarkerDenied
		}

		if err := s.store.UpdateMarker(ctx, reqMarker, checkMarker.GetJob()); err != nil {
			return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
		}

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	}

	reqMarker, err := s.store.GetMarker(ctx, reqMarker.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
	}
	s.enricher.EnrichJobName(reqMarker)
	s.applyMarkerCache(reqMarker)

	return &pblivemap.CreateOrUpdateMarkerResponse{
		Marker: reqMarker,
	}, nil
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

	if err := s.requirePublicMarkerMutationAccess(marker, userInfo); err != nil {
		return nil, err
	}

	if !access.CheckIfHasOwnJobAccess(
		fields.StringList(),
		userInfo,
		marker.GetCreator().GetJob(),
		marker.GetCreator(),
	) {
		return nil, errorslivemap.ErrMarkerDenied
	}

	var deletedAtTime *timestamp.Timestamp
	if marker == nil || marker.GetDeletedAt() == nil || !userInfo.GetJobAdmin() {
		deletedAtTime = timestamp.Now()
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	} else {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_RESTORED)
	}

	if err := s.store.DeleteMarker(ctx, req.GetId(), deletedAtTime); err != nil {
		return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
	}
	marker.SetDeletedAt(deletedAtTime)
	s.applyMarkerCache(marker)

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

	publicMarkers, publicDeletedMarkers := s.markersPublicCache.Snapshot()
	for _, marker := range publicMarkers {
		if marker.GetExpiresAt() == nil || time.Since(marker.GetExpiresAt().AsTime()) < 0 {
			updated = append(updated, marker)
		} else {
			deleted = append(deleted, marker.GetId())
		}
	}

	deleted = append(deleted, publicDeletedMarkers...)

	return updated, deleted
}

func (s *Server) refreshMarkers(ctx context.Context) error {
	dest, err := s.store.ListActiveMarkers(ctx)
	if err != nil {
		return err
	}

	markers := map[string][]*livemapmarkers.MarkerMarker{}
	publicMarkers := []*livemapmarkers.MarkerMarker{}
	for _, m := range dest {
		s.enricher.EnrichJobName(m)

		if m.GetPublic() {
			publicMarkers = append(publicMarkers, m)
			continue
		}

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

	publicDeletedMarkers, err := s.refreshDeletedMarkers(ctx)
	if err != nil {
		return err
	}

	s.markersPublicCache.Replace(publicMarkers, publicDeletedMarkers)

	return nil
}

func (s *Server) refreshDeletedMarkers(ctx context.Context) ([]int64, error) {
	deletedMarkers := map[string][]int64{}
	publicDeletedMarkers := []int64{}

	dest, err := s.store.ListDeletedMarkers(ctx)
	if err != nil {
		return nil, err
	}

	for _, m := range dest {
		if m.GetPublic() {
			publicDeletedMarkers = append(publicDeletedMarkers, m.GetId())
			continue
		}

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

	return publicDeletedMarkers, nil
}

func (s *Server) resolveMarkerPublic(
	existing *livemapmarkers.MarkerMarker,
	marker *livemapmarkers.MarkerMarker,
	userInfo *pbuserinfo.UserInfo,
) error {
	if existing == nil {
		if marker.GetPublic() && !userInfo.GetJobAdmin() {
			return errorslivemap.ErrMarkerDenied
		}

		return nil
	}

	if !marker.HasPublic() {
		marker.SetPublic(existing.GetPublic())
		return nil
	}

	if !userInfo.GetJobAdmin() && marker.GetPublic() != existing.GetPublic() {
		return errorslivemap.ErrMarkerDenied
	}

	return nil
}

func (s *Server) requirePublicMarkerMutationAccess(
	marker *livemapmarkers.MarkerMarker,
	userInfo *pbuserinfo.UserInfo,
) error {
	if marker.GetPublic() && !userInfo.GetJobAdmin() && marker.GetJob() != userInfo.GetJob() {
		return errorslivemap.ErrMarkerDenied
	}

	return nil
}

func (s *Server) applyMarkerCache(marker *livemapmarkers.MarkerMarker) {
	s.markersPublicCache.Apply(marker)

	if marker.GetPublic() {
		s.removeMarkerFromJobCache(marker.GetJob(), marker.GetId())
		s.removeMarkerFromDeletedCache(marker.GetJob(), marker.GetId())
		return
	}

	if marker.GetDeletedAt() == nil {
		s.upsertMarkerInJobCache(marker.GetJob(), marker)
		s.removeMarkerFromDeletedCache(marker.GetJob(), marker.GetId())
		return
	}

	s.removeMarkerFromJobCache(marker.GetJob(), marker.GetId())
	s.upsertMarkerInDeletedCache(marker.GetJob(), marker.GetId())
}

func (s *Server) upsertMarkerInJobCache(job string, marker *livemapmarkers.MarkerMarker) {
	markers, _ := s.markersCache.Load(job)
	markers = slices.Clone(markers)
	markers = slices.DeleteFunc(markers, func(existing *livemapmarkers.MarkerMarker) bool {
		return existing.GetId() == marker.GetId()
	})
	markers = append(markers, marker)
	s.markersCache.Store(job, markers)
}

func (s *Server) removeMarkerFromJobCache(job string, markerID int64) {
	markers, ok := s.markersCache.Load(job)
	if !ok {
		return
	}

	markers = slices.Clone(markers)
	markers = slices.DeleteFunc(markers, func(existing *livemapmarkers.MarkerMarker) bool {
		return existing.GetId() == markerID
	})
	if len(markers) == 0 {
		s.markersCache.Delete(job)
		return
	}

	s.markersCache.Store(job, markers)
}

func (s *Server) upsertMarkerInDeletedCache(job string, markerID int64) {
	markers, _ := s.markersDeletedCache.Load(job)
	markers = slices.Clone(markers)
	markers = slices.DeleteFunc(markers, func(existing int64) bool {
		return existing == markerID
	})
	markers = append(markers, markerID)
	s.markersDeletedCache.Store(job, markers)
}

func (s *Server) removeMarkerFromDeletedCache(job string, markerID int64) {
	markers, ok := s.markersDeletedCache.Load(job)
	if !ok {
		return
	}

	markers = slices.Clone(markers)
	markers = slices.DeleteFunc(markers, func(existing int64) bool {
		return existing == markerID
	})
	if len(markers) == 0 {
		s.markersDeletedCache.Delete(job)
		return
	}

	s.markersDeletedCache.Store(job, markers)
}
