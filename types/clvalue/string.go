package clvalue

import (
	"bytes"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

type String string

func (s *String) Bytes() []byte {
	sizeByte := SizeToBytes(len(*s))
	return append(sizeByte, []byte(*s)...)
}

func (s *String) String() string {
	return string(*s)
}

func NewCLString(val string) *CLValue {
	res := CLValue{}
	res.Type = cltype.String
	v := String(val)
	res.StringVal = &v
	return &res
}

func NewStringFromBytes(src []byte) *String {
	buf := bytes.NewBuffer(src)
	return NewStringFromBuffer(buf)
}

func NewStringFromBuffer(buffer *bytes.Buffer) *String {
	size := TrimByteSize(buffer)
	v := String(buffer.Next(int(size)))
	return &v
}
