package utils

// Taken from kAdor here: https://stackoverflow.com/a/41604514
func StringFirstN(s string, n int) string {
	i := 0
	for j := range s {
		if i == n {
			return s[:j]
		}
		i++
	}
	return s
}

func RemoveDuplicates[T comparable](in []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}

	for _, item := range in {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}

	return list
}
