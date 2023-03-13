package htmlsanitizer

import "github.com/microcosm-cc/bluemonday"

var p *bluemonday.Policy

func init() {
	p = bluemonday.UGCPolicy()
}

func Sanitize(in string) string {
	return p.Sanitize(in)
}
