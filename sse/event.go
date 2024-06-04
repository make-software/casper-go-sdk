package sse

import (
	"encoding/json"
	"time"

	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/key"
	"github.com/make-software/casper-go-sdk/types/keypair"
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
	ShutdownType
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
	ShutdownType:             "Shutdown",
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
		BlockHash string        `json:"block_hash"`
		Block     types.BlockV1 `json:"block"`
	}
	BlockAddedEvent struct {
		BlockAdded BlockAdded `json:"BlockAdded"`
	}
)

type (
	DeployProcessedPayload struct {
		DeployHash      key.Hash                `json:"deploy_hash"`
		Account         string                  `json:"account"`
		Timestamp       time.Time               `json:"timestamp"`
		TTL             string                  `json:"ttl"`
		BlockHash       key.Hash                `json:"block_hash"`
		ExecutionResult types.ExecutionResultV1 `json:"execution_result"`
	}
	DeployProcessedEvent struct {
		DeployProcessed DeployProcessedPayload `json:"DeployProcessed"`
	}
	DeployAcceptedEvent struct {
		DeployAccepted types.Deploy `json:"DeployAccepted"`
	}
	DeployExpiredPayload struct {
		DeployHash key.Hash `json:"deploy_hash"`
	}
	DeployExpiredEvent struct {
		DeployExpired DeployExpiredPayload `json:"DeployExpired"`
	}
)

type (
	FinalitySignaturePayload struct {
		BlockHash key.Hash          `json:"block_hash"`
		EraID     uint64            `json:"era_id"`
		Signature types.HexBytes    `json:"signature"`
		PublicKey keypair.PublicKey `json:"public_key"`
	}
	FinalitySignatureEvent struct {
		FinalitySignature FinalitySignaturePayload `json:"FinalitySignature"`
	}
)

type (
	FaultPayload struct {
		EraID     uint64            `json:"era_id"`
		PublicKey keypair.PublicKey `json:"public_key"`
		Timestamp types.Timestamp   `json:"timestamp"`
	}
	FaultEvent struct {
		Fault FaultPayload `json:"Fault"`
	}
)

type (
	StepPayload struct {
		EraID           uint64       `json:"era_id"`
		ExecutionEffect types.Effect `json:"execution_effect"`
		// Todo: not sure, didn't found example to test
		Operations *[]types.Operation `json:"operations,omitempty"`
		// Todo: not sure, didn't found example to test
		Transform *types.TransformKey `json:"transform,omitempty"`
	}
	StepEvent struct {
		Step StepPayload `json:"step"`
	}
)

func (d *RawEvent) ParseAsAPIVersionEvent() (APIVersionEvent, error) {
	return ParseEvent[APIVersionEvent](d.Data)
}

func (d *RawEvent) ParseAsDeployProcessedEvent() (DeployProcessedEvent, error) {
	return ParseEvent[DeployProcessedEvent](d.Data)
}

func (d *RawEvent) ParseAsBlockAddedEvent() (BlockAddedEvent, error) {
	return ParseEvent[BlockAddedEvent](d.Data)
}

func (d *RawEvent) ParseAsDeployAcceptedEvent() (DeployAcceptedEvent, error) {
	return ParseEvent[DeployAcceptedEvent](d.Data)
}

func (d *RawEvent) ParseAsFinalitySignatureEvent() (FinalitySignatureEvent, error) {
	return ParseEvent[FinalitySignatureEvent](d.Data)
}

func (d *RawEvent) ParseAsDeployExpiredEvent() (DeployExpiredEvent, error) {
	return ParseEvent[DeployExpiredEvent](d.Data)
}

func (d *RawEvent) ParseAsFaultEvent() (FaultEvent, error) {
	return ParseEvent[FaultEvent](d.Data)
}

func (d *RawEvent) ParseAsStepEvent() (StepEvent, error) {
	return ParseEvent[StepEvent](d.Data)
}

func ParseEvent[T interface{}](data []byte) (T, error) {
	var res T
	if err := json.Unmarshal(data, &res); err != nil {
		return res, err
	}
	return res, nil
}
