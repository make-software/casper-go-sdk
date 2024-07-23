package types

import (
	"github.com/make-software/casper-go-sdk/v2/types/key"
	"github.com/make-software/casper-go-sdk/v2/types/keypair"
)

// InitiatorAddr the address of the initiator of a TransactionV1.
type InitiatorAddr struct {
	// The public key of the initiator
	PublicKey *keypair.PublicKey `json:"PublicKey,omitempty"`
	// The account hash derived from the public key of the initiator
	AccountHash *key.AccountHash `json:"AccountHash,omitempty"`
}

func (d InitiatorAddr) Bytes() []byte {
	result := make([]byte, 0, 32)

	if d.AccountHash != nil {
		result = append(result, 1)
		result = append(result, d.AccountHash.Bytes()...)
	} else if d.PublicKey != nil {
		result = append(result, 0)
		result = append(result, d.PublicKey.Bytes()...)
	}

	return result
}
