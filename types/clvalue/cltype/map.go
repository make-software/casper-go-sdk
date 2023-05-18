package cltype

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

var MapJsonParsingError = errors.New("map type parsing error")

type Map struct {
	Key CLType
	Val CLType
}

func (m *Map) Bytes() []byte {
	return append([]byte{TypeIDMap}, append(m.Key.Bytes(), m.Val.Bytes()...)...)
}

func (m *Map) String() string {
	return fmt.Sprintf("%s (%s: %s)", TypeNameMap, m.Key.String(), m.Val.String())
}

func (m *Map) GetTypeID() TypeID {
	return TypeIDMap
}

func (m *Map) Name() TypeName {
	return TypeNameMap
}

func (m *Map) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]CLType{m.Key.String(): m.Val})
}

func NewMap(keyType CLType, valType CLType) *Map {
	return &Map{
		Key: keyType,
		Val: valType,
	}
}

func NewMapFromJson(source interface{}) (*Map, error) {
	mapData, ok := source.(map[string]interface{})
	if !ok {
		return nil, MapJsonParsingError
	}

	key, found := mapData["key"]
	if !found {
		return nil, MapJsonParsingError
	}
	keyType, err := fromInterface(key)
	if err != nil {
		return nil, err
	}

	val, found := mapData["value"]
	if !found {
		return nil, MapJsonParsingError
	}
	valType, err := fromInterface(val)
	if err != nil {
		return nil, err
	}

	return &Map{Key: keyType, Val: valType}, nil
}

func NewMapFromBytes(source []byte) (*Map, error) {
	buf := bytes.NewBuffer(source)
	return NewMapFromBuffer(buf)
}

func NewMapFromBuffer(buf *bytes.Buffer) (*Map, error) {
	key, err := FromBuffer(buf)
	if err != nil {
		return nil, err
	}
	value, err := FromBuffer(buf)
	if err != nil {
		return nil, err
	}
	return &Map{
		Key: key,
		Val: value,
	}, nil
}
