package types

import (
	"encoding/hex"
	"math/big"

	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/key"
)

type ExecutableDeployItemType byte

const (
	ExecutableDeployItemTypeModuleBytes ExecutableDeployItemType = iota
	ExecutableDeployItemTypeStoredContractByHash
	ExecutableDeployItemTypeStoredContractByName
	ExecutableDeployItemTypeStoredVersionedContractByHash
	ExecutableDeployItemTypeStoredVersionedContractByName
	ExecutableDeployItemTypeTransfer
)

// ExecutableDeployItem is a base structure for the possible variants of an executable deploy.
// This structure should contain only one object.
type ExecutableDeployItem struct {
	ModuleBytes                   *ModuleBytes                   `json:"ModuleBytes,omitempty"`
	StoredContractByHash          *StoredContractByHash          `json:"StoredContractByHash,omitempty"`
	StoredContractByName          *StoredContractByName          `json:"StoredContractByName,omitempty"`
	StoredVersionedContractByHash *StoredVersionedContractByHash `json:"StoredVersionedContractByHash,omitempty"`
	StoredVersionedContractByName *StoredVersionedContractByName `json:"StoredVersionedContractByName,omitempty"`
	Transfer                      *TransferDeployItem            `json:"Transfer,omitempty"`
}

func (e ExecutableDeployItem) Bytes() ([]byte, error) {
	if e.ModuleBytes != nil {
		bytes, err := e.ModuleBytes.Bytes()
		if err != nil {
			return nil, err
		}
		return append([]byte{byte(ExecutableDeployItemTypeModuleBytes)}, bytes...), nil
	} else if e.StoredContractByHash != nil {
		bytes, err := e.StoredContractByHash.Bytes()
		if err != nil {
			return nil, err
		}
		return append([]byte{byte(ExecutableDeployItemTypeStoredContractByHash)}, bytes...), nil
	} else if e.StoredContractByName != nil {
		bytes, err := e.StoredContractByName.Bytes()
		if err != nil {
			return nil, err
		}
		return append([]byte{byte(ExecutableDeployItemTypeStoredContractByName)}, bytes...), nil
	} else if e.StoredVersionedContractByHash != nil {
		bytes, err := e.StoredVersionedContractByHash.Bytes()
		if err != nil {
			return nil, err
		}
		return append([]byte{byte(ExecutableDeployItemTypeStoredVersionedContractByHash)}, bytes...), nil
	} else if e.StoredVersionedContractByName != nil {
		bytes, err := e.StoredVersionedContractByName.Bytes()
		if err != nil {
			return nil, err
		}
		return append([]byte{byte(ExecutableDeployItemTypeStoredVersionedContractByName)}, bytes...), nil
	} else if e.Transfer != nil {
		bytes, err := e.Transfer.Bytes()
		if err != nil {
			return nil, err
		}
		return append([]byte{byte(ExecutableDeployItemTypeTransfer)}, bytes...), nil
	}
	return []byte{}, nil
}

// ModuleBytes is a `deploy` item with the capacity to contain executable code (e.g. a contract).
type ModuleBytes struct {
	// WASM Bytes
	ModuleBytes string `json:"module_bytes"`
	Args        *Args  `json:"args,omitempty"`
}

func (m ModuleBytes) Bytes() ([]byte, error) {
	bytes, _ := hex.DecodeString(m.ModuleBytes)
	res := clvalue.NewCLUInt32(uint32(len(bytes))).Bytes()
	res = append(res, clvalue.NewCLByteArray(bytes).Bytes()...)
	argBytes, err := m.Args.Bytes()
	if err != nil {
		return nil, err
	}
	return append(res, argBytes...), nil
}

// StoredContractByHash is a `Deploy` item to call an entry point in a contract. The contract is referenced by its hash.
type StoredContractByHash struct {
	Hash       key.ContractHash `json:"hash"`
	EntryPoint string           `json:"entry_point"`
	Args       *Args            `json:"args"`
}

func (m StoredContractByHash) Bytes() ([]byte, error) {
	argBytes, err := m.Args.Bytes()
	if err != nil {
		return nil, err
	}
	return append(m.Hash.Bytes(), append(clvalue.NewCLString(m.EntryPoint).Bytes(), argBytes...)...), nil
}

// StoredContractByName is a `Deploy` item to call an entry point in a contract. The contract is referenced
// by a named key in the caller account pointing to the contract hash.
type StoredContractByName struct {
	Name       string `json:"name"`
	EntryPoint string `json:"entry_point"`
	Args       *Args  `json:"args"`
}

func (m StoredContractByName) Bytes() ([]byte, error) {
	argBytes, err := m.Args.Bytes()
	if err != nil {
		return nil, err
	}

	return append(
		clvalue.NewCLString(m.Name).Bytes(),
		append(clvalue.NewCLString(m.EntryPoint).Bytes(), argBytes...)...,
	), nil
}

// StoredVersionedContractByHash is a `Deploy` item to call an entry point in a contract. The contract is referenced
// by a contract package hash and a version number.
type StoredVersionedContractByHash struct {
	// Hash of the contract.
	Hash key.ContractHash `json:"hash"`
	// Entry point or method of the contract to call.
	EntryPoint string  `json:"entry_point"`
	Version    *string `json:"version,omitempty"`
	Args       *Args   `json:"args"`
}

func (m StoredVersionedContractByHash) Bytes() ([]byte, error) {
	option := clvalue.Option{}
	if m.Version != nil || *m.Version != "" {
		option.Inner = clvalue.NewCLString(*m.Version)
	}
	argBytes, err := m.Args.Bytes()
	if err != nil {
		return nil, err
	}
	return append(
		m.Hash.Bytes(),
		append(option.Bytes(), append(clvalue.NewCLString(m.EntryPoint).Bytes(), argBytes...)...)...,
	), nil
}

// StoredVersionedContractByName is a `Deploy` item to call an entry point in a contract. The contract is referenced
// by a named key in the caller account pointing to the contract package hash
// and a version number.
type StoredVersionedContractByName struct {
	// Name of a named key in the caller account that stores the contract package hash.
	Name string `json:"name"`
	// Entry point or method of the contract to call.
	EntryPoint string  `json:"entry_point"`
	Version    *string `json:"version,omitempty"`
	Args       *Args   `json:"args"`
}

func (m StoredVersionedContractByName) Bytes() ([]byte, error) {
	option := clvalue.Option{}
	if m.Version != nil || *m.Version != "" {
		option.Inner = clvalue.NewCLString(*m.Version)
	}
	argBytes, err := m.Args.Bytes()
	if err != nil {
		return nil, err
	}
	return append(
		clvalue.NewCLString(m.Name).Bytes(),
		append(option.Bytes(), append(clvalue.NewCLString(m.EntryPoint).Bytes(), argBytes...)...)...,
	), nil
}

// TransferDeployItem is a `Deploy` item for transferring funds to a target account.
type TransferDeployItem struct {
	Args Args `json:"args"`
}

func (m TransferDeployItem) Bytes() ([]byte, error) {
	return m.Args.Bytes()
}

func StandardPayment(amount *big.Int) ExecutableDeployItem {
	return ExecutableDeployItem{
		ModuleBytes: &ModuleBytes{
			Args: (&Args{}).AddArgument("amount", *clvalue.NewCLUInt512(amount)),
		},
	}
}
