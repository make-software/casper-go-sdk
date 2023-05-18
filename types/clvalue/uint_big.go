package clvalue

import (
	"bytes"
	"math/big"
)

func BigToBytes(val *big.Int) []byte {
	data := val.Bytes()
	for i := 0; i < len(data)/2; i++ {
		data[i], data[len(data)-i-1] = data[len(data)-i-1], data[i]
	}
	numberLen := byte(len(data))
	return append([]byte{numberLen}, data...)
}

func BigFromBuffer(buffer *bytes.Buffer) (*big.Int, error) {
	//Todo: check size validation for certain type (U128 <= 16 bytes)
	size, err := buffer.ReadByte()
	if err != nil {
		return nil, err
	}
	data := buffer.Next(int(size))
	numBytes := make([]byte, size)
	for i, b := range data {
		numBytes[len(numBytes)-i-1] = b
	}

	res := new(big.Int)
	res.SetBytes(numBytes[:])
	return res, err
}
