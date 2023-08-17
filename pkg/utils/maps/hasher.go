package maps

import (
	"hash/maphash"
)

// For future use with https://github.com/puzpuzpuz/xsync package to replace the "copied" sync.Map code

func HashString(seed maphash.Seed, in string) uint64 {
	var h maphash.Hash
	h.WriteString(in)
	h.SetSeed(seed)
	h.WriteString(in)
	hash := h.Sum64()
	return 31*hash + h.Sum64()
}
