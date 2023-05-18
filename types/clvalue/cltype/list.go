package cltype

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type List struct {
	ElementsType CLType
}

func (l *List) Bytes() []byte {
	return append([]byte{l.GetTypeID()}, l.ElementsType.Bytes()...)
}

func (l *List) String() string {
	return fmt.Sprintf("(%s of %s)", l.Name(), l.ElementsType.String())
}

func (l *List) GetTypeID() TypeID {
	return TypeIDList
}

func (l *List) Name() TypeName {
	return TypeNameList
}

func (l *List) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]CLType{l.Name(): l.ElementsType})
}

func NewList(clType CLType) *List {
	return &List{ElementsType: clType}
}

func NewListFromJson(source interface{}) (*List, error) {
	inner, err := fromInterface(source)
	if err != nil {
		return nil, err
	}
	return NewList(inner), nil
}

func NewListFromBuffer(buf *bytes.Buffer) (*List, error) {
	inner, err := FromBuffer(buf)
	if err != nil {
		return nil, err
	}
	return NewList(inner), nil
}
