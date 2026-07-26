package utils

// EqualInt32Ptr compares two *int32 pointers for equality.
func EqualInt32Ptr(a, b *int32) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
