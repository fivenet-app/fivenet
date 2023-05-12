package perms

import (
	"errors"
	"fmt"

	"github.com/galexrt/fivenet/gen/go/proto/resources/permissions"
	"github.com/galexrt/fivenet/pkg/utils/dbutils"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	jsoniter "github.com/json-iterator/go"
)

type AttributeTypes string

const (
	StringListAttributeType  AttributeTypes = "StringList"
	JobListAttributeType     AttributeTypes = "JobList"
	JobRankListAttributeType AttributeTypes = "JobRankList"
)

type Key string

type StringList []string

type JobList StringList
type JobRankList map[string]int32

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var (
	tAttrs     = table.FivenetAttrs
	tRoleAttrs = table.FivenetRoleAttrs
)

func (p *Perms) GetAttribute(permId uint64, key Key) (*model.FivenetAttrs, error) {
	stmt := tAttrs.
		SELECT(
			tAttrs.AllColumns,
		).
		FROM(tAttrs).
		WHERE(jet.AND(
			tAttrs.PermissionID.EQ(jet.Uint64(permId)),
			tAttrs.Key.EQ(jet.String(string(key))),
		)).
		LIMIT(1)

	var dest model.FivenetAttrs
	err := stmt.QueryContext(p.ctx, p.db, &dest)
	if err != nil {
		return nil, err
	}

	return &dest, nil
}

func (p *Perms) CreateAttribute(permId uint64, key Key, aType AttributeTypes, validValues string) (uint64, error) {
	validV := jet.NULL
	if validValues != "" {
		validV = jet.String(validValues)
	}

	stmt := tAttrs.
		INSERT(
			tAttrs.PermissionID,
			tAttrs.Key,
			tAttrs.Type,
			tAttrs.ValidValues,
		).
		VALUES(
			permId,
			key,
			aType,
			validV,
		)

	res, err := stmt.ExecContext(p.ctx, p.db)
	if err != nil {
		if !dbutils.IsDuplicateError(err) {
			return 0, err
		}

		attr, err := p.GetAttribute(permId, key)
		if err != nil {
			return 0, err
		}

		return attr.ID, nil
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(lastId), nil
}

func (p *Perms) UpdateAttribute(attributeId uint64, permId uint64, key Key, aType AttributeTypes, validValues string) error {
	validV := jet.StringExp(jet.NULL)
	if validValues != "" {
		validV = jet.String(validValues)
	}

	stmt := tAttrs.
		UPDATE(
			tAttrs.PermissionID,
			tAttrs.Key,
			tAttrs.Type,
			tAttrs.ValidValues,
		).
		SET(
			tAttrs.PermissionID.SET(jet.Uint64(permId)),
			tAttrs.Key.SET(jet.String(string(key))),
			tAttrs.Type.SET(jet.String(string(aType))),
			tAttrs.ValidValues.SET(validV),
		).
		WHERE(
			tAttrs.ID.EQ(jet.Uint64(attributeId)),
		)

	_, err := stmt.ExecContext(p.ctx, p.db)
	return err
}

func (p *Perms) Attr(userId int32, job string, grade int32, category Category, name Name, key Key) (any, error) {
	roleId, ok := p.getRoleIDForJobAndGrade(job, grade)
	if !ok {
		return nil, nil
	}

	stmt := tRoleAttrs.
		SELECT(
			tRoleAttrs.RoleID.AS("roleattribute.role_id"),
			tAttrs.PermissionID.AS("roleattribute.permission_id"),
			tAttrs.Key.AS("roleattribute.key"),
			tAttrs.Type.AS("roleattribute.type"),
			tRoleAttrs.Value.AS("roleattribute.value"),
		).
		FROM(
			tRoleAttrs.
				INNER_JOIN(tAttrs,
					tAttrs.ID.EQ(tRoleAttrs.AttrID),
				),
		).
		WHERE(jet.AND(
			tRoleAttrs.RoleID.EQ(jet.Uint64(roleId)),
			tAttrs.Key.EQ(jet.String(string(key))),
		)).
		LIMIT(1)

	var dest permissions.RoleAttribute
	if err := stmt.QueryContext(p.ctx, p.db, &dest); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}

		return nil, nil
	}

	switch AttributeTypes(dest.Type) {
	case JobListAttributeType:
		fallthrough
	case StringListAttributeType:
		x := StringList{}
		if err := json.UnmarshalFromString(dest.Value, &x); err != nil {
			return nil, err
		}
		return x, nil
	case JobRankListAttributeType:
		x := JobRankList{}
		if err := json.UnmarshalFromString(dest.Value, &x); err != nil {
			return nil, err
		}
		return x, nil
	}

	return nil, fmt.Errorf("invalid permission attribute type: %q", dest.Type)
}

func (p *Perms) GetRoleAttributes(job string, grade int32) ([]*permissions.RoleAttribute, error) {
	roleIds, ok := p.getRoleIDsForJobUpToGrade(job, grade)
	if !ok {
		return nil, nil
	}

	ids := make([]jet.Expression, len(roleIds))
	for i := 0; i < len(roleIds); i++ {
		ids[i] = jet.Uint64(roleIds[i])
	}

	stmt := tRoleAttrs.
		SELECT(
			tAttrs.ID.AS("roleattribute.attr_id"),
			tRoleAttrs.RoleID.AS("roleattribute.role_id"),
			tAttrs.PermissionID.AS("roleattribute.permission_id"),
			tPerms.Category.AS("roleattribute.category"),
			tPerms.Name.AS("roleattribute.name"),
			tAttrs.Key.AS("roleattribute.key"),
			tAttrs.Type.AS("roleattribute.type"),
			tRoleAttrs.Value.AS("roleattribute.value"),
			tAttrs.ValidValues.AS("roleattribute.valid_values"),
		).
		FROM(
			tRoleAttrs.
				LEFT_JOIN(tAttrs,
					tAttrs.ID.EQ(tRoleAttrs.AttrID),
				).
				LEFT_JOIN(tPerms,
					tPerms.ID.EQ(tAttrs.PermissionID),
				),
		).
		WHERE(jet.AND(
			tRoleAttrs.RoleID.IN(ids...),
		))

	fmt.Println(stmt.DebugSql())

	var dest []*permissions.RoleAttribute
	if err := stmt.QueryContext(p.ctx, p.db, &dest); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	attrIds := make([]jet.Expression, len(dest))
	for i := 0; i < len(dest); i++ {
		attrIds[i] = jet.Uint64(dest[i].AttrId)
	}

	othersStmt := tAttrs.
		SELECT(
			tAttrs.ID.AS("roleattribute.attr_id"),
			jet.Uint64(0).AS("roleattribute.role_id"),
			tAttrs.PermissionID.AS("roleattribute.permission_id"),
			tPerms.Category.AS("roleattribute.category"),
			tPerms.Name.AS("roleattribute.name"),
			tAttrs.Key.AS("roleattribute.key"),
			tAttrs.Type.AS("roleattribute.type"),
			jet.String("").AS("roleattribute.value"),
			tAttrs.ValidValues.AS("roleattribute.valid_values"),
		).
		FROM(tAttrs.
			LEFT_JOIN(tPerms,
				tPerms.ID.EQ(tAttrs.PermissionID),
			).
			LEFT_JOIN(tRoleAttrs,
				tRoleAttrs.AttrID.EQ(tAttrs.ID),
			),
		).
		WHERE(jet.OR(
			tRoleAttrs.RoleID.IS_NULL(),
			tRoleAttrs.AttrID.NOT_IN(attrIds...),
		))

	fmt.Println(othersStmt.DebugSql())

	if err := othersStmt.QueryContext(p.ctx, p.db, &dest); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	return dest, nil
}
