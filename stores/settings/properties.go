package settingsstore

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	_ "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens"
	_ "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/vehicles"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	propertyregistry "github.com/fivenet-app/fivenet/v2026/pkg/properties"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	ErrInvalidPropertyDefinition  = errors.New("invalid property definition")
	ErrPropertyDefinitionNotFound = errors.New("property definition not found")
	ErrPropertyGroupNotFound      = errors.New("property group not found")
)

// Kept as local aliases for SQL expectation tests; scope semantics live in
// pkg/properties.ScopeName.
const (
	vehiclePropertyScopeName = "vehicle"
	userPropertyScopeName    = "user"
)

func (s *Store) ListPropertyDefinitions(
	ctx context.Context,
	scope properties.PropertyScope,
	includeDeleted bool,
) ([]*properties.PropertyDefinition, error) {
	scopeName, ok := propertyScopeName(scope)
	if !ok {
		return nil, ErrInvalidPropertyDefinition
	}

	overrides, err := s.systemPropertyOverrides(ctx, scopeName)
	if err != nil {
		return nil, err
	}
	systemDefinitions := propertyregistry.SystemPropertyDefinitionsForScope(scope)
	customDefinitions := make([]*properties.PropertyDefinition, 0, 8)

	tDefinitions := table.FivenetPropertyDefinitions
	stmt := tDefinitions.
		SELECT(
			tDefinitions.PropertyKey,
			tDefinitions.ValueType,
			tDefinitions.Label,
			tDefinitions.Description,
			tDefinitions.Enabled,
			tDefinitions.Searchable,
			tDefinitions.ValueFormat,
			tDefinitions.Behavior,
			tDefinitions.Unit,
			tDefinitions.DecimalPlaces,
			tDefinitions.Currency,
			tDefinitions.GroupKey,
			tDefinitions.SortOrder,
			tDefinitions.Pinned,
			tDefinitions.Renderer,
			tDefinitions.RendererConfig,
			tDefinitions.DisplayMode,
			tDefinitions.Implicit,
			tDefinitions.DeletedAt,
		).
		FROM(tDefinitions).
		WHERE(tDefinitions.Scope.EQ(mysql.String(scopeName)).AND(
			tDefinitions.DeletedAt.IS_NULL().OR(mysql.Bool(includeDeleted)),
		)).
		ORDER_BY(tDefinitions.PropertyKey.ASC())
	rows := []model.FivenetPropertyDefinitions{}
	err = stmt.QueryContext(ctx, s.db, &rows)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		definition, err := propertyregistry.DefinitionFromModel(scope, row)
		if err != nil {
			return nil, err
		}
		customDefinitions = append(customDefinitions, definition)
	}
	definitions := propertyregistry.MergeDefinitions(
		systemDefinitions,
		customDefinitions,
		overrides.enabled,
	)
	definitions = propertyregistry.ApplySortOrderOverrides(definitions, overrides.sortOrders)
	if err := propertyregistry.ApplyGroupSortOrders(
		ctx,
		s.db,
		scope,
		definitionMap(definitions),
	); err != nil {
		return nil, err
	}
	definitions = propertyregistry.SortDefinitions(definitions)
	return definitions, nil
}

func definitionMap(
	definitions []*properties.PropertyDefinition,
) map[string]*properties.PropertyDefinition {
	result := make(map[string]*properties.PropertyDefinition, len(definitions))
	for _, definition := range definitions {
		result[definition.GetKey()] = definition
	}
	return result
}

func (s *Store) ListPropertyGroups(
	ctx context.Context,
	scope properties.PropertyScope,
	includeDeleted bool,
) ([]*properties.PropertyGroup, error) {
	if _, ok := propertyScopeName(scope); !ok {
		return nil, ErrInvalidPropertyDefinition
	}
	groups := make([]*properties.PropertyGroup, 0)
	for _, group := range propertyregistry.SystemPropertyGroups() {
		if group.GetScope() == scope {
			groups = append(groups, group)
		}
	}
	name, _ := propertyScopeName(scope)
	tGroups := table.FivenetPropertyGroupOverrides
	rows := []model.FivenetPropertyGroupOverrides{}
	stmt := tGroups.SELECT(tGroups.GroupKey, tGroups.SortOrder).
		FROM(tGroups).
		WHERE(tGroups.Scope.EQ(mysql.String(name)))
	err := stmt.QueryContext(ctx, s.db, &rows)
	if err != nil {
		return nil, err
	}
	orders := map[string]int32{}
	for _, row := range rows {
		orders[row.GroupKey] = row.SortOrder
	}
	tCustom := table.FivenetPropertyGroups
	customRows := []model.FivenetPropertyGroups{}
	customStmt := tCustom.SELECT(tCustom.GroupKey, tCustom.Label, tCustom.SortOrder, tCustom.DeletedAt).
		FROM(tCustom).
		WHERE(tCustom.Scope.EQ(mysql.String(name)).AND(
			tCustom.DeletedAt.IS_NULL().OR(mysql.Bool(includeDeleted)),
		))
	if err := customStmt.QueryContext(ctx, s.db, &customRows); err != nil {
		return nil, err
	}
	for _, row := range customRows {
		groups = append(groups, propertyGroupFromRow(scope, row))
	}
	for _, group := range groups {
		if order, ok := orders[group.GetKey()]; ok {
			group.SortOrder = order
		}
	}
	sort.SliceStable(
		groups,
		func(i, j int) bool { return groups[i].GetSortOrder() < groups[j].GetSortOrder() },
	)
	return groups, nil
}

func (s *Store) CreatePropertyGroup(
	ctx context.Context,
	group *properties.PropertyGroup,
) (*properties.PropertyGroup, error) {
	scopeName, err := validateCustomPropertyGroup(group)
	if err != nil {
		return nil, err
	}
	t := table.FivenetPropertyGroups
	_, err = t.INSERT(t.Scope, t.GroupKey, t.Label, t.SortOrder).
		VALUES(scopeName, group.GetKey(), group.GetLabel(), group.GetSortOrder()).
		ExecContext(ctx, s.db)
	if err != nil {
		return nil, err
	}
	s.propertyCache.Invalidate(group.GetScope())
	return proto.Clone(group).(*properties.PropertyGroup), nil
}

func (s *Store) UpdatePropertyGroup(
	ctx context.Context,
	group *properties.PropertyGroup,
) (*properties.PropertyGroup, error) {
	scopeName, err := validateCustomPropertyGroup(group)
	if err != nil {
		return nil, err
	}
	t := table.FivenetPropertyGroups
	result, err := t.UPDATE().
		SET(t.Label.SET(mysql.String(group.GetLabel())), t.SortOrder.SET(mysql.Int32(group.GetSortOrder()))).
		WHERE(t.Scope.EQ(mysql.String(scopeName)).AND(t.GroupKey.EQ(mysql.String(group.GetKey()))).AND(t.DeletedAt.IS_NULL())).
		ExecContext(ctx, s.db)
	if err != nil {
		return nil, err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if changed == 0 {
		return nil, ErrPropertyGroupNotFound
	}
	s.propertyCache.Invalidate(group.GetScope())
	return proto.Clone(group).(*properties.PropertyGroup), nil
}

func (s *Store) DeletePropertyGroup(
	ctx context.Context,
	scope properties.PropertyScope,
	key string,
) error {
	scopeName, ok := propertyScopeName(scope)
	if !ok || propertyregistry.ValidateKey(key) != nil {
		return ErrInvalidPropertyDefinition
	}
	if _, system := propertyregistry.SystemPropertyGroup(scope, key); system {
		return ErrInvalidPropertyDefinition
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	t := table.FivenetPropertyGroups
	result, err := t.UPDATE().
		SET(t.DeletedAt.SET(mysql.CURRENT_TIMESTAMP())).
		WHERE(t.Scope.EQ(mysql.String(scopeName)).AND(t.GroupKey.EQ(mysql.String(key))).AND(t.DeletedAt.IS_NULL())).
		ExecContext(ctx, tx)
	if err != nil {
		return err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if changed == 0 {
		return ErrPropertyGroupNotFound
	}
	// Deleting a group never deletes definitions; they become ungrouped.
	defs := table.FivenetPropertyDefinitions
	if _, err = defs.UPDATE().
		SET(defs.GroupKey.SET(mysql.String(""))).
		WHERE(defs.Scope.EQ(mysql.String(scopeName)).AND(defs.GroupKey.EQ(mysql.String(key))).AND(defs.DeletedAt.IS_NULL())).
		ExecContext(ctx, tx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	s.propertyCache.Invalidate(scope)
	return nil
}

func (s *Store) RestorePropertyGroup(
	ctx context.Context,
	scope properties.PropertyScope,
	key string,
) (*properties.PropertyGroup, error) {
	scopeName, ok := propertyScopeName(scope)
	if !ok || propertyregistry.ValidateKey(key) != nil {
		return nil, ErrInvalidPropertyDefinition
	}
	if _, system := propertyregistry.SystemPropertyGroup(scope, key); system {
		return nil, ErrInvalidPropertyDefinition
	}
	t := table.FivenetPropertyGroups
	result, err := t.UPDATE().
		SET(t.DeletedAt.SET(mysql.TimestampExp(mysql.NULL))).
		WHERE(t.Scope.EQ(mysql.String(scopeName)).AND(t.GroupKey.EQ(mysql.String(key))).AND(t.DeletedAt.IS_NOT_NULL())).
		ExecContext(ctx, s.db)
	if err != nil {
		return nil, err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if changed == 0 {
		return nil, ErrPropertyGroupNotFound
	}
	rows := []model.FivenetPropertyGroups{}
	if err := t.SELECT(t.GroupKey, t.Label, t.SortOrder, t.DeletedAt).
		FROM(t).
		WHERE(t.Scope.EQ(mysql.String(scopeName)).AND(t.GroupKey.EQ(mysql.String(key)))).
		QueryContext(ctx, s.db, &rows); err != nil {
		return nil, err
	}
	if len(rows) != 1 {
		return nil, ErrPropertyGroupNotFound
	}
	s.propertyCache.Invalidate(scope)
	return propertyGroupFromRow(scope, rows[0]), nil
}

func propertyGroupFromRow(
	scope properties.PropertyScope,
	row model.FivenetPropertyGroups,
) *properties.PropertyGroup {
	group := &properties.PropertyGroup{
		Scope:     scope,
		Key:       row.GroupKey,
		Label:     row.Label,
		SortOrder: row.SortOrder,
	}
	if row.DeletedAt != nil {
		group.DeletedAt = timestamp.New(*row.DeletedAt)
	}
	return group
}

func validateCustomPropertyGroup(group *properties.PropertyGroup) (string, error) {
	if group == nil || strings.TrimSpace(group.GetLabel()) == "" ||
		propertyregistry.ValidateKey(group.GetKey()) != nil {
		return "", ErrInvalidPropertyDefinition
	}
	group.Label = strings.TrimSpace(group.GetLabel())
	group.DisplayKey = ""
	group.DeletedAt = nil
	if _, ok := propertyregistry.ScopeName(group.GetScope()); !ok {
		return "", ErrInvalidPropertyDefinition
	}
	if _, system := propertyregistry.SystemPropertyGroup(group.GetScope(), group.GetKey()); system {
		return "", ErrInvalidPropertyDefinition
	}
	scopeName, _ := propertyregistry.ScopeName(group.GetScope())
	return scopeName, nil
}

func (s *Store) ReorderPropertyGroups(
	ctx context.Context,
	scope properties.PropertyScope,
	keys []string,
) error {
	if len(keys) == 0 {
		return ErrInvalidPropertyDefinition
	}
	groups, err := s.ListPropertyGroups(ctx, scope, false)
	if err != nil {
		return err
	}
	valid := make(map[string]struct{}, len(groups))
	for _, group := range groups {
		valid[group.GetKey()] = struct{}{}
	}
	seen := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		if _, ok := valid[key]; !ok {
			return ErrInvalidPropertyDefinition
		}
		if _, ok := seen[key]; ok {
			return ErrInvalidPropertyDefinition
		}
		seen[key] = struct{}{}
	}
	if len(seen) != len(groups) {
		return ErrInvalidPropertyDefinition
	}
	name, _ := propertyScopeName(scope)
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	tGroups := table.FivenetPropertyGroupOverrides
	for index, key := range keys {
		order := int32((index + 1) * 100)
		if _, system := propertyregistry.SystemPropertyGroup(scope, key); !system {
			if _, err := table.FivenetPropertyGroups.UPDATE().SET(
				table.FivenetPropertyGroups.SortOrder.SET(mysql.Int32(order)),
			).WHERE(
				table.FivenetPropertyGroups.Scope.EQ(mysql.String(name)).AND(
					table.FivenetPropertyGroups.GroupKey.EQ(mysql.String(key)).AND(
						table.FivenetPropertyGroups.DeletedAt.IS_NULL(),
					),
				),
			).ExecContext(ctx, tx); err != nil {
				return err
			}
			continue
		}
		if _, err := tGroups.
			INSERT(tGroups.Scope, tGroups.GroupKey, tGroups.SortOrder).
			VALUES(name, key, order).
			ON_DUPLICATE_KEY_UPDATE(
				tGroups.SortOrder.SET(mysql.IntExp(mysql.RawString("VALUES(`sort_order`)"))),
			).
			ExecContext(ctx, tx); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	s.propertyCache.Invalidate(scope)
	return nil
}

func (s *Store) ReorderPropertyDefinitions(
	ctx context.Context,
	scope properties.PropertyScope,
	groupKey string,
	keys []string,
) error {
	if len(keys) == 0 {
		return ErrInvalidPropertyDefinition
	}
	name, ok := propertyScopeName(scope)
	if !ok {
		return ErrInvalidPropertyDefinition
	}
	seen := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		if propertyregistry.ValidateKey(key) != nil {
			return ErrInvalidPropertyDefinition
		}
		if systemDefinition, isSystem := propertyregistry.SystemPropertyDefinition(
			scope,
			key,
		); isSystem {
			if systemDefinition.GetGroupKey() != groupKey {
				return ErrInvalidPropertyDefinition
			}
		}
		if _, exists := seen[key]; exists {
			return ErrInvalidPropertyDefinition
		}
		seen[key] = struct{}{}
	}
	definitions, err := s.ListPropertyDefinitions(ctx, scope, false)
	if err != nil {
		return err
	}
	expected := make(map[string]struct{})
	for _, definition := range definitions {
		if definition.GetGroupKey() == groupKey {
			expected[definition.GetKey()] = struct{}{}
		}
	}
	if len(expected) != len(seen) {
		return ErrInvalidPropertyDefinition
	}
	for key := range expected {
		if _, ok := seen[key]; !ok {
			return ErrInvalidPropertyDefinition
		}
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	tDefinitions := table.FivenetPropertyDefinitions
	tOverrides := table.FivenetPropertySystemOverrides
	for index, key := range keys {
		order := int32((index + 1) * 100)
		if propertyregistry.IsSystemProperty(scope, key) {
			if _, err := tOverrides.
				INSERT(
					tOverrides.Scope,
					tOverrides.PropertyKey,
					tOverrides.Enabled,
					tOverrides.SortOrder,
				).
				VALUES(name, key, true, order).
				ON_DUPLICATE_KEY_UPDATE(
					tOverrides.SortOrder.SET(mysql.IntExp(mysql.RawString("VALUES(`sort_order`)"))),
				).
				ExecContext(ctx, tx); err != nil {
				return err
			}
			continue
		}
		result, err := tDefinitions.UPDATE().SET(
			tDefinitions.SortOrder.SET(mysql.Int32(order)),
		).WHERE(
			tDefinitions.Scope.EQ(mysql.String(name)).AND(
				tDefinitions.PropertyKey.EQ(mysql.String(key)).AND(
					tDefinitions.GroupKey.EQ(mysql.String(groupKey)).AND(
						tDefinitions.DeletedAt.IS_NULL(),
					),
				),
			),
		).ExecContext(ctx, tx)
		if err != nil {
			return err
		}
		affected, err := result.RowsAffected()
		if err != nil || affected != 1 {
			return fmt.Errorf("%w: %s", ErrPropertyDefinitionNotFound, key)
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	s.propertyCache.Invalidate(scope)
	return nil
}

func (s *Store) CreatePropertyDefinition(
	ctx context.Context,
	definition *properties.PropertyDefinition,
) (*properties.PropertyDefinition, error) {
	definition, scopeName, err := s.validateCustomPropertyDefinition(ctx, definition)
	if err != nil {
		return nil, err
	}
	rendererConfig, err := rendererConfigJSON(definition.GetRendererConfig())
	if err != nil {
		return nil, err
	}

	tDefinitions := table.FivenetPropertyDefinitions
	_, err = tDefinitions.INSERT(
		tDefinitions.Scope, tDefinitions.PropertyKey, tDefinitions.ValueType, tDefinitions.Label,
		tDefinitions.Description, tDefinitions.Enabled, tDefinitions.Searchable,
		tDefinitions.ValueFormat, tDefinitions.Behavior, tDefinitions.Unit,
		tDefinitions.DecimalPlaces, tDefinitions.Currency,
		tDefinitions.GroupKey, tDefinitions.SortOrder, tDefinitions.Pinned,
		tDefinitions.Renderer, tDefinitions.RendererConfig,
		tDefinitions.DisplayMode,
		tDefinitions.Implicit,
	).VALUES(
		scopeName, definition.GetKey(), definition.GetValueType(), definition.GetLabel(), definition.Description,
		definition.GetEnabled(),
		definition.GetSearchable(), definition.GetValueFormat(), definition.GetBehavior(),
		dbutils.StringPP(definition.Unit),
		dbutils.Uint32PP(definition.DecimalPlaces), dbutils.StringPP(definition.Currency),
		definition.GetGroupKey(), definition.GetSortOrder(), definition.GetPinned(),
		definition.GetRenderer(), rendererConfig,
		definition.GetDisplayMode(),
		definition.GetImplicit(),
	).
		ExecContext(ctx, s.db)
	if err != nil {
		return nil, err
	}
	s.propertyCache.Invalidate(definition.GetScope())
	return customPropertyDefinitionCopy(definition), nil
}

func (s *Store) UpdatePropertyDefinition(
	ctx context.Context,
	definition *properties.PropertyDefinition,
) (*properties.PropertyDefinition, error) {
	definition, scopeName, err := s.validateCustomPropertyDefinition(ctx, definition)
	if err != nil {
		return nil, err
	}
	rendererConfig, err := rendererConfigJSON(definition.GetRendererConfig())
	if err != nil {
		return nil, err
	}

	tDefinitions := table.FivenetPropertyDefinitions
	result, err := tDefinitions.UPDATE().SET(
		tDefinitions.Label.SET(mysql.String(definition.GetLabel())),
		tDefinitions.Description.SET(dbutils.StringPP(definition.Description)),
		tDefinitions.Enabled.SET(mysql.Bool(definition.GetEnabled())),
		tDefinitions.Searchable.SET(mysql.Bool(definition.GetSearchable())),
		tDefinitions.ValueFormat.SET(mysql.Int32(int32(definition.GetValueFormat()))),
		tDefinitions.Behavior.SET(mysql.Int32(int32(definition.GetBehavior()))),
		tDefinitions.Unit.SET(dbutils.StringPP(definition.Unit)),
		tDefinitions.DecimalPlaces.SET(dbutils.Uint32PP(definition.DecimalPlaces)),
		tDefinitions.Currency.SET(dbutils.StringPP(definition.Currency)),
		tDefinitions.GroupKey.SET(mysql.String(definition.GetGroupKey())),
		tDefinitions.SortOrder.SET(mysql.Int32(definition.GetSortOrder())),
		tDefinitions.Pinned.SET(mysql.Bool(definition.GetPinned())),
		tDefinitions.Renderer.SET(mysql.Int32(int32(definition.GetRenderer()))),
		tDefinitions.RendererConfig.SET(rendererConfigValue(rendererConfig)),
		tDefinitions.DisplayMode.SET(mysql.Int32(int32(definition.GetDisplayMode()))),
		tDefinitions.Implicit.SET(mysql.Bool(definition.GetImplicit())),
	).WHERE(
		tDefinitions.Scope.EQ(mysql.String(scopeName)).AND(
			tDefinitions.PropertyKey.EQ(mysql.String(definition.GetKey())).AND(
				tDefinitions.ValueType.EQ(mysql.Int32(int32(definition.GetValueType()))).AND(
					tDefinitions.DeletedAt.IS_NULL(),
				),
			),
		),
	).ExecContext(ctx, s.db)
	if err != nil {
		return nil, err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if changed == 0 {
		return nil, ErrPropertyDefinitionNotFound
	}
	s.propertyCache.Invalidate(definition.GetScope())
	return customPropertyDefinitionCopy(definition), nil
}

func (s *Store) DeletePropertyDefinition(
	ctx context.Context,
	scope properties.PropertyScope,
	key string,
) error {
	scopeName, ok := propertyScopeName(scope)
	if !ok || !isCustomPropertyKey(scope, key) {
		return ErrInvalidPropertyDefinition
	}

	tDefinitions := table.FivenetPropertyDefinitions
	result, err := tDefinitions.UPDATE().SET(
		tDefinitions.DeletedAt.SET(mysql.CURRENT_TIMESTAMP()),
	).WHERE(
		tDefinitions.Scope.EQ(mysql.String(scopeName)).AND(
			tDefinitions.PropertyKey.EQ(mysql.String(key)).AND(
				tDefinitions.DeletedAt.IS_NULL(),
			),
		),
	).
		ExecContext(ctx, s.db)
	if err != nil {
		return err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if changed == 0 {
		return ErrPropertyDefinitionNotFound
	}
	s.propertyCache.Invalidate(scope)
	return nil
}

func (s *Store) RestorePropertyDefinition(
	ctx context.Context,
	scope properties.PropertyScope,
	key string,
) (*properties.PropertyDefinition, error) {
	scopeName, ok := propertyScopeName(scope)
	if !ok || !isCustomPropertyKey(scope, key) {
		return nil, ErrInvalidPropertyDefinition
	}

	tDefinitions := table.FivenetPropertyDefinitions
	result, err := tDefinitions.UPDATE().SET(
		tDefinitions.DeletedAt.SET(mysql.TimestampExp(mysql.NULL)),
	).WHERE(
		tDefinitions.Scope.EQ(mysql.String(scopeName)).AND(
			tDefinitions.PropertyKey.EQ(mysql.String(key)).AND(
				tDefinitions.DeletedAt.IS_NOT_NULL(),
			),
		),
	).ExecContext(ctx, s.db)
	if err != nil {
		return nil, err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if changed == 0 {
		return nil, ErrPropertyDefinitionNotFound
	}
	s.propertyCache.Invalidate(scope)

	definitions, err := s.ListPropertyDefinitions(ctx, scope, true)
	if err != nil {
		return nil, err
	}
	for _, definition := range definitions {
		if definition.GetKey() == key {
			return definition, nil
		}
	}
	return nil, ErrPropertyDefinitionNotFound
}

func (s *Store) SetSystemPropertyEnabled(
	ctx context.Context,
	scope properties.PropertyScope,
	key string,
	enabled bool,
) (*properties.PropertyDefinition, error) {
	scopeName, ok := propertyScopeName(scope)
	definition, found := systemPropertyDefinition(scope, key)
	if !ok || !found {
		return nil, ErrInvalidPropertyDefinition
	}

	tOverrides := table.FivenetPropertySystemOverrides
	_, err := tOverrides.
		INSERT(tOverrides.Scope, tOverrides.PropertyKey, tOverrides.Enabled).
		VALUES(scopeName, key, enabled).
		ON_DUPLICATE_KEY_UPDATE(tOverrides.Enabled.SET(mysql.Bool(enabled))).
		ExecContext(ctx, s.db)
	if err != nil {
		return nil, err
	}
	s.propertyCache.Invalidate(scope)
	defCopy := proto.Clone(definition).(*properties.PropertyDefinition)
	defCopy.Enabled = enabled
	return defCopy, nil
}

type propertySystemOverrides struct {
	enabled    map[string]bool
	sortOrders map[string]int32
}

func (s *Store) systemPropertyOverrides(
	ctx context.Context,
	scopeName string,
) (propertySystemOverrides, error) {
	tOverrides := table.FivenetPropertySystemOverrides
	stmt := tOverrides.
		SELECT(tOverrides.PropertyKey, tOverrides.Enabled, tOverrides.SortOrder).
		FROM(tOverrides).
		WHERE(tOverrides.Scope.EQ(mysql.String(scopeName)))
	rows := []model.FivenetPropertySystemOverrides{}
	err := stmt.QueryContext(ctx, s.db, &rows)
	if err != nil {
		return propertySystemOverrides{}, err
	}

	overrides := propertySystemOverrides{
		enabled:    make(map[string]bool),
		sortOrders: make(map[string]int32),
	}
	for _, row := range rows {
		overrides.enabled[row.PropertyKey] = row.Enabled
		if row.SortOrder != nil {
			overrides.sortOrders[row.PropertyKey] = *row.SortOrder
		}
	}
	return overrides, nil
}

// ReconcileSystemPropertyOverrides removes overrides for system properties
// that no longer exist in the generated registry. Overrides for existing
// properties are intentionally preserved, including when their code default
// changes.
func (s *Store) ReconcileSystemPropertyOverrides(ctx context.Context) error {
	valid := make(map[string]struct{})
	for _, definition := range propertyregistry.SystemPropertyDefinitions() {
		scopeName, ok := propertyregistry.ScopeName(definition.GetScope())
		if !ok {
			continue
		}
		valid[scopeName+"\x00"+definition.GetKey()] = struct{}{}
	}

	tOverrides := table.FivenetPropertySystemOverrides
	rows := []model.FivenetPropertySystemOverrides{}
	stmt := tOverrides.SELECT(tOverrides.Scope, tOverrides.PropertyKey).
		FROM(tOverrides)
	if err := stmt.QueryContext(ctx, s.db, &rows); err != nil {
		return err
	}

	for _, row := range rows {
		if _, ok := valid[row.Scope+"\x00"+row.PropertyKey]; ok {
			continue
		}
		if _, err := tOverrides.DELETE().WHERE(
			tOverrides.Scope.EQ(mysql.String(row.Scope)).AND(
				tOverrides.PropertyKey.EQ(mysql.String(row.PropertyKey)),
			),
		).ExecContext(ctx, s.db); err != nil {
			return err
		}
	}

	validGroups := make(map[string]struct{})
	for _, group := range propertyregistry.SystemPropertyGroups() {
		if scopeName, ok := propertyregistry.ScopeName(group.GetScope()); ok {
			validGroups[scopeName+"\x00"+group.GetKey()] = struct{}{}
		}
	}

	tGroupOverrides := table.FivenetPropertyGroupOverrides
	groupRows := []model.FivenetPropertyGroupOverrides{}
	if err := tGroupOverrides.SELECT(
		tGroupOverrides.Scope,
		tGroupOverrides.GroupKey,
	).FROM(tGroupOverrides).QueryContext(ctx, s.db, &groupRows); err != nil {
		return err
	}
	for _, row := range groupRows {
		if _, ok := validGroups[row.Scope+"\x00"+row.GroupKey]; ok {
			continue
		}
		if _, err := tGroupOverrides.DELETE().WHERE(
			tGroupOverrides.Scope.EQ(mysql.String(row.Scope)).AND(
				tGroupOverrides.GroupKey.EQ(mysql.String(row.GroupKey)),
			),
		).ExecContext(ctx, s.db); err != nil {
			return err
		}
	}

	s.propertyCache.Invalidate(properties.PropertyScope_PROPERTY_SCOPE_USER)
	s.propertyCache.Invalidate(properties.PropertyScope_PROPERTY_SCOPE_VEHICLE)
	return nil
}

func propertyScopeName(scope properties.PropertyScope) (string, bool) {
	return propertyregistry.ScopeName(scope)
}

func systemPropertyDefinition(
	scope properties.PropertyScope,
	key string,
) (*properties.PropertyDefinition, bool) {
	return propertyregistry.SystemPropertyDefinition(scope, key)
}

func (s *Store) validateCustomPropertyDefinition(
	ctx context.Context,
	definition *properties.PropertyDefinition,
) (*properties.PropertyDefinition, string, error) {
	normalized, err := propertyregistry.NormalizeAndValidateCustomDefinition(definition)
	if err != nil {
		return nil, "", ErrInvalidPropertyDefinition
	}
	scopeName, ok := propertyregistry.ScopeName(normalized.GetScope())
	if !ok {
		return nil, "", ErrInvalidPropertyDefinition
	}
	if normalized.GetGroupKey() != "" {
		if _, ok := propertyregistry.SystemPropertyGroup(
			normalized.GetScope(),
			normalized.GetGroupKey(),
		); !ok {
			groups, err := s.ListPropertyGroups(ctx, normalized.GetScope(), false)
			if err != nil {
				return nil, "", err
			}
			found := false
			for _, group := range groups {
				if group.GetKey() == normalized.GetGroupKey() {
					found = true
					break
				}
			}
			if !found {
				return nil, "", ErrInvalidPropertyDefinition
			}
		}
	}
	return normalized, scopeName, nil
}

func isCustomPropertyKey(scope properties.PropertyScope, key string) bool {
	return propertyregistry.IsCustomPropertyKey(scope, key)
}

func customPropertyDefinitionCopy(
	definition *properties.PropertyDefinition,
) *properties.PropertyDefinition {
	copyDef := proto.Clone(definition).(*properties.PropertyDefinition)
	copyDef.Source = properties.PropertyDefinitionSource_PROPERTY_DEFINITION_SOURCE_CUSTOM
	copyDef.DisplayKey = ""
	return copyDef
}

func rendererConfigJSON(config *properties.PropertyRendererConfig) (*string, error) {
	if config == nil {
		return nil, nil
	}
	b, err := protojson.Marshal(config)
	if err != nil {
		return nil, err
	}
	value := string(b)
	return &value, nil
}

func rendererConfigValue(value *string) mysql.StringExpression {
	if value == nil {
		return mysql.StringExp(mysql.NULL)
	}
	return mysql.String(*value)
}
