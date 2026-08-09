package demo

import (
	"context"
	"errors"
	"fmt"
	"time"

	resourcesettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/settings"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"go.uber.org/zap"
)

type demoBannerGenerator struct{}

func (g demoBannerGenerator) Name() string {
	return "demo_banner"
}

func (g demoBannerGenerator) Enabled(_ *Demo) bool {
	return true
}

func (g demoBannerGenerator) Run(ctx context.Context, d *Demo) error {
	return d.seedDemoBanner(ctx)
}

func (d *Demo) seedDemoBanner(ctx context.Context) error {
	if d.appCfg == nil {
		return errors.New("failed to seed demo banner: app config is not available")
	}
	if d.settingsStore == nil {
		return errors.New("failed to seed demo banner: settings store is not available")
	}

	cfg := d.appCfg.Get()
	if cfg == nil {
		return errors.New("failed to seed demo banner: app config is not loaded")
	}

	cfg.Default()
	cfg.GetSystem().BannerMessageEnabled = true
	cfg.GetSystem().BannerMessage = d.buildDemoBannerMessage()

	if err := d.settingsStore.UpdateAppConfig(ctx, cfg); err != nil {
		return fmt.Errorf("failed to persist demo banner config. %w", err)
	}

	if err := d.appCfg.Update(ctx, cfg); err != nil {
		return fmt.Errorf("failed to publish demo banner config update. %w", err)
	}

	d.logger.Info(
		"completed demo banner seeding",
		zap.String("banner_id", cfg.GetSystem().GetBannerMessage().GetId()),
	)

	return nil
}

func (d *Demo) buildDemoBannerMessage() *resourcesettings.BannerMessage {
	title := fmt.Sprintf(
		"Demo credentials: username %s, password %s.",
		demoAccountUsername,
		demoAccountPassword,
	)

	return &resourcesettings.BannerMessage{
		Id:        utils.GetMD5HashFromString(title + "-" + time.Time{}.String()),
		Title:     title,
		CreatedAt: timestamp.Now(),
	}
}
