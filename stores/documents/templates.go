package documentsstore

import (
	"context"
	"errors"

	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	documentstemplates "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/templates"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) ListTemplates(
	ctx context.Context,
	ownJobOnly bool,
	userInfo *userinfo.UserInfo,
) ([]*documentstemplates.TemplateShort, error) {
	tTemplates := table.FivenetDocumentsTemplates.AS("template_short")
	tCategory := table.FivenetDocumentsCategories.AS("category")

	columns := []mysql.Projection{
		tTemplates.ID,
		tTemplates.Weight,
		tCategory.ID,
		tCategory.CreatedAt,
		tCategory.Name,
		tCategory.Description,
		tCategory.Job,
		tCategory.Color,
		tCategory.Icon,
		tTemplates.Title,
		tTemplates.Description,
		tTemplates.Color,
		tTemplates.Icon,
		tTemplates.Schema,
		tTemplates.Workflow,
		tTemplates.Approval,
		tTemplates.CreatorJob,
	}

	if userInfo != nil && userInfo.GetSuperuser() {
		columns = append(columns, tTemplates.DeletedAt)
	}

	selectStmt := tTemplates.
		SELECT(
			columns[0],
			columns[1:]...,
		)

	orderBys := []mysql.OrderByClause{
		tTemplates.Weight.DESC(),
		tTemplates.ID.ASC(),
	}

	var stmt mysql.Statement
	if userInfo != nil && userInfo.GetSuperuser() {
		stmt = selectStmt.
			FROM(tTemplates.
				LEFT_JOIN(tCategory,
					tCategory.ID.EQ(tTemplates.CategoryID),
				),
			).
			ORDER_BY(orderBys...).
			WHERE(tTemplates.CreatorJob.EQ(mysql.String(userInfo.GetJob())))
	} else {
		visibileCondition := table.FivenetDocumentsTemplates.DeletedAt.IS_NULL()
		if ownJobOnly {
			visibileCondition = visibileCondition.AND(
				table.FivenetDocumentsTemplates.CreatorJob.EQ(mysql.String(userInfo.GetJob())),
			)
		}

		visibleIDs := s.templateAccess.VisibleIDsByConditionQuery(
			userInfo,
			int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
			false,
			visibileCondition,
		)
		visibleTemplateID := mysql.IntegerColumn("id").From(visibleIDs.Table)
		stmt = selectStmt.
			FROM(visibleIDs.Table.
				INNER_JOIN(tTemplates,
					tTemplates.ID.EQ(visibleTemplateID),
				).
				LEFT_JOIN(tCategory,
					tCategory.ID.EQ(tTemplates.CategoryID),
				),
			).
			ORDER_BY(orderBys...)

		if len(visibleIDs.CTEs) > 0 {
			stmt = mysql.WITH(visibleIDs.CTEs...)(stmt)
		}
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
	includeDeleted bool,
) (*documentstemplates.Template, error) {
	tDTemplates := table.FivenetDocumentsTemplates.AS("template")
	tDCategory := table.FivenetDocumentsCategories.AS("category")

	condition := tDTemplates.ID.EQ(mysql.Int64(templateID))
	if includeDeleted {
		condition = condition.AND(tDTemplates.DeletedAt.IS_NULL())
	}

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
		WHERE(condition).
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
			dbutils.StringPEmpty(tmpl.Color),
			dbutils.StringPEmpty(tmpl.Icon),
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
	deletedAtTime *timestamp.Timestamp,
) error {
	tDTemplates := table.FivenetDocumentsTemplates
	stmt := tDTemplates.
		UPDATE(
			tDTemplates.DeletedAt,
		).
		SET(
			tDTemplates.DeletedAt.SET(dbutils.TimestampToMySQL(deletedAtTime)),
		).
		WHERE(mysql.AND(
			tDTemplates.CreatorJob.EQ(mysql.String(creatorJob)),
			tDTemplates.ID.EQ(mysql.Int64(templateID)),
		)).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}
