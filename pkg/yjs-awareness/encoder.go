package awareness

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
)

// EncodeAwarenessUpdate encodes the awareness states of the specified clients into a binary update packet.
func EncodeAwarenessUpdate(a *Awareness, clients []uint64) []byte {
	var out []byte
	// Write number of clients
	out = writeVarUint(out, uint64(len(clients)))
	for _, clientID := range clients {
		clock, ok := a.clocks[clientID]
		if !ok {
			// Skip unknown client IDs
			continue
		}
		state := a.States[clientID] // may be nil if offline
		out = writeVarUint(out, clientID)
		out = writeVarUint(out, clock)
		out = writeAny(out, state)
	}
	return out
}

// writeVarUint encodes a non-negative integer using LEB128 (variable-length) encoding.
func writeVarUint(buf []byte, value uint64) []byte {
	for {
		b := byte(value & 0x7F)
		value >>= 7
		if value == 0 {
			buf = append(buf, b)
			break
		} else {
			buf = append(buf, b|0x80)
		}
	}
	return buf
}

// writeVarInt encodes a 32-bit signed integer using zigzag encoding.
func writeVarInt(buf []byte, value int32) []byte {
	// Zigzag encode (negative values mapped to odd numbers)
	zigzag := uint32((value << 1) ^ (value >> 31))
	return writeVarUint(buf, uint64(zigzag))
}

// writeAny encodes an arbitrary JSON-like value (undefined, null, boolean, number, string, object, array, or []byte).
func writeAny(buf []byte, v interface{}) []byte {
	switch v := v.(type) {
	case UndefinedType:
		buf = append(buf, byte(markerUndefined))
		return buf
	case nil:
		buf = append(buf, byte(markerNull))
		return buf
	case bool:
		if v {
			buf = append(buf, byte(markerBoolTrue))
		} else {
			buf = append(buf, byte(markerBoolFalse))
		}
		return buf
	case int:
		vi := int64(v)
		if vi >= math.MinInt32 && vi <= math.MaxInt32 {
			buf = append(buf, byte(markerVarInt))
			buf = writeVarInt(buf, int32(vi))
		} else {
			buf = append(buf, byte(markerBigInt))
			var tmp [8]byte
			binary.LittleEndian.PutUint64(tmp[:], uint64(vi))
			buf = append(buf, tmp[:]...)
		}
		return buf
	case int8, int16, int32:
		// Convert to int32 and encode as varint
		vi := reflect.ValueOf(v).Int()
		buf = append(buf, byte(markerVarInt))
		buf = writeVarInt(buf, int32(vi))
		return buf
	case uint8, uint16, uint32:
		// These fit in positive int32 range
		vi := reflect.ValueOf(v).Uint()
		if vi <= math.MaxInt32 {
			buf = append(buf, byte(markerVarInt))
			buf = writeVarInt(buf, int32(vi))
		} else {
			buf = append(buf, byte(markerBigInt))
			var tmp [8]byte
			binary.LittleEndian.PutUint64(tmp[:], uint64(vi))
			buf = append(buf, tmp[:]...)
		}
		return buf
	case int64:
		vi := v
		if vi >= math.MinInt32 && vi <= math.MaxInt32 {
			buf = append(buf, byte(markerVarInt))
			buf = writeVarInt(buf, int32(vi))
		} else {
			buf = append(buf, byte(markerBigInt))
			var tmp [8]byte
			binary.LittleEndian.PutUint64(tmp[:], uint64(vi))
			buf = append(buf, tmp[:]...)
		}
		return buf
	case uint64:
		// Encode within 64-bit range
		if v <= math.MaxInt64 {
			buf = append(buf, byte(markerBigInt))
			var tmp [8]byte
			binary.LittleEndian.PutUint64(tmp[:], v)
			buf = append(buf, tmp[:]...)
		} else {
			// Values beyond 2^63-1 are truncated to 64 bits
			buf = append(buf, byte(markerBigInt))
			var tmp [8]byte
			binary.LittleEndian.PutUint64(tmp[:], v)
			buf = append(buf, tmp[:]...)
		}
		return buf
	case float32:
		buf = append(buf, byte(markerFloat32))
		var tmp [4]byte
		binary.LittleEndian.PutUint32(tmp[:], math.Float32bits(v))
		buf = append(buf, tmp[:]...)
		return buf
	case float64:
		// Use float32 if no precision loss, otherwise float64
		f32 := float32(v)
		if float64(f32) == v {
			buf = append(buf, byte(markerFloat32))
			var tmp [4]byte
			binary.LittleEndian.PutUint32(tmp[:], math.Float32bits(f32))
			buf = append(buf, tmp[:]...)
		} else {
			buf = append(buf, byte(markerFloat64))
			var tmp [8]byte
			binary.LittleEndian.PutUint64(tmp[:], math.Float64bits(v))
			buf = append(buf, tmp[:]...)
		}
		return buf
	case string:
		buf = append(buf, byte(markerString))
		strBytes := []byte(v)
		buf = writeVarUint(buf, uint64(len(strBytes)))
		buf = append(buf, strBytes...)
		return buf
	case []byte:
		buf = append(buf, byte(markerBinary))
		buf = writeVarUint(buf, uint64(len(v)))
		buf = append(buf, v...)
		return buf
	default:
		rv := reflect.ValueOf(v)
		if !rv.IsValid() {
			buf = append(buf, byte(markerUndefined))
			return buf
		}
		switch rv.Kind() {
		case reflect.Slice, reflect.Array:
			if rv.Kind() == reflect.Slice && rv.Type().Elem().Kind() == reflect.Uint8 {
				// Handle []uint8 (alias types) as binary
				data := v.([]byte)
				buf = append(buf, byte(markerBinary))
				buf = writeVarUint(buf, uint64(len(data)))
				buf = append(buf, data...)
				return buf
			}
			buf = append(buf, byte(markerArray))
			length := rv.Len()
			buf = writeVarUint(buf, uint64(length))
			for i := 0; i < length; i++ {
				elem := rv.Index(i).Interface()
				buf = writeAny(buf, elem)
			}
			return buf
		case reflect.Map:
			buf = append(buf, byte(markerObject))
			keys := rv.MapKeys()
			buf = writeVarUint(buf, uint64(len(keys)))
			for _, key := range keys {
				keyStr := fmt.Sprintf("%v", key.Interface())
				keyBytes := []byte(keyStr)
				buf = writeVarUint(buf, uint64(len(keyBytes)))
				buf = append(buf, keyBytes...)
				val := rv.MapIndex(key).Interface()
				buf = writeAny(buf, val)
			}
			return buf
		default:
			// Unsupported types -> undefined
			buf = append(buf, byte(markerUndefined))
			return buf
		}
	}
}
