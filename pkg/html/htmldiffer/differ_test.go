package htmldiffer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFancyDiff(t *testing.T) {
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
		out, err := differ.FancyDiff(run.old, run.new)
		if run.error {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
		assert.Equal(t, run.result, out, run.msg)
	}
}

func TestPatchDiff(t *testing.T) {
	differ := New()

	for _, run := range []struct {
		old    string
		new    string
		result string
		msg    string
	}{
		{
			old:    "",
			new:    "",
			result: "",
			msg:    "Empty strings should return an empty diff",
		},
		{
			old: "<p>Hello World</p>",
			new: "<p>Hello Universe</p>",
			result: `--- a
+++ b
@@ -1 +1 @@
-<p>Hello World</p>
\ No newline at end of file
+<p>Hello Universe</p>
\ No newline at end of file
`,
			msg: "Basic content change",
		},
		{
			old:    `<img src="data:image/png;base64,ABC123"/>`,
			new:    `<img src="data:image/png;base64,XYZ456"/>`,
			result: "",
			msg:    "Image data should be replaced with placeholder (no diff)",
		},
		{
			old: `<p><img src="data:image/png;base64,ABC123"/> Hello World!</p>`,
			new: `<p><img src="data:image/png;base64,XYZ456"/> Hello Galaxy!</p>`,
			result: `--- a
+++ b
@@ -1 +1 @@
-<p><img src="IMAGE_DATA_OMITTED"/> Hello World!</p>
\ No newline at end of file
+<p><img src="IMAGE_DATA_OMITTED"/> Hello Galaxy!</p>
\ No newline at end of file
`,
			msg: "Image data should be replaced with placeholder (no diff)",
		},
		{
			old:    "<br/>",
			new:    "<br>",
			result: "",
			msg:    "Self-closing <br/> should be normalized to <br> and considered equal",
		},
		{
			old: "<p>Line 1<br/>Line 2</p>",
			new: "<p>Line 1<br>Line 3</p>",
			result: `--- a
+++ b
@@ -1 +1 @@
-<p>Line 1<br>Line 2</p>
\ No newline at end of file
+<p>Line 1<br>Line 3</p>
\ No newline at end of file
`,
			msg: "Line changes with normalized <br> tags",
		},
		{
			old:    "<p>Unchanged</p>",
			new:    "<p>Unchanged</p>",
			result: "",
			msg:    "Identical content should return an empty diff",
		},
	} {
		out := differ.PatchDiff(run.old, run.new)
		assert.Equal(t, run.result, out, run.msg)
	}
}
