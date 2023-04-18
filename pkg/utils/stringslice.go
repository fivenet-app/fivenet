package utils

func InStringSlice(in []string, search string) bool {
	for i := 0; i < len(in); i++ {
		if in[i] == search {
			return true
		}
	}

	return false
}
