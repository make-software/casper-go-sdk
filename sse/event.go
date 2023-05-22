package sse

import (
	"encoding/json"
	"time"

	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/key"
)

const (
	APIVersionEventType EventType = iota + 1
	BlockAddedEventType
	DeployProcessedEventType
	DeployAcceptedEventType
	DeployExpiredEventType
	EventIDEventType
	FinalitySignatureType
	StepEventType
	FaultEventType
)

var AllEventsNames = map[EventType]string{
	APIVersionEventType:      "ApiVersion",
	BlockAddedEventType:      "BlockAdded",
	DeployProcessedEventType: "DeployProcessed",
	DeployAcceptedEventType:  "DeployAccepted",
	DeployExpiredEventType:   "DeployExpired",
	StepEventType:            "Step",
	FaultEventType:           "Fault",
	FinalitySignatureType:    "FinalitySignature",
}

type (
	EventType = int
	EventData = json.RawMessage
)

type RawEvent struct {
	EventType EventType
	Data      EventData
	EventID   uint64
}

type APIVersionEvent struct {
	APIVersion string `json:"ApiVersion"`
}

// BlockAddedEvent definition
type (
	BlockAdded struct {
		BlockHash string      `json:"block_hash"`
		Block     types.Block `json:"block"`
	}
	BlockAddedEvent struct {
		BlockAdded BlockAdded `json:"BlockAdded"`
	}
)

type (
	DeployProcessed struct {
		DeployHash      key.Hash                    `json:"deploy_hash"`
		Account         string                      `json:"account"`
		Timestamp       time.Time                   `json:"timestamp"`
		TTL             string                      `json:"ttl"`
		BlockHash       string                      `json:"block_hash"`
		ExecutionResult types.ExecutionResultStatus `json:"execution_result"`
	}
	DeployProcessedEvent struct {
		DeployProcessed DeployProcessed `json:"DeployProcessed"`
	}
	DeployAcceptedEvent struct {
		DeployAccepted types.Deploy `json:"DeployAccepted"`
	}
)

func (d *RawEvent) ParseAsAPIVersionEvent() (APIVersionEvent, error) {
	res := APIVersionEvent{}
	if err := json.Unmarshal(d.Data, &res); err != nil {
		return APIVersionEvent{}, err
	}
	return res, nil
}

func (d *RawEvent) ParseAsDeployProcessedEvent() (DeployProcessedEvent, error) {
	res := DeployProcessedEvent{}
	if err := json.Unmarshal(d.Data, &res); err != nil {
		return DeployProcessedEvent{}, err
	}
	return res, nil
}

func (d *RawEvent) ParseAsBlockAddedEvent() (BlockAddedEvent, error) {
	res := BlockAddedEvent{}
	if err := json.Unmarshal(d.Data, &res); err != nil {
		return BlockAddedEvent{}, err
	}
	return res, nil
}
