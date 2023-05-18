package clvalue

import (
	"bytes"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

type Option struct {
	Type  *cltype.Option
	Inner *CLValue
}

func (v *Option) Bytes() []byte {
	if v.IsEmpty() {
		return []byte{00}
	}
	return append([]byte{01}, v.Inner.Bytes()...)
}

func (v *Option) String() string {
	if v.IsEmpty() {
		return ""
	}
	return v.Inner.String()
}
func (v *Option) IsEmpty() bool {
	return v.Inner == nil
}

func (v *Option) Value() *CLValue {
	if v.IsEmpty() {
		return nil
	}
	return v.Inner
}

func NewCLOption(inner CLValue) CLValue {
	optionType := cltype.NewOptionType(inner.Type)
	return CLValue{
		Type: optionType,
		Option: &Option{
			Type:  optionType,
			Inner: &inner,
		},
	}
}

func NewOptionFromBytes(source []byte, clType *cltype.Option) (*Option, error) {
	return NewOptionFromBuffer(bytes.NewBuffer(source), clType)
}

func NewOptionFromBuffer(buf *bytes.Buffer, clType *cltype.Option) (*Option, error) {
	result := Option{}
	result.Type = clType
	hasData, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}
	if hasData == 0 {
		return &result, nil
	}
	inner, err := FromBufferByType(buf, clType.Inner)
	if err != nil {
		return nil, err
	}
	result.Inner = &inner
	return &result, nil
}
