package cltype

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

type Tuple1 struct {
	Inner CLType
}

func (t *Tuple1) Bytes() []byte {
	return append([]byte{t.GetTypeID()}, t.Inner.Bytes()...)
}

func (t *Tuple1) String() string {
	return fmt.Sprintf("%s (%s)", t.Name(), t.Inner.String())
}

func (t *Tuple1) GetTypeID() TypeID {
	return TypeIDTuple1
}

func (t *Tuple1) Name() TypeName {
	return TypeNameTuple1
}

func (t *Tuple1) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string][]CLType{t.Name(): []CLType{t.Inner}})
}

func NewTuple1(inner CLType) *Tuple1 {
	return &Tuple1{Inner: inner}
}

func NewTuple1FromJson(source interface{}) (*Tuple1, error) {
	if data, ok := source.([]interface{}); ok {
		if len(data) != 1 {
			return nil, errors.New("invalid tuple1 type format, should be array of 1 element")
		}
		inner, err := fromInterface(data[0])
		if err != nil {
			return nil, err
		}
		return NewTuple1(inner), nil
	}
	return nil, errors.New("invalid tuple1 type format, should be array of 1 element")
}

func NewTuple1FromBuffer(buf *bytes.Buffer) (*Tuple1, error) {
	inner, err := FromBuffer(buf)
	if err != nil {
		return nil, err
	}
	return NewTuple1(inner), nil
}
