package croner

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
)

type NoopCron struct {
	ICron
}

func NewNoopCron() ICron {
	return &NoopCron{}
}

func (c *NoopCron) RegisterCronjob(ctx context.Context, job *cron.Cronjob) error {
	return nil
}

func (c *NoopCron) UnregisterCronjob(ctx context.Context, name string) error {
	return nil
}
