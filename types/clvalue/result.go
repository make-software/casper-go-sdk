package clvalue

import (
	"bytes"
	"fmt"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

type Result struct {
	Type      *cltype.Result
	IsSuccess bool
	Inner     CLValue
}

func (v *Result) Bytes() []byte {
	return append([]byte{01}, v.Inner.Bytes()...)
}

func (v *Result) String() string {
	if v.IsSuccess {
		return fmt.Sprintf("Ok(%s)", v.Inner.String())
	}
	return fmt.Sprintf("Err(%s)", v.Inner.String())
}

func (v *Result) Value() CLValue {
	return v.Inner
}

func NewCLResult(innerOk, innerErr cltype.CLType, value CLValue, isSuccess bool) CLValue {
	resultType := cltype.NewResultType(innerOk, innerErr)
	return CLValue{
		Type: resultType,
		Result: &Result{
			Type:      resultType,
			IsSuccess: isSuccess,
			Inner:     value,
		},
	}
}

func NewResultFromBytes(source []byte, clType *cltype.Result) (*Result, error) {
	return NewResultFromBuffer(bytes.NewBuffer(source), clType)
}

func NewResultFromBuffer(buf *bytes.Buffer, clType *cltype.Result) (*Result, error) {
	val := Result{}
	val.Type = clType
	isSuccess, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}
	val.IsSuccess = isSuccess == 1
	if val.IsSuccess {
		inner, err := FromBufferByType(buf, clType.InnerOk)
		if err != nil {
			return nil, err
		}
		val.Inner = inner
	} else {
		inner, err := FromBufferByType(buf, clType.InnerErr)
		if err != nil {
			return nil, err
		}
		val.Inner = inner
	}

	return &val, nil
}
