package clvalue

import (
	"bytes"
	"encoding/json"
	"math/big"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

type UInt512 struct {
	val *big.Int
}

func (v *UInt512) Bytes() []byte {
	return BigToBytes(v.val)
}

func (v *UInt512) String() string {
	return v.val.String()
}

func (v *UInt512) Value() *big.Int {
	return v.val
}

func (v *UInt512) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Value().String())
}

func (v *UInt512) UnmarshalJSON(b []byte) error {
	var val string
	err := json.Unmarshal(b, &val)
	if err != nil {
		return err
	}
	v.val = &big.Int{}
	v.val.SetString(val, 10)
	return nil
}

func NewCLUInt512(val *big.Int) *CLValue {
	res := CLValue{}
	res.Type = cltype.UInt512
	v := UInt512{val: val}
	res.UI512 = &v
	return &res
}

func NewUint512FromBytes(source []byte) (*UInt512, error) {
	buf := bytes.NewBuffer(source)
	return NewUint512FromBuffer(buf)
}

func NewUint512FromBuffer(buffer *bytes.Buffer) (*UInt512, error) {
	fromBuffer, err := BigFromBuffer(buffer)
	v := UInt512{val: fromBuffer}
	return &v, err
}
