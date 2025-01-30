package types

import (
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/key"
	"github.com/make-software/casper-go-sdk/v2/types/keypair"
)

// EraInfo stores an auction metadata. Intended to be recorded at each era.
type EraInfo struct {
	// List of rewards allocated to delegators and validators.
	SeigniorageAllocations []SeigniorageAllocation `json:"seigniorage_allocations"`
}

// SeigniorageAllocation sores information about a seigniorage allocation
type SeigniorageAllocation struct {
	Validator *ValidatorAllocation `json:"Validator,omitempty"`
	Delegator *DelegatorAllocation `json:"Delegator,omitempty"`
}

type ValidatorAllocation struct {
	// Public key of the validator
	ValidatorPublicKey keypair.PublicKey `json:"validator_public_key"`
	// Amount allocated as a reward.
	Amount clvalue.UInt512 `json:"amount"`
}

type DelegatorAllocation struct {
	// Public key of the delegator
	DelegatorKind DelegatorKind `json:"delegator_kind"`
	// Public key of the validator
	ValidatorPublicKey keypair.PublicKey `json:"validator_public_key"`
	// Amount allocated as a reward.
	Amount clvalue.UInt512 `json:"amount"`
}

func (t *DelegatorAllocation) UnmarshalJSON(data []byte) error {
	if t == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	temp := struct {
		// Public key of the delegator
		DelegatorKind *DelegatorKind `json:"delegator_kind"`
		// Public key of the validator
		DelegatorPublicKey *keypair.PublicKey `json:"delegator_public_key"`
		// Public key of the validator
		ValidatorPublicKey keypair.PublicKey `json:"validator_public_key"`
		// Amount allocated as a reward.
		Amount clvalue.UInt512 `json:"amount"`
	}{}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if temp.DelegatorKind != nil {
		*t = DelegatorAllocation{
			DelegatorKind:      *temp.DelegatorKind,
			ValidatorPublicKey: temp.ValidatorPublicKey,
			Amount:             temp.Amount,
		}
	} else if temp.DelegatorPublicKey != nil {
		*t = DelegatorAllocation{
			DelegatorKind: DelegatorKind{
				PublicKey: temp.DelegatorPublicKey,
			},
			ValidatorPublicKey: temp.ValidatorPublicKey,
			Amount:             temp.Amount,
		}
	} else {
		return errors.New("unexpected DelegatorAllocation format")
	}

	return nil
}

// DelegatorKind Auction bid variants. Kinds of delegation bids.
type DelegatorKind struct {
	// Delegation from public key.
	PublicKey *keypair.PublicKey `json:"PublicKey,omitempty"`
	// Delegation from purse.
	Purse *key.URef `json:"Purse,omitempty"`
}

func (t *DelegatorKind) ToHex() string {
	switch {
	case t.Purse != nil:
		return t.Purse.ToHex()
	case t.PublicKey != nil:
		return t.PublicKey.ToHex()
	}
	return ""
}

func (t *DelegatorKind) UnmarshalJSON(data []byte) error {
	if t == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	temp := struct {
		PublicKey *keypair.PublicKey `json:"PublicKey,omitempty"`
		// purse is represented not in format at uref-{uref-bytes}-{access}
		// but just a hex bytes
		Purse *string `json:"Purse,omitempty"`
	}{}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if temp.PublicKey != nil {
		*t = DelegatorKind{
			PublicKey: temp.PublicKey,
		}
	} else if temp.Purse != nil {
		urefBytes, err := hex.DecodeString(*temp.Purse)
		if err != nil {
			return err
		}

		// added one byte for default access
		uref, err := key.NewURefFromBytes(append(urefBytes, byte(7)))
		if err != nil {
			return err
		}

		*t = DelegatorKind{
			Purse: &uref,
		}
	} else {
		return errors.New("unexpected DelegatorKind format")
	}

	return nil
}
