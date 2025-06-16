// Package collab provides utilities for collaborative features.
package collab

import (
	"strconv"

	"github.com/cespare/xxhash/v2"
)

// maxJSInt is the maximum safe integer in JavaScript (2^53 - 1).
const maxJSInt = uint64(1<<53 - 1)

// MakeClientID generates a Yjs-safe client ID for collaborative sessions.
// The ID is derived from the userID and sessionID, hashed using xxhash,
// and clamped to the maximum safe JavaScript integer to ensure compatibility with Yjs.
func MakeClientID(userID int32, sessionID string) uint64 {
	h := xxhash.New()
	// Hash the userID as a string
	h.Write([]byte(strconv.FormatInt(int64(userID), 10)))
	// Hash the sessionID
	h.Write([]byte(sessionID))
	id := h.Sum64()

	// Clamp to max safe JS integer to be Yjs-safe
	return id & maxJSInt
}
