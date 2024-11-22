package encoding

import (
	"encoding/binary"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
)

type StringFromBytesDecoder struct{}

func (dec *StringFromBytesDecoder) FromBytes(bytes []byte) (string, []byte, error) {
	if len(bytes) < U32SerializedLength {
		return "", nil, ErrInvalidBytesStructure
	}

	length := binary.LittleEndian.Uint32(bytes[:4])

	if len(bytes[4:]) < int(length) {
		return "", nil, ErrInvalidBytesStructure
	}

	strBytes := bytes[U32SerializedLength : U32SerializedLength+length]

	result := string(strBytes)
	return result, bytes[4+length:], nil
}

func StringSerializedLength(val string) int {
	return U32SerializedLength + len(val)
}

func BytesSerializedLength(val []byte) int {
	return U32SerializedLength + len(val)
}

type StringToBytesEncoder struct {
	val string
}

func NewStringToBytesEncoder(val string) StringToBytesEncoder {
	return StringToBytesEncoder{
		val,
	}
}

func (enc StringToBytesEncoder) Bytes() ([]byte, error) {
	return clvalue.NewCLString(enc.val).Bytes(), nil
}
