package wiki

import "github.com/fivenet-app/fivenet/gen/go/proto/resources/wiki"

func mapPagesToNavItems(pages []*wiki.PageShort) map[uint64]*wiki.PageShort {
	if len(pages) == 0 {
		return nil
	}

	var root *wiki.PageShort
	mapping := map[uint64]*wiki.PageShort{}
	for _, page := range pages {
		mapping[page.Id] = page

		if page.ParentId != nil {
			_, ok := mapping[*page.ParentId]
			if !ok {
				continue
			}
			mapping[*page.ParentId].Children = append(mapping[*page.ParentId].Children, page)
		} else {
			root = page
		}
	}

	if root == nil {
		root = pages[0]
	}

	for _, page := range mapping {
		if page.ParentId != nil {
			// Only delete pages for which we have the parent,
			// cause it might be a singular page from a "different" wiki
			if _, ok := mapping[*page.ParentId]; ok {
				delete(mapping, page.Id)
			}
		}

		if len(page.Children) > 0 && page.Id != root.Id {
			// Duplicate page to have an entry without children as a clone
			page.Children = append([]*wiki.PageShort{{
				Id:          page.Id,
				ParentId:    page.ParentId,
				Job:         page.Job,
				JobLabel:    page.JobLabel,
				Slug:        page.Slug,
				Title:       page.Title,
				Description: page.Description,
				Children:    nil,
			}}, page.Children...)
		}
	}

	rootTitle := root.Title
	if root.JobLabel != nil {
		rootTitle = *root.JobLabel + ": " + root.Title
	}
	result := map[uint64]*wiki.PageShort{
		root.Id: {
			Id:          root.Id,
			Job:         root.Job,
			JobLabel:    root.JobLabel,
			Slug:        root.Slug,
			Title:       rootTitle,
			Description: root.Description,
			// Make sure to prepend root page
			Children: append([]*wiki.PageShort{
				{
					Id:          root.Id,
					ParentId:    root.ParentId,
					Job:         root.Job,
					JobLabel:    root.JobLabel,
					Slug:        root.Slug,
					Title:       root.Title,
					Description: root.Description,
					Children:    nil,
				},
			}, root.Children...),
		},
	}

	return result
}
