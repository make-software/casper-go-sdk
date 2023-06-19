package clvalue

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

type Int64 int64

func (v *Int64) Bytes() []byte {
	data := make([]byte, cltype.Int64ByteSize)
	binary.LittleEndian.PutUint64(data, uint64(*v))
	return data
}

func (v *Int64) String() string {
	return fmt.Sprintf("%d", *v)
}

func (v *Int64) Value() int64 {
	return int64(*v)
}

func NewCLInt64(val int64) *CLValue {
	res := CLValue{}
	res.Type = cltype.Int64
	v := Int64(val)
	res.I64 = &v
	return &res
}

func NewInt64FromBytes(source []byte) (*Int64, error) {
	return NewInt64FromBuffer(bytes.NewBuffer(source))
}

func NewInt64FromBuffer(buf *bytes.Buffer) (*Int64, error) {
	if buf.Len() < cltype.Int64ByteSize {
		return nil, errors.New("buffer size is too small")
	}
	byteSlice := buf.Next(cltype.Int64ByteSize)
	val := Int64(int64(binary.LittleEndian.Uint64(byteSlice)))
	return &val, nil
}
