package centrumstate

import (
	"strconv"
	"testing"
)

func BenchmarkUint64ToString(b *testing.B) {
	b.Run("Itoa", func(b *testing.B) {
		for b.Loop() {
			_ = strconv.Itoa(123456789)
		}
	})

	b.Run("FormatUint", func(b *testing.B) {
		for b.Loop() {
			_ = strconv.FormatUint(123456789, 10)
		}
	})
}
