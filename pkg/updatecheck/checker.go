//nolint:tagliatelle // GitHub API JSON response uses snake_case, so we have to use it for the tags.
package updatecheck

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/version"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("updatecheck",
	fx.Provide(New),
)

const minInterval = 30 * time.Minute // Minimum interval for update checks

var ErrNoVersionInfo = errors.New("no version info available")

// Checker periodically checks for updates to the application by querying the GitHub Releases API.
type Checker struct {
	mu sync.Mutex
	wg sync.WaitGroup

	logger *zap.Logger

	interval   time.Duration
	httpClient *http.Client

	currentTag string
	releaseUrl string
	releasedAt time.Time
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Config *config.Config
	Logger *zap.Logger
}

func New(p Params) *Checker {
	if !p.Config.UpdateCheck.Enabled {
		p.Logger.Info("update checker is disabled")
		return nil
	}

	logger := p.Logger.Named("updatechecker")

	interval := p.Config.UpdateCheck.Interval
	if p.Config.UpdateCheck.Interval < minInterval {
		logger.Warn(
			fmt.Sprintf("update check interval is too short, using minimum %s", minInterval),
		)
		interval = minInterval
	}

	c := &Checker{
		logger: logger,

		interval:   interval,
		httpClient: &http.Client{Timeout: 10 * time.Second},

		currentTag: version.Version, // Initialize with the current version to avoid logging on first run
	}

	cancelCtx, cancel := context.WithCancel(context.Background())
	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		c.wg.Go(func() {
			c.Start(cancelCtx)
		})
		return nil
	}))
	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()
		c.wg.Wait()

		return nil
	}))

	return c
}

// Start launches the update loop and blocks until ctx is done.
func (c *Checker) Start(ctx context.Context) {
	if version.Version == version.UnknownVersion {
		c.logger.Warn("no version info found, unable to check for updates")
	}

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	c.logger.Debug(
		fmt.Sprintf(
			"checking github for updates every %s (current=%s)",
			c.interval,
			version.Version,
		),
	)

	for {
		// Add a random delay (between 3 and 30 seconds) before every check to avoid thundering herd problem and "spontaneous synchronization."
		//nolint:gosec // G404 - The random delay is not security sensitive, it's just to avoid all instances collecting metrics at the same time.
		delay := time.Duration(3+rand.IntN(28)) * time.Second
		c.logger.Debug("update check delay", zap.Duration("delay", delay))

		select {
		case <-ctx.Done():
			return

		case <-time.After(delay):
		}

		newTag, htmlURL, isPrerelease, releasedAt, err := c.retrieveLatestTag(
			ctx,
			version.Owner,
			version.Repo,
			version.Version,
		)
		if err != nil {
			c.logger.Debug(fmt.Sprintf("updatechecker fetch error: %v", err))
			// Ignore pre-releases
		} else if !isPrerelease && newTag != c.currentTag {
			c.mu.Lock()
			c.currentTag = newTag
			c.releaseUrl = htmlURL
			c.releasedAt = releasedAt
			c.mu.Unlock()

			// Send log message about the new update
			c.logger.Info("🔔🆕 new update available!",
				zap.String("current", version.Version),
				zap.String("new", newTag),
				zap.String("url", htmlURL),
				zap.Time("released_at", releasedAt),
			)
		}

		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
		}
	}
}

// retrieveLatestTag hits the GitHub Releases API and returns the tag + url.
// It treats API/transport issues as errors but never panics.
func (c *Checker) retrieveLatestTag(
	ctx context.Context,
	owner string,
	repo string,
	currentVersion string,
) (string, string, bool, time.Time, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", "", false, time.Time{}, fmt.Errorf("failed to create request. %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", "", false, time.Time{}, fmt.Errorf(
			"failed to request release info. %w",
			err,
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("GitHub API returned %s", resp.Status)
		return "", "", false, time.Time{}, fmt.Errorf(
			"failed to retrieve latest release info (status code !== 200). %w",
			err,
		)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", false, time.Time{}, fmt.Errorf("failed to read response body. %w", err)
	}
	var payload struct {
		TagName    string `json:"tag_name"`
		HTMLURL    string `json:"html_url"`
		Prerelease bool   `json:"prerelease"`
		Draft      bool   `json:"draft"`
		CreatedAt  string `json:"created_at"`
	}
	if err = json.Unmarshal(body, &payload); err != nil {
		return "", "", false, time.Time{}, fmt.Errorf("failed to unmarshal response body. %w", err)
	}

	if payload.Draft {
		// Treat drafts as non-existent release.
		return currentVersion, "", false, time.Time{}, nil
	}

	releasedAt := time.Time{}
	if payload.CreatedAt != "" {
		if t, perr := time.Parse(time.RFC3339, payload.CreatedAt); perr == nil {
			releasedAt = t
		}
	}

	return payload.TagName, payload.HTMLURL, payload.Prerelease, releasedAt, nil
}

func (c *Checker) GetNewVersionInfo() (string, string, string, time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return version.Version, c.currentTag, c.releaseUrl, c.releasedAt
}
