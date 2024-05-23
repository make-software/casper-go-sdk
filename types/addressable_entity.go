package types

import (
	"encoding/json"
	"errors"
	"github.com/make-software/casper-go-sdk/types/key"
)

var ErrUnknownEntityKind = errors.New("unknown entity kind")

type AddressableEntity struct {
	// The type of Package.
	EntityKind  EntityKind `json:"entity_kind"`
	PackageHash string     `json:"package_hash"`
	// The hash address of the contract wasm
	ByteCodeHash string            `json:"byte_code_hash"`
	EntryPoints  []NamedEntryPoint `json:"entry_points"`
	// A collection of weighted public keys (represented as account hashes) associated with an account.
	AssociatedKeys []AssociatedKey `json:"associated_keys"`
	// Thresholds that have to be met when executing an action of a certain type.
	ActionThresholds EntityActionThresholds `json:"action_thresholds"`
	MainPurse        key.URef               `json:"main_purse"`
	// Casper Platform protocol version
	ProtocolVersion string `json:"protocol_version"`

	MessageTopics []MessageTopic `json:"message_topics"`
}

type MessageTopic struct {
	TopicName     string   `json:"topic_name"`
	TopicNameHash key.Hash `json:"topic_name_hash"`
}

type NamedEntryPoint struct {
	EntryPoint EntryPoint `json:"entry_point"`
	Name       string     `json:"name"`
}

// EntityActionThresholds Thresholds that have to be met when executing an action of a certain type.
type EntityActionThresholds struct {
	// Threshold for deploy execution.
	Deployment uint64 `json:"deployment"`
	// Threshold for upgrading contracts.
	UpgradeManagement uint64 `json:"upgrade_management"`
	// Threshold for managing action threshold.
	KeyManagement uint64 `json:"key_management"`
}

// SystemEntityType System contract types.
type SystemEntityType string

const SmartContractEntityKind = "SmartContract"

type EntityKind struct {
	System        *SystemEntityType `json:"System,omitempty"`
	Account       *key.AccountHash  `json:"Account,omitempty"`
	SmartContract bool              `json:"SmartContract,omitempty"`
}

// EntityKind represent a complex type
// 1. Account ->         "entity_kind": { "Account": "account-hash-7dd64920e60864390c810182b83b53f49310adc8d66e714c57a6e5ff0e3c6552" }
// 2. SmartContract ->  "entity_kind": "SmartContract"
// 2. System ->         "entity_kind": { "System": "Auction" }

func (k *EntityKind) UnmarshalJSON(data []byte) error {
	var kind struct {
		System  *SystemEntityType `json:"System,omitempty"`
		Account *key.AccountHash  `json:"Account,omitempty"`
	}
	if err := json.Unmarshal(data, &kind); err == nil {
		*k = EntityKind{
			System:  kind.System,
			Account: kind.Account,
		}
		return nil
	}

	var value string
	if err := json.Unmarshal(data, &value); err == nil {
		if value == SmartContractEntityKind {
			*k = EntityKind{
				SmartContract: true,
			}
			return nil
		}
	}

	return ErrUnknownEntityKind
}

func (k EntityKind) MarshalJSON() ([]byte, error) {
	if k.System != nil {
		return json.Marshal(map[string]*SystemEntityType{
			"System": k.System,
		})
	}
	if k.Account != nil {
		return json.Marshal(map[string]*key.AccountHash{
			"Account": k.Account,
		})
	}
	if k.SmartContract {
		return json.Marshal("SmartContract")
	}
	return nil, ErrUnknownEntityKind
}
