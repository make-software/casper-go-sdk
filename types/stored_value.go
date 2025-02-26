package types

import (
	"encoding/json"
)

// StoredValue is a wrapper class for different types of values stored in the global state.
type StoredValue struct {
	CLValue           *Argument            `json:"CLValue,omitempty"`
	Account           *Account             `json:"Account,omitempty"`
	Contract          *Contract            `json:"Contract,omitempty"`
	ContractWASM      *json.RawMessage     `json:"ContractWASM,omitempty"`
	ContractPackage   *ContractPackage     `json:"ContractPackage,omitempty"`
	Transfer          *TransferV1          `json:"Transfer,omitempty"`
	DeployInfo        *DeployInfo          `json:"DeployInfo,omitempty"`
	EraInfo           *EraInfo             `json:"EraInfo,omitempty"`
	Bid               *Bid                 `json:"Bid,omitempty"`
	Withdraw          []UnbondingPurse     `json:"Withdraw,omitempty"`
	Unbonding         *UnbondingPurse      `json:"Unbonding,omitempty"`
	AddressableEntity *AddressableEntity   `json:"AddressableEntity,omitempty"`
	BidKind           *BidKind             `json:"BidKind,omitempty"`
	Package           *Package             `json:"Package,omitempty"`
	ByteCode          *ByteCode            `json:"ByteCode,omitempty"`
	MessageTopic      *MessageTopicSummary `json:"MessageTopic,omitempty"`
	Message           *MessageChecksum     `json:"Message,omitempty"`
	NamedKey          *NamedKeyValue       `json:"NamedKey,omitempty"`
	Prepayment        *PrepaymentKind      `json:"Prepayment,omitempty"`
	EntryPoint        *EntryPointValue     `json:"EntryPoint,omitempty"`
	RawBytes          *string              `json:"RawBytes,omitempty"`
}
