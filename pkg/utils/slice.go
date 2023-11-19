package utils

func InSlice[T comparable](in []T, search T) bool {
	for i := 0; i < len(in); i++ {
		if in[i] == search {
			return true
		}
	}

	return false
}

func InSliceIndex[T comparable](in []T, search T) int {
	for i := 0; i < len(in); i++ {
		if in[i] == search {
			return i
		}
	}

	return -1
}

func InSliceFunc[T any](in []T, searchFunc func(in T) bool) bool {
	for i := 0; i < len(in); i++ {
		if searchFunc(in[i]) {
			return true
		}
	}

	return false
}

func InSliceIndexFunc[T any](in []T, searchFunc func(in T) bool) int {
	for i := 0; i < len(in); i++ {
		if searchFunc(in[i]) {
			return i
		}
	}

	return -1
}

func RemoveFromSlice[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// Based on https://stackoverflow.com/a/66751055/2172930
func RemoveDuplicatesFromSlice[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
