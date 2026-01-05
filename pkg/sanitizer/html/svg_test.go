package htmlsanitizer

import (
	"testing"
)

func TestSVGSanitize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid SVG with allowed attributes",
			input:    `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><path d="M10 10 H 90 V 90 H 10 Z" stroke-width="2" stroke="black" fill="none"/></svg>`,
			expected: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><path d="M10 10 H 90 V 90 H 10 Z" stroke-width="2" stroke="black" fill="none"/></svg>`,
		},
		{
			name:     "SVG with disallowed attributes",
			input:    `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" onclick="alert('test')"><path d="M10 10 H 90 V 90 H 10 Z" stroke-width="2" stroke="black" fill="none"/></svg>`,
			expected: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><path d="M10 10 H 90 V 90 H 10 Z" stroke-width="2" stroke="black" fill="none"/></svg>`,
		},
		{
			name:     "SVG with unsupported elements",
			input:    `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><script>alert('test')</script></svg>`,
			expected: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"></svg>`,
		},
		{
			name:     "Empty SVG",
			input:    `<svg xmlns="http://www.w3.org/2000/svg"></svg>`,
			expected: `<svg xmlns="http://www.w3.org/2000/svg"></svg>`,
		},
		{
			name:     "SVG with extra whitespace",
			input:    `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100">   <path d="M10 10 H 90 V 90 H 10 Z" stroke-width="2" stroke="black" fill="none"/>   </svg>`,
			expected: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100">   <path d="M10 10 H 90 V 90 H 10 Z" stroke-width="2" stroke="black" fill="none"/>   </svg>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeSVG(tt.input)
			if tt.expected != result {
				t.Errorf(
					"SanitizeSVG failed for test '%s':\nexpected:\n%s\ngot:\n%s",
					tt.name,
					tt.expected,
					result,
				)
			}
		})
	}
}
