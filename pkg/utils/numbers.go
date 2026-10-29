package utils

import "math"

type signedInteger interface {
	~int | ~int64
}

// ToUint32Saturated converts int to uint32 and clamps out-of-range values.
func ToUint32Saturated(v int) uint32 {
	if v <= 0 {
		return 0
	}
	if uint64(v) > math.MaxUint32 {
		return math.MaxUint32
	}
	return uint32(v)
}

// ToUint32Checked converts int to uint32 and reports whether the value was in range.
func ToUint32Checked[T signedInteger](v T) (uint32, bool) {
	if v < 0 {
		return 0, false
	}
	if uint64(v) > math.MaxUint32 {
		return 0, false
	}
	return uint32(v), true
}

// ToInt32Checked converts an integer to int32 and reports whether the value was in range.
func ToInt32Checked[T signedInteger](v T) (int32, bool) {
	if v < math.MinInt32 || v > math.MaxInt32 {
		return 0, false
	}
	return int32(v), true
}

func SaturatingAddUint32(base uint32, add int) uint32 {
	if add <= 0 {
		return base
	}
	if add > int(math.MaxUint32-base) {
		return math.MaxUint32
	}
	//nolint:gosec // G115: The overflow is checked above, so this is safe.
	return base + uint32(add)
}
