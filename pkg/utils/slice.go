package utils

func InSlice[T comparable](in []T, search T) bool {
	for i := 0; i < len(in); i++ {
		if in[i] == search {
			return true
		}
	}

	return false
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
