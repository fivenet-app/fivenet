package wiki

import (
	"slices"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/wiki"
)

func mapPagesToNavItems(pages []*wiki.PageShort) []*wiki.PageShort {
	if len(pages) == 0 {
		return nil
	}

	mapping := map[int64]*wiki.PageShort{}
	rootPages := []*wiki.PageShort{}

	for _, page := range pages {
		mapping[page.GetId()] = page
	}

	noRoot := false
	for _, page := range pages {
		if page.GetParentId() > 0 {
			if parent, ok := mapping[page.GetParentId()]; ok {
				parent.Children = append(parent.Children, page)
			}
		} else {
			rootPages = append(rootPages, page)
			noRoot = true
		}
	}

	// Handle orphans (pages whose parent is not in the list)
	for _, page := range mapping {
		if page.ParentId != nil && page.GetParentId() > 0 {
			if _, ok := mapping[page.GetParentId()]; !ok {
				// Attach to all roots (or dummy)
				if !noRoot {
					rootPages = append(rootPages, page)
				} else {
					for _, root := range rootPages {
						root.Children = append(root.Children, page)
					}
				}
			}
		}
	}

	mapped := make(map[int64]*wiki.PageShort, len(rootPages))
	for _, root := range rootPages {
		mapped[root.GetId()] = &wiki.PageShort{
			Id:          root.GetId(),
			Job:         root.GetJob(),
			JobLabel:    root.JobLabel,
			Slug:        root.Slug,
			Title:       root.GetTitle(),
			Description: root.GetDescription(),
			Children:    root.GetChildren(),
		}
	}

	result := []*wiki.PageShort{}
	for _, page := range mapped {
		result = append(result, page)
	}

	slices.SortStableFunc(result, func(a, b *wiki.PageShort) int {
		if a.GetStartpage() == b.GetStartpage() {
			return strings.Compare(a.GetTitle(), b.GetTitle())
		}
		if a.GetStartpage() {
			return -1
		}
		return 1
	})

	return result
}
