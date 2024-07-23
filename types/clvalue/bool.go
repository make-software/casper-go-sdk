package clvalue

import (
	"bytes"
	"errors"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
)

var (
	ErrBoolValueIsInvalid = errors.New("invalid bool value")
)

type Bool bool

func (v Bool) Bytes() []byte {
	if v {
		return []byte{1}
	}
	return []byte{0}
}

func (v Bool) String() string {
	if v {
		return "true"
	}
	return "false"
}

func (v Bool) Value() bool {
	return bool(v)
}

func NewCLBool(val bool) CLValue {
	res := CLValue{}
	res.Type = cltype.Bool
	boolVal := Bool(val)
	res.Bool = &boolVal
	return res
}

func NewBoolFromBytes(source []byte) (*Bool, error) {
	buf := bytes.NewBuffer(source)
	return NewBoolFromBuffer(buf)
}

func NewBoolFromBuffer(buffer *bytes.Buffer) (*Bool, error) {
	byteVal, err := buffer.ReadByte()
	if err != nil {
		return nil, err
	}
	var val bool
	if byteVal == 1 {
		val = true
	} else if byteVal == 0 {
		val = false
	} else {
		return nil, ErrBoolValueIsInvalid
	}

	boolVal := Bool(val)
	return &boolVal, nil
}
