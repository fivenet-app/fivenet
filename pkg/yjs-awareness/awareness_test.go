package awareness

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeUpdateGolden(t *testing.T) {
	aw := NewAwareness()
	aw.SetState(42, map[string]any{
		"name": "Alice",
	})

	bytes := aw.EncodeUpdate([]uint64{42})

	want := []byte{
		0x01, // number of entries
		0x2A, // client ID 42
		0x01, // clock = 1
		0x76, // object marker
		0x01, // number of keys
		0x04, 'n', 'a', 'm', 'e',
		0x77, // string marker
		0x05, 'A', 'l', 'i', 'c', 'e',
	}

	assert.Equal(t, want, bytes)
}

func TestDecodeUpdateGolden(t *testing.T) {
	// Encoded state: client 1 with {"name": "Bob"}
	data := []byte{
		0x01, // one entry
		0x01, // client ID = 1
		0x01, // clock = 1
		0x76, // object
		0x01, // one key
		0x04, 'n', 'a', 'm', 'e',
		0x77, 0x03, 'B', 'o', 'b',
	}

	states, err := DecodeAwarenessUpdate(data)
	assert.NoError(t, err)
	assert.Len(t, states, 1)

	state, ok := states[1]
	assert.True(t, ok)
	assert.Equal(t, map[string]any{"name": "Bob"}, state)
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	aw := NewAwareness()
	aw.SetState(1, map[string]any{"user": "one", "color": "blue"})
	aw.SetState(2, map[string]any{"user": "two", "cursor": map[string]any{"x": 100, "y": 200}})

	data := aw.EncodeUpdate([]uint64{1, 2})

	states, err := DecodeAwarenessUpdate(data)
	assert.NoError(t, err)

	assert.Equal(t, map[string]any{"user": "one", "color": "blue"}, states[1])
	assert.Equal(t, map[string]any{
		"user":   "two",
		"cursor": map[string]any{"x": 100, "y": 200},
	}, states[2])
}
