package types

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

var ErrArgumentNotFound = errors.New("argument is not found")

type Args []PairArgument

func (a Args) Bytes() ([]byte, error) {
	var result []byte
	result = append(result, clvalue.SizeToBytes(len(a))...)
	for _, one := range a {
		val, err := one.Value()
		if err != nil {
			return nil, err
		}
		name, err := one.Name()
		if err != nil {
			return nil, err
		}
		result = append(result, clvalue.NewCLString(name).Bytes()...)
		valueBytes, err := clvalue.ToBytesWithType(val)
		if err != nil {
			panic(err)
		}
		result = append(result, valueBytes...)
	}

	return result, nil
}

func (a Args) Find(name string) (*Argument, error) {
	for _, one := range a {
		getName, err := one.Name()
		if err != nil {
			return nil, err
		}
		if getName == name {
			return one.Argument(), nil
		}
	}
	return nil, fmt.Errorf("%w, target: %s", ErrArgumentNotFound, name)
}

func (a *Args) AddArgument(name string, value clvalue.CLValue) *Args {
	pair := PairArgument{}
	pair[0] = &argumentData{name: &name}
	pair[1] = &argumentData{value: &value}
	*a = append(*a, pair)
	return a
}

type argumentData struct {
	rawData json.RawMessage
	name    *string
	value   *clvalue.CLValue
}

func (a *argumentData) UnmarshalJSON(bytes []byte) error {
	a.rawData = bytes
	return nil
}

func (a argumentData) MarshalJSON() ([]byte, error) {
	if a.rawData != nil {
		return a.rawData, nil
	}
	if a.name != nil {
		return json.Marshal(a.name)
	}
	typeName, err := json.Marshal(a.value.Type)
	if err != nil {
		panic(err)
	}
	argData := rawArg{
		CLType: typeName,
		Bytes:  a.value.Bytes(),
	}
	return json.Marshal(argData)
}

func (a *argumentData) Name() (string, error) {
	if a.name != nil {
		return *a.name, nil
	}
	err := json.Unmarshal(a.rawData, &a.name)
	if err != nil {
		return "", err
	}
	return *a.name, nil
}

func (a *argumentData) Value() (clvalue.CLValue, error) {
	if a.value != nil {
		return *a.value, nil
	}
	var err error
	value, err := a.Argument().Value()
	if err != nil {
		return clvalue.CLValue{}, err
	}

	a.value = &value

	return *a.value, nil
}

func (a *argumentData) Argument() *Argument {
	return &Argument{a.rawData}
}

type PairArgument [2]*argumentData

func (r PairArgument) Name() (string, error) {
	return r[0].Name()
}

func (r PairArgument) Value() (clvalue.CLValue, error) {
	return r[1].Value()
}

func (r PairArgument) Argument() *Argument {
	return r[1].Argument()
}

type Argument struct {
	json.RawMessage
}

func (a Argument) Value() (clvalue.CLValue, error) {
	return ArgsFromRawJson(a.RawMessage)
}

func (a Argument) Parsed() (json.RawMessage, error) {
	var rawData rawArg
	err := json.Unmarshal(a.RawMessage, &rawData)
	if err != nil {
		return nil, err
	}
	return rawData.Parsed, nil
}

func (a Argument) Bytes() (HexBytes, error) {
	var rawData rawArg
	err := json.Unmarshal(a.RawMessage, &rawData)
	if err != nil {
		return nil, err
	}
	return rawData.Bytes, nil
}

// rawArg is a type used in deploy input arguments. And it can also be returned as a
// result of a query to the network or a contract call.
type rawArg struct {
	// Type of the value. Can be simple or constructed
	CLType json.RawMessage `json:"cl_type"`
	// Bytes array representation of underlying data
	Bytes HexBytes `json:"bytes"`
	// The optional parsed value of the bytes used when testing
	Parsed json.RawMessage `json:"parsed,omitempty"`
}

func ArgsFromRawJson(raw json.RawMessage) (clvalue.CLValue, error) {
	var rawData rawArg
	err := json.Unmarshal(raw, &rawData)
	if err != nil {
		return clvalue.CLValue{}, err
	}
	valueType, err := cltype.FromRawJson(rawData.CLType)
	if err != nil {
		return clvalue.CLValue{}, err
	}
	return clvalue.FromBytesByType(rawData.Bytes, valueType)
}
