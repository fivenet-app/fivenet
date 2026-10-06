package perms

import (
	"testing"

	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuperuserAttributesUsesValidValues(t *testing.T) {
	attrs := []*permissionsattributes.RoleAttribute{
		{
			Namespace: "citizens",
			ValidValues: &permissionsattributes.AttributeValues{
				ValidValues: &permissionsattributes.AttributeValues_StringList{
					StringList: &permissionsattributes.StringList{Strings: []string{"name", "phone"}},
				},
			},
			Value: &permissionsattributes.AttributeValues{
				ValidValues: &permissionsattributes.AttributeValues_StringList{
					StringList: &permissionsattributes.StringList{Strings: []string{"name"}},
				},
			},
		},
	}

	got := superuserAttributes(attrs)

	require.Len(t, got, 1)
	assert.Equal(t, []string{"name", "phone"}, got[0].GetValue().GetStringList().GetStrings())
	assert.Equal(t, []string{"name", "phone"}, got[0].GetValidValues().GetStringList().GetStrings())

	got[0].GetValue().GetStringList().Strings[0] = "mutated"
	assert.Equal(t, []string{"name", "phone"}, attrs[0].GetValidValues().GetStringList().GetStrings())
}
