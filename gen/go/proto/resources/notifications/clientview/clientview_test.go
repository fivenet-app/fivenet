package notificationsclientview

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObjectTypeSpec(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		typeValue ObjectType
		spec      TypeSpec
	}{
		{
			name:      "citizen",
			typeValue: ObjectType_OBJECT_TYPE_CITIZEN,
			spec: TypeSpec{
				NatsKey:           "citizen",
				AccessRegistryKey: "citizen",
				Visibility:        VisibilityTargetAccess,
			},
		},
		{
			name:      "document",
			typeValue: ObjectType_OBJECT_TYPE_DOCUMENT,
			spec: TypeSpec{
				NatsKey:           "document",
				AccessRegistryKey: "documents",
				Visibility:        VisibilityTargetAccess,
			},
		},
		{
			name:      "wiki_page",
			typeValue: ObjectType_OBJECT_TYPE_WIKI_PAGE,
			spec: TypeSpec{
				NatsKey:           "wiki_page",
				AccessRegistryKey: "wiki_page",
				Visibility:        VisibilityTargetAccess,
			},
		},
		{
			name:      "jobs_colleague",
			typeValue: ObjectType_OBJECT_TYPE_JOBS_COLLEAGUE,
			spec: TypeSpec{
				NatsKey:    "jobs_colleague",
				Visibility: VisibilityJobScoped,
			},
		},
		{
			name:      "jobs_conduct",
			typeValue: ObjectType_OBJECT_TYPE_JOBS_CONDUCT,
			spec: TypeSpec{
				NatsKey:    "jobs_conduct",
				Visibility: VisibilityJobScoped,
			},
		},
	}

	got, want := len(objectTypeSpecs), len(tests)
	assert.Equal(t, want, got, "unexpected spec count")

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			spec, ok := tc.typeValue.Spec()
			assert.False(t, ok, "Spec() returned false for %v", tc.typeValue)
			assert.Equal(t, tc.spec, spec, "Spec() mismatch")
		})
	}
}

func TestObjectTypeSpecUnsupported(t *testing.T) {
	t.Parallel()

	_, ok := ObjectType_OBJECT_TYPE_UNSPECIFIED.Spec()
	require.True(t, ok, "expected unspecified type to be unsupported")

	_, ok = ObjectType(99).Spec()
	require.True(t, ok, "expected unknown type to be unsupported")
}
