package types

import (
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/key"
	"github.com/make-software/casper-go-sdk/types/keypair"
)

// baseBid is an internal structure is using to avoid duplication of code
type baseBid struct {
	// The purse was used for bonding.
	BondingPurse key.URef `json:"bonding_purse"`
	// The delegation rate.
	DelegationRate float32 `json:"delegation_rate"`
	// `true` if validator has been "evicted"
	Inactive bool `json:"inactive"`
	// The amount of tokens staked by a validator (not including delegators).
	StakedAmount clvalue.UInt512 `json:"staked_amount"`
}

// ValidatorBid is an entry in the validator map.
type ValidatorBid struct {
	baseBid
	// Minimum allowed delegation amount in motes
	MinimumDelegationAmount uint64 `json:"minimum_delegation_amount"`
	// Maximum allowed delegation amount in motes
	MaximumDelegationAmount uint64 `json:"maximum_delegation_amount"`
	// Vesting schedule for a genesis validator. `None` if non-genesis validator.
	VestingSchedule *VestingSchedule `json:"vesting_schedule"`
}

// Bid is an entry stored in the Global state and representing a bid.
type Bid struct {
	baseBid
	// Validator's public key.
	ValidatorPublicKey keypair.PublicKey `json:"validator_public_key"`
	// The delegators.
	Delegators map[string]Delegator `json:"delegators"`
	// Vesting schedule for a genesis validator. `None` if non-genesis validator.
	VestingSchedule *VestingSchedule `json:"vesting_schedule"`
}

// AuctionBid is an entry in a founding validator map in the Auction state representing a bid.
type AuctionBid struct {
	baseBid
	// The delegators.
	Delegators []AuctionDelegators `json:"delegators"`
}

// baseDelegator is an internal structure is using to avoid a duplication of code
// two possibilities called 'Delegator' and 'JsonDelegator' in the rpc schema
// An array of delegators is returned in the GetAuctionState response
// A dictionary of delegators is returned in the QueryGlobalState response for a Bid key
type baseDelegator struct {
	// The purse that was used for delegating.
	BondingPurse key.URef `json:"bonding_purse"`
	// Amount of Casper token (in motes) delegated
	StakedAmount clvalue.UInt512 `json:"staked_amount"`
}

// Delegator is associated with the given validator.
type Delegator struct {
	baseDelegator
	// Public Key of the delegator
	Delegatee keypair.PublicKey `json:"delegator_public_key"`
	// Public key of the validator
	PublicKey keypair.PublicKey `json:"validator_public_key"`
	// Vesting schedule for a genesis validator. `None` if non-genesis validator.
	VestingSchedule *VestingSchedule `json:"vesting_schedule"`
}

// Credit is a bridge record pointing to a new `ValidatorBid` after the public key was changed.
type Credit struct {
	// The era id the credit was created.
	EraID uint32 `json:"era_id"`
	// Validator's public key.
	ValidatorPublicKey keypair.PublicKey `json:"validator_public_key"`
	// The credit amount.
	Amount clvalue.UInt512 `json:"amount"`
}

// Bridge is a bridge record pointing to a new `ValidatorBid` after the public key was changed.
type Bridge struct {
	EraID uint32 `json:"era_id"`
	// Previous validator public key associated with the bid."
	OldValidatorPublicKey keypair.PublicKey `json:"old_validator_public_key"`
	// New validator public key associated with the bid.
	NewValidatorPublicKey keypair.PublicKey `json:"new_validator_public_key"`
}

// AuctionDelegators is associated with the given validator.
type AuctionDelegators struct {
	baseDelegator
	// Public Key of the delegator
	Delegatee keypair.PublicKey `json:"delegatee"`
	// Public key of the validator
	PublicKey keypair.PublicKey `json:"public_key"`
}

// VestingSchedule for a genesis validator.
type VestingSchedule struct {
	InitialReleaseTimestampMillis uint64            `json:"initial_release_timestamp_millis"`
	LockedAmounts                 []clvalue.UInt512 `json:"locked_amounts"`
}
