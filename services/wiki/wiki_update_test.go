package wiki

import (
	"testing"

	resourcesaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/content"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/file"
	reswiki "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestPageUpdateChanges(t *testing.T) {
	t.Parallel()

	oldPage := &reswiki.Page{
		ParentId: func() *int64 { id := int64(10); return &id }(),
		Meta: &reswiki.PageMeta{
			Title:       "Original page",
			Description: "Description",
			Toc:         func() *bool { value := true; return &value }(),
		},
		Content: &content.Content{ContentType: content.ContentType_CONTENT_TYPE_HTML},
		Access: &resourcesaccess.Access{
			Users: []*resourcesaccess.UserAccess{{UserId: 1, Access: 3}},
		},
		Files: []*file.File{{Id: 1}},
	}

	tests := []struct {
		name          string
		mutate        func(*reswiki.Page)
		contentChange bool
	}{
		{
			name: "access only",
			mutate: func(page *reswiki.Page) {
				page.Access.Users[0].Access = 4
			},
		},
		{
			name: "content",
			mutate: func(page *reswiki.Page) {
				page.Content.ContentType = content.ContentType_CONTENT_TYPE_TIPTAP_JSON
			},
			contentChange: true,
		},
		{
			name: "page details",
			mutate: func(page *reswiki.Page) {
				page.Meta.Title = "Changed page"
			},
			contentChange: true,
		},
		{
			name: "files",
			mutate: func(page *reswiki.Page) {
				page.Files = append(page.Files, &file.File{Id: 2})
			},
			contentChange: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			newPage := proto.Clone(oldPage).(*reswiki.Page)
			test.mutate(newPage)

			assert.Equal(t, test.contentChange, pageUpdateContentChanged(oldPage, newPage))
		})
	}
}
