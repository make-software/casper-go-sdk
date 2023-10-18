package cltype

import (
	"encoding/json"
)

const Int32ByteSize = 4
const Int64ByteSize = 8

type TypeID = byte

const (
	TypeIDBool TypeID = iota
	TypeIDI32
	TypeIDI64
	TypeIDU8
	TypeIDU32
	TypeIDU64
	TypeIDU128
	TypeIDU256
	TypeIDU512
	TypeIDUnit
	TypeIDString
	TypeIDKey
	TypeIDURef
	TypeIDOption
	TypeIDList
	TypeIDByteArray
	TypeIDResult
	TypeIDMap
	TypeIDTuple1
	TypeIDTuple2
	TypeIDTuple3
	TypeIDAny
	TypeIDPublicKey
)

type TypeName = string

const (
	TypeNameBool      TypeName = "Bool"
	TypeNameI32       TypeName = "I32"
	TypeNameI64       TypeName = "I64"
	TypeNameU8        TypeName = "U8"
	TypeNameU32       TypeName = "U32"
	TypeNameU64       TypeName = "U64"
	TypeNameU128      TypeName = "U128"
	TypeNameU256      TypeName = "U256"
	TypeNameU512      TypeName = "U512"
	TypeNameUnit      TypeName = "Unit"
	TypeNameString    TypeName = "String"
	TypeNameKey       TypeName = "Key"
	TypeNameURef      TypeName = "URef"
	TypeNameOption    TypeName = "Option"
	TypeNameList      TypeName = "List"
	TypeNameByteArray TypeName = "ByteArray"
	TypeNameResult    TypeName = "Result"
	TypeNameMap       TypeName = "Map"
	TypeNameTuple1    TypeName = "Tuple1"
	TypeNameTuple2    TypeName = "Tuple2"
	TypeNameTuple3    TypeName = "Tuple3"
	TypeNameAny       TypeName = "Any"
	TypeNamePublicKey TypeName = "PublicKey"
)

var (
	Bool      = SimpleType{typeID: TypeIDBool, name: TypeNameBool}
	Int32     = SimpleType{typeID: TypeIDI32, name: TypeNameI32}
	Int64     = SimpleType{typeID: TypeIDI64, name: TypeNameI64}
	UInt8     = SimpleType{typeID: TypeIDU8, name: TypeNameU8}
	UInt32    = SimpleType{typeID: TypeIDU32, name: TypeNameU32}
	UInt64    = SimpleType{typeID: TypeIDU64, name: TypeNameU64}
	UInt128   = SimpleType{typeID: TypeIDU128, name: TypeNameU128}
	UInt256   = SimpleType{typeID: TypeIDU256, name: TypeNameU256}
	UInt512   = SimpleType{typeID: TypeIDU512, name: TypeNameU512}
	Unit      = SimpleType{typeID: TypeIDUnit, name: TypeNameUnit}
	String    = SimpleType{typeID: TypeIDString, name: TypeNameString}
	Key       = SimpleType{typeID: TypeIDKey, name: TypeNameKey}
	Uref      = SimpleType{typeID: TypeIDURef, name: TypeNameURef}
	Any       = SimpleType{typeID: TypeIDAny, name: TypeNameAny}
	PublicKey = SimpleType{typeID: TypeIDPublicKey, name: TypeNamePublicKey}
)

type CLType interface {
	Bytes() []byte
	String() string
	GetTypeID() TypeID
	Name() TypeName
}

type SimpleType struct {
	typeID TypeID
	name   TypeName
}

func (c SimpleType) Bytes() []byte {
	return []byte{c.typeID}
}

func (c SimpleType) String() string {
	return c.Name()
}

func (c SimpleType) GetTypeID() TypeID {
	return c.typeID
}

func (c SimpleType) Name() TypeName {
	return c.name
}

func (c SimpleType) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Name())
}
