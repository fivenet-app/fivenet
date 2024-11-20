package utils

// Taken from "KAdot" here: https://stackoverflow.com/a/41604514
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
