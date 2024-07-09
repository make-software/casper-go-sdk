package key

import (
	"encoding/json"
	"strings"
)

type TransferHash struct {
	Hash
	originPrefix string
}

func NewTransferHash(source string) (TransferHash, error) {
	var originPrefix string
	if strings.HasPrefix(source, PrefixNameTransfer) {
		originPrefix = PrefixNameTransfer
	}
	hexBytes, err := NewHash(strings.TrimPrefix(source, originPrefix))
	if err != nil {
		return TransferHash{}, err
	}
	return TransferHash{Hash: hexBytes, originPrefix: originPrefix}, err
}

func (h *TransferHash) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val, err := NewTransferHash(s)
	if err != nil {
		return err
	}
	*h = val
	return nil

}

func (h TransferHash) ToPrefixedString() string {
	return PrefixNameTransfer + h.ToHex()
}

func (h TransferHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.originPrefix + h.ToHex())
}
