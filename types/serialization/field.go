package serialization

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"sort"

	"github.com/make-software/casper-go-sdk/v2/types/serialization/encoding"
)

const FieldSerializedLength = 6

// Field represents an individual field in the serialization envelope.
type Field struct {
	Index  uint16
	Offset uint32
}

func (f Field) SerializedLength() int {
	return FieldSerializedLength
}

func (f Field) Bytes() ([]byte, error) {
	buffer := make([]byte, 6) // 2 bytes for Index + 4 bytes for Offset
	binary.BigEndian.PutUint16(buffer[:2], f.Index)
	binary.BigEndian.PutUint32(buffer[2:], f.Offset)
	return buffer, nil
}

type FieldFromBytesDecoder struct{}

// FromBytes function to deserialize a Field from bytes
func (addr *FieldFromBytesDecoder) FromBytes(inputBytes []byte) (Field, []byte, error) {
	// Ensure the input bytes are long enough for the u16 index and u32 offset
	if len(inputBytes) < 6 {
		return Field{}, nil, ErrEarlyEndOfStream
	}

	u16Decoder := encoding.U16FromBytesDecoder{}
	index, reminder, err := u16Decoder.FromBytes(inputBytes)
	if err != nil {
		return Field{}, nil, err
	}

	u32Decoder := encoding.U32FromBytesDecoder{}
	offset, reminder, err := u32Decoder.FromBytes(reminder)
	if err != nil {
		return Field{}, nil, err
	}

	field := Field{
		Index:  index,
		Offset: offset,
	}

	return field, reminder, nil
}

type Fields map[uint16][]byte

func NewFields() Fields {
	return make(map[uint16][]byte)
}

func (f *Fields) AddField(key uint16, value encoding.ToBytes) error {
	bytesData, err := value.Bytes()
	if err != nil {
		return err
	}
	(*f)[key] = bytesData
	return nil
}

func (t *Fields) UnmarshalJSON(data []byte) error {
	var temp struct {
		Fields map[uint16]string `json:"fields"`
	}
	if err := json.Unmarshal(data, &temp.Fields); err != nil {
		return err
	}

	res := make(map[uint16][]byte, len(temp.Fields))

	for key, value := range temp.Fields {
		decoded, err := hex.DecodeString(value)
		if err != nil {
			return err
		}
		res[key] = decoded
	}
	*t = res
	return nil
}

// Bytes serializes all fields in the slice into a single byte slice.
func (f *Fields) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)

	numItems := int32(len(*f))
	if err := binary.Write(buf, binary.LittleEndian, numItems); err != nil {
		return nil, err
	}

	keys := make([]uint16, 0, len(*f))
	for key := range *f {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, key := range keys {
		data := (*f)[key]
		_, err := buf.Write(data)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (f *Fields) SerializedLength() int {
	length := encoding.U32SerializedLength
	for range *f {
		length += encoding.U16SerializedLength // key u16
		length += Field{}.SerializedLength()   // field
	}

	return length
}
