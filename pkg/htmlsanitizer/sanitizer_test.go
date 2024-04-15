package htmlsanitizer

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestSanitize(t *testing.T) {
	for _, run := range []struct {
		input  string
		result string
		msg    string
	}{
		{
			input:  "",
			result: "",
			msg:    "Empty, stays empty",
		},
		{
			input:  "<h2 class=\"test\">HELLO WORLD!</h2> This is a test.",
			result: "<h2>HELLO WORLD!</h2> This is a test.",
			msg:    "Simple tags with class attribute (unallowed class)",
		},
		{
			input:  "<h2 class=\"ql-link\">HELLO WORLD!</h2> This is a test.",
			result: "<h2>HELLO WORLD!</h2> This is a test.",
			msg:    "Simple tags with class attributes (allowed class)",
		},
		{
			input:  "<h2 style=\"color: #ffffff\" class=\"ql-link\">HELLO WORLD!</h2> <p><span>This is a test.</span></p>",
			result: "<h2 style=\"color: #ffffff\">HELLO WORLD!</h2> <p><span>This is a test.</span></p>",
			msg:    "Simple tags with style attribute",
		},
		{
			input:  "</h2>HELLO WORLD!</p> This<table is a> test</table.",
			result: "</h2>HELLO WORLD!</p> This<table> test",
			msg:    "Broken tags",
		},
		{
			input:  "<script src=\"example.com\"></script>HELLO WORLD!</p> This<table is a> test</table.",
			result: "HELLO WORLD!</p> This<table> test",
			msg:    "Make sure bad tag is removed, even with broken tags",
		},
	} {
		assert.Equal(t, Sanitize(run.input), run.result, run.msg)
	}
}

func TestStripTags(t *testing.T) {
	for _, run := range []struct {
		input  string
		result string
		msg    string
	}{
		{
			input:  "",
			result: "",
			msg:    "Empty, stays empty",
		},
		{
			input:  "<h2 class=\"test\">HELLO WORLD!</h2> This is a test.",
			result: "HELLO WORLD! This is a test.",
			msg:    "Simple tags with attributes",
		},
		{
			input:  "</h2>HELLO WORLD!</p> This<table is a> test</table.",
			result: "HELLO WORLD! This test",
			msg:    "Broken tags",
		},
		{
			input:  "<script src=\"example.com\"></script>HELLO WORLD!</p> This<table is a> test</table.",
			result: "HELLO WORLD! This test",
			msg:    "Make sure bad tag is removed, even with broken tags",
		},
	} {
		assert.Equal(t, StripTags(run.input), run.result, run.msg)
	}
}
