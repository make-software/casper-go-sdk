package encoding

import (
	"errors"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
)

var (
	ErrEmptyBytesSource      = errors.New("empty bytes source")
	ErrInvalidBytesStructure = errors.New("invalid bytes structure")
)

const BoolSerializedLength = 1

type BoolFromBytesDecoder struct{}

func NewBoolFromBytesDecoder() *BoolFromBytesDecoder {
	return &BoolFromBytesDecoder{}
}

func (addr *BoolFromBytesDecoder) FromBytes(bytes []byte) (bool, []byte, error) {
	if len(bytes) == 0 {
		return false, nil, ErrEmptyBytesSource
	}

	switch bytes[0] {
	case 1:
		return true, bytes[1:], nil
	case 0:
		return false, bytes[1:], nil
	default:
		return false, nil, ErrInvalidBytesStructure
	}
}

type BoolToBytesEncoder struct {
	val bool
}

func NewBoolToBytesEncoder(val bool) BoolToBytesEncoder {
	return BoolToBytesEncoder{
		val,
	}
}

func (enc BoolToBytesEncoder) Bytes() ([]byte, error) {
	return clvalue.NewCLBool(enc.val).Bytes(), nil
}
