package types

import (
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
	StakedAmount uint64 `json:"staked_amount,string"`
}

// Bid is an entry stored in the Global state and representing a bid.
type Bid struct {
	baseBid
	// Validator's public key.
	PublicKey keypair.PublicKey `json:"validator_public_key"`
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
	StakedAmount uint64 `json:"staked_amount,string"`
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
	InitialReleaseTimestampMillis uint64   `json:"initial_release_timestamp_millis"`
	LockedAmounts                 []uint64 `json:"locked_amounts"`
}
