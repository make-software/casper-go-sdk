package types

import (
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/key"
)

// Transfer represents a transfer from one purse to another
type Transfer struct {
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
