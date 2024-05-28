package types

import (
	"encoding/json"
	"errors"

	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/key"
	"github.com/make-software/casper-go-sdk/types/keypair"
)

// Transfer a versioned wrapper for a transfer.
type Transfer struct {
	TransferV2

	// source TransferV1, nil if constructed from TransferV2
	OriginV1 *TransferV1
}

func (h *Transfer) UnmarshalJSON(bytes []byte) error {
	versioned := struct {
		Version1 *TransferV1 `json:"Version1,omitempty"`
		Version2 *TransferV2 `json:"Version2,omitempty"`
	}{}

	if err := json.Unmarshal(bytes, &versioned); err != nil {
		return err
	}

	if versioned.Version2 != nil {
		*h = Transfer{
			TransferV2: *versioned.Version2,
		}
		return nil
	}

	if versioned.Version1 != nil {
		*h = NewTransferFromV1(*versioned.Version1)
		return nil
	}

	//v1 compatible
	var v1Compatible = TransferV1{}
	if err := json.Unmarshal(bytes, &v1Compatible); err == nil {
		*h = NewTransferFromV1(v1Compatible)
		return nil
	}

	return errors.New("incorrect RPC response structure")
}

func NewTransferFromV1(transfer TransferV1) Transfer {
	return Transfer{
		TransferV2: TransferV2{
			Amount: transfer.Amount,
			TransactionHash: TransactionHash{
				TransactionV1Hash: &transfer.DeployHash,
			},
			From: InitiatorAddr{
				AccountHash: &transfer.From,
			},
			Gas:    transfer.Gas,
			ID:     transfer.ID,
			Source: transfer.Source,
			Target: transfer.Target,
			To:     transfer.To,
		},
		OriginV1: &transfer,
	}
}

// TransferV1 represents a transfer from one purse to another
type TransferV1 struct {
	// Transfer amount
	Amount clvalue.UInt512 `json:"amount"`
	// Deploy that created the transfer
	DeployHash key.Hash `json:"deploy_hash"`
	// Account hash from which transfer was executed
	From key.AccountHash `json:"from"`
	Gas  uint            `json:"gas,string"`
	// User-defined id
	ID uint64 `json:"id,omitempty"`
	// Source purse
	Source key.URef `json:"source"`
	// Target purse
	Target key.URef `json:"target"`
	// Account to which funds are transferred
	To *key.AccountHash `json:"to"`
}

// TransferV2 represents a version 2 transfer from one purse to another.
type TransferV2 struct {
	// Transfer amount
	Amount clvalue.UInt512 `json:"amount"`
	// Deploy that created the transfer
	TransactionHash TransactionHash `json:"transaction_hash"`
	// Account hash from which transfer was executed
	From InitiatorAddr `json:"from"`
	Gas  uint          `json:"gas,string"`
	// User-defined id
	ID uint64 `json:"id,omitempty"`
	// Source purse
	Source key.URef `json:"source"`
	// Target purse
	Target key.URef `json:"target"`
	// Account to which funds are transferred
	To *key.AccountHash `json:"to"`
}

// InitiatorAddr the address of the initiator of a TransactionV1.
type InitiatorAddr struct {
	// The public key of the initiator
	PublicKey *keypair.PublicKey `json:"PublicKey,omitempty"`
	// The account hash derived from the public key of the initiator
	AccountHash *key.AccountHash `json:"AccountHash,omitempty"`
}
