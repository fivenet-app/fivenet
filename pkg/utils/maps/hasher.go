package maps

import (
	"encoding/binary"
	"hash/maphash"
)

// For future use with https://github.com/puzpuzpuz/xsync package to replace the "copied" sync.Map code

func HashUint64(seed maphash.Seed, id uint64) uint64 {
	var h maphash.Hash
	h.SetSeed(seed)
	binary.Write(&h, binary.LittleEndian, id)
	hash := h.Sum64()
	h.Reset()
	return 31*hash + h.Sum64()
}

func HashString(seed maphash.Seed, in string) uint64 {
	var h maphash.Hash
	h.WriteString(in)
	h.SetSeed(seed)
	h.WriteString(in)
	hash := h.Sum64()
	h.Reset()
	return 31*hash + h.Sum64()
}
