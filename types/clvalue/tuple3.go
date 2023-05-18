package clvalue

import (
	"bytes"
	"fmt"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

type Tuple3 struct {
	innerType *cltype.Tuple3
	Inner1    CLValue
	Inner2    CLValue
	Inner3    CLValue
}

func (v *Tuple3) Bytes() []byte {
	return append(v.Inner1.Bytes(), append(v.Inner2.Bytes(), v.Inner3.Bytes()...)...)
}

func (v *Tuple3) String() string {
	return fmt.Sprintf("(%s, %s, %s)", v.Inner1.String(), v.Inner2.String(), v.Inner3.String())
}
func (v *Tuple3) Value() [3]CLValue {
	return [3]CLValue{v.Inner1, v.Inner2, v.Inner3}
}

func NewCLTuple3(val1 CLValue, val2 CLValue, val3 CLValue) CLValue {
	tupleType := cltype.NewTuple3(val1.Type, val2.Type, val3.Type)
	return CLValue{Type: tupleType, Tuple3: &Tuple3{
		innerType: tupleType,
		Inner1:    val1,
		Inner2:    val2,
		Inner3:    val2,
	}}
}

func NewTuple3FromBytes(source []byte, clType *cltype.Tuple3) (*Tuple3, error) {
	return NewTuple3FromBuffer(bytes.NewBuffer(source), clType)
}

func NewTuple3FromBuffer(buf *bytes.Buffer, clType *cltype.Tuple3) (*Tuple3, error) {
	inner1, err := FromBufferByType(buf, clType.Inner1)
	if err != nil {
		return nil, err
	}
	inner2, err := FromBufferByType(buf, clType.Inner2)
	if err != nil {
		return nil, err
	}
	inner3, err := FromBufferByType(buf, clType.Inner3)
	if err != nil {
		return nil, err
	}

	return &Tuple3{innerType: clType, Inner1: inner1, Inner2: inner2, Inner3: inner3}, nil
}
