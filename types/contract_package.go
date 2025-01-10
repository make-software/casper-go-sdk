package types

import (
	"encoding/json"

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
	key := make([]int, 0, 2)
	if err := json.Unmarshal(data, &key); err != nil {
		var temp struct {
			Version              int `json:"contract_version"`
			ProtocolVersionMajor int `json:"protocol_version_major"`
		}
		if err = json.Unmarshal(data, &temp); err != nil {
			return err
		}

		if temp.Version != 0 && temp.ProtocolVersionMajor != 0 {
			*c = [2]int{temp.Version, temp.ProtocolVersionMajor}
			return nil
		}
	}
	if len(key) == 2 {
		*c = [2]int{key[0], key[1]}
		return nil
	}

	return nil
}
