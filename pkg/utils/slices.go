package utils

// RemoveSliceDuplicates returns a new slice with duplicate values removed, preserving order.
func RemoveSliceDuplicates[T comparable](in []T) []T {
	allKeys := make(map[T]struct{}, len(in))
	list := make([]T, 0, len(in))

	for _, item := range in {
		if _, value := allKeys[item]; !value {
			allKeys[item] = struct{}{}
			list = append(list, item)
		}
	}

	return list
}

// SlicesDifference returns the values added and removed between two slices, ignoring duplicates.
// Duplicates of values are ignored.
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

// SlicesDifferenceFunc returns the values added and removed between two slices, using a key function for comparison.
// Does not handle multiple additions of the same value as values are de-duplicated.
func SlicesDifferenceFunc[T comparable, S comparable](
	a, b []T,
	keyFn func(in T) S,
) ([]T, []T) {
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

	added := []T{}
	removed := []T{}
	for s, v := range temp {
		if v == 0 {
			removed = append(removed, vals[s])
		} else if v < 0 {
			added = append(added, vals[s])
		}
	}

	return added, removed
}
