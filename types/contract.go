package types

import "github.com/make-software/casper-go-sdk/types/key"

// Contract contains information, entry points and named keys belonging to a Contract
type Contract struct {
	// Key to the storage of the ContractPackage object
	ContractPackageHash key.ContractPackageHash `json:"contract_package_hash"`
	// Key to the storage of the ContractWasm object
	ContractWasmHash key.ContractHash `json:"contract_wasm_hash"`
	// List of entry points or methods in the contract.
	EntryPoints []EntryPointV1 `json:"entry_points"`
	// NamedKeys are a collection of String-Key pairs used to easily identify some data on the network.
	NamedKeys NamedKeys `json:"named_keys"`
	// The protocol version when the contract was deployed
	ProtocolVersion string `json:"protocol_version"`
}
