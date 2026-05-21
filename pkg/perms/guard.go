package perms

import (
	"github.com/gosimple/slug"
)

// guard makes the given string into a safe to use guard slug.
// Example: 'create $#% contact' -> 'create-contact'.
func guard(s string) string {
	return slug.Make(s)
}
