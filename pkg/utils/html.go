package utils

// IsHeaderTag returns true if the given tag is a valid HTML header tag (h1-h6).
func IsHeaderTag(tag string) bool {
	return len(tag) == 2 && tag[0] == 'h' && tag[1] >= '1' && tag[1] <= '6'
}
