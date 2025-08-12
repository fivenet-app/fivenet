package perms

import (
	"github.com/gosimple/slug"
)

// Guard makes the given string into a safe to use guard slug.
// example: 'create $#% contact' → 'create-contact'.
//
//	@param	string
//	@param	string
func Guard(b string) string {
	return slug.Make(b)
}
