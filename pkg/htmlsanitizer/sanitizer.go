package htmlsanitizer

import "github.com/microcosm-cc/bluemonday"

var (
	p               *bluemonday.Policy
	policyStripTags *bluemonday.Policy
)

func init() {
	p = bluemonday.UGCPolicy()
	policyStripTags = bluemonday.StripTagsPolicy()
}

func Sanitize(in string) string {
	return p.Sanitize(in)
}

func StripTags(in string) string {
	return policyStripTags.Sanitize(in)
}
