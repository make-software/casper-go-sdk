package clvalue

import (
	"bytes"
	"fmt"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

type UInt8 uint8

func (v *UInt8) Bytes() []byte {
	return []byte{byte(*v)}
}

func (v *UInt8) String() string {
	return fmt.Sprintf("%d", *v)
}

func (v *UInt8) Value() uint8 {
	return uint8(*v)
}

func NewCLUint8(val uint8) *CLValue {
	res := CLValue{}
	res.Type = cltype.UInt8
	v := UInt8(val)
	res.UI8 = &v
	return &res
}

func NewUInt8FromBytes(source []byte) (*UInt8, error) {
	buf := bytes.NewBuffer(source)
	return NewUInt8FromBuffer(buf)
}

func NewUInt8FromBuffer(buffer *bytes.Buffer) (*UInt8, error) {
	buf, err := buffer.ReadByte()
	if err != nil {
		return nil, err
	}
	val := UInt8(buf)
	return &val, nil
}
