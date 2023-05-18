package key

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"strconv"
)

type Era uint64

func (k Era) MarshalJSON() ([]byte, error) {
	intVal := uint64(k)
	return json.Marshal(intVal)
}

func (k *Era) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	tmp, err := NewEraFromString(s)
	if err != nil {
		return err
	}
	*k = tmp

	return nil
}

func (k Era) Bytes() []byte {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, uint64(k))
	return data
}

func NewEra(val uint64) *Era {
	era := Era(val)
	return &era
}

func NewEraFromString(source string) (res Era, err error) {
	intVal, err := strconv.Atoi(source)
	if err != nil {
		return res, err
	}
	return Era(uint32(intVal)), nil
}

func NewEraFromBuffer(buf *bytes.Buffer) (*Era, error) {
	u := binary.LittleEndian.Uint64(buf.Bytes())
	return NewEra(u), nil
}
