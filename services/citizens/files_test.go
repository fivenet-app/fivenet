package citizens

import (
	"testing"

	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	"github.com/stretchr/testify/assert"
)

func TestCitizenFileKindsFollowListProperties(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		keys []string
		want []string
	}{
		{name: "attachments only", keys: []string{"user.files"}, want: []string{"attachment"}},
		{name: "mugshots only", keys: []string{"user.mugshot"}, want: []string{"mugshot"}},
		{
			name: "both",
			keys: []string{"user.files", "user.mugshot"},
			want: []string{"attachment", "mugshot"},
		},
		{name: "neither", keys: []string{"user.labels"}, want: []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				tt.want,
				citizenFileKinds(&permissionsattributes.PropertyList{Keys: tt.keys}),
			)
		})
	}
}
