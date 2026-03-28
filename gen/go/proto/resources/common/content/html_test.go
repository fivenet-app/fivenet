package content

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testsRichTextHtmlNodeFromHTMLNode = []struct {
	HTMLIn   string
	Expected string
}{
	{
		HTMLIn:   `<div id="test-id"><p class="text-orange-500 dark:text-orange-400">Hello World!</p></div>`,
		Expected: `<div id="test-id"><p class="text-orange-500 dark:text-orange-400">Hello World!</p></div>`,
	},
	{
		HTMLIn:   ``,
		Expected: ``,
	},
}

func TestRichTextHtmlNodeFromHTMLNode(t *testing.T) {
	t.Parallel()
	for _, v := range testsRichTextHtmlNodeFromHTMLNode {
		// Parse test input
		h, err := ParseHTML(v.HTMLIn)
		require.NoError(t, err)

		// Parsed HTML to RichTextHtmlNode
		n, err := FromHTMLNode(h)
		require.NoError(t, err)
		require.NotNil(t, n)

		// RichTextHtmlNode to pretty formatted HTML
		hout, err := n.ToHTML()
		require.NoError(t, err)
		assert.Equal(t, v.Expected, hout)
	}
}
