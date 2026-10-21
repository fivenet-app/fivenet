package properties

import (
	"context"
	"slices"

	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ResolveDefinitions returns effective definitions for stored property keys.
// System definitions are resolved from the registry; custom definitions come
// from the database. The result is limited to the requested keys.
func ResolveDefinitions(
	ctx context.Context,
	db qrm.DB,
	scope propertypb.PropertyScope,
	keys []string,
) (map[string]*propertypb.PropertyDefinition, error) {
	result := make(map[string]*propertypb.PropertyDefinition, len(keys))
	if len(keys) == 0 {
		return result, nil
	}

	requested := make([]mysql.Expression, 0, len(keys))
	for _, key := range keys {
		requested = append(requested, mysql.String(key))
	}
	for _, definition := range SystemPropertyDefinitionsForScope(scope) {
		if slices.Contains(keys, definition.GetKey()) {
			result[definition.GetKey()] = proto.Clone(definition).(*propertypb.PropertyDefinition)
		}
	}

	scopeName, ok := ScopeName(scope)
	if !ok {
		return result, nil
	}
	tDefinitions := table.FivenetPropertyDefinitions
	rows := []model.FivenetPropertyDefinitions{}
	stmt := tDefinitions.SELECT(tDefinitions.AllColumns).FROM(tDefinitions).WHERE(
		tDefinitions.Scope.EQ(mysql.String(scopeName)).AND(
			tDefinitions.PropertyKey.IN(requested...).AND(tDefinitions.DeletedAt.IS_NULL()),
		),
	)
	if err := stmt.QueryContext(ctx, db, &rows); err != nil {
		return nil, err
	}
	for _, row := range rows {
		// System definitions are code-owned. A stale or manually inserted
		// database row must not replace their type, renderer, or behavior.
		if IsSystemProperty(scope, row.PropertyKey) {
			continue
		}
		definition, err := DefinitionFromModel(scope, row)
		if err != nil {
			return nil, err
		}
		result[row.PropertyKey] = definition
	}

	overrides := []model.FivenetPropertySystemOverrides{}
	overrideStmt := table.FivenetPropertySystemOverrides.SELECT(
		table.FivenetPropertySystemOverrides.PropertyKey,
		table.FivenetPropertySystemOverrides.Enabled,
		table.FivenetPropertySystemOverrides.SortOrder,
	).FROM(table.FivenetPropertySystemOverrides).WHERE(
		table.FivenetPropertySystemOverrides.Scope.EQ(mysql.String(scopeName)).AND(
			table.FivenetPropertySystemOverrides.PropertyKey.IN(requested...),
		),
	)
	if err := overrideStmt.QueryContext(ctx, db, &overrides); err != nil {
		return nil, err
	}
	for _, override := range overrides {
		if definition, ok := result[override.PropertyKey]; ok &&
			definition.GetSource() == propertypb.PropertyDefinitionSource_PROPERTY_DEFINITION_SOURCE_SYSTEM {
			definition.Enabled = override.Enabled
			if override.SortOrder != nil {
				definition.SortOrder = *override.SortOrder
			}
		}
	}
	if err := ApplyGroupSortOrders(ctx, db, scope, result); err != nil {
		return nil, err
	}
	return result, nil
}

// ApplyGroupSortOrders attaches the effective instance-specific group order
// to resolved definitions so clients can render groups without duplicating
// the system default order or maintaining a hardcoded list.
func ApplyGroupSortOrders(
	ctx context.Context,
	db qrm.DB,
	scope propertypb.PropertyScope,
	definitions map[string]*propertypb.PropertyDefinition,
) error {
	orders := make(map[string]int32)
	for _, group := range SystemPropertyGroups() {
		if group.GetScope() == scope {
			orders[group.GetKey()] = group.GetSortOrder()
		}
	}
	name, ok := ScopeName(scope)
	if !ok {
		return nil
	}
	tGroups := table.FivenetPropertyGroupOverrides
	rows := []model.FivenetPropertyGroupOverrides{}
	if err := tGroups.SELECT(tGroups.GroupKey, tGroups.SortOrder).
		FROM(tGroups).
		WHERE(tGroups.Scope.EQ(mysql.String(name))).
		QueryContext(ctx, db, &rows); err != nil {
		return err
	}
	for _, row := range rows {
		orders[row.GroupKey] = row.SortOrder
	}
	tCustomGroups := table.FivenetPropertyGroups
	customGroups := []model.FivenetPropertyGroups{}
	if err := tCustomGroups.SELECT(tCustomGroups.GroupKey, tCustomGroups.SortOrder).
		FROM(tCustomGroups).
		WHERE(tCustomGroups.Scope.EQ(mysql.String(name)).AND(tCustomGroups.DeletedAt.IS_NULL())).
		QueryContext(ctx, db, &customGroups); err != nil {
		return err
	}
	for _, row := range customGroups {
		orders[row.GroupKey] = row.SortOrder
	}
	for _, definition := range definitions {
		definition.GroupSortOrder = 1 << 30
		if order, ok := orders[definition.GetGroupKey()]; ok {
			definition.GroupSortOrder = order
		}
	}
	return nil
}

// ResolvePropertyGroupsFromDB adds instance-defined group labels to the
// effective groups. System groups remain registry-owned.
func ResolvePropertyGroupsFromDB(
	ctx context.Context,
	db qrm.DB,
	scope propertypb.PropertyScope,
	definitions map[string]*propertypb.PropertyDefinition,
) ([]*propertypb.PropertyGroup, error) {
	groups := ResolvePropertyGroups(definitions)
	name, ok := ScopeName(scope)
	if !ok {
		return groups, nil
	}
	rows := []model.FivenetPropertyGroups{}
	tGroups := table.FivenetPropertyGroups
	if err := tGroups.SELECT(tGroups.GroupKey, tGroups.Label, tGroups.SortOrder).
		FROM(tGroups).
		WHERE(tGroups.Scope.EQ(mysql.String(name)).AND(tGroups.DeletedAt.IS_NULL())).
		QueryContext(ctx, db, &rows); err != nil {
		return nil, err
	}
	labels := make(map[string]string, len(rows))
	for _, row := range rows {
		labels[row.GroupKey] = row.Label
	}
	for _, group := range groups {
		if label, ok := labels[group.GetKey()]; ok {
			group.Label = label
		}
	}
	return groups, nil
}

func uint32Pointer(value *uint8) *uint32 {
	if value == nil {
		return nil
	}
	converted := uint32(*value)
	return &converted
}

// DefinitionFromModel translates the persistence representation to the wire
// resource. This intentionally remains an explicit boundary: the database
// uses string scopes, integer enums, JSON renderer config, and nullable
// narrow integer columns that do not map one-to-one to the protobuf message.
func DefinitionFromModel(
	scope propertypb.PropertyScope,
	row model.FivenetPropertyDefinitions,
) (*propertypb.PropertyDefinition, error) {
	definition := &propertypb.PropertyDefinition{
		Scope:         scope,
		Source:        propertypb.PropertyDefinitionSource_PROPERTY_DEFINITION_SOURCE_CUSTOM,
		Key:           row.PropertyKey,
		ValueType:     propertypb.PropertyValueType(row.ValueType),
		Label:         row.Label,
		Enabled:       row.Enabled,
		Searchable:    row.Searchable,
		ValueFormat:   propertypb.PropertyValueFormat(row.ValueFormat),
		Behavior:      propertypb.PropertyBehavior(row.Behavior),
		Unit:          row.Unit,
		DecimalPlaces: uint32Pointer(row.DecimalPlaces),
		Currency:      row.Currency,
		GroupKey:      row.GroupKey,
		SortOrder:     row.SortOrder,
		Pinned:        row.Pinned,
		Renderer:      propertypb.PropertyRenderer(row.Renderer),
		DisplayMode:   propertypb.PropertyDisplayMode(row.DisplayMode),
		Implicit:      row.Implicit,
	}
	if row.Description != nil {
		definition.Description = row.Description
	}
	if row.DeletedAt != nil {
		definition.DeletedAt = timestamp.New(*row.DeletedAt)
	}
	if row.RendererConfig != nil {
		definition.RendererConfig = &propertypb.PropertyRendererConfig{}
		if err := protojson.Unmarshal(
			[]byte(*row.RendererConfig),
			definition.RendererConfig,
		); err != nil {
			return nil, err
		}
	}
	return definition, nil
}
