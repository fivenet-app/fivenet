package documentsstore

import (
	"context"
	"errors"

	documentscategory "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/category"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) ListCategories(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
) ([]*documentscategory.Category, error) {
	tCategory := table.FivenetDocumentsCategories.AS("category")
	condition := tCategory.Job.EQ(mysql.String(userInfo.GetJob()))
	if !userInfo.GetSuperuser() {
		condition = mysql.AND(
			tCategory.DeletedAt.IS_NULL(),
			tCategory.Job.EQ(mysql.String(userInfo.GetJob())),
		)
	}

	stmt := tCategory.
		SELECT(
			tCategory.ID,
			tCategory.DeletedAt,
			tCategory.Name,
			tCategory.Description,
			tCategory.Job,
			tCategory.Color,
			tCategory.Icon,
		).
		FROM(tCategory).
		WHERE(condition).
		ORDER_BY(tCategory.SortKey.ASC())

	var dest []*documentscategory.Category
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}

func (s *Store) GetCategory(
	ctx context.Context,
	id int64,
	userInfo *userinfo.UserInfo,
) (*documentscategory.Category, error) {
	tCategory := table.FivenetDocumentsCategories.AS("category")
	stmt := tCategory.
		SELECT(
			tCategory.ID,
			tCategory.DeletedAt,
			tCategory.Name,
			tCategory.Description,
			tCategory.Job,
			tCategory.Color,
			tCategory.Icon,
		).
		FROM(tCategory).
		WHERE(mysql.AND(
			tCategory.Job.EQ(mysql.String(userInfo.GetJob())),
			tCategory.ID.EQ(mysql.Int64(id)),
		)).
		LIMIT(1)

	var dest documentscategory.Category
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.GetId() == 0 {
		return nil, nil
	}

	return &dest, nil
}

func (s *Store) CreateCategory(
	ctx context.Context,
	tx qrm.DB,
	category *documentscategory.Category,
	userInfo *userinfo.UserInfo,
) (int64, error) {
	tCategory := table.FivenetDocumentsCategories
	res, err := tCategory.
		INSERT(
			tCategory.Name,
			tCategory.Description,
			tCategory.Job,
			tCategory.Color,
			tCategory.Icon,
		).
		VALUES(
			category.GetName(),
			category.Description,
			userInfo.GetJob(),
			category.Color,
			category.Icon,
		).
		ExecContext(ctx, tx)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (s *Store) UpdateCategory(
	ctx context.Context,
	tx qrm.DB,
	category *documentscategory.Category,
	userInfo *userinfo.UserInfo,
) error {
	tCategory := table.FivenetDocumentsCategories
	stmt := tCategory.
		UPDATE(
			tCategory.Name,
			tCategory.Description,
			tCategory.Job,
			tCategory.Color,
			tCategory.Icon,
		).
		SET(
			category.GetName(),
			category.Description,
			userInfo.GetJob(),
			category.Color,
			category.Icon,
		).
		WHERE(mysql.AND(
			tCategory.ID.EQ(mysql.Int64(category.GetId())),
			tCategory.Job.EQ(mysql.String(userInfo.GetJob())),
		)).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) DeleteCategory(
	ctx context.Context,
	tx qrm.DB,
	id int64,
	userInfo *userinfo.UserInfo,
	deletedAt *timestamp.Timestamp,
) error {
	tCategory := table.FivenetDocumentsCategories
	stmt := tCategory.
		UPDATE(
			tCategory.DeletedAt,
		).
		SET(
			tCategory.DeletedAt.SET(dbutils.TimestampToMySQL(deletedAt)),
		).
		WHERE(mysql.AND(
			tCategory.Job.EQ(mysql.String(userInfo.GetJob())),
			tCategory.ID.EQ(mysql.Int64(id)),
		)).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}
