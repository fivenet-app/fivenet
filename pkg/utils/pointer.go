package utils

type PointerType[T any] interface {
	// This constrains the type to be a pointer of T
	*T
}

func Deref[T any, P PointerType[T]](in P) T {
	if in == nil {
		var zero T
		return zero
	}
	return *in
}
