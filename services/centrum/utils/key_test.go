package centrumutils

import (
	"strconv"
	"testing"
)

func BenchmarkInt64ToString(b *testing.B) {
	b.Run("Itoa", func(b *testing.B) {
		for b.Loop() {
			_ = strconv.Itoa(123456789)
		}
	})

	b.Run("FormatInt", func(b *testing.B) {
		for b.Loop() {
			_ = strconv.FormatInt(123456789, 10)
		}
	})
}
