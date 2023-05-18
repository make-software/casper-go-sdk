package clvalue

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

type Map struct {
	Type        *cltype.Map
	data        []Tuple2
	indexedData map[string]CLValue
}

func (m *Map) Bytes() []byte {
	result := bytes.NewBuffer([]byte{})
	result.Write(SizeToBytes(len(m.indexedData)))
	for _, val := range m.data {
		result.Write(val.Inner1.Bytes())
		result.Write(val.Inner2.Bytes())
	}
	return result.Bytes()
}

func (m *Map) Map() map[string]CLValue {
	result := make(map[string]CLValue, len(m.data))
	for k, v := range m.indexedData {
		result[k] = v
	}
	return result
}

func (m *Map) String() string {
	b := new(bytes.Buffer)
	for key, value := range m.indexedData {
		fmt.Fprintf(b, "(%s=\"%s\")", key, value.String())
	}
	return b.String()
}

func (m *Map) Find(key string) (CLValue, bool) {
	res, ok := m.indexedData[key]
	return res, ok
}

func (m *Map) Get(key string) CLValue {
	return m.indexedData[key]
}

func (m *Map) FindAny(keys []string) (CLValue, bool) {
	for _, key := range keys {
		if value, ok := m.indexedData[key]; ok {
			return value, true
		}
	}

	return CLValue{}, false
}

func (m *Map) Len() int {
	return len(m.indexedData)
}

func (m *Map) Append(key CLValue, val CLValue) error {
	if key.Type != m.Type.Key {
		return errors.New("invalid key type")
	}
	if val.Type != m.Type.Val {
		return errors.New("invalid value type")
	}
	m.data = append(m.data, *NewCLTuple2(key, val).Tuple2)
	if _, found := m.indexedData[key.String()]; found {
		return errors.New("map key is already exist")
	}
	m.indexedData[key.String()] = val
	return nil
}

func NewCLMap(keyType cltype.CLType, valType cltype.CLType) CLValue {
	mapType := cltype.NewMap(keyType, valType)
	return CLValue{
		Type: mapType,
		Map:  newMap(mapType),
	}
}

func newMap(mapType *cltype.Map) *Map {
	return &Map{
		Type:        mapType,
		data:        make([]Tuple2, 0),
		indexedData: make(map[string]CLValue),
	}
}

func NewMapFromBuffer(buffer *bytes.Buffer, mapType *cltype.Map) (*Map, error) {
	result := newMap(mapType)
	var KeyVal CLValue
	var ValVal CLValue
	var err error
	size := TrimByteSize(buffer)
	for i := uint32(0); i < size; i++ {
		if KeyVal, err = FromBufferByType(buffer, mapType.Key); err != nil {
			return nil, err
		}
		if ValVal, err = FromBufferByType(buffer, mapType.Val); err != nil {
			return nil, err
		}

		err = result.Append(KeyVal, ValVal)
		if err != nil {
			return nil, err
		}
	}
	return result, err
}
