package types

import (
	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/keypair"
)

type EraEnd struct {
	// The set of equivocators
	Equivocators []keypair.PublicKey `json:"equivocators"`
	// Validators that haven't produced any unit during the era
	InactiveValidators []keypair.PublicKey `json:"inactive_validators"`
	// The validators for the upcoming era and their respective weights
	NextEraValidatorWeights []ValidatorWeightEraEnd `json:"next_era_validator_weights"`
	// The rewards distributed to the validators
	Rewards         map[string][]clvalue.UInt512 `json:"rewards"`
	NextEraGasPrice uint8                        `json:"next_era_gas_price"`
}

func NewEraEndFromV2(eraEnd *EraEndV2) *EraEnd {
	if eraEnd == nil {
		return nil
	}

	return &EraEnd{
		NextEraGasPrice:         eraEnd.NextEraGasPrice,
		Equivocators:            eraEnd.Equivocators,
		InactiveValidators:      eraEnd.InactiveValidators,
		NextEraValidatorWeights: eraEnd.NextEraValidatorWeights,
		Rewards:                 eraEnd.Rewards,
	}
}

func NewEraEndFromV1(eraEnd *EraEndV1) *EraEnd {
	if eraEnd == nil {
		return nil
	}

	rewards := make(map[string][]clvalue.UInt512, len(eraEnd.EraReport.Rewards))
	for _, reward := range eraEnd.EraReport.Rewards {
		list := rewards[reward.Validator.ToHex()]
		list = append(list, reward.Amount)
		rewards[reward.Validator.ToHex()] = list
	}

	return &EraEnd{
		NextEraGasPrice:         1,
		Equivocators:            eraEnd.EraReport.Equivocators,
		InactiveValidators:      eraEnd.EraReport.InactiveValidators,
		NextEraValidatorWeights: eraEnd.NextEraValidatorWeights,
		Rewards:                 rewards,
	}
}

// EraEndV2 information related to the end of an era, and validator weights for the following era
type EraEndV2 struct {
	// The set of equivocators
	Equivocators []keypair.PublicKey `json:"equivocators"`
	// Validators that haven't produced any unit during the era
	InactiveValidators []keypair.PublicKey `json:"inactive_validators"`
	// The validators for the upcoming era and their respective weights
	NextEraValidatorWeights []ValidatorWeightEraEnd `json:"next_era_validator_weights"`
	// The rewards distributed to the validators
	Rewards         map[string][]clvalue.UInt512 `json:"rewards"`
	NextEraGasPrice uint8                        `json:"next_era_gas_price"`
}

// EraEndV1 contains a report and list of validator weights for the next era
type EraEndV1 struct {
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
	Rewards            []EraReward         `json:"rewards"`
}

type EraReward struct {
	Validator keypair.PublicKey `json:"validator"`
	Amount    clvalue.UInt512   `json:"amount"`
}
