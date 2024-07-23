package types

import (
	"encoding/json"
	"errors"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/key"
)

// Transfer a versioned wrapper for a transfer.
type Transfer struct {
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

	// source originTransferV1, nil if constructed from TransferV2
	originTransferV1 *TransferV1
	// source originTransferV2, nil if constructed from TransferV1
	originTransferV2 *TransferV2
}

func (h *Transfer) GetTransferV1() *TransferV1 {
	return h.originTransferV1
}

func (h *Transfer) GetTransferV2() *TransferV2 {
	return h.originTransferV2
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
			Amount:           versioned.Version2.Amount,
			TransactionHash:  versioned.Version2.TransactionHash,
			From:             versioned.Version2.From,
			Gas:              versioned.Version2.Gas,
			ID:               versioned.Version2.ID,
			Source:           versioned.Version2.Source,
			Target:           versioned.Version2.Target,
			To:               versioned.Version2.To,
			originTransferV2: versioned.Version2,
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
		Amount: transfer.Amount,
		TransactionHash: TransactionHash{
			Deploy: &transfer.DeployHash,
		},
		From: InitiatorAddr{
			AccountHash: &transfer.From,
		},
		Gas:              transfer.Gas,
		ID:               transfer.ID,
		Source:           transfer.Source,
		Target:           transfer.Target,
		To:               transfer.To,
		originTransferV1: &transfer,
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
