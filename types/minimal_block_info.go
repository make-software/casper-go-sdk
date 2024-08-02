package types

import (
	"time"

	"github.com/make-software/casper-go-sdk/v2/types/key"
	"github.com/make-software/casper-go-sdk/v2/types/keypair"
)

type MinimalBlockInfo struct {
	Creator       keypair.PublicKey `json:"creator"`
	EraID         uint64            `json:"era_id"`
	Hash          key.Hash          `json:"hash"`
	Height        uint32            `json:"height"`
	StateRootHash key.Hash          `json:"state_root_hash"`
	Timestamp     time.Time         `json:"timestamp"`
}
