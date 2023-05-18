package clvalue

import (
	"bytes"
	"encoding/hex"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

type ByteArray []byte

func (b ByteArray) Bytes() []byte {
	return b
}

func (b ByteArray) String() string {
	return hex.EncodeToString(b)
}

func NewCLByteArray(val []byte) CLValue {
	res := CLValue{}
	res.Type = &cltype.ByteArray{Size: uint32(len(val))}
	v := ByteArray(val)
	res.ByteArray = &v
	return res
}

func NewByteArrayFromBytes(data []byte, clType *cltype.ByteArray) ByteArray {
	return NewByteArrayFromBuffer(bytes.NewBuffer(data), clType)
}

func NewByteArrayFromBuffer(buf *bytes.Buffer, clType *cltype.ByteArray) ByteArray {
	return buf.Next(clType.Len())
}
