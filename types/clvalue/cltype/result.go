package cltype

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

var ErrInvalidResultJsonFormat = errors.New("invalid json format for Result type")

type Result struct {
	InnerOk  CLType
	InnerErr CLType
}

func (t *Result) Bytes() []byte {
	return append([]byte{t.GetTypeID()}, append((t.InnerOk).Bytes(), t.InnerErr.Bytes()...)...)
}

func (t *Result) String() string {
	return fmt.Sprintf("(%s: Ok(%s), Err(%s)", t.Name(), t.InnerOk.Name(), t.InnerErr.Name())
}

func (t *Result) GetTypeID() TypeID {
	return TypeIDResult
}

func (t *Result) Name() TypeName {
	return TypeNameResult
}

func (t *Result) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]map[string]CLType{t.Name(): {"ok": t.InnerOk, "err": t.InnerErr}})
}

func NewResultType(innerOk CLType, innerErr CLType) *Result {
	return &Result{InnerOk: innerOk, InnerErr: innerErr}
}

func NewResultFromJson(source interface{}) (*Result, error) {
	data, ok := source.(map[string]interface{})
	if !ok {
		return nil, ErrInvalidResultJsonFormat
	}
	okData, found := data["ok"]
	if !found {
		return nil, ErrInvalidResultJsonFormat
	}
	innerOk, err := fromInterface(okData)
	if err != nil {
		return nil, err
	}
	errData, found := data["err"]
	if !found {
		return nil, ErrInvalidResultJsonFormat
	}
	innerErr, err := fromInterface(errData)
	if err != nil {
		return nil, err
	}
	return NewResultType(innerOk, innerErr), nil
}

func NewResultFromBuffer(buf *bytes.Buffer) (*Result, error) {
	innerOk, err := FromBuffer(buf)
	if err != nil {
		return nil, err
	}
	innerErr, err := FromBuffer(buf)
	if err != nil {
		return nil, err
	}
	return NewResultType(innerOk, innerErr), nil
}
