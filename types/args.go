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

func (args Args) Bytes() ([]byte, error) {
	var result []byte
	result = append(result, clvalue.SizeToBytes(len(args))...)
	for _, arg := range args {
		val, err := arg.Value()
		if err != nil {
			return nil, err
		}
		argName, err := arg.Name()
		if err != nil {
			return nil, err
		}
		result = append(result, clvalue.NewCLString(argName).Bytes()...)
		valueBytes, err := clvalue.ToBytesWithType(val)
		if err != nil {
			return nil, err
		}
		result = append(result, valueBytes...)
	}

	return result, nil
}

func (args Args) Find(name string) (*Argument, error) {
	for _, one := range args {
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

func (args *Args) AddArgument(name string, value clvalue.CLValue) *Args {
	pair := PairArgument{}
	pair[0] = &Argument{name: &name}
	pair[1] = &Argument{value: &value}
	*args = append(*args, pair)
	return args
}

type PairArgument [2]*Argument

func (r PairArgument) Name() (string, error) {
	return r[0].Name()
}

func (r PairArgument) Value() (clvalue.CLValue, error) {
	return r.Argument().Value()
}

func (r PairArgument) Argument() *Argument {
	return r[1]
}

type Argument struct {
	rawData json.RawMessage
	name    *string
	value   *clvalue.CLValue
}

func (a *Argument) Value() (clvalue.CLValue, error) {
	if a.value != nil {
		return *a.value, nil
	}
	return ArgsFromRawJson(a.rawData)
}

func (a *Argument) Raw() (RawArg, error) {
	var rawArg RawArg
	if a.rawData == nil {
		return RawArg{}, nil
	}
	err := json.Unmarshal(a.rawData, &rawArg)
	if err != nil {
		return RawArg{}, err
	}
	return rawArg, nil
}

func (a *Argument) Parsed() (json.RawMessage, error) {
	rawArg, err := a.Raw()
	if err != nil {
		return nil, err
	}
	return rawArg.Parsed, nil
}

func (a *Argument) Bytes() (HexBytes, error) {
	if a.value != nil {
		return clvalue.ToBytesWithType(*a.value)
	}
	rawArg, err := a.Raw()
	if err != nil {
		return nil, err
	}
	return rawArg.Bytes, nil
}

func (a *Argument) UnmarshalJSON(bytes []byte) error {
	a.rawData = bytes
	return nil
}

func (a *Argument) MarshalJSON() ([]byte, error) {
	if a.rawData != nil {
		return a.rawData, nil
	}
	if a.name != nil {
		return json.Marshal(a.name)
	}
	typeName, err := json.Marshal(a.value.Type)
	if err != nil {
		return nil, err
	}
	argData := RawArg{
		CLType: typeName,
		Bytes:  a.value.Bytes(),
	}
	return json.Marshal(argData)
}

func (a *Argument) Name() (string, error) {
	if a.name != nil {
		return *a.name, nil
	}
	err := json.Unmarshal(a.rawData, &a.name)
	if err != nil {
		return "", err
	}
	return *a.name, nil
}

// RawArg is a type used in deploy input arguments. And it can also be returned as a
// result of a query to the network or a contract call.
type RawArg struct {
	// Type of the value. Can be simple or constructed
	CLType json.RawMessage `json:"cl_type"`
	// Bytes array representation of underlying data
	Bytes HexBytes `json:"bytes"`
	// The optional parsed value of the bytes used when testing
	Parsed json.RawMessage `json:"parsed,omitempty"`
}

func ArgsFromRawJson(raw json.RawMessage) (clvalue.CLValue, error) {
	var rawData RawArg
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
