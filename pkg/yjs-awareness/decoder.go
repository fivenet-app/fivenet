package awareness

import (
	"encoding/binary"
	"errors"
	"math"
)

// DecodeAwarenessUpdate parses a binary awareness update and returns a map of client IDs to their states.
// If a state is null (client went offline), it will be represented by a nil value in the map.
// If a value is the special Undefined sentinel, it will be represented by awareness.Undefined.
func DecodeAwarenessUpdate(data []byte) (map[uint64]interface{}, error) {
	dec := decoder{data: data}
	count, err := dec.ReadVarUint()
	if err != nil {
		return nil, err
	}
	states := make(map[uint64]interface{}, count)
	for i := uint64(0); i < count; i++ {
		clientID, err := dec.ReadVarUint()
		if err != nil {
			return nil, err
		}
		// Read and ignore clock value
		_, err = dec.ReadVarUint()
		if err != nil {
			return nil, err
		}
		stateVal, err := dec.ReadAny()
		if err != nil {
			return nil, err
		}
		states[clientID] = stateVal
	}
	return states, nil
}

// ModifyAwarenessUpdate applies a transformation function to all state entries in an awareness update.
// The modifyFn is called for each state (which can be any JSON-like value or nil). It should return a new state value.
// The returned update will contain the same client IDs and clocks, with each state replaced by the value returned by modifyFn.
// If modifyFn returns nil, the state will be encoded as null (client offline).
// If modifyFn returns the Undefined sentinel, the state will be encoded as undefined.
func ModifyAwarenessUpdate(update []byte, modifyFn func(interface{}) interface{}) ([]byte, error) {
	dec := decoder{data: update}
	count, err := dec.ReadVarUint()
	if err != nil {
		return nil, err
	}
	var out []byte
	out = writeVarUint(out, count)
	for i := uint64(0); i < count; i++ {
		clientID, err := dec.ReadVarUint()
		if err != nil {
			return nil, err
		}
		clock, err := dec.ReadVarUint()
		if err != nil {
			return nil, err
		}
		stateVal, err := dec.ReadAny()
		if err != nil {
			return nil, err
		}
		// Apply the modification function to the state
		newState := modifyFn(stateVal)
		// Write out the entry with the same client and clock, but modified state
		out = writeVarUint(out, clientID)
		out = writeVarUint(out, clock)
		out = writeAny(out, newState)
	}
	return out, nil
}

// decoder is a helper for reading values from a byte slice.
type decoder struct {
	data []byte
	pos  int
}

// ReadVarUint decodes a LEB128-encoded unsigned integer from the buffer.
func (d *decoder) ReadVarUint() (uint64, error) {
	var result uint64
	var shift uint
	for {
		if d.pos >= len(d.data) {
			return 0, errors.New("unexpected end of data in ReadVarUint")
		}
		b := d.data[d.pos]
		d.pos++
		result |= uint64(b&0x7F) << shift
		if (b & 0x80) == 0 {
			break
		}
		shift += 7
		if shift >= 64 {
			return 0, errors.New("varUint value is too large")
		}
	}
	return result, nil
}

// ReadVarInt decodes a zigzag-encoded 32-bit signed integer from the buffer.
func (d *decoder) ReadVarInt() (int32, error) {
	ui, err := d.ReadVarUint()
	if err != nil {
		return 0, err
	}
	// Zigzag decode: if the least significant bit is 1, the original was negative.
	var n uint64 = ui >> 1
	if ui&1 != 0 {
		n = ^n
	}
	iv := int64(n)
	// Cast down to int32 (value is within 32-bit range)
	return int32(iv), nil
}

// ReadAny decodes the next value from the buffer according to the awareness update encoding.
func (d *decoder) ReadAny() (interface{}, error) {
	if d.pos >= len(d.data) {
		return nil, errors.New("unexpected end of data in ReadAny")
	}
	prefix := d.data[d.pos]
	d.pos++
	switch prefix {
	case markerUndefined:
		return Undefined, nil
	case markerNull:
		return nil, nil
	case markerVarInt:
		iv, err := d.ReadVarInt()
		if err != nil {
			return nil, err
		}
		return int(iv), nil
	case markerFloat32:
		if d.pos+4 > len(d.data) {
			return nil, errors.New("unexpected end of data in float32")
		}
		bits := binary.LittleEndian.Uint32(d.data[d.pos : d.pos+4])
		d.pos += 4
		return math.Float32frombits(bits), nil
	case markerFloat64:
		if d.pos+8 > len(d.data) {
			return nil, errors.New("unexpected end of data in float64")
		}
		bits := binary.LittleEndian.Uint64(d.data[d.pos : d.pos+8])
		d.pos += 8
		return math.Float64frombits(bits), nil
	case markerBigInt:
		if d.pos+8 > len(d.data) {
			return nil, errors.New("unexpected end of data in bigInt")
		}
		bits := binary.LittleEndian.Uint64(d.data[d.pos : d.pos+8])
		d.pos += 8
		i := int64(bits)
		if i <= math.MaxInt32 && i >= math.MinInt32 {
			return int(i), nil
		}
		return i, nil
	case markerBoolFalse:
		return false, nil
	case markerBoolTrue:
		return true, nil
	case markerString:
		length, err := d.ReadVarUint()
		if err != nil {
			return nil, err
		}
		if d.pos+int(length) > len(d.data) {
			return nil, errors.New("unexpected end of data in string")
		}
		strBytes := d.data[d.pos : d.pos+int(length)]
		d.pos += int(length)
		return string(strBytes), nil
	case markerObject:
		length, err := d.ReadVarUint()
		if err != nil {
			return nil, err
		}
		obj := make(map[string]interface{}, length)
		for i := uint64(0); i < length; i++ {
			keyLen, err := d.ReadVarUint()
			if err != nil {
				return nil, err
			}
			if d.pos+int(keyLen) > len(d.data) {
				return nil, errors.New("unexpected end of data in object key")
			}
			key := string(d.data[d.pos : d.pos+int(keyLen)])
			d.pos += int(keyLen)
			val, err := d.ReadAny()
			if err != nil {
				return nil, err
			}
			obj[key] = val
		}
		return obj, nil
	case markerArray:
		length, err := d.ReadVarUint()
		if err != nil {
			return nil, err
		}
		arr := make([]interface{}, 0, length)
		for i := uint64(0); i < length; i++ {
			val, err := d.ReadAny()
			if err != nil {
				return nil, err
			}
			arr = append(arr, val)
		}
		return arr, nil
	case markerBinary:
		length, err := d.ReadVarUint()
		if err != nil {
			return nil, err
		}
		if d.pos+int(length) > len(d.data) {
			return nil, errors.New("unexpected end of data in binary")
		}
		b := make([]byte, length)
		copy(b, d.data[d.pos:d.pos+int(length)])
		d.pos += int(length)
		return b, nil
	default:
		// Unknown prefix type; treat as undefined
		return Undefined, nil
	}
}
