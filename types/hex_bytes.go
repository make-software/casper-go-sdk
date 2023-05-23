package types

import (
	"encoding/hex"
	"encoding/json"
)

type HexBytes []byte

func (h HexBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(h))
}

func (h *HexBytes) UnmarshalJSON(bytes []byte) error {
	var hexString string
	var err error
	if err = json.Unmarshal(bytes, &hexString); err != nil {
		return err
	}

	*h, err = hex.DecodeString(hexString)
	return err
}

func (h HexBytes) String() string {
	return hex.EncodeToString(h)
}
