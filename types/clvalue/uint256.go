package clvalue

import (
	"bytes"
	"math/big"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
)

type UInt256 struct {
	val *big.Int
}

func (v *UInt256) Bytes() []byte {
	return BigToBytes(v.val)
}

func (v *UInt256) String() string {
	return v.val.String()
}

func (v *UInt256) Value() *big.Int {
	return v.val
}

func NewCLUInt256(val *big.Int) *CLValue {
	res := CLValue{}
	res.Type = cltype.UInt256
	v := UInt256{val: val}
	res.UI256 = &v
	return &res
}

func NewUint256FromBytes(source []byte) (*UInt256, error) {
	buf := bytes.NewBuffer(source)
	return NewUint256FromBuffer(buf)
}

func NewUint256FromBuffer(buffer *bytes.Buffer) (*UInt256, error) {
	val, err := BigFromBuffer(buffer)
	v := UInt256{val: val}
	return &v, err
}
