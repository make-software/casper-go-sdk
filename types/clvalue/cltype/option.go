package cltype

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Option struct {
	Inner CLType
}

func (t *Option) Bytes() []byte {
	return append([]byte{t.GetTypeID()}, (t.Inner).Bytes()...)
}

func (t *Option) String() string {
	return fmt.Sprintf("(%s: %s)", t.Name(), t.Inner.Name())
}

func (t *Option) GetTypeID() TypeID {
	return TypeIDOption
}

func (t *Option) Name() TypeName {
	return TypeNameOption
}

func (t *Option) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]CLType{t.Name(): t.Inner})
}

func NewOptionType(inner CLType) *Option {
	return &Option{Inner: inner}
}

func NewOptionFromJson(source interface{}) (*Option, error) {
	inner, err := fromInterface(source)
	if err != nil {
		return nil, err
	}
	return NewOptionType(inner), nil
}

func NewOptionFromBuffer(buf *bytes.Buffer) (*Option, error) {
	inner, err := FromBuffer(buf)
	if err != nil {
		return nil, err
	}
	return NewOptionType(inner), nil
}
