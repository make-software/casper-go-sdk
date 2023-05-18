package clvalue

import (
	"bytes"
	"fmt"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

type Tuple2 struct {
	innerType *cltype.Tuple2
	Inner1    CLValue
	Inner2    CLValue
}

func (v *Tuple2) Bytes() []byte {
	return append(v.Inner1.Bytes(), v.Inner2.Bytes()...)
}

func (v *Tuple2) String() string {
	return fmt.Sprintf("(%s, %s)", v.Inner1.String(), v.Inner2.String())
}
func (v *Tuple2) Value() [2]CLValue {
	return [2]CLValue{v.Inner1, v.Inner2}
}

func NewCLTuple2(val1 CLValue, val2 CLValue) CLValue {
	tupleType := cltype.NewTuple2(val1.Type, val2.Type)
	return CLValue{Type: tupleType, Tuple2: &Tuple2{
		innerType: tupleType,
		Inner1:    val1,
		Inner2:    val2,
	}}
}

func NewTuple2FromBytes(source []byte, clType *cltype.Tuple2) (*Tuple2, error) {
	return NewTuple2FromBuffer(bytes.NewBuffer(source), clType)
}

func NewTuple2FromBuffer(buf *bytes.Buffer, clType *cltype.Tuple2) (*Tuple2, error) {
	inner1, err := FromBufferByType(buf, clType.Inner1)
	if err != nil {
		return nil, err
	}
	inner2, err := FromBufferByType(buf, clType.Inner2)
	if err != nil {
		return nil, err
	}

	return &Tuple2{innerType: clType, Inner1: inner1, Inner2: inner2}, nil
}
