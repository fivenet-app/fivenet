package utils

// SliceDedup removes duplicate values while preserving order.
//
// The returned slice reuses the backing array of in and may overwrite its
// elements. Callers must not rely on in remaining unchanged.
func SliceDedup[T comparable](in []T) []T {
	seen := make(map[T]struct{}, len(in))
	list := make([]T, 0, len(in))

	for _, item := range in {
		if _, value := seen[item]; !value {
			seen[item] = struct{}{}
			list = append(list, item)
		}
	}

	return list
}

// SliceDedupFn returns a new slice with duplicate values removed determined by a function, preserving order.
// keyFn should return true for items that are considered duplicates.
func SliceDedupFn[T comparable, V comparable](in []T, keyFn func(T) V) []T {
	seen := make(map[V]struct{}, len(in))
	list := make([]T, 0, len(in))

	for _, item := range in {
		key := keyFn(item)
		if _, value := seen[key]; !value {
			seen[key] = struct{}{}
			list = append(list, item)
		}
	}

	return list
}

// SliceDiff returns values added to b and removed from a.
//
// Duplicate values are ignored. Result ordering follows the order in which
// values appear in b for added values and a for removed values.
func SliceDiff[T comparable](a, b []T) ([]T, []T) {
	aSet := make(map[T]struct{}, len(a))
	bSet := make(map[T]struct{}, len(b))

	for _, value := range a {
		aSet[value] = struct{}{}
	}

	for _, value := range b {
		bSet[value] = struct{}{}
	}

	added := make([]T, 0)
	addedSeen := make(map[T]struct{})

	for _, value := range b {
		if _, exists := aSet[value]; exists {
			continue
		}
		if _, exists := addedSeen[value]; exists {
			continue
		}

		addedSeen[value] = struct{}{}
		added = append(added, value)
	}

	removed := make([]T, 0)
	removedSeen := make(map[T]struct{})

	for _, value := range a {
		if _, exists := bSet[value]; exists {
			continue
		}
		if _, exists := removedSeen[value]; exists {
			continue
		}

		removedSeen[value] = struct{}{}
		removed = append(removed, value)
	}

	return added, removed
}

// SliceDiffFunc returns values added to b and removed from a,
// comparing values using keyFn.
//
// The first value encountered for each key is retained. Duplicate keys are
// ignored. Result ordering follows b for added values and a for removed values.
func SliceDiffFunc[T comparable, S comparable](
	a, b []T,
	keyFn func(in T) S,
) ([]T, []T) {
	aSet := make(map[T]struct{}, len(a))
	bSet := make(map[T]struct{}, len(b))

	for _, value := range a {
		aSet[value] = struct{}{}
	}

	for _, value := range b {
		bSet[value] = struct{}{}
	}

	added := make([]T, 0)
	addedSeen := make(map[T]struct{})

	for _, value := range b {
		if _, exists := aSet[value]; exists {
			continue
		}
		if _, exists := addedSeen[value]; exists {
			continue
		}

		addedSeen[value] = struct{}{}
		added = append(added, value)
	}

	removed := make([]T, 0)
	removedSeen := make(map[T]struct{})

	for _, value := range a {
		if _, exists := bSet[value]; exists {
			continue
		}
		if _, exists := removedSeen[value]; exists {
			continue
		}

		removedSeen[value] = struct{}{}
		removed = append(removed, value)
	}

	return added, removed
}

// MergeUniqueStrings merges multiple slices of strings into a single slice,
// removing duplicates while preserving order.
func MergeUniqueStrings(lists ...[]string) []string {
	totalLen := 0
	for _, list := range lists {
		totalLen += len(list)
	}

	merged := make([]string, 0, totalLen)
	seen := make(map[string]struct{}, totalLen)

	for _, list := range lists {
		for _, value := range list {
			if _, exists := seen[value]; exists {
				continue
			}

			seen[value] = struct{}{}
			merged = append(merged, value)
		}
	}

	return merged
}
