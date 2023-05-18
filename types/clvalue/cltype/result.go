package cltype

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Result struct {
	Inner CLType
}

func (t *Result) Bytes() []byte {
	return append([]byte{t.GetTypeID()}, (t.Inner).Bytes()...)
}

func (t *Result) String() string {
	return fmt.Sprintf("(%s: %s)", t.Name(), t.Inner.Name())
}

func (t *Result) GetTypeID() TypeID {
	return TypeIDResult
}

func (t *Result) Name() TypeName {
	return TypeNameResult
}

func (t *Result) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]CLType{t.Name(): t.Inner})
}

func NewResultType(inner CLType) *Result {
	return &Result{Inner: inner}
}

func NewResultFromJson(source interface{}) (*Result, error) {
	inner, err := fromInterface(source)
	if err != nil {
		return nil, err
	}
	return NewResultType(inner), nil
}

func NewResultFromBuffer(buf *bytes.Buffer) (*Result, error) {
	inner, err := FromBuffer(buf)
	if err != nil {
		return nil, err
	}
	return NewResultType(inner), nil
}
