package math

import (
	"encoding/binary"
	"fmt"
	"math"
)

// BytesToFloat32 converts a 4 long slice of bytes into a float32. Uses little
// Endian format.
func BytesToFloat32(buf []byte) (float32, error) {
	if len(buf) != 4 {
		return 0, fmt.Errorf("cannot convert %v bytes to float32, 4 bytes are required", len(buf))
	}

	b := binary.LittleEndian.Uint32(buf)
	return math.Float32frombits(b), nil
}

// BytesToUInt32 converts a 4 long slice of bytes into a float16. Uses little
// Endian format.
func BytesToUInt16(buf []byte) (uint16, error) {
	if len(buf) != 2 {
		return 0, fmt.Errorf("cannot convert %v bytes to uint, 2 bytes are required", len(buf))
	}

	b := binary.LittleEndian.Uint16(buf)
	return b, nil
}

func BytesToFloat64(buf []byte) (float64, error) {
	if len(buf) != 8 {
		return 0, fmt.Errorf("cannot convert %v bytes to float64, 8 bytes are required", len(buf))
	}

	b := binary.LittleEndian.Uint64(buf)
	return math.Float64frombits(b), nil
}
