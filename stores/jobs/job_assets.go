package jobsstore

import (
	"context"
	"errors"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	file "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/file"
	jobresources "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type jobAssetRow struct {
	ID              int64     `alias:"job_asset.id"`
	CreatedByUserID *int32    `alias:"job_asset.created_by_user_id"`
	FilePath        string    `alias:"job_asset.file_path"`
	ByteSize        int64     `alias:"job_asset.byte_size"`
	ContentType     string    `alias:"job_asset.content_type"`
	CreatedAt       time.Time `alias:"job_asset.created_at"`
	DisplayName     string    `alias:"job_asset.display_name"`
}

func (s *Store) CountJobAssets(ctx context.Context, db qrm.DB, job string) (int64, error) {
	tAssets := table.FivenetJobAssets
	var count database.DataCount
	if err := tAssets.
		SELECT(
			mysql.COUNT(tAssets.FileID).AS("data_count.total"),
		).
		WHERE(tAssets.Job.EQ(mysql.String(job))).
		QueryContext(ctx, db, &count); err != nil {
		return 0, err
	}
	return count.Total, nil
}

func (s *Store) ListJobAssets(
	ctx context.Context,
	db qrm.DB,
	job string,
	offset, limit int64,
) ([]*jobresources.JobAsset, error) {
	tAssets := table.FivenetJobAssets
	tFiles := table.FivenetFiles
	stmt := tAssets.
		SELECT(
			tAssets.FileID.AS("job_asset.id"),
			tAssets.CreatedByUserID.AS("job_asset.created_by_user_id"),
			tAssets.DisplayName.AS("job_asset.display_name"),
			tFiles.FilePath.AS("job_asset.file_path"),
			tFiles.ByteSize.AS("job_asset.byte_size"),
			tFiles.ContentType.AS("job_asset.content_type"),
			tFiles.CreatedAt.AS("job_asset.created_at"),
		).
		FROM(
			tAssets.INNER_JOIN(
				tFiles,
				tAssets.FileID.EQ(tFiles.ID),
			),
		).
		WHERE(
			tAssets.Job.EQ(mysql.String(job)),
		).
		ORDER_BY(
			tAssets.CreatedAt.DESC(),
		).
		LIMIT(limit).
		OFFSET(offset)

	rows := []*jobAssetRow{}
	if err := stmt.QueryContext(ctx, db, &rows); err != nil {
		return nil, err
	}

	assetsResp := make([]*jobresources.JobAsset, 0, len(rows))
	for _, row := range rows {
		assetsResp = append(assetsResp, jobAssetResource(row))
	}
	return assetsResp, nil
}

func (s *Store) GetJobAsset(
	ctx context.Context,
	db qrm.DB,
	job string,
	id int64,
) (*jobresources.JobAsset, error) {
	tAssets := table.FivenetJobAssets
	tFiles := table.FivenetFiles

	row := &jobAssetRow{}

	err := tAssets.
		SELECT(
			tAssets.FileID.AS("job_asset.id"),
			tAssets.CreatedByUserID.AS("job_asset.created_by_user_id"),
			tAssets.DisplayName.AS("job_asset.display_name"),
			tFiles.FilePath.AS("job_asset.file_path"),
			tFiles.ByteSize.AS("job_asset.byte_size"),
			tFiles.ContentType.AS("job_asset.content_type"),
			tFiles.CreatedAt.AS("job_asset.created_at"),
		).
		FROM(
			tAssets.INNER_JOIN(
				tFiles,
				tAssets.FileID.EQ(tFiles.ID),
			),
		).
		WHERE(mysql.AND(
			tAssets.Job.EQ(mysql.String(job)),
			tAssets.FileID.EQ(mysql.Int64(id)),
		)).
		LIMIT(1).
		QueryContext(ctx, db, row)
	if errors.Is(err, qrm.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return jobAssetResource(row), nil
}

func (s *Store) UpdateJobAssetMetadata(
	ctx context.Context,
	db qrm.DB,
	job string,
	fileID int64,
	createdByUserID int32,
	displayName string,
) error {
	tAssets := table.FivenetJobAssets
	_, err := tAssets.
		UPDATE().
		SET(
			tAssets.CreatedByUserID.SET(mysql.Int32(createdByUserID)),
			tAssets.DisplayName.SET(mysql.String(displayName)),
		).
		WHERE(mysql.AND(
			tAssets.Job.EQ(mysql.String(job)),
			tAssets.FileID.EQ(mysql.Int64(fileID)),
		)).
		ExecContext(ctx, db)
	return err
}

func (s *Store) UpdateJobAsset(
	ctx context.Context,
	db qrm.DB,
	job string,
	fileID int64,
	displayName string,
) error {
	tAssets := table.FivenetJobAssets
	_, err := tAssets.
		UPDATE().
		SET(
			tAssets.DisplayName.SET(mysql.String(displayName)),
		).
		WHERE(mysql.AND(
			tAssets.Job.EQ(mysql.String(job)),
			tAssets.FileID.EQ(mysql.Int64(fileID)),
		)).
		LIMIT(1).
		ExecContext(ctx, db)
	return err
}

func jobAssetResource(row *jobAssetRow) *jobresources.JobAsset {
	return &jobresources.JobAsset{
		Id:              row.ID,
		CreatedByUserId: row.CreatedByUserID,
		DisplayName:     row.DisplayName,
		CreatedAt:       timestamp.New(row.CreatedAt),
		File: &file.File{
			Id:          row.ID,
			FilePath:    row.FilePath,
			ByteSize:    row.ByteSize,
			ContentType: row.ContentType,
			CreatedAt:   timestamp.New(row.CreatedAt),
		},
	}
}
