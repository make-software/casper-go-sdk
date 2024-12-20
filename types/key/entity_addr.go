package key

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
)

var (
	ErrInvalidEntityAddrFormat = errors.New("invalid EntityAddr format")
	ErrInvalidEntityKind       = errors.New("invalid EntityKind")
)

type EntityKind uint8

const (
	// SystemKind `EntityKind::System` variant
	SystemKind EntityKind = iota
	// AccountKind `EntityKind::Account` variant.
	AccountKind
	// SmartContractKind `EntityKind::Package` variant.
	SmartContractKind
)

func NewEntityKindFromByte(tag uint8) (EntityKind, error) {
	kindTag := EntityKind(tag)
	if kindTag != SystemKind && kindTag != AccountKind && kindTag != SmartContractKind {
		return 0, ErrInvalidEntityKind
	}

	return kindTag, nil
}

const (
	SystemKindPrefix        = "system-"
	AccountKindNamePrefix   = "account-"
	SmartContractKindPrefix = "contract-"
)

type EntityAddr struct {
	System        *Hash
	Account       *Hash
	SmartContract *Hash
}

func (h *EntityAddr) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	tmp, err := NewEntityAddr(s)
	if err != nil {
		return err
	}
	*h = tmp
	return nil
}

func (h EntityAddr) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.ToPrefixedString())
}

func (h EntityAddr) ToPrefixedString() string {
	switch {
	case h.System != nil:
		return PrefixNameAddressableEntity + SystemKindPrefix + h.System.ToHex()
	case h.Account != nil:
		return PrefixNameAddressableEntity + AccountKindNamePrefix + h.Account.ToHex()
	case h.SmartContract != nil:
		return PrefixNameAddressableEntity + SmartContractKindPrefix + h.SmartContract.ToHex()
	}
	return ""
}

func NewEntityAddr(source string) (EntityAddr, error) {
	source = strings.TrimPrefix(source, PrefixNameAddressableEntity)

	if strings.HasPrefix(source, SystemKindPrefix) {
		hash, err := NewHash(strings.TrimPrefix(source, SystemKindPrefix))
		return EntityAddr{System: &hash}, err
	} else if strings.HasPrefix(source, AccountKindNamePrefix) {
		hash, err := NewHash(strings.TrimPrefix(source, AccountKindNamePrefix))
		return EntityAddr{Account: &hash}, err
	} else if strings.HasPrefix(source, SmartContractKindPrefix) {
		hash, err := NewHash(strings.TrimPrefix(source, SmartContractKindPrefix))
		return EntityAddr{SmartContract: &hash}, err
	}

	return EntityAddr{}, ErrInvalidEntityAddrFormat
}

func NewEntityAddrFromBuffer(buf *bytes.Buffer) (EntityAddr, error) {
	tag, err := buf.ReadByte()
	if err != nil {
		return EntityAddr{}, err
	}

	entityKindTag, err := NewEntityKindFromByte(tag)
	if err != nil {
		return EntityAddr{}, err
	}

	hash, err := NewByteHashFromBuffer(buf)
	if err != nil {
		return EntityAddr{}, err
	}

	switch entityKindTag {
	case SystemKind:
		return EntityAddr{System: &hash}, nil
	case AccountKind:
		return EntityAddr{Account: &hash}, nil
	case SmartContractKind:
		return EntityAddr{SmartContract: &hash}, nil
	}
	return EntityAddr{}, ErrInvalidEntityAddrFormat
}

func (h EntityAddr) Bytes() []byte {
	switch {
	case h.System != nil:
		res := make([]byte, 0, ByteHashLen)
		res = append(res, byte(SystemKind))
		return append(res, h.System.Bytes()...)
	case h.Account != nil:
		res := make([]byte, 0, ByteHashLen)
		res = append(res, byte(AccountKind))
		return append(res, h.Account.Bytes()...)
	case h.SmartContract != nil:
		res := make([]byte, 0, ByteHashLen)
		res = append(res, byte(SmartContractKind))
		return append(res, h.SmartContract.Bytes()...)
	default:
		panic("Unexpected EntityAddr type")
	}
}
