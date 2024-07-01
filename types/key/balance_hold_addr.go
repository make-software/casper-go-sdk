package key

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrBalanceHoldAddrTag    = errors.New("invalid BalanceHoldAddrTag")
	ErrBalanceHoldAddrFormat = errors.New("invalid BalanceHoldAddr format")
)

type BalanceHoldAddrTag uint8

const (
	// Gas Tag for gas variant
	Gas BalanceHoldAddrTag = iota
	// Processing Tag for processing variant
	Processing
)

const BlockTypeBytesLen = 8

func NewBalanceHoldAddrTagFromByte(tag uint8) (BalanceHoldAddrTag, error) {
	addrTag := BalanceHoldAddrTag(tag)
	if addrTag != Gas && addrTag != Processing {
		return 0, ErrBalanceHoldAddrTag
	}

	return addrTag, nil
}

type (
	URefAddr = [32]byte
)

type Hold struct {
	// The address of the purse this hold is on.
	PurseAddr URefAddr
	// The block time this hold was placed.
	BlockTime time.Time
}

// BalanceHoldAddr  Balance hold address
type BalanceHoldAddr struct {
	Gas        *Hold
	Processing *Hold
}

func NewBalanceHoldAddr(source string) (BalanceHoldAddr, error) {
	decoded, err := hex.DecodeString(source)
	if err != nil {
		return BalanceHoldAddr{}, err
	}

	hold, err := NewBalanceHoldAddrFromBuffer(bytes.NewBuffer(decoded))
	if err != nil {
		return BalanceHoldAddr{}, err
	}
	return hold, nil
}

func NewBalanceHoldAddrFromBuffer(buf *bytes.Buffer) (BalanceHoldAddr, error) {
	if buf.Len() < ByteHashLen+BlockTypeBytesLen {
		return BalanceHoldAddr{}, ErrBalanceHoldAddrFormat
	}

	tag, err := buf.ReadByte()
	if err != nil {
		return BalanceHoldAddr{}, err
	}

	balanceHoldAddrTag, err := NewBalanceHoldAddrTagFromByte(tag)
	if err != nil {
		return BalanceHoldAddr{}, err
	}

	var addr = buf.Next(ByteHashLen)
	var purseAddr [32]byte
	for i := 0; i < 32; i++ {
		purseAddr[i] = addr[i]
	}

	timestamp := int64(binary.LittleEndian.Uint64(buf.Next(BlockTypeBytesLen)))
	blockTime := time.UnixMilli(timestamp)

	switch balanceHoldAddrTag {
	case Gas:
		return BalanceHoldAddr{
			Gas: &Hold{
				PurseAddr: purseAddr,
				BlockTime: blockTime,
			},
		}, nil
	case Processing:
		return BalanceHoldAddr{
			Gas: &Hold{
				PurseAddr: purseAddr,
				BlockTime: blockTime,
			},
		}, nil
	default:
		panic("Unexpected BalanceHoldAddr type")
	}
	return BalanceHoldAddr{}, nil
}

func (h *BalanceHoldAddr) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val, err := NewBalanceHoldAddr(s)
	if err != nil {
		return err
	}
	*h = val
	return nil
}

func (h BalanceHoldAddr) MarshalJSON() ([]byte, error) {
	return []byte(h.ToPrefixedString()), nil
}

func (h BalanceHoldAddr) ToPrefixedString() string {
	return PrefixNameBalanceHold + hex.EncodeToString(h.Bytes())
}

func (h BalanceHoldAddr) Bytes() []byte {
	hold := h.Gas
	holdType := Gas
	if h.Processing != nil {
		holdType = Processing
		hold = h.Processing
	}

	res := make([]byte, 0, ByteHashLen+BlockTypeBytesLen)
	res = append(res, byte(holdType))
	res = append(res, hold.PurseAddr[:]...)

	blockTime := make([]byte, BlockTypeBytesLen)
	binary.LittleEndian.PutUint64(blockTime, uint64(hold.BlockTime.UnixMilli()))
	res = append(res, blockTime...)
	return res
}
