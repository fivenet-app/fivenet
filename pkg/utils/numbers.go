package utils

import "math"

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
func ToUint32Checked(v int) (uint32, bool) {
	if v < 0 {
		return 0, false
	}
	if uint64(v) > math.MaxUint32 {
		return 0, false
	}
	return uint32(v), true
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
