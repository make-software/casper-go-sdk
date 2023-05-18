package cltype

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

type Tuple2 struct {
	Inner1 CLType
	Inner2 CLType
}

func (t *Tuple2) Bytes() []byte {
	return append([]byte{t.GetTypeID()}, append(t.Inner1.Bytes(), t.Inner2.Bytes()...)...)
}

func (t *Tuple2) String() string {
	return fmt.Sprintf("%s (%s, %s)", t.Name(), t.Inner1.String(), t.Inner2.String())
}

func (t *Tuple2) GetTypeID() TypeID {
	return TypeIDTuple2
}

func (t *Tuple2) Name() TypeName {
	return TypeNameTuple2
}

func (t *Tuple2) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string][]CLType{t.Name(): []CLType{t.Inner1, t.Inner2}})
}

func NewTuple2(inner1, inner2 CLType) *Tuple2 {
	return &Tuple2{Inner1: inner1, Inner2: inner2}
}

func NewTuple2FromJson(source interface{}) (*Tuple2, error) {
	if data, ok := source.([]interface{}); ok {
		if len(data) != 2 {
			return nil, errors.New("invalid tuple1 type format, should be array of 2 element")
		}
		inner1, err := fromInterface(data[0])
		if err != nil {
			return nil, err
		}
		inner2, err := fromInterface(data[1])
		if err != nil {
			return nil, err
		}
		return NewTuple2(inner1, inner2), nil
	}
	return nil, errors.New("invalid tuple1 type format, should be array of 2 element")
}

func NewTuple2FromBuffer(buf *bytes.Buffer) (*Tuple2, error) {
	inner1, err := FromBuffer(buf)
	if err != nil {
		return nil, err
	}
	inner2, err := FromBuffer(buf)
	if err != nil {
		return nil, err
	}
	return NewTuple2(inner1, inner2), nil
}
