package perms

import (
	"context"
	"fmt"
	"slices"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	permissionsevents "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/events"
	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	propertyregistry "github.com/fivenet-app/fivenet/v2026/pkg/properties"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/pkg/errors"
	"github.com/puzpuzpuz/xsync/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Key string

var ErrAttrInvalid = errors.New("invalid attributes")

var (
	tAttrs     = table.FivenetRbacAttrs
	tRoleAttrs = table.FivenetRbacRolesAttrs
	tJobAttrs  = table.FivenetRbacJobAttrs
)

func (ps *Perms) GetAttribute(
	namespace Namespace,
	service Service,
	name Name,
	key Key,
) (*permissionsattributes.RoleAttribute, error) {
	permId, ok := ps.lookupPermIDByGuard(BuildGuard(namespace, service, name))
	if !ok {
		return nil, fmt.Errorf(
			"unable to find perm ID for attribute %s.%s/%s/%s",
			namespace,
			service,
			name,
			key,
		)
	}

	attr, ok := ps.lookupAttributeByPermID(permId, key)
	if !ok {
		return nil, errors.New("no attribute found by id")
	}

	return &permissionsattributes.RoleAttribute{
		AttrId:       attr.ID,
		PermissionId: attr.PermissionID,
		Namespace:    string(namespace),
		Service:      string(service),
		Name:         string(name),
		Key:          string(attr.Key),
		Type:         string(attr.Type),
		ValidValues:  attr.ValidValues,
		MaxValues:    nil,
	}, nil
}

func (ps *Perms) GetAttributeByIDs(
	ctx context.Context,
	attrIds ...int64,
) ([]*permissionsattributes.RoleAttribute, error) {
	ids := make([]mysql.Expression, len(attrIds))
	for i := range attrIds {
		ids[i] = mysql.Int64(attrIds[i])
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
		WHERE(mysql.AND(
			tAttrs.ID.IN(ids...),
		)).
		LIMIT(1)

	var dest []*permissionsattributes.RoleAttribute
	if err := stmt.QueryContext(ctx, ps.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	for i := range dest {
		attr, ok := ps.lookupAttributeByID(dest[i].GetAttrId())
		if !ok {
			return nil, errors.New("no attribute found by id")
		}

		dest[i].Namespace = string(attr.Namespace)
		dest[i].Service = string(attr.Service)
		dest[i].Name = string(attr.Name)
	}
	if err := ps.applyDynamicPropertyListOptions(ctx, dest); err != nil {
		return nil, err
	}

	return dest, nil
}

func (ps *Perms) getAttributeFromDatabase(
	ctx context.Context,
	permId int64,
	key Key,
) (*model.FivenetRbacAttrs, error) {
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
		WHERE(mysql.AND(
			tAttrs.PermissionID.EQ(mysql.Int64(permId)),
			tAttrs.Key.EQ(mysql.String(string(key))),
		)).
		LIMIT(1)

	var dest model.FivenetRbacAttrs
	if err := stmt.QueryContext(ctx, ps.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to query attribute from database. %w", err)
		}
	}

	return &dest, nil
}

func (ps *Perms) createAttribute(
	ctx context.Context,
	permId int64,
	key Key,
	aType permissionsattributes.AttributeTypes,
	validValues *permissionsattributes.AttributeValues,
	propertyScope propertypb.PropertyScope,
) (int64, error) {
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

	res, err := stmt.ExecContext(ctx, ps.db)
	if err != nil {
		if !dbutils.IsDuplicateError(err) {
			return 0, fmt.Errorf("failed to insert attribute into database. %w", err)
		}

		attr, err := ps.getAttributeFromDatabase(ctx, permId, key)
		if err != nil {
			return 0, fmt.Errorf(
				"failed to retrieve existing attribute after duplicate error. %w",
				err,
			)
		}

		if err := ps.addOrUpdateAttributeInMap(
			permId,
			attr.ID,
			key,
			aType,
			validValues,
			propertyScope,
		); err != nil {
			return 0, fmt.Errorf(
				"failed to add or update attribute in map after duplicate error. %w",
				err,
			)
		}

		return attr.ID, nil
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID. %w", err)
	}

	if err := ps.addOrUpdateAttributeInMap(
		permId,
		lastId,
		key,
		aType,
		validValues,
		propertyScope,
	); err != nil {
		return 0, fmt.Errorf("failed to add or update attribute in map. %w", err)
	}

	return lastId, nil
}

func (ps *Perms) addOrUpdateAttributeInMap(
	permId int64,
	attrId int64,
	key Key,
	aType permissionsattributes.AttributeTypes,
	validValues *permissionsattributes.AttributeValues,
	propertyScope propertypb.PropertyScope,
) error {
	perm, ok := ps.lookupPermByID(permId)
	if !ok {
		return fmt.Errorf("no permission found by id %d", permId)
	}

	attr := &cacheAttr{
		ID:            attrId,
		PermissionID:  permId,
		Namespace:     perm.Namespace,
		Service:       perm.Service,
		Name:          perm.Name,
		Key:           key,
		Type:          aType,
		PropertyScope: propertyScope,
		ValidValues:   validValues,
	}

	ps.attrsMap.Store(attrId, attr)

	pAttrMap, _ := ps.attrsPermsMap.LoadOrCompute(permId, func() (*xsync.Map[string, int64], bool) {
		return xsync.NewMap[string, int64](), false
	})
	pAttrMap.Store(string(key), attrId)

	return nil
}

func (ps *Perms) updateAttribute(
	ctx context.Context,
	attrId int64,
	permId int64,
	key Key,
	aType permissionsattributes.AttributeTypes,
	validValues *permissionsattributes.AttributeValues,
	propertyScope propertypb.PropertyScope,
) error {
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
			tAttrs.ID.EQ(mysql.Int64(attrId)),
		).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, ps.db); err != nil {
		return fmt.Errorf("failed to execute update statement. %w", err)
	}

	if err := ps.addOrUpdateAttributeInMap(
		permId,
		attrId,
		key,
		aType,
		validValues,
		propertyScope,
	); err != nil {
		return fmt.Errorf("failed to add or update attribute in map. %w", err)
	}

	return nil
}

func (ps *Perms) getClosestRoleAttr(job string, grade int32, permId int64, key Key) *cacheRoleAttr {
	roleIds, ok := ps.lookupRoleIDsForJobUpToGrade(job, grade)
	if !ok {
		return nil
	}

	pAttrs, ok := ps.attrsPermsMap.Load(permId)
	if !ok {
		return nil
	}
	attrId, ok := pAttrs.Load(string(key))
	if !ok {
		return nil
	}

	for i := range slices.Backward(roleIds) {
		val, ok := ps.lookupRoleAttribute(roleIds[i], attrId)
		if !ok {
			continue
		}

		return val
	}

	return nil
}

func (ps *Perms) GetJobAttributeValue(
	ctx context.Context,
	job string,
	attrId int64,
) (*permissionsattributes.AttributeValues, error) {
	tJobAttrs := table.FivenetRbacJobAttrs.AS("role_attribute")
	stmt := tJobAttrs.
		SELECT(
			tJobAttrs.MaxValues.AS("value"),
		).
		FROM(
			tJobAttrs,
		).
		WHERE(mysql.AND(
			tJobAttrs.Job.EQ(mysql.String(job)),
			tJobAttrs.AttrID.EQ(mysql.Int64(attrId)),
		)).
		LIMIT(1)

	dest := &struct {
		Value *permissionsattributes.AttributeValues `alias:"value"`
	}{}
	if err := stmt.QueryContext(ctx, ps.db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest.Value, nil
}

func (ps *Perms) GetJobAttributes(
	ctx context.Context,
	job string,
) ([]*permissionsattributes.RoleAttribute, error) {
	tJobAttrs := table.FivenetRbacJobAttrs.AS("role_attribute")
	tAttrs := table.FivenetRbacAttrs
	stmt := tJobAttrs.
		SELECT(
			tJobAttrs.Job,
			tJobAttrs.AttrID,
			tAttrs.PermissionID.AS("role_attribute.permission_id"),
			tPerms.Namespace.AS("role_attribute.namespace"),
			tPerms.Service.AS("role_attribute.service"),
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
		WHERE(tJobAttrs.Job.EQ(mysql.String(job)))

	dest := []*permissionsattributes.RoleAttribute{}
	if err := stmt.QueryContext(ctx, ps.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}
	if err := ps.applyDynamicPropertyListOptions(ctx, dest); err != nil {
		return nil, err
	}

	return dest, nil
}

// dynamicPropertyListValidValues overlays runtime custom property keys onto
// the generated system-property options for one property-list attribute.
func (ps *Perms) dynamicPropertyListValidValues(
	ctx context.Context,
	scope propertypb.PropertyScope,
	base *permissionsattributes.AttributeValues,
) (*permissionsattributes.AttributeValues, error) {
	scopeName, ok := propertyregistry.ScopeName(scope)
	if !ok {
		return base, nil
	}

	values := &permissionsattributes.AttributeValues{}
	if base != nil {
		values = proto.Clone(base).(*permissionsattributes.AttributeValues)
	}
	keys := append([]string(nil), values.GetPropertyList().GetKeys()...)

	tDefinitions := table.FivenetPropertyDefinitions
	stmt := tDefinitions.SELECT(tDefinitions.PropertyKey, tDefinitions.Label, tDefinitions.Implicit).
		FROM(tDefinitions).
		WHERE(
			tDefinitions.Scope.EQ(mysql.String(scopeName)).AND(tDefinitions.DeletedAt.IS_NULL()),
		)
	// Keep disabled keys selectable so existing grants remain editable; resource
	// handlers independently reject disabled definitions for reads and writes.
	var customDefinitions []model.FivenetPropertyDefinitions
	if err := stmt.QueryContext(ctx, ps.db, &customDefinitions); err != nil {
		return nil, fmt.Errorf("failed to load custom property-list options. %w", err)
	}
	customLabels := make(map[string]string, len(customDefinitions))
	for _, definition := range customDefinitions {
		if propertyregistry.IsCustomPropertyKey(scope, definition.PropertyKey) &&
			!definition.Implicit {
			keys = append(keys, definition.PropertyKey)
			customLabels[definition.PropertyKey] = definition.Label
		}
	}
	keys = slices.DeleteFunc(keys, func(key string) bool {
		definition, ok := propertyregistry.SystemPropertyDefinition(scope, key)
		return ok && definition.GetImplicit()
	})
	slices.Sort(keys)
	keys = slices.Compact(keys)
	propertyDefinitions := make([]*permissionsattributes.PropertyListEntry, 0, len(keys))
	for _, key := range keys {
		label := key
		if definition, ok := propertyregistry.SystemPropertyDefinition(scope, key); ok {
			label = definition.GetLabel()
		} else if customLabel, ok := customLabels[key]; ok {
			label = customLabel
		}
		propertyDefinitions = append(propertyDefinitions, &permissionsattributes.PropertyListEntry{
			Key: key, Label: label,
		})
	}
	values.ValidValues = &permissionsattributes.AttributeValues_PropertyList{
		PropertyList: &permissionsattributes.PropertyList{
			Keys:       keys,
			Properties: propertyDefinitions,
		},
	}

	return values, nil
}

func (ps *Perms) applyDynamicPropertyListOptions(
	ctx context.Context,
	attrs []*permissionsattributes.RoleAttribute,
) error {
	resolved := make(map[propertypb.PropertyScope]*permissionsattributes.AttributeValues)
	for _, roleAttr := range attrs {
		if roleAttr == nil {
			continue
		}
		attr, ok := ps.lookupAttributeByID(roleAttr.GetAttrId())
		if !ok || attr.Type != permissionsattributes.PropertyListAttributeType ||
			attr.PropertyScope == propertypb.PropertyScope_PROPERTY_SCOPE_UNSPECIFIED {
			continue
		}
		validValues, ok := resolved[attr.PropertyScope]
		if !ok {
			var err error
			validValues, err = ps.dynamicPropertyListValidValues(
				ctx,
				attr.PropertyScope,
				attr.ValidValues,
			)
			if err != nil {
				return err
			}
			resolved[attr.PropertyScope] = validValues
		}
		roleAttr.ValidValues = validValues
		roleAttr.MaxValues = propertyListWithDefinitions(roleAttr.MaxValues, validValues)
		roleAttr.Value = filterPropertyListValue(roleAttr.Value, validValues)
	}

	return nil
}

// propertyListWithDefinitions attaches the available definition metadata to a
// property-list value without changing the keys that constrain it.
func propertyListWithDefinitions(
	value *permissionsattributes.AttributeValues,
	validValues *permissionsattributes.AttributeValues,
) *permissionsattributes.AttributeValues {
	if value == nil || value.GetPropertyList() == nil || validValues == nil {
		return value
	}

	definitionsByKey := make(
		map[string]*permissionsattributes.PropertyListEntry,
		len(validValues.GetPropertyList().GetProperties()),
	)
	for _, definition := range validValues.GetPropertyList().GetProperties() {
		if definition != nil {
			definitionsByKey[definition.GetKey()] = definition
		}
	}

	result := proto.Clone(value).(*permissionsattributes.AttributeValues)
	definitions := make(
		[]*permissionsattributes.PropertyListEntry,
		0,
		len(value.GetPropertyList().GetKeys()),
	)
	for _, key := range value.GetPropertyList().GetKeys() {
		if definition, ok := definitionsByKey[key]; ok {
			definitions = append(definitions, definition)
		}
	}
	result.ValidValues = &permissionsattributes.AttributeValues_PropertyList{
		PropertyList: &permissionsattributes.PropertyList{
			Keys:       append([]string(nil), value.GetPropertyList().GetKeys()...),
			Properties: definitions,
		},
	}
	return result
}

func filterPropertyListValue(
	value *permissionsattributes.AttributeValues,
	validValues *permissionsattributes.AttributeValues,
) *permissionsattributes.AttributeValues {
	if value == nil || value.GetPropertyList() == nil || validValues == nil {
		return value
	}

	allowed := make(map[string]struct{}, len(validValues.GetPropertyList().GetKeys()))
	for _, key := range validValues.GetPropertyList().GetKeys() {
		allowed[key] = struct{}{}
	}

	filtered := proto.Clone(value).(*permissionsattributes.AttributeValues)
	keys := make([]string, 0, len(value.GetPropertyList().GetKeys()))
	for _, key := range value.GetPropertyList().GetKeys() {
		if _, ok := allowed[key]; ok {
			keys = append(keys, key)
		}
	}
	filtered.ValidValues = &permissionsattributes.AttributeValues_PropertyList{
		PropertyList: &permissionsattributes.PropertyList{Keys: keys},
	}
	return filtered
}

func (ps *Perms) Attr(
	userInfo *userinfo.UserInfo,
	namespace Namespace,
	service Service,
	name Name,
	key Key,
) (*permissionsattributes.AttributeValues, error) {
	permId, ok := ps.lookupPermIDByGuard(BuildGuard(namespace, service, name))
	if !ok {
		return nil, nil
	}

	rAttr := ps.getClosestRoleAttr(userInfo.GetJob(), userInfo.GetJobGrade(), permId, key)
	if userInfo.GetJobAdmin() {
		attr, ok := ps.lookupAttributeByPermID(permId, key)
		if !ok {
			return nil, nil
		}

		if attr.ValidValues != nil {
			rAttr = &cacheRoleAttr{
				Job:          userInfo.GetJob(),
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

	//nolint:forcetypeassert,errcheck // We know that rAttr.Value is of type permissionsattributes.AttributeValues as we just clone it.
	return proto.Clone(rAttr.Value).(*permissionsattributes.AttributeValues), nil
}

func (ps *Perms) attrStringListRaw(
	userInfo *userinfo.UserInfo,
	namespace Namespace,
	service Service,
	name Name,
	key Key,
) (*permissionsattributes.StringList, error) {
	attrValue, err := ps.Attr(userInfo, namespace, service, name, key)
	if err != nil {
		return &permissionsattributes.StringList{}, err
	}

	if attrValue == nil || attrValue.ValidValues == nil {
		return &permissionsattributes.StringList{}, nil
	}

	switch v := attrValue.GetValidValues().(type) {
	case *permissionsattributes.AttributeValues_StringList:
		return v.StringList, nil

	case *permissionsattributes.AttributeValues_JobList:
		return v.JobList, nil

	default:
		return &permissionsattributes.StringList{}, errors.New("unknown role attribute type")
	}
}

func (ps *Perms) attrJobListRaw(
	userInfo *userinfo.UserInfo,
	namespace Namespace,
	service Service,
	name Name,
	key Key,
) (*permissionsattributes.StringList, error) {
	return ps.attrStringListRaw(userInfo, namespace, service, name, key)
}

func (ps *Perms) attrJobGradeListRaw(
	userInfo *userinfo.UserInfo,
	namespace Namespace,
	service Service,
	name Name,
	key Key,
) (*permissionsattributes.JobGradeList, error) {
	attrValue, err := ps.Attr(userInfo, namespace, service, name, key)
	if err != nil {
		return &permissionsattributes.JobGradeList{
			Jobs:        map[string]int32{},
			FineGrained: false,
			Grades:      map[string]*permissionsattributes.JobGrades{},
		}, err
	}

	if attrValue == nil || attrValue.ValidValues == nil {
		return &permissionsattributes.JobGradeList{
			Jobs:        map[string]int32{},
			FineGrained: false,
			Grades:      map[string]*permissionsattributes.JobGrades{},
		}, nil
	}

	switch v := attrValue.GetValidValues().(type) {
	case *permissionsattributes.AttributeValues_JobGradeList:
		return v.JobGradeList, nil

	default:
		return &permissionsattributes.JobGradeList{
			Jobs:        map[string]int32{},
			FineGrained: false,
			Grades:      map[string]*permissionsattributes.JobGrades{},
		}, errors.New("unknown role attribute type for string list")
	}
}

func (ps *Perms) attrPropertyListRaw(
	userInfo *userinfo.UserInfo,
	namespace Namespace,
	service Service,
	name Name,
	key Key,
) (*permissionsattributes.PropertyList, error) {
	attrValue, err := ps.Attr(userInfo, namespace, service, name, key)
	if err != nil {
		return &permissionsattributes.PropertyList{}, err
	}

	if attrValue == nil || attrValue.ValidValues == nil {
		return &permissionsattributes.PropertyList{}, nil
	}

	switch v := attrValue.GetValidValues().(type) {
	case *permissionsattributes.AttributeValues_PropertyList:
		return proto.Clone(v.PropertyList).(*permissionsattributes.PropertyList), nil

	default:
		return &permissionsattributes.PropertyList{}, errors.New(
			"unknown role attribute type for property list",
		)
	}
}

// AttrPropertyListContext resolves a property-list attribute and appends all
// enabled implicit definitions for its property scope. The context is used
// for the database lookup of custom definitions.
func (ps *Perms) AttrPropertyListContext(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	attr AttrRef[PropertyListAttr],
) (*permissionsattributes.PropertyList, error) {
	result, err := ps.attrPropertyListRaw(
		userInfo,
		attr.perm.namespace,
		attr.perm.service,
		attr.perm.name,
		attr.key,
	)
	if err != nil {
		return result, err
	}
	permID, ok := ps.lookupPermIDByGuard(
		BuildGuard(attr.perm.namespace, attr.perm.service, attr.perm.name),
	)
	if !ok {
		return result, nil
	}
	attribute, ok := ps.lookupAttributeByPermID(permID, attr.key)
	if !ok || attribute.PropertyScope == propertypb.PropertyScope_PROPERTY_SCOPE_UNSPECIFIED {
		return result, nil
	}
	keys, err := ps.implicitPropertyKeys(ctx, attribute.PropertyScope)
	if err != nil {
		return nil, err
	}
	for _, propertyKey := range keys {
		if !slices.Contains(result.GetKeys(), propertyKey) {
			result.Keys = append(result.Keys, propertyKey)
		}
	}
	return result, nil
}

func (ps *Perms) implicitPropertyKeys(
	ctx context.Context,
	scope propertypb.PropertyScope,
) ([]string, error) {
	keys := make([]string, 0)
	enabledOverrides := map[string]bool{}
	name, ok := propertyregistry.ScopeName(scope)
	if !ok {
		return keys, nil
	}
	overrides := []model.FivenetPropertySystemOverrides{}
	tOverrides := table.FivenetPropertySystemOverrides
	if err := tOverrides.SELECT(tOverrides.PropertyKey, tOverrides.Enabled).FROM(tOverrides).
		WHERE(tOverrides.Scope.EQ(mysql.String(name))).
		QueryContext(ctx, ps.db, &overrides); err != nil {
		return nil, err
	}
	for _, override := range overrides {
		enabledOverrides[override.PropertyKey] = override.Enabled
	}
	for _, definition := range propertyregistry.SystemPropertyDefinitionsForScope(scope) {
		enabled := definition.GetEnabled()
		if override, ok := enabledOverrides[definition.GetKey()]; ok {
			enabled = override
		}
		if enabled && definition.GetImplicit() {
			keys = append(keys, definition.GetKey())
		}
	}
	rows := []model.FivenetPropertyDefinitions{}
	tDefinitions := table.FivenetPropertyDefinitions
	if err := tDefinitions.SELECT(tDefinitions.PropertyKey).FROM(tDefinitions).WHERE(
		tDefinitions.Scope.EQ(mysql.String(name)).AND(
			tDefinitions.Implicit.IS_TRUE().
				AND(tDefinitions.Enabled.IS_TRUE()).
				AND(tDefinitions.DeletedAt.IS_NULL()),
		),
	).QueryContext(ctx, ps.db, &rows); err != nil {
		return nil, err
	}
	for _, row := range rows {
		if propertyregistry.IsCustomPropertyKey(scope, row.PropertyKey) {
			keys = append(keys, row.PropertyKey)
		}
	}
	slices.Sort(keys)
	return slices.Compact(keys), nil
}

func (ps *Perms) convertRawValue(
	targetVal *permissionsattributes.AttributeValues,
	rawVal string,
	aType permissionsattributes.AttributeTypes,
) error {
	if err := protojson.Unmarshal([]byte(rawVal), targetVal); err != nil {
		return fmt.Errorf("failed to unmarshal raw value. %w", err)
	}

	targetVal.Default(aType)

	return nil
}

func (ps *Perms) GetAllAttributes(
	ctx context.Context,
) ([]*permissionsattributes.RoleAttribute, error) {
	stmt := tAttrs.
		SELECT(
			tAttrs.ID.AS("role_attribute.attr_id"),
			tAttrs.PermissionID.AS("role_attribute.permission_id"),
			tPerms.Namespace.AS("role_attribute.namespace"),
			tPerms.Service.AS("role_attribute.service"),
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

	var dest []*permissionsattributes.RoleAttribute
	if err := stmt.QueryContext(ctx, ps.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to query all attributes. %w", err)
		}
	}

	for i := range dest {
		if dest[i].GetValue() == nil {
			dest[i].Value = &permissionsattributes.AttributeValues{}
			dest[i].GetValue().Default(permissionsattributes.AttributeTypes(dest[i].GetType()))
		}

		if dest[i].GetValidValues() == nil {
			dest[i].ValidValues = &permissionsattributes.AttributeValues{}
			dest[i].GetValidValues().
				Default(permissionsattributes.AttributeTypes(dest[i].GetType()))
		}

		// MaxValues are not set because we don't know the job context here
	}
	if err := ps.applyDynamicPropertyListOptions(ctx, dest); err != nil {
		return nil, err
	}

	return dest, nil
}

// GetSuperuserAttributes returns all registered attributes with their complete
// set of valid values as the effective values.
func (ps *Perms) GetSuperuserAttributes(
	ctx context.Context,
) ([]*permissionsattributes.RoleAttribute, error) {
	attrs, err := ps.GetAllAttributes(ctx)
	if err != nil {
		return nil, err
	}

	return superuserAttributes(attrs), nil
}

func superuserAttributes(
	attrs []*permissionsattributes.RoleAttribute,
) []*permissionsattributes.RoleAttribute {
	result := make([]*permissionsattributes.RoleAttribute, 0, len(attrs))
	for _, attr := range attrs {
		if attr == nil {
			continue
		}

		cloned := proto.Clone(attr).(*permissionsattributes.RoleAttribute)
		if cloned.GetValidValues() != nil {
			cloned.Value = proto.Clone(cloned.GetValidValues()).(*permissionsattributes.AttributeValues)
		}
		result = append(result, cloned)
	}

	return result
}

func (ps *Perms) GetRoleAttributes(
	ctx context.Context,
	job string,
	grade int32,
) ([]*permissionsattributes.RoleAttribute, error) {
	roleId, ok := ps.lookupRoleIDForJobAndGrade(job, grade)
	if !ok {
		roleId, ok = ps.lookupRoleIDForJobAndGrade(DefaultRoleJob, ps.startJobGrade)
		if !ok {
			return nil, fmt.Errorf(
				"failed to fallback to default role. %w",
				errors.New("no role ID found for job and grade"),
			)
		}
	}

	attrsRoleMap, ok := ps.attrsRoleMap.Load(roleId)
	if !ok {
		return []*permissionsattributes.RoleAttribute{}, nil
	}

	var err error
	attrs := []*permissionsattributes.RoleAttribute{}
	for key := range attrsRoleMap.All() {
		attr, ok := ps.lookupAttributeByID(key)
		if !ok {
			err = fmt.Errorf(
				"no attribute found by id for role. %w",
				errors.New("attribute ID not found"),
			)
			break
		}

		attrVal, ok := attrsRoleMap.Load(attr.ID)
		if !ok {
			err = fmt.Errorf(
				"no role attribute found by id for role. %w",
				errors.New("role attribute ID not found"),
			)
			break
		}

		maxValues, er := ps.GetJobAttributeValue(ctx, job, attr.ID)
		if er != nil {
			err = fmt.Errorf(
				"failed to retrieve max values for attribute %s.%s/%s/%s. %w",
				attr.Namespace,
				attr.Service,
				attr.Name,
				attr.Key,
				er,
			)
			break
		}

		attrs = append(attrs, &permissionsattributes.RoleAttribute{
			RoleId:       roleId,
			AttrId:       attr.ID,
			PermissionId: attr.PermissionID,
			Namespace:    string(attr.Namespace),
			Service:      string(attr.Service),
			Name:         string(attr.Name),
			Key:          string(attr.Key),
			Type:         string(attr.Type),
			Value:        attrVal.Value,
			ValidValues:  attr.ValidValues,
			MaxValues:    maxValues,
		})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to range over attributes. %w", err)
	}
	if err := ps.applyDynamicPropertyListOptions(ctx, attrs); err != nil {
		return nil, err
	}

	return attrs, nil
}

func (ps *Perms) GetEffectiveRoleAttributes(
	ctx context.Context,
	job string,
	grade int32,
) ([]*permissionsattributes.RoleAttribute, error) {
	roleAttrs := map[int64]struct{}{}

	roleIds, ok := ps.lookupRoleIDsForJobUpToGrade(job, grade)
	if !ok {
		return nil, nil
	}

	perms := ps.getRolePermissionsFromCache(roleIds)

	var err error
	attrs := []*permissionsattributes.RoleAttribute{}
	for i := range slices.Backward(roleIds) {
		attrMap, ok := ps.attrsRoleMap.Load(roleIds[i])
		if !ok {
			continue
		}

		for _, value := range attrMap.All() {
			// Skip already added attributes
			if _, ok := roleAttrs[value.AttrID]; ok {
				continue
			}

			// Permission not granted
			if !slices.ContainsFunc(perms, func(p *cachePerm) bool {
				return p.ID == value.PermissionID
			}) {
				continue
			}

			attr, ok := ps.lookupAttributeByID(value.AttrID)
			if !ok {
				err = fmt.Errorf(
					"no attribute found by id for role. %w",
					errors.New("attribute ID not found"),
				)
				break
			}

			attrMap, ok := ps.attrsRoleMap.Load(roleIds[i])
			if !ok {
				continue
			}
			attrVal, ok := attrMap.Load(attr.ID)
			if !ok {
				err = fmt.Errorf(
					"no role attribute found by id for role. %w",
					errors.New("role attribute ID not found"),
				)
				break
			}

			maxVal, _ := ps.GetJobAttributeValue(ctx, job, attr.ID)

			attrs = append(attrs, &permissionsattributes.RoleAttribute{
				RoleId:       roleIds[i],
				AttrId:       value.AttrID,
				PermissionId: value.PermissionID,
				Namespace:    string(attr.Namespace),
				Service:      string(attr.Service),
				Name:         string(attr.Name),
				Key:          string(attr.Key),
				Type:         string(attr.Type),
				Value:        attrVal.Value,
				ValidValues:  attr.ValidValues,
				MaxValues:    maxVal,
			})

			roleAttrs[value.AttrID] = struct{}{}
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to range over attributes. %w", err)
	}
	if err := ps.applyDynamicPropertyListOptions(ctx, attrs); err != nil {
		return nil, err
	}

	return attrs, nil
}

func (ps *Perms) UpdateRoleAttributes(
	ctx context.Context,
	job string,
	roleId int64,
	attrs ...*permissionsattributes.RoleAttribute,
) error {
	for i := range attrs {
		attrs[i].RoleId = roleId

		a, ok := ps.lookupAttributeByID(attrs[i].GetAttrId())
		if !ok {
			return fmt.Errorf(
				"no attribute found by id %d. %w",
				attrs[i].GetAttrId(),
				errors.New("attribute not found"),
			)
		}

		if attrs[i].GetValue() != nil {
			attrs[i].GetValue().Default(permissionsattributes.AttributeTypes(attrs[i].GetType()))

			maxValues, err := ps.GetJobAttributeValue(ctx, job, a.ID)
			if err != nil {
				return fmt.Errorf(
					"failed to retrieve max values for attribute %s/%s/%s. %w",
					a.Namespace,
					a.Name,
					a.Key,
					err,
				)
			}

			validValues := a.ValidValues
			if a.Type == permissionsattributes.PropertyListAttributeType &&
				a.PropertyScope != propertypb.PropertyScope_PROPERTY_SCOPE_UNSPECIFIED {
				validValues, err = ps.dynamicPropertyListValidValues(
					ctx,
					a.PropertyScope,
					a.ValidValues,
				)
				if err != nil {
					return err
				}
				if roleAttrs, ok := ps.attrsRoleMap.Load(roleId); ok {
					if existing, ok := roleAttrs.Load(a.ID); ok {
						validValues = includeExistingPropertyListGrants(
							validValues,
							existing.Value,
							a.PropertyScope,
						)
					}
				}
			}

			valid, _ := attrs[i].GetValue().Check(a.Type, validValues, maxValues)
			if !valid {
				return errors.Wrapf(
					ErrAttrInvalid,
					"attribute %s/%s failed validation",
					a.Key,
					a.Name,
				)
			}
		}
	}

	if err := ps.addOrUpdateAttributesToRole(ctx, roleId, attrs...); err != nil {
		return fmt.Errorf("failed to add or update attributes to role. %w", err)
	}

	if err := ps.publishMessage(ctx, RoleAttrUpdateSubject, &permissionsevents.RoleIDEvent{
		RoleId: roleId,
	}); err != nil {
		return fmt.Errorf("failed to publish role attribute update message. %w", err)
	}

	return nil
}

func includeExistingPropertyListGrants(
	validValues *permissionsattributes.AttributeValues,
	existingValue *permissionsattributes.AttributeValues,
	scope propertypb.PropertyScope,
) *permissionsattributes.AttributeValues {
	if existingValue == nil || len(existingValue.GetPropertyList().GetKeys()) == 0 {
		return validValues
	}

	result := proto.Clone(validValues).(*permissionsattributes.AttributeValues)
	keys := append([]string(nil), result.GetPropertyList().GetKeys()...)
	for _, key := range existingValue.GetPropertyList().GetKeys() {
		if propertyregistry.IsCustomPropertyKey(scope, key) {
			keys = append(keys, key)
		}
	}
	slices.Sort(keys)
	keys = slices.Compact(keys)
	result.ValidValues = &permissionsattributes.AttributeValues_PropertyList{
		PropertyList: &permissionsattributes.PropertyList{Keys: keys},
	}
	return result
}

func (ps *Perms) addOrUpdateAttributesToRole(
	ctx context.Context,
	roleId int64,
	attrs ...*permissionsattributes.RoleAttribute,
) error {
	for i := range attrs {
		a, ok := ps.lookupAttributeByID(attrs[i].GetAttrId())
		if !ok {
			return fmt.Errorf(
				"unable to add role attribute, didn't find attribute by ID %d. %w",
				attrs[i].GetAttrId(),
				errors.New("attribute not found"),
			)
		}

		if attrs[i].GetValue() == nil {
			attrs[i].Value = &permissionsattributes.AttributeValues{}
		}

		if attrs[i].GetValue() != nil {
			attrs[i].GetValue().Default(a.Type)
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
				attrs[i].GetValue(),
			).
			ON_DUPLICATE_KEY_UPDATE(
				tRoleAttrs.Value.SET(mysql.StringExp(mysql.Raw("VALUES(`value`)"))),
			)

		if _, err := stmt.ExecContext(ctx, ps.db); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return fmt.Errorf("failed to execute insert statement for role attributes. %w", err)
			}
		}

		ps.updateRoleAttributeInMap(
			roleId,
			a.PermissionID,
			a.ID,
			a.Key,
			a.Type,
			attrs[i].GetValue(),
		)
	}

	return nil
}

func (ps *Perms) RemoveAttributesFromRole(
	ctx context.Context,
	roleId int64,
	attrs ...*permissionsattributes.RoleAttribute,
) error {
	if len(attrs) == 0 {
		return nil
	}

	ids := make([]mysql.Expression, len(attrs))
	for i := range attrs {
		ids[i] = mysql.Int64(attrs[i].GetAttrId())
	}

	stmt := tRoleAttrs.
		DELETE().
		WHERE(mysql.AND(
			tRoleAttrs.RoleID.EQ(mysql.Int64(roleId)),
			tRoleAttrs.AttrID.IN(ids...),
		)).
		LIMIT(int64(len(ids)))

	if _, err := stmt.ExecContext(ctx, ps.db); err != nil {
		return fmt.Errorf("failed to execute delete statement for role attributes. %w", err)
	}

	for i := range attrs {
		ps.removeRoleAttributeFromMap(roleId, attrs[i].GetAttrId())
	}

	if err := ps.publishMessage(ctx, RoleAttrUpdateSubject, &permissionsevents.RoleIDEvent{
		RoleId: roleId,
	}); err != nil {
		return fmt.Errorf("failed to publish role attribute removal message. %w", err)
	}

	return nil
}

func (ps *Perms) RemoveAttributesFromRoleByPermission(
	ctx context.Context,
	roleId int64,
	permissionId int64,
) error {
	as, ok := ps.attrsPermsMap.Load(permissionId)
	if !ok {
		return nil
	}

	ras := []*permissionsattributes.RoleAttribute{}
	for _, attrId := range as.All() {
		ras = append(ras, &permissionsattributes.RoleAttribute{
			AttrId: attrId,
		})
	}

	if len(ras) == 0 {
		return nil
	}

	if err := ps.RemoveAttributesFromRole(ctx, roleId, ras...); err != nil {
		return fmt.Errorf("failed to remove attributes from role by perm. %w", err)
	}

	if err := ps.publishMessage(ctx, RoleAttrUpdateSubject, &permissionsevents.RoleIDEvent{
		RoleId: roleId,
	}); err != nil {
		return fmt.Errorf("failed to publish role attribute removal message. %w", err)
	}

	return nil
}

func (ps *Perms) UpdateJobAttributes(
	ctx context.Context,
	job string,
	attrs ...*permissionsattributes.RoleAttribute,
) error {
	for _, attr := range attrs {
		a, ok := ps.lookupAttributeByID(attr.GetAttrId())
		if !ok {
			return fmt.Errorf(
				"unable to update role attribute max values, didn't find attribute by ID %d. %w",
				attr.GetAttrId(),
				errors.New("attribute not found"),
			)
		}

		maxVal := mysql.NULL
		if attr.GetMaxValues() != nil {
			attr.GetMaxValues().Default(a.Type)

			out, err := protoutils.MarshalToJSON(attr.GetMaxValues())
			if err != nil {
				return fmt.Errorf("failed to marshal max values. %w", err)
			}

			maxVal = mysql.String(string(out))
		}

		stmt := tJobAttrs.
			INSERT(
				tJobAttrs.Job,
				tJobAttrs.AttrID,
				tJobAttrs.MaxValues,
			).
			VALUES(
				job,
				attr.GetAttrId(),
				maxVal,
			).
			ON_DUPLICATE_KEY_UPDATE(
				tJobAttrs.MaxValues.SET(mysql.StringExp(mysql.Raw("VALUES(`max_values`)"))),
			)

		if _, err := stmt.ExecContext(ctx, ps.db); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return fmt.Errorf("failed to execute insert statement for job attributes. %w", err)
			}
		}
	}

	return nil
}

func (ps *Perms) ClearJobAttributes(ctx context.Context, job string) error {
	var count database.DataCount
	countStmt := tJobAttrs.
		SELECT(mysql.COUNT(tJobAttrs.AttrID).AS("data_count.total")).
		FROM(tJobAttrs).
		WHERE(tJobAttrs.Job.EQ(mysql.String(job)))
	if err := countStmt.QueryContext(ctx, ps.db, &count); err != nil {
		return fmt.Errorf("failed to execute count statement for job attributes. %w", err)
	}

	stmt := tJobAttrs.
		DELETE().
		WHERE(tJobAttrs.Job.EQ(mysql.String(job))).
		LIMIT(count.Total)

	if _, err := stmt.ExecContext(ctx, ps.db); err != nil {
		return fmt.Errorf("failed to execute delete statement for job attributes. %w", err)
	}

	return nil
}

func (ps *Perms) updateRoleAttributeInMap(
	roleId int64,
	permId int64,
	attrId int64,
	key Key,
	aType permissionsattributes.AttributeTypes,
	value *permissionsattributes.AttributeValues,
) {
	job, ok := ps.lookupJobForRoleID(roleId)
	if !ok {
		ps.logger.Error("unable to lookup job for role id", zap.Int64("role_id", roleId))
		return
	}

	attrRoleMap, _ := ps.attrsRoleMap.LoadOrCompute(
		roleId,
		func() (*xsync.Map[int64, *cacheRoleAttr], bool) {
			return xsync.NewMap[int64, *cacheRoleAttr](), false
		},
	)

	attrRoleMap.Store(attrId, &cacheRoleAttr{
		Job:          job,
		AttrID:       attrId,
		PermissionID: permId,
		Key:          key,
		Type:         aType,
		Value:        value,
	})
}

func (ps *Perms) removeRoleAttributeFromMap(roleId int64, attrId int64) {
	attrMap, ok := ps.attrsRoleMap.Load(roleId)
	if !ok {
		return
	}

	attrMap.Delete(attrId)
}
