package types

import (
	"encoding/json"
	"errors"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/key"
	"github.com/make-software/casper-go-sdk/v2/types/keypair"
)

// ValidatorBid is an entry in the validator map.
type ValidatorBid struct {
	// The purse was used for bonding.
	BondingPurse key.URef `json:"bonding_purse"`
	// The delegation rate.
	DelegationRate float32 `json:"delegation_rate"`
	// `true` if validator has been "evicted"
	Inactive bool `json:"inactive"`
	// The amount of tokens staked by a validator (not including delegators).
	StakedAmount clvalue.UInt512 `json:"staked_amount"`
	// Minimum allowed delegation amount in motes
	MinimumDelegationAmount uint64 `json:"minimum_delegation_amount"`
	// Maximum allowed delegation amount in motes
	MaximumDelegationAmount uint64 `json:"maximum_delegation_amount"`
	// Vesting schedule for a genesis validator. `None` if non-genesis validator.
	VestingSchedule *VestingSchedule `json:"vesting_schedule"`
}

// Bid An entry in the validator map.
type Bid struct {
	// The purse was used for bonding.
	BondingPurse key.URef `json:"bonding_purse"`
	// The delegation rate.
	DelegationRate uint8 `json:"delegation_rate"`
	// `true` if validator has been "evicted"
	Inactive bool `json:"inactive"`
	// The amount of tokens staked by a validator (not including delegators).
	StakedAmount clvalue.UInt512 `json:"staked_amount"`
	// Validator's public key.
	ValidatorPublicKey keypair.PublicKey `json:"validator_public_key"`
	// The delegators.
	Delegators Delegators `json:"delegators"`
	// Vesting schedule for a genesis validator. `None` if non-genesis validator.
	VestingSchedule *VestingSchedule `json:"vesting_schedule"`
}

// Delegators the delegators type.
type Delegators []Delegator

func (d *Delegators) UnmarshalJSON(data []byte) error {
	if d == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}

	publicKeyAndDelegators := make([]struct {
		DelegatorPublicKey *keypair.PublicKey `json:"delegator_public_key"`
		Delegator          Delegator          `json:"delegator"`
	}, 0)

	err := json.Unmarshal(data, &publicKeyAndDelegators)
	if err == nil && len(publicKeyAndDelegators) > 0 && publicKeyAndDelegators[0].DelegatorPublicKey != nil {
		delegators := make(Delegators, 0, len(publicKeyAndDelegators))
		for _, item := range publicKeyAndDelegators {
			delegators = append(delegators, item.Delegator)
		}

		*d = delegators
		return nil
	}

	delegatorsV1 := make([]DelegatorV1, 0)
	if err := json.Unmarshal(data, &delegatorsV1); err != nil {
		return err
	}

	delegators := make(Delegators, 0, len(delegatorsV1))
	for _, item := range delegatorsV1 {
		delegators = append(delegators, NewDelegatorFromDelegatorV1(item))
	}

	*d = delegators
	return nil
}

// Delegator is associated with the given validator.
type Delegator struct {
	// The purse that was used for delegating.
	BondingPurse key.URef `json:"bonding_purse"`
	// Amount of Casper token (in motes) delegated
	StakedAmount clvalue.UInt512 `json:"staked_amount"`
	// Public Key of the delegator
	DelegatorPublicKey keypair.PublicKey `json:"delegator_public_key"`
	// Public key of the validator
	ValidatorPublicKey keypair.PublicKey `json:"validator_public_key"`
	// Vesting schedule for a genesis validator. `None` if non-genesis validator.
	VestingSchedule *VestingSchedule `json:"vesting_schedule"`
}

// DelegatorV1 of version 1 which is associated with the given validator.
type DelegatorV1 struct {
	// The purse that was used for delegating.
	BondingPurse key.URef `json:"bonding_purse"`
	// Amount of Casper token (in motes) delegated
	StakedAmount clvalue.UInt512 `json:"staked_amount"`
	// Public Key of the delegator
	Delegatee keypair.PublicKey `json:"delegatee"`
	// Public key of the validator
	ValidatorPublicKey keypair.PublicKey `json:"validator_public_key"`
	// Vesting schedule for a genesis validator. `None` if non-genesis validator.
	VestingSchedule *VestingSchedule `json:"vesting_schedule"`
}

func NewDelegatorFromDelegatorV1(v1 DelegatorV1) Delegator {
	return Delegator{
		BondingPurse:       v1.BondingPurse,
		StakedAmount:       v1.StakedAmount,
		DelegatorPublicKey: v1.Delegatee,
		ValidatorPublicKey: v1.ValidatorPublicKey,
		VestingSchedule:    v1.VestingSchedule,
	}
}

// Credit is a bridge record pointing to a new `ValidatorBid` after the public key was changed.
type Credit struct {
	// The era id the credit was created.
	EraID uint64 `json:"era_id"`
	// Validator's public key.
	ValidatorPublicKey keypair.PublicKey `json:"validator_public_key"`
	// The credit amount.
	Amount clvalue.UInt512 `json:"amount"`
}

// Bridge is a bridge record pointing to a new `ValidatorBid` after the public key was changed.
type Bridge struct {
	EraID uint64 `json:"era_id"`
	// Previous validator public key associated with the bid."
	OldValidatorPublicKey keypair.PublicKey `json:"old_validator_public_key"`
	// New validator public key associated with the bid.
	NewValidatorPublicKey keypair.PublicKey `json:"new_validator_public_key"`
}

// VestingSchedule for a genesis validator.
type VestingSchedule struct {
	InitialReleaseTimestampMillis uint64            `json:"initial_release_timestamp_millis"`
	LockedAmounts                 []clvalue.UInt512 `json:"locked_amounts"`
}
