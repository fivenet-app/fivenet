package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeKey(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"key#value", "key.value"},
		{"key:value", "key.value"},
		{"key/value", "key.value"},
		{"key=value", "key_value"},
		{"key#value:key/value=other", "key.value.value.value_other"},
		{"plainkey", "plainkey"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := SanitizeKey(tt.input)
			assert.Equal(t, tt.expected, result, "SanitizeKey(%q)", tt.input)
		})
	}
}
