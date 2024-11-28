package key

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
)

const ByteHashLen = 32
const StingHashLen = 64

type Hash [ByteHashLen]byte

func (h Hash) Value() (driver.Value, error) {
	return h[:], nil
}

func (h *Hash) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("invalid scan value type")
	}
	if len(b) != len(h) {
		return fmt.Errorf("invalid hash length, expected %d but got %d", len(h), len(b))
	}
	copy(h[:], b)
	return nil
}

func (h *Hash) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	hexBytes, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	copy(h[:], hexBytes)
	return nil
}

func (h *Hash) UnmarshalText(text []byte) error {
	tmp, err := NewHash(string(text))
	if err != nil {
		return err
	}
	*h = tmp
	return nil
}

func (h Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.ToHex())
}

func (h Hash) Bytes() []byte {
	return h[:]
}

func (h Hash) ToHex() string {
	return hex.EncodeToString(h[:])
}

func (h Hash) String() string {
	return h.ToHex()
}

func NewHash(source string) (result Hash, err error) {
	if len(source) != StingHashLen {
		return result, fmt.Errorf("can't parse string to key, source: %s", source)
	}
	hexBytes, err := hex.DecodeString(source)
	if err != nil {
		return result, err
	}
	if len(hexBytes) != len(result) {
		return result, fmt.Errorf("invalid hash length, expected %d but got %d", len(result), len(hexBytes))
	}
	copy(result[:], hexBytes)

	return result, nil
}

func NewHashFromBytes(source []byte) (result Hash, err error) {
	return NewByteHashFromBuffer(bytes.NewBuffer(source))
}

func NewByteHashFromBuffer(buf *bytes.Buffer) (result Hash, err error) {
	if buf.Len() < ByteHashLen {
		return result, errors.New("key length is not equal 32")
	}

	copy(result[:], buf.Next(ByteHashLen))
	return result, nil
}

type HashFromBytesDecoder struct{}

func (addr *HashFromBytesDecoder) FromBytes(bytes []byte) (Hash, []byte, error) {
	panic("implement me")
}
