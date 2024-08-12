package types

import (
	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/key"
	"github.com/make-software/casper-go-sdk/v2/types/keypair"
)

// UnbondingPurse stores information of an unbonding or delegation withdrawal
type UnbondingPurse struct {
	// Unbonding Amount
	Amount clvalue.UInt512 `json:"amount"`
	// Bonding purse
	BondingPurse key.URef `json:"bonding_purse"`
	// Era ID in which this unbonding request was created.
	EraOfCreation uint64 `json:"era_of_creation"`
	// Unbonder public key.
	UnbonderPublicKey keypair.PublicKey `json:"unbonder_public_key"`
	// The original validator's public key.
	ValidatorPublicKey keypair.PublicKey `json:"validator_public_key"`
	// The re-delegated validator's public key.
	NewValidator *keypair.PublicKey `json:"new_validator"`
}
