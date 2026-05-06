package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/creativeprojects/go-selfupdate"
	"github.com/fivenet-app/fivenet/v2026/pkg/version"
)

type UpdateCmd struct {
	CheckOnly bool `help:"Check for updates only and don't update the binary."`
}

func (c *UpdateCmd) Run(ctx *kong.Context) error {
	log.Printf("ℹ️ Current FiveNet version: %s", version.Version)

	latest, found, err := c.checkLatestVersion()
	if err != nil {
		return err
	}
	if !found {
		return nil
	}

	if !c.CheckOnly {
		log.Printf(
			"🚀 Update FiveNet to version %s (size: %dMB)? (y/n): ",
			latest.Version(),
			latest.AssetByteSize/(1024*1024),
		)
		if !promptUserForConfirmation() {
			log.Println("❌ FiveNet Update cancelled.")
			return nil
		}

		log.Printf(
			"🔄 Updating FiveNet binary to version %s. This can take a few minutes based on your internet speed.",
			latest.Version(),
		)

		exe, err := selfupdate.ExecutablePath()
		if err != nil {
			return errors.New("could not locate fivenet executable path")
		}

		log.Printf("📦 Downloading new version from URL: %s", latest.AssetURL)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		if err := selfupdate.UpdateTo(
			ctx,
			latest.AssetURL,
			latest.AssetName,
			exe,
		); err != nil {
			return fmt.Errorf("error occurred while updating binary. %w", err)
		}

		log.Printf("✅ Successfully updated FiveNet to version %s", latest.Version())
	} else {
		log.Println(
			"🔎 Not updating because check only mode is enabled! To update the binary disable `--check-only`",
		)
	}

	return nil
}

func (c *UpdateCmd) checkLatestVersion() (*selfupdate.Release, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	updater, err := selfupdate.NewUpdater(selfupdate.Config{
		// Enable validator for SHA256 checksums
		Validator: &selfupdate.SHAValidator{},
	})
	if err != nil {
		return nil, false, fmt.Errorf("failed to create fivenet updater. %w", err)
	}

	latest, found, err := updater.DetectLatest(
		ctx,
		selfupdate.NewRepositorySlug(version.Owner, version.Repo),
	)
	if err != nil {
		return latest, found, fmt.Errorf(
			"error occurred while checking fivenet latest version. %w",
			err,
		)
	}
	if !found {
		return latest, found, fmt.Errorf(
			"latest version of fivenet for %s/%s could not be found from github repository",
			runtime.GOOS,
			runtime.GOARCH,
		)
	}

	currentVersion := version.Version
	if currentVersion == version.UnknownVersion {
		// Fallback to "sane" value to get the latest release anyways..
		currentVersion = "v0.0.0"
	}

	if latest.LessOrEqual(currentVersion) {
		log.Printf("✅ Current version (%s) is the latest", version.Version)
		return latest, found, nil
	}

	log.Printf("🆕 New FiveNet version %s found!", latest.Version())

	return latest, found, nil
}

func promptUserForConfirmation() bool {
	var input string

	for {
		fmt.Scanln(&input)

		switch strings.ToLower(strings.TrimSpace(input)) {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			fmt.Println("Invalid input, please enter y/n or yes/no.")
		}
	}
}
