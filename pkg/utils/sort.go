package utils

import (
	"sort"

	"golang.org/x/exp/constraints"
)

func SortInt32Slice[T constraints.Ordered](x []T) {
	sort.Sort(Comparable[T](x))
}

// Comparable generic for sorting in increasing order.
type Comparable[T constraints.Ordered] []T

func (x Comparable[T]) Len() int           { return len(x) }
func (x Comparable[T]) Less(i, j int) bool { return x[i] < x[j] }
func (x Comparable[T]) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
