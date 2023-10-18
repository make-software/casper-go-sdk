package types

import (
	"encoding/json"
)

// StoredValue is a wrapper class for different types of values stored in the global state.
type StoredValue struct {
	CLValue         *Argument        `json:"CLValue,omitempty"`
	Account         *Account         `json:"Account,omitempty"`
	Contract        *Contract        `json:"Contract,omitempty"`
	ContractWASM    *json.RawMessage `json:"ContractWASM,omitempty"`
	ContractPackage *ContractPackage `json:"ContractPackage,omitempty"`
	Transfer        *Transfer        `json:"TransferDeployItem,omitempty"`
	DeployInfo      *DeployInfo      `json:"DeployInfo,omitempty"`
	EraInfo         *EraInfo         `json:"EraInfo,omitempty"`
	Bid             *Bid             `json:"Bid,omitempty"`
	Withdraw        []UnbondingPurse `json:"Withdraw,omitempty"`
}
