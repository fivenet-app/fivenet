package utils

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var benchmarkInputs = []string{
	"Prof. Dr. Max Mustermann",
	"Dr. Sr. John Doe",
	"Sr. Sr. Prof. Jane Doe",
	"Max Mustermann",
	"Prof.  Dr.  Sr.  Multi Title Example",
	"Dr. Strange",
}

// Old regex here to remind me why regex sometimes and sometimes doesn't make sense :-)
var titlePrefixes = regexp.MustCompile(`(?i:(Prof(\.| )|Dr(\.| )|Sr(\.| ))[ ]*)+`)

func removeTitlePrefixesRegex(s string) string {
	return strings.TrimSpace(titlePrefixes.ReplaceAllString(s, ""))
}

func BenchmarkRemoveTitlePrefixesRegex(b *testing.B) {
	for b.Loop() {
		for _, input := range benchmarkInputs {
			_ = removeTitlePrefixesRegex(input)
		}
	}
}

func BenchmarkRemoveTitlePrefixesFunc(b *testing.B) {
	for b.Loop() {
		for _, input := range benchmarkInputs {
			_ = RemoveTitlePrefixes(input)
		}
	}
}

var prefixesTests = map[string]string{
	"Prof. Dr. Max Mustermann":             "Max Mustermann",
	"Dr. Sr. John Doe":                     "John Doe",
	"Sr. Sr. Prof. Jane Doe":               "Jane Doe",
	"Max Mustermann":                       "Max Mustermann",
	"Prof.  Dr.  Sr.  Multi Title Example": "Multi Title Example",
	"Dr. Strange":                          "Strange",
	"Dr.    Stranger    ":                  "Stranger",
	// Spaces of the name itself are kept in tact as only the prefix should be changed
	"   Dr.   Sr.   John   Doe ":       "John   Doe",
	"  Prof.   Dr.   John    Smith   ": "John    Smith",
	// Ensure lower case prefixes are removed as well
	"prof. Dr. Max Mustermann": "Max Mustermann",
	"PROF. DR. Strange":        "Strange",
	"sr. Sr. John Doe":         "John Doe",
}

func TestRemoveTitlePrefixes(t *testing.T) {
	for actual, expected := range prefixesTests {
		assert.Equal(t, expected, RemoveTitlePrefixes(actual))
	}
}
