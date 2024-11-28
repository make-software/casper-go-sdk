package encoding

const U8SerializedLength = 1

type U8FromBytesDecoder struct{}

func (addr *U8FromBytesDecoder) FromBytes(inputBytes []byte) (uint8, []byte, error) {
	if len(inputBytes) == 0 {
		return 0, nil, ErrEmptyBytesSource
	}

	byteVal := inputBytes[0]
	remainder := inputBytes[1:]

	return byteVal, remainder, nil
}
