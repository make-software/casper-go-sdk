package types

import "github.com/make-software/casper-go-sdk/v2/types/key"

// Package Entity definition, metadata, and security container.
type Package struct {
	// All versions (enabled & disabled).
	Versions []EntityVersionAndHash `json:"versions"`
	// Collection of disabled entity versions. The runtime will not permit disabled entity versions to be executed.
	DisabledVersions []EntityVersionAndHash `json:"disabled_versions"`
	// A flag that determines whether a entity is locked
	LockStatus string `json:"lock_status"`
	// Mapping maintaining the set of URefs associated with each "user group"
	Groups []string `json:"groups"`
}

type EntityVersionAndHash struct {
	AddressableEntityHash key.AddressableEntityHash `json:"addressable_entity_hash"`
	EntityVersionKey      EntityVersionKey          `json:"entity_version_key"`
}

// EntityVersionKey Major element of `ProtocolVersion` combined with `EntityVersion`.
type EntityVersionKey struct {
	// Automatically incremented value for a contract version within a major `ProtocolVersion`.
	EntityVersion uint32 `json:"entity_version"`
	// Major element of `ProtocolVersion` a `ContractVersion` is compatible with.
	ProtocolVersionMajor uint32 `json:"protocol_version_major"`
}

type NamedUserGroup struct {
	// A (labelled) "user group". Each method of a versioned contract may be associated with one or more user groups which are allowed to call it.
	GroupName  string     `json:"group_name"`
	GroupUsers []key.URef `json:"group_users"`
}
