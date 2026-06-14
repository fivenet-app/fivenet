package mailer

import (
	"context"
	"errors"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	mailertemplates "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/templates"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var tTemplates = table.FivenetMailerTemplates.AS("template")

func (s *Store) ListTemplates(
	ctx context.Context,
	q qrm.DB,
	emailID int64,
	limit int64,
) ([]*mailertemplates.Template, error) {
	stmt := tTemplates.
		SELECT(
			tTemplates.ID,
			tTemplates.CreatedAt,
			tTemplates.UpdatedAt,
			tTemplates.DeletedAt,
			tTemplates.EmailID,
			tTemplates.Title,
			tTemplates.Content,
		).
		FROM(tTemplates).
		WHERE(tTemplates.EmailID.EQ(mysql.Int64(emailID))).
		LIMIT(limit)

	templates := []*mailertemplates.Template{}
	if err := stmt.QueryContext(ctx, s.dbOr(q), &templates); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return templates, nil
}

func (s *Store) GetTemplate(
	ctx context.Context,
	q qrm.DB,
	id int64,
	emailID *int64,
) (*mailertemplates.Template, error) {
	condition := tTemplates.ID.EQ(mysql.Int64(id))
	if emailID == nil || *emailID <= 0 {
		condition = condition.AND(tTemplates.EmailID.IS_NULL())
	} else {
		condition = condition.AND(tTemplates.EmailID.EQ(mysql.Int64(*emailID)))
	}

	stmt := tTemplates.
		SELECT(
			tTemplates.ID,
			tTemplates.CreatedAt,
			tTemplates.UpdatedAt,
			tTemplates.DeletedAt,
			tTemplates.EmailID,
			tTemplates.Title,
			tTemplates.Content,
			tTemplates.CreatorJob,
			tTemplates.CreatorID,
		).
		FROM(tTemplates).
		WHERE(condition).
		LIMIT(1)

	dest := &mailertemplates.Template{}
	if err := stmt.QueryContext(ctx, s.dbOr(q), dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.GetId() == 0 {
		return nil, nil
	}

	return dest, nil
}

func (s *Store) CountTemplatesByCreatorJob(
	ctx context.Context,
	q qrm.DB,
	job string,
) (int64, error) {
	countStmt := tTemplates.
		SELECT(
			mysql.COUNT(tTemplates.ID).AS("data_count.total"),
		).
		FROM(tTemplates).
		WHERE(tTemplates.CreatorJob.EQ(mysql.String(job)))

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.dbOr(q), &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return count.Total, nil
}

func (s *Store) CreateTemplate(
	ctx context.Context,
	q qrm.DB,
	template *mailertemplates.Template,
	creatorID int32,
) (int64, error) {
	stmt := tTemplates.
		INSERT(
			tTemplates.EmailID,
			tTemplates.Title,
			tTemplates.Content,
			tTemplates.CreatorJob,
			tTemplates.CreatorID,
		).
		VALUES(
			template.GetEmailId(),
			template.GetTitle(),
			template.GetContent(),
			template.GetCreatorJob(),
			creatorID,
		)

	res, err := stmt.ExecContext(ctx, s.dbOr(q))
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

func (s *Store) UpdateTemplate(
	ctx context.Context,
	q qrm.DB,
	template *mailertemplates.Template,
) error {
	stmt := tTemplates.
		UPDATE(
			tTemplates.Title,
			tTemplates.Content,
		).
		SET(
			template.GetTitle(),
			template.GetContent(),
		).
		WHERE(tTemplates.ID.EQ(mysql.Int64(template.GetId()))).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.dbOr(q)); err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteTemplate(ctx context.Context, q qrm.DB, id int64) error {
	stmt := tTemplates.
		UPDATE().
		SET(
			tTemplates.DeletedAt.SET(mysql.CURRENT_TIMESTAMP()),
		).
		WHERE(tTemplates.ID.EQ(mysql.Int64(id))).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.dbOr(q)); err != nil {
		return err
	}

	return nil
}
