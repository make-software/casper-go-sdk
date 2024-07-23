package types

import (
	"github.com/make-software/casper-go-sdk/v2/types/keypair"
)

// AuctionState is a data structure summarizing auction contract data.
type AuctionState struct {
	// All bids contained within a vector.
	Bids          []PublicKeyAndBid `json:"bids"`
	BlockHeight   uint64            `json:"block_height"`
	EraValidators []EraValidators   `json:"era_validators"`
	StateRootHash string            `json:"state_root_hash"`
}

// PublicKeyAndBid is an entry in a founding validator map representing a bid.
type PublicKeyAndBid struct {
	// Validator public key
	PublicKey keypair.PublicKey `json:"public_key"`
	Bid       Bid               `json:"bid"`
}

// EraValidators contains validators and weights for an Era.
type EraValidators struct {
	EraID uint64 `json:"era_id"`
	// List of the validator's weight in the Era
	ValidatorWeights []ValidatorWeightAuction `json:"validator_weights"`
}
