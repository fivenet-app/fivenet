package perms

import (
	"fmt"
	"strings"
)

// DefaultPermGuard converts an AppConfig permission category/name pair to the
// canonical permission guard format.
func DefaultPermGuard(category string, name string) (string, error) {
	separator := strings.LastIndexByte(category, '.')
	if separator <= 0 || separator == len(category)-1 || name == "" {
		return "", fmt.Errorf("invalid default permission category/name %q/%q", category, name)
	}

	return BuildGuard(
		Namespace(category[:separator]),
		Service(category[separator+1:]),
		Name(name),
	), nil
}
