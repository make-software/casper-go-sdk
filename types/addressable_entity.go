package types

import (
	"github.com/make-software/casper-go-sdk/v2/types/key"
)

type AddressableEntity struct {
	// The type of Package.
	EntityKind  EntityKind `json:"entity_kind"`
	PackageHash string     `json:"package_hash"`
	// The hash address of the contract wasm
	ByteCodeHash string `json:"byte_code_hash"`
	// A collection of weighted public keys (represented as account hashes) associated with an account.
	AssociatedKeys []AssociatedKey `json:"associated_keys"`
	// Thresholds that have to be met when executing an action of a certain type.
	ActionThresholds EntityActionThresholds `json:"action_thresholds"`
	MainPurse        key.URef               `json:"main_purse"`
	// Casper Platform protocol version
	ProtocolVersion string `json:"protocol_version"`
}

type NamedEntryPoint struct {
	EntryPoint EntryPointV1 `json:"entry_point"`
	Name       string       `json:"name"`
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

type EntityKind struct {
	System        *SystemEntityType   `json:"System,omitempty"`
	Account       *key.AccountHash    `json:"Account,omitempty"`
	SmartContract *TransactionRuntime `json:"SmartContract,omitempty"`
}
