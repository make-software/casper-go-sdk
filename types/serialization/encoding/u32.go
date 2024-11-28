package encoding

import (
	"bytes"
	"encoding/binary"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
)

const U32SerializedLength = 4

type U32FromBytesDecoder struct{}

func NewU32FromBytesDecoder() *U32FromBytesDecoder {
	return &U32FromBytesDecoder{}
}

func (addr *U32FromBytesDecoder) FromBytes(inputBytes []byte) (uint32, []byte, error) {
	if len(inputBytes) < U32SerializedLength {
		return 0, nil, ErrInvalidBytesStructure
	}

	sizeBytes := inputBytes[:U32SerializedLength]
	remainder := inputBytes[U32SerializedLength:]

	var result uint32
	buf := bytes.NewReader(sizeBytes)
	err := binary.Read(buf, binary.LittleEndian, &result)
	if err != nil {
		return 0, nil, ErrInvalidBytesStructure
	}

	return result, remainder, nil
}

type U32ToBytesEncoder struct {
	val uint32
}

func NewU32ToBytesEncoder(val uint32) U32ToBytesEncoder {
	return U32ToBytesEncoder{
		val,
	}
}

func (enc U32ToBytesEncoder) Bytes() ([]byte, error) {
	return clvalue.NewCLUInt32(enc.val).Bytes(), nil
}
