package cltype

import (
	"encoding/json"
)

type Dynamic struct {
	TypeID TypeID
	Inner  CLType
}

func (u Dynamic) Bytes() []byte {
	return u.Inner.Bytes()
}

func (u Dynamic) String() string {
	return u.Inner.String()
}

func (u Dynamic) GetTypeID() TypeID {
	return u.TypeID
}

func (u Dynamic) Name() TypeName {
	return u.Inner.Name()
}

func (u Dynamic) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.Inner)
}
