package utils

import "fmt"

// SliceUint64ToInt64 converts a slice of uint64 values to a slice of int64 values.
func SliceUint64ToInt64(in []uint64) (out []int64) {
	for _, i := range in {
		out = append(out, int64(i))
	}
	return
}

// SliceInt32ToInt converts a slice of int32 values to a slice of int values.
func SliceInt32ToInt(in []int32) (out []int) {
	for _, i := range in {
		out = append(out, int(i))
	}
	return
}

// ToStringMap converts a map with values of any type to a map with string values.
func ToStringMap(params map[string]any) map[string]string {
	out := make(map[string]string, len(params))
	for k, v := range params {
		switch val := v.(type) {
		case string:
			out[k] = val
		case fmt.Stringer:
			out[k] = val.String()
		default:
			out[k] = fmt.Sprintf("%v", v)
		}
	}
	return out
}
