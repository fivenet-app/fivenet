package perms

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/galexrt/fivenet/gen/go/proto/resources/permissions"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/perms/helpers"
	"github.com/galexrt/fivenet/pkg/utils/dbutils"
	"github.com/galexrt/fivenet/pkg/utils/syncx"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	jsoniter "github.com/json-iterator/go"
)

var (
	ErrAttrInvalid = errors.New("invalid attributes")
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var (
	tAttrs     = table.FivenetAttrs
	tRoleAttrs = table.FivenetRoleAttrs
)

func (p *Perms) GetAttribute(category Category, name Name, key Key) (*permissions.RoleAttribute, error) {
	permId, ok := p.lookupPermIDByGuard(BuildGuard(category, name))
	if !ok {
		return nil, fmt.Errorf("unable to find perm ID by attribute")
	}

	attr, ok := p.lookupAttributeByPermID(permId, key)
	if !ok {
		return nil, fmt.Errorf("no attribute found by id")
	}

	return &permissions.RoleAttribute{
		AttrId:       attr.ID,
		PermissionId: attr.PermissionID,
		Category:     string(category),
		Name:         string(name),
		Key:          string(attr.Key),
		Type:         string(attr.Type),
		ValidValues:  attr.ValidValues,
		MaxValues:    nil,
	}, nil
}

func (p *Perms) GetAttributeByIDs(ctx context.Context, attrIds ...uint64) ([]*permissions.RoleAttribute, error) {
	ids := make([]jet.Expression, len(attrIds))
	for i := 0; i < len(attrIds); i++ {
		ids[i] = jet.Uint64(attrIds[i])
	}

	stmt := tAttrs.
		SELECT(
			tAttrs.AllColumns,
		).
		FROM(tAttrs).
		WHERE(jet.AND(
			tAttrs.ID.IN(ids...),
		)).
		LIMIT(1)

	var dest []*model.FivenetAttrs
	err := stmt.QueryContext(ctx, p.db, &dest)
	if err != nil {
		return nil, err
	}

	attrs := make([]*permissions.RoleAttribute, len(dest))
	for i := 0; i < len(dest); i++ {
		attr, ok := p.LookupAttributeByID(dest[i].ID)
		if !ok {
			return nil, fmt.Errorf("no attribute found by id")
		}

		attrs[i] = &permissions.RoleAttribute{
			AttrId:       dest[i].ID,
			PermissionId: dest[i].PermissionID,
			Key:          dest[i].Key,
			Type:         dest[i].Type,
			Category:     string(attr.Category),
			Name:         string(attr.Name),
			ValidValues:  attr.ValidValues,
			MaxValues:    nil,
		}
	}

	return attrs, nil
}

func (p *Perms) getAttributeFromDatabase(ctx context.Context, permId uint64, key Key) (*model.FivenetAttrs, error) {
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
	err := stmt.QueryContext(ctx, p.db, &dest)
	if err != nil {
		return nil, err
	}

	return &dest, nil
}

func (p *Perms) CreateAttribute(ctx context.Context, permId uint64, key Key, aType permissions.AttributeTypes, validValues any, defaultValues any) (uint64, error) {
	validV := jet.NULL
	if validValues != nil {
		out, err := json.MarshalToString(validValues)
		if err != nil {
			return 0, err
		}

		if out != "" && out != "null" {
			validV = jet.String(out)
		}
	}

	defaultV := jet.NULL
	if defaultValues != nil {
		out, err := json.MarshalToString(defaultValues)
		if err != nil {
			return 0, err
		}

		if out != "" && out != "null" {
			defaultV = jet.String(out)
		}
	}

	stmt := tAttrs.
		INSERT(
			tAttrs.PermissionID,
			tAttrs.Key,
			tAttrs.Type,
			tAttrs.ValidValues,
			tAttrs.DefaultValues,
		).
		VALUES(
			permId,
			key,
			aType,
			validV,
			defaultV,
		)

	res, err := stmt.ExecContext(ctx, p.db)
	if err != nil {
		if !dbutils.IsDuplicateError(err) {
			return 0, err
		}

		attr, err := p.getAttributeFromDatabase(ctx, permId, key)
		if err != nil {
			return 0, err
		}

		if err := p.addOrUpdateAttributeInMap(permId, uint64(attr.ID), key, aType, validValues, defaultValues); err != nil {
			return 0, err
		}

		return attr.ID, nil
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if err := p.addOrUpdateAttributeInMap(permId, uint64(lastId), key, aType, validValues, defaultValues); err != nil {
		return 0, err
	}

	return uint64(lastId), nil
}

func (p *Perms) addOrUpdateAttributeInMap(permId uint64, attrId uint64, key Key, aType permissions.AttributeTypes, validValues any, defaultValues any) error {
	var validValsOut string
	var err error
	// if the valid values is a nil or a string, don't do anything extra just set to an empty string
	if validValues != nil {
		vType := reflect.TypeOf(validValues).String()
		if vType == "string" {
			if validValues != "" {
				validValsOut = validValues.(string)
			}
		} else {
			validValsOut, err = json.MarshalToString(validValues)
			if err != nil {
				return err
			}
		}
	}

	var defaultValsOut string
	// if the valid values is a nil or a string, don't do anything extra just set to an empty string
	if validValues != nil {
		vType := reflect.TypeOf(validValues).String()
		if vType == "string" {
			if validValues != "" {
				defaultValsOut = validValues.(string)
			}
		} else {
			defaultValsOut, err = json.MarshalToString(validValues)
			if err != nil {
				return err
			}
		}
	}

	validVals := &permissions.AttributeValues{}
	if err := p.convertRawValue(validVals, validValsOut, aType); err != nil {
		return err
	}

	defaultVals := &permissions.AttributeValues{}
	if err := p.convertRawValue(defaultVals, defaultValsOut, aType); err != nil {
		return err
	}

	if err := p.updateAttributeInMap(permId, attrId, key, aType, validVals, defaultVals); err != nil {
		return err
	}

	return nil
}

func (p *Perms) updateAttributeInMap(permId uint64, attrId uint64, key Key, aType permissions.AttributeTypes, validValues *permissions.AttributeValues, defaultValues *permissions.AttributeValues) error {
	perm, ok := p.lookupPermByID(permId)
	if !ok {
		return fmt.Errorf("no permission found by id")
	}

	attr := &cacheAttr{
		ID:            attrId,
		PermissionID:  permId,
		Category:      perm.Category,
		Name:          perm.Name,
		Key:           key,
		Type:          aType,
		ValidValues:   validValues,
		DefaultValues: defaultValues,
	}

	p.attrsMap.Store(attrId, attr)

	pAttrMap, _ := p.attrsPermsMap.LoadOrStore(permId, &syncx.Map[Key, uint64]{})
	pAttrMap.Store(key, attrId)

	return nil
}

func (p *Perms) UpdateAttribute(ctx context.Context, attrId uint64, permId uint64, key Key, aType permissions.AttributeTypes, validValues any, defaultValues any) error {
	validV := jet.StringExp(jet.NULL)
	if validValues != nil {
		out, err := json.MarshalToString(validValues)
		if err != nil {
			return err
		}

		if strings.ToLower(out) != "null" {
			validV = jet.String(out)
		}
	}

	defaultV := jet.StringExp(jet.NULL)
	if defaultValues != nil {
		out, err := json.MarshalToString(defaultValues)
		if err != nil {
			return err
		}

		if strings.ToLower(out) != "null" {
			defaultV = jet.String(out)
		}
	}

	stmt := tAttrs.
		UPDATE(
			tAttrs.PermissionID,
			tAttrs.Key,
			tAttrs.Type,
			tAttrs.ValidValues,
			tAttrs.DefaultValues,
		).
		SET(
			tAttrs.PermissionID.SET(jet.Uint64(permId)),
			tAttrs.Key.SET(jet.String(string(key))),
			tAttrs.Type.SET(jet.String(string(aType))),
			tAttrs.ValidValues.SET(validV),
			tAttrs.DefaultValues.SET(defaultV),
		).
		WHERE(
			tAttrs.ID.EQ(jet.Uint64(attrId)),
		)

	_, err := stmt.ExecContext(ctx, p.db)
	if err != nil {
		return err
	}

	if err := p.addOrUpdateAttributeInMap(permId, attrId, key, aType, validValues, defaultValues); err != nil {
		return nil
	}

	return nil
}

func (p *Perms) getClosestRoleAttr(job string, grade int32, permId uint64, key Key) *cacheRoleAttr {
	roleIds, ok := p.lookupRoleIDsForJobUpToGrade(job, grade)
	if !ok {
		return nil
	}

	pAttrs, ok := p.attrsPermsMap.Load(permId)
	if !ok {
		return nil
	}
	attrId, ok := pAttrs.Load(key)
	if !ok {
		return nil
	}

	for i := len(roleIds) - 1; i >= 0; i-- {
		val, ok := p.lookupRoleAttribute(roleIds[i], attrId)
		if !ok {
			continue
		}

		return val
	}

	return nil
}

func (p *Perms) GetClosestRoleAttrMaxVals(job string, grade int32, permId uint64, key Key) *permissions.AttributeValues {
	roleIds, ok := p.lookupRoleIDsForJobUpToGrade(job, grade)
	if !ok {
		return nil
	}

	pAttrs, ok := p.attrsPermsMap.Load(permId)
	if !ok {
		return nil
	}
	attrId, ok := pAttrs.Load(key)
	if !ok {
		return nil
	}

	for i := len(roleIds) - 1; i >= 0; i-- {
		val, ok := p.lookupRoleAttribute(roleIds[i], attrId)
		if !ok {
			continue
		}

		if val.Max != nil {
			return val.Max
		}
	}

	return nil
}

func (p *Perms) Attr(userInfo *userinfo.UserInfo, category Category, name Name, key Key) (any, error) {
	permId, ok := p.lookupPermIDByGuard(BuildGuard(category, name))
	if !ok {
		return nil, nil
	}

	var cached *cacheRoleAttr
	cached = p.getClosestRoleAttr(userInfo.Job, userInfo.JobGrade, permId, key)
	if userInfo.SuperUser {
		attr, ok := p.lookupAttributeByPermID(permId, key)
		if !ok {
			return nil, nil
		}

		if attr.ValidValues != nil {
			cached = &cacheRoleAttr{
				AttrID:       attr.ID,
				PermissionID: attr.PermissionID,
				Key:          key,
				Type:         attr.Type,
				Value:        attr.ValidValues,
			}
		}
	}

	if cached == nil {
		return nil, nil
	}

	switch cached.Type {
	case permissions.StringListAttributeType:
		return cached.Value.GetStringList().Strings, nil
	case permissions.JobListAttributeType:
		return cached.Value.GetJobList().Strings, nil
	case permissions.JobGradeListAttributeType:
		return cached.Value.GetJobGradeList().Jobs, nil
	}

	return nil, fmt.Errorf("unknown role attribute type")
}

func (p *Perms) convertRawToAttributeValue(val string, aType permissions.AttributeTypes) (any, error) {
	switch permissions.AttributeTypes(aType) {
	case permissions.StringListAttributeType:
		x := StringList{}
		if err := json.UnmarshalFromString(val, &x); err != nil {
			return nil, err
		}
		return x, nil
	case permissions.JobListAttributeType:
		x := JobList{}
		if err := json.UnmarshalFromString(val, &x); err != nil {
			return nil, err
		}
		return x, nil
	case permissions.JobGradeListAttributeType:
		x := JobGradeList{}
		if err := json.UnmarshalFromString(val, &x); err != nil {
			return nil, err
		}
		return x, nil
	}

	return nil, fmt.Errorf("invalid permission attribute type: %q", aType)
}

func (p *Perms) convertRawToRoleAttributes(in []*permissions.RawRoleAttribute, job string, grade int32) ([]*permissions.RoleAttribute, error) {
	res := make([]*permissions.RoleAttribute, len(in))
	for i := 0; i < len(in); i++ {
		res[i] = &permissions.RoleAttribute{
			RoleId:        in[i].RoleId,
			CreatedAt:     in[i].CreatedAt,
			AttrId:        in[i].AttrId,
			PermissionId:  in[i].PermissionId,
			Category:      in[i].Category,
			Name:          in[i].Name,
			Key:           in[i].Key,
			Type:          in[i].Type,
			Value:         &permissions.AttributeValues{},
			ValidValues:   &permissions.AttributeValues{},
			DefaultValues: &permissions.AttributeValues{},
			MaxValues:     &permissions.AttributeValues{},
		}

		if err := p.convertRawValue(res[i].Value, in[i].RawValue, permissions.AttributeTypes(res[i].Type)); err != nil {
			return nil, err
		}

		if err := p.convertRawValue(res[i].ValidValues, in[i].RawValidValues, permissions.AttributeTypes(res[i].Type)); err != nil {
			return nil, err
		}

		if in[i].RawDefaultValues != nil {
			if err := p.convertRawValue(res[i].DefaultValues, *in[i].RawDefaultValues, permissions.AttributeTypes(res[i].Type)); err != nil {
				return nil, err
			}
		} else {
			res[i].DefaultValues = nil
		}

		res[i].MaxValues = p.GetClosestRoleAttrMaxVals(job, grade, in[i].PermissionId, Key(in[i].Key))
		if res[i].MaxValues == nil {
			res[i].MaxValues = &permissions.AttributeValues{}
			res[i].MaxValues.Default(permissions.AttributeTypes(res[i].Type))
		}
	}

	return res, nil
}

func (p *Perms) convertRawValue(targetVal *permissions.AttributeValues, rawVal string, aType permissions.AttributeTypes) error {
	var val any

	if rawVal != "" {
		var err error
		val, err = p.convertRawToAttributeValue(rawVal, permissions.AttributeTypes(aType))
		if err != nil {
			return err
		}
	}

	switch permissions.AttributeTypes(aType) {
	case permissions.StringListAttributeType:
		pVal := &permissions.AttributeValues_StringList{
			StringList: &permissions.StringList{},
		}
		if val != nil {
			pVal.StringList.Strings = val.(StringList)
		}
		targetVal.ValidValues = pVal
	case permissions.JobListAttributeType:
		pVal := &permissions.AttributeValues_JobList{
			JobList: &permissions.StringList{},
		}
		if val != nil {
			pVal.JobList.Strings = val.(JobList)
		}
		targetVal.ValidValues = pVal
	case permissions.JobGradeListAttributeType:
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

func (p *Perms) GetAllAttributes(ctx context.Context, job string, grade int32) ([]*permissions.RoleAttribute, error) {
	stmt := tAttrs.
		SELECT(
			tAttrs.ID.AS("rawroleattribute.attr_id"),
			tAttrs.PermissionID.AS("rawroleattribute.permission_id"),
			tPerms.Category.AS("rawroleattribute.category"),
			tPerms.Name.AS("rawroleattribute.name"),
			tAttrs.Key.AS("rawroleattribute.key"),
			tAttrs.Type.AS("rawroleattribute.type"),
			tAttrs.ValidValues.AS("rawroleattribute.valid_values"),
			tAttrs.DefaultValues.AS("rawroleattribute.default_values"),
		).
		FROM(tAttrs.
			INNER_JOIN(tPerms,
				tPerms.ID.EQ(tAttrs.PermissionID),
			),
		)

	var dest []*permissions.RawRoleAttribute
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	return p.convertRawToRoleAttributes(dest, job, grade)
}

func (p *Perms) GetRoleAttributes(job string, grade int32) ([]*permissions.RoleAttribute, error) {
	roleId, ok := p.lookupRoleIDForJobAndGrade(job, grade)
	if !ok {
		roleId, ok = p.lookupRoleIDForJobAndGrade(DefaultRoleJob, DefaultRoleJobGrade)
		if !ok {
			return nil, fmt.Errorf("failed to fallback to default role")
		}
	}

	as, ok := p.attrsRoleMap.Load(roleId)
	if !ok {
		return []*permissions.RoleAttribute{}, nil
	}

	dest := []*permissions.RoleAttribute{}
	for _, attrId := range as.Keys() {
		attr, ok := p.LookupAttributeByID(attrId)
		if !ok {
			return nil, fmt.Errorf("no attribute found by id for role")
		}

		v, ok := as.Load(attrId)
		if !ok {
			return nil, fmt.Errorf("no role attribute found by id for role")
		}

		dest = append(dest, &permissions.RoleAttribute{
			RoleId:       roleId,
			AttrId:       attrId,
			PermissionId: attr.PermissionID,
			Category:     string(attr.Category),
			Name:         string(attr.Name),
			Key:          string(attr.Key),
			Type:         string(attr.Type),
			Value:        v.Value,
			ValidValues:  attr.ValidValues,
			MaxValues:    v.Max,
		})
	}

	return dest, nil
}

func (p *Perms) getRoleAttributesFromCache(job string, grade int32) ([]*cacheRoleAttr, error) {
	roleIds, ok := p.lookupRoleIDsForJobUpToGrade(job, grade)
	if !ok {
		return []*cacheRoleAttr{}, nil
	}

	attrs := map[uint64]*cacheRoleAttr{}
	for i := len(roleIds) - 1; i >= 0; i-- {
		attrMap, ok := p.attrsRoleMap.Load(roleIds[i])
		if !ok {
			continue
		}

		attrMap.Range(func(key uint64, value *cacheRoleAttr) bool {
			if _, ok := attrs[key]; !ok {
				attrs[key] = value
			}

			return true
		})
	}

	as := []*cacheRoleAttr{}
	for _, v := range attrs {
		as = append(as, v)
	}

	return as, nil
}

func (p *Perms) FlattenRoleAttributes(job string, grade int32) ([]string, error) {
	attrs, err := p.getRoleAttributesFromCache(job, grade)
	if err != nil {
		return nil, err
	}

	as := []string{}
	for _, rAttr := range attrs {
		attr, ok := p.LookupAttributeByID(rAttr.AttrID)
		if !ok {
			return nil, fmt.Errorf("no attribute found by id")
		}

		switch permissions.AttributeTypes(rAttr.Type) {
		case permissions.StringListAttributeType:
			aKey := BuildGuardWithKey(attr.Category, attr.Name, Key(rAttr.Key))
			for _, v := range rAttr.Value.GetStringList().Strings {
				guard := helpers.Guard(aKey + "." + v)
				as = append(as, guard)
			}
		}
	}

	return as, nil
}

func (p *Perms) AddOrUpdateAttributesToRole(ctx context.Context, job string, grade int32, roleId uint64, attrs ...*permissions.RoleAttribute) error {
	for i := 0; i < len(attrs); i++ {
		if attrs[i].Value != nil {
			attrs[i].Value.Default(permissions.AttributeTypes(attrs[i].Type))

			a, ok := p.LookupAttributeByID(attrs[i].AttrId)
			if !ok {
				return fmt.Errorf("no attribute found by id %d", attrs[i].AttrId)
			}

			max := p.GetClosestRoleAttrMaxVals(job, grade, a.PermissionID, a.Key)
			if !attrs[i].Value.Check(a.Type, a.ValidValues, max) {
				return ErrAttrInvalid
			}
		}
	}

	if err := p.addOrUpdateAttributesToRole(ctx, roleId, attrs...); err != nil {
		return err
	}

	return nil
}

func (p *Perms) addOrUpdateAttributesToRole(ctx context.Context, roleId uint64, attrs ...*permissions.RoleAttribute) error {
	for i := 0; i < len(attrs); i++ {
		a, ok := p.LookupAttributeByID(attrs[i].AttrId)
		if !ok {
			return fmt.Errorf("unable to add role attribute, didn't find attribute by ID %d", attrs[i].AttrId)
		}

		validV := jet.NULL
		if attrs[i].Value == nil {
			attrs[i].Value = a.DefaultValues
		}

		if attrs[i].Value != nil {
			var out string
			var err error

			attrs[i].Value.Default(permissions.AttributeTypes(attrs[i].Type))

			switch permissions.AttributeTypes(attrs[i].Type) {
			case permissions.StringListAttributeType:
				out, err = json.MarshalToString(attrs[i].Value.GetStringList().Strings)
				if err != nil {
					return err
				}
			case permissions.JobListAttributeType:
				out, err = json.MarshalToString(attrs[i].Value.GetJobList().Strings)
				if err != nil {
					return err
				}
			case permissions.JobGradeListAttributeType:
				out, err = json.MarshalToString(attrs[i].Value.GetJobGradeList().Jobs)
				if err != nil {
					return err
				}
			}

			if out != "" && out != "null" {
				validV = jet.String(out)
			}
		}

		stmt := tRoleAttrs.
			INSERT(
				tRoleAttrs.RoleID,
				tRoleAttrs.AttrID,
				tRoleAttrs.Value,
			).
			VALUES(
				roleId,
				a.ID,
				validV,
			).
			ON_DUPLICATE_KEY_UPDATE(
				tRoleAttrs.Value.SET(jet.StringExp(jet.Raw("values(`value`)"))),
			)

		if _, err := stmt.ExecContext(ctx, p.db); err != nil {
			if err != nil && !dbutils.IsDuplicateError(err) {
				return err
			}
		}

		attr, ok := p.LookupAttributeByID(attrs[i].AttrId)
		if !ok {
			return fmt.Errorf("no attribute by id found")
		}

		p.updateRoleAttributeInMap(attrs[i].RoleId, attr.PermissionID, attr.ID, attr.Key, attr.Type, attrs[i].Value, attrs[i].MaxValues)
	}

	if err := p.publishMessage(RoleAttrUpdateSubject, RoleAttrUpdateEvent{
		RoleID: roleId,
	}); err != nil {
		return err
	}

	return nil
}

func (p *Perms) RemoveAttributesFromRole(ctx context.Context, roleId uint64, attrs ...*permissions.RoleAttribute) error {
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

	if _, err := stmt.ExecContext(ctx, p.db); err != nil {
		return err
	}

	for i := 0; i < len(attrs); i++ {
		p.removeRoleAttributeFromMap(roleId, attrs[i].AttrId)

		if err := p.publishMessage(RoleAttrUpdateSubject, RoleAttrUpdateEvent{
			RoleID: roleId,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (p *Perms) GetRoleAttributeByID(roleId uint64, attrId uint64) (*permissions.RoleAttribute, bool) {
	r, ok := p.lookupRoleAttribute(roleId, attrId)
	if !ok {
		return nil, false
	}

	return &permissions.RoleAttribute{
		RoleId:    roleId,
		AttrId:    r.AttrID,
		Key:       string(r.Key),
		Type:      string(r.Type),
		MaxValues: r.Max,
	}, true
}

func (p *Perms) UpdateRoleAttributeMaxValues(ctx context.Context, roleId uint64, attrId uint64, maxValues *permissions.AttributeValues) error {
	a, ok := p.LookupAttributeByID(attrId)
	if !ok {
		return fmt.Errorf("unable to update role attribute max values, didn't find attribute by ID %d", attrId)
	}

	ra, ok := p.lookupRoleAttribute(roleId, attrId)
	if !ok {
		return fmt.Errorf("unable to update role attribute max values, didn't find role attribute by ID %d and role ID %d", attrId, roleId)
	}

	maxVal := jet.NULL
	if maxValues != nil {
		var out string
		var err error

		maxValues.Default(permissions.AttributeTypes(a.Type))

		switch permissions.AttributeTypes(a.Type) {
		case permissions.StringListAttributeType:
			out, err = json.MarshalToString(maxValues.GetStringList().Strings)
			if err != nil {
				return err
			}
		case permissions.JobListAttributeType:
			out, err = json.MarshalToString(maxValues.GetJobList().Strings)
			if err != nil {
				return err
			}
		case permissions.JobGradeListAttributeType:
			out, err = json.MarshalToString(maxValues.GetJobGradeList().Jobs)
			if err != nil {
				return err
			}
		}

		if out != "" && out != "null" {
			maxVal = jet.String(out)
		}
	}

	stmt := tRoleAttrs.
		INSERT(
			tRoleAttrs.RoleID,
			tRoleAttrs.AttrID,
			tRoleAttrs.Value,
			tRoleAttrs.MaxValues,
		).
		VALUES(
			roleId,
			attrId,
			jet.NULL,
			maxVal,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tRoleAttrs.MaxValues.SET(jet.StringExp(jet.Raw("values(`max_values`)"))),
		)

	if _, err := stmt.ExecContext(ctx, p.db); err != nil {
		if err != nil && !dbutils.IsDuplicateError(err) {
			return err
		}
	}

	p.updateRoleAttributeInMap(roleId, a.PermissionID, attrId, a.Key, a.Type, ra.Value, maxValues)

	return nil
}
