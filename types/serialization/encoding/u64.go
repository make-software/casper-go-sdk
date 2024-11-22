package encoding

import (
	"encoding/binary"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
)

const U64SerializedLength = 8

type U64FromBytesDecoder struct{}

// FromBytes deserializes a uint64 from a byte slice.
func (addr *U64FromBytesDecoder) FromBytes(bytes []byte) (uint64, []byte, error) {
	if len(bytes) < U64SerializedLength {
		return 0, nil, ErrInvalidBytesStructure
	}

	result := binary.LittleEndian.Uint64(bytes[:U64SerializedLength])
	remainder := bytes[U64SerializedLength:]

	return result, remainder, nil
}

type U64ToBytesEncoder struct {
	val uint64
}

func NewU64ToBytesEncoder(val uint64) U64ToBytesEncoder {
	return U64ToBytesEncoder{
		val,
	}
}

func (enc U64ToBytesEncoder) Bytes() ([]byte, error) {
	return clvalue.NewCLUInt64(enc.val).Bytes(), nil
}
