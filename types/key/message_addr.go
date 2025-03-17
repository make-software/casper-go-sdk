package key

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"strconv"
	"strings"
)

const TopicPrefix = "topic-"

// MessageAddr MessageTopicAddr
type MessageAddr struct {
	// The hash addr
	EntityAddr EntityAddr
	// The hash of the name of the message topic
	TopicNameHash Hash
	//  The message index.
	MessageIndex *uint32
}

func NewMessageAddr(source string) (MessageAddr, error) {
	var (
		messageAddr MessageAddr
		err         error
	)

	// If no `topic-` prefix, treat it as a message address with an index
	if !strings.HasPrefix(source, TopicPrefix) {
		rawId := source[strings.LastIndex(source, "-")+1:]
		source = source[:strings.LastIndex(source, "-")]
		idx, err := strconv.ParseUint(rawId, 10, 32)
		if err != nil {
			return MessageAddr{}, err
		}
		idx32 := uint32(idx)
		messageAddr.MessageIndex = &idx32
	} else {
		source = strings.TrimPrefix(source, TopicPrefix)
	}

	rawTopicNameHash := source[strings.LastIndex(source, "-")+1:]
	source = source[:strings.LastIndex(source, "-")]

	messageAddr.TopicNameHash, err = NewHash(rawTopicNameHash)
	if err != nil {
		return MessageAddr{}, err
	}

	entityAddr, err := NewEntityAddr(source)
	if err != nil {
		return MessageAddr{}, err
	}

	messageAddr.EntityAddr = entityAddr
	return messageAddr, nil
}

func (h *MessageAddr) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	val, err := NewMessageAddr(s)
	if err != nil {
		return err
	}
	*h = val
	return nil
}

func (h MessageAddr) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.ToPrefixedString())
}

func (h MessageAddr) ToPrefixedString() string {
	res := PrefixNameMessage
	if h.MessageIndex == nil {
		res += TopicPrefix
	}
	res += h.EntityAddr.ToPrefixedString()
	res += "-" + h.TopicNameHash.ToHex()
	if h.MessageIndex != nil {
		res += "-" + strconv.Itoa(int(*h.MessageIndex))
	}
	return res
}

func (h MessageAddr) Bytes() []byte {
	res := make([]byte, 0, ByteHashLen)
	res = append(res, h.EntityAddr.Bytes()...)
	res = append(res, h.TopicNameHash.Bytes()...)

	if h.MessageIndex != nil {
		binary.LittleEndian.PutUint32(res, *h.MessageIndex)
	}
	return res
}

func NewMessageAddrFromBuffer(buf *bytes.Buffer) (MessageAddr, error) {
	entityAddr, err := NewEntityAddrFromBuffer(buf)
	if err != nil {
		return MessageAddr{}, err
	}

	topicNameHash, err := NewByteHashFromBuffer(buf)
	if err != nil {
		return MessageAddr{}, err
	}

	var msgIdx *uint32
	if buf.Len() > 0 {
		idx := binary.LittleEndian.Uint32(buf.Bytes())
		msgIdx = &idx
	}

	return MessageAddr{
		EntityAddr:    entityAddr,
		TopicNameHash: topicNameHash,
		MessageIndex:  msgIdx,
	}, nil
}
