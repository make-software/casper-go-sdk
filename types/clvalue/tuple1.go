package clvalue

import (
	"bytes"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

type Tuple1 struct {
	innerType *cltype.Tuple1
	innerVal  CLValue
}

func (v *Tuple1) Bytes() []byte {
	return v.innerVal.Bytes()
}

func (v *Tuple1) String() string {
	return "(" + v.innerVal.String() + ")"
}
func (v *Tuple1) Value() CLValue {
	return v.innerVal
}

func NewCLTuple1(val CLValue) CLValue {
	tupleType := cltype.NewTuple1(val.Type)
	return CLValue{Type: tupleType, Tuple1: &Tuple1{innerType: tupleType, innerVal: val}}
}

func NewTuple1FromBytes(source []byte, clType *cltype.Tuple1) (*Tuple1, error) {
	return NewTuple1FromBuffer(bytes.NewBuffer(source), clType)
}

func NewTuple1FromBuffer(buf *bytes.Buffer, tuple1Type *cltype.Tuple1) (*Tuple1, error) {
	inner, err := FromBufferByType(buf, tuple1Type.Inner)
	if err != nil {
		return nil, err
	}
	return &Tuple1{innerVal: inner, innerType: tuple1Type}, nil
}
