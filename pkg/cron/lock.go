package cron

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/pkg/utils/protoutils"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	ownerLockTimeout = 15 * time.Second
)

func (cr *Cron) lockLoop() {
	wg := sync.WaitGroup{}

	for {
		select {
		case <-cr.ctx.Done():
			wg.Wait()

			if err := cr.ownerLock.Unlock(cr.ctx, "owner"); err != nil {
				cr.logger.Error("failed to unlock owner lock on shutdown", zap.Error(err))
			}

			return

		case <-time.After(2 * time.Second):
		}

		if err := cr.ownerLock.Lock(cr.ctx, "owner"); err != nil {
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				continue
			}

			cr.logger.Error("error while trying to get owner lock", zap.Error(err))
			continue
		}

		entry, err := cr.ownerKv.Get(cr.ctx, "owner")
		if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
			cr.logger.Error("error getting owner lock entry", zap.Error(err))
			continue
		}

		if entry != nil {
			// Make sure the owner state is really expired and not just a "hiccup" of our nats lock logic
			state := &cron.CronjobLockOwnerState{}
			if err := protojson.Unmarshal(entry.Value(), state); err != nil {
				cr.logger.Error("error getting owner lock entry", zap.Error(err))
				continue
			}

			// Hostname and updated at time is expired
			if time.Since(state.UpdatedAt.AsTime()) < ownerLockTimeout {
				continue
			}
		}

		cr.logger.Info("starting cron scheduler")

		func() {
			ctx, cancel := context.WithCancel(cr.ctx)
			defer cancel()

			// Keep lock owner state uptodate
			wg.Add(1)
			go func() {
				defer wg.Done()

				cr.scheduler.start(ctx)
			}()

			for {
				select {
				case <-ctx.Done():
					wg.Wait()
					return

				case <-time.After(3 * time.Second):
				}

				out, err := protoutils.Marshal(&cron.CronjobLockOwnerState{
					Hostname:  cr.name,
					UpdatedAt: timestamp.Now(),
				})
				if err != nil {
					cr.logger.Error("error marshalling owner lock state", zap.Error(err))
					cancel()
					continue
				}

				if _, err := cr.ownerKv.Put(ctx, "owner", out); err != nil {
					cr.logger.Error("failed to update owner lock state in kv", zap.Error(err))
					cancel()
					continue
				}
			}
		}()
	}
}
