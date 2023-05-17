package utils

import "sort"

func SortInt32Slice(x []int32) {
	sort.Sort(Int32Slice(x))
}

// Int32Slice attaches the methods of Interface to []int32, sorting in increasing order.
type Int32Slice []int32

func (x Int32Slice) Len() int           { return len(x) }
func (x Int32Slice) Less(i, j int) bool { return x[i] < x[j] }
func (x Int32Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
