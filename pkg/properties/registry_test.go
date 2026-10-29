package properties

import (
	"testing"

	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
)

func TestSystemPropertyDefinitionsAreDefensiveAndIdempotent(t *testing.T) {
	t.Parallel()

	const key = "test.registry"
	definition := &propertypb.PropertyDefinition{
		Scope:     propertypb.PropertyScope_PROPERTY_SCOPE_COLLEAGUE,
		Key:       key,
		ValueType: propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_STRING,
		Enabled:   true,
	}

	AddSystemPropertyDefinitions([]*propertypb.PropertyDefinition{definition})
	definition.Enabled = false

	definitions := SystemPropertyDefinitions()
	var got *propertypb.PropertyDefinition
	for _, candidate := range definitions {
		if candidate.GetKey() == key {
			got = candidate
			break
		}
	}
	if got == nil || !got.GetEnabled() {
		t.Fatalf("expected registered definition to be cloned on input")
	}

	got.Enabled = false
	for _, candidate := range SystemPropertyDefinitions() {
		if candidate.GetKey() == key && !candidate.GetEnabled() {
			t.Fatalf("expected returned definitions to be defensive copies")
		}
	}

	updated := &propertypb.PropertyDefinition{
		Scope:     propertypb.PropertyScope_PROPERTY_SCOPE_COLLEAGUE,
		Key:       key,
		ValueType: propertypb.PropertyValueType_PROPERTY_VALUE_TYPE_INTEGER,
	}
	AddSystemPropertyDefinitions([]*propertypb.PropertyDefinition{updated})

	count := 0
	for _, candidate := range SystemPropertyDefinitions() {
		if candidate.GetScope() == updated.GetScope() && candidate.GetKey() == key {
			count++
			if candidate.GetValueType() != updated.GetValueType() {
				t.Fatalf("expected duplicate registration to update definition")
			}
		}
	}
	if count != 1 {
		t.Fatalf("expected one definition after duplicate registration, got %d", count)
	}
}

func TestSystemPropertyGroupsAreDefensiveAndSorted(t *testing.T) {
	t.Parallel()

	AddSystemPropertyGroups([]*propertypb.PropertyGroup{
		{Scope: propertypb.PropertyScope_PROPERTY_SCOPE_COLLEAGUE, Key: "late", SortOrder: 200},
		{Scope: propertypb.PropertyScope_PROPERTY_SCOPE_COLLEAGUE, Key: "early", SortOrder: 100},
	})

	groups := SystemPropertyGroups()
	var colleague []*propertypb.PropertyGroup
	for _, group := range groups {
		if group.GetScope() == propertypb.PropertyScope_PROPERTY_SCOPE_COLLEAGUE {
			colleague = append(colleague, group)
		}
	}
	if len(colleague) < 2 || colleague[len(colleague)-2].GetKey() != "early" ||
		colleague[len(colleague)-1].GetKey() != "late" {
		t.Fatalf("expected colleague groups sorted by order, got %v", colleague)
	}

	colleague[0].Key = "changed"
	group, ok := SystemPropertyGroup(propertypb.PropertyScope_PROPERTY_SCOPE_COLLEAGUE, "early")
	if !ok || group.GetKey() != "early" {
		t.Fatalf("expected defensive group copy")
	}
}
