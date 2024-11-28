package encoding

import (
	"encoding/binary"
)

const U16SerializedLength = 2

type U16FromBytesDecoder struct{}

// FromBytes function to deserialize a u16 from bytes
func (addr *U16FromBytesDecoder) FromBytes(inputBytes []byte) (uint16, []byte, error) {
	if len(inputBytes) < 2 {
		return 0, nil, ErrInvalidBytesStructure
	}

	result := binary.LittleEndian.Uint16(inputBytes[:2])
	return result, inputBytes[2:], nil
}
