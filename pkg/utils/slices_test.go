package utils

import (
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceDedup(t *testing.T) {
	t.Parallel()

	input := []string{"a", "b", "a", "c", "b"}
	expected := []string{"a", "b", "c"}
	result := SliceDedup(input)
	assert.ElementsMatch(t, expected, result)

	// Test with integers
	inputInt := []int{1, 2, 3, 2, 1, 4}
	expectedInt := []int{1, 2, 3, 4}
	resultInt := SliceDedup(inputInt)
	assert.ElementsMatch(t, expectedInt, resultInt)

	// Test with empty slice
	inputEmpty := []string{}
	expectedEmpty := []string{}
	resultEmpty := SliceDedup(inputEmpty)
	assert.ElementsMatch(t, expectedEmpty, resultEmpty)

	// Test with slice having no duplicates
	inputNoDuplicates := []string{"x", "y", "z"}
	expectedNoDuplicates := []string{"x", "y", "z"}
	resultNoDuplicates := SliceDedup(inputNoDuplicates)
	assert.ElementsMatch(t, expectedNoDuplicates, resultNoDuplicates)

	// Test with slice having all identical elements
	inputIdentical := []int{5, 5, 5, 5}
	expectedIdentical := []int{5}
	resultIdentical := SliceDedup(inputIdentical)
	assert.ElementsMatch(t, expectedIdentical, resultIdentical)
}

func BenchmarkSliceDedup(b *testing.B) {
	// Generate test data
	smallSlice := generateRandomSlice(100)   // 100 elements
	mediumSlice := generateRandomSlice(1000) // 1,000 elements
	largeSlice := generateRandomSlice(10000) // 10,000 elements

	b.Run("SmallSlice", func(b *testing.B) {
		for b.Loop() {
			SliceDedup(smallSlice)
		}
	})

	b.Run("MediumSlice", func(b *testing.B) {
		for b.Loop() {
			SliceDedup(mediumSlice)
		}
	})

	b.Run("LargeSlice", func(b *testing.B) {
		for b.Loop() {
			SliceDedup(largeSlice)
		}
	})
}

// Helper function to generate a random slice of integers.
func generateRandomSlice(size int) []int {
	slice := make([]int, size)
	for i := range slice {
		slice[i] = rand.IntN(size / 2) // Introduce duplicates
	}
	return slice
}

func TestSliceDiff(t *testing.T) {
	t.Parallel()

	a := []string{"hello", "example", "abc"}
	b := []string{"hello", "world", "test1", "abc"}

	added, removed := SliceDiff(a, b)
	assert.ElementsMatch(t, []string{"world", "test1"}, added)
	assert.ElementsMatch(t, []string{"example"}, removed)

	a = []string{"hello", "world", "abc"}
	b = []string{"hello", "world", "abc"}

	added, removed = SliceDiff(a, b)
	assert.Equal(t, []string{}, added)
	assert.Equal(t, []string{}, removed)

	a = []string{"hello", "world", "abc"}
	b = []string{"hello", "hello", "world", "abc"}

	added, removed = SliceDiff(a, b)
	assert.Equal(t, []string{}, added)
	assert.Equal(t, []string{}, removed)
}

func TestSliceDiffFunc(t *testing.T) {
	t.Parallel()

	a := []string{"hello", "example", "abc"}
	b := []string{"hello", "world", "test1", "abc"}

	keyFn := func(in string) string {
		return in // Use the string itself as the key
	}

	added, removed := SliceDiffFunc(a, b, keyFn)
	assert.ElementsMatch(t, []string{"world", "test1"}, added)
	assert.ElementsMatch(t, []string{"example"}, removed)

	a = []string{"hello", "world", "abc"}
	b = []string{"hello", "world", "abc"}

	added, removed = SliceDiffFunc(a, b, keyFn)
	assert.Empty(t, added)
	assert.Empty(t, removed)

	// Test with only duplicates added
	a = []string{"hello", "world", "abc"}
	b = []string{"hello", "hello", "world", "abc"}

	added, removed = SliceDiffFunc(a, b, keyFn)
	assert.Empty(t, added)
	assert.Empty(t, removed)

	// Test with slices having duplicates and new values
	a = []string{"hello", "example", "example", "abc"}
	b = []string{"hello", "world", "test1", "abc", "abc"}

	added, removed = SliceDiffFunc(a, b, keyFn)
	assert.ElementsMatch(t, []string{"world", "test1"}, added)
	assert.ElementsMatch(t, []string{"example"}, removed)
}

func TestMergeUniqueStrings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		lists    [][]string
		expected []string
	}{
		{
			name: "merges unique strings preserving first seen order",
			lists: [][]string{
				{"job-admin", "config-admin", "owner"},
				{"owner", "support", "job-admin"},
				{"audit", "support"},
			},
			expected: []string{"job-admin", "config-admin", "owner", "support", "audit"},
		},
		{
			name:     "no lists",
			lists:    nil,
			expected: []string{},
		},
		{
			name: "empty and nil lists",
			lists: [][]string{
				{},
				nil,
				{"user"},
			},
			expected: []string{"user"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expected, MergeUniqueStrings(tt.lists...))
		})
	}
}
