package clvalue

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
)

type UInt32 uint32

func (v *UInt32) Bytes() []byte {
	data := make([]byte, cltype.Int32ByteSize)
	binary.LittleEndian.PutUint32(data, uint32(*v))
	return data
}

func (v *UInt32) String() string {
	return fmt.Sprintf("%d", *v)
}

func (v *UInt32) Value() uint32 {
	return uint32(*v)
}

func NewUint32FromBytes(source []byte) (*UInt32, error) {
	buf := bytes.NewBuffer(source)
	return NewUint32FromBuffer(buf)
}

func NewCLUInt32(val uint32) *CLValue {
	res := CLValue{}
	res.Type = cltype.UInt32
	v := UInt32(val)
	res.UI32 = &v
	return &res
}

func NewUint32FromBuffer(buffer *bytes.Buffer) (*UInt32, error) {
	if buffer.Len() < cltype.Int32ByteSize {
		return nil, errors.New("buffer size is too small")
	}
	buf := buffer.Next(cltype.Int32ByteSize)
	val := UInt32(binary.LittleEndian.Uint32(buf))
	return &val, nil
}

func TrimByteSize(buf *bytes.Buffer) (size uint32, err error) {
	buffer, err := NewUint32FromBuffer(buf)
	if err != nil {
		return 0, err
	}
	return buffer.Value(), nil
}

func SizeToBytes(val int) []byte {
	int32Val := UInt32(val)
	return int32Val.Bytes()
}
