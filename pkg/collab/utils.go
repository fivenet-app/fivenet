package collab

import (
	"strconv"

	"github.com/cespare/xxhash"
)

// Max safe JS integer (2^53 - 1)
const maxJSInt = uint64(1<<53 - 1)

func MakeClientID(userID int32, sessionID string) uint64 {
	h := xxhash.New()
	h.Write([]byte(strconv.FormatInt(int64(userID), 10)))
	h.Write([]byte(sessionID))
	id := h.Sum64()

	// Clamp to max safe JS integer to be Yjs-safe
	return id & maxJSInt
}
