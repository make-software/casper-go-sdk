package clvalue

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

type UInt64 uint64

func (v *UInt64) Bytes() []byte {
	data := make([]byte, cltype.Int64ByteSize)
	binary.LittleEndian.PutUint64(data, uint64(*v))
	return data
}

func (v *UInt64) String() string {
	return fmt.Sprintf("%d", *v)
}

func (v *UInt64) Value() uint64 {
	return uint64(*v)
}

func NewCLUInt64(val uint64) *CLValue {
	res := CLValue{}
	res.Type = cltype.UInt64
	v := UInt64(val)
	res.UI64 = &v
	return &res
}

func NewUint64FromBytes(source []byte) (*UInt64, error) {
	buf := bytes.NewBuffer(source)
	return NewUint64FromBuffer(buf)
}

func NewUint64FromBuffer(buffer *bytes.Buffer) (*UInt64, error) {
	if buffer.Len() < cltype.Int32ByteSize {
		return nil, errors.New("buffer size is too small")
	}

	buf := buffer.Next(cltype.Int64ByteSize)
	val := UInt64(binary.LittleEndian.Uint64(buf))
	return &val, nil
}
