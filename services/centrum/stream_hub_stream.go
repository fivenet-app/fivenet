package centrum

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbcentrum "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"golang.org/x/sync/errgroup"
)

var (
	errFeedResync      = errors.New("centrum stream feed gap")
	errFeedClosed      = errors.New("centrum stream feed closed")
	errAccessChanged   = errors.New("centrum stream access changed")
	errUserInfoChanged = errors.New("centrum stream user info changed")
	errUserInfoResync  = errors.New("centrum stream user info resync")
)

func (s *Server) stream(
	ctx context.Context,
	srv pbcentrum.CentrumService_StreamServer,
	userInfo *userinfo.UserInfo,
	additionalJobs []string,
	feed <-chan *feedEvent,
	userInfoChanges <-chan *userinfo.UserInfoChanged,
	snapshotSequence uint64,
) error {
	jobs := append([]string{userInfo.GetJob()}, additionalJobs...)
	out := make(chan *pbcentrum.StreamResponse, 256)
	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		lastSequence := snapshotSequence
		for {
			select {
			case <-gctx.Done():
				return nil
			case change, ok := <-userInfoChanges:
				if !ok {
					return errUserInfoResync
				}
				if change == nil || change.GetUserId() != userInfo.GetUserId() {
					continue
				}
				userInfo.Job = change.GetNewJob()
				userInfo.JobGrade = change.GetNewJobGrade()
				return errUserInfoChanged
			case event, ok := <-feed:
				if !ok {
					return errFeedClosed
				}
				if event == nil || event.Response == nil || event.Sequence <= snapshotSequence {
					continue
				}
				if event.Sequence != lastSequence+1 {
					return fmt.Errorf(
						"%w: expected %d, received %d",
						errFeedResync,
						lastSequence+1,
						event.Sequence,
					)
				}
				lastSequence = event.Sequence
				if settings := event.Response.GetSettings(); settings != nil &&
					settings.GetJob() == userInfo.GetJob() {
					return errAccessChanged
				}
				if event.Response.GetSettingsDeleted() == userInfo.GetJob() {
					return errAccessChanged
				}
				if !slices.Contains(jobs, event.Job) {
					continue
				}
				select {
				case out <- event.Response:
				case <-gctx.Done():
					return nil
				}
			}
		}
	})

	g.Go(func() error {
		for {
			select {
			case <-gctx.Done():
				return nil
			case response := <-out:
				if response == nil {
					continue
				}
				if err := srv.Send(response); err != nil {
					if protoutils.IsContextCanceled(err) {
						return nil
					}
					return err
				}
			}
		}
	})

	return g.Wait()
}
