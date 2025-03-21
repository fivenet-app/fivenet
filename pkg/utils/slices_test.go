package utils

import (
	"math/rand"
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

func TestRemoveSliceDuplicates(t *testing.T) {
	// Test with strings
	input := []string{"a", "b", "a", "c", "b"}
	expected := []string{"a", "b", "c"}
	result := RemoveSliceDuplicates(input)
	assert.ElementsMatch(t, expected, result)

	// Test with integers
	inputInt := []int{1, 2, 3, 2, 1, 4}
	expectedInt := []int{1, 2, 3, 4}
	resultInt := RemoveSliceDuplicates(inputInt)
	assert.ElementsMatch(t, expectedInt, resultInt)

	// Test with empty slice
	inputEmpty := []string{}
	expectedEmpty := []string{}
	resultEmpty := RemoveSliceDuplicates(inputEmpty)
	assert.ElementsMatch(t, expectedEmpty, resultEmpty)

	// Test with slice having no duplicates
	inputNoDuplicates := []string{"x", "y", "z"}
	expectedNoDuplicates := []string{"x", "y", "z"}
	resultNoDuplicates := RemoveSliceDuplicates(inputNoDuplicates)
	assert.ElementsMatch(t, expectedNoDuplicates, resultNoDuplicates)

	// Test with slice having all identical elements
	inputIdentical := []int{5, 5, 5, 5}
	expectedIdentical := []int{5}
	resultIdentical := RemoveSliceDuplicates(inputIdentical)
	assert.ElementsMatch(t, expectedIdentical, resultIdentical)
}

func BenchmarkRemoveSliceDuplicates(b *testing.B) {
	// Generate test data
	smallSlice := generateRandomSlice(100)   // 100 elements
	mediumSlice := generateRandomSlice(1000) // 1,000 elements
	largeSlice := generateRandomSlice(10000) // 10,000 elements

	b.Run("SmallSlice", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			RemoveSliceDuplicates(smallSlice)
		}
	})

	b.Run("MediumSlice", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			RemoveSliceDuplicates(mediumSlice)
		}
	})

	b.Run("LargeSlice", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			RemoveSliceDuplicates(largeSlice)
		}
	})
}

// Helper function to generate a random slice of integers
func generateRandomSlice(size int) []int {
	slice := make([]int, size)
	for i := range slice {
		slice[i] = rand.Intn(size / 2) // Introduce duplicates
	}
	return slice
}
