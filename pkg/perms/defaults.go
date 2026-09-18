package perms

import (
	"fmt"
	"slices"
	"strings"
)

// DefaultPermGuard converts an AppConfig permission category/name pair to the
// canonical permission guard format.
func DefaultPermGuard(category string, name string) (string, error) {
	parts := strings.Split(category, ".")
	if len(parts) < 2 || name == "" {
		return "", fmt.Errorf("invalid default permission category/name %q/%q", category, name)
	}
	if slices.Contains(parts, "") {
		return "", fmt.Errorf("invalid default permission category/name %q/%q", category, name)
	}

	service := parts[len(parts)-1]
	namespace := strings.Join(parts[:len(parts)-1], ".")

	return BuildGuard(
		Namespace(namespace),
		Service(service),
		Name(name),
	), nil
}
