package perms

import (
	"context"
	"fmt"
	"slices"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/permissions"
	"github.com/fivenet-app/fivenet/pkg/dbutils"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/utils/protoutils"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/puzpuzpuz/xsync/v4"
	"google.golang.org/protobuf/encoding/protojson"
)

var ErrAttrInvalid = errors.New("invalid attributes")

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var (
	tAttrs     = table.FivenetAttrs
	tRoleAttrs = table.FivenetRoleAttrs
	tJobAttrs  = table.FivenetJobAttrs
)

func (p *Perms) GetAttribute(category Category, name Name, key Key) (*permissions.RoleAttribute, error) {
	permId, ok := p.lookupPermIDByGuard(BuildGuard(category, name))
	if !ok {
		return nil, fmt.Errorf("unable to find perm ID for attribute %s/%s/%s", category, name, key)
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
	for i := range attrIds {
		ids[i] = jet.Uint64(attrIds[i])
	}

	stmt := tAttrs.
		SELECT(
			tAttrs.ID,
			tAttrs.CreatedAt,
			tAttrs.PermissionID,
			tAttrs.Key,
			tAttrs.Type,
			tAttrs.ValidValues,
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
	for i := range dest {
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
			tAttrs.ID,
			tAttrs.CreatedAt,
			tAttrs.PermissionID,
			tAttrs.Key,
			tAttrs.Type,
			tAttrs.ValidValues,
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
		return nil, fmt.Errorf("failed to query attribute from database. %w", err)
	}

	return &dest, nil
}

func (p *Perms) CreateAttribute(ctx context.Context, permId uint64, key Key, aType permissions.AttributeTypes, validValues *permissions.AttributeValues) (uint64, error) {
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
			validValues,
		)

	res, err := stmt.ExecContext(ctx, p.db)
	if err != nil {
		if !dbutils.IsDuplicateError(err) {
			return 0, fmt.Errorf("failed to insert attribute into database. %w", err)
		}

		attr, err := p.getAttributeFromDatabase(ctx, permId, key)
		if err != nil {
			return 0, fmt.Errorf("failed to retrieve existing attribute after duplicate error. %w", err)
		}

		if err := p.addOrUpdateAttributeInMap(permId, uint64(attr.ID), key, aType, validValues); err != nil {
			return 0, fmt.Errorf("failed to add or update attribute in map after duplicate error. %w", err)
		}

		return attr.ID, nil
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID. %w", err)
	}

	if err := p.addOrUpdateAttributeInMap(permId, uint64(lastId), key, aType, validValues); err != nil {
		return 0, fmt.Errorf("failed to add or update attribute in map. %w", err)
	}

	return uint64(lastId), nil
}

func (p *Perms) addOrUpdateAttributeInMap(permId uint64, attrId uint64, key Key, aType permissions.AttributeTypes, validValues *permissions.AttributeValues) error {
	perm, ok := p.lookupPermByID(permId)
	if !ok {
		return fmt.Errorf("no permission found by id %d", permId)
	}

	attr := &cacheAttr{
		ID:           attrId,
		PermissionID: permId,
		Category:     perm.Category,
		Name:         perm.Name,
		Key:          key,
		Type:         aType,
		ValidValues:  validValues,
	}

	p.attrsMap.Store(attrId, attr)

	pAttrMap, _ := p.attrsPermsMap.LoadOrCompute(permId, func() (*xsync.Map[string, uint64], bool) {
		return xsync.NewMap[string, uint64](), false
	})
	pAttrMap.Store(string(key), attrId)

	return nil
}

func (p *Perms) UpdateAttribute(ctx context.Context, attrId uint64, permId uint64, key Key, aType permissions.AttributeTypes, validValues *permissions.AttributeValues) error {
	stmt := tAttrs.
		UPDATE(
			tAttrs.PermissionID,
			tAttrs.Key,
			tAttrs.Type,
			tAttrs.ValidValues,
		).
		SET(
			permId,
			string(key),
			string(aType),
			validValues,
		).
		WHERE(
			tAttrs.ID.EQ(jet.Uint64(attrId)),
		)

	if _, err := stmt.ExecContext(ctx, p.db); err != nil {
		return fmt.Errorf("failed to execute update statement. %w", err)
	}

	if err := p.addOrUpdateAttributeInMap(permId, attrId, key, aType, validValues); err != nil {
		return fmt.Errorf("failed to add or update attribute in map. %w", err)
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
	attrId, ok := pAttrs.Load(string(key))
	if !ok {
		return nil
	}

	for i := range slices.Backward(roleIds) {
		val, ok := p.lookupRoleAttribute(roleIds[i], attrId)
		if !ok {
			continue
		}

		return val
	}

	return nil
}

func (p *Perms) GetJobAttrMaxVals(job string, attrId uint64) (*permissions.AttributeValues, bool) {
	jas, ok := p.attrsJobMaxValuesMap.Load(job)
	if !ok {
		return nil, false
	}

	return jas.Load(attrId)
}

func (p *Perms) Attr(userInfo *userinfo.UserInfo, category Category, name Name, key Key) (any, error) {
	permId, ok := p.lookupPermIDByGuard(BuildGuard(category, name))
	if !ok {
		return nil, nil
	}

	cached := p.getClosestRoleAttr(userInfo.Job, userInfo.JobGrade, permId, key)
	if userInfo.SuperUser {
		attr, ok := p.lookupAttributeByPermID(permId, key)
		if !ok {
			return nil, nil
		}

		if attr.ValidValues != nil {
			cached = &cacheRoleAttr{
				Job:          userInfo.Job,
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

func (p *Perms) convertRawToRoleAttributes(in []*permissions.RawRoleAttribute, job string) []*permissions.RoleAttribute {
	res := make([]*permissions.RoleAttribute, len(in))
	for i := range in {
		res[i] = &permissions.RoleAttribute{
			RoleId:       in[i].RoleId,
			CreatedAt:    in[i].CreatedAt,
			AttrId:       in[i].AttrId,
			PermissionId: in[i].PermissionId,
			Category:     in[i].Category,
			Name:         in[i].Name,
			Key:          in[i].Key,
			Type:         in[i].Type,
			Value:        in[i].Value,
			ValidValues:  in[i].ValidValues,
			MaxValues:    &permissions.AttributeValues{},
		}

		if res[i].Value == nil {
			res[i].Value = &permissions.AttributeValues{}
			res[i].Value.Default(permissions.AttributeTypes(in[i].Type))
		}

		if res[i].ValidValues == nil {
			res[i].ValidValues = &permissions.AttributeValues{}
			res[i].ValidValues.Default(permissions.AttributeTypes(in[i].Type))
		}

		res[i].MaxValues, _ = p.GetJobAttrMaxVals(job, in[i].AttrId)
		if res[i].MaxValues != nil {
			res[i].MaxValues.Default(permissions.AttributeTypes(res[i].Type))
		}
	}

	return res
}

func (p *Perms) convertRawValue(targetVal *permissions.AttributeValues, rawVal string, aType permissions.AttributeTypes) error {
	if err := protojson.Unmarshal([]byte(rawVal), targetVal); err != nil {
		return fmt.Errorf("failed to unmarshal raw value. %w", err)
	}

	targetVal.Default(aType)

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
		).
		FROM(tAttrs.
			INNER_JOIN(tPerms,
				tPerms.ID.EQ(tAttrs.PermissionID),
			),
		)

	var dest []*permissions.RawRoleAttribute
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to query all attributes. %w", err)
		}
	}

	attrs := p.convertRawToRoleAttributes(dest, job)

	roleAttrs, err := p.GetRoleAttributes(job, grade)
	if err != nil {
		return nil, fmt.Errorf("failed to get role attributes. %w", err)
	}

	for _, attr := range attrs {
		idx := slices.IndexFunc(roleAttrs, func(ra *permissions.RoleAttribute) bool {
			return ra.AttrId == attr.AttrId
		})
		if idx > -1 {
			attr.Value = roleAttrs[idx].Value
		}
	}

	return attrs, nil
}

func (p *Perms) GetRoleAttributes(job string, grade int32) ([]*permissions.RoleAttribute, error) {
	roleId, ok := p.lookupRoleIDForJobAndGrade(job, grade)
	if !ok {
		roleId, ok = p.lookupRoleIDForJobAndGrade(DefaultRoleJob, p.startJobGrade)
		if !ok {
			return nil, fmt.Errorf("failed to fallback to default role. %w", fmt.Errorf("no role ID found for job and grade"))
		}
	}

	ars, ok := p.attrsRoleMap.Load(roleId)
	if !ok {
		return []*permissions.RoleAttribute{}, nil
	}

	var err error
	dest := []*permissions.RoleAttribute{}
	ars.Range(func(key uint64, value *cacheRoleAttr) bool {
		attr, ok := p.LookupAttributeByID(key)
		if !ok {
			err = fmt.Errorf("no attribute found by id for role. %w", fmt.Errorf("attribute ID not found"))
			return false
		}

		attrVal, ok := ars.Load(attr.ID)
		if !ok {
			err = fmt.Errorf("no role attribute found by id for role. %w", fmt.Errorf("role attribute ID not found"))
			return false
		}

		maxVal, _ := p.GetJobAttrMaxVals(job, attr.ID)

		dest = append(dest, &permissions.RoleAttribute{
			RoleId:       roleId,
			AttrId:       attr.ID,
			PermissionId: attr.PermissionID,
			Category:     string(attr.Category),
			Name:         string(attr.Name),
			Key:          string(attr.Key),
			Type:         string(attr.Type),
			Value:        attrVal.Value,
			ValidValues:  attr.ValidValues,
			MaxValues:    maxVal,
		})

		return true
	})
	if err != nil {
		return nil, fmt.Errorf("failed to range over attributes. %w", err)
	}

	return dest, nil
}

func (p *Perms) getRoleAttributesFromCache(job string, grade int32) ([]*cacheRoleAttr, error) {
	as := []*cacheRoleAttr{}

	roleIds, ok := p.lookupRoleIDsForJobUpToGrade(job, grade)
	if !ok {
		return as, nil
	}

	for i := range slices.Backward(roleIds) {
		attrMap, ok := p.attrsRoleMap.Load(roleIds[i])
		if !ok {
			continue
		}

		attrMap.Range(func(_ uint64, value *cacheRoleAttr) bool {
			as = append(as, value)

			return true
		})
	}

	return as, nil
}

func (p *Perms) FlattenRoleAttributes(job string, grade int32) ([]string, error) {
	attrs, err := p.getRoleAttributesFromCache(job, grade)
	if err != nil {
		return nil, fmt.Errorf("failed to get role attributes from cache. %w", err)
	}

	as := []string{}
	for _, rAttr := range attrs {
		attr, ok := p.LookupAttributeByID(rAttr.AttrID)
		if !ok {
			return nil, fmt.Errorf("no attribute found by id. %w", fmt.Errorf("attribute not found"))
		}

		switch permissions.AttributeTypes(rAttr.Type) {
		case permissions.StringListAttributeType:
			aKey := BuildGuardWithKey(attr.Category, attr.Name, Key(rAttr.Key))
			for _, v := range rAttr.Value.GetStringList().Strings {
				guard := Guard(aKey + "." + v)
				as = append(as, guard)
			}

		case permissions.JobListAttributeType:
			aKey := BuildGuardWithKey(attr.Category, attr.Name, Key(rAttr.Key))
			for _, v := range rAttr.Value.GetJobList().Strings {
				guard := Guard(aKey + "." + v)
				as = append(as, guard)
			}

		case permissions.JobGradeListAttributeType:
			// Only generate jobs as attribute
			aKey := BuildGuardWithKey(attr.Category, attr.Name, Key(rAttr.Key))
			for v := range rAttr.Value.GetJobGradeList().GetJobs() {
				guard := Guard(aKey + "." + v)
				as = append(as, guard)
			}
		}
	}

	return as, nil
}

func (p *Perms) AddOrUpdateAttributesToRole(ctx context.Context, job string, roleId uint64, attrs ...*permissions.RoleAttribute) error {
	for i := range attrs {
		attrs[i].RoleId = roleId

		a, ok := p.LookupAttributeByID(attrs[i].AttrId)
		if !ok {
			return fmt.Errorf("no attribute found by id %d. %w", attrs[i].AttrId, fmt.Errorf("attribute not found"))
		}

		if attrs[i].Value != nil {
			attrs[i].Value.Default(permissions.AttributeTypes(attrs[i].Type))

			max, _ := p.GetJobAttrMaxVals(job, a.ID)

			valid, _ := attrs[i].Value.Check(a.Type, a.ValidValues, max)
			if !valid {
				return errors.Wrapf(ErrAttrInvalid, "attribute %s/%s failed validation", a.Key, a.Name)
			}
		}
	}

	if err := p.addOrUpdateAttributesToRole(ctx, roleId, attrs...); err != nil {
		return fmt.Errorf("failed to add or update attributes to role. %w", err)
	}

	return nil
}

func (p *Perms) addOrUpdateAttributesToRole(ctx context.Context, roleId uint64, attrs ...*permissions.RoleAttribute) error {
	for i := range attrs {
		a, ok := p.LookupAttributeByID(attrs[i].AttrId)
		if !ok {
			return fmt.Errorf("unable to add role attribute, didn't find attribute by ID %d. %w", attrs[i].AttrId, fmt.Errorf("attribute not found"))
		}

		attrValue := jet.NULL
		if attrs[i].Value == nil {
			attrs[i].Value = &permissions.AttributeValues{}
		}

		if attrs[i].Value != nil {
			attrs[i].Value.Default(permissions.AttributeTypes(a.Type))

			out, err := protoutils.Marshal(attrs[i].Value)
			if err != nil {
				return fmt.Errorf("failed to marshal attribute value. %w", err)
			}

			attrValue = jet.String(string(out))
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
				attrValue,
			).
			ON_DUPLICATE_KEY_UPDATE(
				tRoleAttrs.Value.SET(jet.StringExp(jet.Raw("VALUES(`value`)"))),
			)

		if _, err := stmt.ExecContext(ctx, p.db); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return fmt.Errorf("failed to execute insert statement for role attributes. %w", err)
			}
		}

		p.updateRoleAttributeInMap(roleId, a.PermissionID, a.ID, a.Key, a.Type, attrs[i].Value)
	}

	if err := p.publishMessage(ctx, RoleAttrUpdateSubject, RoleIDEvent{
		RoleID: roleId,
	}); err != nil {
		return fmt.Errorf("failed to publish role attribute update message. %w", err)
	}

	return nil
}

func (p *Perms) RemoveAttributesFromRole(ctx context.Context, roleId uint64, attrs ...*permissions.RoleAttribute) error {
	if len(attrs) == 0 {
		return nil
	}

	ids := make([]jet.Expression, len(attrs))
	for i := range attrs {
		ids[i] = jet.Uint64(attrs[i].AttrId)
	}

	stmt := tRoleAttrs.
		DELETE().
		WHERE(jet.AND(
			tRoleAttrs.RoleID.EQ(jet.Uint64(roleId)),
			tRoleAttrs.AttrID.IN(ids...),
		))

	if _, err := stmt.ExecContext(ctx, p.db); err != nil {
		return fmt.Errorf("failed to execute delete statement for role attributes. %w", err)
	}

	for i := range attrs {
		p.removeRoleAttributeFromMap(roleId, attrs[i].AttrId)

		if err := p.publishMessage(ctx, RoleAttrUpdateSubject, RoleIDEvent{
			RoleID: roleId,
		}); err != nil {
			return fmt.Errorf("failed to publish role attribute removal message. %w", err)
		}
	}

	return nil
}

func (p *Perms) GetRoleAttributeByID(roleId uint64, attrId uint64) (*permissions.RoleAttribute, bool) {
	ra, ok := p.lookupRoleAttribute(roleId, attrId)
	if !ok {
		return nil, false
	}

	maxVals, _ := p.GetJobAttrMaxVals(ra.Job, ra.AttrID)

	return &permissions.RoleAttribute{
		RoleId:    roleId,
		AttrId:    attrId,
		Key:       string(ra.Key),
		Type:      string(ra.Type),
		MaxValues: maxVals,
	}, true
}

func (p *Perms) UpdateJobAttributeMaxValues(ctx context.Context, job string, attrId uint64, maxValues *permissions.AttributeValues) error {
	a, ok := p.LookupAttributeByID(attrId)
	if !ok {
		return fmt.Errorf("unable to update role attribute max values, didn't find attribute by ID %d. %w", attrId, fmt.Errorf("attribute not found"))
	}

	maxVal := jet.NULL
	if maxValues != nil {
		maxValues.Default(permissions.AttributeTypes(a.Type))

		out, err := protoutils.Marshal(maxValues)
		if err != nil {
			return fmt.Errorf("failed to marshal max values. %w", err)
		}

		maxVal = jet.String(string(out))
	}

	stmt := tJobAttrs.
		INSERT(
			tJobAttrs.Job,
			tJobAttrs.AttrID,
			tJobAttrs.MaxValues,
		).
		VALUES(
			job,
			attrId,
			maxVal,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tJobAttrs.MaxValues.SET(jet.StringExp(jet.Raw("VALUES(`max_values`)"))),
		)

	if _, err := stmt.ExecContext(ctx, p.db); err != nil {
		if !dbutils.IsDuplicateError(err) {
			return fmt.Errorf("failed to execute insert statement for job attributes. %w", err)
		}
	}

	if err := p.publishMessage(ctx, JobAttrUpdateSubject, JobAttrUpdateEvent{
		Job: job,
	}); err != nil {
		return fmt.Errorf("failed to publish job attribute update message. %w", err)
	}

	return nil
}

func (p *Perms) ClearJobAttributes(ctx context.Context, job string) error {
	stmt := tJobAttrs.
		DELETE().
		WHERE(tJobAttrs.Job.EQ(jet.String(job)))

	if _, err := stmt.ExecContext(ctx, p.db); err != nil {
		return fmt.Errorf("failed to execute delete statement for job attributes. %w", err)
	}

	if err := p.publishMessage(ctx, JobAttrUpdateSubject, JobAttrUpdateEvent{
		Job: job,
	}); err != nil {
		return fmt.Errorf("failed to publish job attribute clear message. %w", err)
	}

	return nil
}
