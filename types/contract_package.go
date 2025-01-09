package types

import "github.com/make-software/casper-go-sdk/v2/types/key"

// ContractPackage contains contract definition, metadata, and security container.
type ContractPackage struct {
	// Access key for this contract.
	AccessKey key.URef `json:"access_key"`
	// List of disabled versions of a contract.
	DisabledVersions []ContractVersionKey `json:"disabled_versions"`
	// Groups associate a set of URefs with a label. Entry points on a contract can be given
	// a list of labels they accept and the runtime will check that a URef from at least one
	// of the allowed groups is present in the caller’s context before execution.
	Groups []ContractGroup `json:"groups"`
	// List of active versions of a contract.
	Versions   []ContractVersion `json:"versions"`
	LockStatus string            `json:"lock_status"`
}

// ContractGroup associate a set of URefs with a label.
type ContractGroup struct {
	// Group label
	Group string `json:"group"`
	// List of URefs associated with the group label.
	Keys []key.URef `json:"keys"`
}

// ContractVersion contains information related to an active version of a contract.
type ContractVersion struct {
	// Hash for this version of the contract.
	Hash key.ContractHash `json:"contract_hash"`
	// Contract version.
	Version uint16 `json:"contract_version"`
	//  The major element of the protocol version this contract is compatible with.
	ProtocolVersionMajor uint16 `json:"protocol_version_major"`
}

// ContractVersionKey Major element of `ProtocolVersion` combined with `ContractVersion`.
type ContractVersionKey [2]int
