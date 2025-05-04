package modules

import "testing"

func TestConstructUserNickname(t *testing.T) {
	tests := []struct {
		name      string
		firstname string
		lastname  string
		prefix    string
		suffix    string
		expected  string
	}{
		{
			name:      "Basic name without prefix or suffix",
			firstname: "John",
			lastname:  "Doe",
			prefix:    "",
			suffix:    "",
			expected:  "John Doe",
		},
		{
			name:      "Name with prefix",
			firstname: "John",
			lastname:  "Doe",
			prefix:    "Dr.",
			suffix:    "",
			expected:  "Dr. John Doe",
		},
		{
			name:      "Name with suffix",
			firstname: "John",
			lastname:  "Doe",
			prefix:    "",
			suffix:    "PhD",
			expected:  "John Doe PhD",
		},
		{
			name:      "Name with prefix and suffix",
			firstname: "John",
			lastname:  "Doe",
			prefix:    "Dr.",
			suffix:    "PhD",
			expected:  "Dr. John Doe PhD",
		},
		{
			name:      "Name exceeding max length",
			firstname: "Jonathan",
			lastname:  "Doe-Smith-Jackson",
			prefix:    "Dr.",
			suffix:    "PhD",
			expected:  "Dr. J. Doe-Smith-Jackson PhD",
		},
		{
			name:      "Name with only prefix exceeding max length",
			firstname: "John",
			lastname:  "Doe",
			prefix:    "A very longp",
			suffix:    "",
			expected:  "A very longp John Doe",
		},
		{
			name:      "Name with only suffix exceeding max length",
			firstname: "John",
			lastname:  "Doe",
			prefix:    "",
			suffix:    "A very longs",
			expected:  "John Doe A very longs",
		},
		{
			name:      "Empty lastname",
			firstname: "Max",
			lastname:  "",
			prefix:    "Dr.",
			suffix:    "PhD",
			expected:  "Dr. Max PhD",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := &UserInfo{}
			result := g.constructUserNickname(test.firstname, test.lastname, test.prefix, test.suffix)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}
