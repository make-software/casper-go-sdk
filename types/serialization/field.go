package serialization

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"

	"github.com/make-software/casper-go-sdk/v2/types/serialization/encoding"
)

const FieldSerializedLength = 6

type Fields map[uint16][]byte

// Field represents an individual field in the serialization envelope.
type Field struct {
	Index  uint16
	Offset uint32
}

func NewFields() Fields {
	return make(map[uint16][]byte)
}

func (f Field) SerializedLength() int {
	return FieldSerializedLength
}

func (f Field) Bytes() ([]byte, error) {
	buffer := make([]byte, 6) // 2 bytes for Index + 4 bytes for Offset
	binary.LittleEndian.PutUint16(buffer[:2], f.Index)
	binary.LittleEndian.PutUint32(buffer[2:], f.Offset)
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

	for key := 0; key < len(*f); key++ {
		data := (*f)[uint16(key)]

		// encode key
		if err := binary.Write(buf, binary.LittleEndian, uint16(key)); err != nil {
			return nil, err
		}

		// encode bytes length
		if err := binary.Write(buf, binary.LittleEndian, uint32(len(data))); err != nil {
			return nil, err
		}

		// encode data
		if _, err := buf.Write(data); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (f *Fields) SerializedLength() int {
	length := encoding.U32SerializedLength
	for _, value := range *f {
		length += encoding.U16SerializedLength // key u16
		length += encoding.U32SerializedLength // key data length u32
		length += len(value)
	}

	return length
}
