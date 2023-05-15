package utils

func InStringSlice(in []string, search string) bool {
	for i := 0; i < len(in); i++ {
		if in[i] == search {
			return true
		}
	}

	return false
}

func RemoveFromStringSlice(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
