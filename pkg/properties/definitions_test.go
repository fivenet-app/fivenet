package properties

import (
	"testing"

	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSortPropertiesByDefinitionSortOrder(t *testing.T) {
	t.Parallel()

	values := []*propertypb.Property{
		{Key: "user.second", Definition: &propertypb.PropertyDefinition{SortOrder: 200}},
		{Key: "user.first", Definition: &propertypb.PropertyDefinition{SortOrder: 100}},
		{Key: "user.unknown"},
	}

	sorted := SortProperties(values)

	assert.Equal(t, []string{"user.first", "user.second", "user.unknown"}, []string{
		sorted[0].GetKey(), sorted[1].GetKey(), sorted[2].GetKey(),
	})
	assert.Equal(t, "user.second", values[0].GetKey())
}

func TestSortPropertiesByGroupThenDefinitionOrder(t *testing.T) {
	t.Parallel()

	values := []*propertypb.Property{
		{
			Key:        "user.later",
			Definition: &propertypb.PropertyDefinition{GroupSortOrder: 200, SortOrder: 10},
		},
		{
			Key:        "user.earlier",
			Definition: &propertypb.PropertyDefinition{GroupSortOrder: 100, SortOrder: 200},
		},
		{
			Key:        "user.earlier_first",
			Definition: &propertypb.PropertyDefinition{GroupSortOrder: 100, SortOrder: 100},
		},
	}

	sorted := SortProperties(values)
	assert.Equal(t, []string{"user.earlier_first", "user.earlier", "user.later"}, []string{
		sorted[0].GetKey(), sorted[1].GetKey(), sorted[2].GetKey(),
	})
}

func TestSortDefinitionsByEffectiveGroupSortOrder(t *testing.T) {
	t.Parallel()

	definitions := []*propertypb.PropertyDefinition{
		{Key: "user.later", GroupKey: "later", GroupSortOrder: 200, SortOrder: 10},
		{Key: "user.earlier", GroupKey: "earlier", GroupSortOrder: 100, SortOrder: 10},
	}

	sorted := SortDefinitions(definitions)

	assert.Equal(t, []string{"user.earlier", "user.later"}, []string{
		sorted[0].GetKey(), sorted[1].GetKey(),
	})
}

func TestResolvePropertiesAttachesDefinitionsAndSorts(t *testing.T) {
	t.Parallel()

	values := []*propertypb.Property{{Key: "user.second"}, {Key: "user.first"}}
	definitions := map[string]*propertypb.PropertyDefinition{
		"user.second": {SortOrder: 200, Enabled: true},
		"user.first":  {SortOrder: 100, Enabled: true},
	}

	resolved := ResolveProperties(values, definitions)

	assert.Same(t, definitions["user.first"], resolved[0].GetDefinition())
	assert.Equal(t, "user.first", resolved[0].GetKey())
	assert.Same(t, definitions["user.second"], resolved[1].GetDefinition())
}

func TestResolvePropertiesOmitsDisabledDefinitions(t *testing.T) {
	t.Parallel()

	values := []*propertypb.Property{{Key: "user.enabled"}, {Key: "user.disabled"}}
	definitions := map[string]*propertypb.PropertyDefinition{
		"user.enabled":  {Enabled: true, SortOrder: 100},
		"user.disabled": {Enabled: false, SortOrder: 200},
	}

	resolved := ResolveProperties(values, definitions)

	require.Len(t, resolved, 1)
	assert.Equal(t, "user.enabled", resolved[0].GetKey())
}

func TestResolvePropertiesIncludesEmptyCustomDefinitions(t *testing.T) {
	t.Parallel()

	definitions := map[string]*propertypb.PropertyDefinition{
		"vehicle.color": {
			Key:     "vehicle.color",
			Enabled: true,
			Source:  propertypb.PropertyDefinitionSource_PROPERTY_DEFINITION_SOURCE_CUSTOM,
		},
		"vehicle.wanted": {
			Key:     "vehicle.wanted",
			Enabled: true,
			Source:  propertypb.PropertyDefinitionSource_PROPERTY_DEFINITION_SOURCE_SYSTEM,
		},
	}

	resolved := ResolveProperties(nil, definitions)

	require.Len(t, resolved, 1)
	assert.Equal(t, "vehicle.color", resolved[0].GetKey())
	assert.Nil(t, resolved[0].GetValue())
	assert.Same(t, definitions["vehicle.color"], resolved[0].GetDefinition())
}

func TestResolvePropertiesResourceIncludesValuesAndGroups(t *testing.T) {
	t.Parallel()

	definitions := map[string]*propertypb.PropertyDefinition{
		"user.email": {
			Scope:          propertypb.PropertyScope_PROPERTY_SCOPE_USER,
			GroupKey:       "contact",
			GroupSortOrder: 200,
			SortOrder:      100,
			Enabled:        true,
		},
	}
	resource := ResolvePropertiesResource([]*propertypb.Property{{Key: "user.email"}}, definitions)

	require.Len(t, resource.GetProperties(), 1)
	assert.Same(t, definitions["user.email"], resource.GetProperties()[0].GetDefinition())
	require.Len(t, resource.GetGroups(), 1)
	assert.Equal(t, "contact", resource.GetGroups()[0].GetKey())
}

func TestResolvePropertyGroupsSortsAndDeduplicates(t *testing.T) {
	t.Parallel()

	groups := ResolvePropertyGroups(map[string]*propertypb.PropertyDefinition{
		"user.contact_email": {
			Scope:          propertypb.PropertyScope_PROPERTY_SCOPE_USER,
			GroupKey:       "contact",
			GroupSortOrder: 200,
			Enabled:        true,
		},
		"user.contact_phone": {
			Scope:          propertypb.PropertyScope_PROPERTY_SCOPE_USER,
			GroupKey:       "contact",
			GroupSortOrder: 200,
			Enabled:        true,
		},
		"user.dateofbirth": {
			Scope:          propertypb.PropertyScope_PROPERTY_SCOPE_USER,
			GroupKey:       "personal",
			GroupSortOrder: 100,
			Enabled:        true,
		},
		"user.ungrouped": {
			Scope:   propertypb.PropertyScope_PROPERTY_SCOPE_USER,
			Enabled: true,
		},
		"user.disabled": {
			Scope:          propertypb.PropertyScope_PROPERTY_SCOPE_USER,
			GroupKey:       "status",
			GroupSortOrder: 300,
			Enabled:        false,
		},
	})

	require.Len(t, groups, 2)
	assert.Equal(t, "personal", groups[0].GetKey())
	assert.Equal(t, int32(100), groups[0].GetSortOrder())
	assert.Equal(t, "contact", groups[1].GetKey())
	assert.Equal(t, int32(200), groups[1].GetSortOrder())
}

func TestResolvePropertyGroupsOnlyIncludesResolvedDefinitions(t *testing.T) {
	t.Parallel()

	groups := ResolvePropertyGroups(map[string]*propertypb.PropertyDefinition{
		"user.email": {
			Scope:          propertypb.PropertyScope_PROPERTY_SCOPE_USER,
			GroupKey:       "contact",
			GroupSortOrder: 200,
			Enabled:        true,
		},
	})

	require.Len(t, groups, 1)
	assert.Equal(t, "contact", groups[0].GetKey())
}

func TestValidateCustomDefinitionAcceptsTypedRendererConfig(t *testing.T) {
	t.Parallel()

	definition := &propertypb.PropertyDefinition{
		Scope:       propertypb.PropertyScope_PROPERTY_SCOPE_USER,
		Key:         "user.height",
		Source:      propertypb.PropertyDefinitionSource_PROPERTY_DEFINITION_SOURCE_CUSTOM,
		ValueType:   propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_NUMBER,
		ValueFormat: propertypb.PropertyValueFormat_PROPERTY_VALUE_FORMAT_PLAIN,
		Behavior:    propertypb.PropertyBehavior_PROPERTY_BEHAVIOR_GENERIC,
		Label:       "Height",
		Renderer:    propertypb.PropertyRenderer_PROPERTY_RENDERER_NUMBER,
		RendererConfig: &propertypb.PropertyRendererConfig{
			Config: &propertypb.PropertyRendererConfig_Number{
				Number: &propertypb.PropertyNumberRendererConfig{Unit: stringPointer("cm")},
			},
		},
	}

	require.NoError(t, ValidateCustomDefinition(definition))
}

func TestValidateCustomDefinitionRejectsMismatchedRendererConfig(t *testing.T) {
	t.Parallel()

	definition := &propertypb.PropertyDefinition{
		Scope:       propertypb.PropertyScope_PROPERTY_SCOPE_USER,
		Key:         "user.height",
		Source:      propertypb.PropertyDefinitionSource_PROPERTY_DEFINITION_SOURCE_CUSTOM,
		ValueType:   propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_NUMBER,
		ValueFormat: propertypb.PropertyValueFormat_PROPERTY_VALUE_FORMAT_PLAIN,
		Behavior:    propertypb.PropertyBehavior_PROPERTY_BEHAVIOR_GENERIC,
		Label:       "Height",
		Renderer:    propertypb.PropertyRenderer_PROPERTY_RENDERER_NUMBER,
		RendererConfig: &propertypb.PropertyRendererConfig{
			Config: &propertypb.PropertyRendererConfig_String_{
				String_: &propertypb.PropertyStringRendererConfig{Suffix: "cm"},
			},
		},
	}

	require.ErrorIs(t, ValidateCustomDefinition(definition), ErrInvalidPropertyDefinition)
}

func TestNormalizeDefinitionReturnsCopy(t *testing.T) {
	t.Parallel()

	definition := &propertypb.PropertyDefinition{
		Scope:     propertypb.PropertyScope_PROPERTY_SCOPE_VEHICLE,
		Key:       "vehicle.example",
		ValueType: propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_STRING,
		Source:    propertypb.PropertyDefinitionSource_PROPERTY_DEFINITION_SOURCE_CUSTOM,
		Label:     "Example",
	}

	normalized := NormalizeDefinition(definition)

	assert.NotSame(t, definition, normalized)
	assert.Equal(
		t,
		propertypb.PropertyValueFormat_PROPERTY_VALUE_FORMAT_UNSPECIFIED,
		definition.GetValueFormat(),
	)
	assert.Equal(
		t,
		propertypb.PropertyValueFormat_PROPERTY_VALUE_FORMAT_PLAIN,
		normalized.GetValueFormat(),
	)
	assert.Equal(t, propertypb.PropertyBehavior_PROPERTY_BEHAVIOR_GENERIC, normalized.GetBehavior())
	assert.Equal(t, propertypb.PropertyRenderer_PROPERTY_RENDERER_DEFAULT, normalized.GetRenderer())
	assert.Equal(
		t,
		propertypb.PropertyDisplayMode_PROPERTY_DISPLAY_MODE_DEFAULT,
		normalized.GetDisplayMode(),
	)
}

func stringPointer(value string) *string {
	return &value
}
