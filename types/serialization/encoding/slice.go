package encoding

type SliceFromBytesDecoder[T any, D FromBytes[T]] struct {
	Decoder D
}

// FromBytes decodes a slice of type T using the FromBytes interface
func (d *SliceFromBytesDecoder[T, D]) FromBytes(inputBytes []byte) ([]T, []byte, error) {
	if len(inputBytes) < 4 {
		return nil, nil, ErrInvalidBytesStructure
	}

	count, remainder, err := NewU32FromBytesDecoder().FromBytes(inputBytes)
	if err != nil {
		return nil, nil, err
	}

	if count == 0 {
		return nil, remainder, nil
	}

	result := make([]T, 0, count)
	for i := uint32(0); i < count; i++ {
		elem, rem, err := d.Decoder.FromBytes(remainder)
		if err != nil {
			return nil, nil, err
		}
		result = append(result, elem)
		remainder = rem
	}

	return result, remainder, nil
}

type SliceToBytesEncoder[E ToBytes] struct {
	values []E
}

func NewSliceToBytesEncoder[E ToBytes](values []E) *SliceToBytesEncoder[E] {
	return &SliceToBytesEncoder[E]{
		values,
	}
}

func (enc *SliceToBytesEncoder[E]) Bytes() ([]byte, error) {
	var estimatedSize int
	for _, el := range enc.values {
		estimatedSize += el.SerializedLength()
	}
	result := make([]byte, 0, estimatedSize)

	lengthBytes, err := NewU32ToBytesEncoder(uint32(len(enc.values))).Bytes()
	if err != nil {
		return nil, err
	}

	result = append(result, lengthBytes...)
	for _, el := range enc.values {
		elBytes, err := el.Bytes()
		if err != nil {
			return nil, err
		}
		result = append(result, elBytes...)
	}

	return result, nil
}
