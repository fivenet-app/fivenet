package croner

import (
	"context"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
)

type NoopCron struct {
	IRegistry
}

func NewNoopRegistry() IRegistry {
	return &NoopCron{}
}

func (c *NoopCron) RegisterCronjob(ctx context.Context, job *cron.Cronjob) error {
	return nil
}

func (c *NoopCron) UnregisterCronjob(ctx context.Context, name string) error {
	return nil
}
