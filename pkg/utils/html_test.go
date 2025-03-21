package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsHeaderTag(t *testing.T) {
	for actual, expected := range map[string]bool{
		"div":   false,
		"s pan": false,
		"h3":    true,
		"h7":    false,
		"h6":    true,
		"h1":    true,
		"h0":    false,
	} {
		assert.Equal(t, expected, IsHeaderTag(actual), actual)
	}
}
