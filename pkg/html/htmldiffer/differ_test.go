package htmldiffer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Diff(t *testing.T) {
	differ := New()

	for _, run := range []struct {
		old string
		new string

		result string
		msg    string
		error  bool
	}{
		{
			old:    "",
			new:    "",
			result: "",
			msg:    "Empty, stays empty",
		},
		{
			old:    "",
			new:    "<p>Hello <b>World</b>!</p>",
			result: `<p><span class="htmldiff bg-success-600">Hello </span><b><span class="htmldiff bg-success-600">World</span></b><span class="htmldiff bg-success-600">!</span></p>`,
			msg:    "Empty to basic html",
		},
		{
			old:    "<p>Hello <b>World</b>!</p>",
			new:    "<p>Goodbye <b>Tom</b>?</p>",
			result: `<p><span class="htmldiff bg-error-600">H</span><span class="htmldiff bg-success-600">Goodby</span>e<span class="htmldiff bg-error-600">llo</span> <b><span class="htmldiff bg-error-600">W</span><span class="htmldiff bg-success-600">T</span>o<span class="htmldiff bg-error-600">rld</span></b><span class="htmldiff bg-error-600">!</span><b><span class="htmldiff bg-success-600">m</span></b><span class="htmldiff bg-success-600">?</span></p>`,
			msg:    "Basic addition/update/removal of HTML content",
		},
		{
			old:    "<p>No <b>Changes</b></p>",
			new:    "<p>No <b>Changes</b></p>",
			result: "",
			msg:    "No content changes",
		},
		{
			old:    "<p>No Changes",
			new:    "",
			result: `<p><span class="htmldiff bg-error-600">No Changes</span></p>`,
			msg:    "\"Invalid\" HTML",
		},
	} {
		out, err := differ.Diff(run.old, run.new)
		if run.error {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, run.result, out, run.msg)
	}
}
