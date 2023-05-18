package cltype

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

type Tuple3 struct {
	Inner1 CLType
	Inner2 CLType
	Inner3 CLType
}

func (t *Tuple3) Bytes() []byte {
	return append([]byte{t.GetTypeID()}, append(t.Inner1.Bytes(), append(t.Inner2.Bytes(), t.Inner3.Bytes()...)...)...)
}

func (t *Tuple3) String() string {
	return fmt.Sprintf("%s (%s, %s, %s)", t.Name(), t.Inner1.String(), t.Inner2.String(), t.Inner3.String())
}

func (t *Tuple3) GetTypeID() TypeID {
	return TypeIDTuple3
}

func (t *Tuple3) Name() TypeName {
	return TypeNameTuple3
}

func (t *Tuple3) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string][]CLType{t.Name(): []CLType{t.Inner1, t.Inner2, t.Inner3}})
}

func NewTuple3(inner1, inner2, inner3 CLType) *Tuple3 {
	return &Tuple3{Inner1: inner1, Inner2: inner2, Inner3: inner3}
}

func NewTuple3FromJson(source interface{}) (*Tuple3, error) {
	if data, ok := source.([]interface{}); ok {
		if len(data) != 3 {
			return nil, errors.New("invalid tuple1 type format, should be array of 3 element")
		}
		inner1, err := fromInterface(data[0])
		if err != nil {
			return nil, err
		}
		inner2, err := fromInterface(data[1])
		if err != nil {
			return nil, err
		}
		inner3, err := fromInterface(data[2])
		if err != nil {
			return nil, err
		}
		return NewTuple3(inner1, inner2, inner3), nil
	}
	return nil, errors.New("invalid tuple1 type format, should be array of 3 element")
}

func NewTuple3FromBuffer(buf *bytes.Buffer) (*Tuple3, error) {
	inner1, err := FromBuffer(buf)
	if err != nil {
		return nil, err
	}
	inner2, err := FromBuffer(buf)
	if err != nil {
		return nil, err
	}
	inner3, err := FromBuffer(buf)
	if err != nil {
		return nil, err
	}
	return NewTuple3(inner1, inner2, inner3), nil
}
