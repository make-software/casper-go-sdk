package key

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
)

var (
	ErrInvalidBlockGlobalAddrFormat = errors.New("invalid BlockGlobalAddr format")
	ErrInvalidBlockGlobalAddrTag    = errors.New("invalid BlockGlobalAddrTag")
)

type BlockGlobalAddrTag uint8

const (
	// BlockTime Tag for block time variant
	BlockTime BlockGlobalAddrTag = iota
	// MessageCount Tag for message count variant
	MessageCount
)

func NewBlockGlobalAddrTagFromByte(tag uint8) (BlockGlobalAddrTag, error) {
	addrTag := BlockGlobalAddrTag(tag)
	if addrTag != BlockTime && addrTag != MessageCount {
		return 0, ErrInvalidBlockGlobalAddrTag
	}

	return addrTag, nil
}

const (
	BlockTimePrefix    = "time-"
	MessageCountPrefix = "message-count-"
)

type BlockGlobalAddr struct {
	BlockTime    *struct{}
	MessageCount *struct{}
}

func (h *BlockGlobalAddr) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	tmp, err := NewBlockGlobalAddr(s)
	if err != nil {
		return err
	}
	*h = tmp
	return nil
}

func (h BlockGlobalAddr) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.ToPrefixedString())
}

func (h BlockGlobalAddr) ToPrefixedString() string {
	res := PrefixNameBlockGlobal
	if h.BlockTime != nil {
		res += BlockTimePrefix
	} else {
		res += MessageCountPrefix
	}
	emptyHash, _ := NewHashFromBytes(make([]byte, ByteHashLen))
	return res + emptyHash.ToHex()
}

func NewBlockGlobalAddr(source string) (BlockGlobalAddr, error) {
	if strings.HasPrefix(source, BlockTimePrefix) {
		return BlockGlobalAddr{BlockTime: &struct{}{}}, nil
	} else if strings.HasPrefix(source, MessageCountPrefix) {
		return BlockGlobalAddr{MessageCount: &struct{}{}}, nil
	}

	return BlockGlobalAddr{}, ErrInvalidBlockGlobalAddrFormat
}

func (h BlockGlobalAddr) Bytes() []byte {
	tag := BlockTime
	if h.MessageCount != nil {
		tag = MessageCount
	}
	return []byte{byte(tag)}
}

func NewBlockGlobalAddrFrom(buf *bytes.Buffer) (BlockGlobalAddr, error) {
	tagByte, err := buf.ReadByte()
	if err != nil {
		return BlockGlobalAddr{}, err
	}

	tag, err := NewBlockGlobalAddrTagFromByte(tagByte)
	if err != nil {
		return BlockGlobalAddr{}, err
	}
	switch tag {
	case BlockTime:
		return BlockGlobalAddr{BlockTime: &struct{}{}}, nil
	case MessageCount:
		return BlockGlobalAddr{MessageCount: &struct{}{}}, nil
	default:
		panic("Unexpected BlockGlobalAddr type")
	}
}
