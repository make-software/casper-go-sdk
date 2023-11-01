package types

import (
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/keypair"
)

// EraEnd contains a report and list of validator weights for the next era
type EraEnd struct {
	EraReport EraReport `json:"era_report"`
	// A list of validator weights for the next era
	NextEraValidatorWeights []ValidatorWeightEraEnd `json:"next_era_validator_weights"`
}

// EraReport is an equivocation and reward information to be included in the terminal block.
type EraReport struct {
	// List of public keys of the equivocators
	Equivocators []keypair.PublicKey `json:"equivocators"`
	// List of public keys of inactive validators
	InactiveValidators []keypair.PublicKey `json:"inactive_validators"`
	Rewards            []Reward            `json:"rewards"`
}

type Reward struct {
	Validator keypair.PublicKey `json:"validator"`
	Amount    clvalue.UInt512   `json:"amount"`
}
