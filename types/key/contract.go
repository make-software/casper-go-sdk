package key

import (
	"encoding/json"
	"strings"
)

type ContractHash struct {
	Hash
	originPrefix string
}

func (h *ContractHash) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	tmp, err := NewContract(s)
	if err != nil {
		return err
	}
	*h = tmp
	return nil
}

func (h ContractHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.originPrefix + h.ToHex())
}

func (h ContractHash) ToPrefixedWasmString() string {
	return PrefixNameContractWasm + h.ToHex()
}

func (h ContractHash) ToPrefixedString() string {
	return PrefixNameContract + h.ToHex()
}

func NewContract(source string) (ContractHash, error) {
	var originPrefix string
	if strings.HasPrefix(source, PrefixNameHash) {
		originPrefix = PrefixNameHash
	} else if strings.HasPrefix(source, PrefixNameContractWasm) {
		originPrefix = PrefixNameContractWasm
	} else if strings.HasPrefix(source, PrefixNameContract) {
		originPrefix = PrefixNameContract
	} else if strings.HasPrefix(source, PrefixNameEntityContract) {
		originPrefix = PrefixNameEntityContract
	}
	hexBytes, err := NewHash(strings.TrimPrefix(source, originPrefix))
	if err != nil {
		return ContractHash{}, err
	}

	return ContractHash{Hash: hexBytes, originPrefix: originPrefix}, nil
}
