package types

// MessageTopicSummary Summary of a message topic that will be stored in global state.
type MessageTopicSummary struct {
	// Number of messages in this topic.
	MessageCount uint32 `json:"message_count"`
	// Block timestamp in which these messages were emitted.
	BlockTime uint64 `json:"blocktime"`
}

// MessageChecksum Variant that stores a message digest.
type MessageChecksum string
