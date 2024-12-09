package key

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
)

var (
	ErrInvalidBidAddrTag             = errors.New("invalid BidAddrTag")
	ErrUnexpectedBidAddrTagInBidAddr = errors.New("unexpected BidAddrTag in BidAddr")
	ErrInvalidBidAddrFormat          = errors.New("invalid BidAddr format")
)

type BidAddrTag uint8

const (
	UnifiedTag BidAddrTag = iota
	ValidatorTag
	DelegatedAccountTag
	DelegatedPurseTag
	CreditTag
	ReservedDelegationAccountTag
	ReservedDelegationPurseTag
	UnbondAccountTag
	UnbondPurseTag
)

var allowedBidAddrTags = map[BidAddrTag]struct{}{
	UnifiedTag:                   {},
	ValidatorTag:                 {},
	DelegatedAccountTag:          {},
	DelegatedPurseTag:            {},
	CreditTag:                    {},
	ReservedDelegationAccountTag: {},
	ReservedDelegationPurseTag:   {},
	UnbondAccountTag:             {},
	UnbondPurseTag:               {},
}

func NewBidAddrTagFromByte(tag uint8) (BidAddrTag, error) {
	addrTag := BidAddrTag(tag)

	if _, ok := allowedBidAddrTags[addrTag]; !ok {
		return 0, ErrInvalidBidAddrTag
	}

	return addrTag, nil
}

const (
	// UnifiedOrValidatorAddrLen BidAddrTag(1) + Hash(32)
	UnifiedOrValidatorAddrLen = 33
	// CreditAddrLen BidAddrTag(1) + Hash(32) + EraId(8)
	CreditAddrLen = 41
	// ValidatorHashDelegatorHashAddrLen BidAddrTag(1) + Hash(32) + Hash(32)
	ValidatorHashDelegatorHashAddrLen = 65
	// ValidatorHashDelegatorUrefAddrLen BidAddrTag(1) + Hash(32) + URef(32)
	ValidatorHashDelegatorUrefAddrLen = 65
)

// BidAddr  Bid Address
type BidAddr struct {
	Unified          *Hash
	Validator        *Hash
	DelegatedAccount *struct {
		Validator Hash
		Delegator Hash
	}
	DelegatedPurse *struct {
		Validator Hash
		Delegator URef
	}

	Credit *struct {
		Validator Hash
		EraId     uint64
	}

	ReservedDelegationAccount *struct {
		Validator Hash
		Delegator Hash
	}

	ReservedDelegationPurse *struct {
		Validator Hash
		Delegator URef
	}

	UnbondAccount *struct {
		Validator Hash
		Delegator Hash
	}

	UnbondPurse *struct {
		Validator Hash
		Delegator URef
	}
}

func NewBidAddr(source string) (BidAddr, error) {
	hexBytes, err := hex.DecodeString(source)
	if err != nil {
		return BidAddr{}, err
	}
	buf := bytes.NewBuffer(hexBytes)
	return NewBidAddrFromBuffer(buf)
}

func NewBidAddrFromBuffer(buf *bytes.Buffer) (BidAddr, error) {
	tag, err := buf.ReadByte()
	if err != nil {
		return BidAddr{}, err
	}

	bitAddrTag, err := NewBidAddrTagFromByte(tag)
	if err != nil {
		return BidAddr{}, err
	}

	switch bitAddrTag {
	case UnifiedTag:
		hash, err := NewHashFromBytes(buf.Next(ByteHashLen))
		if err != nil {
			return BidAddr{}, err
		}
		return BidAddr{Unified: &hash}, nil
	case ValidatorTag:
		hash, err := NewHashFromBytes(buf.Next(ByteHashLen))
		if err != nil {
			return BidAddr{}, err
		}
		return BidAddr{Validator: &hash}, nil
	case DelegatedAccountTag:
		validator, delegator, err := readValidatorDelegatorHash(buf)
		if err != nil {
			return BidAddr{}, err
		}
		return BidAddr{DelegatedAccount: &struct {
			Validator Hash
			Delegator Hash
		}{Validator: validator, Delegator: delegator}}, nil
	case DelegatedPurseTag:
		validator, delegator, err := readValidatorHashDelegatorUref(buf)
		if err != nil {
			return BidAddr{}, err
		}
		return BidAddr{DelegatedPurse: &struct {
			Validator Hash
			Delegator URef
		}{Validator: validator, Delegator: delegator}}, nil
	case CreditTag:
		validator, eraID, err := readValidatorHashEraID(buf)
		if err != nil {
			return BidAddr{}, err
		}
		return BidAddr{Credit: &struct {
			Validator Hash
			EraId     uint64
		}{Validator: validator, EraId: eraID}}, nil
	case ReservedDelegationAccountTag:
		validator, delegator, err := readValidatorDelegatorHash(buf)
		if err != nil {
			return BidAddr{}, err
		}
		return BidAddr{ReservedDelegationAccount: &struct {
			Validator Hash
			Delegator Hash
		}{Validator: validator, Delegator: delegator}}, nil
	case ReservedDelegationPurseTag:
		validator, delegator, err := readValidatorHashDelegatorUref(buf)
		if err != nil {
			return BidAddr{}, err
		}
		return BidAddr{ReservedDelegationPurse: &struct {
			Validator Hash
			Delegator URef
		}{Validator: validator, Delegator: delegator}}, nil
	case UnbondAccountTag:
		validator, delegator, err := readValidatorDelegatorHash(buf)
		if err != nil {
			return BidAddr{}, err
		}
		return BidAddr{UnbondAccount: &struct {
			Validator Hash
			Delegator Hash
		}{Validator: validator, Delegator: delegator}}, nil
	case UnbondPurseTag:
		validator, delegator, err := readValidatorHashDelegatorUref(buf)
		if err != nil {
			return BidAddr{}, err
		}
		return BidAddr{UnbondPurse: &struct {
			Validator Hash
			Delegator URef
		}{Validator: validator, Delegator: delegator}}, nil
	default:
		return BidAddr{}, ErrUnexpectedBidAddrTagInBidAddr
	}
}

func (h *BidAddr) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val, err := NewBidAddr(s)
	if err != nil {
		return err
	}
	*h = val
	return nil
}

func (h BidAddr) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.ToPrefixedString())
}

func (h BidAddr) ToPrefixedString() string {
	return PrefixNameBidAddr + hex.EncodeToString(h.Bytes())
}

func (h BidAddr) Bytes() []byte {
	switch {
	case h.Unified != nil:
		res := make([]byte, 0, UnifiedOrValidatorAddrLen)
		res = append(res, byte(UnifiedTag))
		return append(res, h.Unified.Bytes()...)
	case h.Validator != nil:
		res := make([]byte, 0, UnifiedOrValidatorAddrLen)
		res = append(res, byte(ValidatorTag))
		return append(res, h.Validator.Bytes()...)
	case h.DelegatedAccount != nil:
		res := make([]byte, 0, ValidatorHashDelegatorHashAddrLen)
		res = append(res, byte(DelegatedAccountTag))
		res = append(res, h.DelegatedAccount.Validator.Bytes()...)
		return append(res, h.DelegatedAccount.Delegator.Bytes()...)
	case h.DelegatedPurse != nil:
		res := make([]byte, 0, ValidatorHashDelegatorUrefAddrLen)
		res = append(res, byte(DelegatedPurseTag))
		res = append(res, h.DelegatedPurse.Validator.Bytes()...)
		return append(res, h.DelegatedPurse.Delegator.DataBytes()...)
	case h.Credit != nil:
		res := make([]byte, 0, CreditAddrLen)
		res = append(res, byte(CreditTag))
		res = append(res, h.Credit.Validator.Bytes()...)
		data := make([]byte, cltype.Int64ByteSize)
		binary.LittleEndian.PutUint64(data, h.Credit.EraId)
		return append(res, data...)
	case h.ReservedDelegationAccount != nil:
		res := make([]byte, 0, ValidatorHashDelegatorHashAddrLen)
		res = append(res, byte(ReservedDelegationAccountTag))
		res = append(res, h.ReservedDelegationAccount.Validator.Bytes()...)
		return append(res, h.ReservedDelegationAccount.Delegator.Bytes()...)
	case h.ReservedDelegationPurse != nil:
		res := make([]byte, 0, ValidatorHashDelegatorUrefAddrLen)
		res = append(res, byte(ReservedDelegationPurseTag))
		res = append(res, h.ReservedDelegationPurse.Validator.Bytes()...)
		return append(res, h.ReservedDelegationPurse.Delegator.DataBytes()...)
	case h.UnbondAccount != nil:
		res := make([]byte, 0, ValidatorHashDelegatorHashAddrLen)
		res = append(res, byte(UnbondAccountTag))
		res = append(res, h.UnbondAccount.Validator.Bytes()...)
		return append(res, h.UnbondAccount.Delegator.Bytes()...)
	case h.UnbondPurse != nil:
		res := make([]byte, 0, ValidatorHashDelegatorUrefAddrLen)
		res = append(res, byte(UnbondPurseTag))
		res = append(res, h.UnbondPurse.Validator.Bytes()...)
		return append(res, h.UnbondPurse.Delegator.DataBytes()...)
	default:
		panic("Unexpected BidAddr type")
	}
}

func readValidatorDelegatorHash(buf *bytes.Buffer) (Hash, Hash, error) {
	if buf.Len() < ByteHashLen {
		return Hash{}, Hash{}, ErrInvalidBidAddrFormat
	}

	validator := make([]byte, ByteHashLen)
	copy(validator[:], buf.Next(ByteHashLen))

	if buf.Len() < ByteHashLen {
		return Hash{}, Hash{}, ErrInvalidBidAddrFormat
	}

	delegator := make([]byte, ByteHashLen)
	copy(delegator[:], buf.Next(ByteHashLen))

	return Hash(validator), Hash(delegator), nil
}

func readValidatorHashDelegatorUref(buf *bytes.Buffer) (Hash, URef, error) {
	if buf.Len() < ByteHashLen {
		return Hash{}, URef{}, ErrInvalidBidAddrFormat
	}

	validator := make([]byte, ByteHashLen)
	copy(validator[:], buf.Next(ByteHashLen))

	if buf.Len() < ByteHashLen {
		return Hash{}, URef{}, ErrInvalidBidAddrFormat
	}

	uref := make([]byte, ByteHashLen)
	copy(uref[:], buf.Next(ByteHashLen))

	urefRes, err := NewURefFromBytes(append(uref, byte(07)))
	if err != nil {
		return Hash{}, URef{}, err
	}

	return Hash(validator), urefRes, nil
}

func readValidatorHashEraID(buf *bytes.Buffer) (Hash, uint64, error) {
	if buf.Len() < ByteHashLen {
		return Hash{}, 0, ErrInvalidBidAddrFormat
	}

	validator := make([]byte, ByteHashLen)
	copy(validator[:], buf.Next(ByteHashLen))

	if buf.Len() < 8 {
		return Hash{}, 0, ErrInvalidBidAddrFormat
	}

	var eraID uint64
	if err := binary.Read(bytes.NewReader(buf.Next(8)), binary.LittleEndian, &eraID); err != nil {
		return Hash{}, 0, ErrInvalidBidAddrFormat
	}

	return Hash(validator), eraID, nil
}
