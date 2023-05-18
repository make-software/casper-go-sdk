package key

import (
	"encoding/json"
	"strings"
)

type AccountHash struct {
	Hash
	originPrefix PrefixName
}

func NewAccountHash(source string) (AccountHash, error) {
	var originPrefix string
	if len(source) == 66 && strings.HasPrefix(source, "00") {
		originPrefix = "00"
	} else if strings.HasPrefix(source, PrefixNameAccount) {
		originPrefix = PrefixNameAccount
	}

	hexBytes, err := NewHash(strings.TrimPrefix(source, originPrefix))
	if err != nil {
		return AccountHash{}, err
	}

	return AccountHash{Hash: hexBytes, originPrefix: originPrefix}, err
}

func (h *AccountHash) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val, err := NewAccountHash(s)
	if err != nil {
		return err
	}
	*h = val
	return nil
}

func (h AccountHash) ToPrefixedString() string {
	return PrefixNameAccount + h.ToHex()
}

func (h AccountHash) MarshalJSON() ([]byte, error) {
	return []byte(`"` + h.originPrefix + h.ToHex() + `"`), nil
}
