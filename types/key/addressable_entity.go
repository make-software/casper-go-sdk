package key

import (
	"encoding/json"
	"strings"
)

type AddressableEntityHash struct {
	Hash
	originPrefix PrefixName
}

func NewAddressableEntityHash(source string) (AddressableEntityHash, error) {
	var originPrefix string
	if strings.HasPrefix(source, PrefixAddressableEntity) {
		originPrefix = PrefixAddressableEntity
	}

	hexBytes, err := NewHash(strings.TrimPrefix(source, originPrefix))
	if err != nil {
		return AddressableEntityHash{}, err
	}

	return AddressableEntityHash{Hash: hexBytes, originPrefix: originPrefix}, nil
}

func (h *AddressableEntityHash) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val, err := NewAddressableEntityHash(s)
	if err != nil {
		return err
	}
	*h = val
	return nil
}

func (h AddressableEntityHash) ToPrefixedString() string {
	return PrefixNameAccount + h.ToHex()
}

func (h AddressableEntityHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.originPrefix + h.ToHex())
}
