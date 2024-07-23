package clvalue

import (
	"bytes"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
)

type Any []byte

func (v Any) Bytes() []byte {
	return v
}

func (v Any) String() string {
	return string(v)
}

func NewCLAny(data []byte) CLValue {
	var res CLValue
	res.Type = cltype.Any
	val := Any(data)
	res.Any = &val
	return res
}

func NewAnyFromBytes(data []byte) *CLValue {
	val := NewCLAny(data)
	return &val
}

func NewAnyFromBuffer(buf *bytes.Buffer) Any {
	return buf.Next(buf.Len())
}
