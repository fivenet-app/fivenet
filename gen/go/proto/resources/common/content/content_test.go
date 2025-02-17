package content

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testsJSONNodeFromHTMLNode = []struct {
	HTMLIn   string
	Expected string
}{
	{
		HTMLIn: `<div id="test-id"><p class="text-orange-500 dark:text-orange-400">Hello World!</p></div>`,
		Expected: `<div id="test-id">
  <p class="text-orange-500 dark:text-orange-400">
    Hello World!
  </p>
</div>`,
	},
	{
		HTMLIn:   ``,
		Expected: ``,
	},
}

func TestJSONNodeFromHTMLNode(t *testing.T) {
	for _, v := range testsJSONNodeFromHTMLNode {
		// Parse test input
		h, err := ParseHTML(v.HTMLIn)
		require.NoError(t, err)

		// Parsed HTML to JSONNode
		n, err := FromHTMLNode(h)
		require.NoError(t, err)
		require.NotNil(t, n)

		// JSONNode to pretty formatted HTML
		hout, err := n.ToHTMLP()
		require.NoError(t, err)
		assert.Equal(t, v.Expected, hout)
	}
}
