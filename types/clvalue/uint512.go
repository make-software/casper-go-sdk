package clvalue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
)

type UInt512 struct {
	val         *big.Int
	isStringFmt bool
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

func (v UInt512) MarshalJSON() ([]byte, error) {
	if v.isStringFmt {
		return json.Marshal(v.String())
	}
	return []byte(v.String()), nil
}

func (v *UInt512) UnmarshalJSON(b []byte) error {
	var num json.Number
	err := json.Unmarshal(b, &num)
	if err != nil {
		return err
	}

	// Convert json.Number to string
	s := num.String()

	// Check if the original data was a quoted string
	if b[0] == '"' {
		v.isStringFmt = true
	}

	v.val = new(big.Int)
	val, ok := v.val.SetString(s, 10)
	if !ok {
		return fmt.Errorf("invalid integer string: %s", s)
	}
	v.val = val

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
