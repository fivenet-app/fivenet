package perms

import (
	"context"
	"fmt"
	"slices"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/protoutils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/puzpuzpuz/xsync/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Key string

var ErrAttrInvalid = errors.New("invalid attributes")

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var (
	tAttrs     = table.FivenetRbacAttrs
	tRoleAttrs = table.FivenetRbacRolesAttrs
	tJobAttrs  = table.FivenetRbacJobAttrs
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

	tAttrs := table.FivenetRbacAttrs.AS("role_attribute")
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

	var dest []*permissions.RoleAttribute
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	for i := range dest {
		attr, ok := p.LookupAttributeByID(dest[i].AttrId)
		if !ok {
			return nil, fmt.Errorf("no attribute found by id")
		}

		dest[i].Category = string(attr.Category)
		dest[i].Name = string(attr.Name)
	}

	return dest, nil
}

func (p *Perms) getAttributeFromDatabase(ctx context.Context, permId uint64, key Key) (*model.FivenetRbacAttrs, error) {
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

	var dest model.FivenetRbacAttrs
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to query attribute from database. %w", err)
		}
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

func (p *Perms) GetJobAttributeValue(ctx context.Context, job string, attrId uint64) (*permissions.AttributeValues, error) {
	tJobAttrs := table.FivenetRbacJobAttrs.AS("role_attribute")
	stmt := tJobAttrs.
		SELECT(
			tJobAttrs.MaxValues.AS("value"),
		).
		FROM(
			tJobAttrs,
		).
		WHERE(jet.AND(
			tJobAttrs.Job.EQ(jet.String(job)),
			tJobAttrs.AttrID.EQ(jet.Uint64(attrId)),
		)).
		LIMIT(1)

	dest := &struct {
		Value *permissions.AttributeValues `alias:"value"`
	}{}
	if err := stmt.QueryContext(ctx, p.db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest.Value, nil
}

func (p *Perms) GetJobAttributes(ctx context.Context, job string) ([]*permissions.RoleAttribute, error) {
	tJobAttrs := table.FivenetRbacJobAttrs.AS("role_attribute")
	tAttrs := table.FivenetRbacAttrs
	stmt := tJobAttrs.
		SELECT(
			tJobAttrs.Job,
			tJobAttrs.AttrID,
			tAttrs.PermissionID.AS("role_attribute.permission_id"),
			tPerms.Category.AS("role_attribute.category"),
			tPerms.Name.AS("role_attribute.name"),
			tAttrs.Key.AS("role_attribute.key"),
			tAttrs.Type.AS("role_attribute.type"),
			tAttrs.ValidValues.AS("role_attribute.valid_values"),
			tJobAttrs.MaxValues,
		).
		FROM(
			tJobAttrs.
				INNER_JOIN(tAttrs,
					tAttrs.ID.EQ(tJobAttrs.AttrID),
				).
				INNER_JOIN(tPerms,
					tPerms.ID.EQ(tAttrs.PermissionID),
				),
		).
		WHERE(tJobAttrs.Job.EQ(jet.String(job)))

	dest := []*permissions.RoleAttribute{}
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}

func (p *Perms) Attr(userInfo *userinfo.UserInfo, category Category, name Name, key Key) (*permissions.AttributeValues, error) {
	permId, ok := p.lookupPermIDByGuard(BuildGuard(category, name))
	if !ok {
		return nil, nil
	}

	rAttr := p.getClosestRoleAttr(userInfo.Job, userInfo.JobGrade, permId, key)
	if userInfo.Superuser {
		attr, ok := p.lookupAttributeByPermID(permId, key)
		if !ok {
			return nil, nil
		}

		if attr.ValidValues != nil {
			rAttr = &cacheRoleAttr{
				Job:          userInfo.Job,
				AttrID:       attr.ID,
				PermissionID: attr.PermissionID,
				Key:          key,
				Type:         attr.Type,
				Value:        attr.ValidValues,
			}
		}
	}

	if rAttr == nil {
		return nil, nil
	}

	return proto.Clone(rAttr.Value).(*permissions.AttributeValues), nil
}

func (p *Perms) AttrStringList(userInfo *userinfo.UserInfo, category Category, name Name, key Key) (*permissions.StringList, error) {
	attrValue, err := p.Attr(userInfo, category, name, key)
	if err != nil {
		return &permissions.StringList{}, err
	}

	if attrValue == nil || attrValue.ValidValues == nil {
		return &permissions.StringList{}, nil
	}

	switch v := attrValue.ValidValues.(type) {
	case *permissions.AttributeValues_StringList:
		return v.StringList, nil

	case *permissions.AttributeValues_JobList:
		return v.JobList, nil

	default:
		return &permissions.StringList{}, fmt.Errorf("unknown role attribute type")
	}
}

func (p *Perms) AttrJobList(userInfo *userinfo.UserInfo, category Category, name Name, key Key) (*permissions.StringList, error) {
	return p.AttrStringList(userInfo, category, name, key)
}

func (p *Perms) AttrJobGradeList(userInfo *userinfo.UserInfo, category Category, name Name, key Key) (*permissions.JobGradeList, error) {
	attrValue, err := p.Attr(userInfo, category, name, key)
	if err != nil {
		return &permissions.JobGradeList{
			Jobs:        map[string]int32{},
			FineGrained: false,
			Grades:      map[string]*permissions.JobGrades{},
		}, err
	}

	if attrValue == nil || attrValue.ValidValues == nil {
		return &permissions.JobGradeList{
			Jobs:        map[string]int32{},
			FineGrained: false,
			Grades:      map[string]*permissions.JobGrades{},
		}, nil
	}

	switch v := attrValue.ValidValues.(type) {
	case *permissions.AttributeValues_JobGradeList:
		return v.JobGradeList, nil

	default:
		return &permissions.JobGradeList{
			Jobs:        map[string]int32{},
			FineGrained: false,
			Grades:      map[string]*permissions.JobGrades{},
		}, fmt.Errorf("unknown role attribute type for string list")
	}
}

func (p *Perms) convertRawValue(targetVal *permissions.AttributeValues, rawVal string, aType permissions.AttributeTypes) error {
	if err := protojson.Unmarshal([]byte(rawVal), targetVal); err != nil {
		return fmt.Errorf("failed to unmarshal raw value. %w", err)
	}

	targetVal.Default(aType)

	return nil
}

func (p *Perms) GetAllAttributes(ctx context.Context) ([]*permissions.RoleAttribute, error) {
	stmt := tAttrs.
		SELECT(
			tAttrs.ID.AS("role_attribute.attr_id"),
			tAttrs.PermissionID.AS("role_attribute.permission_id"),
			tPerms.Category.AS("role_attribute.category"),
			tPerms.Name.AS("role_attribute.name"),
			tAttrs.Key.AS("role_attribute.key"),
			tAttrs.Type.AS("role_attribute.type"),
			tAttrs.ValidValues.AS("role_attribute.valid_values"),
		).
		FROM(tAttrs.
			INNER_JOIN(tPerms,
				tPerms.ID.EQ(tAttrs.PermissionID),
			),
		)

	var dest []*permissions.RoleAttribute
	if err := stmt.QueryContext(ctx, p.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to query all attributes. %w", err)
		}
	}

	for i := range dest {
		if dest[i].Value == nil {
			dest[i].Value = &permissions.AttributeValues{}
			dest[i].Value.Default(permissions.AttributeTypes(dest[i].Type))
		}

		if dest[i].ValidValues == nil {
			dest[i].ValidValues = &permissions.AttributeValues{}
			dest[i].ValidValues.Default(permissions.AttributeTypes(dest[i].Type))
		}

		// MaxValues are not set because we don't know the job context here
	}

	return dest, nil
}

func (p *Perms) GetRoleAttributes(ctx context.Context, job string, grade int32) ([]*permissions.RoleAttribute, error) {
	roleId, ok := p.lookupRoleIDForJobAndGrade(job, grade)
	if !ok {
		roleId, ok = p.lookupRoleIDForJobAndGrade(DefaultRoleJob, p.startJobGrade)
		if !ok {
			return nil, fmt.Errorf("failed to fallback to default role. %w", fmt.Errorf("no role ID found for job and grade"))
		}
	}

	attrsRoleMap, ok := p.attrsRoleMap.Load(roleId)
	if !ok {
		return []*permissions.RoleAttribute{}, nil
	}

	var err error
	attrs := []*permissions.RoleAttribute{}
	attrsRoleMap.Range(func(key uint64, value *cacheRoleAttr) bool {
		attr, ok := p.LookupAttributeByID(key)
		if !ok {
			err = fmt.Errorf("no attribute found by id for role. %w", fmt.Errorf("attribute ID not found"))
			return false
		}

		attrVal, ok := attrsRoleMap.Load(attr.ID)
		if !ok {
			err = fmt.Errorf("no role attribute found by id for role. %w", fmt.Errorf("role attribute ID not found"))
			return false
		}

		maxValues, er := p.GetJobAttributeValue(ctx, job, attr.ID)
		if er != nil {
			err = fmt.Errorf("failed to retrieve max values for attribute %s/%s/%s. %w", attr.Category, attr.Name, attr.Key, er)
			return false
		}

		attrs = append(attrs, &permissions.RoleAttribute{
			RoleId:       roleId,
			AttrId:       attr.ID,
			PermissionId: attr.PermissionID,
			Category:     string(attr.Category),
			Name:         string(attr.Name),
			Key:          string(attr.Key),
			Type:         string(attr.Type),
			Value:        attrVal.Value,
			ValidValues:  attr.ValidValues,
			MaxValues:    maxValues,
		})

		return true
	})
	if err != nil {
		return nil, fmt.Errorf("failed to range over attributes. %w", err)
	}

	return attrs, nil
}

func (p *Perms) GetEffectiveRoleAttributes(ctx context.Context, job string, grade int32) ([]*permissions.RoleAttribute, error) {
	roleAttrs := map[uint64]interface{}{}

	roleIds, ok := p.lookupRoleIDsForJobUpToGrade(job, grade)
	if !ok {
		return nil, nil
	}

	perms := p.getRolePermissionsFromCache(roleIds)

	var err error
	attrs := []*permissions.RoleAttribute{}
	for i := range slices.Backward(roleIds) {
		attrMap, ok := p.attrsRoleMap.Load(roleIds[i])
		if !ok {
			continue
		}

		attrMap.Range(func(_ uint64, value *cacheRoleAttr) bool {
			// Skip already added attributes
			if _, ok := roleAttrs[value.AttrID]; ok {
				return true
			}

			// Permission not granted
			if !slices.ContainsFunc(perms, func(p *cachePerm) bool {
				return p.ID == value.PermissionID
			}) {
				return true
			}

			attr, ok := p.LookupAttributeByID(value.AttrID)
			if !ok {
				err = fmt.Errorf("no attribute found by id for role. %w", fmt.Errorf("attribute ID not found"))
				return false
			}

			attrMap, ok := p.attrsRoleMap.Load(roleIds[i])
			if !ok {
				return true
			}
			attrVal, ok := attrMap.Load(attr.ID)
			if !ok {
				err = fmt.Errorf("no role attribute found by id for role. %w", fmt.Errorf("role attribute ID not found"))
				return false
			}

			maxVal, _ := p.GetJobAttributeValue(ctx, job, attr.ID)

			attrs = append(attrs, &permissions.RoleAttribute{
				RoleId:       roleIds[i],
				AttrId:       value.AttrID,
				PermissionId: value.PermissionID,
				Category:     string(attr.Category),
				Name:         string(attr.Name),
				Key:          string(attr.Key),
				Type:         string(attr.Type),
				Value:        attrVal.Value,
				ValidValues:  attr.ValidValues,
				MaxValues:    maxVal,
			})

			roleAttrs[value.AttrID] = nil

			return true
		})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to range over attributes. %w", err)
	}

	return attrs, nil
}

func (p *Perms) getRoleAttributesFromCache(job string, grade int32) ([]*cacheRoleAttr, error) {
	roleAttrs := map[uint64]*cacheRoleAttr{}

	roleIds, ok := p.lookupRoleIDsForJobUpToGrade(job, grade)
	if !ok {
		return nil, nil
	}

	for i := range slices.Backward(roleIds) {
		attrMap, ok := p.attrsRoleMap.Load(roleIds[i])
		if !ok {
			continue
		}

		attrMap.Range(func(_ uint64, value *cacheRoleAttr) bool {
			// Skip already added attributes
			if _, ok := roleAttrs[value.AttrID]; ok {
				return true
			}

			roleAttrs[value.AttrID] = value

			return true
		})
	}

	as := make([]*cacheRoleAttr, 0, len(roleAttrs))
	for _, v := range roleAttrs {
		as = append(as, v)
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

func (p *Perms) UpdateRoleAttributes(ctx context.Context, job string, roleId uint64, attrs ...*permissions.RoleAttribute) error {
	for i := range attrs {
		attrs[i].RoleId = roleId

		a, ok := p.LookupAttributeByID(attrs[i].AttrId)
		if !ok {
			return fmt.Errorf("no attribute found by id %d. %w", attrs[i].AttrId, fmt.Errorf("attribute not found"))
		}

		if attrs[i].Value != nil {
			attrs[i].Value.Default(permissions.AttributeTypes(attrs[i].Type))

			maxValues, err := p.GetJobAttributeValue(ctx, job, a.ID)
			if err != nil {
				return fmt.Errorf("failed to retrieve max values for attribute %s/%s/%s. %w", a.Category, a.Name, a.Key, err)
			}

			valid, _ := attrs[i].Value.Check(a.Type, a.ValidValues, maxValues)
			if !valid {
				return errors.Wrapf(ErrAttrInvalid, "attribute %s/%s failed validation", a.Key, a.Name)
			}
		}
	}

	if err := p.addOrUpdateAttributesToRole(ctx, roleId, attrs...); err != nil {
		return fmt.Errorf("failed to add or update attributes to role. %w", err)
	}

	if err := p.publishMessage(ctx, RoleAttrUpdateSubject, &permissions.RoleIDEvent{
		RoleId: roleId,
	}); err != nil {
		return fmt.Errorf("failed to publish role attribute update message. %w", err)
	}

	return nil
}

func (p *Perms) addOrUpdateAttributesToRole(ctx context.Context, roleId uint64, attrs ...*permissions.RoleAttribute) error {
	for i := range attrs {
		a, ok := p.LookupAttributeByID(attrs[i].AttrId)
		if !ok {
			return fmt.Errorf("unable to add role attribute, didn't find attribute by ID %d. %w", attrs[i].AttrId, fmt.Errorf("attribute not found"))
		}

		if attrs[i].Value == nil {
			attrs[i].Value = &permissions.AttributeValues{}
		}

		if attrs[i].Value != nil {
			attrs[i].Value.Default(permissions.AttributeTypes(a.Type))
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
				attrs[i].Value,
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
	}

	if err := p.publishMessage(ctx, RoleAttrUpdateSubject, &permissions.RoleIDEvent{
		RoleId: roleId,
	}); err != nil {
		return fmt.Errorf("failed to publish role attribute removal message. %w", err)
	}

	return nil
}

func (p *Perms) RemoveAttributesFromRoleByPermission(ctx context.Context, roleId uint64, permissionId uint64) error {
	as, ok := p.attrsPermsMap.Load(permissionId)
	if !ok {
		return nil
	}

	ras := []*permissions.RoleAttribute{}
	as.Range(func(key string, attrId uint64) bool {
		ras = append(ras, &permissions.RoleAttribute{
			AttrId: attrId,
		})
		return true
	})

	if len(ras) == 0 {
		return nil
	}

	if err := p.RemoveAttributesFromRole(ctx, roleId, ras...); err != nil {
		return fmt.Errorf("failed to remove attributes from role by perm. %w", err)
	}

	if err := p.publishMessage(ctx, RoleAttrUpdateSubject, &permissions.RoleIDEvent{
		RoleId: roleId,
	}); err != nil {
		return fmt.Errorf("failed to publish role attribute removal message. %w", err)
	}

	return nil
}

func (p *Perms) UpdateJobAttributes(ctx context.Context, job string, attrs ...*permissions.RoleAttribute) error {
	for _, attr := range attrs {
		a, ok := p.LookupAttributeByID(attr.AttrId)
		if !ok {
			return fmt.Errorf("unable to update role attribute max values, didn't find attribute by ID %d. %w", attr.AttrId, fmt.Errorf("attribute not found"))
		}

		maxVal := jet.NULL
		if attr.MaxValues != nil {
			attr.MaxValues.Default(permissions.AttributeTypes(a.Type))

			out, err := protoutils.MarshalToPJSON(attr.MaxValues)
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
				attr.AttrId,
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

	return nil
}

func (p *Perms) updateRoleAttributeInMap(roleId uint64, permId uint64, attrId uint64, key Key, aType permissions.AttributeTypes, value *permissions.AttributeValues) {
	job, ok := p.lookupJobForRoleID(roleId)
	if !ok {
		p.logger.Error("unable to lookup job for role id", zap.Uint64("role_id", roleId))
		return
	}

	attrRoleMap, _ := p.attrsRoleMap.LoadOrCompute(roleId, func() (*xsync.Map[uint64, *cacheRoleAttr], bool) {
		return xsync.NewMap[uint64, *cacheRoleAttr](), false
	})

	attrRoleMap.Store(attrId, &cacheRoleAttr{
		Job:          job,
		AttrID:       attrId,
		PermissionID: permId,
		Key:          key,
		Type:         aType,
		Value:        value,
	})
}

func (p *Perms) removeRoleAttributeFromMap(roleId uint64, attrId uint64) {
	attrMap, ok := p.attrsRoleMap.Load(roleId)
	if !ok {
		return
	}

	attrMap.Delete(attrId)
}
