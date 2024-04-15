package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlicesDifference(t *testing.T) {
	a := []string{"hello", "example", "abc"}
	b := []string{"hello", "world", "test1", "abc"}

	added, removed := SlicesDifference(a, b)
	assert.ElementsMatch(t, []string{"world", "test1"}, added)
	assert.ElementsMatch(t, []string{"example"}, removed)

	a = []string{"hello", "world", "abc"}
	b = []string{"hello", "world", "abc"}

	added, removed = SlicesDifference(a, b)
	assert.Equal(t, []string{}, added)
	assert.Equal(t, []string{}, removed)

	a = []string{"hello", "world", "abc"}
	b = []string{"hello", "hello", "world", "abc"}

	added, removed = SlicesDifference(a, b)
	assert.Equal(t, []string{}, added)
	assert.Equal(t, []string{}, removed)
}
