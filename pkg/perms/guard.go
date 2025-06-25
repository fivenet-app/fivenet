package perms

import (
	"github.com/gosimple/slug"
)

// Guard edits the given string.
// example: 'create $#% contact' → 'create-contact'.
// @param string
// @param string
// return bool
func Guard(b string) string {
	return slug.Make(b)
}
