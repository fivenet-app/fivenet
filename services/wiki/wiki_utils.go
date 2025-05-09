package wiki

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/wiki"
)

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

	// If no root page is found, use "dummy"
	if root == nil {
		firstPage := pages[0]
		root = &wiki.PageShort{
			Id:       0,
			Title:    "",
			Job:      firstPage.Job,
			JobLabel: firstPage.JobLabel,
			Children: []*wiki.PageShort{},
		}
	}

	for _, page := range mapping {
		if page.ParentId != nil {
			// Only delete pages for which we have the parent,
			// cause it might be a singular page from a "different" wiki
			// that we add to the root page
			if _, ok := mapping[*page.ParentId]; !ok {
				root.Children = append(root.Children, page)
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
	if root.JobLabel != nil && rootTitle != "" {
		rootTitle = *root.JobLabel + ": " + rootTitle
	}

	// Make sure to prepend root page (if it isn't our dummy)
	if root.Id != 0 {
		root.Children = append([]*wiki.PageShort{
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
		}, root.Children...)
	}

	return map[uint64]*wiki.PageShort{
		root.Id: {
			Id:          root.Id,
			Job:         root.Job,
			JobLabel:    root.JobLabel,
			Slug:        root.Slug,
			Title:       rootTitle,
			Description: root.Description,
			Children:    root.Children,
		},
	}
}
