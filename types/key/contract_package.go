package key

import (
	"encoding/json"
	"strings"
)

type ContractPackageHash struct {
	Hash
	originPrefix string
}

func (h ContractPackageHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.originPrefix + h.ToHex())
}

func (h ContractPackageHash) ToPrefixedString() string {
	return PrefixNameContractPackage + h.ToHex()
}

func (h *ContractPackageHash) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	tmp, err := NewContractPackage(s)
	if err != nil {
		return err
	}
	*h = tmp
	return nil
}

func NewContractPackage(source string) (ContractPackageHash, error) {
	var originPrefix string
	if strings.HasPrefix(source, PrefixNameHash) {
		originPrefix = PrefixNameHash
	} else if strings.HasPrefix(source, PrefixNameContractPackageWasm) {
		originPrefix = PrefixNameContractPackageWasm
	} else if strings.HasPrefix(source, PrefixNameContractPackage) {
		originPrefix = PrefixNameContractPackage
	}
	hexBytes, err := NewHash(strings.TrimPrefix(source, originPrefix))
	if err != nil {
		return ContractPackageHash{}, err
	}

	return ContractPackageHash{Hash: hexBytes, originPrefix: originPrefix}, nil
}
