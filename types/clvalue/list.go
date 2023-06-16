package clvalue

import (
	"bytes"
	"strings"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

type List struct {
	Type     *cltype.List
	Elements []CLValue
}

func (v *List) Bytes() []byte {
	var byteData []byte
	for _, one := range v.Elements {
		byteData = append(byteData, one.Bytes()...)
	}
	return append(SizeToBytes(v.Len()), byteData...)
}

func (v *List) String() string {
	var strData []string
	for _, one := range v.Elements {
		strData = append(strData, "\""+one.String()+"\"")
	}
	return "[" + strings.Join(strData, ",") + "]"
}

func (v *List) IsEmpty() bool {
	return v.Len() == 0
}

func (v *List) Len() int {
	return len(v.Elements)
}

func (v *List) Append(value CLValue) {
	v.Elements = append(v.Elements, value)
}

func NewCLList(elementType cltype.CLType) CLValue {
	listType := cltype.NewList(elementType)
	return CLValue{
		Type: listType,
		List: &List{
			Type:     listType,
			Elements: []CLValue{},
		},
	}
}

func NewListFromBytes(source []byte, clType *cltype.List) (*List, error) {
	return NewListFromBuffer(bytes.NewBuffer(source), clType)
}

func NewListFromBuffer(buf *bytes.Buffer, clType *cltype.List) (*List, error) {
	size, err := TrimByteSize(buf)
	if err != nil {
		return nil, err
	}
	listSize := int(size)
	elements := make([]CLValue, 0)
	for i := 0; i < listSize; i += 1 {
		one, err := FromBufferByType(buf, clType.ElementsType)
		if err != nil {
			return nil, err
		}
		elements = append(elements, one)
	}
	result := List{Type: clType, Elements: elements}

	return &result, nil
}
