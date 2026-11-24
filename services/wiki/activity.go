package wiki

import (
	"context"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki"
	wikiaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki/access"
	wikiactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki/activity"
	pbwiki "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/wiki"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/textdiff"
	errorswiki "github.com/fivenet-app/fivenet/v2026/services/wiki/errors"
	citizenshydrator "github.com/fivenet-app/fivenet/v2026/stores/citizens/hydrator"
	wikistore "github.com/fivenet-app/fivenet/v2026/stores/wiki"
	logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) ListPageActivity(
	ctx context.Context,
	req *pbwiki.ListPageActivityRequest,
) (*pbwiki.ListPageActivityResponse, error) {
	logging.InjectFields(ctx, logging.Fields{pageIDLogFieldKey, req.GetPageId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetPageId(),
		userInfo,
		int32(wikiaccess.AccessLevel_ACCESS_LEVEL_VIEW),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	if !check {
		return nil, errorswiki.ErrPageDenied
	}

	count, err := s.store.CountPageActivity(
		ctx,
		wikistore.PageActivityQuery{PageID: req.GetPageId()},
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	pag, limit := req.GetPagination().
		GetResponseWithPageSize(count, database.WikiPageActivityPageSize)
	resp := &pbwiki.ListPageActivityResponse{
		Pagination: pag,
		Activity:   []*wikiactivity.PageActivity{},
	}
	if count <= 0 {
		return resp, nil
	}

	activity, err := s.store.ListPageActivity(ctx, wikistore.PageActivityQuery{
		PageID: req.GetPageId(),
		Offset: req.GetPagination().GetOffset(),
		Limit:  limit,
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}
	resp.Activity = activity

	hydrateShort := s.hydrator.HydrateBasicTargetsSafeFunc(userInfo)
	targets := make([]citizenshydrator.BasicTarget, 0, len(resp.GetActivity()))
	for i := range resp.GetActivity() {
		if resp.GetActivity()[i].GetCreatorId() <= 0 {
			continue
		}
		targets = append(targets, citizenshydrator.BasicTarget{
			UserID: resp.GetActivity()[i].GetCreatorId(),
			Set:    resp.Activity[i].SetCreator,
		})
	}
	if err := hydrateShort(ctx, nil, targets); err != nil {
		return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	return resp, nil
}

// generatePageDiff produces activity changes from page metadata and extracted content text.
func (s *Server) generatePageDiff(
	old *wiki.Page,
	new *wiki.Page,
) (*wikiactivity.PageUpdated, error) {
	diff := &wikiactivity.PageUpdated{}

	if old.GetMeta().GetTitle() != new.GetMeta().GetTitle() {
		if titleDiff := textdiff.DiffText(
			old.GetMeta().GetTitle(),
			new.GetMeta().GetTitle(),
		); titleDiff.HasChanges() {
			diff.TitleCdiff = titleDiff
		}
	}

	if old.GetMeta().GetDescription() != new.GetMeta().GetDescription() {
		if descriptionDiff := textdiff.DiffText(
			old.GetMeta().GetDescription(),
			new.GetMeta().GetDescription(),
		); descriptionDiff.HasChanges() {
			diff.DescriptionCdiff = descriptionDiff
		}
	}

	if cd := textdiff.DiffText(
		old.GetContent().Extract().GetText(),
		new.GetContent().Extract().GetText(),
	); cd.HasChanges() {
		diff.ContentCdiff = cd
	}

	return diff, nil
}

func pageUpdatedHasChanges(diff *wikiactivity.PageUpdated) bool {
	if diff == nil {
		return false
	}

	return diff.GetTitleCdiff() != nil ||
		diff.GetDescriptionCdiff() != nil ||
		diff.GetContentCdiff() != nil ||
		diff.GetTitleDiff() != "" ||
		diff.GetDescriptionDiff() != "" ||
		diff.GetContentDiff() != "" ||
		(diff.GetFilesChange() != nil &&
			(diff.GetFilesChange().GetAdded() > 0 || diff.GetFilesChange().GetDeleted() > 0))
}
