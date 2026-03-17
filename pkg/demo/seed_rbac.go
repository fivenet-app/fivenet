package demo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

func (d *Demo) upsertDemoTargetJobHighestGradeRolePerms(ctx context.Context) error {
	job := d.targetJobName()

	highestGrade, ok, err := d.lookupHighestJobGrade(ctx, job)
	if err != nil {
		return err
	}
	if !ok {
		d.logger.Warn(
			"skipping demo rbac seeding, no job grades found for target job",
			zap.String("job", job),
		)
		return nil
	}
	if highestGrade <= 0 {
		fallback := d.highestJobGrade(job)
		if fallback <= 0 {
			return fmt.Errorf("invalid highest grade %d for target job %s", highestGrade, job)
		}
		d.logger.Warn(
			"rbac highest-grade lookup returned non-positive value, using fallback from demo catalog",
			zap.String("job", job),
			zap.Int32("lookup_grade", highestGrade),
			zap.Int32("fallback_grade", fallback),
		)
		highestGrade = fallback
	}

	roleID, err := d.upsertRoleForJobGrade(ctx, job, highestGrade)
	if err != nil {
		return err
	}

	permissionIDs, err := d.lookupAllPermissionIDs(ctx)
	if err != nil {
		return err
	}
	if len(permissionIDs) == 0 {
		d.logger.Warn("skipping demo rbac role permission seeding, no permissions found")
		return nil
	}

	if err := d.upsertDemoTargetJobPerms(ctx, job, permissionIDs); err != nil {
		return err
	}

	attrs, err := d.lookupRbacAttrs(ctx)
	if err != nil {
		return err
	}
	if err := d.upsertDemoTargetJobAttrs(ctx, job, highestGrade, attrs); err != nil {
		return err
	}

	if err := d.upsertDemoTargetJobRolePerms(
		ctx,
		roleID,
		job,
		highestGrade,
		permissionIDs,
	); err != nil {
		return err
	}
	if err := d.upsertDemoTargetJobRoleAttrs(ctx, roleID, job, highestGrade, attrs); err != nil {
		return err
	}

	d.logger.Info(
		"completed demo rbac role seeding for target job",
		zap.String("job", job),
		zap.Int32("grade", highestGrade),
		zap.Int64("role_id", roleID),
		zap.Int("permissions_count", len(permissionIDs)),
	)

	return nil
}

func (d *Demo) upsertDemoTargetJobRolePerms(
	ctx context.Context,
	roleID int64,
	job string,
	grade int32,
	permissionIDs []int64,
) error {
	stmt := tRbacRolePerms.
		INSERT(
			tRbacRolePerms.RoleID,
			tRbacRolePerms.PermissionID,
			tRbacRolePerms.Val,
		)
	for _, permissionID := range permissionIDs {
		stmt = stmt.VALUES(roleID, permissionID, true)
	}

	stmt = stmt.ON_DUPLICATE_KEY_UPDATE(
		tRbacRolePerms.Val.SET(mysql.RawBool("VALUES(`val`)")),
	)

	if _, err := stmt.ExecContext(ctx, d.db); err != nil {
		return fmt.Errorf(
			"failed to upsert demo rbac role permissions for job %s grade %d. %w",
			job,
			grade,
			err,
		)
	}

	return nil
}

func (d *Demo) upsertDemoTargetJobRoleAttrs(
	ctx context.Context,
	roleID int64,
	job string,
	grade int32,
	attrs []rbacAttrDef,
) error {
	if len(attrs) == 0 {
		return nil
	}

	stmt := tRbacRoleAttrs.
		INSERT(
			tRbacRoleAttrs.RoleID,
			tRbacRoleAttrs.AttrID,
			tRbacRoleAttrs.Value,
		)

	for _, attr := range attrs {
		valueJSON, err := d.buildDemoRoleAttrValueJSON(attr.AttrType, attr.ValidValues, job, grade)
		if err != nil {
			return err
		}
		stmt = stmt.VALUES(roleID, attr.AttrID, valueJSON)
	}

	stmt = stmt.ON_DUPLICATE_KEY_UPDATE(
		tRbacRoleAttrs.Value.SET(mysql.RawString("VALUES(`value`)")),
	)

	if _, err := stmt.ExecContext(ctx, d.db); err != nil {
		return fmt.Errorf("failed to upsert demo rbac role attributes. %w", err)
	}

	return nil
}

func (d *Demo) upsertDemoTargetJobPerms(
	ctx context.Context,
	job string,
	permissionIDs []int64,
) error {
	stmt := tRbacJobPerms.
		INSERT(
			tRbacJobPerms.Job,
			tRbacJobPerms.PermissionID,
			tRbacJobPerms.Val,
		)
	for _, permissionID := range permissionIDs {
		stmt = stmt.VALUES(job, permissionID, true)
	}

	stmt = stmt.ON_DUPLICATE_KEY_UPDATE(
		tRbacJobPerms.Val.SET(mysql.RawBool("VALUES(`val`)")),
	)

	if _, err := stmt.ExecContext(ctx, d.db); err != nil {
		return fmt.Errorf("failed to upsert demo rbac job permissions for job %s. %w", job, err)
	}

	return nil
}

func (d *Demo) upsertDemoTargetJobAttrs(
	ctx context.Context,
	job string,
	grade int32,
	attrs []rbacAttrDef,
) error {
	if len(attrs) == 0 {
		return nil
	}

	stmt := tRbacJobAttrs.
		INSERT(
			tRbacJobAttrs.Job,
			tRbacJobAttrs.AttrID,
			tRbacJobAttrs.MaxValues,
		)

	for _, attr := range attrs {
		valueJSON, err := d.buildDemoRoleAttrValueJSON(attr.AttrType, attr.ValidValues, job, grade)
		if err != nil {
			return err
		}
		stmt = stmt.VALUES(job, attr.AttrID, valueJSON)
	}

	stmt = stmt.ON_DUPLICATE_KEY_UPDATE(
		tRbacJobAttrs.MaxValues.SET(mysql.RawString("VALUES(`max_values`)")),
	)

	if _, err := stmt.ExecContext(ctx, d.db); err != nil {
		return fmt.Errorf("failed to upsert demo rbac job attributes for job %s. %w", job, err)
	}

	return nil
}

func (d *Demo) lookupHighestJobGrade(ctx context.Context, job string) (int32, bool, error) {
	stmt := tJobsGrades.
		SELECT(tJobsGrades.Grade.AS("highest_grade")).
		FROM(tJobsGrades).
		WHERE(tJobsGrades.JobName.EQ(mysql.String(job))).
		ORDER_BY(tJobsGrades.Grade.DESC()).
		LIMIT(1)

	var row struct {
		HighestGrade int32 `alias:"highest_grade"`
	}
	if err := stmt.QueryContext(ctx, d.db, &row); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("failed to lookup highest job grade for %s. %w", job, err)
	}

	return row.HighestGrade, true, nil
}

func (d *Demo) upsertRoleForJobGrade(ctx context.Context, job string, grade int32) (int64, error) {
	stmt := tRbacRoles.
		INSERT(
			tRbacRoles.Job,
			tRbacRoles.Grade,
		).
		VALUES(job, grade).
		ON_DUPLICATE_KEY_UPDATE(
			tRbacRoles.Job.SET(mysql.RawString("VALUES(`job`)")),
			tRbacRoles.Grade.SET(mysql.RawInt("VALUES(`grade`)")),
		)

	if _, err := stmt.ExecContext(ctx, d.db); err != nil {
		return 0, fmt.Errorf("failed to upsert rbac role for job %s grade %d. %w", job, grade, err)
	}

	selectStmt := tRbacRoles.
		SELECT(tRbacRoles.ID.AS("role_id")).
		FROM(tRbacRoles).
		WHERE(mysql.AND(
			tRbacRoles.Job.EQ(mysql.String(job)),
			tRbacRoles.Grade.EQ(mysql.Int32(grade)),
		)).
		LIMIT(1)

	var row struct {
		RoleID int64 `alias:"role_id"`
	}
	if err := selectStmt.QueryContext(ctx, d.db, &row); err != nil {
		return 0, fmt.Errorf(
			"failed to lookup rbac role id for job %s grade %d. %w",
			job,
			grade,
			err,
		)
	}
	if row.RoleID <= 0 {
		return 0, fmt.Errorf("failed to resolve valid rbac role id for job %s grade %d", job, grade)
	}

	return row.RoleID, nil
}

func (d *Demo) lookupAllPermissionIDs(ctx context.Context) ([]int64, error) {
	stmt := tRbacPermissions.
		SELECT(tRbacPermissions.ID).
		FROM(tRbacPermissions).
		ORDER_BY(tRbacPermissions.ID.ASC())

	ids := []int64{}
	if err := stmt.QueryContext(ctx, d.db, &ids); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return ids, nil
		}
		return nil, fmt.Errorf("failed to lookup rbac permission ids. %w", err)
	}

	return ids, nil
}

type rbacAttrDef struct {
	AttrID      int64   `alias:"attr_id"`
	AttrType    string  `alias:"attr_type"`
	ValidValues *string `alias:"valid_values"`
}

func (d *Demo) lookupRbacAttrs(ctx context.Context) ([]rbacAttrDef, error) {
	stmt := tRbacAttrs.
		SELECT(
			tRbacAttrs.ID.AS("rbacattrdef.attr_id"),
			tRbacAttrs.Type.AS("rbacattrdef.attr_type"),
			tRbacAttrs.ValidValues.AS("rbacattrdef.valid_values"),
		).
		FROM(tRbacAttrs).
		ORDER_BY(tRbacAttrs.ID.ASC())

	attrs := []rbacAttrDef{}
	if err := stmt.QueryContext(ctx, d.db, &attrs); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return attrs, nil
		}
		return nil, fmt.Errorf("failed to lookup rbac attributes. %w", err)
	}

	return attrs, nil
}

func (d *Demo) buildDemoRoleAttrValueJSON(
	attrTypeRaw string,
	validValuesRaw *string,
	targetJob string,
	targetGrade int32,
) (string, error) {
	attrType := permissionsattributes.AttributeTypes(strings.TrimSpace(attrTypeRaw))

	value := &permissionsattributes.AttributeValues{}
	switch attrType {
	case permissionsattributes.StringListAttributeType:
		valid := &permissionsattributes.AttributeValues{}
		if validValuesRaw != nil && strings.TrimSpace(*validValuesRaw) != "" {
			if err := protoutils.UnmarshalPartialJSON([]byte(*validValuesRaw), valid); err != nil {
				return "", fmt.Errorf("failed to parse rbac attr valid values. %w", err)
			}
		}
		valid.Default(attrType)
		value.Default(attrType)
		value.GetStringList().Strings = append(
			value.GetStringList().Strings,
			valid.GetStringList().GetStrings()...,
		)

	case permissionsattributes.JobListAttributeType:
		value.Default(attrType)
		value.GetJobList().Strings = []string{targetJob}

	case permissionsattributes.JobGradeListAttributeType:
		value.Default(attrType)
		value.GetJobGradeList().FineGrained = false
		value.GetJobGradeList().Jobs = map[string]int32{
			targetJob: targetGrade,
		}

	default:
		return "", fmt.Errorf("unsupported rbac attribute type %q", attrTypeRaw)
	}

	out, err := protoutils.MarshalToJSON(value)
	if err != nil {
		return "", fmt.Errorf("failed to marshal rbac role attribute value. %w", err)
	}

	return string(out), nil
}
