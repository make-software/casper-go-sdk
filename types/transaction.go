package types

import "github.com/make-software/casper-go-sdk/types/key"

// TransactionHash A versioned wrapper for a transaction hash or deploy hash
type TransactionHash struct {
	Deploy      *key.Hash `json:"Deploy,omitempty"`
	Transaction *key.Hash `json:"Version1,omitempty"`
}
