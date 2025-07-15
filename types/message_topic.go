package types

import "github.com/make-software/casper-go-sdk/v2/types/key"

// MessageTopicSummary Summary of a message topic that will be stored in global state.
type MessageTopicSummary struct {
	// Number of messages in this topic.
	MessageCount uint32 `json:"message_count"`
	// Block timestamp in which these messages were emitted.
	BlockTime uint64 `json:"blocktime"`
}

// MessageChecksum Variant that stores a message digest.
type MessageChecksum string

type MessageTopic struct {
	// The name of the topic on which the message was emitted on.
	TopicName string `json:"topic_name"`
	// The hash of the name of the topic.
	TopicNameHash key.Hash `json:"topic_name_hash"`
}

// MessagePayload The payload of the message.
type MessagePayload struct {
	String *string `json:"String"`
	Bytes  *string `json:"Bytes"`
}

// Message that was emitted by an addressable entity during execution.
type Message struct {
	// The payload of the message.
	Message MessagePayload `json:"message"`
	// The name of the topic on which the message was emitted on.
	TopicName string `json:"topic_name"`
	// The hash of the name of the topic.
	TopicNameHash key.Hash `json:"topic_name_hash"`
	// The identity of the entity that produced the message.
	HashAddr key.Hash `json:"hash_addr"`
	// Message index in the block.
	BlockIndex uint64 `json:"block_index"`
	// Message index in the topic.
	TopicIndex uint32 `json:"topic_index"`
}
