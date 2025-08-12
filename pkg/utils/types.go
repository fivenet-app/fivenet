package utils

import "fmt"

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
