package awareness

// Type marker constants for lib0 encoding
const (
	markerUndefined = 127
	markerNull      = 126
	markerVarInt    = 125
	markerFloat32   = 124
	markerFloat64   = 123
	markerBigInt    = 122
	markerBoolFalse = 121
	markerBoolTrue  = 120
	markerString    = 119
	markerObject    = 118
	markerArray     = 117
	markerBinary    = 116
)

// UndefinedType is a special type to represent a JavaScript undefined value in awareness state.
type UndefinedType struct{}

// Undefined is a value that can be used in awareness state to represent an 'undefined' value.
var Undefined = UndefinedType{}

// Awareness holds awareness state for multiple clients, including logical clocks for each.
type Awareness struct {
	States map[uint64]interface{}
	clocks map[uint64]uint64
}

// NewAwareness creates a new Awareness instance.
func NewAwareness() *Awareness {
	return &Awareness{
		States: make(map[uint64]interface{}),
		clocks: make(map[uint64]uint64),
	}
}

// SetState sets or updates the awareness state for a given client ID.
// If state is nil, the client is marked as offline (the state is removed).
// Otherwise, the state must be a JSON-encodable value (e.g., map, slice, string, number, bool, nil).
// The internal clock for the client is incremented whenever the state is updated or removed.
func (a *Awareness) SetState(clientID uint64, state interface{}) {
	currClock, existed := a.clocks[clientID]
	if !existed {
		currClock = 0
	}
	// Increment clock on any state change
	currClock++
	a.clocks[clientID] = currClock
	if state == nil {
		// Remove state (mark offline)
		delete(a.States, clientID)
	} else {
		// Store new state
		a.States[clientID] = state
	}
}

// EncodeUpdate produces a binary awareness update for the specified clients.
// The update is encoded to include each client's latest awareness state (or null if the client is offline).
// Returns a byte slice containing the encoded update.
func (a *Awareness) EncodeUpdate(clients []uint64) []byte {
	return EncodeAwarenessUpdate(a, clients)
}
