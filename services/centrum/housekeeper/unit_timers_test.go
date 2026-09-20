package housekeeper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseUnitPingKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		key string
		id  int64
		ok  bool
	}{
		{key: "ping.42", id: 42, ok: true},
		{key: "ping.0"},
		{key: "ping.-1"},
		{key: "ping.nope"},
		{key: "unit.42"},
		{key: "ping."},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			id, err := parseUnitPingKey(tt.key)
			if tt.ok {
				assert.NoError(t, err)
				assert.Equal(t, tt.id, id)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
