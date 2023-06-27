package cltype

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

var (
	BufferConstructorNotDetectedError = errors.New("buffer constructor not detected")
	ComplexTypeFormatInvalidError     = errors.New("complex type format is invalid")
	ErrComplexTypeFormatNotDetected   = errors.New("complex type format is not detected")
	ErrJsonConstructorNotFound        = errors.New("json type constructor is not found")
)

var simpleTypeByName = map[TypeName]CLType{
	TypeNameBool:      Bool,
	TypeNameI32:       Int32,
	TypeNameI64:       Int64,
	TypeNameU8:        UInt8,
	TypeNameU32:       UInt32,
	TypeNameU64:       UInt64,
	TypeNameU128:      UInt128,
	TypeNameU256:      UInt256,
	TypeNameU512:      UInt512,
	TypeNameUnit:      Unit,
	TypeNameString:    String,
	TypeNameKey:       Key,
	TypeNameURef:      Uref,
	TypeNameAny:       Any,
	TypeNamePublicKey: PublicKey,
}

var simpleTypeByID = map[TypeID]CLType{
	TypeIDBool:      Bool,
	TypeIDI32:       Int32,
	TypeIDI64:       Int64,
	TypeIDU8:        UInt8,
	TypeIDU32:       UInt32,
	TypeIDU64:       UInt64,
	TypeIDU128:      UInt128,
	TypeIDU256:      UInt256,
	TypeIDU512:      UInt512,
	TypeIDUnit:      Unit,
	TypeIDString:    String,
	TypeIDKey:       Key,
	TypeIDURef:      Uref,
	TypeIDAny:       Any,
	TypeIDPublicKey: PublicKey,
}

func GetSimpleTypeByName(typeName TypeName) (CLType, error) {
	result, ok := simpleTypeByName[typeName]
	if !ok {
		return nil, fmt.Errorf("type name is not registered, source: %s", typeName)
	}
	return result, nil
}

func FromRawJson(source json.RawMessage) (CLType, error) {
	var rawData interface{}
	err := json.Unmarshal(source, &rawData)
	if err != nil {
		return GetSimpleTypeByName(string(source))
	}
	return fromInterface(rawData)
}

func FromBytes(source []byte) (CLType, error) {
	buf := bytes.NewBuffer(source)
	return FromBuffer(buf)
}

func FromBuffer(buf *bytes.Buffer) (CLType, error) {
	dest, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}

	switch dest {
	case
		TypeIDBool,
		TypeIDI32,
		TypeIDI64,
		TypeIDU8,
		TypeIDU32,
		TypeIDU64,
		TypeIDU128,
		TypeIDU256,
		TypeIDU512,
		TypeIDUnit,
		TypeIDString,
		TypeIDKey,
		TypeIDURef,
		TypeIDAny,
		TypeIDPublicKey:
		return simpleTypeByID[dest], nil
	case TypeIDOption:
		return NewOptionFromBuffer(buf)
	case TypeIDList:
		return NewListFromBuffer(buf)
	case TypeIDByteArray:
		return NewByteArrayFromBuffer(buf)
	case TypeIDResult:
		return NewResultFromBuffer(buf)
	case TypeIDMap:
		return NewMapFromBuffer(buf)
	case TypeIDTuple1:
		return NewTuple1FromBuffer(buf)
	case TypeIDTuple2:
		return NewTuple2FromBuffer(buf)
	case TypeIDTuple3:
		return NewTuple3FromBuffer(buf)
	}

	return nil, BufferConstructorNotDetectedError
}

func fromInterface(rawData interface{}) (CLType, error) {
	if data, ok := rawData.(string); ok {
		return GetSimpleTypeByName(data)
	}

	return fromComplexStruct(rawData)
}

func fromComplexStruct(rawData interface{}) (CLType, error) {
	if data, ok := rawData.(map[string]interface{}); ok {
		if len(data) > 1 {
			return nil, ComplexTypeFormatInvalidError
		}
		for key, val := range data {
			switch key {
			case TypeNameOption:
				return NewOptionFromJson(val)
			case TypeNameList:
				return NewListFromJson(val)
			case TypeNameByteArray:
				return NewByteArrayFromJson(val)
			case TypeNameResult:
				return NewResultFromJson(val)
			case TypeNameMap:
				return NewMapFromJson(val)
			case TypeNameTuple1:
				return NewTuple1FromJson(val)
			case TypeNameTuple2:
				return NewTuple2FromJson(val)
			case TypeNameTuple3:
				return NewTuple3FromJson(val)
			}
			return nil, ErrJsonConstructorNotFound
		}
	}
	return nil, ErrComplexTypeFormatNotDetected
}
