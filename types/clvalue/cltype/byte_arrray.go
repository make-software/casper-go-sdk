package cltype

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
)

type ByteArray struct {
	Size uint32
}

func (b *ByteArray) Bytes() []byte {
	sizeInBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sizeInBytes, b.Size)
	return append([]byte{b.GetTypeID()}, sizeInBytes...)
}

func (b *ByteArray) String() string {
	return fmt.Sprintf("%s: %d", b.Name(), b.Size)
}

func (b *ByteArray) GetTypeID() TypeID {
	return TypeIDByteArray
}

func (b *ByteArray) Name() TypeName {
	return TypeNameByteArray
}

func (b *ByteArray) Len() int {
	return int(b.Size)
}

func (b ByteArray) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]uint32{b.Name(): b.Size})
}

func NewByteArray(size uint32) *ByteArray {
	return &ByteArray{Size: size}
}

func NewByteArrayFromJson(source interface{}) (*ByteArray, error) {
	val, ok := source.(float64)
	if !ok {
		return nil, errors.New("invalid json parsing to ByteArray type")
	}
	return NewByteArray(uint32(val)), nil
}

func NewByteArrayFromBuffer(buf *bytes.Buffer) (*ByteArray, error) {
	if buf.Len() < Int32ByteSize {
		return nil, fmt.Errorf("buffer len is less then size = %d", Int32ByteSize)
	}
	return NewByteArray(binary.LittleEndian.Uint32(buf.Next(Int32ByteSize))), nil
}
