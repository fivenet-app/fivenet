package documentsstore

import (
	"context"
	"errors"

	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	documentstemplates "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/templates"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) ListTemplates(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
) ([]*documentstemplates.TemplateShort, error) {
	if userInfo == nil {
		userInfo = &userinfo.UserInfo{}
	}

	tDTemplates := table.FivenetDocumentsTemplates.AS("template_short")
	tDTemplatesFilter := table.FivenetDocumentsTemplates.AS("template_filter")
	tDCategory := table.FivenetDocumentsCategories.AS("category")
	visibleQuery := s.templateAccess.VisibleIDsByConditionQuery(
		userInfo,
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		tDTemplatesFilter.DeletedAt.IS_NULL(),
	)
	templateID := mysql.IntegerColumn("id").From(visibleQuery.Table)

	innerStmt := tDTemplates.
		SELECT(
			tDTemplates.ID,
			tDTemplates.Weight,
			tDCategory.ID,
			tDCategory.Name,
			tDCategory.Description,
			tDCategory.Job,
			tDCategory.Color,
			tDCategory.Icon,
			tDTemplates.Title,
			tDTemplates.Description,
			tDTemplates.Color,
			tDTemplates.Icon,
			tDTemplates.Schema,
			tDTemplates.Workflow,
			tDTemplates.Approval,
			tDTemplates.CreatorJob,
		).
		FROM(
			visibleQuery.Table.
				INNER_JOIN(tDTemplates,
					tDTemplates.ID.EQ(templateID),
				).
				LEFT_JOIN(tDCategory,
					tDCategory.ID.EQ(tDTemplates.CategoryID),
				),
		).
		ORDER_BY(
			tDTemplates.Weight.DESC(),
			tDTemplates.ID.ASC(),
		)

	var stmt mysql.Statement = innerStmt
	if len(visibleQuery.CTEs) > 0 {
		stmt = mysql.WITH(visibleQuery.CTEs...)(innerStmt)
	}

	var dest []*documentstemplates.TemplateShort
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}

func (s *Store) GetTemplate(
	ctx context.Context,
	templateID int64,
) (*documentstemplates.Template, error) {
	tDTemplates := table.FivenetDocumentsTemplates.AS("template")
	tDCategory := table.FivenetDocumentsCategories.AS("category")

	stmt := tDTemplates.
		SELECT(
			tDTemplates.ID,
			tDTemplates.Weight,
			tDTemplates.CreatedAt,
			tDTemplates.UpdatedAt,
			tDCategory.ID,
			tDCategory.Name,
			tDCategory.Description,
			tDCategory.Job,
			tDCategory.Color,
			tDCategory.Icon,
			tDTemplates.Title,
			tDTemplates.Description,
			tDTemplates.Color,
			tDTemplates.Icon,
			tDTemplates.ContentTitle,
			tDTemplates.Content,
			tDTemplates.State,
			tDTemplates.Access,
			tDTemplates.Schema,
			tDTemplates.Workflow,
			tDTemplates.Approval,
			tDTemplates.CreatorJob,
		).
		FROM(
			tDTemplates.
				LEFT_JOIN(tDCategory,
					tDCategory.ID.EQ(tDTemplates.CategoryID),
				),
		).
		WHERE(mysql.AND(
			tDTemplates.DeletedAt.IS_NULL(),
			tDTemplates.ID.EQ(mysql.Int64(templateID)),
		)).
		LIMIT(1)

	var dest documentstemplates.Template
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.GetId() <= 0 {
		return nil, nil
	}

	return &dest, nil
}

func (s *Store) CreateTemplate(
	ctx context.Context,
	tx qrm.DB,
	tmpl *documentstemplates.Template,
	creatorJob string,
	categoryID *int64,
) (int64, error) {
	tDTemplates := table.FivenetDocumentsTemplates
	stmt := tDTemplates.
		INSERT(
			tDTemplates.Weight,
			tDTemplates.CategoryID,
			tDTemplates.Title,
			tDTemplates.Description,
			tDTemplates.Color,
			tDTemplates.Icon,
			tDTemplates.ContentTitle,
			tDTemplates.Content,
			tDTemplates.State,
			tDTemplates.Access,
			tDTemplates.Schema,
			tDTemplates.Workflow,
			tDTemplates.Approval,
			tDTemplates.CreatorJob,
		).
		VALUES(
			tmpl.GetWeight(),
			categoryID,
			tmpl.GetTitle(),
			tmpl.GetDescription(),
			tmpl.Color,
			tmpl.Icon,
			tmpl.GetContentTitle(),
			tmpl.GetContent(),
			tmpl.GetState(),
			tmpl.GetContentAccess(),
			tmpl.GetSchema(),
			tmpl.GetWorkflow(),
			tmpl.GetApproval(),
			creatorJob,
		)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return 0, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (s *Store) UpdateTemplate(
	ctx context.Context,
	tx qrm.DB,
	tmpl *documentstemplates.Template,
	categoryID *int64,
) error {
	tDTemplates := table.FivenetDocumentsTemplates
	stmt := tDTemplates.
		UPDATE(
			tDTemplates.Weight,
			tDTemplates.CategoryID,
			tDTemplates.Title,
			tDTemplates.Description,
			tDTemplates.Color,
			tDTemplates.Icon,
			tDTemplates.ContentTitle,
			tDTemplates.Content,
			tDTemplates.State,
			tDTemplates.Access,
			tDTemplates.Schema,
			tDTemplates.Workflow,
			tDTemplates.Approval,
		).
		SET(
			tmpl.GetWeight(),
			categoryID,
			tmpl.GetTitle(),
			tmpl.GetDescription(),
			tmpl.Color,
			tmpl.Icon,
			tmpl.GetContentTitle(),
			tmpl.GetContent(),
			tmpl.GetState(),
			tmpl.GetContentAccess(),
			tmpl.GetSchema(),
			tmpl.GetWorkflow(),
			tmpl.GetApproval(),
		).
		WHERE(
			tDTemplates.ID.EQ(mysql.Int64(tmpl.GetId())),
		).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) DeleteTemplate(
	ctx context.Context,
	tx qrm.DB,
	templateID int64,
	creatorJob string,
) error {
	tDTemplates := table.FivenetDocumentsTemplates
	stmt := tDTemplates.
		UPDATE(
			tDTemplates.DeletedAt,
		).
		SET(
			tDTemplates.DeletedAt.SET(mysql.CURRENT_TIMESTAMP()),
		).
		WHERE(mysql.AND(
			tDTemplates.CreatorJob.EQ(mysql.String(creatorJob)),
			tDTemplates.ID.EQ(mysql.Int64(templateID)),
		)).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}
