package perms

import (
	"errors"
	"fmt"

	"github.com/galexrt/fivenet/gen/go/proto/resources/permissions"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/pkg/utils/dbutils"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	jsoniter "github.com/json-iterator/go"
)

type AttributeTypes string

const (
	StringListAttributeType   AttributeTypes = "StringList"
	JobListAttributeType      AttributeTypes = "JobList"
	JobGradeListAttributeType AttributeTypes = "JobGradeList"
)

type Key string

type StringList []string

type JobList []string
type JobGradeList map[string]int32

func (l StringList) Validate(validVals []string) bool {
	// If more values than valid values in the list, it can't be valid
	if len(l) > len(validVals) {
		return false
	}

	for i := 0; i < len(l); i++ {
		if !utils.InStringSlice(validVals, l[i]) {
			return false
		}
	}

	return true
}

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

func (p *Perms) CreateAttribute(permId uint64, key Key, aType AttributeTypes, validValues any) (uint64, error) {
	validV := jet.NULL
	if validValues != nil {
		out, err := json.MarshalToString(validValues)
		if err != nil {
			return 0, err
		}

		validV = jet.String(out)
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

func (p *Perms) UpdateAttribute(attrId uint64, permId uint64, key Key, aType AttributeTypes, validValues any) error {
	validV := jet.StringExp(jet.NULL)
	if validValues != nil {
		out, err := json.MarshalToString(validValues)
		if err != nil {
			return err
		}

		validV = jet.String(out)
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
			tAttrs.ID.EQ(jet.Uint64(attrId)),
		)

	_, err := stmt.ExecContext(p.ctx, p.db)
	if err != nil {
		return err
	}

	return nil
}

func (p *Perms) Attr(userId int32, job string, grade int32, category Category, name Name, key Key) (any, error) {
	roleId, ok := p.getRoleIDForJobAndGrade(job, grade)
	if !ok {
		return nil, nil
	}

	permId, ok := p.lookupPermIDByGuard(BuildGuard(category, name))
	if !ok {
		return nil, nil
	}

	// TODO Attributes should be inheritable use a cache like for role perms

	stmt := tRoleAttrs.
		SELECT(
			tAttrs.Key.AS("key"),
			tAttrs.Type.AS("type"),
			tRoleAttrs.Value.AS("value"),
		).
		FROM(
			tRoleAttrs.
				INNER_JOIN(tAttrs,
					tAttrs.ID.EQ(tRoleAttrs.AttrID).
						AND(tAttrs.PermissionID.EQ(jet.Uint64(permId))),
				),
		).
		WHERE(jet.AND(
			tRoleAttrs.RoleID.EQ(jet.Uint64(roleId)),
			tAttrs.Key.EQ(jet.String(string(key))),
		)).
		LIMIT(1)

	var dest struct {
		Type  AttributeTypes
		Value string
	}
	if err := stmt.QueryContext(p.ctx, p.db, &dest); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}

		return nil, nil
	}

	return p.convertAttributeValue(dest.Value, dest.Type)
}

func (p *Perms) convertAttributeValue(val string, aType AttributeTypes) (any, error) {
	switch AttributeTypes(aType) {
	case JobListAttributeType:
		fallthrough
	case StringListAttributeType:
		x := StringList{}
		if err := json.UnmarshalFromString(val, &x); err != nil {
			return nil, err
		}
		return x, nil
	case JobGradeListAttributeType:
		x := JobGradeList{}
		if err := json.UnmarshalFromString(val, &x); err != nil {
			return nil, err
		}
		return x, nil
	}

	return nil, fmt.Errorf("invalid permission attribute type: %q", aType)
}

func (p *Perms) convertRawToRoleAttributes(in []*permissions.RawAttribute) ([]*permissions.RoleAttribute, error) {
	res := make([]*permissions.RoleAttribute, len(in))
	for i := 0; i < len(in); i++ {
		res[i] = &permissions.RoleAttribute{
			RoleId:       in[i].RoleId,
			CreatedAt:    in[i].CreatedAt,
			AttrId:       in[i].AttrId,
			PermissionId: in[i].PermissionId,
			Category:     in[i].Category,
			Name:         in[i].Name,
			Key:          in[i].Key,
			Type:         in[i].Type,
			Value:        &permissions.AttributeValues{},
			ValidValues:  &permissions.AttributeValues{},
		}

		if err := p.convertRawValue(res[i].Value, in[i].RawValue, AttributeTypes(res[i].Type)); err != nil {
			return nil, err
		}

		if err := p.convertRawValue(res[i].ValidValues, in[i].RawValidValues, AttributeTypes(res[i].Type)); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (p *Perms) convertRawValue(targetVal *permissions.AttributeValues, rawVal string, aType AttributeTypes) error {
	var val any
	var err error

	if rawVal != "" {
		val, err = p.convertAttributeValue(rawVal, AttributeTypes(aType))
		if err != nil {
			return err
		}
	}

	switch AttributeTypes(aType) {
	case StringListAttributeType:
		pVal := &permissions.AttributeValues_StringList{
			StringList: &permissions.StringList{},
		}
		if val != nil {
			pVal.StringList.Strings = val.(StringList)
		}
		targetVal.ValidValues = pVal
	case JobListAttributeType:
		pVal := &permissions.AttributeValues_JobList{
			JobList: &permissions.StringList{},
		}
		if val != nil {
			pVal.JobList.Strings = val.(JobList)
		}
		targetVal.ValidValues = pVal
	case JobGradeListAttributeType:
		pVal := &permissions.AttributeValues_JobGradeList{
			JobGradeList: &permissions.JobGradeList{},
		}
		if val != nil {
			pVal.JobGradeList.Jobs = val.(JobGradeList)
		}
		targetVal.ValidValues = pVal
	}

	return nil
}

func (p *Perms) GetAllAttributes(job string) ([]*permissions.RoleAttribute, error) {
	stmt := tAttrs.
		SELECT(
			tAttrs.ID.AS("rawattribute.attr_id"),
			tAttrs.PermissionID.AS("rawattribute.permission_id"),
			tPerms.Category.AS("rawattribute.category"),
			tPerms.Name.AS("rawattribute.name"),
			tAttrs.Key.AS("rawattribute.key"),
			tAttrs.Type.AS("rawattribute.type"),
			tAttrs.ValidValues.AS("rawattribute.valid_values"),
		).
		FROM(tAttrs.
			INNER_JOIN(tPerms,
				tPerms.ID.EQ(tAttrs.PermissionID),
			),
		)

	var dest []*permissions.RawAttribute
	if err := stmt.QueryContext(p.ctx, p.db, &dest); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	return p.convertRawToRoleAttributes(dest)
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
			tAttrs.ID.AS("rawattribute.attr_id"),
			tRoleAttrs.RoleID.AS("rawattribute.role_id"),
			tAttrs.PermissionID.AS("rawattribute.permission_id"),
			tPerms.Category.AS("rawattribute.category"),
			tPerms.Name.AS("rawattribute.name"),
			tAttrs.Key.AS("rawattribute.key"),
			tAttrs.Type.AS("rawattribute.type"),
			tRoleAttrs.Value.AS("rawattribute.raw_value"),
			tAttrs.ValidValues.AS("rawattribute.raw_valid_values"),
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

	var dest []*permissions.RawAttribute
	if err := stmt.QueryContext(p.ctx, p.db, &dest); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	return p.convertRawToRoleAttributes(dest)
}

func (p *Perms) AddAttributesToRole(roleId uint64, attrs ...*permissions.RoleAttribute) error {
	stmt := tRoleAttrs.
		INSERT().
		MODELS(attrs).
		ON_DUPLICATE_KEY_UPDATE(
			tRoleAttrs.Value.SET(jet.StringExp(jet.Raw("values(`value`)"))),
		)

	if _, err := stmt.ExecContext(p.ctx, p.db); err != nil {
		if err != nil && !dbutils.IsDuplicateError(err) {
			return err
		}
	}

	return nil
}

func (p *Perms) UpdateRoleAttributes(roleId uint64, attrs ...*permissions.RoleAttribute) error {
	ids := make([]jet.Expression, len(attrs))
	for i := 0; i < len(attrs); i++ {
		ids[i] = jet.Uint64(attrs[i].AttrId)
	}

	stmt := tRoleAttrs.
		UPDATE(
			tRoleAttrs.Value,
		).
		SET(
			tRoleAttrs.Value,
		).
		WHERE(jet.AND(
			tRoleAttrs.RoleID.EQ(jet.Uint64(roleId)),
			tRoleAttrs.AttrID.IN(ids...),
		))

	if _, err := stmt.ExecContext(p.ctx, p.db); err != nil {
		if err != nil && !dbutils.IsDuplicateError(err) {
			return err
		}
	}

	return nil
}

func (p *Perms) RemoveAttributesFromRole(roleId uint64, attrs ...*permissions.RoleAttribute) error {
	ids := make([]jet.Expression, len(attrs))
	for i := 0; i < len(attrs); i++ {
		ids[i] = jet.Uint64(attrs[i].AttrId)
	}

	stmt := tRoleAttrs.
		DELETE().
		WHERE(jet.AND(
			tRoleAttrs.RoleID.EQ(jet.Uint64(roleId)),
			tRoleAttrs.AttrID.IN(ids...),
		))

	if _, err := stmt.ExecContext(p.ctx, p.db); err != nil {
		return err
	}

	return nil
}
