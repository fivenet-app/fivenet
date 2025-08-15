package wiki

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/wiki"
)

func mapPagesToNavItems(pages []*wiki.PageShort) map[int64]*wiki.PageShort {
	if len(pages) == 0 {
		return nil
	}

	var root *wiki.PageShort
	mapping := map[int64]*wiki.PageShort{}
	for _, page := range pages {
		mapping[page.GetId()] = page

		if page.ParentId != nil {
			_, ok := mapping[page.GetParentId()]
			if !ok {
				continue
			}
			mapping[page.GetParentId()].Children = append(
				mapping[page.GetParentId()].Children,
				page,
			)
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
			Job:      firstPage.GetJob(),
			JobLabel: firstPage.JobLabel,
			Children: []*wiki.PageShort{},
		}
	}

	for _, page := range mapping {
		if page.ParentId != nil {
			// Only delete pages for which we have the parent,
			// cause it might be a singular page from a "different" wiki
			// that we add to the root page
			if _, ok := mapping[page.GetParentId()]; !ok {
				root.Children = append(root.Children, page)
			}
		}

		if len(page.GetChildren()) > 0 && page.GetId() != root.GetId() {
			// Duplicate page to have an entry without children as a clone
			page.Children = append([]*wiki.PageShort{{
				Id:          page.GetId(),
				ParentId:    page.ParentId,
				Job:         page.GetJob(),
				JobLabel:    page.JobLabel,
				Slug:        page.Slug,
				Title:       page.GetTitle(),
				Description: page.GetDescription(),
				Children:    nil,
			}}, page.GetChildren()...)
		}
	}

	rootTitle := root.GetTitle()
	if root.JobLabel != nil && rootTitle != "" {
		rootTitle = root.GetJobLabel() + ": " + rootTitle
	}

	// Make sure to prepend root page (if it isn't our dummy)
	if root.GetId() != 0 {
		root.Children = append([]*wiki.PageShort{
			{
				Id:          root.GetId(),
				ParentId:    root.ParentId,
				Job:         root.GetJob(),
				JobLabel:    root.JobLabel,
				Slug:        root.Slug,
				Title:       root.GetTitle(),
				Description: root.GetDescription(),
				Children:    nil,
			},
		}, root.GetChildren()...)
	}

	return map[int64]*wiki.PageShort{
		root.GetId(): {
			Id:          root.GetId(),
			Job:         root.GetJob(),
			JobLabel:    root.JobLabel,
			Slug:        root.Slug,
			Title:       rootTitle,
			Description: root.GetDescription(),
			Children:    root.GetChildren(),
		},
	}
}
