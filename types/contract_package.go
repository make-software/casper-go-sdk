package types

import (
	"encoding/json"
	"fmt"

	"github.com/make-software/casper-go-sdk/v2/types/key"
)

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
	GroupName string `json:"group_name"`
	// List of URefs associated with the group label.
	GroupUsers []key.URef `json:"group_users"`
}

func (c *ContractGroup) UnmarshalJSON(data []byte) error {
	var temp struct {
		Group      string     `json:"group"`
		GroupName  string     `json:"group_name"`
		Keys       []key.URef `json:"keys"`
		GroupUsers []key.URef `json:"group_users"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	var group = ContractGroup{
		GroupName:  temp.GroupName,
		GroupUsers: temp.GroupUsers,
	}

	if temp.Group != "" {
		group.GroupName = temp.Group
	}

	if len(temp.Keys) != 0 {
		group.GroupUsers = temp.Keys
	}

	*c = group
	return nil
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

func (c *ContractVersionKey) UnmarshalJSON(data []byte) error {
	var key []int
	if err := json.Unmarshal(data, &key); err == nil {
		if len(key) != 2 {
			return fmt.Errorf("contract version key must contain exactly 2 elements")
		}
		*c = ContractVersionKey{key[0], key[1]}
		return nil
	}

	var value struct {
		Version              *int `json:"contract_version"`
		ProtocolVersionMajor *int `json:"protocol_version_major"`
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	if value.Version == nil || value.ProtocolVersionMajor == nil {
		return fmt.Errorf("contract version key must contain contract_version and protocol_version_major")
	}

	*c = ContractVersionKey{*value.ProtocolVersionMajor, *value.Version}
	return nil
}
