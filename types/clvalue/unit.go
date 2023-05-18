package clvalue

import (
	"bytes"
	"errors"
	"go/types"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

var ErrUnitByteNotEmpty = errors.New("byte source for unit type should be empty")

type Unit struct {
	obj *types.Nil
}

func (v *Unit) Bytes() []byte {
	return []byte{}
}

func (v *Unit) String() string {
	return v.obj.String()
}

func (v *Unit) Value() *types.Nil {
	return v.obj
}

func NewCLUnit() *CLValue {
	res := CLValue{}
	res.Type = cltype.Unit
	v := Unit{obj: new(types.Nil)}
	res.Unit = &v
	return &res
}

func NewUnitFromBytes(source []byte) (*Unit, error) {
	return NewUnitFromBuffer(bytes.NewBuffer(source))
}

func NewUnitFromBuffer(buffer *bytes.Buffer) (*Unit, error) {
	if buffer.Len() > 0 {
		return nil, ErrUnitByteNotEmpty
	}
	v := Unit{obj: new(types.Nil)}
	return &v, nil
}
