package wiki

import "github.com/fivenet-app/fivenet/gen/go/proto/resources/wiki"

func mapPagesToNavItems(pages []*wiki.PageShort) map[uint64]*wiki.PageShort {
	mapping := map[uint64]*wiki.PageShort{}
	for _, page := range pages {
		mapping[page.Id] = page

		if page.ParentId != nil {
			_, ok := mapping[*page.ParentId]
			if !ok {
				continue
			}
			mapping[*page.ParentId].Children = append(mapping[*page.ParentId].Children, page)
		}
	}

	for key, page := range mapping {
		if page.ParentId != nil {
			delete(mapping, key)
		}

		if len(page.Children) > 0 {
			// Duplicate page to have an entry without children
			page.Children = append([]*wiki.PageShort{{
				Id:          page.Id,
				ParentId:    page.ParentId,
				Job:         page.Job,
				Slug:        page.Slug,
				Title:       page.Title,
				Description: page.Description,
				Children:    nil,
			}}, page.Children...)
		}
	}

	return mapping
}
