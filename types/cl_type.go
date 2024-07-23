package types

import (
	"encoding/json"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
)

type CLTypeRaw struct {
	json.RawMessage
}

func (t CLTypeRaw) ParseCLType() (cltype.CLType, error) {
	return cltype.FromRawJson(t.RawMessage)
}
