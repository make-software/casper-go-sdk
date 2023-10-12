package key

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type PrefixName = string

const (
	PrefixNameAccount                PrefixName = "account-hash-"
	PrefixNameHash                   PrefixName = "hash-"
	PrefixNameContractPackageWasm    PrefixName = "contract-package-wasm"
	PrefixNameContractPackage        PrefixName = "contract-package-"
	PrefixNameContractWasm           PrefixName = "contract-wasm-"
	PrefixNameContract               PrefixName = "contract-"
	PrefixNameURef                   PrefixName = "uref-"
	PrefixNameTransfer               PrefixName = "transfer-"
	PrefixNameDeployInfo             PrefixName = "deploy-"
	PrefixNameEraId                  PrefixName = "era-"
	PrefixNameBid                    PrefixName = "bid-"
	PrefixNameBalance                PrefixName = "balance-"
	PrefixNameWithdraw               PrefixName = "withdraw-"
	PrefixNameDictionary             PrefixName = "dictionary-"
	PrefixNameSystemContractRegistry PrefixName = "system-contract-registry-"
	PrefixNameEraSummary             PrefixName = "era-summary-"
	PrefixNameUnbond                 PrefixName = "unbond-"
	PrefixNameChainspecRegistry      PrefixName = "chainspec-registry-"
	PrefixNameChecksumRegistry       PrefixName = "checksum-registry-"
)

var ErrNotFoundPrefix = errors.New("prefix is not found")

type TypeID = byte

const (
	TypeIDAccount TypeID = iota
	TypeIDHash
	TypeIDURef
	TypeIDTransfer
	TypeIDDeployInfo
	TypeIDEraId
	TypeIDBalance
	TypeIDBid
	TypeIDWithdraw
	TypeIDDictionary
	TypeIDSystemContractRegistry
	TypeIDEraSummary
	TypeIDUnbond
	TypeIDChainspecRegistry
	TypeIDChecksumRegistry
)

type TypeName = string

const (
	TypeNameAccount                TypeName = "Account"
	TypeNameHash                   TypeName = "Hash"
	TypeNameURef                   TypeName = "URef"
	TypeNameTransfer               TypeName = "Transfer"
	TypeNameDeployInfo             TypeName = "Deploy"
	TypeNameEraId                  TypeName = "Era"
	TypeNameBid                    TypeName = "Bid"
	TypeNameBalance                TypeName = "Balance"
	TypeNameWithdraw               TypeName = "Withdraw"
	TypeNameDictionary             TypeName = "Dictionary"
	TypeNameSystemContractRegistry TypeName = "SystemContractRegistry"
	TypeNameEraSummary             TypeName = "EraSummary"
	TypeNameUnbond                 TypeName = "Unbond"
	TypeNameChainspecRegistry      TypeName = "ChainspecRegistry"
	TypeNameChecksumRegistry       TypeName = "ChecksumRegistry"
)

var typeIDbyNames = map[TypeName]TypeID{
	TypeNameAccount:                TypeIDAccount,
	TypeNameHash:                   TypeIDHash,
	TypeNameURef:                   TypeIDURef,
	TypeNameTransfer:               TypeIDTransfer,
	TypeNameDeployInfo:             TypeIDDeployInfo,
	TypeNameEraId:                  TypeIDEraId,
	TypeNameBid:                    TypeIDBid,
	TypeNameBalance:                TypeIDBalance,
	TypeNameWithdraw:               TypeIDWithdraw,
	TypeNameDictionary:             TypeIDDictionary,
	TypeNameSystemContractRegistry: TypeIDSystemContractRegistry,
	TypeNameEraSummary:             TypeIDEraSummary,
	TypeNameUnbond:                 TypeIDUnbond,
	TypeNameChainspecRegistry:      TypeIDChainspecRegistry,
	TypeNameChecksumRegistry:       TypeIDChecksumRegistry,
}

var keyIDbyPrefix = map[PrefixName]TypeID{
	PrefixNameAccount:                TypeIDAccount,
	PrefixNameHash:                   TypeIDHash,
	PrefixNameTransfer:               TypeIDTransfer,
	PrefixNameURef:                   TypeIDURef,
	PrefixNameDeployInfo:             TypeIDDeployInfo,
	PrefixNameEraId:                  TypeIDEraId,
	PrefixNameBid:                    TypeIDBid,
	PrefixNameBalance:                TypeIDBalance,
	PrefixNameWithdraw:               TypeIDWithdraw,
	PrefixNameDictionary:             TypeIDDictionary,
	PrefixNameSystemContractRegistry: TypeIDSystemContractRegistry,
	PrefixNameEraSummary:             TypeIDEraSummary,
	PrefixNameUnbond:                 TypeIDUnbond,
	PrefixNameChainspecRegistry:      TypeIDChainspecRegistry,
	PrefixNameChecksumRegistry:       TypeIDChecksumRegistry,
}

type Key struct {
	Type TypeID
	// A `Key` under which a user account is stored.
	Account *AccountHash
	// A `Key` under which a smart contract is stored and which is the pseudo-hash of the contract.
	Hash *Hash
	// A `Key` which is a [`URef`], under which most types of data can be stored.
	URef *URef
	// A `Key` under which we store a transfer.
	Transfer *TransferHash
	// A `Key` under which we store a deploy info.
	Deploy *Hash
	// A `Key` under which we store an era info.
	Era *Era
	// A `Key` under which we store a purse balance.
	Balance *Hash
	// A `Key` under which we store bid information
	Bid *AccountHash
	// A `Key` under which we store withdraw information.
	Withdraw *AccountHash
	// A `Key` variant whose value is derived by hashing [`URef`]s address and arbitrary data.
	Dictionary *Hash
	// A `Key` variant under which system contract hashes are stored.
	SystemContactRegistry *Hash
	// A `Key` under which we store current era info.
	EraSummary *Hash
	// A `Key` under which we store unbond information.
	Unbond *AccountHash
	// A `Key` variant under which chainspec and other hashes are stored.
	ChainspecRegistry *Hash
	// A `Key` variant under which we store a registry of checksums.
	ChecksumRegistry *Hash
}

func (k Key) Bytes() []byte {
	switch k.Type {
	case TypeIDBalance:
		return append([]byte{TypeIDBalance}, k.Balance.Bytes()...)
	case TypeIDBid:
		return append([]byte{TypeIDBid}, k.Bid.Bytes()...)
	case TypeIDWithdraw:
		return append([]byte{TypeIDWithdraw}, k.Withdraw.Bytes()...)
	case TypeIDSystemContractRegistry:
		return append([]byte{TypeIDSystemContractRegistry}, k.SystemContactRegistry.Bytes()...)
	case TypeIDUnbond:
		return append([]byte{TypeIDUnbond}, k.Unbond.Bytes()...)
	case TypeIDChainspecRegistry:
		return append([]byte{TypeIDChainspecRegistry}, k.ChainspecRegistry.Bytes()...)
	case TypeIDChecksumRegistry:
		return append([]byte{TypeIDChecksumRegistry}, k.ChecksumRegistry.Bytes()...)
	case TypeIDEraSummary:
		return append([]byte{TypeIDEraSummary}, k.EraSummary.Bytes()...)
	case TypeIDAccount:
		return append([]byte{TypeIDAccount}, k.Account.Bytes()...)
	case TypeIDHash:
		return append([]byte{TypeIDHash}, k.Hash.Bytes()...)
	case TypeIDEraId:
		return append([]byte{TypeIDEraId}, k.Era.Bytes()...)
	case TypeIDURef:
		return append([]byte{TypeIDURef}, k.URef.Bytes()...)
	case TypeIDTransfer:
		return append([]byte{TypeIDTransfer}, k.Transfer.Bytes()...)
	case TypeIDDeployInfo:
		return append([]byte{TypeIDDeployInfo}, k.Deploy.Bytes()...)
	case TypeIDDictionary:
		return append([]byte{TypeIDDictionary}, k.Dictionary.Bytes()...)

	default:
		return []byte{}
	}
}

func (k Key) String() string {
	return k.ToPrefixedString()
}

func (k Key) ToPrefixedString() string {
	switch k.Type {
	case TypeIDAccount:
		return k.Account.ToPrefixedString()
	case TypeIDHash:
		return PrefixNameHash + k.Hash.ToHex()
	case TypeIDEraId:
		return PrefixNameEraId + strconv.Itoa(int(*k.Era))
	case TypeIDURef:
		return k.URef.ToPrefixedString()
	case TypeIDTransfer:
		return k.Transfer.ToPrefixedString()
	case TypeIDDeployInfo:
		return PrefixNameDeployInfo + k.Deploy.ToHex()
	case TypeIDDictionary:
		return PrefixNameDictionary + k.Dictionary.ToHex()
	case TypeIDBalance:
		return PrefixNameBalance + k.Balance.ToHex()
	case TypeIDBid:
		return PrefixNameBid + k.Bid.ToHex()
	case TypeIDWithdraw:
		return PrefixNameWithdraw + k.Withdraw.ToHex()
	case TypeIDSystemContractRegistry:
		return PrefixNameSystemContractRegistry + k.SystemContactRegistry.ToHex()
	case TypeIDEraSummary:
		return PrefixNameEraSummary + k.EraSummary.ToHex()
	case TypeIDUnbond:
		return PrefixNameUnbond + k.Unbond.ToHex()
	case TypeIDChainspecRegistry:
		return PrefixNameChainspecRegistry + k.ChainspecRegistry.ToHex()
	case TypeIDChecksumRegistry:
		return PrefixNameChecksumRegistry + k.ChecksumRegistry.ToHex()
	default:
		return ""
	}
}

func (k Key) MarshalJSON() ([]byte, error) {
	return []byte(`"` + k.ToPrefixedString() + `"`), nil
}

func (k *Key) UnmarshalJSON(i []byte) error {
	var s string
	err := json.Unmarshal(i, &s)
	if err != nil {
		return err
	}
	tmp, err := NewKey(s)
	if err != nil {
		return err
	}
	*k = tmp
	return nil
}

func (k *Key) Scan(value any) (err error) {
	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid scan value type")
	}

	dst := make([]byte, len(data))
	copy(dst, data)

	*k, err = NewKeyFromBuffer(bytes.NewBuffer(dst))
	return err
}

func (k Key) Value() (driver.Value, error) {
	return k.Bytes(), nil
}

func findPrefixByMap(source string, prefixes map[PrefixName]TypeID) PrefixName {
	var prefix PrefixName
	for one := range prefixes {
		if strings.HasPrefix(source, one) {
			// handle the special case when prefix era- is the part of the prefix era-summary-
			if one == PrefixNameEraId && strings.HasPrefix(source, PrefixNameEraSummary) {
				return PrefixNameEraSummary
			}
			prefix = one
			break
		}
	}
	return prefix
}

func NewKey(source string) (result Key, err error) {
	if len(source) == StingHashLen {
		defaultHash, err := NewHash(source)
		result.Type = TypeIDHash
		result.Hash = &defaultHash
		return result, err
	}
	// In the Rust implementation these are prefixes for display, that shouldn't use in API
	if strings.HasPrefix(source, "Key::") {
		return parseTypeByString(source)
	}
	if strings.HasPrefix(source, "00") && len(source) == (StingHashLen+2) {
		return createByType(source[2:], TypeIDAccount)
	}
	prefix := findPrefixByMap(source, keyIDbyPrefix)
	if prefix == "" {
		return result, fmt.Errorf("%w, source: %s", ErrNotFoundPrefix, source)
	}
	return createByType(source, keyIDbyPrefix[prefix])
}

// And with text prefixes Example: "Key::AccountHash(<hash>)", "Key::Hash(<hash>)"
func parseTypeByString(source string) (result Key, err error) {
	openBracketIndex := strings.Index(source, "(")
	if openBracketIndex == -1 {
		return result, errors.New("invalid key format")
	}

	const typeIndexStart = 5
	keyTypeStr := source[typeIndexStart:openBracketIndex]
	keyValue := source[openBracketIndex+1 : len(source)-1]
	typeID, ok := typeIDbyNames[keyTypeStr]
	if !ok {
		return result, errors.New("unexpected KeyType")
	}

	return createByType(keyValue, typeID)
}

func createByType(source string, typeID TypeID) (result Key, err error) {
	result.Type = typeID
	switch result.Type {
	case TypeIDEraId:
		hash, err := NewEraFromString(strings.TrimPrefix(source, PrefixNameEraId))
		result.Era = &hash
		return result, err
	case TypeIDHash:
		hash, err := NewHash(strings.TrimPrefix(source, PrefixNameHash))
		result.Hash = &hash
		return result, err
	case TypeIDURef:
		hash, err := NewURef(source)
		result.URef = &hash
		return result, err
	case TypeIDAccount:
		hash, err := NewAccountHash(source)
		result.Account = &hash
		return result, err
	case TypeIDTransfer:
		hash, err := NewTransferHash(source)
		result.Transfer = &hash
		return result, err
	case TypeIDDeployInfo:
		hash, err := NewHash(strings.TrimPrefix(source, PrefixNameDeployInfo))
		result.Deploy = &hash
		return result, err
	case TypeIDBalance:
		hash, err := NewHash(strings.TrimPrefix(source, PrefixNameBalance))
		result.Balance = &hash
		return result, err
	case TypeIDBid:
		hash, err := NewAccountHash(strings.TrimPrefix(source, PrefixNameBid))
		result.Bid = &hash
		return result, err
	case TypeIDWithdraw:
		hash, err := NewAccountHash(strings.TrimPrefix(source, PrefixNameWithdraw))
		result.Withdraw = &hash
		return result, err
	case TypeIDDictionary:
		hash, err := NewHash(strings.TrimPrefix(source, PrefixNameDictionary))
		result.Dictionary = &hash
		return result, err
	case TypeIDSystemContractRegistry:
		hash, err := NewHash(strings.TrimPrefix(source, PrefixNameSystemContractRegistry))
		result.SystemContactRegistry = &hash
		return result, err
	case TypeIDEraSummary:
		hash, err := NewHash(strings.TrimPrefix(source, PrefixNameEraSummary))
		result.EraSummary = &hash
		return result, err
	case TypeIDUnbond:
		data, err := NewAccountHash(strings.TrimPrefix(source, PrefixNameUnbond))
		result.Unbond = &data
		return result, err
	case TypeIDChainspecRegistry:
		data, err := NewHash(strings.TrimPrefix(source, PrefixNameChainspecRegistry))
		result.ChainspecRegistry = &data
		return result, err
	case TypeIDChecksumRegistry:
		data, err := NewHash(strings.TrimPrefix(source, PrefixNameChecksumRegistry))
		result.ChecksumRegistry = &data
		return result, err
	default:
		err = errors.New("type is not found")
	}
	return result, err
}

func NewKeyFromBytes(data []byte) (Key, error) {
	return NewKeyFromBuffer(bytes.NewBuffer(data))
}

func NewKeyFromBuffer(buffer *bytes.Buffer) (result Key, err error) {
	keyType, err := buffer.ReadByte()
	if err != nil {
		return result, err
	}
	result.Type = keyType
	switch keyType {
	case TypeIDEraId:
		result.Era, err = NewEraFromBuffer(buffer)
		return result, err
	case TypeIDAccount:
		data, err := NewByteHashFromBuffer(buffer)
		result.Account = &AccountHash{Hash: data}
		return result, err
	case TypeIDHash:
		data, err := NewByteHashFromBuffer(buffer)
		result.Hash = &data
		return result, err
	case TypeIDTransfer:
		data, err := NewByteHashFromBuffer(buffer)
		result.Transfer = &TransferHash{Hash: data}
		return result, err
	case TypeIDDeployInfo:
		data, err := NewByteHashFromBuffer(buffer)
		result.Deploy = &data
		return result, err
	case TypeIDBalance:
		data, err := NewByteHashFromBuffer(buffer)
		result.Balance = &data
		return result, err
	case TypeIDBid:
		data, err := NewByteHashFromBuffer(buffer)
		result.Bid = &AccountHash{Hash: data}
		return result, err
	case TypeIDWithdraw:
		data, err := NewByteHashFromBuffer(buffer)
		result.Withdraw = &AccountHash{Hash: data}
		return result, err
	case TypeIDDictionary:
		data, err := NewByteHashFromBuffer(buffer)
		result.Dictionary = &data
		return result, err
	case TypeIDSystemContractRegistry:
		data, err := NewByteHashFromBuffer(buffer)
		result.SystemContactRegistry = &data
		return result, err
	case TypeIDURef:
		data, err := NewURefFromBuffer(buffer)
		result.URef = &data
		return result, err
	case TypeIDEraSummary:
		data, err := NewByteHashFromBuffer(buffer)
		result.EraSummary = &data
		return result, err
	case TypeIDUnbond:
		data, err := NewByteHashFromBuffer(buffer)
		result.Unbond = &AccountHash{Hash: data}
		return result, err
	case TypeIDChainspecRegistry:
		data, err := NewByteHashFromBuffer(buffer)
		result.ChainspecRegistry = &data
		return result, err
	case TypeIDChecksumRegistry:
		data, err := NewByteHashFromBuffer(buffer)
		result.ChecksumRegistry = &data
		return result, err
	default:
		return result, errors.New("type is not found")
	}
}
