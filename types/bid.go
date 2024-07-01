package types

import (
	"encoding/json"
	"errors"

	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/key"
	"github.com/make-software/casper-go-sdk/types/keypair"
)

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
		DelegatorPublicKey keypair.PublicKey `json:"delegator_public_key"`
		Delegator          Delegator         `json:"delegator"`
	}, 0)

	if err := json.Unmarshal(data, &publicKeyAndDelegators); err == nil && len(publicKeyAndDelegators) > 0 {
		delegators := make(Delegators, 0, len(publicKeyAndDelegators))
		for _, item := range publicKeyAndDelegators {
			delegators = append(delegators, item.Delegator)
		}

		*d = delegators
		return nil
	}

	delegators := make([]Delegator, 0)
	if err := json.Unmarshal(data, &publicKeyAndDelegators); err != nil {
		return err
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

// DelegatorV1 is associated with the given validator for V1 network version.
type DelegatorV1 struct {
	// The purse that was used for delegating.
	BondingPurse key.URef `json:"bonding_purse"`
	// Amount of Casper token (in motes) delegated
	StakedAmount clvalue.UInt512 `json:"staked_amount"`
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
