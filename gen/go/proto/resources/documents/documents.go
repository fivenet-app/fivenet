package documents

import documentscategory "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/category"

func (x *DocumentShort) SetCategory(cat *documentscategory.Category) {
	x.Category = cat
}

func (x *Document) SetCategory(cat *documentscategory.Category) {
	x.Category = cat
}
