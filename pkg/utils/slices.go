package utils

import (
	"cmp"
	"slices"
)

func RemoveSliceDuplicates[T ~[]E, E cmp.Ordered](in T) T {
	slices.Sort(in)
	return slices.Compact(in)
}

// SlicesDifference duplicates of values are ignored
func SlicesDifference[T comparable](a, b []T) ([]T, []T) {
	temp := map[T]int{}
	for _, s := range a {
		if _, ok := temp[s]; !ok {
			temp[s] = 0
		}
	}
	for _, s := range b {
		if _, ok := temp[s]; !ok {
			temp[s] = -1
		} else {
			temp[s] = 1
		}
	}

	added, removed := []T{}, []T{}
	for s, v := range temp {
		if v == 0 {
			removed = append(removed, s)
		} else if v < 0 {
			added = append(added, s)
		}
	}

	return added, removed
}

func SlicesDifferenceFunc[T comparable, S comparable](a, b []T, keyFn func(in T) S) (added []T, removed []T) {
	temp := map[S]int{}
	vals := map[S]T{}
	for _, i := range a {
		s := keyFn(i)
		if _, ok := temp[s]; !ok {
			temp[s] = 0
			vals[s] = i
		}
	}
	for _, i := range b {
		s := keyFn(i)
		if _, ok := temp[s]; !ok {
			temp[s] = -1
			vals[s] = i
		} else {
			temp[s] = 1
		}
	}

	for s, v := range temp {
		if v == 0 {
			removed = append(removed, vals[s])
		} else if v < 0 {
			added = append(added, vals[s])
		}
	}

	return
}
