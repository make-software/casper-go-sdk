package serialization

import (
	"encoding/binary"
	"errors"

	"github.com/make-software/casper-go-sdk/v2/types/serialization/encoding"
)

var (
	ErrFormatting       = errors.New("formatting error")
	ErrEarlyEndOfStream = errors.New("early end of stream")
)

// CallTableSerializationEnvelopeBuilder is used to construct a serialization envelope.
type CallTableSerializationEnvelopeBuilder struct {
	fields               []Field
	expectedPayloadSizes []int
	bytes                []byte
	currentFieldIndex    int
	currentOffset        int
}

func NewCallTableSerializationEnvelopeBuilder(expectedPayloadSizes []int) (*CallTableSerializationEnvelopeBuilder, error) {
	numberOfFields := len(expectedPayloadSizes)
	fieldsSize := serializedVecSize(numberOfFields)

	bytesOfPayloadSize := sum(expectedPayloadSizes)
	payloadAndVecOverhead := encoding.U32SerializedLength + bytesOfPayloadSize

	buffer := make([]byte, 0, fieldsSize+payloadAndVecOverhead)
	buffer = append(buffer, make([]byte, fieldsSize)...) // Making room for the call table

	uint32Buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(uint32Buf, uint32(bytesOfPayloadSize))

	buffer = append(buffer, uint32Buf...)

	return &CallTableSerializationEnvelopeBuilder{
		fields:               make([]Field, 0, numberOfFields),
		expectedPayloadSizes: expectedPayloadSizes,
		bytes:                buffer,
		currentFieldIndex:    0,
		currentOffset:        0,
	}, nil
}

func (b *CallTableSerializationEnvelopeBuilder) AddField(fieldIndex uint16, value []byte) error {
	if b.currentFieldIndex >= len(b.expectedPayloadSizes) {
		return errors.New("more fields than expected")
	}
	if len(b.fields) > 0 && b.fields[b.currentFieldIndex-1].Index >= fieldIndex {
		return errors.New("fields must be in ascending order")
	}

	size := b.expectedPayloadSizes[b.currentFieldIndex]

	if len(value) == 0 {
		return errors.New("empty fields are not allowed")
	}

	if len(value) != size {
		return errors.New("field size does not match expected size")
	}

	b.bytes = append(b.bytes, value...)
	b.fields = append(b.fields, Field{
		Index:  fieldIndex,
		Offset: uint32(b.currentOffset),
	})

	b.currentFieldIndex++
	b.currentOffset += size
	return nil
}

func (b *CallTableSerializationEnvelopeBuilder) BinaryPayloadBytes() ([]byte, error) {
	if b.currentFieldIndex != len(b.expectedPayloadSizes) {
		return nil, errors.New("not all fields have been added")
	}

	fieldBytes, err := encoding.NewSliceToBytesEncoder(b.fields).Bytes()
	if err != nil {
		return nil, err
	}

	copy(b.bytes[:len(fieldBytes)], fieldBytes)
	return b.bytes, nil
}

func sum(slice []int) int {
	total := 0
	for _, val := range slice {
		total += val
	}
	return total
}

func serializedVecSize(numberOfFields int) int {
	return encoding.U32SerializedLength + numberOfFields*FieldSerializedLength
}

type CallTableSerializationEnvelope struct {
	Fields []Field
	Bytes  []byte
}

// FromBytes deserializes the envelope
func (env *CallTableSerializationEnvelope) FromBytes(maxExpectedFields uint32, inputBytes []byte) (*CallTableSerializationEnvelope, []byte, error) {
	if len(inputBytes) < 4 {
		// The first "thing" in the bytes of the payload should be a `fields` vector.
		// Check the number of entries in that vector to avoid field pumping.
		// If the payload doesn't have a u32 size of bytes in it, then it's malformed.
		return nil, nil, ErrFormatting
	}

	// Deserialize the number of fields (u32)
	u32Decoder := encoding.U32FromBytesDecoder{}
	numberOfFields, _, err := u32Decoder.FromBytes(inputBytes)
	if err != nil {
		return nil, nil, err
	}

	if numberOfFields > maxExpectedFields {
		return nil, nil, ErrFormatting
	}

	fieldsDecoder := &encoding.SliceFromBytesDecoder[Field, *FieldFromBytesDecoder]{
		Decoder: &FieldFromBytesDecoder{},
	}
	fields, remainder, err := fieldsDecoder.FromBytes(inputBytes)
	if err != nil {
		return nil, nil, err
	}

	bytesDecoder := &encoding.SliceFromBytesDecoder[uint8, *encoding.U8FromBytesDecoder]{
		Decoder: &encoding.U8FromBytesDecoder{},
	}
	bytes, remainder, err := bytesDecoder.FromBytes(remainder)
	if err != nil {
		return nil, nil, err
	}
	return &CallTableSerializationEnvelope{Fields: fields, Bytes: bytes}, remainder, nil
}

func (env *CallTableSerializationEnvelope) EstimateSize(fieldSizes []int) int {
	numberOfFields := len(fieldSizes)
	payloadBytes := 0
	for _, size := range fieldSizes {
		payloadBytes += size
	}

	size := encoding.U32SerializedLength + encoding.U32SerializedLength // Overhead for fields vec and bytes vec
	size += numberOfFields * 6                                          // Each Field has 6 bytes (2 for Index, 4 for Offset)
	size += payloadBytes
	return size
}

func (env *CallTableSerializationEnvelope) StartConsuming() (*CallTableFieldsIterator, error) {
	if len(env.Fields) == 0 {
		return nil, nil
	}

	field := env.Fields[0]
	expectedSize := len(env.Bytes)
	if len(env.Fields) > 1 {
		expectedSize = int(env.Fields[1].Offset)
	}

	return &CallTableFieldsIterator{
		IndexInFieldsVec: 0,
		ExpectedSize:     expectedSize,
		Field:            &field,
		Bytes:            env.Bytes,
		Parent:           env,
	}, nil
}

type CallTableFieldsIterator struct {
	IndexInFieldsVec int
	ExpectedSize     int
	Field            *Field
	Bytes            []byte
	Parent           *CallTableSerializationEnvelope
}

// VerifyIndex verifies the index of the current field
func (it *CallTableFieldsIterator) VerifyIndex(expectedIndex uint16) error {
	if it.Field.Index != expectedIndex {
		return ErrFormatting
	}
	return nil
}

// DeserializeAndMaybeNext deserializes the current field and returns the next iterator if available
func DeserializeAndMaybeNext[T any, D encoding.FromBytes[T]](it *CallTableFieldsIterator, decoder D) (T, *CallTableFieldsIterator, error) {
	data, nextIt, err := step(it, decoder)
	return data, nextIt, err
}

// step deserializes the current field and prepares the next iterator
func step[T any, D encoding.FromBytes[T]](it *CallTableFieldsIterator, decoder D) (T, *CallTableFieldsIterator, error) {
	data, remainder, err := decoder.FromBytes(it.Bytes)
	if err != nil {
		var zero T
		return zero, nil, err
	}

	isLastField := it.IndexInFieldsVec == len(it.Parent.Fields)-1
	if len(remainder)+it.ExpectedSize != len(it.Bytes) {
		var zero T
		return zero, nil, ErrFormatting
	}

	// if not the last field, prepare the next iterator
	if !isLastField {
		nextIndex := it.IndexInFieldsVec + 1
		nextField := &it.Parent.Fields[nextIndex]
		isNextFieldLast := nextIndex == len(it.Parent.Fields)-1

		expectedSize := len(remainder)
		if !isNextFieldLast {
			expectedSize = int(it.Parent.Fields[nextIndex+1].Offset - nextField.Offset)
		}

		nextIterator := &CallTableFieldsIterator{
			IndexInFieldsVec: nextIndex,
			ExpectedSize:     expectedSize,
			Field:            nextField,
			Bytes:            remainder,
			Parent:           it.Parent,
		}

		return data, nextIterator, nil
	}

	if len(remainder) != 0 {
		var zero T
		return zero, nil, ErrFormatting
	}
	return data, nil, nil
}
