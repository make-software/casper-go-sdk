package clvalue

import (
	"bytes"
	"errors"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
	"github.com/make-software/casper-go-sdk/types/key"
	"github.com/make-software/casper-go-sdk/types/keypair"
)

var (
	ErrUnsupportedCLType = errors.New("buffer constructor is not found")
)

func FromBytes(source []byte) (CLValue, error) {
	buffer := bytes.NewBuffer(source)
	valueLength, err := TrimByteSize(buffer)
	if err != nil {
		return CLValue{}, err
	}
	clType, err := cltype.FromBytes(buffer.Bytes()[valueLength:])
	if err != nil {
		return CLValue{}, err
	}

	return FromBytesByType(buffer.Bytes()[:valueLength], clType)
}

func FromBuffer(buffer *bytes.Buffer) (CLValue, error) {
	valueLength, err := TrimByteSize(buffer)
	if err != nil {
		return CLValue{}, err
	}
	data := buffer.Next(int(valueLength))
	clType, err := cltype.FromBuffer(buffer)
	if err != nil {
		return CLValue{}, err
	}

	return FromBytesByType(data, clType)
}

func ToBytesWithType(value CLValue) ([]byte, error) {
	valueBytes := value.Bytes()
	valueSize := SizeToBytes(len(valueBytes))
	return append(valueSize, append(valueBytes, value.Type.Bytes()...)...), nil
}

func FromBytesByType(source []byte, clType cltype.CLType) (CLValue, error) {
	buf := bytes.NewBuffer(source)
	return FromBufferByType(buf, clType)
}

func FromBufferByType(buf *bytes.Buffer, sourceType cltype.CLType) (result CLValue, err error) {
	result.Type = sourceType
	switch clType := result.Type.(type) {
	case cltype.SimpleType:
		switch clType.GetTypeID() {
		case cltype.TypeIDBool:
			result.Bool, err = NewBoolFromBuffer(buf)
			return result, err
		case cltype.TypeIDI32:
			result.I32, err = NewInt32FromBuffer(buf)
			return result, err
		case cltype.TypeIDI64:
			result.I64, err = NewInt64FromBuffer(buf)
			return result, err
		case cltype.TypeIDU8:
			result.UI8, err = NewUInt8FromBuffer(buf)
			return result, err
		case cltype.TypeIDU32:
			result.UI32, err = NewUint32FromBuffer(buf)
			return result, err
		case cltype.TypeIDU64:
			result.UI64, err = NewUint64FromBuffer(buf)
			return result, err
		case cltype.TypeIDU128:
			result.UI128, err = NewUint128FromBuffer(buf)
			return result, err
		case cltype.TypeIDU256:
			result.UI256, err = NewUint256FromBuffer(buf)
			return result, err
		case cltype.TypeIDU512:
			result.UI512, err = NewUint512FromBuffer(buf)
			return result, err
		case cltype.TypeIDString:
			result.StringVal, err = NewStringFromBuffer(buf)
			return result, err
		case cltype.TypeIDUnit:
			result.Unit, err = NewUnitFromBuffer(buf)
			return result, err
		case cltype.TypeIDKey:
			keyFromBuffer, err := key.NewKeyFromBuffer(buf)
			if err != nil {
				return result, err
			}
			result.Key = &keyFromBuffer
			return result, nil
		case cltype.TypeIDURef:
			uRef, err := key.NewURefFromBuffer(buf)
			if err != nil {
				return result, err
			}
			result.Uref = &uRef
			return result, nil
		case cltype.TypeIDAny:
			buffer := NewAnyFromBuffer(buf)
			result.Any = &buffer
			return result, nil
		case cltype.TypeIDPublicKey:
			publicKey, err := keypair.NewPublicKeyFromBuffer(buf)
			if err != nil {
				return result, err
			}
			result.PublicKey = &publicKey
			return result, err
		}
	case *cltype.Option:
		result.Option, err = NewOptionFromBuffer(buf, clType)
		return result, err
	case *cltype.List:
		result.List, err = NewListFromBuffer(buf, clType)
		return result, err
	case *cltype.ByteArray:
		buffer := NewByteArrayFromBuffer(buf, clType)
		result.ByteArray = &buffer
		return result, nil
	case *cltype.Result:
		result.Result, err = NewResultFromBuffer(buf, clType)
		return result, err
	case *cltype.Map:
		result.Map, err = NewMapFromBuffer(buf, clType)
		return result, err
	case *cltype.Tuple1:
		result.Tuple1, err = NewTuple1FromBuffer(buf, clType)
		return result, err
	case *cltype.Tuple2:
		result.Tuple2, err = NewTuple2FromBuffer(buf, clType)
		return result, err
	case *cltype.Tuple3:
		result.Tuple3, err = NewTuple3FromBuffer(buf, clType)
		return result, err
	case *cltype.Dynamic:
		typeData, err := cltype.FromBuffer(buf)
		if err != nil {
			return result, err
		}
		result.Type = &cltype.Dynamic{
			TypeID: typeData.GetTypeID(),
			Inner:  typeData,
		}
		return result, err
	}
	return result, ErrUnsupportedCLType
}
