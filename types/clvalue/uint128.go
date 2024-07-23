package clvalue

import (
	"bytes"
	"math/big"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
)

type UInt128 struct {
	val *big.Int
}

func (v *UInt128) Bytes() []byte {
	return BigToBytes(v.val)
}

func (v *UInt128) String() string {
	return v.val.String()
}

func (v *UInt128) Value() *big.Int {
	return v.val
}

func NewCLUInt128(val *big.Int) *CLValue {
	res := CLValue{}
	res.Type = cltype.UInt128
	v := UInt128{val: val}
	res.UI128 = &v
	return &res
}

func NewUint128FromBytes(source []byte) (*UInt128, error) {
	buf := bytes.NewBuffer(source)
	return NewUint128FromBuffer(buf)
}

func NewUint128FromBuffer(buffer *bytes.Buffer) (*UInt128, error) {
	val, err := BigFromBuffer(buffer)
	v := UInt128{val: val}
	return &v, err
}
