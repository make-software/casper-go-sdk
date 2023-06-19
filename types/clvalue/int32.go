package clvalue

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

type Int32 int32

func (v *Int32) Bytes() []byte {
	data := make([]byte, cltype.Int32ByteSize)
	binary.LittleEndian.PutUint32(data, uint32(*v))
	return data
}

func (v *Int32) String() string {
	return fmt.Sprintf("%d", *v)
}

func (v *Int32) Value() int32 {
	return int32(*v)
}

func NewCLInt32(val int32) CLValue {
	res := CLValue{}
	res.Type = cltype.Int32
	v := Int32(val)
	res.I32 = &v
	return res
}

func NewInt32FromBytes(source []byte) (*Int32, error) {
	buf := bytes.NewBuffer(source)
	return NewInt32FromBuffer(buf)
}

func NewInt32FromBuffer(buffer *bytes.Buffer) (*Int32, error) {
	if buffer.Len() < cltype.Int32ByteSize {
		return nil, errors.New("buffer size is too small")
	}
	buf := buffer.Next(cltype.Int32ByteSize)
	val := Int32(int32(binary.LittleEndian.Uint32(buf)))
	return &val, nil
}
