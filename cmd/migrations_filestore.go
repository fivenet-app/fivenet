package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/cmd/envs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/storage"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/fx"
	"go.uber.org/multierr"
)

type FilestoreCmd struct {
	db      *sql.DB
	storage storage.IStorage
}

func (c *FilestoreCmd) Run() error {
	fxOpts := getFxBaseOpts(Cli.StartTimeout, false)

	if err := os.Setenv(envs.SkipDBMigrationsEnv, "true"); err != nil {
		return err
	}

	fxOpts = append(fxOpts,
		fx.Invoke(func(lifecycle fx.Lifecycle, cfg *config.Config, db *sql.DB, storage storage.IStorage, shutdowner fx.Shutdowner) {
			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						exitCode := 0
						c.db = db
						c.storage = storage
						if err := c.run(ctx); err != nil {
							// handle error, set non-zero exit code so caller knows the job failed
							exitCode = 1
							fmt.Println("Error running filestore command:", err)
						}
						_ = shutdowner.Shutdown(fx.ExitCode(exitCode))
					}()
					return nil
				},
			})
		}),
	)

	app := fx.New(fxOpts...)
	app.Run()

	return nil
}

func (c *FilestoreCmd) run(ctx context.Context) error {
	var errs error
	if err := c.migrateJobLogos(ctx); err != nil {
		errs = multierr.Append(errs, fmt.Errorf("failed to migrate job logos: %w", err))
	}

	if err := c.migrateAvatars(ctx); err != nil {
		errs = multierr.Append(errs, fmt.Errorf("failed to migrate avatars: %w", err))
	}

	if err := c.migrateMugshots(ctx); err != nil {
		errs = multierr.Append(errs, fmt.Errorf("failed to migrate mugshots: %w", err))
	}

	return errs
}

func (c *FilestoreCmd) migrateJobLogos(ctx context.Context) error {
	// Query for jobs with non-null logo URLs
	tJobProps := table.FivenetJobProps
	tFiles := table.FivenetFiles

	stmt := tJobProps.
		SELECT(
			tJobProps.Job.AS("job"),
			tJobProps.LogoURL.AS("logo_url"),
		).
		WHERE(
			tJobProps.LogoURL.IS_NOT_NULL(),
		)

	var rows []struct {
		Job     string `alias:"job"`
		LogoURL string `alias:"logo_url"`
	}
	if err := stmt.QueryContext(ctx, c.db, &rows); err != nil {
		return err
	}

	var errs error
	for _, row := range rows {
		row.LogoURL = strings.TrimPrefix(row.LogoURL, "/api/filestore")

		objectInfo, err := c.storage.Stat(ctx, row.LogoURL)
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to stat logo file %q. %w", row.LogoURL, err))
			continue
		}

		fmt.Println(objectInfo)

		insertStmt := tFiles.
			INSERT(
				tFiles.FilePath,
				tFiles.ByteSize,
				tFiles.ContentType,
				tFiles.CreatedAt,
			).
			VALUES(
				strings.TrimPrefix(objectInfo.GetName(), "/"),
				objectInfo.GetSize(),
				objectInfo.GetContentType(),
				timestamp.New(objectInfo.GetLastModified()),
			)

		res, err := insertStmt.ExecContext(ctx, c.db)
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to insert file %q for job %q. %w", objectInfo.GetName(), row.Job, err))
			continue
		}
		lastInsertID, err := res.LastInsertId()
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to get last insert ID for file %q for job %q. %w", objectInfo.GetName(), row.Job, err))
			continue
		}

		updateStmt := tJobProps.
			UPDATE(
				tJobProps.LogoURL,
				tJobProps.LogoFileID,
			).
			SET(
				jet.NULL,
				lastInsertID,
			).
			WHERE(
				tJobProps.Job.EQ(jet.String(row.Job)),
			)

		if _, err := updateStmt.ExecContext(ctx, c.db); err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to update job %q with new logo file ID %d. %w", row.Job, lastInsertID, err))
			continue
		}

		fmt.Printf("Updated job %q with new logo file ID %d\n", row.Job, lastInsertID)
	}

	return nil
}

func (c FilestoreCmd) migrateAvatars(ctx context.Context) error {
	var errs error

	// Query for users with non-null avatar columns
	tUserProps := table.FivenetUserProps
	tFiles := table.FivenetFiles

	stmt := tUserProps.
		SELECT(
			tUserProps.UserID.AS("user_id"),
			tUserProps.Avatar.AS("avatar"),
		).
		WHERE(
			tUserProps.Avatar.IS_NOT_NULL(),
		)

	var rows []struct {
		UserId int32  `alias:"userId"`
		Avatar string `alias:"avatar"`
	}
	if err := stmt.QueryContext(ctx, c.db, &rows); err != nil {
		return err
	}

	for _, row := range rows {
		row.Avatar = strings.TrimPrefix(row.Avatar, "/api/filestore")

		objectInfo, err := c.storage.Stat(ctx, row.Avatar)
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to stat avatar file %q. %w", row.Avatar, err))
			continue
		}

		insertStmt := tFiles.
			INSERT(
				tFiles.FilePath,
				tFiles.ByteSize,
				tFiles.ContentType,
				tFiles.CreatedAt,
			).
			VALUES(
				strings.TrimPrefix(objectInfo.GetName(), "/"),
				objectInfo.GetSize(),
				objectInfo.GetContentType(),
				timestamp.New(objectInfo.GetLastModified()),
			)

		res, err := insertStmt.ExecContext(ctx, c.db)
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to insert avatar file %q for user %d. %w", objectInfo.GetName(), row.UserId, err))
			continue
		}
		lastInsertID, err := res.LastInsertId()
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to get last insert ID for avatar file %q for job %q. %w", objectInfo.GetName(), row.UserId, err))
			continue
		}

		updateStmt := tUserProps.
			UPDATE(
				tUserProps.Avatar,
				tUserProps.AvatarFileID,
			).
			SET(
				jet.NULL,
				lastInsertID,
			).
			WHERE(
				tUserProps.UserID.EQ(jet.Int32(row.UserId)),
			)

		if _, err := updateStmt.ExecContext(ctx, c.db); err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to update user %d with new avatar file ID %d. %w", row.UserId, lastInsertID, err))
			continue
		}

		fmt.Printf("Updated user %d with new avatar file ID %d\n", row.UserId, lastInsertID)
	}

	return errs
}

func (c FilestoreCmd) migrateMugshots(ctx context.Context) error {
	var errs error

	// Query for users with non-null avatar columns
	tUserProps := table.FivenetUserProps
	tFiles := table.FivenetFiles

	stmt := tUserProps.
		SELECT(
			tUserProps.UserID.AS("user_id"),
			tUserProps.MugShot.AS("mugshot"),
		).
		WHERE(
			tUserProps.MugShot.IS_NOT_NULL(),
		)

	var rows []struct {
		UserId  int32  `alias:"userId"`
		Mugshot string `alias:"mugshot"`
	}
	if err := stmt.QueryContext(ctx, c.db, &rows); err != nil {
		return err
	}

	for _, row := range rows {
		row.Mugshot = strings.TrimPrefix(row.Mugshot, "/api/filestore")

		objectInfo, err := c.storage.Stat(ctx, row.Mugshot)
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to stat mugshot file %q. %w", row.Mugshot, err))
			continue
		}

		insertStmt := tFiles.
			INSERT(
				tFiles.FilePath,
				tFiles.ByteSize,
				tFiles.ContentType,
				tFiles.CreatedAt,
			).
			VALUES(
				strings.TrimPrefix(objectInfo.GetName(), "/"),
				objectInfo.GetSize(),
				objectInfo.GetContentType(),
				timestamp.New(objectInfo.GetLastModified()),
			)

		res, err := insertStmt.ExecContext(ctx, c.db)
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to insert mugshot file %q for user %d. %w", objectInfo.GetName(), row.UserId, err))
			continue
		}
		lastInsertID, err := res.LastInsertId()
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to get last insert ID for mugshot file %q for job %q. %w", objectInfo.GetName(), row.UserId, err))
			continue
		}

		updateStmt := tUserProps.
			UPDATE(
				tUserProps.MugShot,
				tUserProps.MugshotFileID,
			).
			SET(
				jet.NULL,
				lastInsertID,
			).
			WHERE(
				tUserProps.UserID.EQ(jet.Int32(row.UserId)),
			)

		if _, err := updateStmt.ExecContext(ctx, c.db); err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to update user %d with new mugshot file ID %d. %w", row.UserId, lastInsertID, err))
			continue
		}

		fmt.Printf("Updated user %d with new mugshot file ID %d\n", row.UserId, lastInsertID)
	}

	return errs
}
