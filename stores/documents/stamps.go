package documentsstore

import (
	"context"
	"errors"

	resourcesdatabase "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	documentsstamps "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/stamps"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

const stampDefaultPageSize = 20

func (s *Store) CheckJobStampCount(ctx context.Context, job string) (int64, error) {
	tStamp := table.FivenetDocumentsStamps.AS("stamp")
	countStmt := tStamp.
		SELECT(mysql.COUNT(tStamp.ID).AS("data_count.total")).
		FROM(tStamp).
		WHERE(tStamp.Name.EQ(mysql.String(job)))

	var count struct {
		Total int64 `alias:"total"`
	}
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return count.Total, nil
}

func (s *Store) ListUsableStamps(
	ctx context.Context,
	q ListUsableStampsQuery,
) (*resourcesdatabase.PaginationResponse, []*documentsstamps.Stamp, error) {
	tStamp := table.FivenetDocumentsStamps.AS("stamp")

	condition := mysql.Bool(true)
	if !q.UserInfo.GetSuperuser() {
		condition = condition.AND(tStamp.DeletedAt.IS_NULL())
		condition = condition.AND(mysql.OR(
			tStamp.Name.EQ(mysql.String(q.UserInfo.GetJob())),
			s.subjectAccess.ACLAccessExistsCondition(
				tStamp.ID,
				q.UserInfo,
				int32(documentsstamps.StampAccessLevel_STAMP_ACCESS_LEVEL_USE),
			),
		))
	}

	countStmt := tStamp.
		SELECT(mysql.COUNT(tStamp.ID).AS("data_count.total")).
		FROM(tStamp).
		WHERE(condition)

	var count resourcesdatabase.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, nil, err
		}
	}
	pag, limit := q.Pagination.GetResponseWithPageSize(count.Total, stampDefaultPageSize)
	if count.Total <= 0 {
		return pag, []*documentsstamps.Stamp{}, nil
	}

	stmt := mysql.
		SELECT(
			tStamp.ID,
			tStamp.Name,
			tStamp.SvgTemplate,
			tStamp.VariantsJSON,
			tStamp.CreatedAt,
		).
		FROM(tStamp).
		WHERE(condition).
		OFFSET(q.Pagination.GetOffset()).
		ORDER_BY(tStamp.SortKey.ASC(), tStamp.CreatedAt.DESC()).
		LIMIT(limit)

	var stamps []*documentsstamps.Stamp
	if err := stmt.QueryContext(ctx, s.db, &stamps); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, nil, err
		}
	}

	return pag, stamps, nil
}

func (s *Store) GetStamp(ctx context.Context, stampID int64) (*documentsstamps.Stamp, error) {
	tStamp := table.FivenetDocumentsStamps.AS("stamp")
	stmt := mysql.
		SELECT(
			tStamp.ID,
			tStamp.CreatedAt,
			tStamp.Name,
			tStamp.SvgTemplate,
			tStamp.VariantsJSON,
		).
		FROM(tStamp).
		WHERE(tStamp.ID.EQ(mysql.Int64(stampID))).
		LIMIT(1)

	var stamp documentsstamps.Stamp
	if err := stmt.QueryContext(ctx, s.db, &stamp); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if stamp.Id == 0 {
		return nil, nil
	}

	return &stamp, nil
}

func (s *Store) CreateStamp(
	ctx context.Context,
	tx qrm.DB,
	stamp *documentsstamps.Stamp,
) (int64, error) {
	tStamp := table.FivenetDocumentsStamps
	res, err := tStamp.
		INSERT(
			tStamp.Name,
			tStamp.SvgTemplate,
			tStamp.VariantsJSON,
		).
		VALUES(
			stamp.GetName(),
			stamp.GetSvgTemplate(),
			mysql.StringExp(mysql.NULL),
		).
		ExecContext(ctx, tx)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (s *Store) UpdateStamp(ctx context.Context, tx qrm.DB, stamp *documentsstamps.Stamp) error {
	tStamp := table.FivenetDocumentsStamps
	stmt := tStamp.
		UPDATE(
			tStamp.Name,
			tStamp.SvgTemplate,
			tStamp.VariantsJSON,
		).
		SET(
			stamp.GetName(),
			stamp.GetSvgTemplate(),
			mysql.StringExp(mysql.NULL),
		).
		WHERE(tStamp.ID.EQ(mysql.Int64(stamp.GetId()))).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) DeleteStamp(ctx context.Context, tx qrm.DB, stampID int64) error {
	tStamp := table.FivenetDocumentsStamps.AS("stamp")
	stmt := tStamp.DELETE().WHERE(tStamp.ID.EQ(mysql.Int64(stampID))).LIMIT(1)
	_, err := stmt.ExecContext(ctx, tx)
	return err
}
