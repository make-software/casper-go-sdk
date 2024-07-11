package key

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
)

var (
	ErrInvalidBidAddrTag             = errors.New("invalid BidAddrTag")
	ErrUnexpectedBidAddrTagInBidAddr = errors.New("unexpected BidAddrTag in BidAddr")
	ErrInvalidBidAddrFormat          = errors.New("invalid BidAddr format")
)

type BidAddrTag uint8

const (
	// Unified BidAddr for legacy unified bid.
	Unified BidAddrTag = iota
	// Validator BidAddr for validator bid.
	Validator
	// Delegator BidAddr for delegator bid.
	Delegator
	// Credit BidAddr for auction credit.
	Credit = 4
)

func NewBidAddrTagFromByte(tag uint8) (BidAddrTag, error) {
	addrTag := BidAddrTag(tag)
	if addrTag != Unified && addrTag != Validator && addrTag != Delegator && addrTag != Credit {
		return 0, ErrInvalidBidAddrTag
	}

	return addrTag, nil
}

const (
	// UnifiedOrValidatorAddrLen BidAddrTag(1) + Hash(32)
	UnifiedOrValidatorAddrLen = 33
	// CreditAddrLen BidAddrTag(1) + Hash(32) + EraId(8)
	CreditAddrLen = 41
	// DelegatorAddrLen BidAddrTag(1) + Hash(32) + Hash(32)
	DelegatorAddrLen = 65
)

// BidAddr  Bid Address
type BidAddr struct {
	Unified   *Hash
	Validator *Hash
	Delegator *struct {
		Validator Hash
		Delegator Hash
	}

	Credit *struct {
		Validator Hash
		EraId     uint64
	}
}

func NewBidAddr(source string) (BidAddr, error) {
	hexBytes, err := hex.DecodeString(source)
	if err != nil {
		return BidAddr{}, err
	}

	if len(source) < UnifiedOrValidatorAddrLen {
		return BidAddr{}, ErrInvalidBidAddrFormat
	}

	bitAddrTag, err := NewBidAddrTagFromByte(hexBytes[0])
	if err != nil {
		return BidAddr{}, err
	}

	if len(hexBytes) == UnifiedOrValidatorAddrLen {
		hash, err := NewHashFromBytes(hexBytes[1:])
		if err != nil {
			return BidAddr{}, err
		}
		switch bitAddrTag {
		case Unified:
			return BidAddr{Unified: &hash}, nil
		case Validator:
			return BidAddr{Validator: &hash}, nil
		default:
			return BidAddr{}, ErrUnexpectedBidAddrTagInBidAddr
		}
	}

	validatorHash, err := NewHashFromBytes(hexBytes[1:34])
	if err != nil {
		return BidAddr{}, err
	}

	if len(hexBytes) == CreditAddrLen {
		var eraID uint64
		if err := binary.Read(bytes.NewReader(hexBytes[33:]), binary.LittleEndian, &eraID); err != nil {
			return BidAddr{}, err
		}

		return BidAddr{Credit: &struct {
			Validator Hash
			EraId     uint64
		}{Validator: validatorHash, EraId: eraID}}, nil
	}

	delegatorHash, err := NewHashFromBytes(hexBytes[33:])
	if err != nil {
		return BidAddr{}, err
	}

	return newDelegatorBidAddr(validatorHash, delegatorHash), nil
}

func NewBidAddrFromBuffer(buf *bytes.Buffer) (BidAddr, error) {
	if buf.Len() < UnifiedOrValidatorAddrLen {
		return BidAddr{}, ErrInvalidBidAddrFormat
	}

	tag, err := buf.ReadByte()
	if err != nil {
		return BidAddr{}, err
	}

	bitAddrTag, err := NewBidAddrTagFromByte(tag)
	if err != nil {
		return BidAddr{}, err
	}

	if bitAddrTag == Unified {
		hash, err := NewByteHashFromBuffer(buf)
		if err != nil {
			return BidAddr{}, err
		}
		return BidAddr{
			Unified: &hash,
		}, nil
	}

	if bitAddrTag == Validator {
		hash, err := NewByteHashFromBuffer(buf)
		if err != nil {
			return BidAddr{}, err
		}
		return BidAddr{
			Validator: &hash,
		}, nil
	}

	validatorHash, err := NewByteHashFromBuffer(buf)
	if err != nil {
		return BidAddr{}, err
	}

	delegatorHash, err := NewByteHashFromBuffer(buf)
	if err != nil {
		return BidAddr{}, err
	}
	return newDelegatorBidAddr(validatorHash, delegatorHash), nil
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
		res = append(res, byte(Unified))
		return append(res, h.Unified.Bytes()...)
	case h.Validator != nil:
		res := make([]byte, 0, UnifiedOrValidatorAddrLen)
		res = append(res, byte(Validator))
		return append(res, h.Validator.Bytes()...)
	case h.Delegator != nil:
		res := make([]byte, 0, DelegatorAddrLen)
		res = append(res, byte(Delegator))
		res = append(res, h.Delegator.Validator.Bytes()...)
		return append(res, h.Delegator.Delegator.Bytes()...)
	case h.Credit != nil:
		res := make([]byte, 0, CreditAddrLen)
		res = append(res, byte(Credit))
		res = append(res, h.Credit.Validator.Bytes()...)

		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, h.Credit.EraId)

		return append(res, buf.Bytes()...)
	default:
		panic("Unexpected BidAddr type")
	}
}

func newDelegatorBidAddr(validatorHash, delegatorHash Hash) BidAddr {
	delegator := struct {
		Validator Hash
		Delegator Hash
	}{
		Validator: validatorHash,
		Delegator: delegatorHash,
	}
	return BidAddr{
		Delegator: &delegator,
	}
}
