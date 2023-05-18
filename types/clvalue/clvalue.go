package clvalue

import (
	"errors"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
	"github.com/make-software/casper-go-sdk/types/key"
	"github.com/make-software/casper-go-sdk/types/keypair"
)

type IValue interface {
	Bytes() []byte
	String() string
}

type CLValue struct {
	Type      cltype.CLType
	Bool      *Bool
	I32       *Int32
	I64       *Int64
	UI8       *UInt8
	UI32      *UInt32
	UI64      *UInt64
	UI128     *UInt128
	UI256     *UInt256
	UI512     *UInt512
	Unit      *Unit
	Uref      *key.URef
	Key       *key.Key
	Option    *Option
	List      *List
	ByteArray *ByteArray
	Result    *Result
	StringVal *String
	Map       *Map
	Tuple1    *Tuple1
	Tuple2    *Tuple2
	Tuple3    *Tuple3
	Any       *Any
	PublicKey *keypair.PublicKey
}

func (c CLValue) GetType() cltype.CLType {
	if val, ok := c.Type.(*cltype.Dynamic); ok {
		return val.Inner
	}
	return c.Type
}

func (c CLValue) ToBytesWithType() ([]byte, error) {
	return ToBytesWithType(c)
}

func (c CLValue) String() string {
	return c.GetValueByType().String()
}

func (c CLValue) Bytes() []byte {
	return c.GetValueByType().Bytes()
}

func (c CLValue) GetValueByType() IValue {
	switch c.Type.GetTypeID() {
	case cltype.TypeIDBool:
		return c.Bool
	case cltype.TypeIDI32:
		return c.I32
	case cltype.TypeIDI64:
		return c.I64
	case cltype.TypeIDU8:
		return c.UI8
	case cltype.TypeIDU32:
		return c.UI32
	case cltype.TypeIDU64:
		return c.UI64
	case cltype.TypeIDU128:
		return c.UI128
	case cltype.TypeIDU256:
		return c.UI256
	case cltype.TypeIDU512:
		return c.UI512
	case cltype.TypeIDUnit:
		return c.Unit
	case cltype.TypeIDString:
		return c.StringVal
	case cltype.TypeIDKey:
		return c.Key
	case cltype.TypeIDURef:
		return c.Uref
	case cltype.TypeIDOption:
		return c.Option
	case cltype.TypeIDList:
		return c.List
	case cltype.TypeIDByteArray:
		return c.ByteArray
	case cltype.TypeIDResult:
		return c.Result
	case cltype.TypeIDMap:
		return c.Map
	case cltype.TypeIDTuple1:
		return c.Tuple1
	case cltype.TypeIDTuple2:
		return c.Tuple2
	case cltype.TypeIDTuple3:
		return c.Tuple3
	case cltype.TypeIDAny:
		return c.Any
	case cltype.TypeIDPublicKey:
		return c.PublicKey
	}

	panic("type in GetValueByType method of CLValue is not implemented")
}

func (c *CLValue) GetKey() (*key.Key, error) {
	if c.Key == nil {
		return nil, errors.New("Key property is empty in CLValue, type is " + c.Type.String())
	}
	return c.Key, nil
}
