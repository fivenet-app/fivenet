package wiki

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki"
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

	seen := make(map[int64]struct{}, len(rootPages))
	result := make([]*wiki.PageShort, 0, len(rootPages))
	for _, root := range rootPages {
		if root == nil || root.GetId() <= 0 {
			continue
		}
		if _, ok := seen[root.GetId()]; ok {
			continue
		}
		seen[root.GetId()] = struct{}{}
		result = append(result, root)
	}

	return result
}
